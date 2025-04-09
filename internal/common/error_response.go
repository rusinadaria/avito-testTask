package common

import (
    "encoding/json"
    "net/http"
	"avito-testTask/models"
)

func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string) {
    w.WriteHeader(statusCode)
    errorResponse := models.Error{Message: message}
    json.NewEncoder(w).Encode(errorResponse)
}