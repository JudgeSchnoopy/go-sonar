package client

import (
	"fmt"
	"time"
)

// StartDependencyChecks begins a scheduler to check in with Sonar
func (client *Client) StartDependencyChecks(interval time.Duration) {
	client.scheduleStopper = make(chan bool)
	ticker := time.NewTicker(interval)
	go func() {
		fmt.Printf("checking dependencies every %v\n", interval)
		for {
			select {
			case <-client.scheduleStopper:
				fmt.Println("Stopping dependency checking")
				ticker.Stop()
				return
			case <-ticker.C:
				client.checkAllDependencies(false)
			}
		}
	}()
}

func (client *Client) StopDependdencyChecks() {
	client.scheduleStopper <- true
}
