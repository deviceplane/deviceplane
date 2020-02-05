package main

import (
	"database/sql"
	"net/http"
	"net/url"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/controller/connman"
	"github.com/deviceplane/deviceplane/pkg/controller/runner"
	"github.com/deviceplane/deviceplane/pkg/controller/runner/datadog"
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

//commit for test
var (
	addr = kingpin.
		Flag("addr", "").
		Default(":8080").
		OverrideDefaultFromEnvar("ADDR").
		String()
	mysql = kingpin.
		Flag("mysql", "").
		Default("deviceplane:deviceplane@tcp(localhost:3306)/deviceplane?parseTime=true").
		OverrideDefaultFromEnvar("MYSQL_CONNECTION_STRING").
		String()
	enableSignUp = kingpin.
			Flag("enableSignUp", "").
			OverrideDefaultFromEnvar("ENABLE_SING_UP").
			Bool()
	statsdAddress = kingpin.
			Flag("statsd", "").
			Default("127.0.0.1:8125").
			OverrideDefaultFromEnvar("STATS_ADDRESS").
			String()
	allowedOrigins = kingpin.
			Flag("allowed-origin", "").
			OverrideDefaultFromEnvar("ALLOWED_ORIGINS").
			Strings()
	emailProvider = kingpin.
			Flag("email-provider", "").
			Default("none").
			OverrideDefaultFromEnvar("EMAIL_PROVIDER").
			String()
	emailFromName = kingpin.
			Flag("email-from-name", "").
			Default("Deviceplane").
			OverrideDefaultFromEnvar("EMAIL_FROM_NAME").
			String()
	emailFromAddress = kingpin.
				Flag("email-from-address", "").
				OverrideDefaultFromEnvar("EMAIL_FROM_ADDRESS").
				String()
	allowedEmailDomains = kingpin.
				Flag("allowed-email-domain", "").
				OverrideDefaultFromEnvar("ALLOWED_EMAIL_DOMAINS").
				Strings()
	smtpServer = kingpin.
			Flag("smtp-server", "").
			OverrideDefaultFromEnvar("SMTP_SERVER").
			String()
	smtpPort = kingpin.
			Flag("smtp-port", "").
			OverrideDefaultFromEnvar("SMTP_PORT").
			Int()
	smtpUsername = kingpin.
			Flag("smtp-username", "").
			OverrideDefaultFromEnvar("SMTP_USERNAME").
			String()
	smtpPassword = kingpin.
			Flag("smtp-password", "").
			OverrideDefaultFromEnvar("SMTP_PASSWORD").
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

	runnerManager := runner.NewManager([]runner.Runner{
		datadog.NewRunner(sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, st, connman),
	})
	runnerManager.Start()

	svc := service.NewService(sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore,
		sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, *enableSignUp,
		emailProvider, *emailFromName, *emailFromAddress, *allowedEmailDomains, statikFS, st, connman, allowedOriginURLs)

	server := &http.Server{
		Addr: *addr,
		Handler: handlers.CORS(
			handlers.AllowCredentials(),
			handlers.AllowedHeaders([]string{"Content-Type"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE"}),
			handlers.AllowedOrigins(*allowedOrigins),
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
