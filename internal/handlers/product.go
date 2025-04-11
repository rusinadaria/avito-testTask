package handlers

import (
	"github.com/google/uuid"
	"net/http"
	"avito-testTask/internal/handlers/middleware"
	"avito-testTask/internal/common"
	"avito-testTask/models"
	"fmt"
	"encoding/json"
	// "github.com/go-chi/chi"
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
		common.WriteErrorResponse(w, http.StatusForbidden, "Доступ запрещен: No role in context")
		return
	}

	if role != models.RoleEmployee {
		common.WriteErrorResponse(w, http.StatusForbidden, "Доступ запрещен: неверная роль")
		return
	} else {
		createdProduct, err := h.services.AddProduct(requestProduct.Type, requestProduct.PvzId)
		if err != nil {
			common.WriteErrorResponse(w, http.StatusInternalServerError, "Не удалось добавить товар") // тут должна быть 400 ошибка
        	return
		} else {
			w.WriteHeader(http.StatusCreated) // 201
			json.NewEncoder(w).Encode(models.Product {Id: createdProduct.Id, DateTime: createdProduct.DateTime, Type: createdProduct.Type, ReceptionId: createdProduct.ReceptionId})
		}
		// При этом товар должен привязываться к последнему незакрытому приёму товаров в рамках текущего ПВЗ.

		// if checkStatusLastReception(pvzId) == in_progress {
		// 	addProduct(requestProduct)
		// }

		// Если же нет новой незакрытой приёмки товаров, то в таком случае должна возвращаться ошибка, и товар не должен добавляться в систему.

		// if checkStatusLastReception == close {
		// 	return err
		// }


		// Если последняя приёмка товара все ещё не была закрыта, то результатом должна стать привязка товара к текущему ПВЗ и текущей приёмке с последующем добавлением данных в хранилище.
	}
	// 201 + schemas/Product
}