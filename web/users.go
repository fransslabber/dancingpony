package rest_api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	sqldb "biz.orcshack/menu/db"
	"github.com/golang-jwt/jwt/v5"
)

type Error_Response struct {
	Error ErrorDetails `json:"error"`
}

type ErrorDetails struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type LoginUser_ResponseUser struct {
	JWT string `json:"jwt"`
}

// Register a new user, no auth required
func Register(w http.ResponseWriter, r *http.Request) {
	var user sqldb.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 004, Message: fmt.Sprintf("Malformed JSON in register request: %v", err)}})
	} else {

		err := sqldb.Global_db.Create_user(user.Name, user.Email, user.Password_hash, r.PathValue("restaurant"))
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 005, Message: fmt.Sprintf("Register user failed: %v", err)}})
		}

	}
}

// Login validation, checks email and password, if ok, returns a JWT, for use in all other calls
func Login(w http.ResponseWriter, r *http.Request) {
	var user sqldb.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 001, Message: fmt.Sprintf("malformed JSON in login request: %v", err)}})
	} else {
		is_authenticated, user, err := sqldb.Global_db.Login_user(user.Email, user.Password_hash, r.PathValue("restaurant"))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 002, Message: fmt.Sprintf("authentication failed: %v", err)}})
		} else {
			if is_authenticated {
				// Return JWT
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"user_id": user.Id,
					"role":    user.Role,
					"exp":     jwt.NewNumericDate(time.Now().Add(time.Hour)),
				})

				tokenStr, err := token.SignedString(jwtSecret)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 002, Message: fmt.Sprintf("JWT token error: %v", err)}})
					return
				}

				json.NewEncoder(w).Encode(LoginUser_ResponseUser{JWT: tokenStr})
				return

			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 401, Message: "authentication failed."}})
			}
		}
	}
}
