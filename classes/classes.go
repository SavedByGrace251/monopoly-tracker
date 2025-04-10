package classes

import "go.mongodb.org/mongo-driver/bson/primitive"

type AppClient struct {
	ClientId    string `bson:"_id,omitempty" json:"client_id"`
	CurrentGame *Game              `bson:"current_game,omitempty" json:"current_game,omitempty"`
	RecentGames []*Game             `bson:"recent_games" json:"recent_games"`
}

type Game struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Players  []Player           `bson:"players" json:"players"`
	JoinCode string             `bson:"join_code" json:"join_code"`
}

type Player struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ClientId primitive.ObjectID `bson:"client_id" json:"client_id"`
	Name     string             `bson:"name" json:"name"`
	Money    int                `bson:"money" json:"money"`
}
