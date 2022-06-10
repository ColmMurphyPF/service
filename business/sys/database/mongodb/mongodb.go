package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Config struct {
	User     string
	Password string
	Host     string
	DBName   string
}

func Open(cfg Config) (*mongo.Database, error) {
	credential := options.Credential{
		AuthSource: "admin",
		Username:   cfg.User,
		Password:   cfg.Password,
	}
	clientOptions := options.Client().ApplyURI("mongodb://" + cfg.Host + "/?ssl=false").SetAuth(credential)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")

	test := client.Database(cfg.DBName)

	return test, nil
}
