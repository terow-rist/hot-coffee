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
		if strings.HasPrefix(path, "/orders/") {
			h.GetOrderByID(w, r)
		} else {
			h.GetAllOrders(w, r)
		}
	case http.MethodPut:
		if strings.HasPrefix(path, "/orders/") {
			h.UpdateOrder(w, r)
		} else {
			respondWithError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case http.MethodDelete:
		if strings.HasPrefix(path, "/orders/") {
			h.DeleteOrder(w, r)
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
		respondWithError(w, err.Error(), http.StatusNotFound)
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

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	orderID := strings.TrimPrefix(r.URL.Path, "/orders/")

	if err := h.orderService.DeleteOrder(orderID); err != nil {
		if err.Error() == "order not found" {
			respondWithError(w, "Order not found", http.StatusNotFound)
		} else {
			respondWithError(w, "Failed to delete order", http.StatusInternalServerError)
		}
		return
	}

	respondWithJSON(w, "Order deleted successfully", http.StatusNoContent)
}

func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	orderID := strings.TrimPrefix(r.URL.Path, "/orders/")
	var updatedOrder models.Order

	// Decode the incoming JSON into the updatedOrder struct
	if err := json.NewDecoder(r.Body).Decode(&updatedOrder); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate that the order exists
	existingOrder, err := h.orderService.GetOrderByID(orderID)
	if err != nil {
		respondWithError(w, "Order not found", http.StatusNotFound)
		return
	}

	// Only update the fields that are provided in the request body (not empty)
	if updatedOrder.CustomerName != "" {
		existingOrder.CustomerName = updatedOrder.CustomerName
	}
	if len(updatedOrder.Items) > 0 {
		existingOrder.Items = updatedOrder.Items
	}

	// Automatically set the CreatedAt to the current time
	existingOrder.CreatedAt = time.Now().Format(time.RFC3339)

	// Call the service to save the updated order
	if err := h.orderService.UpdateOrder(existingOrder); err != nil {
		respondWithError(w, "Failed to update order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the updated order
	respondWithJSON(w, existingOrder, http.StatusOK)
}

func respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
