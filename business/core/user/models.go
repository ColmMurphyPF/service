package user

import (
	"github.com/colmmurphy91/go-service/business/core/user/db"
	"time"
)

// User represents an individual user.
type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Roles        []string  `json:"roles"`
	PasswordHash []byte    `json:"-"`
	DateCreated  time.Time `json:"date_created"`
	DateUpdated  time.Time `json:"date_updated"`
	ConfirmHash  int64     `json:"confirm_hash"`
	Confirmed    bool      `json:"confirmed"`
}

// NewUser contains information needed to create a new User.
type NewUser struct {
	Name            string   `json:"name" validate:"required"`
	Email           string   `json:"email" validate:"required,email"`
	Roles           []string `json:"roles" validate:"required"`
	Password        string   `json:"password" validate:"required"`
	PasswordConfirm string   `json:"password_confirm" validate:"eqfield=Password"`
}

// UpdateUser defines what information may be provided to modify an existing
// User. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types but we make exceptions around
// marshalling/unmarshalling.
type UpdateUser struct {
	Name            *string  `json:"name"`
	Email           *string  `json:"email" validate:"omitempty,email"`
	Roles           []string `json:"roles"`
	Password        *string  `json:"password"`
	PasswordConfirm *string  `json:"password_confirm" validate:"omitempty,eqfield=Password"`
	ConfirmHash     *int64   `json:"confirm_hash"`
	Confirmed       *bool    `json:"confirmed"`
}

// =============================================================================

func toUser(dbUsr db.User) User {
	var confirmHash int64
	if dbUsr.ConfirmHash.Valid {
		confirmHash = dbUsr.ConfirmHash.Int64
	} else {
		confirmHash = 0
	}
	return User{
		ID:          dbUsr.ID,
		Name:        dbUsr.Name,
		Email:       dbUsr.Email,
		Roles:       dbUsr.Roles,
		DateCreated: dbUsr.DateCreated,
		DateUpdated: dbUsr.DateUpdated,
		ConfirmHash: confirmHash,
		Confirmed:   dbUsr.Confirmed,
	}
}

func toUserSlice(dbUsrs []db.User) []User {
	users := make([]User, len(dbUsrs))
	for i, dbUsr := range dbUsrs {
		users[i] = toUser(dbUsr)
	}
	return users
}
