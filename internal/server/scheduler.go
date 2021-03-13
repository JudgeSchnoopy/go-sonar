package server

import (
	"fmt"
	"time"
)

func (server *Server) startScheduler(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		fmt.Println("starting scheduler")
		for {
			select {
			case <-server.scheduleStopper:
				fmt.Println("Stopping scheduler")
				ticker.Stop()
				return
			case <-ticker.C:
				server.Registry.CheckAll()
			}
		}
	}()
}
