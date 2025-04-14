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

type ReceptionsRequest struct {
	PvzId uuid.UUID `json:"pvzId"`
}

func (h *Handler) Receptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var receptionsRequest ReceptionsRequest

	if err := json.NewDecoder(r.Body).Decode(&receptionsRequest); err != nil {
		common.WriteErrorResponse(w, http.StatusBadRequest, "Неверный запрос")
        return
    }

	// проверка, закрыта ли предыдушая приемка

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
		createdReception, err := h.services.CreateReception(receptionsRequest.PvzId)
		if err != nil {
			h.logger.Error("Ошибка при создании Приемки")
			common.WriteErrorResponse(w, http.StatusInternalServerError, "Ошибка при создании Приемки")
			return
		} else {
			w.WriteHeader(http.StatusCreated) // 201
			json.NewEncoder(w).Encode(models.Reception {Id: createdReception.Id, DateTime: createdReception.DateTime, PvzId: createdReception.PvzId, Status: createdReception.Status})
		}
	}
}


func (h *Handler) CloseReception(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	pvzIdStr := chi.URLParam(r, "pvzId")
	if pvzIdStr == "" {
		h.logger.Error("Нет параметра")
		common.WriteErrorResponse(w, http.StatusBadRequest, "Неверный запрос")
		return
	}

	pvzId, err := uuid.Parse(pvzIdStr)
	if err != nil {
		h.logger.Error("Некорректный параметр")
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
		updatedReception, err := h.services.CheckReception(pvzId)
		if err != nil {
			h.logger.Error("Приемка уже закрыта")
			common.WriteErrorResponse(w, http.StatusBadRequest, "Неверный запрос")
			return
		} else {
			w.WriteHeader(http.StatusOK) // 200
			json.NewEncoder(w).Encode(models.Reception {Id: updatedReception.Id, DateTime: updatedReception.DateTime, PvzId: updatedReception.PvzId, Status: updatedReception.Status})
		}
	}
}