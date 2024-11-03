package handler

import (
	"encoding/json"
	"net/http"

	"hot-coffee/internal/service"
	"hot-coffee/models"
)

type InventoryHandler struct {
	service *service.InventoryService // Use a pointer to InventoryService
}

func NewInventoryHandler(service *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: service}
}

func (h *InventoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.AddInventoryItem(w, r)
	case http.MethodGet:
		h.GetAllInventoryItems(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *InventoryHandler) AddInventoryItem(w http.ResponseWriter, r *http.Request) {
	var item models.InventoryItem

	// Decode the incoming JSON
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Call the service to add the item
	if err := h.service.AddItem(&item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // Return the actual error message
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (h *InventoryHandler) GetAllInventoryItems(w http.ResponseWriter, r *http.Request) {
	// Call the service to get all items
	items, err := h.service.GetAllItems() // Using the service method here
	if err != nil {
		http.Error(w, "Failed to retrieve inventory items", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)     // Set the response status code
	json.NewEncoder(w).Encode(items) // Encode and return the items in JSON format
}

// You can add more methods for handling other inventory operations (e.g., GET, PUT, DELETE) here
