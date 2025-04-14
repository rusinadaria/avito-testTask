package tests

import (
	"avito-testTask/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func (s *APITestSuite) TestAuth() {

	user := map[string]string{
		"role": "moderator",
	}

	requestBody, err := json.Marshal(user)
	assert.NoError(s.T(), err)

	req, err := http.NewRequest("POST", "/dummyLogin", bytes.NewBuffer(requestBody))
	assert.NoError(s.T(), err)

	req.Header.Set("Content-Type", "application/json")

	// rr := httptest.NewRecorder()
	// s.r.ServeHTTP(rr, req)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handler.DummyLogin)
	handler.ServeHTTP(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code)

	var resp struct {
		Token string `json:"token"`
	}

	err = json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), resp.Token)

	parsedToken, err := jwt.Parse(resp.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte("qrkjk#4#%35FSFJlja#4353KSFjH"), nil
	})
	assert.NoError(s.T(), err)
	assert.True(s.T(), parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), "moderator", claims["role"])
}

func (s *APITestSuite) TestCreatePVZ() {

	loginPayload := map[string]string{
		"role": "employee",
	}
	loginBody, err := json.Marshal(loginPayload)
	assert.NoError(s.T(), err)

	loginReq, err := http.NewRequest("POST", "/dummyLogin", bytes.NewBuffer(loginBody))
	assert.NoError(s.T(), err)
	loginReq.Header.Set("Content-Type", "application/json")

	loginRR := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handler.DummyLogin)
	handler.ServeHTTP(loginRR, loginReq)

	assert.Equal(s.T(), http.StatusOK, loginRR.Code)

	var loginResp struct {
		Token string `json:"token"`
	}
	err = json.NewDecoder(loginRR.Body).Decode(&loginResp)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), loginResp.Token)


	id := uuid.MustParse("28b403b3-c7b0-4cb9-9c42-1f5d0cc2cdd4")
	registrationDate, err := time.Parse(time.RFC3339, "2025-04-11T14:30:00Z")
	assert.NoError(s.T(), err)

	pvz := models.PVZ{
		Id:               id,
		RegistrationDate: registrationDate,
		City:             "Москва",
	}
	requestBody, err := json.Marshal(pvz)
	assert.NoError(s.T(), err)

	req, err := http.NewRequest("POST", "/pvz", bytes.NewBuffer(requestBody))
	assert.NoError(s.T(), err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)

	rr := httptest.NewRecorder()
	handler = http.HandlerFunc(s.handler.PVZCreate)
	handler.ServeHTTP(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code)
}

func (s *APITestSuite) TestCreateReception() {
	loginPayload := map[string]string{
		"role": "employee",
	}
	loginBody, err := json.Marshal(loginPayload)
	assert.NoError(s.T(), err)

	loginReq, err := http.NewRequest("POST", "/dummyLogin", bytes.NewBuffer(loginBody))
	assert.NoError(s.T(), err)
	loginReq.Header.Set("Content-Type", "application/json")

	loginRR := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handler.DummyLogin)
	handler.ServeHTTP(loginRR, loginReq)

	assert.Equal(s.T(), http.StatusOK, loginRR.Code)

	var loginResp struct {
		Token string `json:"token"`
	}
	err = json.NewDecoder(loginRR.Body).Decode(&loginResp)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), loginResp.Token)

	pvzId := uuid.MustParse("28b403b3-c7b0-4cb9-9c42-1f5d0cc2cdd4")

	receptionRequest := map[string]string{
		"pvzId": pvzId.String(),
	}

	receptionRequestBody, err := json.Marshal(receptionRequest)
	assert.NoError(s.T(), err)

	req, err := http.NewRequest("POST", "/receptions", bytes.NewBuffer(receptionRequestBody))
	assert.NoError(s.T(), err)
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler = http.HandlerFunc(s.handler.Receptions)
	handler.ServeHTTP(rr, req)

	assert.Equal(s.T(), http.StatusCreated, rr.Code)

	var createdReception models.Reception
	err = json.NewDecoder(rr.Body).Decode(&createdReception)
	assert.NoError(s.T(), err)

	assert.NotEmpty(s.T(), createdReception.Id)
	assert.Equal(s.T(), pvzId, createdReception.PvzId)
	assert.NotEmpty(s.T(), createdReception.Status)

	createdDateTime := createdReception.DateTime.Format(time.RFC3339)
	_, err = time.Parse(time.RFC3339, createdDateTime)
	assert.NoError(s.T(), err)
}

