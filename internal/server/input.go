package server

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

func readInput(r *http.Request, v interface{}) error {
	body, err := ioutil.ReadAll(io.Reader(r.Body))
	if err != nil {
		return err
	}

	if err := r.Body.Close(); err != nil {
		return err
	}

	if err := json.Unmarshal(body, v); err != nil {
		return err
	}

	return nil
}
