package handlers

import (
	"net/http"
	"avito-testTask/internal/handlers/middleware"
	"avito-testTask/internal/common"
	"avito-testTask/models"
	"fmt"
	"encoding/json"
	"strconv"
	"time"
	// "github.com/go-chi/chi"
)

func (h *Handler) PVZCreate(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var pvz models.PVZ

	if err := json.NewDecoder(r.Body).Decode(&pvz); err != nil {
		common.WriteErrorResponse(w, http.StatusBadRequest, "Неверный запрос")
        return
    }

	role, ok := r.Context().Value(middleware.ContextKeyRole).(models.Role)
	fmt.Println(role)
	if !ok {
		common.WriteErrorResponse(w, http.StatusForbidden, "Доступ запрещен: No role in context") // обработать
		return
	}

	if role != models.RoleModerator {
		common.WriteErrorResponse(w, http.StatusForbidden, "Доступ запрещен")
		return
	} else {
		createdPVZ, err := h.services.CreatePVZ(pvz)
		if err != nil {
			fmt.Println("Ошибка в хэндлере при создании ПВЗ")
			common.WriteErrorResponse(w, http.StatusInternalServerError, "Ошибка при создании ПВЗ")
        	return
		} else {
			w.WriteHeader(http.StatusCreated) // 201 ПВЗ создан
			json.NewEncoder(w).Encode(models.PVZ {Id: createdPVZ.Id, RegistrationDate: createdPVZ.RegistrationDate, City: createdPVZ.City})
		}
	}
}

func (h *Handler) GetPvz(w http.ResponseWriter, r *http.Request) {
	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	var (
		startDate time.Time
		endDate   time.Time
		err       error
	)

	if startDateStr != "" {
		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			http.Error(w, "Invalid startDate format", http.StatusBadRequest)
			return
		}
	}

	if endDateStr != "" {
		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			http.Error(w, "Invalid endDate format", http.StatusBadRequest)
			return
		}
	}

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 30 {
			limit = l
		}
	}

	pvzList, err := h.services.GetPvzList(&startDate, &endDate, page, limit)
	if err != nil {
		http.Error(w, "Failed to fetch PVZ list", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pvzList)
}