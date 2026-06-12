package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/mingassing/app/backend/internal/api"
	"github.com/mingassing/app/backend/internal/store"
)

// The entrypoint, runs api router for mingas backend
// creates a connection to the database, creates a router instance,
// listens and serves the http router.
//
// example usage / running (from the backend/ directory):
// `DATABASE_URL=placeholder ADDR=placeholder go run ./cmd/server`
func main() {
	// this context only governs startup work (connecting to / pinging the db).
	// individual http requests each get their own context automatically, thats
	// the r.Context() you see in handlers
	ctx := context.Background()

	// this gets the DATABASE_URL from environment. Getting a variable from environment
	// prevents you from hard coding in sensitive data like "bitcoin_password := "123""
	//
	// the env var is passed in at run time, so for example youd run it with
	// DATABASE_URL=... go run ./cmd/server
	dbURL := os.Getenv("DATABASE_URL")
	// if none is provided, default to this
	if dbURL == "" {
		dbURL = "postgres://mingas:mingas@localhost:5432/mingas?sslmode=disable"
	}

	// migrate db (create and alter tables in order). if that produces an error, log it
	if err := store.Migrate(dbURL); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	// connect to the db instance via store module
	db, err := store.Connect(ctx, dbURL)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	// defer takes the rest of the line of code after it (i.e. db.Close())
	// and moves it to the end of the scope (like the function)
	defer db.Close()

	// create the router (decides what data gets sent back based on url thats hit)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// add the routes from api/
	api.Routes(r, db)

	// grab address from environment, default to 3000
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":3000"
	}
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
