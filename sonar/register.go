package sonar

import (
	"fmt"
	"strings"
	"sync"
)

type Registry struct {
	Servers map[string]string `json:"servers"`
	lock    *sync.Mutex
}

type Entry struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func NewRegistry() Registry {
	return Registry{
		Servers: make(map[string]string),
		lock:    &sync.Mutex{},
	}
}

func (reg *Registry) Register(entry Entry) {
	if !reg.checkRegistry(entry) {
		reg.lock.Lock()

		reg.Servers[entry.Name] = entry.Address

		reg.lock.Unlock()
	}
}

func (reg *Registry) checkRegistry(entry Entry) bool {
	currentEntry, ok := reg.Servers[entry.Name]
	if ok && strings.EqualFold(currentEntry, entry.Address) {
		fmt.Printf("entry for %v already exists and matches address %v\n", entry.Name, entry.Address)
		return true
	} else if ok {
		fmt.Printf("entry for %v exists, updating address to %v\n", entry.Name, entry.Address)
		return false
	}
	fmt.Printf("creating entry for %v at address %v", entry.Name, entry.Address)
	return false
}
