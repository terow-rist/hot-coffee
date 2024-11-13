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
	UpdateOrder(order *models.Order) error
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

func (s *OrderService) UpdateOrder(order *models.Order) error {
	// Get the existing order to restore inventory
	existingOrder, err := s.orderRepo.GetOrderByID(order.ID)
	if err != nil {
		return errors.New("order not found")
	}
	// Restrict updates to closed orders
	if existingOrder.Status == "closed" {
		return errors.New("cannot update a closed order")
	}
	// Return previous quantities to the inventory
	if err := s.returnInventoryForOrder(existingOrder); err != nil {
		return err
	}

	// Check and deduct inventory for the new order data
	if err := s.checkAndDeductInventoryForOrder(order); err != nil {
		return err
	}

	// Update the order's created time
	order.CreatedAt = time.Now().Format(time.RFC3339)

	// Save the updated order
	return s.orderRepo.UpdateOrder(order)
}

func (s *OrderService) CreateOrder(order *models.Order) error {
	// Check inventory and deduct quantities
	if err := s.checkAndDeductInventoryForOrder(order); err != nil {
		return err
	}

	// Generate unique order ID and set the created time
	order.ID = generateOrderID()
	order.CreatedAt = time.Now().Format(time.RFC3339)

	// Save order
	return s.orderRepo.SaveOrder(order)
}

func (s *OrderService) checkAndDeductInventoryForOrder(order *models.Order) error {
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
			err := s.inventoryService.DeductInventory(ingredient.IngredientID, requiredQty)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *OrderService) returnInventoryForOrder(order *models.Order) error {
	for _, item := range order.Items {
		menuItem, err := s.menuService.GetMenuItemByID(item.ProductID)
		if err != nil {
			return errors.New("product not found in menu")
		}

		for _, ingredient := range menuItem.Ingredients {
			returnQty := ingredient.Quantity * float64(item.Quantity)
			err := s.inventoryService.AddInventory(ingredient.IngredientID, returnQty)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func generateOrderID() string {
	// Generate a unique ID; you may use a more complex logic
	return "order_" + time.Now().Format("20060102150405")
}
