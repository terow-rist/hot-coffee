package dal

import (
	"encoding/json"
	"errors"
	"fmt"
	"hot-coffee/models"
	"os"
)

type OrderRepository interface {
	SaveOrder(order *models.Order) error
	GetAllOrders() ([]models.Order, error)
	GetOrderByID(id string) (*models.Order, error)
	UpdateOrder(order *models.Order) error
}

type FileOrderRepository struct{}

func (repo *FileOrderRepository) SaveOrder(order *models.Order) error {
	var orders []models.Order

	file, err := os.OpenFile("data/orders.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if err = json.NewDecoder(file).Decode(&orders); err != nil && err != os.ErrNotExist {
		return err
	}

	orders = append(orders, *order)
	file.Seek(0, 0)
	return json.NewEncoder(file).Encode(orders)
}

func (r *FileOrderRepository) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order
	file, err := os.Open("data/orders.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *FileOrderRepository) GetOrderByID(id string) (*models.Order, error) {
	orders, err := r.GetAllOrders()
	if err != nil {
		return nil, err
	}

	for _, order := range orders {
		if order.ID == id {
			return &order, nil
		}
	}
	return nil, errors.New("order not found")
}

func (r *FileOrderRepository) UpdateOrder(order *models.Order) error {
	// Load the existing orders from the file
	orders, err := r.LoadOrders()
	if err != nil {
		return err
	}

	// Find the order to update by ID
	var updated bool
	for i, o := range orders {
		if o.ID == order.ID {
			// Update only the specified fields
			if order.CustomerName != "" {
				orders[i].CustomerName = order.CustomerName
			}
			if len(order.Items) > 0 {
				orders[i].Items = order.Items
			}
			if order.CreatedAt != "" { // Check if CreatedAt is set (not empty string)
				orders[i].CreatedAt = order.CreatedAt
			}
			updated = true
			break
		}
	}

	// If no order found with the given ID, return an error
	if !updated {
		return fmt.Errorf("order not found")
	}

	// Save the updated list of orders back to the file, only modified orders
	return r.SaveOrders(orders)
}

func (r *FileOrderRepository) LoadOrders() ([]models.Order, error) {
	var orders []models.Order

	file, err := os.Open("data/orders.json")
	if err != nil {
		// If the file does not exist, return an empty slice (this is valid)
		if os.IsNotExist(err) {
			return orders, nil
		}
		return nil, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&orders)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *FileOrderRepository) SaveOrders(orders []models.Order) error {
	file, err := os.OpenFile("data/orders.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the updated orders to the file
	return json.NewEncoder(file).Encode(orders)
}

func (r *FileOrderRepository) DeleteOrder(orderID string) error {
	orders, err := r.GetAllOrders()
	if err != nil {
		return err
	}

	// Filter out the order to delete
	var updatedOrders []models.Order
	for _, order := range orders {
		if order.ID != orderID {
			updatedOrders = append(updatedOrders, order)
		}
	}

	// Check if order was found
	if len(orders) == len(updatedOrders) {
		return errors.New("order not found")
	}

	// Write the updated orders back to file
	file, err := os.OpenFile("data/orders.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(updatedOrders)
}
