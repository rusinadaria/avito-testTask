package handlers

import (
	// "log"
	"log/slog"
	"github.com/go-chi/chi"
	"avito-testTask/internal/services"
	"net/http"
	"avito-testTask/internal/handlers/middleware"
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

	r.Post("/dummyLogin", h.DummyLogin) // Получение тестового токена
	r.Group(func(r chi.Router) {
        r.Use(middleware.JWTMiddleware)

        r.Post("/pvz", h.PVZCreate) // Создание ПВЗ (только для модераторов)
		r.Post("/receptions", h.Receptions) // Создание новой приемки товаров (только для сотрудников ПВЗ)
		r.Post("/pvz/{pvzId}/close_last_reception", h.CloseReception) // Закрытие последней открытой приемки товаров в рамках ПВЗ (только для сотрудников ПВЗ)
		r.Post("/products", h.AddProduct) // Добавление товара в текущую приемку (только для сотрудников ПВЗ)
		r.Post("/pvz/{pvzId}/delete_last_product", h.DeleteProduct) // Удаление последнего добавленного товара из текущей приемки (LIFO, только для сотрудников ПВЗ)
		r.Get("/pvz", h.GetPvz) // Получение списка ПВЗ с фильтрацией по дате приемки и пагинацией (для сотрудников и модераторов)
    })
	return r
}