package commands

import (
	"context"
	"fmt"
	"github.com/colmmurphy91/go-service/business/sys/database/sql"
	"time"

	"github.com/colmmurphy91/go-service/business/data/dbschema"
)

// Seed loads test data into the database.
func Seed(cfg sql.Config) error {
	db, err := sql.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := dbschema.Seed(ctx, db); err != nil {
		return fmt.Errorf("seed database: %w", err)
	}

	fmt.Println("seed data complete")
	return nil
}
