package sonar

import (
	"fmt"
	"strings"
	"sync"
)

// Registry is a collection of monitored servers
type Registry struct {
	Servers map[string]RegEntry `json:"servers"`
	lock    *sync.Mutex
}

// NewRegistr generates a new registry
func NewRegistry() Registry {
	return Registry{
		Servers: make(map[string]RegEntry),
		lock:    &sync.Mutex{},
	}
}

// Register adds an entry to the registry
// This checks the register for similar entries and queries the service to ensure it responds
func (reg *Registry) Register(entry RegEntry) error {
	err := reg.checkRegistry(entry)
	if err != nil {
		return err
	}

	err = entry.validateEntry()
	if err != nil {
		return err
	}

	values := entry.Get()

	reg.lock.Lock()

	reg.Servers[values.Name] = entry

	reg.lock.Unlock()

	return nil
}

// checkRegistry determines whether the entry already exists in the registry
func (reg *Registry) checkRegistry(entry RegEntry) error {
	values := entry.Get()
	regEntry, ok := reg.Servers[values.Name]

	if ok {
		currentEntry := regEntry.Get()
		if strings.EqualFold(currentEntry.Address, values.Address) {
			return fmt.Errorf("entry for %v already exists and matches address %v", values.Name, values.Address)
		} else {
			fmt.Printf("entry for %v exists, updating address to %v\n", values.Name, values.Address)
			return nil
		}
	}

	return nil
}

// CheckAll loops through all registry entries and runs a check-in
func (reg *Registry) CheckAll() {
	fmt.Println("checking all services")

	for i, v := range reg.Servers {
		v.Checkin()

		reg.lock.Lock()

		reg.Servers[i] = v

		reg.lock.Unlock()

		fmt.Printf("%+v\n", reg.Servers[i])
	}
}
