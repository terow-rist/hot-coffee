package main

import (
	"fmt"
	"net/http"

	"hot-coffee/internal/dal"
	"hot-coffee/internal/handler"
	"hot-coffee/internal/service"
)

func main() {
	// Initialize the inventory repository and service
	inventoryRepo := &dal.FileInventoryRepository{}                   // Ensure this is defined correctly in dal
	inventoryService := service.NewInventoryService(inventoryRepo)    // Use the NewInventoryService constructor
	inventoryHandler := handler.NewInventoryHandler(inventoryService) // Pass the service to the handler

	// Register the /inventory route
	http.HandleFunc("/inventory", inventoryHandler.AddInventoryItem)

	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
