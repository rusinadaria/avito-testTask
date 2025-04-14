package middleware

import (
	"net/http"
	"fmt"
	"strings"
	"github.com/golang-jwt/jwt"
	"context"
	"avito-testTask/models"
)


type contextKey string

const ContextKeyRole contextKey = "role"

const SigningKey = "qrkjk#4#%35FSFJlja#4353KSFjH"

type TokenClaims struct {
	jwt.StandardClaims
	Role models.Role `json:"role"`
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		http.Error(w, "Invalid token format", http.StatusUnauthorized)
		return
	}

	tokenStr := parts[1]
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SigningKey), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	ctx := context.WithValue(r.Context(), ContextKeyRole, claims.Role)
	next.ServeHTTP(w, r.WithContext(ctx))
	})
}