package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/JudgeSchnoopy/go-sonar/sonar"
	"github.com/gorilla/mux"
)

// Server serves http responses
type Server struct {
	http            *http.Server
	Registry        sonar.Registry
	scheduleStopper chan bool
}

type Config struct {
	ScheduleInterval time.Duration
	Port             int
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
}

// New generates a new server
func New(config Config) (Server, error) {
	server := Server{
		http: &http.Server{
			Addr:         fmt.Sprintf(":%v", config.Port),
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
		},
		Registry:        sonar.NewRegistry(),
		scheduleStopper: make(chan bool),
	}

	server.http.Handler = server.router()
	server.startScheduler(config.ScheduleInterval)

	return server, nil
}

// Start begins the listening service.
func (server *Server) Start() error {
	err := server.http.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

// Stop shuts down the listening service
func (server *Server) Stop(ctx context.Context) {
	server.scheduleStopper <- true
	server.http.Shutdown(ctx)
}

func (server *Server) router() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.HandleFunc("/docs", docsHandler).Methods("GET")
	r.HandleFunc("/registry", server.showRegistryHandler).Methods("GET")
	r.HandleFunc("/register", server.registerHandler).Methods("POST")

	return r
}
