package server

import (
	"encoding/json"
	"net/http"
)

func (server *Server) Respond(w http.ResponseWriter, msg interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(msg); err != nil {
		panic(err)
	}
}
