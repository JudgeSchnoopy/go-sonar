package server

import (
	"fmt"
	"time"
)

func (server *Server) startScheduler(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		fmt.Printf("starting scheduler for every %v\n", interval)
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
