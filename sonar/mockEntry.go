package sonar

import (
	"fmt"
	"time"
)

type mockEntry struct {
	entry           *Entry
	newStatus       string
	newStatusCode   int
	newHealthStatus bool
	validateFailure bool
}

func (m mockEntry) Get() *Entry {
	return m.entry
}

func (m mockEntry) validateEntry() error {
	if m.validateFailure {
		return fmt.Errorf("validateEntry returned failure")
	}
	m.Checkin()

	return nil
}

func (m mockEntry) Checkin() {
	m.entry.LastCheck = time.Now()
	m.entry.Status = m.newStatus
	m.entry.Healthy = m.newHealthStatus
	m.entry.StatusCode = m.newStatusCode
}
