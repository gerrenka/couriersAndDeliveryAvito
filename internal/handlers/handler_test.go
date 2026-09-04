package handlers

import (
    "net/http"
    "net/http/httptest"
    "testing"
	"bytes"
	// "context"
	"avito/internal/model"
	"strings"
	"avito/internal/service"
	"errors"
	"github.com/go-chi/chi/v5"
)

func TestHealthcheck(t *testing.T) {
    req := httptest.NewRequest("GET", "/healthcheck", nil)
    rr := httptest.NewRecorder()

    crh := &CourierHandler{}
    crh.Healthcheck(rr, req)

    if rr.Code != http.StatusNoContent {
        t.Errorf("Ожидался статус 204, но получили %d", rr.Code)
    }

    if rr.Body.String() != "" {
        t.Errorf("Ожидался пустой ответ, но получили %s", rr.Body.String())
    }
}


func TestCreate(t *testing.T) {
    jsonBody := []byte(`{"name":"Иван","phone":"+79001234567","status":"available","transport_type":"bike"}`)
    
    req := httptest.NewRequest("POST", "/courier", bytes.NewBuffer(jsonBody))
    req.Header.Set("Content-Type", "application/json") 
    
    rr := httptest.NewRecorder()

    mockService := &MockCourierService{
    ReturnCourier: model.Courier{ID: 1, Name: "Иван"},
}
    crh := &CourierHandler{
        serv: mockService, 
    }

    crh.Create(rr, req)

    if rr.Code != http.StatusCreated { 
        t.Errorf("Ожидался статус 201, но получили %d", rr.Code)
    }

    if mockService.CalledName != "Иван" {
        t.Errorf("Сервис не был вызван с именем Иван")
    }
}

func TestCreate_ServiceValidationError(t *testing.T) {
	t.Parallel()

	body := `{"name":"Иван","phone":"+79001234567","status":"летит","transport_type":"bike"}`
	req := httptest.NewRequest(http.MethodPost, "/courier", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mockService := &MockCourierService{ReturnError: service.ErrInvalidStatus}
	crh := &CourierHandler{serv: mockService}

	crh.Create(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Ожидался статус 400, получили %d", rr.Code)
	}
}


func TestCreate_BadJSON(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, "/courier", strings.NewReader("{"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mockService := &MockCourierService{}
	crh := &CourierHandler{serv: mockService}

	crh.Create(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Ожидался статус 400, получили %d", rr.Code)
	}

	if mockService.CallCount != 0 {
		t.Errorf("Сервис не должен вызываться при битом JSON, вызван %d раз", mockService.CallCount)
	}
}
func TestCreate_ServiceInternalError(t *testing.T) {
	t.Parallel()

	body := `{"name":"Иван","phone":"+79001234567","status":"available","transport_type":"bike"}`
	req := httptest.NewRequest(http.MethodPost, "/courier", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mockService := &MockCourierService{ReturnError: errors.New("connection refused")}
	crh := &CourierHandler{serv: mockService}

	crh.Create(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус 500, получили %d", rr.Code)
	}
}


func TestGetById(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/courier/123", nil)
	rr := httptest.NewRecorder()
	mockService := &MockCourierService{
		ReturnCourier: model.Courier{ID: 123, Name: "Иван"},
	}
	crh := &CourierHandler{serv: mockService}

	router := chi.NewRouter()          
	router.Get("/courier/{id}", crh.GetById)        
	router.ServeHTTP(rr, req)                       

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получили %d", rr.Code)
	}
}

func TestGetById_InvalidIDFormat (t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/courier/abc", nil)
	rr := httptest.NewRecorder()
	mockService := &MockCourierService{
	}
	crh := &CourierHandler{serv: mockService}

	router := chi.NewRouter()          
	router.Get("/courier/{id}", crh.GetById)        
	router.ServeHTTP(rr, req)                       
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Ожидался 400, получили %d", rr.Code)
	}
	if mockService.CallCount != 0 {
		t.Errorf("Сервис не должен вызываться, вызван %d раз", mockService.CallCount)
	}
}