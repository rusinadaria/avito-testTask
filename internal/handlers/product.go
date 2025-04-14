package handlers

import (
	"github.com/google/uuid"
	"net/http"
	"avito-testTask/internal/handlers/middleware"
	"avito-testTask/internal/common"
	"avito-testTask/models"
	"fmt"
	"encoding/json"
	"github.com/go-chi/chi"
)

type RequestProduct struct {
	Type models.Type `json:"type"`
	PvzId uuid.UUID `json:"pvzId"`
}

func (h *Handler) AddProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var requestProduct RequestProduct

	if err := json.NewDecoder(r.Body).Decode(&requestProduct); err != nil {
		common.WriteErrorResponse(w, http.StatusBadRequest, "Неверный запрос")
        return
    }

	role, ok := r.Context().Value(middleware.ContextKeyRole).(models.Role)
	fmt.Println(role)
	if !ok {
		h.logger.Error("Нет роли")
		common.WriteErrorResponse(w, http.StatusForbidden, "Доступ запрещен")
		return
	}

	if role != models.RoleEmployee {
		h.logger.Error("Неверная роль")
		common.WriteErrorResponse(w, http.StatusForbidden, "Доступ запрещен")
		return
	} else {
		createdProduct, err := h.services.AddProduct(requestProduct.Type, requestProduct.PvzId)
		if err != nil {
			h.logger.Error("Не удалось добавить товар")
			common.WriteErrorResponse(w, http.StatusBadRequest, "Неверный запрос")
        	return
		} else {
			w.WriteHeader(http.StatusCreated) // 201
			json.NewEncoder(w).Encode(models.Product {Id: createdProduct.Id, DateTime: createdProduct.DateTime, Type: createdProduct.Type, ReceptionId: createdProduct.ReceptionId})
		}
	}
	// 201 + schemas/Product
}

func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	pvzIdStr := chi.URLParam(r, "pvzId")
	if pvzIdStr == "" {
		h.logger.Error("Нет параметра")
		common.WriteErrorResponse(w, http.StatusBadRequest, "Неверный запрос")
		return
	}

	pvzId, err := uuid.Parse(pvzIdStr)
	if err != nil {
		h.logger.Error("Неверный параметр")
		common.WriteErrorResponse(w, http.StatusBadRequest, "Неверный запрос")
		return
	}

	role, ok := r.Context().Value(middleware.ContextKeyRole).(models.Role)
	fmt.Println(role)
	if !ok {
		h.logger.Error("Нет роли")
		common.WriteErrorResponse(w, http.StatusForbidden, "Доступ запрещен")
		return
	}

	if role != models.RoleEmployee {
		h.logger.Error("Неверная роль")
		common.WriteErrorResponse(w, http.StatusForbidden, "Доступ запрещен")
		return
	} else {
		err := h.services.DeleteProduct(pvzId)
		if err != nil {
			common.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		} else {
			w.WriteHeader(http.StatusOK) // 200
		}
	}
	// 200
}