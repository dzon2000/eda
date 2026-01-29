package web

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/dzon2000/eda/order/internal/db"
	"github.com/dzon2000/eda/order/internal/payload"
)

type Handler struct {
	orderRepository db.OrderRepository
	db              *sql.DB
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", h.CreateOrder)
	mux.HandleFunc("GET /health", h.Health)
	return mux
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req payload.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondJSON(w, http.StatusCreated, payload.CreateOrderResponse{
		OrderID: req.OrderID,
		Status:  "CREATED",
	})
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func respondJSON(w http.ResponseWriter, httpStatus int, createOrderResponse payload.CreateOrderResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(createOrderResponse)
}
