package handlers

import (
	"encoding/json"
	"time"
	"github.com/golang-jwt/jwt"
	"avito-testTask/internal/common"
	"avito-testTask/models"
	"net/http"
	"avito-testTask/internal/handlers/middleware"
)

type RoleRequest struct {
	Role models.Role `json:"role"`
}

func (h *Handler) DummyLogin(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var role RoleRequest

	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		common.WriteErrorResponse(w, http.StatusBadRequest, "Неверный запрос")
        return
    }

	if role.Role != models.RoleEmployee && role.Role != models.RoleModerator {
		common.WriteErrorResponse(w, http.StatusBadRequest, "Неверный запрос") // Недопустимая роль
        return
	}

	var token string
	var err error
	if role.Role == models.RoleEmployee {
		token, err = generateToken(models.RoleEmployee)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Не удалось сгенерировать токен") // обработать
			return
		}
	} else {
		token, err = generateToken(models.RoleModerator)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Не удалось сгенерировать токен") // обработать
			return
		}
	}

	w.Header().Set("Authorization", "Bearer "+token)
    w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct { // ????????
		Token string `json:"token"`
	}{
		Token: token,
	})
}

// const signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"

// type TokenClaims struct {
// 	jwt.StandardClaims
// 	// UserId int `json:"user_id"`
// 	Role models.Role `json:"role"`
// }

func generateToken(role models.Role) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &middleware.TokenClaims{
		// jwt.StandardClaims{
		// 	ExpiresAt: time.Now().Add(12 * time.Hour).Unix(),
		// 	IssuedAt:  time.Now().Unix(),
		// },
		// role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(12 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Role: role,
	})

	return token.SignedString([]byte(middleware.SigningKey))
}
