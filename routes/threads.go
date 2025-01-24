package routes

import (
	"captintheconvo-backend/controllers"
	"captintheconvo-backend/middleware"

	"github.com/gin-gonic/gin"
)

func ThreadRoutes(router *gin.Engine) {
	// Public route for retrieving threads
	router.GET("/threads", controllers.GetThreads)

	// Public route for retrieving a single thread by ID
	router.GET("/threads/:id", controllers.GetThreadByID)

	// Protected routes for thread creation
	protected := router.Group("/threads")
	protected.Use(middleware.AuthMiddleware()) // Protect these routes with authentication
	{
		protected.POST("", controllers.CreateThread)
		protected.DELETE("/:id", controllers.DeleteThread)
	}
}
