package main

import (
	"context"
	"fmt"
	"log"
	"monopoly-tracker/api"
	"monopoly-tracker/ui"
	"monopoly-tracker/utils"
	"time"

	
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/labstack/echo/v4"
)

var client *mongo.Client

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
	if err := utils.CreateDbIfNotExists(client); err != nil {
		log.Fatal(err)
	}

	// Create injector for mongo client
	injectClient := utils.CreateClientInjector(client)

	e := echo.New()
	// serve files from static folder
	e.Static("/static", "static")
	ui.RegisterRoutes(e, injectClient)
	api.RegisterRoutes(e, injectClient)

	fmt.Println("Server running on :8080")
	log.Fatal(e.Start(":8080"))
}

