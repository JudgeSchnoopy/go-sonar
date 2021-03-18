package server

import (
	"fmt"
	"net/http"

	"github.com/JudgeSchnoopy/go-sonar/client"
	"github.com/JudgeSchnoopy/go-sonar/sonar"
)

// get /docs
func docsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "docs")
}

// get /registry
func (server *Server) showRegistryHandler(w http.ResponseWriter, r *http.Request) {
	server.Respond(w, server.Registry, http.StatusOK)
}

// post /register
func (server *Server) registerHandler(w http.ResponseWriter, r *http.Request) {
	var post client.Response
	err := readInput(r, &post)
	if err != nil {
		server.Respond(w, err, http.StatusBadRequest)
		return
	}

	entry := sonar.NewEntry(post)
	err = server.Registry.Register(entry)
	if err != nil {
		server.Respond(w, err, http.StatusBadRequest)
		return
	}

	server.Respond(w, server.Registry, http.StatusAccepted)
}

// delete /register
func (server *Server) removeHandler(w http.ResponseWriter, r *http.Request) {
	var post sonar.Entry
	err := readInput(r, &post)
	if err != nil {
		server.Respond(w, err, http.StatusBadRequest)
		return
	}

	server.Registry.Remove(post)
	server.Respond(w, "ok", http.StatusOK)
}

// post /report
func (server *Server) reportHandler(w http.ResponseWriter, r *http.Request) {
	var post client.Response
	err := readInput(r, &post)
	if err != nil {
		server.Respond(w, err, http.StatusBadRequest)
		return
	}

	entry, err := server.Registry.Get(post.Name)
	if err != nil {
		server.Respond(w, err, http.StatusBadRequest)
	}

	entry.Update(post)

	server.Respond(w, entry, http.StatusAccepted)
}
