package test

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"net/mail"
	"sales-api/business/core/user"
	"sales-api/business/core/user/stores/userdb"
	"sales-api/business/web/v1/auth"
	"sales-api/foundation/logger"
	"sales-api/foundation/web"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
)

// Test owns state for running and shutting down tests.
type Test struct {
	TestDatabase
	Log      *logger.Logger
	CoreAPIs CoreAPIs
	TearDown func()
	tb       testing.TB
	Auth     *auth.Auth
}

func New(tb testing.TB) *Test {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	const letterByte = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, 5)
	for i := range b {
		b[i] = letterByte[rand.Intn(len(letterByte))]
	}
	dbName := string(b)

	tb.Log("Waiting for database to be ready ...")

	testDB := SetUpTestDatabase(ctx, dbName)

	tb.Log("Database ready")

	//  ------------------------------------------------------------
	var buf bytes.Buffer
	log := logger.New(&buf, logger.LevelInfo, "TEST", func(ctx context.Context) string {
		return web.GetTraceID(ctx)
	})

	coreAPIs := newCoreAPIs(log, testDB.DB)

	tb.Log("Ready for testing ...")
	//  ------------------------------------------------------------

	cfg := auth.Config{
		Log:       log,
		KeyLookup: &keyStore{},
	}

	auth, err := auth.New(cfg)
	if err != nil {
		tb.Fatal(err)
	}

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		tb.Helper()
		testDB.TearDown()
		fmt.Println("******************** LOGS ********************")
		fmt.Print(buf.String())
		fmt.Println("******************** LOGS ********************")
	}

	test := Test{
		TestDatabase: *testDB,
		Log:          log,
		CoreAPIs:     coreAPIs,
		TearDown:     teardown,
		tb:           tb,
		Auth:         auth,
	}
	return &test
}

// ====================================================================
// CoreAPIs represents all the core api's needed for testing.
type CoreAPIs struct {
	User *user.Core
}

func newCoreAPIs(log *logger.Logger, db *sqlx.DB) CoreAPIs {
	usrCore := user.NewCore(log, userdb.NewRepository(log, db))
	return CoreAPIs{
		User: usrCore,
	}
}

// ============================================================

func (test *Test) TokenV1(email, password string) (string, error) {
	test.tb.Logf("Generating %q token for test ...", email)

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return "", err
	}

	store := userdb.NewRepository(test.Log, test.DB)
	dbUsr, err := store.QueryByEmail(context.Background(), *addr)
	if err != nil {
		return "", err
	}

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   dbUsr.ID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: dbUsr.Roles,
	}

	token, err := test.Auth.GenerateToken(kid, claims)

	if err != nil {
		test.tb.Fatal(err)
	}

	return token, nil

}

// ==================================================================

type keyStore struct{}

func (ks *keyStore) PrivateKey(kid string) (string, error) {
	return privateKeyPEM, nil
}

func (ks *keyStore) PublicKey(kid string) (string, error) {
	return publicKeyPEM, nil
}

