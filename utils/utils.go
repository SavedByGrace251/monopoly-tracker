package utils

import (
	"context"
	"math/rand"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

func GetLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	logger, err := config.Build()
	if err != nil {
		panic("Failed to create logger: " + err.Error())
	}

	zap.ReplaceGlobals(logger)

	return logger
}

func GenerateJoinCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 6)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}