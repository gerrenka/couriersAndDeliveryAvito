package handlers

import (
	"avito/internal/model"
	"context"
	"encoding/json"
	"errors"
	"net/http"

)

type DeliveryUsecase interface {
	Assign(ctx context.Context, orderID string) (model.AssignResponse, error)
	Unassign(ctx context.Context, orderID string) (model.UnassignResponse, error)
}

type DeliveryHandler struct {
	use DeliveryUsecase
}

func NewDeliveryHandler(use DeliveryUsecase) *DeliveryHandler {
	return &DeliveryHandler{
		use: use,
	}
}

func (h *DeliveryHandler) Assign (w http.ResponseWriter, r *http.Request) {
	var req model.OrderIDRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	deli, err := h.use.Assign(r.Context(), req.OrderID)
	if errors.Is(err, model.ErrCourierNotAvailable) {
    http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deli)

}

func (h *DeliveryHandler) Unassign (w http.ResponseWriter, r *http.Request){
	var req model.OrderIDRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	deli, err := h.use.Unassign(r.Context(), req.OrderID)
	if errors.Is(err, model.ErrCourierNotOrder) {
    http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deli)

}