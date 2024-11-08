package dal

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"hot-coffee/models"
)

type MenuRepository interface {
	AddItem(item *models.MenuItem) error
	GetAllItems() ([]models.MenuItem, error)
	SaveItems(items []models.MenuItem) error
}

type FileMenuRepository struct{}

func (r *FileMenuRepository) AddItem(item *models.MenuItem) error {
	log.Printf("Adding item with ID %s", item.ID)
	items, err := r.GetAllItems()
	if err != nil {
		log.Printf("Error retrieving items: %v", err)
		return err
	}

	for _, existingItem := range items {
		if existingItem.ID == item.ID {
			log.Printf("Item with ID %s already exists", item.ID)
			return fmt.Errorf("item with this ID %s already exists", item.ID)
		}
	}
	items = append(items, *item)
	err = r.saveItems(items)
	if err != nil {
		log.Printf("Error saving items: %v", err)
	} else {
		log.Printf("Item with ID %s successfully added", item.ID)
	}
	return err
}

func (r *FileMenuRepository) GetAllItems() ([]models.MenuItem, error) {
	log.Println("Retrieving all menu items")
	var items []models.MenuItem
	file, err := os.Open("data/menu_items.json")
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("No menu items file found, returning empty list")
			return items, nil
		}
		log.Printf("Error opening file: %v", err)
		return nil, err
	}
	defer file.Close()

	if err = json.NewDecoder(file).Decode(&items); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		return nil, err
	}
	log.Printf("Retrieved %d items", len(items))
	return items, nil
}

func (r *FileMenuRepository) saveItems(items []models.MenuItem) error {
	log.Println("Saving menu items")
	file, err := os.OpenFile("data/menu_items.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		log.Printf("Error opening file for saving: %v", err)
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(items)
	if err != nil {
		log.Printf("Error encoding JSON: %v", err)
	} else {
		log.Println("Menu items successfully saved")
	}
	return err
}

func (r *FileMenuRepository) SaveItems(items []models.MenuItem) error {
	return r.saveItems(items)
}
