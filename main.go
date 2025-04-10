package main

import (
	"context"
	"log"
	"monopoly-tracker/api"
	appMiddleware "monopoly-tracker/middleware"
	"monopoly-tracker/utils"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// Setup zap logger for loggin to console
var logger = utils.GetLogger()

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(options.Credential{
		Username: "root",
		Password: "example",
	})
	logger.Debug("MongoDB client options", zap.String("uri", clientOptions.GetURI()))

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	if err := utils.CreateDbIfNotExists(client); err != nil {
		log.Fatal(err)
	}

	// Create injector for mongo client
	injectClient := utils.CreateClientInjector(client)

	e := echo.New()

	// Set up middleware
	e.Pre(appMiddleware.SetClientID())
	e.Use(appMiddleware.ZapLogger(logger))
	e.Use(middleware.Recover())

	e.Renderer = utils.NewTemplate()
	e.Static("/static", "static") // serve files from static folder
	api.RegisterRoutes(e, injectClient)

	logger.Info("Server Starting",
		zap.String("address", ":8080"),
		zap.String("env", "development"),
		zap.String("version", "1.0.0"),
	)

	log.Fatal(e.Start(":8080"))
}
