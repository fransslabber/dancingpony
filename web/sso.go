package rest_api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOAuthConfig = &oauth2.Config{
	ClientID:     "YOUR_GOOGLE_CLIENT_ID",     // Replace with your Google Client ID
	ClientSecret: "YOUR_GOOGLE_CLIENT_SECRET", // Replace with your Google Client Secret
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

func generateJWT(userInfo map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userInfo["id"], // Use Google user ID or email
		"role":    "customer",
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := googleOAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Get the authorization code from the query
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Authorization code not provided", http.StatusBadRequest)
		return
	}

	// Exchange the code for an access token
	token, err := googleOAuthConfig.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Use the access token to get user info
	client := googleOAuthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Parse user info
	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to parse user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Log or process user info
	log.Printf("User Info: %v\n", userInfo)

	// Generate a JWT for the user
	tokenStr, err := generateJWT(userInfo)
	if err != nil {
		http.Error(w, "Failed to generate JWT: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the JWT to the user
	json.NewEncoder(w).Encode(map[string]string{"token": tokenStr})
}
