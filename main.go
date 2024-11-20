package main

import (
	"log"
	"net/http"

	sqldb "biz.orcshack/menu/db"
	rest_api "biz.orcshack/menu/web"
)

// Redirect http calls to https
func redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
	log.Printf("Redirect https://%v", r.Host+r.URL.String())
	http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
}

// Instantiate server with routes and self signed certificate for POC
func main() {
	err := sqldb.CreateDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer sqldb.CloseDB()

	// HTTP redirect server
	go func() {
		http.HandleFunc("/", redirectToHTTPS)
		log.Println("HTTP redirect server running on http://localhost:8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// routes
	http.HandleFunc("GET /{restaurant}/v1/register", rest_api.Register)
	http.HandleFunc("GET /{restaurant}/v1/login", rest_api.Login)

	http.HandleFunc("GET /{restaurant}/v1/list_dishes", rest_api.List_Dishes)
	http.HandleFunc("GET /{restaurant}/v1/view_dish", rest_api.View_Dish)
	http.HandleFunc("GET /{restaurant}/v1/create_dish", rest_api.Create_Dish)
	http.HandleFunc("GET /{restaurant}/v1/delete_dish", rest_api.Delete_Dish)
	http.HandleFunc("GET /{restaurant}/v1/search_dish", rest_api.Search_Dish)

	http.HandleFunc("GET /{restaurant}/v1/rate_dish", rest_api.Rate_Dish)
	http.HandleFunc("POST /{restaurant}/v1/add_dish_image", rest_api.Add_Dish_Image)

	http.HandleFunc("GET /{restaurant}/v1/create_review", rest_api.Create_Review)

	http.HandleFunc("/auth/google/login", rest_api.GoogleLoginHandler)
	http.HandleFunc("/auth/google/callback", rest_api.GoogleCallbackHandler)

	server := &http.Server{
		Addr:    ":4443",
		Handler: http.DefaultServeMux,
	}
	log.Println("HTTPS server running on https://localhost:4443")
	log.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))
}
