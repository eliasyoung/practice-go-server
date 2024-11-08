package db

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	QueryTimeoutDuration = time.Second * 5
)

type Store struct {
	connPool *pgxpool.Pool
	Queries  *Queries
}

func NewPostgresStore(connPool *pgxpool.Pool) *Store {
	return &Store{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
