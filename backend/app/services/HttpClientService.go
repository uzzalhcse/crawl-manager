package services

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type HttpClientService interface {
	DoRequest(method, url string, payload interface{}, headers map[string]string) ([]byte, error)
}

type HttpClientServiceImpl struct {
	client *http.Client
}

func NewHttpClientService() *HttpClientServiceImpl {
	return &HttpClientServiceImpl{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// DoRequest handles any HTTP request and returns a generic JSON response as a map
func (s *HttpClientServiceImpl) DoRequest(method, url string, payload interface{}, headers map[string]string) ([]byte, error) {
	var req *http.Request
	var err error

	// If payload is provided, marshal it into JSON for POST/PUT requests
	if payload != nil {
		body, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, err
	}

	// Add headers to the request
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Return response body as byte slice
	return io.ReadAll(resp.Body)
}
