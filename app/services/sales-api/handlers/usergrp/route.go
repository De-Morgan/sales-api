package usergrp

import (
	"sales-api/business/core/user"
	"sales-api/business/core/user/stores/userdb"
	"sales-api/business/data/dbsql/pgx"
	"sales-api/business/web/v1/auth"
	"sales-api/business/web/v1/mid"
	"sales-api/foundation/logger"
	"sales-api/foundation/web"

	"github.com/jmoiron/sqlx"
)

type Config struct {
	Build string
	Log   *logger.Logger
	DB    *sqlx.DB
	Auth  *auth.Auth
}

func Route(app *web.App, cfg Config) {

	usrCore := user.NewCore(cfg.Log, userdb.NewRepository(cfg.Log, cfg.DB))

	authMid := mid.Authenticate(cfg.Auth)
	ruleAdmin := mid.Authorize(cfg.Auth, auth.RuleAdminOnly)
	ruleAdminOrSubject := mid.Authorize(cfg.Auth, auth.RuleAdminOrSubject)

	tran := mid.ExecuteInTransaction(cfg.Log, pgx.NewBeginner(cfg.DB))

	hdl := New(usrCore, cfg.Auth)
	// POST===========================================================================
	app.HandleFunc("/users", hdl.Create).Methods("POST")
	app.HandleFunc("/users/login", hdl.Login).Methods("POST")

	// PUT===========================================================================
	app.HandleFunc("/users/{user_id}", hdl.UpdateByID, authMid, ruleAdminOrSubject, tran).Methods("PUT")

	// GET===========================================================================

	app.HandleFunc("/users/{user_id}", hdl.QueryByID, authMid, ruleAdminOrSubject).Methods("GET")
	app.HandleFunc("/users", hdl.Query, authMid, ruleAdmin).Methods("GET")

	// DELETE===========================================================================
	app.HandleFunc("/users/{user_id}", hdl.DeleteByID, authMid, ruleAdmin).Methods("DELETE")

}
