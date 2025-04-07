package ui

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(e *echo.Echo, injectClient func(f func(context echo.Context, client *mongo.Client) error) echo.HandlerFunc) {
	e.GET("/ui/", ServeHTML)
}

func ServeHTML(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}
