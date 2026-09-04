package handlers

import (
    // "net/http"
    // "net/http/httptest" // Импортируем инструмент для HTTP-тестов
    // "testing"
	// "bytes"
	"context"
	"avito/internal/model"
// 	"strings"
// 	"avito/internal/service"
// 	"errors"
)
type MockCourierService struct {
	ReturnCourier model.Courier
	ReturnError   error
	CalledID int
	CalledName      string
	CalledPhone     string
	CalledStatus    string
	CalledTransport string
	CallCount       int
}

func (m *MockCourierService) CreateCourier(ctx context.Context, name, phone, status, transport string) (model.Courier, error) {
	m.CallCount++
	m.CalledName, m.CalledPhone = name, phone
	m.CalledStatus, m.CalledTransport = status, transport
	return m.ReturnCourier, m.ReturnError
}

func (m *MockCourierService) GetAllCourier(ctx context.Context) ([]model.Courier, error){
	return nil, nil 
}

func (m *MockCourierService) GetById(ctx context.Context, id int) (model.Courier, error) {
	m.CallCount++
	m.CalledID = id
	return m.ReturnCourier, m.ReturnError
}


func (m *MockCourierService) UpdateCourier(ctx context.Context, id int, name, phone, status, transport string) (model.Courier, error) {
	return m.ReturnCourier, m.ReturnError
}

func (m *MockCourierService) DeleteCourier(ctx context.Context, id int) error {
	return m.ReturnError
}
