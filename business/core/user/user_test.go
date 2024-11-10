package user_test

import (
	"context"
	"net/mail"
	"sales-api/business/core/user"
	"sales-api/business/data/test"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
	test *test.Test
}

func (s *UserTestSuite) SetupSuite() {
	s.test = test.New(s.T())
}
func (s *UserTestSuite) TearDownSuite() {
	s.test.TearDown()
}

// ==================================================

func (suite *UserTestSuite) TestCreate() {
	email, err := mail.ParseAddress("morgan@gmail.com")
	suite.NoError(err)
	nu := user.NewUser{
		Name:  "Morgan",
		Email: *email,
		Roles: []user.Role{
			user.RoleUser,
		},
		Department: "Technology",
		Password:   "password",
	}
	suite.createUser(nu)
	// Test duplicate entry
	_, err = suite.test.CoreAPIs.User.Create(context.Background(), nu)
	suite.Error(err)
	suite.ErrorIs(err, user.ErrUniqueEmail)
}
func (suite *UserTestSuite) TestQueryByID() {

	email, err := mail.ParseAddress("morgan2@gmail.com")
	suite.NoError(err)
	nu := user.NewUser{
		Name:  "Morgan",
		Email: *email,
		Roles: []user.Role{
			user.RoleUser,
		},
		Department: "Technology",
		Password:   "password",
	}
	// Create new user
	usr := suite.createUser(nu)
	// Test query by id
	qusr, err := suite.test.CoreAPIs.User.QueryByID(context.Background(), usr.ID)
	suite.NoError(err)
	suite.Equal(usr.Email, qusr.Email)
	suite.Equal(usr.Name, qusr.Name)
	suite.Equal(usr.ID.String(), qusr.ID.String())

	// Test query by id not found
	_, err = suite.test.CoreAPIs.User.QueryByID(context.Background(), uuid.New())
	suite.Error(err)
	suite.ErrorIs(err, user.ErrNotFound)

}

func (suite *UserTestSuite) createUser(nu user.NewUser) user.User {
	usr, err := suite.test.CoreAPIs.User.Create(context.Background(), nu)
	suite.NoError(err)
	suite.NotEmpty(usr)
	suite.Equal(nu.Name, usr.Name)
	suite.Equal(nu.Department, usr.Department)
	return usr
}

// ================================================
func TestUser(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
