package routes

import (
	"captintheconvo-backend/controllers"
	"captintheconvo-backend/middleware"

	"github.com/gin-gonic/gin"
)

func CommentRoutes(router *gin.Engine) {
	// Public route for retrieving comments
	router.GET("/comments/:thread_id", controllers.GetCommentsByThread) // Public route

	// Protected routes for thread creation
	protected := router.Group("/comments")
	protected.Use(middleware.AuthMiddleware()) // Protect these routes with authentication
	{
		protected.POST("", controllers.CreateComment)
		protected.DELETE("/:id", controllers.DeleteComment)
	}

}
