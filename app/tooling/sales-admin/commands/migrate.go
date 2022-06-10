package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/colmmurphy91/go-service/business/sys/database/sql"
	"time"

	"github.com/colmmurphy91/go-service/business/data/dbschema"
)

// ErrHelp provides context that help was given.
var ErrHelp = errors.New("provided help")

// Migrate creates the schema in the database.
func Migrate(cfg sql.Config) error {
	db, err := sql.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := dbschema.Migrate(ctx, db); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	fmt.Println("migrations complete")
	return nil
}
