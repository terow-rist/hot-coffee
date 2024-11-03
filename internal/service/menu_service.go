package service

import "hot-coffee/models"

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
	return s.repo.AddItem(item)
}

func (s *MenuService) GetAllItems() ([]models.MenuItem, error) {
	return s.repo.GetAllItems()
}

func (s *MenuService) GetMenuItemByID(id string) (*models.MenuItem, error) {
	items, err := s.repo.GetAllItems()
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if item.ID == id {
			return &item, nil
		}
	}

	return nil, models.ErrItemNotFound // You need to define this error
}

func (s *MenuService) UpdateMenuItem(item *models.MenuItem) error {
	items, err := s.repo.GetAllItems()
	if err != nil {
		return err
	}

	for i, existingItem := range items {
		if existingItem.ID == item.ID {
			// Update the existing item with fields from the updated item
			items[i].Name = item.Name
			items[i].Description = item.Description
			items[i].Price = item.Price
			items[i].Ingredients = item.Ingredients
			return s.repo.SaveItems(items) // Save updated items back to the repository
		}
	}

	return models.ErrItemNotFound // Return error if the item is not found
}

func (s *MenuService) DeleteMenuItem(id string) error {
	items, err := s.repo.GetAllItems()
	if err != nil {
		return err
	}

	for i, existingItem := range items {
		if existingItem.ID == id {
			// Remove the item from the slice
			items = append(items[:i], items[i+1:]...) // Remove item at index i
			return s.repo.SaveItems(items)            // Save updated items back to the repository
		}
	}

	return models.ErrItemNotFound // Return error if the item is not found
}
