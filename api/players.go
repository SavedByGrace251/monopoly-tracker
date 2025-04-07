package api

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Player struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name  string             `bson:"name" json:"name"`
	Money int                `bson:"money" json:"money"`
}

func GetPlayers(c echo.Context, client *mongo.Client) error {
	collection := client.Database("monopoly").Collection("players")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("Database query error:", err)
		return c.JSON(http.StatusInternalServerError, "Database error")
	}
	defer cursor.Close(ctx)

	var players []Player
	if err = cursor.All(ctx, &players); err != nil {
		log.Println("Cursor decoding error:", err)
		return c.JSON(http.StatusInternalServerError, "Decoding error")
	}

	return c.Render(200, )
}

func AddPlayer(c echo.Context, client *mongo.Client) error {
	if err := c.Request().ParseForm(); err != nil {
		log.Println("Form parsing error:", err)
		return c.JSON(http.StatusBadRequest, "Invalid form data")
	}

	name := c.FormValue("name")
	moneyStr := c.FormValue("money")
	money, err := strconv.Atoi(moneyStr)
	if err != nil {
		log.Println("Invalid money value:", err)
		return c.JSON(http.StatusBadRequest, "Money must be a valid number")
	}

	player := Player{
		ID:    primitive.NewObjectID(),
		Name:  name,
		Money: money,
	}

	collection := client.Database("monopoly").Collection("players")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, player)
	if err != nil {
		log.Println("Database insert error:", err)
		return c.JSON(http.StatusInternalServerError, "Database error")
	}

	return c.JSON(http.StatusOK, player)
}
