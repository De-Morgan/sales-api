package usergrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"sales-api/business/core/user"
	"sales-api/business/data/page"
	"sales-api/business/data/transaction"
	"sales-api/business/web/v1/auth"
	"sales-api/business/web/v1/mid"
	"sales-api/business/web/v1/response"
	"sales-api/foundation/validate"
	"sales-api/foundation/web"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func (h *Handlers) executeUnderTransaction(ctx context.Context) (*Handlers, error) {
	if tx, ok := transaction.Get(ctx); ok {
		user, err := h.user.ExecuteUnderTransaction(tx)
		if err != nil {
			return nil, err
		}
		h = &Handlers{
			user: user,
			auth: h.auth,
		}
		return h, nil
	}
	return h, nil
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
			return response.NewError(user.ErrUniqueEmail, http.StatusConflict)
		}
		return fmt.Errorf("create: usr[%+v]: %w", usr, err)
	}

	return web.Respond(ctx, w, userResponse(usr), http.StatusCreated)
}

// Create adds a new user to the system.
func (h *Handlers) Login(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	kid := r.Header.Get("kid")
	if kid == "" {
		return validate.NewFieldsError("kid", errors.New("missing kid"))
	}
	var app AppLoginRequest
	if err := web.Decode(r, &app); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}
	email, err := mail.ParseAddress(app.Email)
	if err != nil {
		return validate.NewFieldsError("email", errors.New("invalid email"))
	}
	usr, err := h.user.Authenticate(ctx, *email, app.Password)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrAuthenticationFailure):
			return auth.NewAuthError(err.Error())
		case errors.Is(err, user.ErrNotFound):
			return response.NewError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("login: usr[%+v]: %w", usr, err)
		}
	}
	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   usr.ID.String(),
			Issuer:    "service project",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
		},
		Roles: usr.Roles,
	}
	token, err := h.auth.GenerateToken(kid, claims)

	if err != nil {
		return fmt.Errorf("generatetoken: %w", err)
	}

	return web.Respond(ctx, w, response.NewSuccess(AppLoginResponse{
		User:  toAppUser(usr),
		Token: token,
	}), http.StatusCreated)
}

// QueryByID returns a user by its ID.
func (h *Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := web.Param(r, "user_id")
	if id == "" {
		return response.NewError(errors.New("id is required"), http.StatusBadRequest)
	}
	uuid, err := uuid.Parse(id)
	if err != nil {
		return response.NewError(mid.ErrInvalidID, http.StatusBadRequest)
	}
	usr, err := h.user.QueryByID(ctx, uuid)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotFound):
			return response.NewError(user.ErrNotFound, http.StatusNotFound)
		default:
			return fmt.Errorf("QueryByID: id:[%q] :%w", id, err)
		}
	}
	return web.Respond(ctx, w, userResponse(usr), http.StatusOK)
}

// Query returns a list of users with paging.
func (h *Handlers) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	page, err := page.Parse(r)
	if err != nil {
		return err
	}
	filter, err := parseFilter(r)
	if err != nil {
		return err
	}

	orderBy, err := parseOrder(r)

	if err != nil {
		return err
	}

	users, err := h.user.Query(ctx, filter, orderBy, page.Page, page.PageSize)

	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	total, err := h.user.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	return web.Respond(ctx, w, response.NewPageDocument(toAppUsers(users), total, page.Page, page.PageSize), http.StatusOK)
}
