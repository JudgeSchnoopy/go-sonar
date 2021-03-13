package server

import (
	"fmt"
	"net/http"

	"github.com/JudgeSchnoopy/go-sonar/sonar"
)

func docsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "docs")
}

func (server *Server) showRegistryHandler(w http.ResponseWriter, r *http.Request) {
	server.Respond(w, server.Registry, http.StatusOK)
}

func (server *Server) registerHandler(w http.ResponseWriter, r *http.Request) {
	var post sonar.Entry
	err := readInput(r, &post)
	if err != nil {
		server.Respond(w, err, http.StatusBadRequest)
		return
	}

	entry := sonar.NewEntry(post.Name, post.Address)
	err = server.Registry.Register(entry)
	if err != nil {
		server.Respond(w, err, http.StatusBadRequest)
		return
	}

	server.Respond(w, server.Registry, http.StatusAccepted)
}
