package rest_api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("Eyeball of the Gargoyle") // Secret key for JWT

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		fmt.Printf("Auth token %v\n", authHeader)
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
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
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized: Invalid claims", http.StatusUnauthorized)
			return
		}

		// Validate `user_id`
		userID, ok := claims["user_id"].(float64) // JWT stores numbers as float64
		if !ok || userID <= 0 {
			http.Error(w, "Unauthorized: Invalid user ID", http.StatusUnauthorized)
			return
		}

		// Validate expiration (`exp`) claim
		exp, ok := claims["exp"].(float64) // JWT stores numeric values as float64
		if !ok {
			http.Error(w, "Unauthorized: Expiration claim missing", http.StatusUnauthorized)
			return
		}

		// Check if the token has expired
		if time.Now().Unix() > int64(exp) {
			http.Error(w, "Unauthorized: Token has expired", http.StatusUnauthorized)
			return
		}

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

func AuthAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		fmt.Printf("Auth token %v\n", authHeader)
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
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
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized: Invalid claims", http.StatusUnauthorized)
			return
		}

		// Validate `user_id`
		userID, ok := claims["user_id"].(float64) // JWT stores numbers as float64
		if !ok || userID <= 0 {
			http.Error(w, "Unauthorized: Invalid user ID", http.StatusUnauthorized)
			return
		}

		// Mock user role validation (e.g., check against database)
		role, ok := claims["role"].(string)
		if !ok || role == "customer" {
			http.Error(w, "Unauthorized: No permission not found", http.StatusUnauthorized)
			return
		}

		// Validate expiration (`exp`) claim
		exp, ok := claims["exp"].(float64) // JWT stores numeric values as float64
		if !ok {
			http.Error(w, "Unauthorized: Expiration claim missing", http.StatusUnauthorized)
			return
		}

		// Check if the token has expired
		if time.Now().Unix() > int64(exp) {
			http.Error(w, "Unauthorized: Token has expired", http.StatusUnauthorized)
			return
		}

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}
