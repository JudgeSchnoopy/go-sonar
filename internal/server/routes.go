package server

import "github.com/gorilla/mux"

func (server *Server) router() *mux.Router {
	r := mux.NewRouter()

	r.Use(loggingMiddleware)

	r.HandleFunc("/docs", docsHandler).Methods("GET")
	r.HandleFunc("/registry", server.showRegistryHandler).Methods("GET")
	r.HandleFunc("/register", server.registerHandler).Methods("POST")
	r.HandleFunc("/register", server.removeHandler).Methods("DELETE")

	return r
}
