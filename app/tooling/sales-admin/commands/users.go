package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/colmmurphy91/go-service/business/sys/database/sql"
	"os"
	"strconv"
	"time"

	"github.com/colmmurphy91/go-service/business/core/user"
	"go.uber.org/zap"
)

// Users retrieves all users from the database.
func Users(log *zap.SugaredLogger, cfg sql.Config, pageNumber string, rowsPerPage string) error {
	db, err := sql.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	page, err := strconv.Atoi(pageNumber)
	if err != nil {
		return fmt.Errorf("converting page number: %w", err)
	}

	rows, err := strconv.Atoi(rowsPerPage)
	if err != nil {
		return fmt.Errorf("converting rows per page: %w", err)
	}

	user := user.NewCore(log, db)

	users, err := user.Query(ctx, page, rows)
	if err != nil {
		return fmt.Errorf("retrieve users: %w", err)
	}

	return json.NewEncoder(os.Stdout).Encode(users)
}
