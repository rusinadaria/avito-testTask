package middleware

import (
	"fmt"
	"net/http"
	// "fmt"
	"avito-testTask/models"
	// "avito-testTask/internal/common"
)

func CheckRoleMiddleware(allowedRoles ...models.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, _ := r.Context().Value(ContextKeyRole).(models.Role)
			fmt.Println(role)

			for _, allowed := range allowedRoles {
				if role == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, "Access Denied", http.StatusForbidden)
		})
	}
}