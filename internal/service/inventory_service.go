package service

import (
	"hot-coffee/models"
)

type InventoryRepository interface {
	AddItem(item *models.InventoryItem) error
	GetAllItems() ([]models.InventoryItem, error)
	// Add other methods as needed
}

type InventoryService struct {
	repo InventoryRepository // Ensure this field exists
}

func NewInventoryService(repo InventoryRepository) *InventoryService {
	return &InventoryService{repo: repo}
}

func (s *InventoryService) AddItem(item *models.InventoryItem) error {
	return s.repo.AddItem(item)
}

// service/inventory_service.go
func (s *InventoryService) GetAllItems() ([]models.InventoryItem, error) {
	return s.repo.GetAllItems()
}

// Add other service methods as needed
