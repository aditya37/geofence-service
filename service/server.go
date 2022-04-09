package service

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	delivemux "github.com/aditya37/geofence-service/delivery/mux"
	"github.com/aditya37/geofence-service/infra"
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
	// infra

	infra.NewRedisInstance(infra.RedisConfigParam{
		Port:     getenv.GetInt("REDIS_PORT", 6379),
		Host:     getenv.GetString("REDIS_HOST", ""),
		Database: getenv.GetInt("REDIS_DATABASE", 0),
		Password: getenv.GetString("REDIS_PASSWORD", ""),
	})
	redisClient := infra.GetRedisInstance()

	// if nil instance
	if redisClient == nil {
		// create new instance
		infra.NewRedisInstance(infra.RedisConfigParam{
			Port:     getenv.GetInt("REDIS_PORT", 6379),
			Host:     getenv.GetString("REDIS_HOST", "127.0.0.1"),
			Database: getenv.GetInt("REDIS_DATABASE", 0),
			Password: getenv.GetString("REDIS_PASSWORD", ""),
		})
		redisClient = infra.GetRedisInstance()
	}

	// repository

	// usecase
	eventStateCase, err := eventstate.NewEventStateUsecase()
	if err != nil {
		return nil, err
	}

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
