package dal

import (
	"encoding/json"
	"hot-coffee/models"
	"os"
)

type OrderRepository interface {
	SaveOrder(order *models.Order) error
}

type FileOrderRepository struct {
}

func NewOrderRepository() OrderRepository {
	return &FileOrderRepository{}
}

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
