package main

import (
	"fmt"
	"net/http"

	"hot-coffee/internal/dal"
	"hot-coffee/internal/handler"
	"hot-coffee/internal/service"
)

func main() {
	inventoryRepo := &dal.FileInventoryRepository{}
	menuRepo := &dal.FileMenuRepository{}

	inventoryService := service.NewInventoryService(inventoryRepo)
	menuService := service.NewMenuService(menuRepo)

	inventoryHandler := handler.NewInventoryHandler(inventoryService)
	menuHandler := handler.NewMenuHandler(menuService)

	// Register the /inventory route to handle all methods through ServeHTTP
	http.Handle("/inventory", inventoryHandler)
	http.Handle("/inventory/", inventoryHandler)
	http.Handle("/menu", menuHandler)

	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
