package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func Register(r Response, server string) error {
	jsonResponse, err := json.Marshal(r)
	if err != nil {
		return err
	}

	resp, err := http.Post(server, "application/json", bytes.NewBuffer(jsonResponse))
	if err != nil {
		return err
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("registration unsuccessful: %v: %v", resp.StatusCode, resp.Status)
	}

	return nil
}
