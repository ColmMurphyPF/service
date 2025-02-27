package commands

import (
	"context"
	"fmt"
	"github.com/colmmurphy91/go-service/business/sys/database/sql"
	"time"

	"github.com/colmmurphy91/go-service/business/core/user"
	"github.com/colmmurphy91/go-service/business/sys/auth"
	"go.uber.org/zap"
)

// UserAdd adds new users into the database.
func UserAdd(log *zap.SugaredLogger, cfg sql.Config, name, email, password string) error {
	if name == "" || email == "" || password == "" {
		fmt.Println("help: useradd <name> <email> <password>")
		return ErrHelp
	}

	db, err := sql.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	core := user.NewCore(log, db)

	nu := user.NewUser{
		Name:            name,
		Email:           email,
		Password:        password,
		PasswordConfirm: password,
		Roles:           []string{auth.RoleAdmin, auth.RoleUser},
	}

	usr, err := core.Create(ctx, nu, time.Now())
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	fmt.Println("user id:", usr.ID)
	return nil
}
