package user_test

import (
	"context"
	"net/mail"
	"sales-api/business/core/user"
	"sales-api/business/data/test"
	"testing"

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

	usr, err := suite.test.CoreAPIs.User.Create(context.Background(), nu)
	suite.NoError(err)
	suite.NotEmpty(usr)
	suite.Equal(nu.Name, usr.Name)
	suite.Equal(nu.Department, usr.Department)
	// Test duplicate entry
	_, err = suite.test.CoreAPIs.User.Create(context.Background(), nu)
	suite.Error(err)
	suite.ErrorIs(err, user.ErrUniqueEmail)
}

// ================================================
func TestUser(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
