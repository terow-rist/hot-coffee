package dal

import (
	"encoding/json"
	"errors"
	"hot-coffee/models"
	"os"
)

type OrderRepository interface {
	SaveOrder(order *models.Order) error
	GetAllOrders() ([]models.Order, error)
	GetOrderByID(id string) (*models.Order, error)
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
