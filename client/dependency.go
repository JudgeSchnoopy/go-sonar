package client

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
)

// cependency details one endpoint dependency and current validation status
type dependency struct {
	Name      string      `json:"name"`
	Address   string      `json:"address"`
	Want      interface{} `json:"want"`
	Validated bool        `json:"validated"`
}

// Dependencies is a list of Dependency
type dependencies map[string][]dependency

// CheckAll runs through each Dependency and runs the CheckDependency method
func (client *Client) checkAllDependencies(noReport bool) {
	healthy := true
	for _, v := range client.Response.Dependencies {
		for _, d := range v {
			d.checkDependency()
			if !d.Validated {
				healthy = false
			}
		}
	}

	if !healthy {
		client.Response.Healthy = false
		err := client.Report()
		if err != nil {
			fmt.Printf("failed to report to Sonar: %v\n", err)
		}
	}
}

// AddDependency creates a new Dependency and adds it to the Dependencies map
func (client *Client) AddDependency(service, name, address string, want interface{}) {
	newDep := dependency{
		Name:    name,
		Address: address,
		Want:    want,
	}

	currentMap := client.Response.Dependencies

	list, ok := currentMap[service]
	if ok {
		list = append(list, newDep)

		currentMap[service] = list
	} else {
		currentMap[service] = []dependency{newDep}
	}

	client.Response.Dependencies = currentMap
}

// checkDependency executes the address and compares results to expectations
func (dep *dependency) checkDependency() {
	resp, err := http.Get(dep.Address)
	if err != nil {
		fmt.Errorf("failed to validate dependency %v: %v", dep.Name, err)
		dep.Validated = false
	}

	var got interface{}

	if resp != nil && resp.Body != nil {
		body, err := ioutil.ReadAll(io.Reader(resp.Body))
		if err != nil {
			dep.Validated = false
		}

		if err := resp.Body.Close(); err != nil {
			dep.Validated = false
		}

		if err := json.Unmarshal(body, &got); err != nil {
			dep.Validated = false
		}
	} else {
		got = nil
	}

	if !reflect.DeepEqual(got, dep.Want) {
		dep.Validated = false
	}

	dep.Validated = true
}
