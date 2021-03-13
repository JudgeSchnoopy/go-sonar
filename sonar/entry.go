package sonar

import (
	"net/http"
	"time"
)

type Entry struct {
	Name       string      `json:"name"`
	Address    string      `json:"address"`
	LastCheck  time.Time   `json:"lastCheck"`
	Healthy    bool        `json:"healthy"`
	StatusCode int         `json:"statusCode"`
	Status     interface{} `json:"status"`
}

func (entry *Entry) Checkin() {
	entry.LastCheck = time.Now()
	response, err := http.Get(entry.Address)
	if err != nil {
		entry.Healthy = false
		entry.StatusCode = response.StatusCode
		entry.Status = err.Error()
		return
	}

	entry.Healthy = true
	entry.StatusCode = response.StatusCode
	entry.Status = response.Body
}
