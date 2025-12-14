package config

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo(cfg *Config) *mongo.Collection {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}

	log.Println("Connected to MongoDB")

	db := client.Database(cfg.MongoDB)
	collection := db.Collection(cfg.MongoCollection)

	return collection
}
