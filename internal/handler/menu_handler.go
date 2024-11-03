package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"hot-coffee/internal/service"
	"hot-coffee/models"
)

type MenuHandler struct {
	service *service.MenuService
}

func NewMenuHandler(service *service.MenuService) *MenuHandler {
	return &MenuHandler{service: service}
}

func (h *MenuHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	path := r.URL.Path
	switch r.Method {
	case http.MethodPost:
		if strings.HasPrefix(path, "/menu/") {
			http.Error(w, "Invalid request", http.StatusBadRequest)
		} else {
			h.AddMenuItem(w, r)
		}
	case http.MethodGet:
		if strings.HasPrefix(path, "/menu/") {
			id := strings.TrimPrefix(path, "/menu/")
			h.GetMenuItem(w, r, id)
		} else {
			h.GetAllMenuItems(w, r)
		}
	case http.MethodPut:
		if strings.HasPrefix(path, "/menu/") {
			id := strings.TrimPrefix(path, "/menu/")
			h.UpdateMenuItem(w, r, id)
		} else {
			http.Error(w, "Invalid request", http.StatusBadRequest)
		}
	case http.MethodDelete:
		if strings.HasPrefix(path, "/menu/") {
			id := strings.TrimPrefix(path, "/menu/")
			h.DeleteMenuItem(w, r, id)
		} else {
			http.Error(w, "Invalid request", http.StatusBadRequest)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *MenuHandler) AddMenuItem(w http.ResponseWriter, r *http.Request) {
	var item models.MenuItem

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.AddItem(&item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (h *MenuHandler) GetAllMenuItems(w http.ResponseWriter, r *http.Request) {
	// Call the service to get all items
	items, err := h.service.GetAllItems() // Using the service method here
	if err != nil {
		http.Error(w, "Failed to retrieve menu items", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)     // Set the response status code
	json.NewEncoder(w).Encode(items) // Encode and return the items in JSON format
}

func (h *MenuHandler) GetMenuItem(w http.ResponseWriter, r *http.Request, id string) {
	item, err := h.service.GetMenuItemByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}

func (h *MenuHandler) UpdateMenuItem(w http.ResponseWriter, r *http.Request, id string) {
	var updatedItem models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Ensure the ID is set correctly
	updatedItem.ID = id // Use the ID from the URL

	// Call the existing UpdateItem method in the service
	if err := h.service.UpdateMenuItem(&updatedItem); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedItem)
}

func (h *MenuHandler) DeleteMenuItem(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.service.DeleteMenuItem(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent) // No content response for successful deletion
}
