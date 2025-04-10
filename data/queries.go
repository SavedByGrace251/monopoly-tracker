package data

import (
	"context"
	"monopoly-tracker/classes"
	"monopoly-tracker/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var logger = utils.GetLogger()

func GetClientWithGames(ctx context.Context, db *mongo.Database, clientId string) (*classes.AppClient, error) {
	clientCollection := db.Collection("clients")

	// Fetch client information
	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: clientId}}}},
		bson.D{{Key: "$lookup", 
			Value: bson.D{
				{Key: "from", Value: "games"},
				{Key: "localField", Value: "currentGame"},
				{Key: "foreignField", Value: "_id"},
				{Key: "as", Value: "currentGame"},
		}}},
		bson.D{{Key: "$unwind", Value: "$currentGame"}},
	}

	cursor, err := clientCollection.Aggregate(ctx, pipeline)
	if err != nil {
		logger.Error("Failed to aggregate client data", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var client classes.AppClient
	if cursor.Next(ctx) {
		err := cursor.Decode(&client)
		if err != nil {
			logger.Error("Failed to decode client data", zap.Error(err))
			return nil, err
		}
		return &client, nil
	}

	return nil, mongo.ErrNoDocuments
}
