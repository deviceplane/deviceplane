package main

import (
	"database/sql"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/controller/connman"
	"github.com/deviceplane/deviceplane/pkg/controller/runner"
	"github.com/deviceplane/deviceplane/pkg/controller/runner/datadog"
	"github.com/deviceplane/deviceplane/pkg/controller/service"
	mysql_store "github.com/deviceplane/deviceplane/pkg/controller/store/mysql"
	"github.com/deviceplane/deviceplane/pkg/email"
	sendgrid_email "github.com/deviceplane/deviceplane/pkg/email/sendgrid"
	_ "github.com/deviceplane/deviceplane/pkg/statik"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/rakyll/statik/fs"
	"github.com/segmentio/conf"
	"github.com/sendgrid/sendgrid-go"
)

var version = "dev"
var name = "deviceplane-controller"

var config struct {
	Addr           string   `conf:"addr"`
	MySQLPrimary   string   `conf:"mysql-primary"`
	Statsd         string   `conf:"statsd"`
	AllowedOrigins []string `conf:"allowed-origins"`
	EmailProvider  string   `conf:"email-provider"`
}

func init() {
	config.Addr = ":8080"
	config.MySQLPrimary = "deviceplane:deviceplane@tcp(localhost:3306)/deviceplane?parseTime=true"
	config.Statsd = "127.0.0.1:8125"
	config.AllowedOrigins = []string{}
	config.EmailProvider = "none"
}

func main() {
	conf.Load(&config)

	var allowedOriginURLs []url.URL
	for _, origin := range config.AllowedOrigins {
		originURL, err := url.Parse(origin)
		if err != nil {
			log.WithError(err).Fatal("parsing allowed origin url: " + origin)
		}
		allowedOriginURLs = append(allowedOriginURLs, *originURL)
	}

	statikFS, err := fs.New()
	if err != nil {
		log.WithError(err).Fatal("statik")
	}

	db, err := tryConnect(config.MySQLPrimary)
	if err != nil {
		log.WithError(err).Fatal("connect to MySQL primary")
	}

	sqlStore := mysql_store.NewStore(db)

	st, err := statsd.New(config.Statsd,
		statsd.WithNamespace("deviceplane."),
	)
	if err != nil {
		log.WithError(err).Fatal("statsd")
	}

	emailProvider := getEmailProvider(config.EmailProvider)

	connman := connman.New()

	runnerManager := runner.NewManager([]runner.Runner{
		datadog.NewRunner(sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, st, connman),
	})
	runnerManager.Start()

	svc := service.NewService(sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore,
		sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore,
		emailProvider, statikFS, st, connman, allowedOriginURLs)

	server := &http.Server{
		Addr: config.Addr,
		Handler: handlers.CORS(
			handlers.AllowCredentials(),
			handlers.AllowedHeaders([]string{"Content-Type"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE"}),
			handlers.AllowedOrigins(config.AllowedOrigins),
		)(svc),
	}

	if err := server.ListenAndServe(); err != nil {
		log.WithError(err).Fatal("listen and serve")
	}
}

func tryConnect(uri string) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 0; i < 30; i++ {
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

func getEmailProvider(emailProvider string) email.Interface {
	switch emailProvider {
	case "sendgrid":
		sendgridClient := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
		return sendgrid_email.NewEmail(sendgridClient)
	default:
		return nil
	}
}
