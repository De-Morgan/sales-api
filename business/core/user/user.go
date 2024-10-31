package user

import (
	"context"
	"errors"
	"fmt"
	"sales-api/foundation/logger"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound              = errors.New("user not found")
	ErrUniqueEmail           = errors.New("email is not unique")
	ErrAuthenticationFailure = errors.New("authentication failed")
)

// Repository interface declares the behavior this package needs to perists and
// retrieve data.
type Repository interface {
	Create(ctx context.Context, usr User) error
	// Update(ctx context.Context, usr User) error
	// Delete(ctx context.Context, usr User) error
	// Query(ctx context.Context, filter QueryFilter, orderBy order.By, page int, pageSize int) ([]User, error)
	// Count(ctx context.Context, filter QueryFilter) (int, error)
	// QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
	// QueryByIDs(ctx context.Context, userID []uuid.UUID) ([]User, error)
	// QueryByEmail(ctx context.Context, email mail.Address) (User, error)
}

// =============================================================================

// Core manages the set of APIs for user access.
type Core struct {
	repository Repository
	log        *logger.Logger
}

// NewCore constructs a core for user api access.
func NewCore(log *logger.Logger, repository Repository) *Core {
	return &Core{
		repository: repository,
		log:        log,
	}
}

func (c *Core) Create(ctx context.Context, nu NewUser) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generatefrompassword: %w", err)
	}

	now := time.Now()

	usr := User{
		ID:           uuid.New(),
		Name:         nu.Name,
		Email:        nu.Email,
		Roles:        nu.Roles,
		PasswordHash: hash,
		Department:   nu.Department,
		Enabled:      true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := c.repository.Create(ctx, usr); err != nil {
		return User{}, fmt.Errorf("create: %w", err)
	}

	return usr, nil
}
