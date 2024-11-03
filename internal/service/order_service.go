package service

import (
	"encoding/json"
	"fmt"
	"os"

	"hot-coffee/models" // Adjust the import path
)

type OrderService struct {
	// You can add dependencies here if needed
}

func (s *OrderService) CreateOrder(order *models.Order) error {
	// Load existing orders
	orders, err := loadOrders()
	if err != nil {
		return err
	}

	// Append the new order
	orders = append(orders, *order)

	// Save the updated orders back to the JSON file
	return saveOrders(orders)
}

// Generate a new unique order ID based on the last stored ID in the file.
func GenerateOrderID() (string, error) {
	orders, err := loadOrders()
	if err != nil {
		return "", err
	}

	// Check the highest order number used
	maxID := 0
	for _, order := range orders {
		var idNum int
		if _, err := fmt.Sscanf(order.ID, "order%d", &idNum); err == nil && idNum > maxID {
			maxID = idNum
		}
	}

	return fmt.Sprintf("order%d", maxID+1), nil
}

func loadOrders() ([]models.Order, error) {
	var orders []models.Order
	file, err := os.Open("data/orders.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&orders)
	return orders, err
}

func saveOrder(order models.Order) error {
	// Load existing orders
	orders, err := loadOrders()
	if err != nil {
		return err
	}

	// Check if the order already exists to avoid duplication
	for _, existingOrder := range orders {
		if existingOrder.ID == order.ID {
			return fmt.Errorf("order with ID %s already exists", order.ID)
		}
	}

	// Append the new order
	orders = append(orders, order)

	// Save the updated list back to the file
	return saveOrders(orders)
}

func saveOrders(orders []models.Order) error {
	file, err := os.OpenFile("data/orders.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(orders)
}
