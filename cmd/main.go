package main

import (
	"fmt"
	"net/http"

	"hot-coffee/internal/handler"
	"hot-coffee/internal/service" // Adjust the import path
)

func main() {
	orderService := &service.OrderService{}
	orderHandler := handler.NewOrderHandler(orderService) // Pass pointer directly
	fmt.Println("Server is running on port 8080")         // Change the message
	http.HandleFunc("/orders", orderHandler.CreateOrder)

	err := http.ListenAndServe(":8080", nil) // Handle potential errors
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
