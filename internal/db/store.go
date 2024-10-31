package db

import "github.com/jackc/pgx/v5/pgxpool"

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
