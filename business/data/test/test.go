package test

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"sales-api/business/core/user"
	"sales-api/business/core/user/stores/userdb"
	"sales-api/foundation/logger"
	"sales-api/foundation/web"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
)

// Test owns state for running and shutting down tests.
type Test struct {
	TestDatabase
	Log      *logger.Logger
	CoreAPIs CoreAPIs
	TearDown func()
	tb       testing.TB
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
