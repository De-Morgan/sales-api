package usergrp

import (
	"sales-api/business/core/user"
	"sales-api/business/core/user/stores/userdb"
	"sales-api/business/web/v1/auth"
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
	hdl := New(usrCore, cfg.Auth)
	app.HandleFunc("/users", hdl.Create).Methods("POST")
}
