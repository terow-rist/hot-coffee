package service

import (
	"log"

	"hot-coffee/models"
)

type MenuRepository interface {
	AddItem(item *models.MenuItem) error
	GetAllItems() ([]models.MenuItem, error)
	SaveItems(items []models.MenuItem) error
}

type MenuService struct {
	repo MenuRepository
}

func NewMenuService(repo MenuRepository) *MenuService {
	return &MenuService{repo: repo}
}

func (s *MenuService) AddItem(item *models.MenuItem) error {
	log.Printf("Attempting to add menu item with ID %s", item.ID)
	err := s.repo.AddItem(item)
	if err != nil {
		log.Printf("Error adding menu item with ID %s: %v", item.ID, err)
	}
	return err
}

func (s *MenuService) GetAllItems() ([]models.MenuItem, error) {
	log.Println("Fetching all menu items")
	items, err := s.repo.GetAllItems()
	if err != nil {
		log.Printf("Error fetching menu items: %v", err)
	}
	return items, err
}

func (s *MenuService) GetMenuItemByID(id string) (*models.MenuItem, error) {
	log.Printf("Fetching menu item with ID %s", id)
	items, err := s.repo.GetAllItems()
	if err != nil {
		log.Printf("Error fetching menu items: %v", err)
		return nil, err
	}

	for _, item := range items {
		if item.ID == id {
			log.Printf("Menu item with ID %s found", id)
			return &item, nil
		}
	}

	log.Printf("Menu item with ID %s not found", id)
	return nil, models.ErrItemNotFound // You need to define this error
}

func (s *MenuService) UpdateMenuItem(item *models.MenuItem) error {
	log.Printf("Attempting to update menu item with ID %s", item.ID)
	items, err := s.repo.GetAllItems()
	if err != nil {
		log.Printf("Error fetching menu items: %v", err)
		return err
	}

	for i, existingItem := range items {
		if existingItem.ID == item.ID {
			// Update the existing item with fields from the updated item
			log.Printf("Updating menu item with ID %s", item.ID)
			items[i].Name = item.Name
			items[i].Description = item.Description
			items[i].Price = item.Price
			items[i].Ingredients = item.Ingredients
			err = s.repo.SaveItems(items) // Save updated items back to the repository
			if err != nil {
				log.Printf("Error saving updated menu item with ID %s: %v", item.ID, err)
			}
			return err
		}
	}

	log.Printf("Menu item with ID %s not found for update", item.ID)
	return models.ErrItemNotFound // Return error if the item is not found
}

func (s *MenuService) DeleteMenuItem(id string) error {
	log.Printf("Attempting to delete menu item with ID %s", id)
	items, err := s.repo.GetAllItems()
	if err != nil {
		log.Printf("Error fetching menu items: %v", err)
		return err
	}

	for i, existingItem := range items {
		if existingItem.ID == id {
			// Remove the item from the slice
			log.Printf("Deleting menu item with ID %s", id)
			items = append(items[:i], items[i+1:]...) // Remove item at index i
			err = s.repo.SaveItems(items)             // Save updated items back to the repository
			if err != nil {
				log.Printf("Error saving after deletion of menu item with ID %s: %v", id, err)
			}
			return err
		}
	}

	log.Printf("Menu item with ID %s not found for deletion", id)
	return models.ErrItemNotFound // Return error if the item is not found
}
