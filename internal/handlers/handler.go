package handlers

import (
	"log/slog"
	"github.com/go-chi/chi"
	"avito-testTask/internal/services"
	"net/http"
	"avito-testTask/internal/handlers/middleware"
	"avito-testTask/models"
)

type Handler struct {
	services *services.Service
	logger   *slog.Logger
}

func NewHandler(services *services.Service, logger *slog.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}

func (h *Handler) InitRoutes(logger *slog.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.LoggerMiddlewareWrapper(logger))

	r.Post("/dummyLogin", h.DummyLogin)
	r.Group(func(r chi.Router) {
        r.Use(middleware.JWTMiddleware)

			// r.Get("/pvz", h.GetPvz)

			r.Group(func(r chi.Router) {
				r.Use(middleware.CheckRoleMiddleware(models.RoleModerator))
	
				r.Post("/pvz", h.PVZCreate)
			})
	
			r.Group(func(r chi.Router) {
				r.Use(middleware.CheckRoleMiddleware(models.RoleEmployee))
	
				r.Post("/receptions", h.Receptions)
				r.Post("/pvz/{pvzId}/close_last_reception", h.CloseReception)
				r.Post("/products", h.AddProduct)
				r.Post("/pvz/{pvzId}/delete_last_product", h.DeleteProduct)
			})
	
			r.Group(func(r chi.Router) {
				r.Use(middleware.CheckRoleMiddleware(models.RoleModerator, models.RoleEmployee))
	
				r.Get("/pvz", h.GetPvz)
			})
    })
	return r
}