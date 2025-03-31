package api

import (
    "github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
    r.Get("/players", GetPlayers)
    r.Post("/players", AddPlayer)
}