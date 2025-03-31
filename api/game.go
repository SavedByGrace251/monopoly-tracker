package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

type Game struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Players  []Player           `bson:"players" json:"players"`
	JoinCode string             `bson:"join_code" json:"join_code"`
}

func GetGames(w http.ResponseWriter, r *http.Request) {
	collection := Client.Database("monopoly").Collection("games")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Database query error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var games []Game
	if err = cursor.All(ctx, &games); err != nil {
		log.Println("Cursor decoding error:", err)
		http.Error(w, "Decoding error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(games)
}

func CreateGame(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("Form parsing error:", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	collection := Client.Database("monopoly").Collection("games")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Make sure there's a unique index on join_code
	_, _ = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"join_code": 1},
		Options: options.Index().SetUnique(true),
	})

	var game Game
	game.Name = name

	var maxTries = 10

	for i := 0; i < maxTries; i++ {
		joinCode := fmt.Sprintf("%06d", rand.Intn(1000000))
		game.JoinCode = joinCode
		game.ID = primitive.NewObjectID()

		_, err := collection.InsertOne(ctx, game)
		if err != nil {
			// Check if it's a duplicate key error:
			if writeErr, ok := err.(mongo.WriteException); ok && len(writeErr.WriteErrors) > 0 {
				// Retry on duplicate join code
				if writeErr.WriteErrors[0].Code == 11000 {
					continue
				}
			}
			log.Println("Database insertion error:", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		// Successfully inserted
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(game)
		return
	}

	http.Error(w, "Failed to generate unique join code", http.StatusInternalServerError)
}
