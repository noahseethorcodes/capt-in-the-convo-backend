package main

import (
	"captintheconvo-backend/database"
	"captintheconvo-backend/models"
	"captintheconvo-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

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

	// router.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:3000"},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// }))

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"}, // Allow all headers for simplicity
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
	}))

	// Register routes
	routes.AuthRoutes(router)
	routes.ThreadRoutes(router)
	routes.AdminRoutes(router)
	routes.CommentRoutes(router)
	routes.TagRoutes(router)

	// Start the server
	router.Run() // Default port: 8080
}
