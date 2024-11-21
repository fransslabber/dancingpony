package rest_api

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	sqldb "biz.orcshack/menu/db"
	"github.com/golang-jwt/jwt/v5"
	"github.com/prometheus/client_golang/prometheus"
)

// Prometheus logging
var (
	// Define Prometheus metrics
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of response times for HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

func init() {
	// Register metrics with Prometheus
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

//
// Rate limiting functions
//

func GetClientIP(r *http.Request) string {
	// Check for headers that might contain the client's IP
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// X-Forwarded-For may contain a comma-separated list of IPs, pick the first one
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check for the X-Real-IP header
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// Fallback to using the RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // If parsing fails, return the full RemoteAddr
	}

	return ip
}

// Client holds request timestamps for rate limiting
type Client struct {
	Requests []time.Time
}

// RateLimiter manages rate-limiting logic
type RateLimiter struct {
	clients map[string]*Client
	mu      sync.Mutex
	limit   int           // Maximum requests allowed
	window  time.Duration // Sliding window duration
}

var RateLimit *RateLimiter

// NewRateLimiter initializes a new RateLimiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*Client),
		limit:   limit,
		window:  window,
	}
}

// Allow checks if a client is within the rate limit
func (rl *RateLimiter) Allow(clientIP string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	client, exists := rl.clients[clientIP]
	if !exists {
		// First request from this client
		rl.clients[clientIP] = &Client{Requests: []time.Time{now}}
		return true
	}

	// Filter out timestamps outside the sliding window
	var validRequests []time.Time
	for _, t := range client.Requests {
		if now.Sub(t) <= rl.window {
			validRequests = append(validRequests, t)
		}
	}
	client.Requests = validRequests

	// Check if the client is within the rate limit
	if len(client.Requests) < rl.limit {
		client.Requests = append(client.Requests, now)
		return true
	}

	return false
}

// Cleanup removes old client data periodically
func (rl *RateLimiter) Cleanup(interval time.Duration) {
	for {
		time.Sleep(interval)
		rl.mu.Lock()
		for ip, client := range rl.clients {
			// Remove clients with no recent requests
			if len(client.Requests) == 0 || time.Since(client.Requests[len(client.Requests)-1]) > rl.window {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}

type ResponseWriterWrapper struct {
	http.ResponseWriter
	StatusCode int
}

// Middleware for rate limiting implementation and prometheus metrics
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := GetClientIP(r)
		if !RateLimit.Allow(clientIP) {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Error_Response{ErrorDetails{Code: 429, Message: "Too many requests"}})
			return
		}

		start := time.Now()

		wrapper := &ResponseWriterWrapper{
			ResponseWriter: w,
			StatusCode:     http.StatusOK, // Default to 200
		}

		next.ServeHTTP(wrapper, r)

		// Record metrics after request is complete
		duration := time.Since(start).Seconds()
		status := wrapper.StatusCode
		endpoint := r.URL.Path // Use the full route path as the endpoint

		httpRequestsTotal.WithLabelValues(r.Method, endpoint, http.StatusText(status)).Inc()
		httpRequestDuration.WithLabelValues(r.Method, endpoint).Observe(duration)
	})
}

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
