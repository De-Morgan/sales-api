package userdb

import (
	"database/sql"
	"fmt"
	"net/mail"
	"sales-api/business/core/user"
	"sales-api/business/data/dbsql/pgx/dbarray"
	"time"

	"github.com/google/uuid"
)

// dbUser represent the structure we need for moving data
// between the app and the database.
type dbUser struct {
	ID           uuid.UUID      `db:"user_id"`
	Name         string         `db:"name"`
	Email        string         `db:"email"`
	Roles        dbarray.String `db:"roles"`
	PasswordHash []byte         `db:"password_hash"`
	Department   sql.NullString `db:"department"`
	Enabled      bool           `db:"enabled"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
}

func toDBUser(usr user.User) dbUser {
	roles := make([]string, len(usr.Roles))
	for i, role := range usr.Roles {
		roles[i] = role.Name()
	}
	return dbUser{
		ID:           usr.ID,
		Name:         usr.Name,
		Email:        usr.Email.Address,
		Roles:        roles,
		PasswordHash: usr.PasswordHash,
		Department: sql.NullString{
			String: usr.Department,
			Valid:  usr.Department != "",
		},
		Enabled:   usr.Enabled,
		CreatedAt: usr.CreatedAt.UTC(),
		UpdatedAt: usr.UpdatedAt.UTC(),
	}

}

func toCoreUser(dbUsr dbUser) (user.User, error) {
	addr := mail.Address{
		Address: dbUsr.Email,
	}

	roles := make([]user.Role, len(dbUsr.Roles))
	for i, value := range dbUsr.Roles {
		var err error
		roles[i], err = user.ParseRole(value)
		if err != nil {
			return user.User{}, fmt.Errorf("parse role: %w", err)
		}
	}

	usr := user.User{
		ID:           dbUsr.ID,
		Name:         dbUsr.Name,
		Email:        addr,
		Roles:        roles,
		PasswordHash: dbUsr.PasswordHash,
		Enabled:      dbUsr.Enabled,
		Department:   dbUsr.Department.String,
		CreatedAt:    dbUsr.CreatedAt.In(time.Local),
		UpdatedAt:    dbUsr.UpdatedAt.In(time.Local),
	}

	return usr, nil
}

func toCoreUserSlice(dbUsers []dbUser) ([]user.User, error) {
	usrs := make([]user.User, len(dbUsers))
	for i, dbUsr := range dbUsers {
		var err error
		usrs[i], err = toCoreUser(dbUsr)
		if err != nil {
			return nil, err
		}
	}
	return usrs, nil
}