package sonar

import (
	"reflect"
	"testing"
)

func TestNewEntry(t *testing.T) {
	type newEntryTest struct {
		name    string
		address string
		want    Entry
	}
	tests := []newEntryTest{
		{
			name:    "test001",
			address: "http://fake/api/v1/test",
			want: Entry{
				Name:    "test001",
				Address: "http://fake/api/v1/test",
				caller:  httpCaller{},
			},
		},
	}
	for _, v := range tests {
		result := NewEntry(v.name, v.address)
		if !reflect.DeepEqual(result, v.want) {
			t.Errorf("newEntryTest failed, want %+v, got %+v", v.want, result)
		}
	}
}
