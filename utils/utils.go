package utils

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateClientInjector(client *mongo.Client) func(f func(contextIn echo.Context, clientIn *mongo.Client) error) echo.HandlerFunc {
	return func(f func(contextIn echo.Context, clientIn *mongo.Client) error) echo.HandlerFunc {
		return func(context echo.Context) error {
			return f(context, client)
		}
	}
}

func CreateDbIfNotExists(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check databases
	dbs, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		return err
	}

	// If monopoly DB is missing, create a collection to finalize creation
	for _, dbName := range dbs {
		if dbName == "monopoly" {
			return nil
		}
	}

	db := client.Database("monopoly")
	if err := db.CreateCollection(ctx, "players"); err != nil {
		return err
	}
	return nil
}
