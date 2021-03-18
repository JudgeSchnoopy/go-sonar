// +build integration

package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/JudgeSchnoopy/go-sonar/client"
)

func TestSonar(t *testing.T) {
	sonar, err := New(
		WithCustomPort(8081),
		WithCustomSchedule(5*time.Second),
	)
	if err != nil {
		t.Errorf("server.New() failed: %v", err)
	}

	go func() {
		err := sonar.Start()
		if err != nil {
			t.Errorf("server.Start() failed: %v", err)
		}
	}()

	clientServer1, err := New(
		WithCustomPort(8082),
		WithCustomSchedule(1*time.Hour),
	)
	if err != nil {
		t.Errorf("server.New() failed: %v", err)
	}

	go func() {
		err := clientServer1.Start()
		if err != nil {
			t.Errorf("server.Start() failed: %v", err)
		}
	}()

	// wait for the Sonar service to be up
	_, err = http.Get("http://localhost:8081/docs")
	for err != nil {
		_, err = http.Get("http://localhost:8081/docs")
	}

	// wait for the cliet1 service to be up
	_, err = http.Get("http://localhost:8082/docs")
	for err != nil {
		_, err = http.Get("http://localhost:8082/docs")
	}

	client1 := client.New("http://localhost:8081", "http://localhost:8082", "client1",
		client.WithSelfRegistration(),
	)
	client1.AddDependency("client1", "Sonar", "http://localhost:8080/registry", sonar.Registry)
	client1.StartDependencyChecks(1 * time.Second)
	time.Sleep(1 * time.Second)
	client1.Report()
	client1.StopDependdencyChecks()

	t.Errorf("Registry: %+v", sonar.Registry)
}
