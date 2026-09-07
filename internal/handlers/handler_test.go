package handlers

import (
    "net/http"
    "net/http/httptest"
    "testing"

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


func TestCourierHandler_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		body          string
		mockErr       error
		wantCode      int
		wantCallCount int
	}{
		{
			name:          "успех",
			body:          `{"name":"Иван","phone":"+79001234567","status":"available","transport_type":"bike"}`,
			wantCode:      http.StatusCreated,
			wantCallCount: 1,
		},
		{
			name:          "битый json",
			body:          `{`,
			wantCode:      http.StatusBadRequest,
			wantCallCount: 0,
		},
		{
			name:          "невалидный статус",
			body:          `{"name":"Иван","phone":"+79001234567","status":"летит"}`,
			mockErr:       service.ErrInvalidStatus,
			wantCode:      http.StatusBadRequest,
			wantCallCount: 1,
		},
		{
			name:          "ошибка бд",
			body:          `{"name":"Иван","phone":"+79001234567","status":"available"}`,
			mockErr:       errors.New("connection refused"),
			wantCode:      http.StatusInternalServerError,
			wantCallCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/courier", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			mockService := &MockCourierService{
				ReturnCourier: model.Courier{ID: 1, Name: "Иван"},
				ReturnError:   tt.mockErr,
			}
			crh := &CourierHandler{serv: mockService}

			crh.Create(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("Ожидался %d, получили %d", tt.wantCode, rr.Code)
			}
			if mockService.CallCount != tt.wantCallCount {
				t.Errorf("Ожидали %d вызовов, было %d", tt.wantCallCount, mockService.CallCount)
			}
		})
	}
}


func TestCourierHandler_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		id            string
		body          string
		mockErr       error
		wantCode      int
		wantCallCount int
	}{
		{
			name:          "успех",
			id:            "123",
			body:          `{"name":"Иван","phone":"+79001234567","status":"available","transport_type":"bike"}`,
			wantCode:      http.StatusOK,
			wantCallCount: 1,
		},
		{
			name:          "нечисловой id",
			id:            "abc",
			body:          `{"name":"Иван","phone":"+79001234567","status":"available"}`,
			wantCode:      http.StatusBadRequest,
			wantCallCount: 0,
		},
		{
			name:          "битый json",
			id:            "123",
			body:          `{`,
			wantCode:      http.StatusBadRequest,
			wantCallCount: 0,
		},
		{
			name:          "невалидный статус",
			id:            "123",
			body:          `{"name":"Иван","phone":"+79001234567","status":"летит"}`,
			mockErr:       service.ErrInvalidStatus,
			wantCode:      http.StatusBadRequest,
			wantCallCount: 1,
		},
		{
			name:          "ошибка бд",
			id:            "123",
			body:          `{"name":"Иван","phone":"+79001234567","status":"available"}`,
			mockErr:       errors.New("connection refused"),
			wantCode:      http.StatusInternalServerError,
			wantCallCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockService := &MockCourierService{
				ReturnCourier: model.Courier{ID: 123, Name: "Иван"},
				ReturnError:   tt.mockErr,
			}
			crh := &CourierHandler{serv: mockService}

			router := chi.NewRouter()
			router.Put("/courier/{id}", crh.Update)

			req := httptest.NewRequest(http.MethodPut, "/courier/"+tt.id, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("Ожидался %d, получили %d, тело: %s", tt.wantCode, rr.Code, rr.Body.String())
			}
			if mockService.CallCount != tt.wantCallCount {
				t.Errorf("Ожидали %d вызовов сервиса, было %d", tt.wantCallCount, mockService.CallCount)
			}
		})
	}
}

func TestCourierHandler_GetById(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		id            string
		mockErr       error
		wantCode      int
		wantCallCount int
	}{
		{
			name:          "успех",
			id:            "123",
			wantCode:      http.StatusOK,
			wantCallCount: 1,
		},
		{
			name:          "нечисловой id",
			id:            "abc",
			wantCode:      http.StatusBadRequest,
			wantCallCount: 0,
		},
		{
			name:          "курьер не найден",
			id:            "999",
			mockErr:       model.ErrCourierNotFound,
			wantCode:      http.StatusNotFound,
			wantCallCount: 1,
		},
		{
			name:          "ошибка бд",
			id:            "123",
			mockErr:       errors.New("connection refused"),
			wantCode:      http.StatusInternalServerError,
			wantCallCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockService := &MockCourierService{
				ReturnCourier: model.Courier{ID: 123, Name: "Иван"},
				ReturnError:   tt.mockErr,
			}
			crh := &CourierHandler{serv: mockService}

			router := chi.NewRouter()
			router.Get("/courier/{id}", crh.GetById)

			req := httptest.NewRequest(http.MethodGet, "/courier/"+tt.id, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("Ожидался %d, получили %d, тело: %s", tt.wantCode, rr.Code, rr.Body.String())
			}
			if mockService.CallCount != tt.wantCallCount {
				t.Errorf("Ожидали %d вызовов, было %d", tt.wantCallCount, mockService.CallCount)
			}
		})
	}
}

func TestCourierHandler_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		id            string
		mockErr       error
		wantCode      int
		wantCallCount int
	}{
		{
			name:          "успех",
			id:            "123",
			wantCode:      http.StatusNoContent,
			wantCallCount: 1,
		},
		{
			name:          "нечисловой id",
			id:            "abc",
			wantCode:      http.StatusBadRequest,
			wantCallCount: 0,
		},
		{
			name:          "курьер не найден",
			id:            "999",
			mockErr:       model.ErrCourierNotFound,
			wantCode:      http.StatusNotFound,
			wantCallCount: 1,
		},
		{
			name:          "ошибка бд",
			id:            "123",
			mockErr:       errors.New("connection refused"),
			wantCode:      http.StatusInternalServerError,
			wantCallCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockService := &MockCourierService{ReturnError: tt.mockErr}
			crh := &CourierHandler{serv: mockService}

			router := chi.NewRouter()
			router.Delete("/courier/{id}", crh.Delete)

			req := httptest.NewRequest(http.MethodDelete, "/courier/"+tt.id, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("Ожидался %d, получили %d, тело: %s", tt.wantCode, rr.Code, rr.Body.String())
			}
			if mockService.CallCount != tt.wantCallCount {
				t.Errorf("Ожидали %d вызовов, было %d", tt.wantCallCount, mockService.CallCount)
			}
			if tt.wantCode == http.StatusNoContent && rr.Body.String() != "" {
				t.Errorf("При 204 тело должно быть пустым, получили %s", rr.Body.String())
			}
		})
	}
}

func TestCourierHandler_List(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		couriers      []model.Courier
		mockErr       error
		wantCode      int
	}{
		{
			name:          "успех",
			couriers: []model.Courier{{ID: 1, Name: "Иван"}, {ID: 2, Name: "Пётр"}},
			wantCode: http.StatusOK,
		},
		{
			name:     "пустой список",
			couriers: []model.Courier{},
			wantCode: http.StatusOK,
		},
		{
			name:     "ошибка бд",
			mockErr:  errors.New("connection refused"),
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockService := &MockCourierService{ReturnCouriers: tt.couriers, ReturnError: tt.mockErr}
			crh := &CourierHandler{serv: mockService}

			req := httptest.NewRequest(http.MethodGet, "/courier", nil)
			rr := httptest.NewRecorder()

			crh.List(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("Ожидался %d, получили %d", tt.wantCode, rr.Code)
			}
		})
	}
}