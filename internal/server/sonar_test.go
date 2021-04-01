package server

import (
	"fmt"
	"net/http"
	"strconv"
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
	defer client1.Unregister()
	client1.AddDependency("client1", "Sonar", "http://localhost:8081/registry", sonarService.Registry)
	client1.AddDependency("client1", "Service1", "http://localhost:8082/docs", docsReturn)
	client1.AddDependency("client1", "Service2", "http://localhost:8083/docs", docsReturn)
	client1.AddDependency("client1", "Service3", "http://localhost:8084/docs", docsReturn)
	client1.StartDependencyChecks(1 * time.Second)
	time.Sleep(2 * time.Second)
	fmt.Printf("client1 is sending report: %+v", client1.Response)
	client1.Report()
	client1.StopDependdencyChecks()

	if len(client1.Response.Dependencies) != 4 {
		t.Errorf("failed adding dependencies to client: wanted 4, got %v.  %+v", len(client1.Response.Dependencies), client1.Response.Dependencies)
	}

	client1Entry, err := sonarService.Registry.Get("client1")
	if err != nil {
		t.Errorf("could not retrieve client1 from registry: %v", err)
	}

	if !client1Entry.Healthy {
		t.Errorf("client1 failed healthy checks: %v", client1Entry.Response.Dependencies)
	}

	if len(client1Entry.Response.Dependencies) != 4 {
		t.Errorf("failed full client registration: wanted 4 dependencies, got %v", len(client1Entry.Response.Dependencies))

		fmt.Printf("total entry is: %+v", client1Entry)
	}

	for _, v := range client1Entry.Response.Dependencies {
		for _, d := range v {
			if !d.Validated {
				t.Errorf("Dependency failed health checks, but client didn't report it: %v, address: %v", d.Name, d.Address)
			}
		}
	}
}

func newTestClientServer(port int) error {
	clientServer, err := New(
		WithCustomPort(port),
		WithCustomSchedule(1*time.Hour),
		WithSonarClient("http://localhost:8081", "http://localhost:"+strconv.Itoa(port), strconv.Itoa(port),
			client.WithSelfRegistration(),
		),
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
