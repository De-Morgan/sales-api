package user

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"sales-api/business/data/order"
	"sales-api/business/data/transaction"
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
	ExecuteUnderTransaction(tx transaction.Transaction) (Repository, error)
	Create(ctx context.Context, usr User) error
	// Update(ctx context.Context, usr User) error
	// Delete(ctx context.Context, usr User) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page int, pageSize int) ([]User, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
	// QueryByIDs(ctx context.Context, userID []uuid.UUID) ([]User, error)
	QueryByEmail(ctx context.Context, email mail.Address) (User, error)
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

// ExecuteUnderTransaction constructs a new Core value that will use the
// specified transaction in any store related calls.
func (c *Core) ExecuteUnderTransaction(tx transaction.Transaction) (*Core, error) {
	trs, err := c.repository.ExecuteUnderTransaction(tx)
	if err != nil {
		return nil, err
	}
	c = &Core{
		repository: trs,
		log:        c.log,
	}

	return c, nil
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

func (c *Core) QueryByEmail(ctx context.Context, email mail.Address) (User, error) {
	user, err := c.repository.QueryByEmail(ctx, email)
	if err != nil {
		return User{}, fmt.Errorf("query: email[%s]: %w", email, err)
	}
	return user, nil
}

// QueryByID returns the user by it ID,
// returns "ErrNotFound" if the user record is not found
func (c *Core) QueryByID(ctx context.Context, userID uuid.UUID) (User, error) {
	user, err := c.repository.QueryByID(ctx, userID)
	if err != nil {
		return User{}, fmt.Errorf("query: user_id[%s]: %w", userID, err)
	}
	return user, nil
}

func (c *Core) Query(ctx context.Context, filter QueryFilter, orderBy order.By, page int, pageSize int) ([]User, error) {
	users, err := c.repository.Query(ctx, filter, orderBy, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	return users, nil
}

// Count returns the total number of users.
func (c *Core) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return c.repository.Count(ctx, filter)
}

// ============================================================================

// Authenticate finds a user by their email and verifies their password. On
// success it returns a Claims User representing this user. The claims can be
// used to generate a token for future authentication.

func (c *Core) Authenticate(ctx context.Context, email mail.Address, pass string) (User, error) {
	usr, err := c.QueryByEmail(ctx, email)
	if err != nil {
		return User{}, fmt.Errorf("query: email[%s]: %w", email, err)
	}
	if err := bcrypt.CompareHashAndPassword(usr.PasswordHash, []byte(pass)); err != nil {
		return User{}, fmt.Errorf("comparehashandpassword: %w", ErrAuthenticationFailure)
	}
	return usr, nil
}
