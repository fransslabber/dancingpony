package rest_api

import (
	"encoding/json"
	"fmt"
	"net/http"
	sqldb "wo-infield-service/db"
)

// ///////////////////////////////////////////////////////////////////////////////////
// // Register Device
type RegisterDevice_Reponse struct {
	MobileDeviceId string `json:"mobile_device_id"`
}

type RegisterDevice_Request struct {
	RegistrationToken string `json:"registration_token"`
}

func registerDevice(w http.ResponseWriter, r *http.Request) {
	var req RegisterDevice_Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("Failed to comprehend request: %v", err)})
	}

	if err := sqldb.CheckRegToken(req.RegistrationToken); err == nil {

		err, uuid_ := sqldb.RegisterDevice()
		if err != nil {
			w.WriteHeader(400)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("Failed to register Device(DB): %v", err)})
		} else {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(RegisterDevice_Reponse{uuid_})
		}

	} else {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Error_Response{fmt.Sprintf("Failed to register Device: %v", err.Error())})
	}
}
