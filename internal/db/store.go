package db

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	QueryTimeoutDuration = time.Second * 5
	// common errors
	ErrNotFound = errors.New("resource not found")
	ErrConflict = errors.New("resource already exists")
	// users related error
	ErrDuplicateEmail    = errors.New("a user with that email already exists")
	ErrDuplicateUsername = errors.New("a user with that username already exists")
)

type Store struct {
	ConnPool *pgxpool.Pool
	Queries  *Queries
}

func NewPostgresStore(connPool *pgxpool.Pool) *Store {
	return &Store{
		ConnPool: connPool,
		Queries:  New(connPool),
	}
}

func ExecWithTx(ctx context.Context, s *Store, fn func(*Queries) error) error {
	conn, err := s.ConnPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	txq := s.Queries.WithTx(tx)

	if err := fn(txq); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
