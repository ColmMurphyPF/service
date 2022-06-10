package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Store struct {
	log *zap.SugaredLogger
	db  *mongo.Database
	trx *mongo.Session
	ctx context.Context
}

func (s Store) Save(c Cafe) (Cafe, error) {
	c.ID = primitive.NewObjectID()
	one, err := s.db.Collection("cafes").InsertOne(s.ctx, c)
	if err != nil {
		return Cafe{}, err
	}
	c.ID = one.InsertedID.(primitive.ObjectID)
	return c, err
}

func (s Store) FindById(id string) (Cafe, error) {
	var cafe Cafe
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Cafe{}, err
	}
	one := s.db.Collection("cafes").FindOne(s.ctx, bson.M{
		"_id": hex,
	})
	err = one.Decode(&cafe)
	if err != nil {
		return Cafe{}, err
	}
	return cafe, nil
}

func (s Store) UpdateCafe(cafe Cafe) error {
	filter := bson.D{{"_id", cafe.ID}}
	update := bson.D{{"$set", bson.M{
		"name":         cafe.Name,
		"address":      cafe.Address,
		"phone_number": cafe.PhoneNumber,
	}}}
	_, err := s.db.Collection("cafes").UpdateOne(s.ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (s Store) DeleteByID(id string) error {
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = s.db.Collection("cafes").DeleteOne(s.ctx, bson.M{"_id": hex})
	if err != nil {
		return err
	}
	return nil
}

func (s Store) FindAll() ([]Cafe, error) {
	var cafes []Cafe
	cur, err := s.db.Collection("cafes").Find(s.ctx, bson.M{})
	err = cur.All(s.ctx, &cafes)
	if err != nil {
		return nil, err
	}
	return cafes, nil
}

// NewStore constructs a data for api access.
func NewStore(log *zap.SugaredLogger, db *mongo.Database, ctx context.Context) Store {
	return Store{
		log: log,
		db:  db,
		ctx: ctx,
	}
}
