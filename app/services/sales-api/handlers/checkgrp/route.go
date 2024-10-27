package checkgrp

import (
	"sales-api/foundation/logger"
	"sales-api/foundation/web"
)

type Config struct {
	Build  string
	Logger *logger.Logger
}

func Route(app *web.App, cfg Config) {
	hdl := New(cfg.Build, cfg.Logger)

	app.HandleNoMiddleWareFunc("/readiness", hdl.readiness)
	app.HandleNoMiddleWareFunc("/liveness", hdl.liviness)

}
