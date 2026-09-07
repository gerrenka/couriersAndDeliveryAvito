package handlers

import (
	"context"
	"avito/internal/model"
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
	ReturnCouriers  []model.Courier
}

type MockDeliveryUsecase struct {
	ReturnAssign model.AssignResponse
	ReturnUnassign model.UnassignResponse
	ReturnError   error

	CalledOrderID string
	CallCount     int
} 

func (m *MockCourierService) CreateCourier(ctx context.Context, name, phone, status, transport string) (model.Courier, error) {
	m.CallCount++
	m.CalledName, m.CalledPhone = name, phone
	m.CalledStatus, m.CalledTransport = status, transport
	return m.ReturnCourier, m.ReturnError
}

func (m *MockCourierService) GetAllCourier(ctx context.Context) ([]model.Courier, error) {
	m.CallCount++
	return m.ReturnCouriers, m.ReturnError
}

func (m *MockCourierService) GetById(ctx context.Context, id int) (model.Courier, error) {
	m.CallCount++
	m.CalledID = id
	return m.ReturnCourier, m.ReturnError
}


func (m *MockCourierService) DeleteCourier(ctx context.Context, id int) error {
	m.CallCount++
	m.CalledID = id
	return m.ReturnError
}

func (m *MockCourierService) UpdateCourier(ctx context.Context, id int, name, phone, status, transport string) (model.Courier, error) {
	m.CallCount++
	m.CalledID = id
	m.CalledName, m.CalledPhone = name, phone
	m.CalledStatus, m.CalledTransport = status, transport
	return m.ReturnCourier, m.ReturnError
}

func (m *MockDeliveryUsecase) Assign(ctx context.Context, orderID string) (model.AssignResponse, error) {
	m.CallCount++
	m.CalledOrderID = orderID
	return m.ReturnAssign, m.ReturnError
}

func (m *MockDeliveryUsecase) Unassign(ctx context.Context, orderID string) (model.UnassignResponse, error) {
	m.CallCount++
	m.CalledOrderID = orderID
	return m.ReturnUnassign, m.ReturnError
}