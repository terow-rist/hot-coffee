package service

import (
	"fmt"
	"hot-coffee/models"
)

type ReportsService struct {
	orderRepo   OrderRepository
	menuService MenuService
}

func NewReportsService(orderRepo OrderRepository, menuService MenuService) *ReportsService {
	return &ReportsService{
		orderRepo:   orderRepo,
		menuService: menuService,
	}
}

// GetTotalSales calculates the total sales amount from all orders
func (s *ReportsService) GetTotalSales() (float64, error) {
	orders, err := s.orderRepo.GetAllOrders()
	if err != nil {
		return 0, err
	}

	var totalSales float64
	for _, order := range orders {
		// Assuming each order item has a price in the menu
		for _, item := range order.Items {
			menuItem, err := s.menuService.GetMenuItemByID(item.ProductID)
			if err != nil {
				return 0, err
			}
			totalSales += menuItem.Price * float64(item.Quantity)
		}
	}
	return totalSales, nil
}

// GetPopularItems returns the list of most ordered menu items
func (s *ReportsService) GetPopularItems() ([]models.MenuItem, error) {
	orders, err := s.orderRepo.GetAllOrders()
	if err != nil {
		return nil, err
	}

	itemCounts := make(map[string]int)

	// Count the occurrences of each item in the orders
	for _, order := range orders {
		for _, item := range order.Items {
			itemCounts[item.ProductID] += item.Quantity
		}
	}

	var popularItems []models.MenuItem
	for productID, count := range itemCounts {
		menuItem, err := s.menuService.GetMenuItemByID(productID)
		if err != nil {
			return nil, err
		}
		// Add the sold quantity to the menu item for reporting purposes
		menuItem.Description = fmt.Sprintf("Quantity Sold: %d", count)
		popularItems = append(popularItems, *menuItem)
	}

	return popularItems, nil
}
