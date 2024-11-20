package rest_api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("Eyeball of the Gargoyle") // Secret key for JWT

// Validate JWT token, return user id and role from token
func ValidateJWT(r *http.Request) (uint32, string, error) {
	authHeader := r.Header.Get("Authorization")
	fmt.Printf("Auth token %v\n", authHeader)
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return 0, "", fmt.Errorf("authorization not present")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	// Parse and validate the token
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return 0, "", fmt.Errorf("unauthorized: Invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", fmt.Errorf("unauthorized: Invalid claims")
	}

	// Validate `user_id`
	userID, ok := claims["user_id"].(float64) // JWT stores numbers as float64
	if !ok || userID <= 0 {
		return 0, "", fmt.Errorf("unauthorized: Invalid user ID")
	}

	// Validate expiration (`exp`) claim
	exp, ok := claims["exp"].(float64) // JWT stores numeric values as float64
	if !ok {
		return 0, "", fmt.Errorf("unauthorized: Expiration claim missing")
	}

	// Check if the token has expired
	if time.Now().Unix() > int64(exp) {
		return 0, "", fmt.Errorf("unauthorized: Token has expired")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return 0, "", fmt.Errorf("unauthorized: No role not found")
	}

	return uint32(userID), role, nil
}
