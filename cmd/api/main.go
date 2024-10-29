package main

import (
	"database/sql"
	"log"

	"github.com/eliasyoung/go-backend-server-practice/internal/env"
	"github.com/eliasyoung/go-backend-server-practice/internal/store"
)

func main() {
	cfg := config{
		addr: env.GetDotEnvConfigWithFallback("ADDR", ":8080"),
	}

	store := store.NewPostgresStorage(db * sql.DB)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Fatal((app.run(mux)))
}
