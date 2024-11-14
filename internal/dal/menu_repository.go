package dal

import (
	"encoding/json"
	"fmt"
	"hot-coffee/config"
	"hot-coffee/models"
	"os"
)

type MenuRepository interface {
	AddItem(item *models.MenuItem) error
	GetAllItems() ([]models.MenuItem, error)
	SaveItems(items []models.MenuItem) error
}

type FileMenuRepository struct{}

func (r *FileMenuRepository) AddItem(item *models.MenuItem) error {
	items, err := r.GetAllItems()
	if err != nil {
		return err
	}

	for _, existingItem := range items {
		if existingItem.ID == item.ID {
			return fmt.Errorf("item with this ID %s already exists", item.ID)
		}
	}
	items = append(items, *item)
	return r.saveItems(items)
}

func (r *FileMenuRepository) GetAllItems() ([]models.MenuItem, error) {
	var items []models.MenuItem
	file, err := os.Open(config.Directory + "/menu_items.json")
	if err != nil {
		if os.IsNotExist(err) {
			return items, nil
		}
		return nil, err
	}
	defer file.Close()

	if err = json.NewDecoder(file).Decode(&items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *FileMenuRepository) saveItems(items []models.MenuItem) error {
	file, err := os.OpenFile(config.Directory+"/menu_items.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(items)
}

func (r *FileMenuRepository) SaveItems(items []models.MenuItem) error {
	return r.saveItems(items)
}
