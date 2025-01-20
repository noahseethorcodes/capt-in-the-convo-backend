package routes

import (
	"captintheconvo-backend/controllers"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(router *gin.Engine) {
	admin := router.Group("/admin")
	{
		admin.POST("/tags", controllers.CreateTag)
		admin.PUT("/tags/:id", controllers.UpdateTag)
		admin.DELETE("/tags/:id", controllers.DeleteTag)
	}
}
