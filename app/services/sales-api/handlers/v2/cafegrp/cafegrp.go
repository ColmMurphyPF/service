package cafegrp

import (
	"context"
	"errors"
	"fmt"
	"github.com/colmmurphy91/go-service/business/core/cafev2"
	"github.com/colmmurphy91/go-service/business/sys/auth"
	v1Web "github.com/colmmurphy91/go-service/business/web/v1"
	"github.com/colmmurphy91/go-service/foundation/web"
	"net/http"
)

type Handlers struct {
	Cafe cafev2.Core
}

func (h Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var newCafe cafev2.NewCafe
	err := web.Decode(r, &newCafe)
	if err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}
	claims, err := auth.GetClaims(ctx)

	createCafe, err := h.Cafe.Create(ctx, newCafe, claims.Subject)
	if err != nil {
		return fmt.Errorf("creating new cafe, nc[%+v]: %w", newCafe, err)
	}
	return web.Respond(ctx, w, createCafe, http.StatusCreated)
}

func (h Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := web.Param(r, "id")
	claims, err := auth.GetClaims(ctx)
	caf, err := h.Cafe.QueryByOwnerID(ctx, claims.Subject)
	if claims.Subject != caf.OwnerID {
		return v1Web.NewRequestError(errors.New("user is not part of cafe"), http.StatusForbidden)
	}
	if err != nil {
		switch {
		case errors.Is(err, cafev2.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s]: %w", id, err)
		}
	}

	return web.Respond(ctx, w, caf, http.StatusOK)
}

func (h Handlers) Hello(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return web.Respond(ctx, w, "{'status':'ok'}", http.StatusOK)
}
