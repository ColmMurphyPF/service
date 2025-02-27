// Package db contains product related CRUD functionality.
package db

import (
	"context"
	"fmt"
	"github.com/colmmurphy91/go-service/business/sys/database/sql"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Store manages the set of APIs for user access.
type Store struct {
	log          *zap.SugaredLogger
	tr           sql.Transactor
	db           sqlx.ExtContext
	isWithinTran bool
}

// NewStore constructs a data for api access.
func NewStore(log *zap.SugaredLogger, db *sqlx.DB) Store {
	return Store{
		log: log,
		tr:  db,
		db:  db,
	}
}

// WithinTran runs passed function and do commit/rollback at the end.
func (s Store) WithinTran(ctx context.Context, fn func(sqlx.ExtContext) error) error {
	if s.isWithinTran {
		fn(s.db)
	}
	return sql.WithinTran(ctx, s.log, s.tr, fn)
}

// Tran return new Store with transaction in it.
func (s Store) Tran(tx sqlx.ExtContext) Store {
	return Store{
		log:          s.log,
		tr:           s.tr,
		db:           tx,
		isWithinTran: true,
	}
}

func (s Store) Create(ctx context.Context, caf Cafe) error {
	const q = `
	INSERT INTO cafes
		(cafe_id, owner_id, cafe_name, address, logo_url)
	VALUES
		(:cafe_id, :owner_id, :cafe_name, :address, :logo_url)`

	if err := sql.NamedExecContext(ctx, s.log, s.db, q, caf); err != nil {
		return fmt.Errorf("inserting cafe: #{err}")
	}
	return nil
}

func (s Store) QueryByOwnerID(ctx context.Context, ownerID string) (Cafe, error) {
	data := struct {
		OwnerID string `db:"owner_id"`
	}{
		OwnerID: ownerID,
	}

	const q = `
	SELECT 
		* 
	FROM 
		cafes 
	WHERE 
		owner_id=:owner_id`

	var cafe Cafe
	if err := sql.NamedQueryStruct(ctx, s.log, s.db, q, data, &cafe); err != nil {
		return Cafe{}, fmt.Errorf("selecting ownerid[%q]: %w", ownerID, err)
	}
	return cafe, nil
}

//// Update modifies data about a Product. It will error if the specified ID is
//// invalid or does not reference an existing Product.
//func (s Store) Update(ctx context.Context, prd Product) error {
//	const q = `
//	UPDATE
//		products
//	SET
//		"name" = :name,
//		"cost" = :cost,
//		"quantity" = :quantity,
//		"date_updated" = :date_updated
//	WHERE
//		product_id = :product_id`
//
//	if err := sql.NamedExecContext(ctx, s.log, s.db, q, prd); err != nil {
//		return fmt.Errorf("updating product productID[%s]: %w", prd.ID, err)
//	}
//
//	return nil
//}
//
//// Delete removes the product identified by a given ID.
//func (s Store) Delete(ctx context.Context, productID string) error {
//	data := struct {
//		ProductID string `db:"product_id"`
//	}{
//		ProductID: productID,
//	}
//
//	const q = `
//	DELETE FROM
//		products
//	WHERE
//		product_id = :product_id`
//
//	if err := sql.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
//		return fmt.Errorf("deleting product productID[%s]: %w", productID, err)
//	}
//
//	return nil
//}
//
//// Query gets all Products from the database.
//func (s Store) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]Product, error) {
//	data := struct {
//		Offset      int `db:"offset"`
//		RowsPerPage int `db:"rows_per_page"`
//	}{
//		Offset:      (pageNumber - 1) * rowsPerPage,
//		RowsPerPage: rowsPerPage,
//	}
//
//	const q = `
//	SELECT
//		p.*,
//		COALESCE(SUM(s.quantity) ,0) AS sold,
//		COALESCE(SUM(s.paid), 0) AS revenue
//	FROM
//		products AS p
//	LEFT JOIN
//		sales AS s ON p.product_id = s.product_id
//	GROUP BY
//		p.product_id
//	ORDER BY
//		user_id
//	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`
//
//	var prds []Product
//	if err := sql.NamedQuerySlice(ctx, s.log, s.db, q, data, &prds); err != nil {
//		return nil, fmt.Errorf("selecting products: %w", err)
//	}
//
//	return prds, nil
//}
//
//// QueryByID finds the product identified by a given ID.
//func (s Store) QueryByID(ctx context.Context, productID string) (Product, error) {
//	data := struct {
//		ProductID string `db:"product_id"`
//	}{
//		ProductID: productID,
//	}
//
//	const q = `
//	SELECT
//		p.*,
//		COALESCE(SUM(s.quantity), 0) AS sold,
//		COALESCE(SUM(s.paid), 0) AS revenue
//	FROM
//		products AS p
//	LEFT JOIN
//		sales AS s ON p.product_id = s.product_id
//	WHERE
//		p.product_id = :product_id
//	GROUP BY
//		p.product_id`
//
//	var prd Product
//	if err := sql.NamedQueryStruct(ctx, s.log, s.db, q, data, &prd); err != nil {
//		return Product{}, fmt.Errorf("selecting product productID[%q]: %w", productID, err)
//	}
//
//	return prd, nil
//}
//
//// QueryByUserID finds the product identified by a given User ID.
//func (s Store) QueryByUserID(ctx context.Context, userID string) ([]Product, error) {
//	data := struct {
//		UserID string `db:"user_id"`
//	}{
//		UserID: userID,
//	}
//
//	const q = `
//	SELECT
//		p.*,
//		COALESCE(SUM(s.quantity), 0) AS sold,
//		COALESCE(SUM(s.paid), 0) AS revenue
//	FROM
//		products AS p
//	LEFT JOIN
//		sales AS s ON p.product_id = s.product_id
//	WHERE
//		p.user_id = :user_id
//	GROUP BY
//		p.product_id`
//
//	var prds []Product
//	if err := sql.NamedQuerySlice(ctx, s.log, s.db, q, data, &prds); err != nil {
//		return nil, fmt.Errorf("selecting products userID[%s]: %w", userID, err)
//	}
//
//	return prds, nil
//}
