package handlers

import (
	// "log"
	"log/slog"
	"github.com/go-chi/chi"
	// "github.com/go-chi/chi/v5"
	"avito-testTask/internal/services"
	"net/http"
	"avito-testTask/internal/handlers/middleware"
)

type Handler struct {
	services *services.Service
}

func NewHandler(services *services.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes(logger *slog.Logger) http.Handler {
	r := chi.NewRouter()

	r.Post("/dummyLogin", h.DummyLogin) // Получение тестового токена

	r.Group(func(r chi.Router) {
        r.Use(middleware.JWTMiddleware)

        r.Post("/pvz", h.PVZCreate) // Создание ПВЗ (только для модераторов)
    })
	return r
}