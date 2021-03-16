package client

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
)

type Dependency struct {
	Name      string      `json:"name"`
	Address   string      `json:"address"`
	Want      interface{} `json:"want"`
	Validated bool        `json:"validated"`
}

type Dependencies []Dependency

func (deps Dependencies) CheckAll() {
	for _, v := range deps {
		v.CheckDependency()
	}
}

func (dep *Dependency) CheckDependency() {
	resp, err := http.Get(dep.Address)
	if err != nil {
		dep.Validated = false
	}

	body, err := ioutil.ReadAll(io.Reader(resp.Body))
	if err != nil {
		dep.Validated = false
	}

	if err := resp.Body.Close(); err != nil {
		dep.Validated = false
	}

	var got interface{}

	if err := json.Unmarshal(body, &got); err != nil {
		dep.Validated = false
	}

	if !reflect.DeepEqual(got, dep.Want) {
		dep.Validated = false
	}
}
