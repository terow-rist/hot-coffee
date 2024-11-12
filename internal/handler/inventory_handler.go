package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"hot-coffee/internal/service"
	"hot-coffee/models"
	"log/slog"
)

type InventoryHandler struct {
	service *service.InventoryService // Use a pointer to InventoryService
}

func NewInventoryHandler(service *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: service}
}

// respondWithError formats an error message in JSON format and writes it to the response.
func respondWithError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (h *InventoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	path := r.URL.Path
	slog.Info("Request received", "method", r.Method, "path", path)
	switch r.Method {
	case http.MethodPost:
		if strings.HasPrefix(path, "/inventory/") {
			slog.Error("Invalid POST request - unexpected URL path")
			respondWithError(w, "Invalid request", http.StatusBadRequest)
		} else {
			h.AddInventoryItem(w, r)
		}

	case http.MethodGet:
		if strings.HasPrefix(path, "/inventory/") {
			id := strings.TrimPrefix(path, "/inventory/")
			h.GetInventoryItem(w, r, id)
		} else {
			h.GetAllInventoryItems(w, r)
		}

	case http.MethodPut:
		if strings.HasPrefix(path, "/inventory/") {
			id := strings.TrimPrefix(path, "/inventory/")
			h.UpdateInventoryItem(w, r, id)
		} else {
			slog.Error("Invalid PUT request - unexpected URL path")
			respondWithError(w, "Invalid request", http.StatusBadRequest)
		}

	case http.MethodDelete:
		if strings.HasPrefix(path, "/inventory/") {
			id := strings.TrimPrefix(path, "/inventory/")
			h.DeleteInventoryItem(w, r, id)
		} else {
			slog.Error("Invalid DELETE request - unexpected URL path")
			respondWithError(w, "Invalid request", http.StatusBadRequest)
		}

	default:
		slog.Warn("Method not allowed", "method", r.Method)
		respondWithError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *InventoryHandler) DeleteInventoryItem(w http.ResponseWriter, r *http.Request, id string) {
	slog.Info("Deleting inventory item", "id", id)
	if err := h.service.DeleteItem(id); err != nil {
		slog.Error("Error deleting inventory item", "id", id, "error", err)
		respondWithError(w, err.Error(), http.StatusNotFound)
		return
	}
	slog.Info("Inventory item deleted", "id", id)
	w.WriteHeader(http.StatusNoContent)
}

func (h *InventoryHandler) AddInventoryItem(w http.ResponseWriter, r *http.Request) {
	var item models.InventoryItem

	slog.Info("Decoding JSON for new inventory item")
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		slog.Error("Error decoding JSON", "error", err)
		respondWithError(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.service.AddItem(&item); err != nil {
		slog.Error("Error adding inventory item", "error", err)
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Inventory item added", "id", item.IngredientID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (h *InventoryHandler) GetInventoryItem(w http.ResponseWriter, r *http.Request, id string) {
	slog.Info("Retrieving inventory item", "id", id)
	item, err := h.service.GetItemByID(id)
	if err != nil {
		slog.Error("Error retrieving inventory item", "id", id, "error", err)
		respondWithError(w, err.Error(), http.StatusNotFound)
		return
	}

	slog.Info("Inventory item retrieved", "id", id)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}

func (h *InventoryHandler) UpdateInventoryItem(w http.ResponseWriter, r *http.Request, id string) {
	var updatedItem models.InventoryItem

	slog.Info("Updating inventory item", "id", id)
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		slog.Error("Error decoding JSON", "error", err)
		respondWithError(w, "Invalid input", http.StatusBadRequest)
		return
	}

	updatedItem.IngredientID = id
	if err := h.service.UpdateItem(&updatedItem); err != nil {
		slog.Error("Error updating inventory item", "id", id, "error", err)
		respondWithError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Inventory item updated", "id", id)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedItem)
}

func (h *InventoryHandler) GetAllInventoryItems(w http.ResponseWriter, r *http.Request) {
	slog.Info("Retrieving all inventory items")
	items, err := h.service.GetAllItems()
	if err != nil {
		slog.Error("Error retrieving all items", "error", err)
		respondWithError(w, "Failed to retrieve inventory items", http.StatusInternalServerError)
		return
	}

	slog.Info("Inventory items retrieved", "count", len(items))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}
