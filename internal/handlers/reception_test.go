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
)


func TestHandler_Receptions(t *testing.T) {
	type mockBehavior func(s *mock_services.MockReception, pvzId uuid.UUID)

	testReception := models.Reception{
		Id:       uuid.New(),
		// DateTime: time.Now(),
		DateTime: time.Now().Truncate(time.Second),
		PvzId:    uuid.New(),
		Status:   models.InProgress,
	}

	testCases := []struct {
		name                 string
		inputBody            string
		role                 models.Role
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			inputBody: fmt.Sprintf(`{"pvzId":"%s"}`, testReception.PvzId),
			role: models.RoleEmployee,
			mockBehavior: func(s *mock_services.MockReception, pvzId uuid.UUID) {
				s.EXPECT().CreateReception(pvzId).Return(testReception, nil)
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponseBody: fmt.Sprintf(`{
				"id": "%s",
				"dateTime": "%s",
				"pvzId": "%s",
				"status": "%s"
			}`, testReception.Id, testReception.DateTime.Format(time.RFC3339), testReception.PvzId, testReception.Status),
		},
		{
			name: "Forbidden Role",
			inputBody: fmt.Sprintf(`{"pvzId":"%s"}`, testReception.PvzId),
			role: models.RoleModerator,
			mockBehavior: func(s *mock_services.MockReception, pvzId uuid.UUID) {
				// Не вызывается
			},
			expectedStatusCode:   http.StatusForbidden,
			expectedResponseBody: `{"message":"Доступ запрещен: неверная роль"}`,
		},
		{
			name: "Bad Request - invalid JSON",
			inputBody:          `{bad json}`,
			role:               models.RoleEmployee,
			mockBehavior:       func(s *mock_services.MockReception, pvzId uuid.UUID) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: `{"message":"Неверный запрос"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockReception := mock_services.NewMockReception(ctrl)
			tc.mockBehavior(mockReception, testReception.PvzId)

			services := &services.Service{Reception: mockReception}
			handler := NewHandler(services)

			r := chi.NewRouter()
			r.With(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					ctx := context.WithValue(r.Context(), middleware.ContextKeyRole, tc.role)
					next.ServeHTTP(w, r.WithContext(ctx))
				})
			}).Post("/receptions", handler.Receptions)

			req := httptest.NewRequest("POST", "/receptions", bytes.NewBufferString(tc.inputBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.JSONEq(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_CloseReception(t *testing.T) {
	type mockBehavior func(s *mock_services.MockReception, pvzId uuid.UUID)

	testReception := models.Reception{
		Id:       uuid.New(),
		DateTime: time.Now().Truncate(time.Second),
		PvzId:    uuid.New(),
		Status:   models.InProgress,
	}

	testCases := []struct {
		name                 string
		pvzId                string
		role                 models.Role
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			pvzId:     testReception.PvzId.String(),
			role:      models.RoleEmployee,
			mockBehavior: func(s *mock_services.MockReception, pvzId uuid.UUID) {
				s.EXPECT().CheckReception(pvzId).Return(testReception, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: fmt.Sprintf(`{
				"id": "%s",
				"dateTime": "%s",
				"pvzId": "%s",
				"status": "%s"
			}`, testReception.Id, testReception.DateTime.Format(time.RFC3339), testReception.PvzId, testReception.Status),
		},
		{
			name:      "Forbidden Role",
			pvzId:     testReception.PvzId.String(),
			role:      models.RoleModerator,
			mockBehavior: func(s *mock_services.MockReception, pvzId uuid.UUID) {
				// Не вызывается
			},
			expectedStatusCode:   http.StatusForbidden,
			expectedResponseBody: `{"message":"Доступ запрещен: неверная роль"}`,
		},
		{
			name:      "Bad Request - Reception already closed",
			pvzId:     testReception.PvzId.String(),
			role:      models.RoleEmployee,
			mockBehavior: func(s *mock_services.MockReception, pvzId uuid.UUID) {
				s.EXPECT().CheckReception(pvzId).Return(models.Reception{}, fmt.Errorf("Приемка уже закрыта"))
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"Приемка уже закрыта"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockReception := mock_services.NewMockReception(ctrl)
			tc.mockBehavior(mockReception, testReception.PvzId)

			services := &services.Service{Reception: mockReception}
			handler := NewHandler(services)

			r := chi.NewRouter()
			r.With(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					ctx := context.WithValue(r.Context(), middleware.ContextKeyRole, tc.role)
					next.ServeHTTP(w, r.WithContext(ctx))
				})
			}).Post("/pvz/{pvzId}/close_last_reception", handler.CloseReception)

			req := httptest.NewRequest("POST", fmt.Sprintf("/pvz/%s/close_last_reception", tc.pvzId), nil)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.JSONEq(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}
