package store

import (
	"context"
	"database/sql"
)

type PostgresUsersStore struct {
	db *sql.DB
}

func (s *PostgresUsersStore) Create(ctx context.Context) error {
	return nil
}
