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
		common.WriteErrorResponse(w, http.StatusForbidden, "Доступ запрещен: No role in context") // обработать
		return
	}

	if role != models.RoleEmployee {
		common.WriteErrorResponse(w, http.StatusForbidden, "Доступ запрещен: неверная роль")
		return
	} else {
		createdReception, err := h.services.CreateReception(receptionsRequest.PvzId)
		if err != nil {
			common.WriteErrorResponse(w, http.StatusInternalServerError, "Ошибка при создании Приемки")
			return
		} else {
			w.WriteHeader(http.StatusCreated) // 201
			json.NewEncoder(w).Encode(models.Reception {Id: createdReception.Id, DateTime: createdReception.DateTime, PvzId: createdReception.PvzId, Status: createdReception.Status})
		}
		// 201 + models.Reception
	}
}


func (h *Handler) CloseReception(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	pvzIdStr := chi.URLParam(r, "pvzId")
	if pvzIdStr == "" {
		http.Error(w, "pvzId not provided", http.StatusBadRequest)
		return
	}

	pvzId, err := uuid.Parse(pvzIdStr)
	if err != nil {
		http.Error(w, "invalid pvzId", http.StatusBadRequest)
		return
	}

	role, ok := r.Context().Value(middleware.ContextKeyRole).(models.Role)
	fmt.Println(role)
	if !ok {
		common.WriteErrorResponse(w, http.StatusForbidden, "Доступ запрещен: No role in context") // обработать
		return
	}

	if role != models.RoleEmployee {
		common.WriteErrorResponse(w, http.StatusForbidden, "Доступ запрещен: неверная роль")
		return
	} else {
		// проверить на закрытие
		updatedReception, err := h.services.CheckReception(pvzId)
		if err != nil {
			common.WriteErrorResponse(w, http.StatusBadRequest, "Приемка уже закрыта")
			return
		} else {
			w.WriteHeader(http.StatusOK) // 200
			json.NewEncoder(w).Encode(models.Reception {Id: updatedReception.Id, DateTime: updatedReception.DateTime, PvzId: updatedReception.PvzId, Status: updatedReception.Status})
		}
	}
	// 200 + schemas/Reception
}