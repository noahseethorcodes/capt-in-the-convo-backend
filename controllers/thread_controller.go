package controllers

import (
	"captintheconvo-backend/database"
	"captintheconvo-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EnrichedThread struct {
	ID            uint         `json:"ID"`
	Title         string       `json:"Title"`
	Content       string       `json:"Content"`
	UserID        uint         `json:"UserID"`
	Username      string       `json:"Username"`
	CommentsCount int          `json:"CommentsCount"`
	CreatedAt     string       `json:"CreatedAt"`
	Tags          []models.Tag `json:"Tags"`
}

func CreateThread(context *gin.Context) {
	userID, userExists := context.Get("user_id")
	if !userExists {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		Title   string   `json:"title"`
		Content string   `json:"content"`
		UserID  uint     `json:"user_id"`
		Tags    []string `json:"tags"`
	}

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate tags
	var validTags []models.Tag
	if len(input.Tags) > 0 {
		if err := database.DB.Where("name IN ? AND is_active = ?", input.Tags, true).Find(&validTags).Error; err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tags"})
			return
		}

		// Check if all provided tags exist and are valid
		if len(validTags) != len(input.Tags) {
			context.JSON(http.StatusBadRequest, gin.H{"error": "Some tags are invalid or inactive"})
			return
		}
	}

	// Create thread
	thread := models.Thread{
		Title:   input.Title,
		Content: input.Content,
		UserID:  userID.(uint),
	}

	if err := database.DB.Create(&thread).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create thread"})
		return
	}

	// Associate tags
	database.DB.Model(&thread).Association("Tags").Append(&validTags)

	context.JSON(http.StatusOK, gin.H{"message": "Thread created successfully", "thread": thread})
}

func GetThreads(context *gin.Context) {
	var threads []models.Thread
	tags := context.QueryArray("tag") // Retrieve multiple `tag` parameters
	userID := context.Query("user_id")

	query := database.DB.Preload("Tags", "is_active = ?", true)

	// Filter by tags if provided
	if len(tags) > 0 {
		query = query.Joins("JOIN thread_tags ON thread_tags.thread_id = threads.id").
			Joins("JOIN tags ON tags.id = thread_tags.tag_id").
			Where("tags.name IN ?", tags)
	}

	// Filter by user ID if provided
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// Fetch threads
	if err := query.Find(&threads).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve threads"})
		return
	}

	// Enrich threads with additional data
	var enrichedThreads []gin.H
	for _, thread := range threads {
		// Fetch username
		var user models.User
		if err := database.DB.First(&user, thread.UserID).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user for thread"})
			return
		}

		// Count comments for the thread
		var commentsCount int64
		if err := database.DB.
			Model(&models.Comment{}).
			Where("thread_id = ?", thread.ID).
			Count(&commentsCount).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count comments"})
			return
		}

		// Simplify Tags
		var simplifiedTags []gin.H
		for _, tag := range thread.Tags {
			simplifiedTags = append(simplifiedTags, gin.H{
				"ID":   tag.ID,
				"Name": tag.Name,
			})
		}

		// Append enriched thread data
		enrichedThreads = append(enrichedThreads, gin.H{
			"ID":            thread.ID,
			"Title":         thread.Title,
			"Content":       thread.Content,
			"UserID":        thread.UserID,
			"Username":      user.Username,
			"CommentsCount": commentsCount,
			"CreatedAt":     thread.CreatedAt,
			"Tags":          simplifiedTags,
		})
	}

	context.JSON(http.StatusOK, enrichedThreads)
}

func GetThreadByID(context *gin.Context) {
	id := context.Param("id")
	var thread models.Thread
	if err := database.DB.Preload("Tags", "is_active = ?", true).First(&thread, id).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Thread not found"})
		return
	}

	// Fetch username
	var user models.User
	if err := database.DB.First(&user, thread.UserID).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user for thread"})
		return
	}

	// Count comments for the thread
	var commentsCount int64
	if err := database.DB.Model(&models.Comment{}).Where("thread_id = ?", thread.ID).Count(&commentsCount).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count comments"})
		return
	}

	// Simplify Tags
	var simplifiedTags []gin.H
	for _, tag := range thread.Tags {
		simplifiedTags = append(simplifiedTags, gin.H{
			"ID":   tag.ID,
			"Name": tag.Name,
		})
	}

	// Build enriched thread data
	enrichedThread := gin.H{
		"ID":            thread.ID,
		"Title":         thread.Title,
		"Content":       thread.Content,
		"UserID":        thread.UserID,
		"Username":      user.Username,
		"CommentsCount": commentsCount,
		"CreatedAt":     thread.CreatedAt,
		"Tags":          simplifiedTags,
	}

	context.JSON(http.StatusOK, enrichedThread)
}

func DeleteThread(context *gin.Context) {
	threadID := context.Param("id")

	// Fetch the thread to verify ownership
	var thread models.Thread
	if err := database.DB.First(&thread, threadID).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Thread not found"})
		return
	}

	// Verify the user is the owner of the thread
	userID, exists := context.Get("user_id")
	if !exists || thread.UserID != userID.(uint) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to delete this thread"})
		return
	}

	// Delete the thread
	if err := database.DB.Delete(&thread).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete thread"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Thread deleted successfully"})
}
