package controllers

import (
	"captintheconvo-backend/database"
	"captintheconvo-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateComment(context *gin.Context) {
	_, userExists := context.Get("user_id")
	if !userExists {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		Content  string `json:"content"`
		UserID   uint   `json:"user_id"`
		ThreadID uint   `json:"thread_id"`
	}

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate if the thread exists
	var thread models.Thread
	if err := database.DB.First(&thread, input.ThreadID).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Thread not found"})
		return
	}

	// Create the comment
	comment := models.Comment{
		Content:  input.Content,
		UserID:   input.UserID,
		ThreadID: input.ThreadID,
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Comment added successfully", "comment": comment})
}

func GetCommentsByThread(context *gin.Context) {
	threadID := context.Param("thread_id")

	// Fetch comments for the thread
	var comments []models.Comment
	if err := database.DB.Where("thread_id = ?", threadID).Order("created_at DESC").Find(&comments).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comments"})
		return
	}

	// Enrich comments with Username and simplify the response
	enrichedComments := []gin.H{}
	for _, comment := range comments {
		// Fetch username
		var user models.User
		if err := database.DB.First(&user, comment.UserID).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user for comment"})
			return
		}

		// Build simplified comment structure
		enrichedComments = append(enrichedComments, gin.H{
			"ID":        comment.ID,
			"CreatedAt": comment.CreatedAt,
			"Content":   comment.Content,
			"UserID":    comment.UserID,
			"Username":  user.Username,
			"ThreadID":  comment.ThreadID,
		})
	}

	context.JSON(http.StatusOK, enrichedComments)
}

func DeleteComment(context *gin.Context) {
	commentID := context.Param("id")

	// Fetch the comment to verify ownership
	var comment models.Comment
	if err := database.DB.First(&comment, commentID).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Verify the user is the owner of the comment
	userID, exists := context.Get("user_id")
	if !exists || comment.UserID != userID.(uint) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to delete this comment"})
		return
	}

	// Delete the comment
	if err := database.DB.Delete(&comment).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
