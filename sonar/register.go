package sonar

import (
	"fmt"
	"strings"
	"sync"
)

// Registry is a collection of monitored servers
type Registry struct {
	Servers map[string]Entry `json:"servers"`
	lock    *sync.Mutex
}

// NewRegistr generates a new registry
func NewRegistry() Registry {
	return Registry{
		Servers: make(map[string]Entry),
		lock:    &sync.Mutex{},
	}
}

// Register adds an entry to the registry
// This checks the register for similar entries and queries the service to ensure it responds
func (reg *Registry) Register(entry Entry) error {
	found := reg.checkRegistry(entry)
	if found {
		return fmt.Errorf("entry for %v already exists and matches address %v", entry.Name, entry.Address)
	}

	err := entry.validateEntry()
	if err != nil {
		return err
	}

	reg.lock.Lock()

	reg.Servers[entry.Name] = entry

	reg.lock.Unlock()

	return nil
}

// checkRegistry determines whether the entry already exists in the registry
func (reg *Registry) checkRegistry(entry Entry) bool {
	regEntry, ok := reg.Servers[entry.Name]

	if ok {
		if strings.EqualFold(regEntry.Address, entry.Address) {
			return true
		} else {
			fmt.Printf("entry %v exists but doesn't match address %v\n", entry.Name, entry.Address)
		}
	}

	return false
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

func (reg *Registry) Remove(entry Entry) error {
	found := reg.checkRegistry(entry)
	if !found {
		return fmt.Errorf("entry %v not found - cannot remove", entry.Name)
	}

	delete(reg.Servers, entry.Name)
	return nil
}
