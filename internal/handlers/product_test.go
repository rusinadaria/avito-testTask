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

func TestHandler_AddProduct(t *testing.T) {
	type mockBehavior func(s *mock_services.MockProduct, prodType models.Type, pvzId uuid.UUID)

	testPvzId := uuid.New()
	testReceptionId := uuid.New()
	testProductId := uuid.New()
	testTime := time.Date(2025, 4, 14, 16, 40, 50, 0, time.UTC)

	testProduct := models.Product{
		Id:          testProductId,
		DateTime:    testTime,
		Type:        models.Clothes,
		ReceptionId: testReceptionId,
	}

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
			inputBody: fmt.Sprintf(`{"type":"одежда", "pvzId":"%s"}`, testPvzId),
			role:      models.RoleEmployee,
			mockBehavior: func(s *mock_services.MockProduct, prodType models.Type, pvzId uuid.UUID) {
				s.EXPECT().AddProduct(prodType, pvzId).Return(testProduct, nil)
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponseBody: fmt.Sprintf(`{
				"id": "%s",
				"dateTime": "%s",
				"type": "одежда",
				"receptionId": "%s"
			}`, testProductId, testTime.Format(time.RFC3339), testReceptionId),
		},
		{
			name:      "Forbidden role",
			inputBody: fmt.Sprintf(`{"type":"одежда", "pvzId":"%s"}`, testPvzId),
			role:      models.RoleModerator,
			mockBehavior: func(s *mock_services.MockProduct, prodType models.Type, pvzId uuid.UUID) {},
			expectedStatusCode:   http.StatusForbidden,
			expectedResponseBody: `{"message":"Доступ запрещен: неверная роль"}`,
		},
		{
			name:      "No role in context",
			inputBody: fmt.Sprintf(`{"type":"clothing", "pvzId":"%s"}`, testPvzId),
			role:      "",
			mockBehavior: func(s *mock_services.MockProduct, prodType models.Type, pvzId uuid.UUID) {},
			expectedStatusCode:   http.StatusForbidden,
			expectedResponseBody: `{"message":"Доступ запрещен: неверная роль"}`,
		},
		{
			name:      "Invalid JSON",
			inputBody: `{invalid_json}`,
			role:      models.RoleEmployee,
			mockBehavior: func(s *mock_services.MockProduct, prodType models.Type, pvzId uuid.UUID) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"Неверный запрос"}`,
		},
		{
			name:      "Internal service error",
			inputBody: fmt.Sprintf(`{"type":"одежда", "pvzId":"%s"}`, testPvzId),
			role:      models.RoleEmployee,
			mockBehavior: func(s *mock_services.MockProduct, prodType models.Type, pvzId uuid.UUID) {
				s.EXPECT().AddProduct(prodType, pvzId).Return(models.Product{}, fmt.Errorf("error"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"Не удалось добавить товар"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockProduct := mock_services.NewMockProduct(ctrl)
			testCase.mockBehavior(mockProduct, models.Clothes, testPvzId)

			services := &services.Service{Product: mockProduct}
			handler := &Handler{services: services}

			r := chi.NewRouter()
			r.Post("/products", func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), middleware.ContextKeyRole, testCase.role)
				handler.AddProduct(w, r.WithContext(ctx))
			})

			req := httptest.NewRequest("POST", "/products", bytes.NewBufferString(testCase.inputBody))
			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)

			assert.Equal(t, testCase.expectedStatusCode, rr.Code)
			assert.JSONEq(t, testCase.expectedResponseBody, rr.Body.String())
		})
	}
}

func TestHandler_DeleteProduct(t *testing.T) {
	type mockBehavior func(s *mock_services.MockProduct, pvzId uuid.UUID)

	testPvzId := uuid.New()

	testTable := []struct {
		name                 string
		pvzIdParam           string
		role                 models.Role
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:       "OK",
			pvzIdParam: testPvzId.String(),
			role:       models.RoleEmployee,
			mockBehavior: func(s *mock_services.MockProduct, pvzId uuid.UUID) {
				s.EXPECT().DeleteProduct(pvzId).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: ``,
		},
		{
			name:       "Forbidden: no role",
			pvzIdParam: testPvzId.String(),
			role:       "",
			mockBehavior: func(s *mock_services.MockProduct, pvzId uuid.UUID) {},
			expectedStatusCode:   http.StatusForbidden,
			expectedResponseBody: `{"message":"Доступ запрещен: неверная роль"}`,
		},
		{
			name:       "Forbidden: wrong role",
			pvzIdParam: testPvzId.String(),
			role:       models.RoleModerator,
			mockBehavior: func(s *mock_services.MockProduct, pvzId uuid.UUID) {},
			expectedStatusCode:   http.StatusForbidden,
			expectedResponseBody: `{"message":"Доступ запрещен: неверная роль"}`,
		},
		{
			name:       "DeleteProduct returns error",
			pvzIdParam: testPvzId.String(),
			role:       models.RoleEmployee,
			mockBehavior: func(s *mock_services.MockProduct, pvzId uuid.UUID) {
				s.EXPECT().DeleteProduct(pvzId).Return(fmt.Errorf("product not found"))
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"product not found"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			productService := mock_services.NewMockProduct(c)
			testPvzIdParsed, _ := uuid.Parse(testCase.pvzIdParam)
			testCase.mockBehavior(productService, testPvzIdParsed)

			services := &services.Service{Product: productService}
			handler := &Handler{services: services}

			r := chi.NewRouter()
			r.Delete("/products/{pvzId}", func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), middleware.ContextKeyRole, testCase.role)
				handler.DeleteProduct(w, r.WithContext(ctx))
			})

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/products/%s", testCase.pvzIdParam), nil)
			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)

			assert.Equal(t, testCase.expectedStatusCode, rr.Code)
			if testCase.expectedResponseBody != `` {
				assert.JSONEq(t, testCase.expectedResponseBody, rr.Body.String())
			} else {
				assert.Empty(t, rr.Body.String())
			}
		})
	}
}