func (s *APITestSuite) getTokenWithRole(role string) string {
	payload := map[string]string{"role": role}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/dummyLogin", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handler.DummyLogin)
	handler.ServeHTTP(rr, req)

	var resp struct {
		Token string `json:"token"`
	}
	_ = json.NewDecoder(rr.Body).Decode(&resp)

	return resp.Token
}

type RequestProduct struct {
	Type models.Type `json:"type"`
	PvzId uuid.UUID `json:"pvzId"`
}

func (s *APITestSuite) TestAddProducts() {
	// mockProduct := models.Product{
	// 	Id: uuid.New(),
	// 	DateTime: time.Now(),
	// 	Type: models.Type("product_type"),
	// 	ReceptionId: uuid.New(),
	// }

	pvzId := uuid.MustParse("28b403b3-c7b0-4cb9-9c42-1f5d0cc2cdd4")

	// s.mockService.
	// 	EXPECT().
	// 	AddProduct(models.Type("product_type"), pvzId).
	// 	Return(mockProduct, nil).
	// 	Times(50)


	token := s.getTokenWithRole("employee")
	fmt.Println(token)

	for i := 0; i < 50; i++ {
		productRequest := RequestProduct{
			Type: models.Type("product_type"),
			PvzId: pvzId,
		}
		body, err := json.Marshal(productRequest)
		assert.NoError(s.T(), err)

		req, err := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
		assert.NoError(s.T(), err)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(s.handler.AddProduct)
		handler.ServeHTTP(rr, req)

		assert.Equal(s.T(), http.StatusCreated, rr.Code)
	}
}

func (s *APITestSuite) TestCloseReception() {
	loginPayload := map[string]string{
		"role": "employee",
	}
	loginBody, err := json.Marshal(loginPayload)
	assert.NoError(s.T(), err)

	loginReq, err := http.NewRequest("POST", "/dummyLogin", bytes.NewBuffer(loginBody))
	assert.NoError(s.T(), err)
	loginReq.Header.Set("Content-Type", "application/json")

	loginRR := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handler.DummyLogin)
	handler.ServeHTTP(loginRR, loginReq)

	assert.Equal(s.T(), http.StatusOK, loginRR.Code)

	var loginResp struct {
		Token string `json:"token"`
	}
	err = json.NewDecoder(loginRR.Body).Decode(&loginResp)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), loginResp.Token)

	pvzId := uuid.MustParse("28b403b3-c7b0-4cb9-9c42-1f5d0cc2cdd4")

	receptionRequest := map[string]string{
		"pvzId": pvzId.String(),
	}
	receptionRequestBody, err := json.Marshal(receptionRequest)
	assert.NoError(s.T(), err)

	req, err := http.NewRequest("POST", "/receptions", bytes.NewBuffer(receptionRequestBody))
	assert.NoError(s.T(), err)
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)
	req.Header.Set("Content-Type", "application/json")


	rr := httptest.NewRecorder()
	handler = http.HandlerFunc(s.handler.Receptions)
	handler.ServeHTTP(rr, req)


	assert.Equal(s.T(), http.StatusCreated, rr.Code)

	var createdReception models.Reception
	err = json.NewDecoder(rr.Body).Decode(&createdReception)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), createdReception.Id)
	assert.Equal(s.T(), pvzId, createdReception.PvzId)


	closeReq, err := http.NewRequest("POST", "/receptions/close/"+createdReception.PvzId.String(), nil)
	assert.NoError(s.T(), err)
	closeReq.Header.Set("Authorization", "Bearer "+loginResp.Token)

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(s.handler.CloseReception)
	handler.ServeHTTP(rr, closeReq)

	assert.Equal(s.T(), http.StatusOK, rr.Code)

	var closedReception models.Reception
	err = json.NewDecoder(rr.Body).Decode(&closedReception)
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), closedReception.Id)
	assert.Equal(s.T(), createdReception.PvzId, closedReception.PvzId)
	assert.Equal(s.T(), models.Close, closedReception.Status)
}
