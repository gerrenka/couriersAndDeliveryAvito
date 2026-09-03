package handlers

import (
	"avito/internal/model"
	"avito/internal/service"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CourierService interface {
	GetById(ctx context.Context, id int) (model.Courier, error)
	GetAllCourier(ctx context.Context) ([]model.Courier, error)
	CreateCourier(ctx context.Context, name, phone, status, transport_type string) (model.Courier, error)
	UpdateCourier(ctx context.Context, id int, name, phone, status, transport_type string) (model.Courier, error)
	DeleteCourier(ctx context.Context, id int) error
}

type CourierHandler struct {
	serv CourierService
}

func NewCourierHandler(serv CourierService) *CourierHandler {
	return &CourierHandler{
		serv: serv,
	}
}

func (h *CourierHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "pong"}
	json.NewEncoder(w).Encode(response)
}

func (h *CourierHandler) Healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (h *CourierHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateCourierRequest /// курьер с таким телефоном сущ
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c, err := h.serv.CreateCourier(r.Context(), req.Name, req.Phone, req.Status, req.TransportType)
	if errors.Is(err, service.ErrInvalidStatus) ||
		errors.Is(err, service.ErrInvalidName) ||
		errors.Is(err, service.ErrInvalidPhone) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//
	if errors.Is(err, model.ErrPhoneExists) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	//
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func (h *CourierHandler) List(w http.ResponseWriter, r *http.Request) {
	couriers, err := h.serv.GetAllCourier(r.Context())
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError) // нужен ли?
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(couriers)
}

func (h *CourierHandler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	c, err := h.serv.GetById(r.Context(), id)

	if errors.Is(err, model.ErrCourierNotFound) {
		http.Error(w, "Courier not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError) //под вопросом нужен ли
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(c)
}

func (h *CourierHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var req model.CreateCourierRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c, err := h.serv.UpdateCourier(r.Context(), id, req.Name, req.Phone, req.Status, req.TransportType)
	if errors.Is(err, service.ErrInvalidStatus) ||
		errors.Is(err, service.ErrInvalidName) ||
		errors.Is(err, service.ErrInvalidPhone) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if errors.Is(err, model.ErrCourierNotFound) {
		http.Error(w, "Courier not found", http.StatusNotFound)
		return
	}
	//
	if errors.Is(err, model.ErrPhoneExists) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	//
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError) // возможно не нужен
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(c)
}

func (h *CourierHandler) Delete(w http.ResponseWriter, r *http.Request) { // в задании не написано писать
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	err = h.serv.DeleteCourier(r.Context(), id)
	if errors.Is(err, model.ErrCourierNotFound) {
		http.Error(w, "Courier not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
