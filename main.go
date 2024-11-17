package main

import (
	"crypto/ecdsa"
	"net/http"
)

type Error_Response struct {
	Msg string `json:"msg"`
}

var ecdsakey *ecdsa.PrivateKey

func routes() {
	http.HandleFunc("GET /{tenant}/v1/registration", register)
	http.HandleFunc("GET /{tenant}/v1/login", loginUser)
	http.HandleFunc("GET /{tenant}/v1/logout", logoutUser)

	// http.HandleFunc("/api/v1/audit", auditTrail)
	// http.HandleFunc("/api/v1/form", formData)
	// //http.HandleFunc("/api/v1/oauthcallback", OauthCallback)
	// http.HandleFunc("/api/v1/plan_data", planData)
	// http.HandleFunc("/api/v1/form_options", formOptions)
	// http.HandleFunc("/api/v1/form_generator", formGenerator)

}

func main() {
	var err error

	// Set up paths and handlers
	routes()
	// Setup JWT signing key
	// ecdsakey, err = GenerateECDSAKey()
	if err != nil {
		//sugar.Errorf("Could not generate ECDSA key: %+v", err)
		return
	}

	srv := &http.Server{Addr: ":80", Handler: http.DefaultServeMux}

	//sugar.Infof("Web Server started at localhost:80")
	err = srv.ListenAndServe()

	// if err != http.ErrServerClosed {
	// 	sugar.Fatalf("listen: %s", err)
	// }
	// <-exit
	// sugar.Info("Web Server gracefully stopped")
	// exitWebServer <- 1

}
