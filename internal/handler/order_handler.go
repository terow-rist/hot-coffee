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
		if strings.HasSuffix(path, "/close") {
			h.CloseOrder(w, r)
		} else {
			h.CreateOrder(w, r)
		}
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
		slog.Error("Failed to decode order", slog.String("error", err.Error()))
		respondWithError(w, "Invalid input", http.StatusBadRequest)
		return
	}

	order.Status = "open"
	order.CreatedAt = time.Now().Format(time.RFC3339)

	if err := h.orderService.CreateOrder(&order); err != nil {
		slog.Error("Failed to create order", slog.String("error", err.Error()))
		respondWithError(w, err.Error(), http.StatusNotFound)
		return
	}

	slog.Info("Order created successfully", slog.String("orderID", order.ID))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) CloseOrder(w http.ResponseWriter, r *http.Request) {
	orderID := strings.TrimPrefix(strings.TrimSuffix(r.URL.Path, "/close"), "/orders/")
	slog.Info("Closing order", slog.String("orderID", orderID))

	if err := h.orderService.CloseOrder(orderID); err != nil {
		slog.Error("Failed to close order", slog.String("orderID", orderID), slog.String("error", err.Error()))
		respondWithError(w, err.Error(), http.StatusNotFound)
		return
	}

	slog.Info("Order closed successfully", slog.String("orderID", orderID))
	respondWithJSON(w, map[string]string{"message": "Order closed successfully"}, http.StatusOK)
}

func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	slog.Info("Fetching all orders")
	orders, err := h.orderService.GetAllOrders()
	if err != nil {
		slog.Error("Failed to retrieve orders", slog.String("error", err.Error()))
		http.Error(w, "Failed to retrieve orders", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/orders/")
	slog.Info("Fetching order by ID", slog.String("orderID", id))
	order, err := h.orderService.GetOrderByID(id)
	if err != nil {
		slog.Error("Order not found", slog.String("orderID", id))
		respondWithError(w, "Order not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	orderID := strings.TrimPrefix(r.URL.Path, "/orders/")
	slog.Info("Deleting order", slog.String("orderID", orderID))

	if err := h.orderService.DeleteOrder(orderID); err != nil {
		slog.Error("Failed to delete order", slog.String("orderID", orderID), slog.String("error", err.Error()))
		if err.Error() == "order not found" {
			respondWithError(w, "Order not found", http.StatusNotFound)
		} else {
			respondWithError(w, "Failed to delete order", http.StatusInternalServerError)
		}
		return
	}

	slog.Info("Order deleted successfully", slog.String("orderID", orderID))
	respondWithJSON(w, "Order deleted successfully", http.StatusNoContent)
}

func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	orderID := strings.TrimPrefix(r.URL.Path, "/orders/")
	var updatedOrder models.Order
	slog.Info("Updating order", slog.String("orderID", orderID))

	if err := json.NewDecoder(r.Body).Decode(&updatedOrder); err != nil {
		slog.Error("Failed to decode update data", slog.String("error", err.Error()))
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	existingOrder, err := h.orderService.GetOrderByID(orderID)
	if err != nil {
		slog.Error("Order not found", slog.String("orderID", orderID))
		respondWithError(w, "Order not found", http.StatusNotFound)
		return
	}

	if updatedOrder.CustomerName != "" {
		existingOrder.CustomerName = updatedOrder.CustomerName
	}
	if len(updatedOrder.Items) > 0 {
		existingOrder.Items = updatedOrder.Items
	}

	existingOrder.CreatedAt = time.Now().Format(time.RFC3339)

	if err := h.orderService.UpdateOrder(existingOrder); err != nil {
		slog.Error("Failed to update order", slog.String("orderID", orderID), slog.String("error", err.Error()))
		respondWithError(w, "Failed to update order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Order updated successfully", slog.String("orderID", orderID))
	respondWithJSON(w, existingOrder, http.StatusOK)
}

func respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
