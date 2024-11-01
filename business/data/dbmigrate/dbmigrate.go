package dbmigrate

import (
	"context"
	"fmt"
	"sales-api/business/data/dbsql/pgx"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

func Migration(ctx context.Context, source string, db *sqlx.DB) error {
	if err := pgx.StatusCheck(ctx, db); err != nil {
		return fmt.Errorf("status check database: %w", err)
	}
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("postgres instance: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(source, "postgres", driver)
	if err != nil {
		return fmt.Errorf("migration: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration up: %w", err)
	}
	return err
}
