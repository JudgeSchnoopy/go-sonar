// +build integration

package server

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/JudgeSchnoopy/go-sonar/client"
)

func TestSonar(t *testing.T) {
	sonarService, err := New(
		WithCustomPort(8081),
		WithCustomSchedule(5*time.Second),
	)
	if err != nil {
		t.Errorf("server.New() failed: %v", err)
	}

	go func() {
		err := sonarService.Start()
		if err != nil {
			t.Errorf("server.Start() failed: %v", err)
		}
	}()

	err = newTestClientServer(8082)
	if err != nil {
		t.Errorf("failed service1: %v", err)
	}
	err = newTestClientServer(8083)
	if err != nil {
		t.Errorf("failed service1: %v", err)
	}
	err = newTestClientServer(8084)
	if err != nil {
		t.Errorf("failed service1: %v", err)
	}

	// wait for the Sonar service to be up
	_, err = http.Get("http://localhost:8081/docs")
	for err != nil {
		_, err = http.Get("http://localhost:8081/docs")
	}

	// wait for the client1 service to be up
	_, err = http.Get("http://localhost:8082/docs")
	for err != nil {
		_, err = http.Get("http://localhost:8082/docs")
	}

	// wait for the client2 service to be up
	_, err = http.Get("http://localhost:8083/docs")
	for err != nil {
		_, err = http.Get("http://localhost:8083/docs")
	}

	// wait for the client3 service to be up
	_, err = http.Get("http://localhost:8084/docs")
	for err != nil {
		_, err = http.Get("http://localhost:8084/docs")
	}

	var docsReturn []string
	client1 := client.New("http://localhost:8081", "http://localhost:8082/docs", "client1",
		client.WithSelfRegistration(),
	)
	client1.AddDependency("client1", "Sonar", "http://localhost:8081/registry", sonarService.Registry)
	client1.AddDependency("client1", "Service1", "http://localhost:8082/docs", docsReturn)
	client1.AddDependency("client1", "Service2", "http://localhost:8083/docs", docsReturn)
	client1.AddDependency("client1", "Service3", "http://localhost:8084/docs", docsReturn)
	client1.StartDependencyChecks(1 * time.Second)
	time.Sleep(1 * time.Second)
	client1.Report()
	client1.StopDependdencyChecks()

	client1Entry, err := sonarService.Registry.Get("client1")
	if err != nil {
		t.Errorf("could not retrieve client1 from registry: %v", err)
	}

	if !client1Entry.Healthy {
		t.Errorf("client1 failed healthy checks: %v", client1Entry.Response.Dependencies)
	}
}

func newTestClientServer(port int) error {
	clientServer, err := New(
		WithCustomPort(port),
		WithCustomSchedule(1*time.Hour),
	)
	if err != nil {
		return err
	}

	var errorChan chan error

	go func(chan error) {
		err := clientServer.Start()
		if err != nil {
			fmt.Printf("service failure: %v\n", err)
			errorChan <- err
		}
	}(errorChan)

	if len(errorChan) > 0 {
		return <-errorChan
	}

	return nil
}
