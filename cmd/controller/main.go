package main

import (
	"database/sql"
	"net/http"
	"net/url"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/controller/connman"
	"github.com/deviceplane/deviceplane/pkg/controller/service"
	mysql_store "github.com/deviceplane/deviceplane/pkg/controller/store/mysql"
	"github.com/deviceplane/deviceplane/pkg/email"
	"github.com/deviceplane/deviceplane/pkg/email/smtp"
	_ "github.com/deviceplane/deviceplane/pkg/statik"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/rakyll/statik/fs"
	"gopkg.in/alecthomas/kingpin.v2"
)

var version = "dev"
var name = "deviceplane-controller"

var (
	addr = kingpin.
		Flag("addr", "").
		Default(":8080").
		String()
	mysql = kingpin.
		Flag("mysql", "").
		Default("deviceplane:deviceplane@tcp(localhost:3306)/deviceplane?parseTime=true").
		String()
	statsdAddress = kingpin.
			Flag("statsd", "").
			Default("127.0.0.1:8125").
			String()
	allowedOrigins = kingpin.
			Flag("allowed-origin", "").
			Strings()
	emailProvider = kingpin.
			Flag("email-provider", "").
			Default("none").
			String()
	emailFromName = kingpin.
			Flag("email-from-name", "").
			Default("Deviceplane").
			String()
	emailFromAddress = kingpin.
				Flag("email-from-address", "").
				String()
	allowedEmailDomains = kingpin.
				Flag("allowed-email-domain", "").
				Strings()
	smtpServer = kingpin.
			Flag("smtp-server", "").
			String()
	smtpPort = kingpin.
			Flag("smtp-port", "").
			Int()
	smtpUsername = kingpin.
			Flag("smtp-username", "").
			String()
	smtpPassword = kingpin.
			Flag("smtp-password", "").
			String()
	dbMaxOpenConns = kingpin.
			Flag("db-max-open-conns", "50").
			Int()
	dbMaxIdleConns = kingpin.
			Flag("db-max-idle-conns", "25").
			Int()
	dbMaxConnLifetime = kingpin.
				Flag("db-max-conn-lifetime", "60m").
				Duration()
	auth0Domain = kingpin.
			Flag("auth0-domain", "").
			URL()
	auth0Audience = kingpin.
			Flag("auth0-audience", "").
			String()
)

func main() {
	kingpin.Parse()

	var allowedOriginURLs []url.URL
	for _, origin := range *allowedOrigins {
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

	db, err := tryConnect(*mysql)
	if err != nil {
		log.WithError(err).Fatal("connect to MySQL")
	}

	sqlStore := mysql_store.NewStore(db)

	st, err := statsd.New(*statsdAddress,
		statsd.WithNamespace("deviceplane."),
	)
	if err != nil {
		log.WithError(err).Fatal("statsd")
	}

	emailProvider := getEmailProvider(*emailProvider)

	connman := connman.New()

	svc := service.NewService(sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore,
		sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore,
		sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore,
		emailProvider, *emailFromName, *emailFromAddress, *allowedEmailDomains,
		*auth0Domain, *auth0Audience,
		statikFS, st, connman, allowedOriginURLs)

	server := &http.Server{
		Addr: *addr,
		Handler: handlers.CORS(
			handlers.AllowCredentials(),
			handlers.AllowedHeaders([]string{"Content-Type"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE"}),
			handlers.AllowedOrigins(*allowedOrigins),
		)(svc),
	}

	log.Info("Server will now listen on " + *addr)
	if err := server.ListenAndServe(); err != nil {
		log.WithError(err).Fatal("listen and serve")
	}
}

func tryConnect(uri string) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 0; i < 30; i++ {
		log.Info("attempting to connect to db")
		if db, err = sql.Open("mysql", uri); err != nil {
			time.Sleep(time.Second)
			continue
		}

		log.Info("attempting to ping db")
		if err = db.Ping(); err != nil {
			db.Close()
			time.Sleep(time.Second)
			continue
		}

		log.Info("connected to db")
		break
	}

	db.SetMaxOpenConns(*dbMaxOpenConns)
	db.SetMaxIdleConns(*dbMaxIdleConns)
	db.SetConnMaxLifetime(*dbMaxConnLifetime)

	return db, err
}

func getEmailProvider(emailProvider string) email.Interface {
	switch emailProvider {
	case "smtp":
		return smtp.NewEmail(
			*smtpServer,
			*smtpPort,
			*smtpUsername,
			*smtpPassword,
		)
	default:
		return nil
	}
}
