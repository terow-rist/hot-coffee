package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"hot-coffee/internal/service" // Adjust the import path
	"hot-coffee/models"           // Adjust the import path
)

type OrderHandler struct {
	orderService *service.OrderService // Change this to a pointer
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler { // Change argument type
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order

	// Parse the request body
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate the order
	if order.CustomerName == "" || len(order.Items) == 0 {
		http.Error(w, "Customer name and items are required", http.StatusBadRequest)
		return
	}

	// Generate a unique order ID
	order.ID, err = service.GenerateOrderID() // Ensure this has access to the correct context
	if err != nil {
		http.Error(w, "Failed to generate order ID", http.StatusInternalServerError)
		return
	}

	// Set additional fields
	order.Status = "open"
	order.CreatedAt = time.Now().Format(time.RFC3339)

	// Call the service to save the order
	err = h.orderService.CreateOrder(&order)
	if err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	// Respond with the created order
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}
