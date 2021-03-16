// Package client provides a Sonar client for Go microservices to implement
package client

import (
	"fmt"
	"time"
)

type Client struct {
	SonarAddress    string
	Response        Response
	scheduleStopper chan bool
}

type ClientOptions func(*Client)

func WithScheduler(interval time.Duration) func(*Client) {
	return func(client *Client) {
		client.StartDependencyChecks(interval)
	}
}

func WithSelfRegistration() func(*Client) {
	return func(client *Client) {
		err := client.Report
		if err != nil {
			fmt.Printf("failed to register to Sonar: %v\n", err)
		}
		fmt.Println("Sonar registration successful")
	}
}

//
func NewClient(sonarAddress, selfAddress, serviceName string) Client {
	return Client{
		SonarAddress: sonarAddress,
		Response: Response{
			Name:         serviceName,
			Address:      selfAddress,
			Dependencies: make(map[string][]dependency),
		},
	}
}
