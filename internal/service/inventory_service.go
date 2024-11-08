package service

import (
	"log"

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
	log.Printf("Attempting to add item with ID %s", item.IngredientID)
	err := s.repo.AddItem(item)
	if err != nil {
		log.Printf("Error adding item with ID %s: %v", item.IngredientID, err)
	}
	return err
}

func (s *InventoryService) GetAllItems() ([]models.InventoryItem, error) {
	log.Println("Fetching all inventory items")
	items, err := s.repo.GetAllItems()
	if err != nil {
		log.Printf("Error fetching inventory items: %v", err)
	}
	return items, err
}

func (s *InventoryService) GetItemByID(id string) (*models.InventoryItem, error) {
	log.Printf("Fetching item with ID %s", id)
	items, err := s.repo.GetAllItems()
	if err != nil {
		log.Printf("Error fetching inventory items: %v", err)
		return nil, err
	}

	for _, item := range items {
		if item.IngredientID == id {
			log.Printf("Item with ID %s found", id)
			return &item, nil
		}
	}

	log.Printf("Item with ID %s not found", id)
	return nil, models.ErrItemNotFound // You need to define this error
}

func (s *InventoryService) UpdateItem(item *models.InventoryItem) error {
	log.Printf("Attempting to update item with ID %s", item.IngredientID)
	items, err := s.repo.GetAllItems()
	if err != nil {
		log.Printf("Error fetching inventory items: %v", err)
		return err
	}

	for i, existingItem := range items {
		if existingItem.IngredientID == item.IngredientID {
			// Update the existing item with fields from the updated item
			log.Printf("Updating item with ID %s", item.IngredientID)
			items[i].Name = item.Name         // Assuming "Name" is a field in InventoryItem
			items[i].Quantity = item.Quantity // Assuming "Quantity" is a field in InventoryItem
			// Update other fields as needed
			err = s.repo.SaveItems(items) // Save updated items back to the repository
			if err != nil {
				log.Printf("Error saving updated item with ID %s: %v", item.IngredientID, err)
			}
			return err
		}
	}

	log.Printf("Item with ID %s not found for update", item.IngredientID)
	return models.ErrItemNotFound // Return error if the item is not found
}

func (s *InventoryService) DeleteItem(id string) error {
	log.Printf("Attempting to delete item with ID %s", id)
	items, err := s.repo.GetAllItems()
	if err != nil {
		log.Printf("Error fetching inventory items: %v", err)
		return err
	}

	for i, existingItem := range items {
		if existingItem.IngredientID == id {
			// Remove the item from the slice
			log.Printf("Deleting item with ID %s", id)
			items = append(items[:i], items[i+1:]...) // Remove item at index i
			err = s.repo.SaveItems(items)             // Save updated items back to the repository
			if err != nil {
				log.Printf("Error saving after deletion of item with ID %s: %v", id, err)
			}
			return err
		}
	}

	log.Printf("Item with ID %s not found for deletion", id)
	return models.ErrItemNotFound // Return error if the item is not found
}
