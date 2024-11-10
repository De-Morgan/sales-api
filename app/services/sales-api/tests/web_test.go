package tests

import (
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"sales-api/app/services/sales-api/handlers"
	"sales-api/business/data/test"
	v1 "sales-api/business/web/v1"
	"sales-api/business/web/v1/auth"
	"syscall"
	"testing"
)

// WebTests holds methods for each subtest. This type allows passing
// dependencies for tests while still providing a convenient syntax when
// subtests are registered.
type WebTest struct {
	app        http.Handler
	userToken  string
	adminToken string
	coreAPIs   test.CoreAPIs
	auth       *auth.Auth
	teardown   func()
}

func NewWebTest(t *testing.T) *WebTest {
	test := test.New(t)
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.Error(string(debug.Stack()))
		}
	}()

	api := test.CoreAPIs

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	handler := v1.APIMux(v1.APIMuxConfig{
		Shutdown: shutdown,
		Log:      test.Log,
		DB:       test.DB,
		Auth:     test.Auth,
	}, handlers.Routes())

	usrToken, err := test.TokenV1("user@example.com", "gophers")
	if err != nil {
		t.Fatal(err.Error())
	}
	adminToken, err := test.TokenV1("admin@example.com", "gophers")
	if err != nil {
		t.Fatal(err.Error())
	}
	tests := WebTest{
		app:        handler,
		userToken:  usrToken,
		adminToken: adminToken,
		coreAPIs:   api,
		auth:       test.Auth,
		teardown:   test.TearDown,
	}
	return &tests

}
func (t *WebTest) TearDown() {
	t.teardown()
}
