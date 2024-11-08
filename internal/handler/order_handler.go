package handler

import (
	"encoding/json"
	"log"
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

	log.Println("Received request to create a new order")

	// Parse the request body
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		log.Printf("Error decoding order request body: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate the order
	if order.CustomerName == "" || len(order.Items) == 0 {
		log.Println("Invalid order data: Customer name and items are required")
		http.Error(w, "Customer name and items are required", http.StatusBadRequest)
		return
	}

	// Generate a unique order ID
	order.ID, err = service.GenerateOrderID() // Ensure this has access to the correct context
	if err != nil {
		log.Printf("Error generating order ID: %v", err)
		http.Error(w, "Failed to generate order ID", http.StatusInternalServerError)
		return
	}

	// Set additional fields
	order.Status = "open"
	order.CreatedAt = time.Now().Format(time.RFC3339)
	log.Printf("Order ID %s generated, status set to 'open', and timestamp created", order.ID)

	// Call the service to save the order
	err = h.orderService.CreateOrder(&order)
	if err != nil {
		log.Printf("Error saving order with ID %s: %v", order.ID, err)
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	// Respond with the created order
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(order); err != nil {
		log.Printf("Error encoding response for order ID %s: %v", order.ID, err)
		http.Error(w, "Failed to encode order response", http.StatusInternalServerError)
		return
	}

	log.Printf("Order with ID %s successfully created", order.ID)
}
