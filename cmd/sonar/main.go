package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/JudgeSchnoopy/go-sonar/internal/server"
)

func main() {
	sonar, err := server.New(
		server.WithCustomSchedule(time.Second*5),
		server.WithDebugEndpoints(),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("starting sonar")

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := sonar.Start(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	var wait time.Duration
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	sonar.Stop(ctx)
	// Optionally, you could run sonar.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
