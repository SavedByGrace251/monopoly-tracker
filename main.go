package main

import (
	"context"
	"fmt"
	"log"
	"monopoly-tracker/api"
	"monopoly-tracker/ui"
	"net/http"
	"text/template"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/labstack/echo/v4"
)

var client *mongo.Client

type Templates struct {
	templates *template.Template
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(options.Credential{
		Username: "root",
		Password: "example",
	})

	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ensure monopoly DB/collection exist
	if err := createDbIfNotExists(client); err != nil {
		log.Fatal(err)
	}

	api.Client = client

	e := echo.New()
	// serve files from static folder

	ui.RegisterRoutes(r)
	api.RegisterRoutes(r)

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func createDbIfNotExists(client *mongo.Client) error {
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
