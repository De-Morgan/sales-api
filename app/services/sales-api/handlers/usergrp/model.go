package usergrp

import (
	"fmt"
	"net/mail"
	"sales-api/business/core/user"
	"sales-api/foundation/validate"
	"time"
)

// AppUser represents information about an individual user.
type AppUser struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	Roles        []string `json:"roles"`
	PasswordHash []byte   `json:"-"`
	Department   string   `json:"department"`
	Enabled      bool     `json:"enabled"`
	CreatedAt    string   `json:"createdAt"`
	UpdatedAt    string   `json:"updatedAt"`
}

func toAppUser(usr user.User) AppUser {
	roles := make([]string, len(usr.Roles))
	for i, role := range usr.Roles {
		roles[i] = role.Name()
	}

	return AppUser{
		ID:           usr.ID.String(),
		Name:         usr.Name,
		Email:        usr.Email.Address,
		Roles:        roles,
		PasswordHash: usr.PasswordHash,
		Department:   usr.Department,
		Enabled:      usr.Enabled,
		CreatedAt:    usr.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    usr.UpdatedAt.Format(time.RFC3339),
	}
}

// AppNewUser contains information needed to create a new user.
type AppNewUser struct {
	Name            string   `json:"name" validate:"required"`
	Email           string   `json:"email" validate:"required,email"`
	Roles           []string `json:"roles" validate:"required"`
	Department      string   `json:"department"`
	Password        string   `json:"password" validate:"required"`
	PasswordConfirm string   `json:"passwordConfirm" validate:"eqfield=Password"`
}

func toCoreNewUser(app AppNewUser) (user.NewUser, error) {
	roles := make([]user.Role, len(app.Roles))
	for i, roleStr := range app.Roles {
		role, err := user.ParseRole(roleStr)
		if err != nil {
			return user.NewUser{}, fmt.Errorf("parsing role: %w", err)
		}
		roles[i] = role
	}

	addr, err := mail.ParseAddress(app.Email)
	if err != nil {
		return user.NewUser{}, fmt.Errorf("parsing email: %w", err)
	}

	usr := user.NewUser{
		Name:       app.Name,
		Email:      *addr,
		Roles:      roles,
		Department: app.Department,
		Password:   app.Password,
	}

	return usr, nil
}

// Validate checks the data in the model is considered clean.
func (app AppNewUser) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}
	return nil
}

func toAppUsers(users []user.User) []AppUser {
	items := make([]AppUser, len(users))
	for i, usr := range users {
		items[i] = toAppUser(usr)
	}

	return items
}

// AppLoginRequest contains information needed to login a new user.
type AppLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Validate checks the data in the model is considered clean.
func (app AppLoginRequest) Validate() error {
	if err := validate.Check(app); err != nil {
		return err
	}
	return nil
}

// AppLoginResponse contains information returned after user login.
type AppLoginResponse struct {
	User  AppUser `json:"user"`
	Token string  `json:"token"`
}
