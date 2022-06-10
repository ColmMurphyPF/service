package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type Cafe struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name"`
	Address     string             `bson:"address"`
	PhoneNumber string             `bson:"phone_number"`
}
