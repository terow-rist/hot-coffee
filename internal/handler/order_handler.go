package handler

import (
	"encoding/json"
	"hot-coffee/internal/service"
	"hot-coffee/models"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	path := r.URL.Path
	slog.Info("Received request", slog.String("method", r.Method), slog.String("path", path))

	switch r.Method {
	case http.MethodPost:
		h.CreateOrder(w, r)
	case http.MethodGet:
		// Handle request for all orders
		if path == "/orders" {
			h.GetAllOrders(w, r)
		} else if strings.HasPrefix(path, "/orders/") {
			// Handle request for a specific order by ID
			h.GetOrderByID(w, r)
		} else {
			respondWithError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		respondWithError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		respondWithError(w, "Invalid input", http.StatusBadRequest)
		return
	}

	order.Status = "open"
	order.CreatedAt = time.Now().Format(time.RFC3339)

	if err := h.orderService.CreateOrder(&order); err != nil {
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// / Handle GET /orders
func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.orderService.GetAllOrders()
	if err != nil {
		http.Error(w, "Failed to retrieve orders", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/orders/") // Extract the ID from the URL path
	order, err := h.orderService.GetOrderByID(id)
	if err != nil {
		respondWithError(w, "Order not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
