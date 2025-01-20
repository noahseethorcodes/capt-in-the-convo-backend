package controllers

import (
	"captintheconvo-backend/database"
	"captintheconvo-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllTags(context *gin.Context) {
	var tags []models.Tag
	isActive := context.DefaultQuery("is_active", "") // Optional: "true" or "false"
	name := context.Query("name")                     // Optional: Filter by tag name

	query := database.DB

	// Filter by is_active if provided
	if isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}

	// Filter by name if provided
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	// Fetch tags
	if err := query.Find(&tags).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tags"})
		return
	}

	// Simplify the response to only include ID and Name
	simplifiedTags := []gin.H{}
	for _, tag := range tags {
		simplifiedTags = append(simplifiedTags, gin.H{
			"ID":   tag.ID,
			"Name": tag.Name,
		})
	}

	context.JSON(http.StatusOK, simplifiedTags)
}
