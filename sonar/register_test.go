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
		entry        mockEntry
		wantRegistry Registry
		wantErr      bool
	}

	tests := []registerTest{
		{
			name:     "goodTest",
			registry: NewRegistry(),
			entry: mockEntry{
				entry: &Entry{
					Name:    "server01",
					Address: ":8080",
				},
				newStatus:       "ok",
				newStatusCode:   200,
				newHealthStatus: true,
				validateFailure: false,
			},
			wantRegistry: Registry{
				Servers: map[string]RegEntry{
					"server01": Entry{
						Name:       "server01",
						Address:    ":8080",
						StatusCode: 200,
						Healthy:    true,
						Status:     "ok",
					},
				},
				lock: &sync.Mutex{},
			},
			wantErr: false,
		},
	}

	for _, v := range tests {
		err := v.registry.Register(v.entry)
		if err != nil && !v.wantErr {
			t.Errorf("register test failed with err %v", err)
		}

		for _, server := range v.registry.Servers {
			entry := server.Get()
			entry.LastCheck = time.Time{}
			compare := v.wantRegistry.Servers[entry.Name]

			if !reflect.DeepEqual(*entry, compare) {
				t.Errorf("register test failed, want %+v, got %+v", compare, *entry)
			}
		}
	}
}

func TestCheckRegistry(t *testing.T) {
	type checkRegistryTest struct {
		name     string
		registry Registry
		entry    Entry
		wantErr  bool
	}

	tests := []checkRegistryTest{
		{
			name: "Existing entry test",
			registry: Registry{
				Servers: map[string]RegEntry{
					"test001": Entry{
						Name:    "test001",
						Address: ":8080",
					},
				},
				lock: &sync.Mutex{},
			},
			entry: Entry{
				Name:    "test001",
				Address: ":8080",
			},
			wantErr: true,
		},
		{
			name: "New entry test",
			registry: Registry{
				Servers: map[string]RegEntry{
					"test001": Entry{
						Name:    "test001",
						Address: ":8080",
					},
				},
				lock: &sync.Mutex{},
			},
			entry: Entry{
				Name:    "test002",
				Address: ":8080",
			},
			wantErr: false,
		},
		{
			name: "Update address test",
			registry: Registry{
				Servers: map[string]RegEntry{
					"test001": Entry{
						Name:    "test001",
						Address: ":8080",
					},
				},
				lock: &sync.Mutex{},
			},
			entry: Entry{
				Name:    "test001",
				Address: ":8081",
			},
			wantErr: false,
		},
	}

	for _, v := range tests {
		err := v.registry.checkRegistry(v.entry)
		if err != nil && !v.wantErr {
			t.Errorf("checkRegistry test failed with err %v", err)
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
					"mock001": mockEntry{
						entry: &Entry{
							Name:       "mock001",
							Address:    ":8080",
							Healthy:    true,
							StatusCode: 200,
						},
						newStatus:       "bad request",
						newStatusCode:   500,
						newHealthStatus: false,
					},
				},
				lock: &sync.Mutex{},
			},
			wantRegistry: Registry{
				Servers: map[string]RegEntry{
					"mock001": mockEntry{
						entry: &Entry{
							Name:       "mock001",
							Address:    ":8080",
							Healthy:    false,
							StatusCode: 500,
							Status:     "bad request",
						},
						newStatus:       "bad request",
						newStatusCode:   500,
						newHealthStatus: false,
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
		}

		if !reflect.DeepEqual(values, wantedValues) {
			t.Errorf("checkall test failed, want %v, got %v", wantedValues, values)
		}
	}
}
