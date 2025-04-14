package handlers

import (
	mock_services "avito-testTask/internal/services/mocks"
	"avito-testTask/internal/services"
	"avito-testTask/models"
	"bytes"
	"net/http/httptest"
	"testing"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	// "strings"
	"github.com/google/uuid"
	"time"
	"fmt"
	"context"
	"avito-testTask/internal/handlers/middleware"
	"log/slog"
	"os"
)

func TestHandler_PVZCreate(t *testing.T) {
	type mockBehavior func(s *mock_services.MockPVZ, input models.PVZ)

	// fixedTime := time.Date(2025, 4, 14, 16, 40, 50, 0, time.FixedZone("UTC+7", 7*60*60))
	// testID := uuid.New()

	// testPVZ := models.PVZ{
	// 	Id:               testID,
	// 	RegistrationDate: fixedTime,
	// 	City:             models.Moscow,
	// }
	fixedTime := time.Date(2025, 4, 14, 16, 40, 50, 0, time.UTC)
	testPVZ := models.PVZ{
		Id:               uuid.New(),
		RegistrationDate: fixedTime,
		City:             models.Moscow,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	testTable := []struct {
		name                 string
		inputBody            string
		role                 models.Role
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: fmt.Sprintf(`{"id":"%s","registrationDate":"%s","city":"%s"}`, testPVZ.Id, testPVZ.RegistrationDate.Format(time.RFC3339), testPVZ.City),
			role:      models.RoleModerator,
			mockBehavior: func(s *mock_services.MockPVZ, input models.PVZ) {
				s.EXPECT().CreatePVZ(input).Return(input, nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: fmt.Sprintf(`{"id":"%s","registrationDate":"%s","city":"%s"}`, testPVZ.Id, testPVZ.RegistrationDate.Format(time.RFC3339), testPVZ.City),
		},
		{
			name:      "Forbidden role",
			inputBody: `{}`,
			role:      models.RoleEmployee,
			mockBehavior: func(s *mock_services.MockPVZ, input models.PVZ) {},
			expectedStatusCode:   http.StatusForbidden,
			expectedResponseBody: `{"message":"Доступ запрещен"}`,
		},
		{
			name:      "Invalid JSON",
			inputBody: `{invalid_json}`,
			role:      models.RoleModerator,
			mockBehavior: func(s *mock_services.MockPVZ, input models.PVZ) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"Неверный запрос"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			pvzService := mock_services.NewMockPVZ(c)
			testCase.mockBehavior(pvzService, testPVZ)

			services := &services.Service{PVZ: pvzService}

			handler := NewHandler(services, logger)

			r := chi.NewRouter()
			r.Post("/pvz", func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), middleware.ContextKeyRole, testCase.role)
				handler.PVZCreate(w, r.WithContext(ctx))
			})

			req := httptest.NewRequest("POST", "/pvz", bytes.NewBufferString(testCase.inputBody))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.JSONEq(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}


