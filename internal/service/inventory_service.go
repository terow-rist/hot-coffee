package service

import (
	"hot-coffee/models"
)

type InventoryRepository interface {
	AddItem(item *models.InventoryItem) error
	GetAllItems() ([]models.InventoryItem, error)
	SaveItems(items []models.InventoryItem) error // Add this line

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

// service/inventory_service.go

func (s *InventoryService) GetItemByID(id string) (*models.InventoryItem, error) {
	items, err := s.repo.GetAllItems()
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if item.IngredientID == id {
			return &item, nil
		}
	}

	return nil, models.ErrItemNotFound // You need to define this error
}

func (s *InventoryService) UpdateItem(item *models.InventoryItem) error {
	items, err := s.repo.GetAllItems()
	if err != nil {
		return err
	}

	for i, existingItem := range items {
		if existingItem.IngredientID == item.IngredientID {
			// Update the existing item with fields from the updated item
			items[i].Name = item.Name         // Assuming "Name" is a field in InventoryItem
			items[i].Quantity = item.Quantity // Assuming "Quantity" is a field in InventoryItem
			// Update other fields as needed
			return s.repo.SaveItems(items) // Save updated items back to the repository
		}
	}

	return models.ErrItemNotFound // Return error if the item is not found
}

// Add other service methods as needed
