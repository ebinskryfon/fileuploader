package main

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func main() {
	// Default values
	userID := "dev-user-123"
	secret := "your-secret-key-change-in-production"
	expiration := 24 * time.Hour

	// Create claims
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "fileuploader-service",
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatal("Failed to generate token:", err)
	}

	fmt.Println("Generated JWT Token for development:")
	fmt.Println("User ID:", userID)
	fmt.Println("Expires:", claims.ExpiresAt.Time.Format(time.RFC3339))
	fmt.Println()
	fmt.Println("Token:")
	fmt.Println(tokenString)
	fmt.Println()
	fmt.Println("Usage example:")
	fmt.Printf("curl -H \"Authorization: Bearer %s\" http://localhost:8080/api/v1/upload\n", tokenString)
}
