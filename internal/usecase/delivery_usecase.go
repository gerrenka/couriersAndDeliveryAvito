package usecase

import (
	"avito/internal/factory"
	"avito/internal/model"

	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CourierRepository interface{
	UpdateStatus (ctx context.Context, tx pgx.Tx, id int, status string) (error)
	FindAvailable (ctx context.Context, tx pgx.Tx) (model.Courier, error) 
}

type DeliveryRepository interface{
	CreateDelivery(ctx context.Context, tx pgx.Tx, courierID int, orderID string, deadline time.Time) (error)
	FindByOrderID(ctx context.Context, tx pgx.Tx, orderID string) (model.Delivery, error)
	DeleteByOrderId(ctx context.Context, tx pgx.Tx, orderID string) error
}

type DeliveryUsecase struct{
	pool *pgxpool.Pool
	courierRepo CourierRepository
	deliveryRepo DeliveryRepository
	factory factory.DeliveryTimeFactory
}

func NewDeliveryUsecase(pool *pgxpool.Pool,  courierRepo CourierRepository, deliveryRepo DeliveryRepository, factory factory.DeliveryTimeFactory) *DeliveryUsecase{
	return &DeliveryUsecase{
		pool: pool,
		courierRepo: courierRepo,
		deliveryRepo: deliveryRepo,
		factory: factory,
	}
}

func (u *DeliveryUsecase) Assign(ctx context.Context, orderID string) (model.AssignResponse, error){
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return model.AssignResponse{}, err
	}
	defer tx.Rollback(ctx)

	courier, err:=u.courierRepo.FindAvailable(ctx, tx)
	if err != nil {
    return model.AssignResponse{}, err
}
	deadline, err:=u.factory.CalculateDeadline(courier.TransportType)
	if err != nil {
    return model.AssignResponse{}, err
}
	err = u.courierRepo.UpdateStatus(ctx, tx, courier.ID, "busy")
	if err != nil {
    return model.AssignResponse{}, err
}
	err = u.deliveryRepo.CreateDelivery(ctx, tx, courier.ID, orderID, deadline)
	if err != nil {
    return model.AssignResponse{}, err
}

	err = tx.Commit(ctx)
	return model.AssignResponse{CourierID: courier.ID, OrderID: orderID, TransportType: courier.TransportType, Deadline: deadline}, err
}



func (u *DeliveryUsecase) Unassign(ctx context.Context, orderID string) (model.UnassignResponse, error){
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return model.UnassignResponse{}, err
	}
	defer tx.Rollback(ctx)
	delivery, err := u.deliveryRepo.FindByOrderID(ctx, tx, orderID)
		if err != nil {
		return model.UnassignResponse{}, err
	}

	err = u.deliveryRepo.DeleteByOrderId(ctx, tx, orderID)
		if err != nil {
		return model.UnassignResponse{}, err
	}
	err = u.courierRepo.UpdateStatus(ctx, tx, delivery.CourierID, "available")
		if err != nil {
		return model.UnassignResponse{}, err
	}
	
	err = tx.Commit(ctx)
	return model.UnassignResponse{OrderID: orderID, Status: "unassigned", CourierID: delivery.CourierID}, err
}
