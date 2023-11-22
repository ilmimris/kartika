package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type EchoPlugin struct{}

type Response struct {
	Data       any `json:"data"`
	StatusCode int `json:"statusCode"`
	Latency    int `json:"latency"`
}

// HandleRequest handles the incoming HTTP request and generates a response.
// It takes the request object, path parameters, query parameters, and request body as input.
// It returns the generated HTTP response and an error, if any.
func (p *EchoPlugin) HandleRequest(r *http.Request, pathParams map[string]string, queryParams map[string]string, body []byte) (*http.Response, error) {
	// Get request parameters
	_ = pathParams
	params := queryParams
	_ = body

	// Get the query parameter "q"
	q := params["q"]

	// Generate a response body
	resp := Response{
		Data:       q,
		StatusCode: http.StatusOK,
		Latency:    0,
	}

	// Marshal the response body
	// This is a placeholder. Replace with your own logic.
	res, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	// Wrap the response body in a bytes.Reader
	respBody := bytes.NewReader(res)

	// Handle the request and generate a response
	// This is a placeholder. Replace with your own logic.
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(respBody),
	}, nil
}

// Export a variable named Plugin that implements the Plugin interface
var Plugin EchoPlugin
