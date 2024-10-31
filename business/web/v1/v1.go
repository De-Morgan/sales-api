package v1

import (
	"os"
	"sales-api/business/web/v1/auth"
	"sales-api/business/web/v1/mid"
	"sales-api/foundation/logger"
	"sales-api/foundation/web"

	"github.com/jmoiron/sqlx"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Build    string
	Shutdown chan os.Signal
	Log      *logger.Logger
	Auth     *auth.Auth
	DB       *sqlx.DB
}

type RouteAdder interface {
	Add(*web.App, APIMuxConfig)
}

func APIMux(cfg APIMuxConfig, routeAdder RouteAdder) *web.App {
	const version = "/v1"
	app := web.NewApp(cfg.Shutdown, version, mid.Logger(cfg.Log), mid.Errors(cfg.Log), mid.Metrics(), mid.Panics())
	routeAdder.Add(app, cfg)
	return app
}
