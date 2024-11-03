package main

import (
	"fmt"
	"net/http"

	"hot-coffee/internal/dal"
	"hot-coffee/internal/handler"
	"hot-coffee/internal/service"
)

// main.go
func main() {
	// Initialize the inventory repository and service
	inventoryRepo := &dal.FileInventoryRepository{}
	inventoryService := service.NewInventoryService(inventoryRepo)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)

	// Register the /inventory route to handle both POST and GET requests
	http.Handle("/inventory", inventoryHandler)

	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
