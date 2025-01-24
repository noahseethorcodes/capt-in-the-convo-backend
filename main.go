package main

import (
	"captintheconvo-backend/database"
	"captintheconvo-backend/models"
	"captintheconvo-backend/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	// Connect to the database
	database.ConnectDB()

	// Auto-migrate models
	database.DB.AutoMigrate(
		&models.User{},
		&models.Thread{},
		&models.Tag{},
		&models.Comment{},
	)

	// Create a router
	router := gin.Default()

	// Register routes
	routes.AuthRoutes(router)
	routes.ThreadRoutes(router)
	routes.AdminRoutes(router)
	routes.CommentRoutes(router)
	routes.TagRoutes(router)

	// Start the server
	router.Run() // Default port: 8080
}
