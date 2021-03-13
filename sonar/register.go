package sonar

import (
	"fmt"
	"strings"
	"sync"
)

type Registry struct {
	Servers map[string]Entry `json:"servers"`
	lock    *sync.Mutex
}

func NewRegistry() Registry {
	return Registry{
		Servers: make(map[string]Entry),
		lock:    &sync.Mutex{},
	}
}

func (reg *Registry) Register(entry Entry) error {
	err := reg.checkRegistry(entry)
	if err != nil {
		return err
	}

	reg.lock.Lock()

	reg.Servers[entry.Name] = entry

	reg.lock.Unlock()

	return nil
}

func (reg *Registry) checkRegistry(entry Entry) error {
	currentEntry, ok := reg.Servers[entry.Name]

	if ok && strings.EqualFold(currentEntry.Address, entry.Address) {
		return fmt.Errorf("entry for %v already exists and matches address %v", entry.Name, entry.Address)
	} else if ok {
		fmt.Printf("entry for %v exists, updating address to %v\n", entry.Name, entry.Address)
		return nil
	}

	fmt.Println("checking service")

	entry.Checkin()
	if !entry.Healthy {
		return fmt.Errorf("server %v did not respond at %v and will not be added", entry.Name, entry.Address)
	}

	fmt.Printf("no entry found for %v\n", entry.Name)

	return nil
}

func (reg *Registry) CheckAll() {
	fmt.Println("checking all services")
	for _, v := range reg.Servers {
		v.Checkin()
		fmt.Printf("%+v\n", v)
	}
}
