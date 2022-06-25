package config

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"

	"github.com/Lameaux/core/httpclient"
	"github.com/Lameaux/core/logger"
	"github.com/Lameaux/core/utils"
)

var (
	EnvDevelopment = "development"
	EnvTest        = "test"
	EnvProduction  = "production"
)

type AppConfig struct {
	Env string

	Port    string
	AppName string
	Version string

	DBPool     *pgxpool.Pool
	HTTPClient *httpclient.HTTPClient

	WaitTimeout       time.Duration
	ConnectionTimeout time.Duration
	TLSTimeout        time.Duration
	ReadTimeout       time.Duration
}

const (
	dbPingTimeout = 2 * time.Second
)

func NewAppConfig(env string) *AppConfig {
	logger.Infow("loading env configuration", "env", env)

	if err := godotenv.Load(".env."+env, ".env"); err != nil {
		logger.Fatalw("failed to load env", "env", env, "error", err)

		return nil
	}

	port := utils.GetEnv("PORT")

	waitTimeout, err := strconv.Atoi(utils.GetEnv("WAIT_TIMEOUT"))
	if err != nil {
		logger.Fatal(err)
	}

	connectionTimeout, err := strconv.Atoi(utils.GetEnv("CONNECTION_TIMEOUT"))
	if err != nil {
		logger.Fatal(err)
	}

	tlsTimeout, err := strconv.Atoi(utils.GetEnv("TLS_TIMEOUT"))
	if err != nil {
		logger.Fatal(err)
	}

	readTimeout, err := strconv.Atoi(utils.GetEnv("READ_TIMEOUT"))
	if err != nil {
		logger.Fatal(err)
	}

	app := &AppConfig{
		Env: env,

		Port: port,

		WaitTimeout:       time.Duration(waitTimeout) * time.Second,
		ConnectionTimeout: time.Duration(connectionTimeout) * time.Second,
		TLSTimeout:        time.Duration(tlsTimeout) * time.Second,
		ReadTimeout:       time.Duration(readTimeout) * time.Second,
	}

	app.configureHTTPClient()

	databaseURI := os.Getenv("DATABASE_URI")
	if databaseURI != "" {
		app.configureDBPool(databaseURI)
	}

	return app
}

func (app *AppConfig) IsProduction() bool {
	return app.Env == EnvProduction
}

func (app *AppConfig) configureHTTPClient() {
	app.HTTPClient = httpclient.NewBuilder().
		ConnectionTimeout(app.ConnectionTimeout).
		TLSTimeout(app.TLSTimeout).
		ReadTimeout(app.ReadTimeout).
		Build()
}

func (app *AppConfig) configureDBPool(uri string) {
	logger.Infow("connecting to db", "database_uri", uri)

	pool, err := pgxpool.Connect(context.Background(), uri)
	if err != nil {
		logger.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbPingTimeout)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		logger.Fatal(err)
	}

	app.DBPool = pool
}

func (app *AppConfig) closeDBPool() {
	logger.Infow("closing db pool")
	app.DBPool.Close()
}

func (app *AppConfig) Shutdown() {
	if app.DBPool != nil {
		app.closeDBPool()
	}
}
