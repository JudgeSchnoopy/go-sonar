package server

import (
	"log"
	"net/http"
)

// implement zap logging here

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v : %v : %v\n", r.RemoteAddr, r.Method, r.RequestURI)

		next.ServeHTTP(w, r)
	})
}
