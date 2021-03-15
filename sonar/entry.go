package sonar

import (
	"bytes"
	"fmt"
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
	caller     caller
}

type RegEntry interface {
	Checkin()
	validateEntry() error
	Get() *Entry
}

// NewEntry generates a new entry object
func NewEntry(name, address string) Entry {
	entry := Entry{
		Name:       name,
		Address:    address,
		LastCheck:  time.Time{},
		Healthy:    false,
		StatusCode: 0,
		Status:     nil,
		caller:     httpCaller{},
	}

	return entry
}

func (entry *Entry) Get() *Entry {
	return entry
}

// Checkin queries the monitored server and records it's new status
func (entry *Entry) Checkin() {
	entry.LastCheck = time.Now()

	response, err := entry.caller.call(entry)
	if err != nil {
		entry.Healthy = false
		entry.StatusCode = response.StatusCode
		entry.Status = err.Error()
		return
	}

	entry.StatusCode = response.StatusCode

	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	entry.Status = buf.String()

	if entry.StatusCode > 299 {
		entry.Healthy = false
		return
	}

	entry.Healthy = true
}

func (entry *Entry) validateEntry() error {
	fmt.Println("checking service")

	entry.Checkin()

	if !entry.Healthy {
		return fmt.Errorf("server %v did not respond at %v and will not be added", entry.Name, entry.Address)
	}

	fmt.Printf("no entry found for %v\n", entry.Name)

	return nil
}
