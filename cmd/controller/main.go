package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/deviceplane/deviceplane/pkg/controller/service"
	mysql_store "github.com/deviceplane/deviceplane/pkg/controller/store/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/segmentio/conf"
)

var version = "dev"
var name = "deviceplane-controller"

var config struct {
	Addr         string `conf:"addr"`
	MySQLPrimary string `conf:"mysql-primary"`
}

func init() {
	config.Addr = ":8080"
	config.MySQLPrimary = "user:pass@tcp(localhost:3313)/deviceplane?parseTime=true"
}

func main() {
	conf.Load(&config)

	db, err := tryConnect(config.MySQLPrimary)
	if err != nil {
		panic(err)
	}

	store := mysql_store.NewStore(db)

	svc := service.NewService(store, store, store, store, store)
	server := &http.Server{
		Addr:    config.Addr,
		Handler: svc,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func tryConnect(uri string) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		if db, err = sql.Open("mysql", uri); err != nil {
			time.Sleep(time.Second)
			continue
		}

		if err = db.Ping(); err != nil {
			db.Close()
			time.Sleep(time.Second)
			continue
		}

		break
	}

	return db, err
}
