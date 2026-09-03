package factory

import (
	"errors"
	"time"
)

type DeliveryTimeFactory struct {
}

func (d DeliveryTimeFactory) CalculateDeadline(transportType string) (time.Time, error) {
	switch transportType {
	case "on_foot":
		return time.Now().Add(time.Minute * 30), nil
	case "scooter":
		return time.Now().Add(time.Minute * 15), nil
	case "car":
		return time.Now().Add(time.Minute * 5), nil
	default:
		return time.Time{}, errors.New("неизвестный тип транспорта")
	}
}
