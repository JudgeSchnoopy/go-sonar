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

func (client *Client) Unregister() error {
	httpClient := &http.Client{}

	jsonResponse, err := json.Marshal(client.Response)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, client.SonarAddress+"/register", bytes.NewBuffer(jsonResponse))
	if err != nil {
		return err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("sonar unregistration unsuccessful: %v: %v", resp.StatusCode, resp.Status)
	}

	client.Registered = false
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
		return resp, fmt.Errorf("sonar call unsuccessful: %v: %v", resp.StatusCode, resp.Status)
	}

	return resp, err
}
