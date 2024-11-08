package handler

import (
	"encoding/json"
	"log"
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
	w.Header().Set("Content-Type", "application/json")
	path := r.URL.Path
	log.Printf("Received %s request at %s", r.Method, path)
	switch r.Method {
	case http.MethodPost:
		if strings.HasPrefix(path, "/inventory/") {
			log.Println("Invalid POST request - unexpected URL path")
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
			log.Println("Invalid PUT request - unexpected URL path")
			http.Error(w, "Invalid request", http.StatusBadRequest)
		}

	case http.MethodDelete:
		if strings.HasPrefix(path, "/inventory/") {
			// Extract the ID from the URL
			id := strings.TrimPrefix(path, "/inventory/")
			h.DeleteInventoryItem(w, r, id) // Pass the ID to the handler
		} else {
			log.Println("Invalid DELETE request - unexpected URL path")
			http.Error(w, "Invalid request", http.StatusBadRequest)
		}

	default:
		log.Printf("Method %s not allowed", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *InventoryHandler) DeleteInventoryItem(w http.ResponseWriter, r *http.Request, id string) {
	log.Printf("Attempting to delete item with ID %s", id)
	if err := h.service.DeleteItem(id); err != nil {
		log.Printf("Error deleting item with ID %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	log.Printf("Item with ID %s successfully deleted", id)
	w.WriteHeader(http.StatusNoContent) // No content response for successful deletion
}

func (h *InventoryHandler) AddInventoryItem(w http.ResponseWriter, r *http.Request) {
	var item models.InventoryItem

	// Decode the incoming JSON
	log.Println("Decoding JSON for new inventory item")
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Call the service to add the item
	if err := h.service.AddItem(&item); err != nil {
		log.Printf("Error adding inventory item: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError) // Return the actual error message
		return
	}

	log.Printf("Inventory item with ID %s successfully added", item.IngredientID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (h *InventoryHandler) GetInventoryItem(w http.ResponseWriter, r *http.Request, id string) {
	log.Printf("Retrieving item with ID %s", id)
	item, err := h.service.GetItemByID(id)
	if err != nil {
		log.Printf("Error retrieving item with ID %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Printf("Item with ID %s retrieved successfully", id)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}

func (h *InventoryHandler) UpdateInventoryItem(w http.ResponseWriter, r *http.Request, id string) {
	var updatedItem models.InventoryItem

	log.Printf("Updating item with ID %s", id)
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Ensure the ID is set correctly
	updatedItem.IngredientID = id // Use the ID from the URL

	// Call the existing UpdateItem method in the service
	if err := h.service.UpdateItem(&updatedItem); err != nil {
		log.Printf("Error updating item with ID %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Item with ID %s successfully updated", id)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedItem)
}

func (h *InventoryHandler) GetAllInventoryItems(w http.ResponseWriter, r *http.Request) {
	// Call the service to get all items
	log.Println("Retrieving all inventory items")
	items, err := h.service.GetAllItems() // Using the service method here
	if err != nil {
		log.Printf("Error retrieving all items: %v", err)
		http.Error(w, "Failed to retrieve inventory items", http.StatusInternalServerError)
		return
	}

	log.Printf("Retrieved %d inventory items", len(items))
	w.WriteHeader(http.StatusOK)     // Set the response status code
	json.NewEncoder(w).Encode(items) // Encode and return the items in JSON format
}
