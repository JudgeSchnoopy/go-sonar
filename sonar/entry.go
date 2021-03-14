package sonar

import (
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

// NewEntry generates a new entry object
func NewEntry(name, address string) Entry {
	return Entry{
		Name:    name,
		Address: address,
	}
}

// Checkin queries the monitored server and records it's new status
func (entry *Entry) Checkin() {
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
