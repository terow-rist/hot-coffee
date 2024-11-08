package service

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"hot-coffee/models"
)

type OrderService struct {
	// You can add dependencies here if needed
}

func (s *OrderService) CreateOrder(order *models.Order) error {
	log.Printf("Creating new order with ID %s", order.ID)
	// Load existing orders
	orders, err := loadOrders()
	if err != nil {
		log.Printf("Error loading orders: %v", err)
		return err
	}

	// Append the new order
	orders = append(orders, *order)

	// Save the updated orders back to the JSON file
	err = saveOrders(orders)
	if err != nil {
		log.Printf("Error saving orders after creating order with ID %s: %v", order.ID, err)
	}
	return err
}

// Generate a new unique order ID based on the last stored ID in the file.
func GenerateOrderID() (string, error) {
	log.Println("Generating a new unique order ID")
	orders, err := loadOrders()
	if err != nil {
		log.Printf("Error loading orders: %v", err)
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

	newID := fmt.Sprintf("order%d", maxID+1)
	log.Printf("Generated new order ID: %s", newID)
	return newID, nil
}

func loadOrders() ([]models.Order, error) {
	log.Println("Loading orders from file")
	var orders []models.Order
	file, err := os.Open("data/orders.json")
	if err != nil {
		log.Printf("Error opening orders file: %v", err)
		return nil, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&orders)
	if err != nil {
		log.Printf("Error decoding orders JSON: %v", err)
	}
	return orders, err
}

func saveOrder(order models.Order) error {
	log.Printf("Saving order with ID %s", order.ID)
	// Load existing orders
	orders, err := loadOrders()
	if err != nil {
		log.Printf("Error loading orders: %v", err)
		return err
	}

	// Check if the order already exists to avoid duplication
	for _, existingOrder := range orders {
		if existingOrder.ID == order.ID {
			errMsg := fmt.Sprintf("Order with ID %s already exists", order.ID)
			log.Println(errMsg)
			return fmt.Errorf(errMsg)
		}
	}

	// Append the new order
	orders = append(orders, order)

	// Save the updated list back to the file
	err = saveOrders(orders)
	if err != nil {
		log.Printf("Error saving orders after adding order with ID %s: %v", order.ID, err)
	}
	return err
}

func saveOrders(orders []models.Order) error {
	log.Println("Saving updated list of orders to file")
	file, err := os.OpenFile("data/orders.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		log.Printf("Error opening orders file for writing: %v", err)
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(orders)
	if err != nil {
		log.Printf("Error encoding orders to JSON: %v", err)
	}
	return err
}
