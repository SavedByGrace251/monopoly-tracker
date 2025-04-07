package api

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(r *echo.Echo, injectClient func(f func(c echo.Context, client *mongo.Client) error) echo.HandlerFunc) {
	r.GET("/players", injectClient(GetPlayers))
	r.POST("/players", injectClient(AddPlayer))
}
