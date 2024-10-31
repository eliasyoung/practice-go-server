package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetConnPool(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) *pgxpool.Pool {
	connConfig, err := pgxpool.ParseConfig(addr)
	if err != nil {
		log.Panicf("Database url not valid: %s", addr)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connPool, err := pgxpool.NewWithConfig(ctx, connConfig)

	if err != nil {
		log.Panicf("Database connect with addr %s failed", addr)
	}

	log.Println("DB connection pool established successfully!")

	return connPool
}
