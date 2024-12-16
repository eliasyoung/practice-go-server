package main

import (
	"time"

	"github.com/eliasyoung/go-backend-server-practice/internal/db"
	"github.com/eliasyoung/go-backend-server-practice/internal/env"
	"github.com/eliasyoung/go-backend-server-practice/internal/service"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			GO BACKEND API
//	@description	This is a sample server with chi
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host						petstore.swagger.io
// @BasePath					/v1
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	cfg := config{
		addr: env.GetDotEnvConfigWithFallback("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetDotEnvConfigWithFallback("DATABASE_URL", "postgres://admin:adminpassword@localhost/social?sslmode=disbale"),
			maxOpenConns: env.GetIntDotEnvConfigWithFallback("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetIntDotEnvConfigWithFallback("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetDotEnvConfigWithFallback("DB_MAX_IDLE_TIME", "15min"),
		},
		env:    env.GetDotEnvConfigWithFallback("ENV", "development"),
		apiURL: env.GetDotEnvConfigWithFallback("EXTERNAL_URL", "localhost:8080"),
		mail: mailConfig{
			exp: time.Hour * 24 * 3,
		},
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// setup logger for package db first before calling any db func
	db.SetupLogger(logger)
	service.SetupLogger(logger)

	conn := db.GetConnPool(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	defer conn.Close()

	dbb := stdlib.OpenDBFromPool(conn)

	mig := db.NewMigrator(dbb)

	if err := mig.Migrate(); err != nil {
		logger.Panicf("Error occured when migrating: %s", err)
	}

	store := db.NewPostgresStore(conn)

	services := service.Services{
		User: service.NewUserService(store),
		Post: service.NewPostService(store),
	}

	app := &application{
		config:  cfg,
		store:   *store,
		service: services,
		logger:  logger,
	}

	mux := app.mount()
	logger.Fatal((app.run(mux)))
}
