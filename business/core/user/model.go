package user

import (
	"net/mail"
	"time"

	"github.com/google/uuid"
)

// User represents information about an individual user.
type User struct {
	ID           uuid.UUID
	Name         string
	Email        mail.Address
	Roles        []Role
	PasswordHash []byte
	Department   string
	Enabled      bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewUser contains information needed to create a new user.
type NewUser struct {
	Name       string
	Email      mail.Address
	Roles      []Role
	Department string
	Password   string
}

type UpdateUser struct {
	Name       *string
	Email      *mail.Address
	Roles      []Role
	Department *string
	Enabled    *bool
}
