package main

import (
	"log"
	"net/http"
	"time"

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

	// Create a new rate limiter (5 requests per 10 seconds)
	rest_api.RateLimit = rest_api.NewRateLimiter(5, 10*time.Second)

	// Start cleanup goroutine
	go rest_api.RateLimit.Cleanup(1 * time.Minute)

	// HTTP redirect server
	go func() {
		http.HandleFunc("/", redirectToHTTPS)
		log.Println("HTTP redirect server running on http://localhost:8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// routes with rate limiting per IP address
	http.Handle("GET /{restaurant}/v1/register", rest_api.RateLimitMiddleware(http.HandlerFunc(rest_api.Register)))
	http.Handle("GET /{restaurant}/v1/login", rest_api.RateLimitMiddleware(http.HandlerFunc(rest_api.Login)))

	http.Handle("GET /{restaurant}/v1/list_dishes", rest_api.RateLimitMiddleware(http.HandlerFunc(rest_api.List_Dishes)))
	http.Handle("GET /{restaurant}/v1/view_dish", rest_api.RateLimitMiddleware(http.HandlerFunc(rest_api.View_Dish)))
	http.Handle("GET /{restaurant}/v1/create_dish", rest_api.RateLimitMiddleware(http.HandlerFunc(rest_api.Create_Dish)))
	http.Handle("GET /{restaurant}/v1/delete_dish", rest_api.RateLimitMiddleware(http.HandlerFunc(rest_api.Delete_Dish)))
	http.Handle("GET /{restaurant}/v1/search_dish", rest_api.RateLimitMiddleware(http.HandlerFunc(rest_api.Search_Dish)))

	http.Handle("GET /{restaurant}/v1/rate_dish", rest_api.RateLimitMiddleware(http.HandlerFunc(rest_api.Rate_Dish)))
	http.Handle("POST /{restaurant}/v1/add_dish_image", rest_api.RateLimitMiddleware(http.HandlerFunc(rest_api.Add_Dish_Image)))

	http.Handle("GET /{restaurant}/v1/create_review", rest_api.RateLimitMiddleware(http.HandlerFunc(rest_api.Create_Review)))

	http.HandleFunc("/auth/google/login", rest_api.GoogleLoginHandler)
	http.HandleFunc("/auth/google/callback", rest_api.GoogleCallbackHandler)

	server := &http.Server{
		Addr:    ":4443",
		Handler: http.DefaultServeMux,
	}
	log.Println("HTTPS server running on https://localhost:4443")
	log.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))
}
