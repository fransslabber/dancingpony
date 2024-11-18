package main

import (
	"crypto/ecdsa"
	"fmt"
	"net/http"

	sqldb "biz.orcshack/menu/db"
	rest_api "biz.orcshack/menu/web"
)

type Error_Response struct {
	Msg string `json:"msg"`
}

var ecdsakey *ecdsa.PrivateKey

func routes() {
	http.HandleFunc("GET /{restaurant}/v1/register", rest_api.Register)
	http.HandleFunc("GET /{restaurant}/v1/login", rest_api.Login)

	http.HandleFunc("GET /{restaurant}/v1/list_dishes", rest_api.List_Dishes)
	http.HandleFunc("GET /{restaurant}/v1/view_dish", rest_api.View_Dish)
	http.HandleFunc("GET /{restaurant}/v1/create_dish", rest_api.Create_Dish)
	http.HandleFunc("GET /{restaurant}/v1/delete_dish", rest_api.Delete_Dish)
	http.HandleFunc("GET /{restaurant}/v1/search_dish", rest_api.Search_Dish)
	http.HandleFunc("GET /{restaurant}/v1/rate_dish", rest_api.Rate_Dish)

}

func main() {
	var err error

	// Set up paths and handlers
	routes()
	// Setup JWT signing key
	// ecdsakey, err = GenerateECDSAKey()
	// if err != nil {
	// 	//sugar.Errorf("Could not generate ECDSA key: %+v", err)
	// 	return
	// }
	err = sqldb.CreateDB()
	if err != nil {
		//sugar.Errorf("Could not generate ECDSA key: %+v", err)
		return
	}
	defer sqldb.CloseDB()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}

	// if err != http.ErrServerClosed {
	// 	sugar.Fatalf("listen: %s", err)
	// }
	// <-exit
	// sugar.Info("Web Server gracefully stopped")
	// exitWebServer <- 1

}
