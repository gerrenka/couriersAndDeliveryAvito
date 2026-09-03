package service

import (
	"avito/internal/model"
	"context"
	"errors"
)

type CourierRepository interface {
	GetById(ctx context.Context, id int) (model.Courier, error)
	GetAllCourier(ctx context.Context) ([]model.Courier, error)
	CreateCourier(ctx context.Context, name, phone, status, transport_type string) (model.Courier, error)
	UpdateCourier(ctx context.Context, id int, name, phone, status, transport_type string) (model.Courier, error)
	DeleteCourier(ctx context.Context, id int) error
}

type CourierService struct {
	repo CourierRepository
}

func NewCourierService(repo CourierRepository) *CourierService {
	return &CourierService{
		repo: repo,
	}
}

var (
	ErrInvalidStatus = errors.New("невалидный статус")
	ErrInvalidName   = errors.New("невалидный имя")
	ErrInvalidPhone  = errors.New("невалидный телефон")
)

func (s *CourierService) CreateCourier(ctx context.Context, name, phone, status, transport_type string) (model.Courier, error) {
	switch status {
	case "available", "busy", "paused":
	default:
		return model.Courier{}, ErrInvalidStatus
	}
	if name == "" {
		return model.Courier{}, ErrInvalidName
	}
	if phone == "" {
		return model.Courier{}, ErrInvalidPhone
	}
	if transport_type == "" {
		transport_type = "on_foot"
	}
	return s.repo.CreateCourier(ctx, name, phone, status, transport_type)
}

func (s *CourierService) GetById(ctx context.Context, id int) (model.Courier, error) {
	return s.repo.GetById(ctx, id)
}

func (s *CourierService) GetAllCourier(ctx context.Context) ([]model.Courier, error) {
	return s.repo.GetAllCourier(ctx)
}

func (s *CourierService) DeleteCourier(ctx context.Context, id int) error {
	return s.repo.DeleteCourier(ctx, id)
}

func (s *CourierService) UpdateCourier(ctx context.Context, id int, name, phone, status, transport_type string) (model.Courier, error) {
	switch status {
	case "available", "busy", "paused":
	default:
		return model.Courier{}, ErrInvalidStatus
	}
	if name == "" {
		return model.Courier{}, ErrInvalidName
	}
	if phone == "" {
		return model.Courier{}, ErrInvalidPhone
	}
	return s.repo.UpdateCourier(ctx, id, name, phone, status, transport_type)
}
