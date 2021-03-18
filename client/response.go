package client

import (
	"encoding/json"
	"net/http"
)

// Response is the expected response from a /sonar endpoint
type Response struct {
	Name         string
	Address      string
	Healthy      bool
	Dependencies dependencies
}

// SonarHandler provides a pre-built response handler for a /sonar endpoint
func (resp *Response) SonarHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	jsonResponse, err := json.MarshalIndent(*resp, "", "    ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		w.Write([]byte(err.Error()))
	}

	w.Write(jsonResponse)
}
