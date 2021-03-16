package client

import (
	"fmt"
	"time"
)

func (deps Dependencies) startDependencyChecks(interval time.Duration, stopper chan bool) {
	ticker := time.NewTicker(interval)
	go func() {
		fmt.Printf("checking dependencies every %v\n", interval)
		for {
			select {
			case <-stopper:
				fmt.Println("Stopping dependency checking")
				ticker.Stop()
				return
			case <-ticker.C:
				deps.CheckAll()
			}
		}
	}()
}
