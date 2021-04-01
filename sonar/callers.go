package sonar

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
)

type caller interface {
	call(*Entry) (*http.Response, error)
}

type httpCaller struct {
	client *http.Client
}

func (caller httpCaller) call(entry *Entry) (*http.Response, error) {
	if caller.client == nil {
		caller.client = &http.Client{}
	}

	response, err := caller.client.Get(entry.Address)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type mockCaller struct {
	body       interface{}
	statusCode int
}

func (caller mockCaller) call(entry *Entry) (*http.Response, error) {
	mockResponse, err := json.Marshal(caller.body)
	if err != nil {
		return nil, err
	}

	// generate a test server so we can capture and inspect the request
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(caller.statusCode)
		res.Write(mockResponse)
	}))
	defer func() { testServer.Close() }()

	req, err := http.NewRequest(http.MethodGet, testServer.URL, nil)
	if err != nil {
		fmt.Printf("failed to make new request: %v", err)
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("failed to execute new request: %v", err)
		return nil, err
	}

	return res, nil
}
