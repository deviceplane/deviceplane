package main

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/deviceplane/deviceplane/pkg/controller/service"
	sendgrid_email "github.com/deviceplane/deviceplane/pkg/email/sendgrid"

	mysql_store "github.com/deviceplane/deviceplane/pkg/controller/store/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/segmentio/conf"
	"github.com/sendgrid/sendgrid-go"
)

var version = "dev"
var name = "deviceplane-controller"

var config struct {
	Addr          string `conf:"addr"`
	MySQLPrimary  string `conf:"mysql-primary"`
	CookieDomain  string `conf:"cookie-domain"`
	CookieSecure  bool   `conf:"cookie-secure"`
	AllowedOrigin string `conf:"allowed-origin"`
}

func init() {
	config.Addr = ":8080"
	config.MySQLPrimary = "user:pass@tcp(localhost:3306)/deviceplane?parseTime=true"
	config.CookieDomain = "localhost"
	config.CookieSecure = false
	config.AllowedOrigin = "http://localhost:3000"
}

func main() {
	conf.Load(&config)

	db, err := tryConnect(config.MySQLPrimary)
	if err != nil {
		panic(err)
	}

	store := mysql_store.NewStore(db)

	sendgridClient := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	sendgridEmail := sendgrid_email.NewEmail(sendgridClient)

	svc := service.NewService(store, store, store, store, store, store, store, store, store, store,
		store, store, store, store, sendgridEmail, config.CookieDomain, config.CookieSecure)

	server := &http.Server{
		Addr: config.Addr,
		Handler: handlers.CORS(
			handlers.AllowCredentials(),
			handlers.AllowedHeaders([]string{"Content-Type"}),
			handlers.AllowedMethods([]string{"GET", "POST", "DELETE"}),
			handlers.AllowedOrigins([]string{config.AllowedOrigin}),
		)(svc),
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
