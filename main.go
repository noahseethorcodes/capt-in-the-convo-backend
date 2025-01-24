package main

import (
	"captintheconvo-backend/database"
	"captintheconvo-backend/models"
	"captintheconvo-backend/routes"
	"fmt"
	"log"
	"os"

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

	host := os.Getenv("APP_HOST") // Render sets the PORT environment variable
	if host == "" {
		host = "localhost" // Default to 8080 for local testing
	}
	port := os.Getenv("APP_PORT") // Render sets the PORT environment variable
	if port == "" {
		port = "8080" // Default to 8080 for local testing
	}

	fmt.Println("APP_HOST:", os.Getenv("APP_HOST"))
	fmt.Println("APP_PORT:", os.Getenv("APP_PORT"))

	// Start the server
	router.Run(fmt.Sprintf("%s:%s", host, port))
}
