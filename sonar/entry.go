package sonar

import (
	"fmt"
	"net/http"
	"time"
)

// Entry represents a monitored server in the registry
type Entry struct {
	Name       string      `json:"name"`
	Address    string      `json:"address"`
	LastCheck  time.Time   `json:"lastCheck"`
	Healthy    bool        `json:"healthy"`
	StatusCode int         `json:"statusCode"`
	Status     interface{} `json:"status"`
}

type RegEntry interface {
	Checkin()
	validateEntry() error
	Get() *Entry
}

// NewEntry generates a new entry object
func NewEntry(name, address string) Entry {
	return Entry{
		Name:    name,
		Address: address,
	}
}

func (entry Entry) Get() *Entry {
	return &entry
}

// Checkin queries the monitored server and records it's new status
func (entry Entry) Checkin() {
	entry.LastCheck = time.Now()

	response, err := http.Get(entry.Address)
	if err != nil {
		entry.Healthy = false
		entry.StatusCode = response.StatusCode
		entry.Status = err.Error()
		return
	}

	entry.StatusCode = response.StatusCode
	entry.Status = response.Body

	if entry.StatusCode > 299 {
		entry.Healthy = false
		return
	}

	entry.Healthy = true
}

func (entry Entry) validateEntry() error {
	fmt.Println("checking service")

	entry.Checkin()

	if !entry.Healthy {
		return fmt.Errorf("server %v did not respond at %v and will not be added", entry.Name, entry.Address)
	}

	fmt.Printf("no entry found for %v\n", entry.Name)

	return nil
}
