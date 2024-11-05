package checkgrp

import (
	"sales-api/foundation/logger"
	"sales-api/foundation/web"

	"github.com/jmoiron/sqlx"
)

type Config struct {
	Build  string
	Logger *logger.Logger
	DB     *sqlx.DB
}

func Route(app *web.App, cfg Config) {
	hdl := New(cfg.Build, cfg.Logger, cfg.DB)

	app.HandleNoMiddleWareFunc("/readiness", hdl.readiness)
	app.HandleNoMiddleWareFunc("/liveness", hdl.liviness)

}
