package cafegrp

import (
	"context"
	"errors"
	"fmt"
	"github.com/colmmurphy91/go-service/business/core/cafe"
	v1Web "github.com/colmmurphy91/go-service/business/web/v1"
	"github.com/colmmurphy91/go-service/foundation/web"
	"net/http"
)

type Handlers struct {
	Cafe cafe.Core
}

func (h Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var newCafe cafe.NewCafe
	err := web.Decode(r, &newCafe)
	if err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	createCafe, err := h.Cafe.CreateCafe(newCafe)
	if err != nil {
		return fmt.Errorf("creating new cafe, nc[%+v]: %w", newCafe, err)
	}
	return web.Respond(ctx, w, createCafe, http.StatusCreated)
}

func (h Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := web.Param(r, "id")
	caf, err := h.Cafe.FindCafe(id)
	if err != nil {
		switch {
		case errors.Is(err, cafe.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s]: %w", id, err)
		}
	}

	return web.Respond(ctx, w, caf, http.StatusOK)
}

func (h Handlers) DeleteByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := web.Param(r, "id")
	err := h.Cafe.DeleteCafe(id)
	fmt.Println("here")
	if err != nil {
		fmt.Println(err.Error())
		switch {
		case errors.Is(err, cafe.ErrDeletion):
			fmt.Println("here 3")
			return web.Respond(ctx, w, nil, http.StatusNoContent)
		default:
			fmt.Println("here 4")
			return fmt.Errorf("ID[%s]: %w", id, err)
		}
	}
	fmt.Println("here 3")
	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (h Handlers) UpdateByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := web.Param(r, "id")
	var uc cafe.UpdateCafe
	err := web.Decode(r, &uc)
	if err != nil {
		return err
	}
	uc.ID = id
	err = h.Cafe.UpdateCafe(uc)
	if err != nil {
		switch {
		case errors.Is(err, cafe.ErrNotFound):
			return v1Web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("ID[%s]: %w", id, err)
		}
	}
	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (h Handlers) GetAll(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	allCafes, _ := h.Cafe.FindAll()
	return web.Respond(ctx, w, allCafes, http.StatusOK)
}
