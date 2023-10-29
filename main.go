package main

import (
	"price-tracker/database"
	"price-tracker/handler"
	"price-tracker/router"
)

func main() {
	// Initialize repository
	productDB := database.NewDB()

	// Initialize use case
	productHandler := handler.NewHandler(productDB)

	// Set up router
	router := router.SetupRouter(productHandler)

	// Run the server
	router.Run(":8080")
}
