package main

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

type (
	HTTPClient interface {
		Do(req *http.Request) (*http.Response, error)
	}
	MockClient struct {
		DoFunc func(req *http.Request) (*http.Response, error)
	}
)

var (
	Client    HTTPClient
	GetDoFunc func(req *http.Request) (*http.Response, error)
)

func Post(url string, body interface{}, headers http.Header) (*http.Response, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}
	request.Header = headers
	return Client.Do(request)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}

func pay() (*http.Response, error) {
	Client = &MockClient{}
	time.Sleep(2000 * time.Millisecond)
	rand.Seed(time.Now().UnixNano())
	// TODO Revert timeout to 2 secs
	success := rand.Intn(2)
	code := http.StatusOK
	if success == 0 {
		code = http.StatusGatewayTimeout
	}
	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: code,
		}, nil
	}
	return Post("http://pay.dummmy.com", nil, nil)
}



