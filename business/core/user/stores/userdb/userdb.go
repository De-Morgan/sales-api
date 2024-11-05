package userdb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/mail"
	"sales-api/business/core/user"
	"sales-api/business/data/dbsql/pgx"
	"sales-api/business/data/order"
	"sales-api/foundation/logger"

	"github.com/google/uuid"
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

func (r *PostgresRepository) QueryByID(ctx context.Context, userID uuid.UUID) (user.User, error) {
	data := struct {
		ID uuid.UUID `db:"user_id"`
	}{
		ID: userID,
	}
	const q = `
		SELECT
        user_id, name, email, password_hash, roles, enabled, department, created_at, updated_at
	FROM
		users
	WHERE
		user_id = :user_id`
	return r.queryUser(ctx, q, data)
}

func (r *PostgresRepository) QueryByEmail(ctx context.Context, email mail.Address) (user.User, error) {
	data := struct {
		Email string `db:"email"`
	}{
		Email: email.Address,
	}

	const q = `
		SELECT
        user_id, name, email, password_hash, roles, enabled, department, created_at, updated_at
	FROM
		users
	WHERE
		email = :email`

	return r.queryUser(ctx, q, data)

}

// Query retrieves a list of existing users from the database.
func (r *PostgresRepository) Query(ctx context.Context, filter user.QueryFilter, orderBy order.By, page int, pageSize int) ([]user.User, error) {
	data := map[string]any{
		"offset": (page - 1) * pageSize,
		"limit":  pageSize,
	}

	const q = `
	SELECT
        user_id, name, email, password_hash, roles, enabled, department, created_at, updated_at
	FROM
		users`

	buf := bytes.NewBufferString(q)
	r.applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}
	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY")

	var dbUsrs []dbUser
	if err := pgx.NamedQuerySlice(ctx, r.log, r.db, buf.String(), data, &dbUsrs); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	usrs, err := toCoreUserSlice(dbUsrs)
	if err != nil {
		return nil, err
	}
	return usrs, nil
}

// =======================================================================================================
func (r *PostgresRepository) queryUser(ctx context.Context, q string, data any) (user.User, error) {
	var dbUsr dbUser
	if err := pgx.NamedQueryStruct(ctx, r.log, r.db, q, data, &dbUsr); err != nil {
		if errors.Is(err, pgx.ErrDBNotFound) {
			return user.User{}, fmt.Errorf("namedquerystruct: %w", user.ErrNotFound)
		}
		return user.User{}, fmt.Errorf("namedquerystruct: %w", err)

	}
	usr, err := toCoreUser(dbUsr)
	if err != nil {
		return user.User{}, err
	}
	return usr, nil
}
