package repository

import (
	"avito/internal/model"
	"context"
	"github.com/jackc/pgx/v5"
	"errors"
	"time"
)

type DeliveryRepository struct{

}

func NewDeliveryRepository () *DeliveryRepository{
	return &DeliveryRepository{}
}


func (repos *DeliveryRepository) CreateDelivery(ctx context.Context, tx pgx.Tx, courierID int, orderID string, deadline time.Time) (error){
	query := `INSERT INTO delivery (courier_id, order_id, deadline) VALUES($1, $2, $3)`
	_, err  := tx.Exec(ctx, query, courierID, orderID, deadline)
	return err

}

func (repos *DeliveryRepository) FindByOrderID(ctx context.Context, tx pgx.Tx, orderID string) (model.Delivery, error){
	var del model.Delivery
	query := `SELECT id, courier_id, order_id, assigned_at, deadline FROM delivery where order_ID =$1`
	err  := tx.QueryRow(ctx, query, orderID).Scan(&del.ID, &del.CourierID, &del.OrderID, &del.AssignedAt, &del.Deadline)
	if errors.Is(err, pgx.ErrNoRows) {
    	return model.Delivery{}, model.ErrCourierNotOrder
	}
	if err != nil {
		return model.Delivery{}, err
	}
	return del, nil
}

func (repos *DeliveryRepository) UpdateStatusByOrderId(ctx context.Context, tx pgx.Tx, status, orderID string) error {
	cmd, err := tx.Exec(ctx, "UPDATE delivery SET status = $1 WHERE order_id = $2 AND AND status = 'active", status, orderID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
    return model.ErrCourierNotOrder
}
	return nil
}

