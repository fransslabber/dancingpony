package main

import (
	"log"
	"net/http"

	sqldb "biz.orcshack/menu/db"
	rest_api "biz.orcshack/menu/web"
)

func main() {
	err := sqldb.CreateDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer sqldb.CloseDB()

	// public routes
	http.HandleFunc("GET /{restaurant}/v1/register", rest_api.Register)
	http.HandleFunc("GET /{restaurant}/v1/login", rest_api.Login)
	http.HandleFunc("/auth/google/login", rest_api.GoogleLoginHandler)
	http.HandleFunc("/auth/google/callback", rest_api.GoogleCallbackHandler)

	// Authenticated routes
	protected := http.NewServeMux()
	protected.HandleFunc("GET /{restaurant}/v1/list_dishes", rest_api.List_Dishes)
	protected.HandleFunc("GET /{restaurant}/v1/view_dish", rest_api.View_Dish)
	protected.HandleFunc("GET /{restaurant}/v1/create_dish", rest_api.Create_Dish)
	protected.HandleFunc("GET /{restaurant}/v1/delete_dish", rest_api.Delete_Dish)
	protected.HandleFunc("GET /{restaurant}/v1/search_dish", rest_api.Search_Dish)
	protected.HandleFunc("GET /{restaurant}/v1/rate_dish", rest_api.Rate_Dish)

	http.Handle("GET /{restaurant}/v1/list_dishes", rest_api.AuthMiddleware(protected))
	http.Handle("GET /{restaurant}/v1/view_dish", rest_api.AuthMiddleware(protected))
	http.Handle("GET /{restaurant}/v1/create_dish", rest_api.AuthAdminMiddleware(protected))
	http.Handle("GET /{restaurant}/v1/delete_dish", rest_api.AuthAdminMiddleware(protected))
	http.Handle("GET /{restaurant}/v1/search_dish", rest_api.AuthMiddleware(protected))
	http.Handle("GET /{restaurant}/v1/rate_dish", rest_api.AuthMiddleware(protected))

	log.Fatal(http.ListenAndServe(":8080", nil))

}
