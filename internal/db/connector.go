package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetConnPool(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) *pgxpool.Pool {
	connConfig, err := pgxpool.ParseConfig(addr)
	if err != nil {
		logger.Panicf("Database url not valid: %s", addr)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connPool, err := pgxpool.NewWithConfig(ctx, connConfig)

	if err != nil {
		logger.Panicf("Database connect with addr %s failed", addr)
	}

	logger.Info("DB connection pool established successfully!")

	return connPool
}
