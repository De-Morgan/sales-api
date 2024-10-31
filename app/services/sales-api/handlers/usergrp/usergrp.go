package usergrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sales-api/business/core/user"
	"sales-api/business/web/v1/auth"
	"sales-api/business/web/v1/response"
	"sales-api/foundation/web"
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	user *user.Core
	auth *auth.Auth
}

// New constructs a handlers for route access.
func New(user *user.Core, auth *auth.Auth) *Handlers {
	return &Handlers{
		user: user,
		auth: auth,
	}
}

// Create adds a new user to the system.
func (h *Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app AppNewUser
	if err := web.Decode(r, &app); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	nc, err := toCoreNewUser(app)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	usr, err := h.user.Create(ctx, nc)
	if err != nil {
		if errors.Is(err, user.ErrUniqueEmail) {
			return response.NewError(err, http.StatusConflict)
		}
		return fmt.Errorf("create: usr[%+v]: %w", usr, err)
	}

	return web.Respond(ctx, w, toAppUser(usr), http.StatusCreated)
}
