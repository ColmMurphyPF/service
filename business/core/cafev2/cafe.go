// Package product provides an example of a core business API. Right now these
// calls are just wrapping the data/store layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package cafev2

import (
	"context"
	"errors"
	"fmt"
	"github.com/colmmurphy91/go-service/business/core/cafev2/db"
	"github.com/colmmurphy91/go-service/business/sys/validate"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound  = errors.New("cafe not found")
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// Core manages the set of APIs for product access.
type Core struct {
	store db.Store
}

// NewCore constructs a core for product api access.
func NewCore(log *zap.SugaredLogger, sqlxDB *sqlx.DB) Core {
	return Core{
		store: db.NewStore(log, sqlxDB),
	}
}

// Create adds a Product to the database. It returns the created Product with
// fields like ID and DateCreated populated.
func (c Core) Create(ctx context.Context, np NewCafe, ownerID string) (Cafe, error) {
	if err := validate.Check(np); err != nil {
		return Cafe{}, fmt.Errorf("validating data: %w", err)
	}

	_, err := c.store.QueryByOwnerID(ctx, ownerID)

	if err == nil {
		return Cafe{}, errors.New("owner already has cafe created")
	}

	dbCaf := db.Cafe{
		ID:      validate.GenerateID(),
		Name:    np.Name,
		Address: np.Address,
		LogoURL: np.LogoURL,
		OwnerID: ownerID,
	}

	if err := c.store.Create(ctx, dbCaf); err != nil {
		return Cafe{}, fmt.Errorf("create: %w", err)
	}

	return toCafe(dbCaf), nil
}

//// Update modifies data about a Product. It will error if the specified ID is
//// invalid or does not reference an existing Product.
//func (c Core) Update(ctx context.Context, productID string, up UpdateProduct, now time.Time) error {
//	if err := validate.CheckID(productID); err != nil {
//		return ErrInvalidID
//	}
//
//	if err := validate.Check(up); err != nil {
//		return fmt.Errorf("validating data: %w", err)
//	}
//
//	dbPrd, err := c.store.QueryByID(ctx, productID)
//	if err != nil {
//		if errors.Is(err, sql.ErrDBNotFound) {
//			return ErrNotFound
//		}
//		return fmt.Errorf("updating product productID[%s]: %w", productID, err)
//	}
//
//	if up.Name != nil {
//		dbPrd.Name = *up.Name
//	}
//	if up.Cost != nil {
//		dbPrd.Cost = *up.Cost
//	}
//	if up.Quantity != nil {
//		dbPrd.Quantity = *up.Quantity
//	}
//	dbPrd.DateUpdated = now
//
//	if err := c.store.Update(ctx, dbPrd); err != nil {
//		return fmt.Errorf("update: %w", err)
//	}
//
//	return nil
//}
//
//// Delete removes the product identified by a given ID.
//func (c Core) Delete(ctx context.Context, productID string) error {
//	if err := validate.CheckID(productID); err != nil {
//		return ErrInvalidID
//	}
//
//	if err := c.store.Delete(ctx, productID); err != nil {
//		return fmt.Errorf("delete: %w", err)
//	}
//
//	return nil
//}
//
//// Query gets all Products from the database.
//func (c Core) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]Product, error) {
//	dbPrds, err := c.store.Query(ctx, pageNumber, rowsPerPage)
//	if err != nil {
//		return nil, fmt.Errorf("query: %w", err)
//	}
//
//	return toProductSlice(dbPrds), nil
//}
//
//// QueryByID finds the product identified by a given ID.
//func (c Core) QueryByID(ctx context.Context, productID string) (Product, error) {
//	if err := validate.CheckID(productID); err != nil {
//		return Product{}, ErrInvalidID
//	}
//
//	dbPrd, err := c.store.QueryByID(ctx, productID)
//	if err != nil {
//		if errors.Is(err, sql.ErrDBNotFound) {
//			return Product{}, ErrNotFound
//		}
//		return Product{}, fmt.Errorf("query: %w", err)
//	}
//
//	return toProduct(dbPrd), nil
//}
//
// QueryByUserID finds the products identified by a given User ID.
func (c Core) QueryByOwnerID(ctx context.Context, ownerID string) (Cafe, error) {
	if err := validate.CheckID(ownerID); err != nil {
		return Cafe{}, ErrInvalidID
	}

	dbCafe, err := c.store.QueryByOwnerID(ctx, ownerID)
	if err != nil {
		return Cafe{}, fmt.Errorf("query: %w", err)
	}

	return toCafe(dbCafe), nil
}
