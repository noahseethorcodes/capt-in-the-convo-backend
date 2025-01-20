package routes

import (
	"captintheconvo-backend/controllers"

	"github.com/gin-gonic/gin"
)

func TagRoutes(router *gin.Engine) {
	tags := router.Group("/tags")
	{
		tags.GET("/", controllers.GetAllTags)
	}
}
