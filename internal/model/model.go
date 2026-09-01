package model

import (
	"time"
	"errors"
)

var ErrCourierNotFound = errors.New("Курьер не найден")
var ErrCourierNotAvailable = errors.New("Нет свободных курьеров")
var ErrCourierNotOrder = errors.New("Связь курьера с заказом не найдена,")
var ErrPhoneExists = errors.New("курьер с таким телефоном уже существует")
type CreateCourierRequest struct {
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Status string `json:"status"`
	TransportType  string `json:"transport_type"`
}

type Courier struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TransportType  string `json:"transport_type"`
}
 

type Delivery struct{
	ID            int       `json:"id"`
    CourierID    int       `json:"courier_id"`
    OrderID      string    `json:"order_id"`
    AssignedAt   time.Time `json:"assigned_at"`
    Deadline      time.Time `json:"deadline"`

}

type OrderIDRequest struct { 
	OrderID string  `json:"order_id"`
}

type AssignResponse struct {
	CourierID   int       	`json:"courier_id"`
  	OrderID 	string  	`json:"order_id"`
 	TransportType  string `json:"transport_type"`
  	Deadline      time.Time `json:"delivery_deadline"`
}

type UnassignResponse struct {
	OrderID 	string  	`json:"order_id"`
	Status    string    `json:"status"`
	CourierID   int       	`json:"courier_id"`
}