package handlers

import (
	"net/http"
	"avito-testTask/internal/handlers/middleware"
	"avito-testTask/internal/common"
	"avito-testTask/models"
	"fmt"
	"encoding/json"
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
	}

    w.WriteHeader(http.StatusCreated) // 201 ПВЗ создан
	json.NewEncoder(w).Encode(models.PVZ {Id: pvz.Id, RegistrationDate: pvz.RegistrationDate, City: pvz.City})
}