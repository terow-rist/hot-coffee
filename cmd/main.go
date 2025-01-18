package main

import (
	"fmt"
	"hot-coffee/config"
	"hot-coffee/internal/dal"
	"hot-coffee/internal/handler"
	"hot-coffee/internal/service"
	"net/http"
	"os"
	"strconv"
)

// ..
func main() {
	// Validate directory and port
	if err := config.ValidateDirectory(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	// Check if the port is valid
	port, err := strconv.Atoi(config.PortNumber)
	if err != nil || port < 1024 || port > 49151 {
		fmt.Printf("Error: Invalid port number %s\n", config.PortNumber)
		os.Exit(1)

	}
	inventoryRepo := &dal.FileInventoryRepository{}
	menuRepo := &dal.FileMenuRepository{}
	orderRepo := &dal.FileOrderRepository{}

	inventoryService := service.NewInventoryService(inventoryRepo)
	menuService := service.NewMenuService(menuRepo)
	orderService := service.NewOrderService(orderRepo, *menuService, *inventoryService)
	reportsService := service.NewReportsService(orderRepo, *menuService)

	reportsHandler := handler.NewReportsHandler(reportsService)
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

	// Register the new report routes
	http.Handle("/reports/total-sales", reportsHandler)
	http.Handle("/reports/popular-items", reportsHandler)

	fmt.Println("Server is running on port " + config.PortNumber)
	if err := http.ListenAndServe(":"+config.PortNumber, nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
