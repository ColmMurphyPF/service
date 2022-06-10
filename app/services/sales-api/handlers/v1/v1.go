// Package v1 contains the full set of handler functions and routes
// supported by the v1 web api.
package v1

import (
	"github.com/colmmurphy91/go-service/app/services/sales-api/handlers/v1/cafegrp"
	"github.com/colmmurphy91/go-service/business/core/cafe"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"

	"github.com/colmmurphy91/go-service/app/services/sales-api/handlers/v1/productgrp"
	"github.com/colmmurphy91/go-service/app/services/sales-api/handlers/v1/usergrp"
	"github.com/colmmurphy91/go-service/business/core/product"
	"github.com/colmmurphy91/go-service/business/core/user"
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
	MDB  *mongo.Database
}

// Routes binds all the version 1 routes.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.Auth)
	admin := mid.Authorize(auth.RoleAdmin)

	// Register user management and authentication endpoints.
	ugh := usergrp.Handlers{
		User: user.NewCore(cfg.Log, cfg.DB),
		Auth: cfg.Auth,
	}
	app.Handle(http.MethodPut, version, "/users/confirm", ugh.Confirm)
	app.Handle(http.MethodGet, version, "/users/token", ugh.Token)
	app.Handle(http.MethodGet, version, "/users/:page/:rows", ugh.Query, authen, admin)
	app.Handle(http.MethodGet, version, "/users/:id", ugh.QueryByID, authen)
	app.Handle(http.MethodPost, version, "/users", ugh.Create, authen, admin)
	app.Handle(http.MethodPut, version, "/users/:id", ugh.Update, authen, admin)
	app.Handle(http.MethodDelete, version, "/users/:id", ugh.Delete, authen, admin)

	// Register product and sale endpoints.
	pgh := productgrp.Handlers{
		Product: product.NewCore(cfg.Log, cfg.DB),
	}
	app.Handle(http.MethodGet, version, "/products/:page/:rows", pgh.Query, authen)
	app.Handle(http.MethodGet, version, "/products/:id", pgh.QueryByID, authen)
	app.Handle(http.MethodPost, version, "/products", pgh.Create, authen)
	app.Handle(http.MethodPut, version, "/products/:id", pgh.Update, authen)
	app.Handle(http.MethodDelete, version, "/products/:id", pgh.Delete, authen)

	cgh := cafegrp.Handlers{
		Cafe: cafe.NewCore(cfg.Log, cfg.MDB),
	}

	app.Handle(http.MethodPost, version, "/cafes", cgh.Create, authen)
	app.Handle(http.MethodGet, version, "/cafes/:id", cgh.QueryByID, authen)
	app.Handle(http.MethodGet, version, "/cafes", cgh.GetAll, authen)
	app.Handle(http.MethodDelete, version, "/cafes/:id", cgh.DeleteByID, authen)
	app.Handle(http.MethodPut, version, "/cafes/:id", cgh.UpdateByID, authen)
}
