package handler

import (
	"encoding/json"
	"log/slog"
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
	slog.Info("Received request", slog.String("method", r.Method), slog.String("path", path))

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
		slog.Error("Error decoding request body", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.AddItem(&item); err != nil {
		slog.Error("Error adding menu item", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Menu item successfully added", slog.String("itemID", item.ID))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (h *MenuHandler) GetAllMenuItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetAllItems()
	if err != nil {
		slog.Error("Error retrieving menu items", slog.Any("error", err))
		http.Error(w, "Failed to retrieve menu items", http.StatusInternalServerError)
		return
	}

	slog.Info("Retrieved all menu items", slog.Int("itemCount", len(items)))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}

func (h *MenuHandler) GetMenuItem(w http.ResponseWriter, r *http.Request, id string) {
	item, err := h.service.GetMenuItemByID(id)
	if err != nil {
		slog.Error("Error retrieving menu item", slog.String("itemID", id), slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	slog.Info("Menu item retrieved", slog.String("itemID", id))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}

func (h *MenuHandler) UpdateMenuItem(w http.ResponseWriter, r *http.Request, id string) {
	var updatedItem models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		slog.Error("Error decoding request body for update", slog.Any("error", err))
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	updatedItem.ID = id

	if err := h.service.UpdateMenuItem(&updatedItem); err != nil {
		slog.Error("Error updating menu item", slog.String("itemID", id), slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Menu item updated successfully", slog.String("itemID", id))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedItem)
}

func (h *MenuHandler) DeleteMenuItem(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.service.DeleteMenuItem(id); err != nil {
		slog.Error("Error deleting menu item", slog.String("itemID", id), slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	slog.Info("Menu item deleted", slog.String("itemID", id))
	w.WriteHeader(http.StatusNoContent)
}
