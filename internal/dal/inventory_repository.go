package dal

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"hot-coffee/models"
)

type InventoryRepository interface {
	AddItem(item *models.InventoryItem) error
	GetAllItems() ([]models.InventoryItem, error)
	SaveItems(items []models.InventoryItem) error
}

type FileInventoryRepository struct{}

// AddItem adds a new inventory item to the repository.
// It reads the existing items, checks for duplicates, and then saves the new item.
func (r *FileInventoryRepository) AddItem(item *models.InventoryItem) error {
	log.Printf("Attempting to add item: %+v", item)

	items, err := r.GetAllItems()
	if err != nil {
		log.Printf("Error getting all items: %v", err)
		return err
	}

	// Check for duplicate ingredient ID
	for _, existingItem := range items {
		if existingItem.IngredientID == item.IngredientID {
			err := fmt.Errorf("item with Ingredient ID %s already exists", item.IngredientID)
			log.Printf("Duplicate item found: %v", err)
		}
	}

	// Append the new item and save
	items = append(items, *item)
	err = r.saveItems(items)
	if err != nil {
		log.Printf("Error saving items: %v", err)
	}
	return err
}

// GetAllItems retrieves all inventory items from the repository.
func (r *FileInventoryRepository) GetAllItems() ([]models.InventoryItem, error) {
	log.Println("Retrieving all inventory items")

	var items []models.InventoryItem
	file, err := os.Open("data/inventory.json")
	if err != nil {
		if os.IsNotExist(err) {
			// If the file doesn't exist, return an empty slice without an error

			log.Println("Inventory file does not exist, returning empty slice")
			return items, nil
		}
		log.Printf("Error opening inventory file: %v", err)
		return nil, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&items)
	if err != nil {
		log.Printf("Error decoding inventory items: %v", err)
		return nil, err
	}

	log.Printf("Retrieved %d inventory items", len(items))
	return items, nil
}

// saveItems saves the inventory items to the JSON file.
func (r *FileInventoryRepository) saveItems(items []models.InventoryItem) error {
	log.Printf("Saving %d inventory items", len(items))

	file, err := os.OpenFile("data/inventory.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		log.Printf("Error opening inventory file for writing: %v", err)
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(items)
	if err != nil {
		log.Printf("Error encoding inventory items to JSON: %v", err)
	}
	return err
}

// dal/file_inventory_repository.go

func (r *FileInventoryRepository) SaveItems(items []models.InventoryItem) error {
	return r.saveItems(items)
}