const (
	kid = "s4sKIjD9kIRjxs2tulPqGLdxSfgPErRN1Mu3Hd9k9NQ"

	privateKeyPEM = `-----BEGIN PRIVATE KEY-----
MIIEpQIBAAKCAQEAvMAHb0IoLvoYuW2kA+LTmnk+hfnBq1eYIh4CT/rMPCxgtzjq
U0guQOMnLg69ydyA5uu37v6rbS1+stuBTEiMQl/bxAhgLkGrUhgpZ10Bt6GzSEgw
QNloZoGaxe4p20wMPpT4kcMKNHkQds3uONNcLxPUmfjbbH64g+seg28pbgQPwKFK
tF7bIsOBgz0g5Ptn5mrkdzqMPUSy9k9VCu+R42LH9c75JsRzz4FeN+VzwMAL6yQn
ZvOi7/zOgNyxeVia8XVKykrnhgcpiOn5oaLRBzQGN00Z7TuBRIfDJWU21qQN4Cq7
keZmMP4gqCVWjYneK4bzrG/+H2w9BJ2TsmMGvwIDAQABAoIBAFQmQKpHkmavNYql
6POaksBRwaA1YzSijr7XJizGIXvKRSwqgb2zdnuTSgpspAx09Dr/aDdy7rZ0DAJt
fk2mInINDottOIQm3txwzTS58GQQAT/+fxTKWJMqwPfxYFPWqbbU76T8kXYna0Gs
OcK36GdMrgIfQqQyMs0Na8MpMg1LmkAxuqnFCXS/NMyKl9jInaaTS+Kz+BSzUMGQ
zebfLFsf2N7sLZuimt9zlRG30JJTfBlB04xsYMo734usA2ITe8U0XqG6Og0qc6ev
6lsoM8hpvEUsQLcjQQ5up7xx3S2stZJ8o0X8GEX5qUMaomil8mZ7X5xOlEqf7p+v
lXQ46cECgYEA2lbZQON6l3ZV9PCn9j1rEGaXio3SrAdTyWK3D1HF+/lEjClhMkfC
XrECOZYj+fiI9n+YpSog+tTDF7FTLf7VP21d2gnhQN6KAXUnLIypzXxodcC6h+8M
ZGJh/EydLvC7nPNoaXx96bohxzS8hrOlOlkCbr+8gPYKf8qkbe7HyxECgYEA3U6e
x9g4FfTvI5MGrhp2BIzoRSn7HlNQzjJ71iMHmM2kBm7TsER8Co1PmPDrP8K/UyGU
Q25usTsPSrHtKQEV6EsWKaP/6p2Q82sDkT9bZlV+OjRvOfpdO5rP6Q95vUmMGWJ/
S6oimbXXL8p3gDafw3vC1PCAhoaxMnGyKuZwlM8CgYEAixT1sXr2dZMg8DV4mMfI
8pqXf+AVyhWkzsz+FVkeyAKiIrKdQp0peI5C/5HfevVRscvX3aY3efCcEfSYKt2A
07WEKkdO4LahrIoHGT7FT6snE5NgfwTMnQl6p2/aVLNun20CHuf5gTBbIf069odr
Af7/KLMkjfWs/HiGQ6zuQjECgYEAv+DIvlDz3+Wr6dYyNoXuyWc6g60wc0ydhQo0
YKeikJPLoWA53lyih6uZ1escrP23UOaOXCDFjJi+W28FR0YProZbwuLUoqDW6pZg
U3DxWDrL5L9NqKEwcNt7ZIDsdnfsJp5F7F6o/UiyOFd9YQb7YkxN0r5rUTg7Lpdx
eMyv0/UCgYEAhX9MPzmTO4+N8naGFof1o8YP97pZj0HkEvM0hTaeAQFKJiwX5ijQ
xumKGh//G0AYsjqP02ItzOm2mWnbI3FrNlKmGFvR6VxIZMOyXvpLofHucjJ5SWli
eYjPklKcXaMftt1FVO4n+EKj1k1+Tv14nytq/J5WN+r4FBlNEYj/6vg=
-----END PRIVATE KEY-----
`
	publicKeyPEM = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvMAHb0IoLvoYuW2kA+LT
mnk+hfnBq1eYIh4CT/rMPCxgtzjqU0guQOMnLg69ydyA5uu37v6rbS1+stuBTEiM
Ql/bxAhgLkGrUhgpZ10Bt6GzSEgwQNloZoGaxe4p20wMPpT4kcMKNHkQds3uONNc
LxPUmfjbbH64g+seg28pbgQPwKFKtF7bIsOBgz0g5Ptn5mrkdzqMPUSy9k9VCu+R
42LH9c75JsRzz4FeN+VzwMAL6yQnZvOi7/zOgNyxeVia8XVKykrnhgcpiOn5oaLR
BzQGN00Z7TuBRIfDJWU21qQN4Cq7keZmMP4gqCVWjYneK4bzrG/+H2w9BJ2TsmMG
vwIDAQAB
-----END PUBLIC KEY-----`
)
