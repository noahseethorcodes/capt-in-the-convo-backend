package controllers

import (
	"captintheconvo-backend/database"
	"captintheconvo-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = "your_secret_key" // Replace with your actual secret key

func Register(context *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	user := models.User{Username: input.Username, Password: string(hashedPassword)}

	if err := database.DB.Create(&user).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func Login(context *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	var user models.User
	// Check if the username exists
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Username not found"})
		return
	}

	// Check if the password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Generate short-lived access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(), // Token valid for 15 minutes
	})

	accessTokenString, err := accessToken.SignedString([]byte(secretKey))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	// Generate long-lived refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(), // Token valid for 7 days
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(secretKey))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Return both tokens in the response
	context.JSON(http.StatusOK, gin.H{
		"accessToken":  accessTokenString,
		"refreshToken": refreshTokenString,
		"userID":       user.ID,
	})
}

func Refresh(context *gin.Context) {
	var input struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Parse and validate the refresh token
	token, err := jwt.Parse(input.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	userID := uint(claims["user_id"].(float64))

	// Issue a new access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(), // Token valid for 15 minutes
	})

	accessTokenString, err := accessToken.SignedString([]byte(secretKey))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new access token"})
		return
	}

	// Issue a new refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(), // Extend token validity
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(secretKey))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new refresh token"})
		return
	}

	// Return the new tokens
	context.JSON(http.StatusOK, gin.H{
		"accessToken":  accessTokenString,
		"refreshToken": refreshTokenString,
	})
}
