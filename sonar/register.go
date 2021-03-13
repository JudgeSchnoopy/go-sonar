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

func NewRegistry() Registry {
	return Registry{
		Servers: make(map[string]string),
	}
}

func (reg *Registry) Register(name, address string) {
	if !reg.checkRegistry(name, address) {
		reg.lock.Lock()

		reg.Servers[name] = address

		reg.lock.Unlock()
	}
}

func (reg *Registry) checkRegistry(name, address string) bool {
	entry, ok := reg.Servers[name]
	if ok && strings.EqualFold(entry, address) {
		fmt.Printf("entry for %v already exists and matches address %v\n", name, address)
		return true
	} else if ok {
		fmt.Printf("entry for %v exists, updating address to %v\n", name, address)
		return false
	}
	fmt.Printf("creating entry for %v at address %v", name, address)
	return false
}
