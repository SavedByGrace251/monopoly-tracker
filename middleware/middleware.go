package middleware

import (
	"monopoly-tracker/utils"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

var logger = utils.GetLogger()

func SetClientID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("clientId")
			if err != nil || cookie == nil || cookie.Value == "" {
				logger.Debug("ClientId cookie not found, creating a new one")
				newID := uuid.New().String()
				// Store the clientId in the current request context
				c.Set("clientId", newID)
				logger.Debug("New ClientId generated", zap.String("clientId", newID))
				// Set it as a response cookie
				c.SetCookie(&http.Cookie{
					Name:  "clientId",
					Value: newID,
					Path:  "/",
				})
			} else {
				logger.Debug("ClientId cookie found", zap.String("clientId", cookie.Value))
				c.Set("clientId", cookie.Value)
			}
			return next(c)
		}
	}
}

// ZapLogger is a custom middleware for Echo that logs requests using Zap
func ZapLogger(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Process request
			err := next(c)

			// Prepare log fields
			req := c.Request()
			res := c.Response()

			fields := []zap.Field{
				zap.String("method", req.Method),
				zap.String("uri", req.RequestURI),
				zap.Int("status", res.Status),
				zap.Duration("latency", time.Since(start)),
				zap.String("remote_ip", c.RealIP()),
				zap.String("host", req.Host),
				zap.String("user_agent", req.UserAgent()),
			}

			// Add error field if there's an error
			if err != nil {
				fields = append(fields, zap.Error(err))
			}

			// Log based on status code
			switch {
			case res.Status >= 500:
				logger.Error("Server error", fields...)
			case res.Status >= 400:
				logger.Warn("Client error", fields...)
			case res.Status >= 300:
				logger.Info("Redirect", fields...)
			default:
				logger.Info("Request completed", fields...)
			}

			return err
		}
	}
}
