package db

import (
	"database/sql"

	"github.com/eliasyoung/go-backend-server-practice/internal/env"
	"github.com/pressly/goose/v3"
)

type migrator struct {
	db *sql.DB
}

func NewMigrator(db *sql.DB) *migrator {
	return &migrator{
		db: db,
	}
}

func (m *migrator) Migrate() error {

	db := m.db

	path := env.GetDotEnvConfig("MIGRATION_PATH")

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, path); err != nil {
		return err
	}
	if err := db.Close(); err != nil {
		return err
	}

	return nil
}
