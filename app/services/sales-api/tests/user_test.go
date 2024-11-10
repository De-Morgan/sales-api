package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sales-api/app/services/sales-api/handlers/usergrp"
	"sales-api/business/web/v1/response"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
	web *WebTest
}

func (s *UserTestSuite) SetupSuite() {
	s.web = NewWebTest(s.T())
}
func (s *UserTestSuite) TearDownSuite() {
	s.web.TearDown()
}

// ==================================================

func (suite *UserTestSuite) TestQueryUsers() {
	url := "/v1/users?page=1&page_size=2&orderBy=user_id,DESC"
	r := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+suite.web.adminToken)

	suite.web.app.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		suite.Errorf(http.ErrServerClosed, "Should receive a status code of 200 for the response : %d", w.Code)
	}

	var resp response.PageDocument[usergrp.AppUser]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		suite.T().Errorf("Should be able to unmarshal the response : %s", err)
	}
	//suite.NotEmpty(resp)

	//suite.T().Logf("\nresponse is: %#v", resp)

}

// ================================================
func TestUser(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
