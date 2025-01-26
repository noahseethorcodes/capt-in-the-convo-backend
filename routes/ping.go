package routes

import (
	"captintheconvo-backend/controllers"

	"github.com/gin-gonic/gin"
)

func PingRoute(router *gin.Engine) {
	// Add your other routes here
	router.GET("/ping", controllers.PingHandler)
}
