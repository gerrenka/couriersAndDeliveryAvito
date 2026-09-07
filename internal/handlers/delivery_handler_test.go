package handlers

import (
    "net/http"
    "net/http/httptest"
    "testing"

	"avito/internal/model"
	"strings"
	"errors"
)

func TestDeliveryHandler_Assign(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		body          string
		mockErr       error
		wantCode      int
		wantCallCount int
		wantOrderID   string
	}{
		{
			name:          "успех",
			body:          `{"order_id":"order-1"}`,
			wantCode:      http.StatusOK,
			wantCallCount: 1,
			wantOrderID:   "order-1",
		},
		{
			name:          "битый json",
			body:          `{`,
			wantCode:      http.StatusBadRequest,
			wantCallCount: 0,
		},
		{
			name:          "нет свободных курьеров",
			body:          `{"order_id":"order-1"}`,
			mockErr:       model.ErrCourierNotAvailable,
			wantCode:      http.StatusConflict,
			wantCallCount: 1,
		},
		{
			name:          "ошибка бд",
			body:          `{"order_id":"order-1"}`,
			mockErr:       errors.New("connection refused"),
			wantCode:      http.StatusInternalServerError,
			wantCallCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUsecase := &MockDeliveryUsecase{
				ReturnAssign: model.AssignResponse{CourierID: 1, OrderID: "order-1"},
				ReturnError:  tt.mockErr,
			}
			dh := &DeliveryHandler{use: mockUsecase}

			req := httptest.NewRequest(http.MethodPost, "/assign", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			dh.Assign(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("Ожидался %d, получили %d, тело: %s", tt.wantCode, rr.Code, rr.Body.String())
			}
			if mockUsecase.CallCount != tt.wantCallCount {
				t.Errorf("Ожидали %d вызовов, было %d", tt.wantCallCount, mockUsecase.CallCount)
			}
			if tt.wantOrderID != "" && mockUsecase.CalledOrderID != tt.wantOrderID {
				t.Errorf("В usecase ушёл orderID %q вместо %q", mockUsecase.CalledOrderID, tt.wantOrderID)
			}
		})
	}
}

func TestDeliveryHandler_Unassign(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		body          string
		mockErr       error
		wantCode      int
		wantCallCount int
		wantOrderID   string
	}{
		{
			name:          "успех",
			body:          `{"order_id":"order-1"}`,
			wantCode:      http.StatusOK,
			wantCallCount: 1,
			wantOrderID:   "order-1",
		},
		{
			name:          "битый json",
			body:          `{`,
			wantCode:      http.StatusBadRequest,
			wantCallCount: 0,
		},
		{
			name:          "заказ не найден",
			body:          `{"order_id":"order-1"}`,
			mockErr:       model.ErrCourierNotOrder,
			wantCode:      http.StatusNotFound,
			wantCallCount: 1,
		},
		{
			name:          "ошибка бд",
			body:          `{"order_id":"order-1"}`,
			mockErr:       errors.New("connection refused"),
			wantCode:      http.StatusInternalServerError,
			wantCallCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockUsecase := &MockDeliveryUsecase{
				ReturnUnassign: model.UnassignResponse{CourierID: 1, OrderID: "order-1"},
				ReturnError:  tt.mockErr,
			}
			dh := &DeliveryHandler{use: mockUsecase}

			req := httptest.NewRequest(http.MethodPost, "/unassign", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			dh.Unassign(rr, req)

			if rr.Code != tt.wantCode {
				t.Errorf("Ожидался %d, получили %d, тело: %s", tt.wantCode, rr.Code, rr.Body.String())
			}
			if mockUsecase.CallCount != tt.wantCallCount {
				t.Errorf("Ожидали %d вызовов, было %d", tt.wantCallCount, mockUsecase.CallCount)
			}
			if tt.wantOrderID != "" && mockUsecase.CalledOrderID != tt.wantOrderID {
				t.Errorf("В usecase ушёл orderID %q вместо %q", mockUsecase.CalledOrderID, tt.wantOrderID)
			}
		})
	}
}