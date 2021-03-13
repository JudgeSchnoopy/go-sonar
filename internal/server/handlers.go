package server

import (
	"fmt"
	"net/http"
)

func docsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "docs")
}

func (server *Server) showRegistryHandler(w http.ResponseWriter, r *http.Request) {
	server.Respond(w, server.Registry, http.StatusOK)
}
