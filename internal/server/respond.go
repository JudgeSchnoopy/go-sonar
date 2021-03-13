package server

import (
	"encoding/json"
	"net/http"
)

func (server *Server) Respond(w http.ResponseWriter, msg interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(statusCode)

	jsonResponse, err := json.MarshalIndent(msg, "", "    ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		w.Write([]byte(err.Error()))
	}

	w.Write(jsonResponse)
}
