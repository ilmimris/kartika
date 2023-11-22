package main

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestHandleRequest(t *testing.T) {
	// Create a new instance of the EchoPlugin
	p := &EchoPlugin{}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/path", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set query parameters
	q := "test"
	qParams := req.URL.Query()
	qParams.Add("q", q)
	req.URL.RawQuery = qParams.Encode()

	// Create empty pathParams and body
	pathParams := make(map[string]string)
	// Convert url.Values to map[string]string
	queryParams := make(map[string]string)
	for key, values := range qParams {
		if len(values) > 0 {
			queryParams[key] = values[0]
		}
	}

	// Call the HandleRequest function
	resp, err := p.HandleRequest(req, pathParams, queryParams, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Unmarshal the response body
	var respData Response
	err = json.Unmarshal(respBody, &respData)
	if err != nil {
		t.Fatal(err)
	}

	// Check the response data
	if respData.Data != q {
		t.Errorf("Expected response data %s, got %s", q, respData.Data)
	}
}
