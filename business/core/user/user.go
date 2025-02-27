// Package user provides an example of a core business API. Right now these
// calls are just wrapping the data/data layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package user

import (
	"context"
	crand "crypto/rand"
	"database/sql"
	"encoding/binary"
	"errors"
	"fmt"
	csql "github.com/colmmurphy91/go-service/business/sys/database/sql"
	"log"
	"math/rand"
	"time"

	"github.com/colmmurphy91/go-service/business/core/user/db"
	"github.com/colmmurphy91/go-service/business/sys/auth"
	"github.com/colmmurphy91/go-service/business/sys/validate"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound              = errors.New("user not found")
	ErrInvalidID             = errors.New("ID is not in its proper form")
	ErrInvalidEmail          = errors.New("email is not valid")
	ErrUniqueEmail           = errors.New("email is not unique")
	ErrAuthenticationFailure = errors.New("authentication failed")
	ErrUserNotConfirmed      = errors.New("user is not confirmed")
	ErrAlreadyConfirmed      = errors.New("user is already confirmed")
)

// Core manages the set of APIs for user access.
type Core struct {
	store db.Store
	log   *zap.SugaredLogger
}

// NewCore constructs a core for user api access.
func NewCore(log *zap.SugaredLogger, sqlxDB *sqlx.DB) Core {
	return Core{
		store: db.NewStore(log, sqlxDB),
		log:   log,
	}
}

// Create inserts a new user into the database.
func (c Core) Create(ctx context.Context, nu NewUser, now time.Time) (User, error) {
	if err := validate.Check(nu); err != nil {
		return User{}, fmt.Errorf("validating data: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generating password hash: %w", err)
	}

	var src cryptoSource
	rnd := rand.New(src)

	myNum := rnd.Intn(100000)

	toSave := sql.NullInt64{
		Int64: int64(myNum),
		Valid: true,
	}

	dbUsr := db.User{
		ID:           validate.GenerateID(),
		Name:         nu.Name,
		Email:        nu.Email,
		PasswordHash: hash,
		Roles:        nu.Roles,
		DateCreated:  now,
		DateUpdated:  now,
		Confirmed:    false,
		ConfirmHash:  toSave,
	}

	// This provides an example of how to execute a transaction if required.
	tran := func(tx sqlx.ExtContext) error {
		if err := c.store.Tran(tx).Create(ctx, dbUsr); err != nil {
			if errors.Is(err, csql.ErrDBDuplicatedEntry) {
				return fmt.Errorf("create: %w", ErrUniqueEmail)
			}
			return fmt.Errorf("create: %w", err)
		}
		return nil
	}

	if err := c.store.WithinTran(ctx, tran); err != nil {
		return User{}, fmt.Errorf("tran: %w", err)
	}

	if err != nil {
		return User{}, err
	}

	return toUser(dbUsr), nil
}

