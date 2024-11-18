package rest_api

import (
	"encoding/json"
	"fmt"
	"net/http"

	sqldb "biz.orcshack/menu/db"
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

func Register(w http.ResponseWriter, r *http.Request) {
	var user sqldb.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("Malformed JSON in register request: %v", err)}})
	} else {
		err := sqldb.Global_db.Create_user(user.Name, user.Email, user.Password_hash, r.PathValue("restaurant"))
		if err != nil {
			w.WriteHeader(400)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("Register user failed: %v", err)}})
		}

	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user sqldb.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("malformed JSON in login request: %v", err)}})
	} else {
		is_authenticated, err := sqldb.Global_db.Login_user(user.Email, user.Password_hash, r.PathValue("restaurant"))
		if err != nil {
			w.WriteHeader(400)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 400, Message: fmt.Sprintf("Authentication failed: %v", err)}})
		} else {
			if is_authenticated {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(LoginUser_ResponseUser{JWT: "JWT"})
			} else {
				w.WriteHeader(401)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 401, Message: "Authentication failed."}})
			}
		}
	}
}
