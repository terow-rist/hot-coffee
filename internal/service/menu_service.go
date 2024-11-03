package service

import "hot-coffee/models"

type MenuRepository interface {
	AddItem(item *models.MenuItem) error
	GetAllItems() ([]models.MenuItem, error)
	//...
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
