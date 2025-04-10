package api

import (
	"monopoly-tracker/classes"
	"monopoly-tracker/data"
	"monopoly-tracker/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"context"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var logger = utils.GetLogger()

func RegisterRoutes(e *echo.Echo, injectClient func(f func(c echo.Context, client *mongo.Client) error) echo.HandlerFunc) {
	e.GET("/", injectClient(index))
	e.POST("/create-game", injectClient(CreateGame))
	e.GET("/players", injectClient(GetPlayers))
	e.POST("/players", injectClient(AddPlayer))
}

func index(c echo.Context, client *mongo.Client) error {
	clientId := c.Get("clientId").(string)

	// Get information about the clientId
	collection := client.Database("monopoly").Collection("clients")
	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	appClient, err := data.GetClientWithGames(ctx, client.Database("monopoly"), clientId)
	if err == mongo.ErrNoDocuments {
		// If no client found, create a new one
		if err == mongo.ErrNoDocuments {
			appClient = &classes.AppClient{
				ClientId:    clientId,
				CurrentGame: nil,
				RecentGames: []*classes.Game{},
			}
			_, err = collection.InsertOne(ctx, *appClient)
			if err != nil {
				logger.Error("Failed to insert new client", zap.Error(err))
				return c.Render(http.StatusInternalServerError, "error", "Failed to create new client")
			}
		} else {
			logger.Error("Database query error", zap.Error(err))
			return c.Render(http.StatusInternalServerError, "error", "Database error")
		}
	} else if err != nil {
		logger.Error("Database query error", zap.Error(err))
		return c.Render(http.StatusInternalServerError, "error", "Database error")
	}

	return c.Render(http.StatusOK, "index", appClient)
}

func CreateGame(c echo.Context, client *mongo.Client) error {
	clientId := c.Get("clientId").(string)

	// Parse form data
	if err := c.Request().ParseForm(); err != nil {
		logger.Error("Form parsing error", zap.Error(err))
		return c.JSON(http.StatusBadRequest, "Invalid form data")
	}

	gameName := c.FormValue("game_name")
	joinCode := utils.GenerateJoinCode()
	codeCreate := false
	for i := 0; i < 10; i++ {
		// check to see if any other code exists in the database
		collection := client.Database("monopoly").Collection("games")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var existingGame classes.Game
		collection.FindOne(ctx, bson.M{"join_code": joinCode}).Decode(&existingGame)
		if existingGame.JoinCode != "" && existingGame.JoinCode != joinCode {
			// If the code already exists, generate a new one
			joinCode = utils.GenerateJoinCode()
		} else {
			codeCreate = true
			break
		}
	}

	if !codeCreate {
		logger.Error("Failed to create unique join code")
		return c.JSON(http.StatusInternalServerError, "Failed to create unique join code")
	}

	if gameName == "" {
		logger.Error("Game name is required")
		return c.JSON(http.StatusBadRequest, "Game name is required")
	}

	// Create a new game object
	game := classes.Game{
		ID:       primitive.NewObjectID(),
		Name:     gameName,
		JoinCode: joinCode,
		Players:  []classes.Player{},
	}

	// Insert the game into the database
	collection := client.Database("monopoly").Collection("games")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, game)
	if err != nil {
		logger.Error("Database insert error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, "Database error")
	}

	// Add the game to the client's recent games
	clientCollection := client.Database("monopoly").Collection("clients")
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = clientCollection.UpdateOne(ctx, bson.M{"_id": clientId}, bson.M{
		"$push": bson.M{"recent_games": game.ID},
		"$set":  bson.M{"current_game": game.ID},
	})
	if err != nil {
		logger.Error("Failed to update client with new game", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, "Failed to update client with new game")
	}

	return c.JSON(http.StatusOK, game)
}

func GetPlayers(c echo.Context, client *mongo.Client) error {
	collection := client.Database("monopoly").Collection("players")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		logger.Error("Database query error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, "Database error")
	}
	defer cursor.Close(ctx)

	var players []classes.Player
	if err = cursor.All(ctx, &players); err != nil {
		logger.Error("Cursor decoding error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, "Decoding error")
	}

	return c.Render(200, "players", players)
}

func AddPlayer(c echo.Context, client *mongo.Client) error {
	if err := c.Request().ParseForm(); err != nil {
		logger.Error("Form parsing error", zap.Error(err))
		return c.JSON(http.StatusBadRequest, "Invalid form data")
	}

	name := c.FormValue("name")
	moneyStr := c.FormValue("money")
	money, err := strconv.Atoi(moneyStr)
	if err != nil {
		logger.Error("Money conversion error", zap.String("moneyStr", moneyStr), zap.Error(err))
		return c.JSON(http.StatusBadRequest, "Money must be a valid number")
	}

	player := classes.Player{
		ID:    primitive.NewObjectID(),
		Name:  name,
		Money: money,
	}

	collection := client.Database("monopoly").Collection("players")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, player)
	if err != nil {
		logger.Error("Database insert error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, "Database error")
	}

	return c.JSON(http.StatusOK, player)
}
