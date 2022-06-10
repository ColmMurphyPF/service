package cafe

import (
	"context"
	"errors"
	"fmt"
	"github.com/colmmurphy91/go-service/business/core/cafe/db"
	"github.com/colmmurphy91/go-service/business/sys/validate"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var (
	ErrNotFound = errors.New("cafe not found")
	ErrDeletion = errors.New("unable to delete cafe")
)

// Core manages the set of APIs for product access.
type Core struct {
	store db.Store
}

// NewCore constructs a core for product api access.
func NewCore(log *zap.SugaredLogger, monDB *mongo.Database) Core {
	return Core{
		store: db.NewStore(log, monDB, context.Background()),
	}
}

func (c Core) CreateCafe(nc NewCafe) (Cafe, error) {
	err := validate.Check(nc)
	if err != nil {
		return Cafe{}, fmt.Errorf("validating data: %w", err)
	}
	dbCafe := db.Cafe{
		Name:        nc.Name,
		Address:     nc.Address,
		PhoneNumber: nc.PhoneNumber,
	}
	cafe, err := c.store.Save(dbCafe)
	if err != nil {
		return Cafe{}, err
	}
	return toCafe(cafe), nil
}

func (c Core) FindCafe(id string) (Cafe, error) {
	cafe, err := c.store.FindById(id)
	if err != nil {
		return Cafe{}, ErrNotFound
	}
	return toCafe(cafe), nil
}

func (c Core) UpdateCafe(uc UpdateCafe) error {
	id, err := primitive.ObjectIDFromHex(uc.ID)
	if err != nil {
		fmt.Println("failing here")
		return err
	}
	byId, err := c.store.FindById(id.Hex())
	if err != nil {
		return ErrNotFound
	}
	if uc.Name != nil {
		byId.Name = *uc.Name
	}

	if uc.Address != nil {
		byId.Address = *uc.Address
	}

	if uc.PhoneNumber != nil {
		byId.PhoneNumber = *uc.PhoneNumber
	}
	err = c.store.UpdateCafe(byId)
	if err != nil {
		return err
	}
	return nil
}

func (c Core) DeleteCafe(id string) error {
	err := c.store.DeleteByID(id)
	if err != nil {
		return ErrDeletion
	}
	return nil
}

func (c Core) FindAll() ([]Cafe, error) {
	all, err := c.store.FindAll()
	if err != nil {
		return nil, err
	}
	var cafes []Cafe
	for _, j := range all {
		cafes = append(cafes, toCafe(j))
	}
	if cafes == nil {
		return []Cafe{}, nil
	}
	return cafes, nil
}

func toCafe(cafe db.Cafe) Cafe {
	return Cafe{
		ID:          cafe.ID.Hex(),
		Name:        cafe.Name,
		Address:     cafe.Address,
		PhoneNumber: cafe.PhoneNumber,
	}
}
