package handlers

import (
	"sales-api/app/services/sales-api/handlers/checkgrp"
	"sales-api/app/services/sales-api/handlers/usergrp"
	v1 "sales-api/business/web/v1"
	"sales-api/foundation/web"
)

func Routes() *add {
	return &add{}
}

var _ v1.RouteAdder = (*add)(nil)

type add struct{}

func (a *add) Add(app *web.App, cfg v1.APIMuxConfig) {

	checkgrp.Route(app, checkgrp.Config{Build: cfg.Build, Logger: cfg.Log, DB: cfg.DB})
	usergrp.Route(app, usergrp.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
		Auth:  cfg.Auth,
	})
}
