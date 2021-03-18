package sonar

import (
	"reflect"
	"testing"

	"github.com/JudgeSchnoopy/go-sonar/client"
)

func TestNewEntry(t *testing.T) {
	type newEntryTest struct {
		response client.Response
		want     Entry
	}
	tests := []newEntryTest{
		{
			response: client.Response{
				Name:    "test001",
				Address: "http://fake/api/v1/test",
			},
			want: Entry{
				Name:    "test001",
				Address: "http://fake/api/v1/test",
				caller:  httpCaller{},
				Response: client.Response{
					Name:    "test001",
					Address: "http://fake/api/v1/test",
				},
			},
		},
	}
	for _, v := range tests {
		result := NewEntry(v.response)
		// reset the caller for test comparison
		result.caller = httpCaller{}
		if !reflect.DeepEqual(result, v.want) {
			t.Errorf("newEntryTest failed, want %+v, got %+v", v.want, result)
		}
	}
}
