package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Report sends current status to Sonar
func (client *Client) Report() error {
	jsonResponse, err := json.Marshal(client.Response)
	if err != nil {
		return err
	}

	resp, err := http.Post(client.SonarAddress, "application/json", bytes.NewBuffer(jsonResponse))
	if err != nil {
		return err
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("registration unsuccessful: %v: %v", resp.StatusCode, resp.Status)
	}

	return nil
}
