package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	delivemux "github.com/aditya37/geofence-service/delivery/mux"
	"github.com/aditya37/geofence-service/infra"
	geofencemanager "github.com/aditya37/geofence-service/repository/mysql/geofence"
	"github.com/aditya37/geofence-service/repository/pubsub"
	eventmanager "github.com/aditya37/geofence-service/repository/redis/event-manager"
	tile38Channel "github.com/aditya37/geofence-service/repository/tile38/channel-manager"
	"github.com/aditya37/geofence-service/usecase/eventstate"
	getenv "github.com/aditya37/get-env"

	// multiplex server (http,grpc)
	"github.com/soheilhy/cmux"
)

type (
	server struct {
		httpHandler *httpServer
		close       func()
	}
	Server interface {
		Run()
	}
)

func NewServer() (Server, error) {
	ctx := context.Background()
	// infra config
	configRedis := infra.RedisConfigParam{
		Port:     getenv.GetInt("REDIS_PORT", 6379),
		Host:     getenv.GetString("REDIS_HOST", ""),
		Database: getenv.GetInt("REDIS_DATABASE", 0),
		Password: getenv.GetString("REDIS_PASSWORD", ""),
	}
	configMysqlClient := infra.MysqlConfigParam{
		Host:     getenv.GetString("DB_HOST", "127.0.0.1"),
		Port:     getenv.GetInt("DB_PORT", 3306),
		Name:     getenv.GetString("DB_NAME", "db_geofencing"),
		User:     getenv.GetString("DB_USER", "root"),
		Password: getenv.GetString("DB_PASSWORD", "admin"),
	}
	configTile38Client := infra.Tile38ClientParam{
		Host: getenv.GetString("TILE38_HOST", "127.0.0.1"),
		Port: getenv.GetInt("TILE38_PORT", 9851),
	}

	// infra instance
	// redis
	infra.NewRedisInstance(configRedis)
	redisClient := infra.GetRedisInstance()
	// if redis instance nil
	if redisClient == nil {
		// create new instance
		infra.NewRedisInstance(configRedis)
		redisClient = infra.GetRedisInstance()
	}
	// MySQL
	if err := infra.NewMysqlClient(configMysqlClient); err != nil {
		return nil, err
	}
	mysqlClient := infra.GetMysqlClientInstance()
	if mysqlClient == nil {
		infra.NewMysqlClient(configMysqlClient)
		mysqlClient = infra.GetMysqlClientInstance()
	}
	// TIle38
	if err := infra.NewTile38Client(configTile38Client); err != nil {
		return nil, err
	}
	tile38Client := infra.GetTile38ClientInstance()
	if tile38Client == nil {
		infra.NewTile38Client(configTile38Client)
		tile38Client = infra.GetTile38ClientInstance()
	}

	// repository
	gpubsub, err := pubsub.NewGcpPubsub(
		ctx,
		getenv.GetString("GCP_PROJECT_ID", ""),
	)
	if err != nil {
		return nil, err
	}
	// cacheEventManager
	eventManager, err := eventmanager.NewCacheEventManager(redisClient)
	if err != nil {
		return nil, err
	}

	// geofenceManager
	geofenceManager, err := geofencemanager.NewGeofenceManager(mysqlClient)
	if err != nil {
		return nil, err
	}
	// tile38 channel Manager
	t38cChannelManager, err := tile38Channel.NewChannelManager(tile38Client)
	if err != nil {
		return nil, err
	}

	// usecase
	eventStateCase, err := eventstate.NewEventStateUsecase(
		gpubsub,
		eventManager,
		geofenceManager,
		t38cChannelManager,
	)
	if err != nil {
		return nil, err
	}

	// subscribe eventState...
	go eventStateCase.ConsumeEventState(
		ctx,
		getenv.GetString("GEOFENCE_EVENT_STATE_TOPIC", "geofence-event-state"),
		getenv.GetString("SERVICE_NAME", "geofence-service"),
	)

	// http/mux delivery
	eventStateDeliver := delivemux.NewEventStateDelivery(eventStateCase)

	// http handler
	handler, err := NewHttpServer(eventStateDeliver)
	if err != nil {
		return nil, err
	}

	// grcp handler

	return &server{
		httpHandler: handler,
		close: func() {
			gpubsub.Close()
			eventManager.Close()
			geofenceManager.Close()
			t38cChannelManager.Close()
		},
	}, nil
}

// server runner and listen
func (s *server) Run() {
	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
		errs <- fmt.Errorf("%s", <-c)
	}()
	go func() {
		errs <- s.listen(s.httpHandler.handler())
	}()
	defer s.close()
	log.Fatalf("Stop server with error detail: %v", <-errs)
}

// listen
func (s *server) listen(httpHandler http.Handler) error {
	tcpListen, err := net.Listen("tcp", fmt.Sprintf(":%s", getenv.GetString("SERVICE_PORT", "7778")))
	if err != nil {
		return err
	}
	// multiplex server
	m := cmux.New(tcpListen)
	log.Println(fmt.Sprintf(
		"%s run on Port %s",
		getenv.GetString("SERVICE_NAME", "geofence-service"),
		getenv.GetString("SERVICE_PORT", "7778"),
	),
	)
	// serve http1
	httpl := m.Match(cmux.HTTP1Fast())
	// server config
	httpConf := &http.Server{
		Handler: httpHandler,
	}
	go httpConf.Serve(httpl)
	return m.Serve()
}
