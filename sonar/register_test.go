package sonar

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	type registerTest struct {
		name         string
		registry     *Registry
		entry        Entry
		wantRegistry *Registry
		wantErr      bool
	}

	tests := []registerTest{
		{
			name:     "goodTest",
			registry: NewRegistry(),
			entry: Entry{
				Name:    "server01",
				Address: ":8080",
				caller: mockCaller{
					statusCode: 200,
					body:       "ok",
				},
			},
			wantRegistry: &Registry{
				Servers: map[string]Entry{
					"server01": {
						Name:       "server01",
						Address:    ":8080",
						StatusCode: 200,
						Healthy:    true,
						Status:     `"ok"`,
					},
				},
				lock: &sync.Mutex{},
			},
			wantErr: false,
		},
		{
			name: "duplicateTest",
			registry: &Registry{
				Servers: map[string]Entry{
					"server02": {
						Name:       "server02",
						Address:    ":8080",
						StatusCode: 200,
						Healthy:    true,
						Status:     `"ok"`,
					},
				},
			},
			entry: Entry{
				Name:    "server02",
				Address: ":8080",
				caller: mockCaller{
					statusCode: 200,
					body:       "ok",
				},
			},
			wantRegistry: &Registry{
				Servers: map[string]Entry{
					"server02": {
						Name:       "server02",
						Address:    ":8080",
						StatusCode: 200,
						Healthy:    true,
						Status:     `"ok"`,
					},
				},
				lock: &sync.Mutex{},
			},
			wantErr: true,
		},
		{
			name: "updateTest",
			registry: &Registry{
				Servers: map[string]Entry{
					"server03": {
						Name:       "server03",
						Address:    ":8080",
						StatusCode: 200,
						Healthy:    true,
						Status:     `"ok"`,
						caller: mockCaller{
							statusCode: 200,
							body:       "ok",
						},
					},
				},
				lock: &sync.Mutex{},
			},
			entry: Entry{
				Name:    "server03",
				Address: ":8081",
				caller: mockCaller{
					statusCode: 200,
					body:       "ok",
				},
			},
			wantRegistry: &Registry{
				Servers: map[string]Entry{
					"server03": {
						Name:       "server03",
						Address:    ":8081",
						StatusCode: 200,
						Healthy:    true,
						Status:     `"ok"`,
					},
				},
				lock: &sync.Mutex{},
			},
			wantErr: true,
		},
		{
			name:     "failedValidationTest",
			registry: NewRegistry(),
			entry: Entry{
				Name:    "server04",
				Address: ":8080",
				caller: mockCaller{
					statusCode: 500,
					body:       nil,
				},
			},
			wantRegistry: NewRegistry(),
			wantErr:      true,
		},
	}

	for _, v := range tests {
		err := v.registry.Register(v.entry)
		if err != nil && !v.wantErr {
			t.Errorf("register test %v failed with err %v", v.name, err)
		}

		for _, server := range v.registry.Servers {
			server.LastCheck = time.Time{}

			server.caller = nil

			compare := v.wantRegistry.Servers[server.Name]

			if !reflect.DeepEqual(server, compare) {
				t.Errorf("register test failed, want %+v, got %+v", compare, server)
			}
		}
	}
}

func TestCheckAll(t *testing.T) {
	type checkAllTest struct {
		name         string
		registry     Registry
		wantRegistry Registry
	}

	tests := []checkAllTest{
		{
			name: "mock001",
			registry: Registry{
				Servers: map[string]Entry{
					"mock001": {
						Name:       "mock001",
						Address:    ":8080",
						Healthy:    true,
						StatusCode: 200,
						caller: mockCaller{
							statusCode: 500,
							body:       "bad request",
						},
					},
				},
				lock: &sync.Mutex{},
			},
			wantRegistry: Registry{
				Servers: map[string]Entry{
					"mock001": {
						Name:       "mock001",
						Address:    ":8080",
						Healthy:    false,
						StatusCode: 500,
						Status:     `"bad request"`,
					},
				},
				lock: &sync.Mutex{},
			},
		},
	}

	for _, v := range tests {
		v.registry.CheckAll()
		values := v.registry.Servers[v.name]
		wantedValues := v.wantRegistry.Servers[v.name]

		// We can't replicate the LastCheck time accurately
		// Instead we'll validate that it's been updated to a non-zero value
		// If validation passes, we can reset it to a zero value to validate the rest of the entry
		if values.LastCheck == wantedValues.LastCheck {
			t.Errorf("checkall test failed - LastCheck was not updated")
		} else {
			values.LastCheck = time.Time{}
			values.caller = nil
		}

		if !reflect.DeepEqual(values, wantedValues) {
			t.Errorf("checkall test failed, want %v, got %v", wantedValues, values)
		}
	}
}

func TestRemove(t *testing.T) {
	type removeTest struct {
		name         string
		registry     Registry
		entry        Entry
		wantRegistry Registry
	}

	tests := []removeTest{
		{
			name: "removeTest",
			registry: Registry{
				Servers: map[string]Entry{
					"server05": {
						Name:    "server05",
						Address: ":8080",
					},
				},
			},
			entry: Entry{
				Name:    "server05",
				Address: ":8080",
			},
			wantRegistry: Registry{
				Servers: map[string]Entry{},
			},
		},
		{
			name: "notFoundTest",
			registry: Registry{
				Servers: map[string]Entry{
					"server06": {
						Name:    "server06",
						Address: ":8080",
					},
				},
			},
			entry: Entry{
				Name:    "server05",
				Address: ":8080",
			},
			wantRegistry: Registry{
				Servers: map[string]Entry{
					"server06": {
						Name:    "server06",
						Address: ":8080",
					},
				},
			},
		},
	}

	for _, v := range tests {
		v.registry.Remove(v.entry)

		if !reflect.DeepEqual(v.wantRegistry, v.registry) {
			t.Errorf("remove test failed, want %+v, got %+v", v.wantRegistry, v.registry)
		}
	}
}
