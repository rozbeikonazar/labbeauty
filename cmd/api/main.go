package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"sync"
	"time"

	"cosmetcab.dp.ua/internal/data"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	blobURL       = goDotEnvVariable("BLOB_URL")
	containerName = goDotEnvVariable("CONTAINER_NAME")
	botToken      = goDotEnvVariable("botToken")
	chatID        = goDotEnvVariable("chatID")
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

type application struct {
	config           config
	logger           *slog.Logger
	models           data.Models
	azureBlobStorage *AzureBlobStorage
	wg               sync.WaitGroup
	sessionManager   *sessions.CookieStore
}

func goDotEnvVariable(key string) string {
	_ = godotenv.Load(".env")

	return os.Getenv(key)
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", goDotEnvVariable("LABBEAUTY_DB_DSN"), "PostgreSQL DSN")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connectiond")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable limiter")
	flag.Parse()
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	azureBlobStorage, err := NewAzureBlobStorage(blobURL, credential, ctx)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)

	}
	defer db.Close()

	logger.Info("DB connection pool established")

	var store = sessions.NewCookieStore([]byte(goDotEnvVariable("TOKEN")))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60,
		HttpOnly: true,
	}

	app := &application{
		config:           cfg,
		logger:           logger,
		models:           data.NewModels(db),
		azureBlobStorage: azureBlobStorage,
		sessionManager:   store,
	}
	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
