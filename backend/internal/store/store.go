package store

import (
	"context"
	"embed"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// this compiler directive tells it to bake the migrations into the exe, its not
// a normal comment!

//go:embed migrations/*.sql
var migrations embed.FS

// Store owns the pool, because the pool is lowercase it is private
// this means the api routes or any external module which isnt Store cannot
// write sql into the store, only here can we actually interact with the db
type Store struct {
	pool *pgxpool.Pool
}

// create a connection to the postgres database instance
func Connect(ctx context.Context, url string) (*Store, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return &Store{pool: pool}, nil
}

// Close the connection to the postgres database instance
func (s *Store) Close() { s.pool.Close() }

// Ping the database to ensure that it is still running and queryable
func (s *Store) Ping(ctx context.Context) error { return s.pool.Ping(ctx) }

// run the migrations into the postgres db so that the tables are up to date
func Migrate(url string) error {
	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	db, err := goose.OpenDBWithDriver("pgx", url)
	if err != nil {
		return err
	}
	defer db.Close()
	return goose.Up(db, "migrations")
}
