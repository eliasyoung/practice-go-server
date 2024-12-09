package db

import (
	"context"
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
	ConnPool *pgxpool.Pool
	Queries  *Queries
}

func NewPostgresStore(connPool *pgxpool.Pool) *Store {
	return &Store{
		ConnPool: connPool,
		Queries:  New(connPool),
	}
}

func (s *Store) ExecWithTx(ctx context.Context, fn func(*Queries) error) error {
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
