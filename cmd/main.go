package main

import (
	"fmt"
	"hot-coffee/internal/dal"
	"hot-coffee/internal/handler"
	"hot-coffee/internal/service"
	"net/http"
)

// ..
func main() {
	inventoryRepo := &dal.FileInventoryRepository{}
	menuRepo := &dal.FileMenuRepository{}
	orderRepo := &dal.FileOrderRepository{}

	inventoryService := service.NewInventoryService(inventoryRepo)
	menuService := service.NewMenuService(menuRepo)
	orderService := service.NewOrderService(orderRepo, *menuService, *inventoryService)

	inventoryHandler := handler.NewInventoryHandler(inventoryService)
	menuHandler := handler.NewMenuHandler(menuService)
	orderHandler := handler.NewOrderHandler(orderService)

	// Register the /inventory route to handle all methods through ServeHTTP
	http.Handle("/inventory", inventoryHandler)
	http.Handle("/inventory/", inventoryHandler)
	http.Handle("/menu", menuHandler)
	http.Handle("/menu/", menuHandler)
	http.Handle("/orders", orderHandler)
	http.Handle("/orders/", orderHandler)

	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
