package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	mongoDbName = "sampleapp"

	mongoDogCollection = "dogs"
)

type Repository struct {
	client *mongo.Client
}

func New(client *mongo.Client) *Repository {
	return &Repository{
		client,
	}
}

func (r *Repository) GetAllDogs(ctx context.Context) ([]Dog, error) {
	dogCollection := r.client.Database(mongoDbName).Collection(mongoDogCollection)
	cursor, err := dogCollection.Find(ctx, bson.D{})

	if err != nil {
		return nil, fmt.Errorf("dogCollection.Find: %w", err)
	}

	var results []Dog
	err = cursor.All(ctx, &results)

	if err != nil {
		return nil, fmt.Errorf("cursor.All: %w", err)
	}

	return results, nil
}
