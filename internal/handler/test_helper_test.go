package handler

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
    t.Helper()

    ctx := context.Background()

    pgContainer, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:15-alpine"),
        postgres.WithDatabase("test-db"),
        postgres.WithUsername("user"),
        postgres.WithPassword("password"),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(2).
                WithStartupTimeout(1*time.Minute),
        ),
    )
    if err != nil {
        t.Fatalf("failed to start postgres container: %s", err)
    }

    connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
    if err != nil {
        t.Fatalf("failed to get connection string: %s", err)
    }

    db, err := sql.Open("pgx", connStr)
    if err != nil {
        t.Fatalf("failed to connect to database: %s", err)
    }

    if err := db.Ping(); err != nil {
        t.Fatalf("failed to ping database: %s", err)
    }

    migrationsDir := "../migrations"

    if err := goose.SetDialect("postgres"); err != nil {
        t.Fatalf("failed to set goose dialect: %s", err)
    }

    if err := goose.Up(db, migrationsDir); err != nil {
        t.Fatalf("failed to apply goose migrations: %s", err)
    }

    log.Println("Test database successfully prepared and migrations applied.")

    cleanup := func() {
        log.Println("Cleaning up test environment: stopping container.")
        if err := db.Close(); err != nil {
            t.Logf("warning: failed to close test DB connection: %s", err)
        }
        if err := pgContainer.Terminate(ctx); err != nil {
            t.Fatalf("failed to terminate postgres container: %s", err)
        }
    }

    return db, cleanup
}