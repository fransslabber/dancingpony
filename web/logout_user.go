package rest_api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	sqldb "wo-infield-service/db"
)

// ///////////////////////////////////////////////////////////////////////////////////
// // Login User
type LogoutUser_Request struct {
	Username       string `json:"username"`
	MobileDeviceId string `json:"mobile_device_id"`
}

func logoutUser(w http.ResponseWriter, r *http.Request) {
	var req LogoutUser_Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("Failed to comprehend request: %v", err)})
	} else {
		if userid, err := AuthenticateJWT(strings.Split(r.Header["Authorization"][0], " ")[1], req.MobileDeviceId, ecdsakey); err != nil {
			w.WriteHeader(400)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Error_Response{err.Error()})
		} else {
			if err := sqldb.LogoutUser(userid, req.MobileDeviceId); err != nil {
				w.WriteHeader(400)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(Error_Response{err.Error()})
			}
		}
	}
}
