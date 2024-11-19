package db

import (
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	QueryTimeoutDuration = time.Second * 5
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
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
