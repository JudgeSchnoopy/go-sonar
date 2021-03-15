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
		registry     Registry
		entry        Entry
		wantRegistry Registry
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
			wantRegistry: Registry{
				Servers: map[string]RegEntry{
					"server01": &Entry{
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
			registry: Registry{
				Servers: map[string]RegEntry{
					"server02": &Entry{
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
			wantRegistry: Registry{
				Servers: map[string]RegEntry{
					"server02": &Entry{
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
			name:     "failedValidationTest",
			registry: NewRegistry(),
			entry: Entry{
				Name:    "server03",
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
		err := v.registry.Register(&v.entry)
		if err != nil && !v.wantErr {
			t.Errorf("register test %v failed with err %v", v.name, err)
		}

		for _, server := range v.registry.Servers {
			entry := server.Get()
			entry.LastCheck = time.Time{}

			entry.caller = nil

			compare := v.wantRegistry.Servers[entry.Name]

			if !reflect.DeepEqual(entry, compare) {
				t.Errorf("register test failed, want %+v, got %+v", compare, entry)
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
				Servers: map[string]RegEntry{
					"mock001": &Entry{
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
				Servers: map[string]RegEntry{
					"mock001": &Entry{
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
		values := *v.registry.Servers[v.name].Get()
		wantedValues := *v.wantRegistry.Servers[v.name].Get()

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