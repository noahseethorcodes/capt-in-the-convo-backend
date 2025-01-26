package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PingHandler returns a simple response for testing
func PingHandler(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
