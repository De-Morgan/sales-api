package userdb

import (
	"context"
	"errors"
	"fmt"
	"sales-api/business/core/user"
	"sales-api/business/data/dbsql/pgx"
	"sales-api/foundation/logger"

	"github.com/jmoiron/sqlx"
)

type PostgresRepository struct {
	log *logger.Logger
	db  *sqlx.DB
}

var _ user.Repository = (*PostgresRepository)(nil)

func NewRepository(log *logger.Logger, db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{
		log: log,
		db:  db,
	}
}

func (s *PostgresRepository) Create(ctx context.Context, usr user.User) error {
	const q = `
	INSERT INTO users
		(user_id, name, email, password_hash, roles, enabled, department, created_at, updated_at)
	VALUES
		(:user_id, :name, :email, :password_hash, :roles, :enabled, :department, :created_at, :updated_at)`

	if err := pgx.NamedExecContext(ctx, s.log, s.db, q, toDBUser(usr)); err != nil {
		if errors.Is(err, pgx.ErrDBDuplicatedEntry) {
			return fmt.Errorf("namedexeccontext: %w", user.ErrUniqueEmail)
		}
		return fmt.Errorf("namedexeccontext: %w", err)

	}
	return nil
}
