package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func main() {
	// Define the secret key (should match the one in your middleware)
	secretKey := "your_secret_key"

	// Define the claims (payload)
	claims := jwt.MapClaims{
		"user_id": 1,                                // Replace with dynamic user ID
		"exp":     time.Now().Add(time.Hour).Unix(), // Token expires in 1 hour
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		fmt.Println("Error generating token:", err)
		return
	}

	// Print the token
	fmt.Println("Generated Token:", tokenString)
}
