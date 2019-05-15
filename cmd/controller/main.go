package main

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/controller/service"
	sendgrid_email "github.com/deviceplane/deviceplane/pkg/email/sendgrid"
	"github.com/gomodule/redigo/redis"

	mysql_store "github.com/deviceplane/deviceplane/pkg/controller/store/mysql"
	redis_store "github.com/deviceplane/deviceplane/pkg/controller/store/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/segmentio/conf"
	"github.com/sendgrid/sendgrid-go"
)

var version = "dev"
var name = "deviceplane-controller"

var config struct {
	Addr          string        `conf:"addr"`
	MySQLPrimary  string        `conf:"mysql-primary"`
	Redis         string        `conf:"redis"`
	RedisTimeout  time.Duration `conf:"redis-timeout"`
	RedisConns    int           `conf:"redis-conns"`
	CookieDomain  string        `conf:"cookie-domain"`
	CookieSecure  bool          `conf:"cookie-secure"`
	AllowedOrigin string        `conf:"allowed-origin"`
}

func init() {
	config.Addr = ":8080"
	config.MySQLPrimary = "user:pass@tcp(localhost:3306)/deviceplane?parseTime=true"
	config.Redis = "localhost:6379"
	config.RedisTimeout = 5 * time.Second
	config.RedisConns = 10
	config.CookieDomain = "localhost"
	config.CookieSecure = false
	config.AllowedOrigin = "http://localhost:3000"
}

func main() {
	conf.Load(&config)

	db, err := tryConnect(config.MySQLPrimary)
	if err != nil {
		log.WithError(err).Fatal("connect to MySQL primary")
	}

	sqlStore := mysql_store.NewStore(db)

	redisPool := &redis.Pool{
		MaxIdle:   config.RedisConns,
		MaxActive: config.RedisConns,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", config.Redis,
				redis.DialConnectTimeout(config.RedisTimeout),
				redis.DialReadTimeout(config.RedisTimeout),
				redis.DialWriteTimeout(config.RedisTimeout),
			)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) (err error) {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err = c.Do("PING")
			return err
		},
	}

	redisStore := redis_store.NewStore(redisPool)

	sendgridClient := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	sendgridEmail := sendgrid_email.NewEmail(sendgridClient)

	svc := service.NewService(sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, redisStore,
		sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sqlStore, sendgridEmail, config.CookieDomain, config.CookieSecure)

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
		log.WithError(err).Fatal("listen and serve")
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
