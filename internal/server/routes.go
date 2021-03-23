package server

import (
	"github.com/gorilla/mux"
)

func (server *Server) sonarRouter() *mux.Router {
	r := mux.NewRouter()

	r.Use(
		loggingMiddleware,
		NewTimeoutMiddleware(server.http.WriteTimeout),
	)

	r.HandleFunc("/docs", docsHandler).Methods("GET")
	r.HandleFunc("/registry", server.showRegistryHandler).Methods("GET")
	r.HandleFunc("/register", server.registerHandler).Methods("POST")
	r.HandleFunc("/register", server.removeHandler).Methods("DELETE")
	r.HandleFunc("/report", server.reportHandler).Methods("POST")

	return r
}
