package controllers

import (
	"captintheconvo-backend/database"
	"captintheconvo-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateTag(context *gin.Context) {
	var input struct {
		Name string `json:"name"`
	}
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag := models.Tag{Name: input.Name}
	if err := database.DB.Create(&tag).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Tag created successfully", "tag": tag})
}

func UpdateTag(context *gin.Context) {
	var tag models.Tag
	if err := database.DB.First(&tag, context.Param("id")).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	var input struct {
		Name     *string `json:"name"`
		IsActive *bool   `json:"is_active"`
	}
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != nil {
		tag.Name = *input.Name
	}
	if input.IsActive != nil {
		tag.IsActive = *input.IsActive
	}

	if err := database.DB.Save(&tag).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tag"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Tag updated successfully", "tag": tag})
}

func DeleteTag(context *gin.Context) {
	if err := database.DB.Delete(&models.Tag{}, context.Param("id")).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tag"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Tag deleted successfully"})
}
