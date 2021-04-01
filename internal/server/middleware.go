package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// implement zap logging here

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v : %v : %v\n", r.RemoteAddr, r.Method, r.RequestURI)

		next.ServeHTTP(w, r)
	})
}

// Returns a new TimeoutMiddleware with specified server timeout
func NewTimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			toHandler := http.TimeoutHandler(next, timeout-1*time.Second, "error: request timed out")

			fmt.Println("passing through timeout middleware")

			toHandler.ServeHTTP(w, r)
		})
	}
}
