package sonar

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/JudgeSchnoopy/go-sonar/client"
)

// Entry represents a monitored server in the registry
type Entry struct {
	Name         string          `json:"name"`
	Address      string          `json:"address"`
	LastCheck    time.Time       `json:"lastCheck"`
	Healthy      bool            `json:"healthy"`
	StatusCode   int             `json:"statusCode"`
	Status       interface{}     `json:"status"`
	Dependencies client.Response `json:"dependencies"`
	caller       caller
}

// NewEntry generates a new entry object
func NewEntry(response client.Response) Entry {
	entry := Entry{
		Name:       response.Name,
		Address:    response.Address,
		LastCheck:  time.Time{},
		Healthy:    response.Healthy,
		StatusCode: 0,
		Status:     nil,
		caller: httpCaller{
			client: http.DefaultClient,
		},
		Dependencies: response,
	}

	return entry
}

// Checkin queries the monitored server and records it's new status
func (entry *Entry) Checkin() {
	entry.LastCheck = time.Now()

	response, err := entry.caller.call(entry)

	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	entry.Status = buf.String()

	if err != nil || response.StatusCode > 299 {
		entry.Healthy = false
		entry.StatusCode = response.StatusCode
		return
	}

	entry.StatusCode = response.StatusCode

	entry.Healthy = true
}

func (entry *Entry) validateEntry() error {
	fmt.Printf("checking service %v\n", entry.Name)

	entry.Checkin()

	if !entry.Healthy {
		return fmt.Errorf("server %v did not respond at %v and will not be added", entry.Name, entry.Address)
	}

	fmt.Printf("%v passed validation\n", entry.Name)

	return nil
}
