package infra

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/aditya37/geofence-service/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose"
)

type MysqlConfigParam struct {
	Host              string
	Port              int
	Option            string
	Name              string
	User              string
	Password          string
	MaxConnection     int
	MaxIdleConnection int
}

// create singleton..
var mysqlSingleton sync.Once

// instance
var mysqlClientInstance *sql.DB = nil
var retErr error

func NewMysqlClient(param MysqlConfigParam) error {
	mysqlSingleton.Do(func() {
		connURL := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?multiStatements=true&parseTime=true&loc=%s",
			param.User,
			param.Password,
			param.Host,
			param.Port,
			param.Name,
			"Asia%2FJakarta",
		)
		log.Println(fmt.Sprintf(
			"MySQL Connection %s:%s@tcp(%s:%d)/%s",
			param.User,
			"********************",
			param.Host,
			param.Port,
			param.Name,
		))
		db, err := sql.Open("mysql", connURL)
		if err != nil {
			util.Logger().Error(err)
			retErr = err
			return
		}
		if param.MaxConnection > 0 {
			db.SetMaxOpenConns(param.MaxConnection)
		}
		if param.MaxIdleConnection > 0 {
			db.SetMaxIdleConns(param.MaxIdleConnection)
		}

		// migration
		if err := goose.SetDialect("mysql"); err != nil {
			util.Logger().Error(err)
			retErr = err
			return
		}
		if err := goose.Up(db, "migration"); err != nil {
			util.Logger().Error(err)
			retErr = err
			return
		}

		mysqlClientInstance = db

	})
	if retErr != nil {
		return retErr
	}

	return nil
}

// GetMysqlClientInstance.
func GetMysqlClientInstance() *sql.DB {
	return mysqlClientInstance
}
