package handler

import (
	"encoding/json"
	"net/http"
	"strings"

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
	path := r.URL.Path
	switch r.Method {
	case http.MethodPost:
		if strings.HasPrefix(path, "/inventory/") {
			http.Error(w, "Invalid request", http.StatusBadRequest)
		} else {
			h.AddInventoryItem(w, r)
		}

	case http.MethodGet:
		if strings.HasPrefix(path, "/inventory/") {
			// Extract the ID from the URL
			id := strings.TrimPrefix(path, "/inventory/")
			h.GetInventoryItem(w, r, id) // Pass the ID to the handler
		} else {
			h.GetAllInventoryItems(w, r)
		}
	case http.MethodPut:
		if strings.HasPrefix(path, "/inventory/") {
			// Extract the ID from the URL
			id := strings.TrimPrefix(path, "/inventory/")
			h.UpdateInventoryItem(w, r, id) // Pass the ID to the handler
		} else {
			http.Error(w, "Invalid request", http.StatusBadRequest)
		}
	case http.MethodDelete:
		if strings.HasPrefix(path, "/inventory/") {
			// Extract the ID from the URL
			id := strings.TrimPrefix(path, "/inventory/")
			h.DeleteInventoryItem(w, r, id) // Pass the ID to the handler
		} else {
			http.Error(w, "Invalid request", http.StatusBadRequest)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *InventoryHandler) DeleteInventoryItem(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.service.DeleteItem(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent) // No content response for successful deletion
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

func (h *InventoryHandler) GetInventoryItem(w http.ResponseWriter, r *http.Request, id string) {
	item, err := h.service.GetItemByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}

func (h *InventoryHandler) UpdateInventoryItem(w http.ResponseWriter, r *http.Request, id string) {
	var updatedItem models.InventoryItem
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Ensure the ID is set correctly
	updatedItem.IngredientID = id // Use the ID from the URL

	// Call the existing UpdateItem method in the service
	if err := h.service.UpdateItem(&updatedItem); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedItem)
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
