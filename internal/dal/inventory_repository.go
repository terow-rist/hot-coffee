package dal

import (
	"encoding/json"
	"fmt"
	"os"

	"hot-coffee/models"
)

type InventoryRepository interface {
	AddItem(item *models.InventoryItem) error
	GetAllItems() ([]models.InventoryItem, error)
}

type FileInventoryRepository struct{}

// AddItem adds a new inventory item to the repository.
// It reads the existing items, checks for duplicates, and then saves the new item.
func (r *FileInventoryRepository) AddItem(item *models.InventoryItem) error {
	items, err := r.GetAllItems()
	if err != nil {
		return err
	}

	// Check for duplicate ingredient ID
	for _, existingItem := range items {
		if existingItem.IngredientID == item.IngredientID {
			return fmt.Errorf("item with Ingredient ID %s already exists", item.IngredientID)
		}
	}

	// Append the new item and save
	items = append(items, *item)
	return r.saveItems(items)
}

// GetAllItems retrieves all inventory items from the repository.
func (r *FileInventoryRepository) GetAllItems() ([]models.InventoryItem, error) {
	var items []models.InventoryItem
	file, err := os.Open("data/inventory.json")
	if err != nil {
		// If the file doesn't exist, return an empty slice without an error
		if os.IsNotExist(err) {
			return items, nil
		}
		return nil, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&items)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// saveItems saves the inventory items to the JSON file.
func (r *FileInventoryRepository) saveItems(items []models.InventoryItem) error {
	file, err := os.OpenFile("data/inventory.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the inventory items back to the file as JSON
	return json.NewEncoder(file).Encode(items)
}
