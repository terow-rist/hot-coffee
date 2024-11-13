package service

import (
	"errors"
	"hot-coffee/models"
	"time"
)

type OrderRepository interface {
	SaveOrder(order *models.Order) error // Add this line
	GetAllOrders() ([]models.Order, error)
	GetOrderByID(id string) (*models.Order, error)
}

type OrderService struct {
	orderRepo        OrderRepository
	menuService      MenuService
	inventoryService InventoryService
}

func NewOrderService(orderRepo OrderRepository, menuService MenuService, inventoryService InventoryService) *OrderService {
	return &OrderService{orderRepo, menuService, inventoryService}
}

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	return s.orderRepo.GetAllOrders()
}

func (s *OrderService) GetOrderByID(id string) (*models.Order, error) {
	return s.orderRepo.GetOrderByID(id)
}

func (s *OrderService) CreateOrder(order *models.Order) error {
	// Check if all items are available in the menu
	for _, item := range order.Items {
		menuItem, err := s.menuService.GetMenuItemByID(item.ProductID)
		if err != nil {
			return errors.New("product not found in menu")
		}

		// Check inventory availability
		for _, ingredient := range menuItem.Ingredients {
			inventoryItem, err := s.inventoryService.GetInventoryItemByID(ingredient.IngredientID)
			if err != nil {
				return errors.New("ingredient not found in inventory")
			}
			requiredQty := ingredient.Quantity * float64(item.Quantity)
			if inventoryItem.Quantity < requiredQty {
				return errors.New("insufficient ingredient quantity for " + ingredient.IngredientID)
			}
		}
	}

	// Deduct inventory quantities
	for _, item := range order.Items {
		menuItem, _ := s.menuService.GetMenuItemByID(item.ProductID)
		for _, ingredient := range menuItem.Ingredients {
			requiredQty := ingredient.Quantity * float64(item.Quantity)
			s.inventoryService.DeductInventory(ingredient.IngredientID, requiredQty)
		}
	}

	// Generate unique order ID
	order.ID = generateOrderID()
	order.CreatedAt = time.Now().Format(time.RFC3339)

	// Save order
	return s.orderRepo.SaveOrder(order)
}

func generateOrderID() string {
	// Generate a unique ID; you may use a more complex logic
	return "order_" + time.Now().Format("20060102150405")
}
