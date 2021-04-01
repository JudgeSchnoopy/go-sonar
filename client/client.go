// Package client provides a Sonar client for Go microservices to implement
package client

import (
	"fmt"
	"time"
)

// Client runs Sonar services locally
type Client struct {
	SonarAddress    string
	Response        Response
	scheduleStopper chan bool
}

// ClientOptions provide customizations to the client
type ClientOptions func(*Client)

// WithScheduler initializes the client Sonar scheduled checks
func WithScheduler(interval time.Duration) func(*Client) {
	return func(client *Client) {
		client.StartDependencyChecks(interval)
	}
}

// WithSelfRegistration checks in with the Sonar server on client initialization
func WithSelfRegistration() func(*Client) {
	return func(client *Client) {
		err := client.Report
		if err != nil {
			fmt.Printf("failed to register to Sonar")
			return
		}
		fmt.Println("Sonar registration successful")
	}
}

// New generaes a new Sonar client
func New(sonarAddress, selfAddress, serviceName string) Client {
	return Client{
		SonarAddress: sonarAddress,
		Response: Response{
			Name:         serviceName,
			Address:      selfAddress,
			Dependencies: make(map[string][]dependency),
		},
	}
}
