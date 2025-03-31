package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var Client *mongo.Client

type Player struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name  string             `bson:"name" json:"name"`
	Money int                `bson:"money" json:"money"`
}

func GetPlayers(w http.ResponseWriter, r *http.Request) {
	collection := Client.Database("monopoly").Collection("players")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Database query error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var players []Player
	if err = cursor.All(ctx, &players); err != nil {
		log.Println("Cursor decoding error:", err)
		http.Error(w, "Decoding error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(players)
}

func AddPlayer(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println("Form parsing error:", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	moneyStr := r.FormValue("money")
	money, err := strconv.Atoi(moneyStr)
	if err != nil {
		log.Println("Invalid money value:", err)
		http.Error(w, "Money must be a valid number", http.StatusBadRequest)
		return
	}

	player := Player{
		ID:    primitive.NewObjectID(),
		Name:  name,
		Money: money,
	}

	collection := Client.Database("monopoly").Collection("players")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, player)
	if err != nil {
		log.Println("Database insert error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<li>%s - $%d</li>", player.Name, player.Money)
}
