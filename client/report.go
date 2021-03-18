package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Report sends current status to Sonar
func (client *Client) Report() error {
	if !client.Registered {
		client.Register()
	}

	client.sonarPost(client.SonarAddress + "/report")

	return nil
}

func (client *Client) Register() error {
	_, err := client.sonarPost(client.SonarAddress + "/register")
	if err != nil {
		fmt.Printf("Sonar registration unsuccessful: %v\n", err.Error())
		return err
	}

	client.Registered = true
	return nil
}

func (client *Client) sonarPost(address string) (*http.Response, error) {
	jsonResponse, err := json.Marshal(client.Response)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(address, "application/json", bytes.NewBuffer(jsonResponse))
	if err != nil {
		return resp, err
	}

	if resp.StatusCode > 299 {
		return resp, fmt.Errorf("sonar call unsuccessful: %v: %v\n", resp.StatusCode, resp.Status)
	}

	return resp, err
}
