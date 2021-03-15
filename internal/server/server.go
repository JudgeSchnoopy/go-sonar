package server

import (
	"context"
	"net/http"
	"time"

	"github.com/JudgeSchnoopy/go-sonar/sonar"
)

// Server serves http responses
type Server struct {
	http              *http.Server
	Registry          sonar.Registry
	scheduleStopper   chan bool
	scheduledInterval time.Duration
}

type ServerOption func(*Server)

// New generates a new server
func New(options ...ServerOption) (Server, error) {
	server := Server{
		http: &http.Server{
			Addr:         ":8080",
			ReadTimeout:  time.Second * 15,
			WriteTimeout: time.Second * 15,
		},
		Registry:          sonar.NewRegistry(),
		scheduleStopper:   make(chan bool),
		scheduledInterval: time.Minute * 5,
	}

	// runs the options that can override server defaults
	for _, v := range options {
		v(&server)
	}

	server.http.Handler = server.router()

	return server, nil
}

// Start begins services.
func (server *Server) Start() error {
	server.startScheduler(server.scheduledInterval)

	err := server.http.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

// Stop ends all running Server processes
func (server *Server) Stop(ctx context.Context) {
	server.scheduleStopper <- true
	server.http.Shutdown(ctx)
}
