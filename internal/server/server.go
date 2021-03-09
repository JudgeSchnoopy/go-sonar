package server

import "net/http"

// Server serves http responses
type Server struct {
	*http.Server
}

// New generates a new server
func New() (Server, error) {
	server := Server{
		&http.Server{},
	}
	return server, nil
}
