// Package v1 contains the full set of handler functions and routes
// supported by the v1 web api.
package v1

import (
	"github.com/colmmurphy91/go-service/app/services/sales-api/handlers/v2/cafegrp"
	"github.com/colmmurphy91/go-service/business/core/cafev2"
	"net/http"

	"github.com/colmmurphy91/go-service/business/sys/auth"
	"github.com/colmmurphy91/go-service/business/web/v1/mid"
	"github.com/colmmurphy91/go-service/foundation/web"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log  *zap.SugaredLogger
	Auth *auth.Auth
	DB   *sqlx.DB
}

// Routes binds all the version 2 routes.
func Routes(app *web.App, cfg Config) {
	const version = "v2"

	authen := mid.Authenticate(cfg.Auth)

	cgh := cafegrp.Handlers{
		Cafe: cafev2.NewCore(cfg.Log, cfg.DB),
	}

	app.Handle(http.MethodPost, version, "/cafes", cgh.Create, authen)
	app.Handle(http.MethodGet, version, "/cafes/:id", cgh.QueryByID, authen)
	app.Handle(http.MethodGet, version, "/hello", cgh.Hello)
}
