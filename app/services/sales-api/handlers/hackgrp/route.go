package hackgrp

import (
	"sales-api/business/web/v1/auth"
	"sales-api/business/web/v1/mid"
	"sales-api/foundation/web"
)

type Config struct {
	Auth *auth.Auth
}

func Route(app *web.App, cfg Config) {

	authen := mid.Authenticate(cfg.Auth)
	roleAdmin := mid.Authorize(cfg.Auth, auth.RuleAdminOnly)

	app.HandleFunc("/hack", hack)
	app.HandleFunc("/hack/auth", hack, authen)
	app.HandleFunc("/hack/auth/admin", hack, authen, roleAdmin)

}
