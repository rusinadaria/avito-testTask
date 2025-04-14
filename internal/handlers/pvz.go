package handlers

import (
	"net/http"
	// "avito-testTask/internal/handlers/middleware"
	"avito-testTask/internal/common"
	"avito-testTask/models"
	"encoding/json"
	"strconv"
	"time"
)

func (h *Handler) PVZCreate(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var pvz models.PVZ

	if err := json.NewDecoder(r.Body).Decode(&pvz); err != nil {
		common.WriteErrorResponse(w, http.StatusBadRequest, "Неверный запрос")
        return
    }

	createdPVZ, err := h.services.CreatePVZ(pvz)
	if err != nil {
		h.logger.Error("Ошибка при создании ПВЗ")
		common.WriteErrorResponse(w, http.StatusBadRequest, "Неверный запрос")
		return
	} else {
		w.WriteHeader(http.StatusCreated) // 201 ПВЗ создан
		json.NewEncoder(w).Encode(models.PVZ {Id: createdPVZ.Id, RegistrationDate: createdPVZ.RegistrationDate, City: createdPVZ.City})
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
			h.logger.Error("Invalid startDate format")
			common.WriteErrorResponse(w, http.StatusBadRequest, "Неверный запрос")
        	return
		}
	}

	if endDateStr != "" {
		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			h.logger.Error("Invalid endDate format")
			common.WriteErrorResponse(w, http.StatusBadRequest, "Неверный запрос")
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
		h.logger.Error("Ошибка при получении списка ПВЗ")
		common.WriteErrorResponse(w, http.StatusBadRequest, "Неверный запрос")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pvzList)
}