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
	// UserId int `json:"user_id"`
	Role models.Role `json:"role"`
}


func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// получить headers authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	// проверить, что там лежит bearer jwt токен
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		http.Error(w, "Invalid token format", http.StatusUnauthorized)
		return
	}

	// распарсить токен
	tokenStr := parts[1]
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SigningKey), nil
	})
	fmt.Println(token)


	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}
	fmt.Printf("Parsed claims: %+v\n", claims)
	fmt.Printf("Role: %v\n", claims.Role)
	// fmt.Println(claims.Role)

	ctx := context.WithValue(r.Context(), ContextKeyRole, claims.Role)
	next.ServeHTTP(w, r.WithContext(ctx))

	})
}