// Update replaces a user document in the database.
func (c Core) Update(ctx context.Context, userID string, uu UpdateUser, now time.Time) error {
	if err := validate.CheckID(userID); err != nil {
		return ErrInvalidID
	}

	if err := validate.Check(uu); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	dbUsr, err := c.store.QueryByID(ctx, userID)
	if err != nil {
		if errors.Is(err, csql.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("updating user userID[%s]: %w", userID, err)
	}

	if uu.Name != nil {
		dbUsr.Name = *uu.Name
	}
	if uu.Email != nil {
		dbUsr.Email = *uu.Email
	}
	if uu.Roles != nil {
		dbUsr.Roles = uu.Roles
	}
	if uu.Password != nil {
		pw, err := bcrypt.GenerateFromPassword([]byte(*uu.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("generating password hash: %w", err)
		}
		dbUsr.PasswordHash = pw
	}
	if uu.ConfirmHash == nil {
		dbUsr.ConfirmHash = sql.NullInt64{
			Valid: false,
		}
	}

	if uu.Confirmed != nil {
		dbUsr.Confirmed = true
	}

	dbUsr.DateUpdated = now

	if err := c.store.Update(ctx, dbUsr); err != nil {
		if errors.Is(err, csql.ErrDBDuplicatedEntry) {
			return fmt.Errorf("updating user userID[%s]: %w", userID, ErrUniqueEmail)
		}
		return fmt.Errorf("update: %w", err)
	}

	return nil
}

// Delete removes a user from the database.
func (c Core) Delete(ctx context.Context, userID string) error {
	if err := validate.CheckID(userID); err != nil {
		return ErrInvalidID
	}

	if err := c.store.Delete(ctx, userID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query retrieves a list of existing users from the database.
func (c Core) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]User, error) {
	dbUsers, err := c.store.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return toUserSlice(dbUsers), nil
}

// QueryByID gets the specified user from the database.
func (c Core) QueryByID(ctx context.Context, userID string) (User, error) {
	if err := validate.CheckID(userID); err != nil {
		return User{}, ErrInvalidID
	}

	dbUsr, err := c.store.QueryByID(ctx, userID)
	if err != nil {
		if errors.Is(err, csql.ErrDBNotFound) {
			return User{}, ErrNotFound
		}
		return User{}, fmt.Errorf("query: %w", err)
	}

	return toUser(dbUsr), nil
}

// QueryByEmail gets the specified user from the database by email.
func (c Core) QueryByEmail(ctx context.Context, email string) (User, error) {

	// Email Validate function in validate.
	if !validate.CheckEmail(email) {
		return User{}, ErrInvalidEmail
	}

	dbUsr, err := c.store.QueryByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, csql.ErrDBNotFound) {
			return User{}, ErrNotFound
		}
		return User{}, fmt.Errorf("query: %w", err)
	}

	return toUser(dbUsr), nil
}

// Authenticate finds a user by their email and verifies their password. On
// success it returns a Claims User representing this user. The claims can be
// used to generate a token for future authentication.
func (c Core) Authenticate(ctx context.Context, now time.Time, email, password string) (auth.Claims, error) {
	dbUsr, err := c.store.QueryByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, csql.ErrDBNotFound) {
			return auth.Claims{}, ErrNotFound
		}
		return auth.Claims{}, fmt.Errorf("query: %w", err)
	}

	if !dbUsr.Confirmed {
		return auth.Claims{}, ErrUserNotConfirmed
	}

	// Compare the provided password with the saved hash. Use the bcrypt
	// comparison function so it is cryptographically secure.
	if err := bcrypt.CompareHashAndPassword(dbUsr.PasswordHash, []byte(password)); err != nil {
		return auth.Claims{}, ErrAuthenticationFailure
	}

	// If we are this far the request is valid. Create some claims for the user
	// and generate their token.
	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   dbUsr.ID,
			Issuer:    "service project",
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: dbUsr.Roles,
	}

	return claims, nil
}

func (c Core) Confirm(ctx context.Context, email string, token int64) error {
	dbUser, err := c.store.QueryByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, csql.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("query: %w", err)
	}

	if dbUser.Confirmed {
		return ErrAlreadyConfirmed
	}

	if dbUser.ConfirmHash.Int64 != token {
		return errors.New("user token is not equal")
	}

	if dbUser.ConfirmHash.Int64 == token {
		t := true
		err := c.Update(ctx, dbUser.ID, UpdateUser{
			Confirmed:   &t,
			ConfirmHash: nil,
		}, time.Now())

		if err != nil {
			return err
		}
	}
	return nil
}

type cryptoSource struct{}

func (s cryptoSource) Seed(seed int64) {}

func (s cryptoSource) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}

func (s cryptoSource) Uint64() (v uint64) {
	err := binary.Read(crand.Reader, binary.BigEndian, &v)
	if err != nil {
		log.Fatal(err)
	}
	return v
}
