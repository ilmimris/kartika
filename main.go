package main

import (
	"fmt"
	"io"
	"net/http"
	"plugin"

	"github.com/getkin/kin-openapi/openapi3"
)

// Define the interface that all plugins must implement
type Plugin interface {
	HandleRequest(*http.Request, map[string]string, map[string]string, []byte) (*http.Response, error)
}

// Load the OpenAPI Spec
func loadSpec(path string) (*openapi3.T, error) {
	loader := openapi3.NewLoader()
	spec, err := loader.LoadFromFile(path)
	if err != nil {
		return nil, err
	}
	return spec, nil
}

// Load the Go plugin
func loadPlugin(operationID string) (Plugin, error) {
	p, err := plugin.Open(operationID + ".so")
	if err != nil {
		return nil, err
	}
	symPlugin, err := p.Lookup("Plugin")
	if err != nil {
		return nil, err
	}
	var plugin Plugin
	plugin, ok := symPlugin.(Plugin)
	if !ok {
		return nil, fmt.Errorf("unexpected type from module symbol")
	}
	return plugin, nil
}

// Handle the HTTP request
func handleRequest(spec *openapi3.T, plugins map[string]Plugin) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the operationID from the request
		operationID := spec.Paths[r.URL.Path].Get.OperationID

		// Get the plugin
		plugin, ok := plugins[operationID]
		if !ok {
			http.Error(w, "Operation not supported", http.StatusNotFound)
			return
		}

		// Extract the path parameters and query parameters
		pathParams := map[string]string{
			r.URL.Path: r.URL.Path,
		}
		queryParams := make(map[string]string)
		for k, v := range r.URL.Query() {
			queryParams[k] = v[0]
		}

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Handle the request
		resp, err := plugin.HandleRequest(r, pathParams, queryParams, body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Read the response body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write the response
		w.WriteHeader(resp.StatusCode)
		w.Write(respBody)
	}
}

func main() {
	// Load the OpenAPI Spec
	spec, err := loadSpec("oas_echo.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Map to store operationID to plugin mapping
	plugins := make(map[string]Plugin)

	// Iterate over the paths
	for _, pathItem := range spec.Paths {
		// Iterate over the operations
		for _, operation := range pathItem.Operations() {
			// Get the operationID
			operationID := operation.OperationID

			// Load the Go plugin
			p, err := loadPlugin(operationID)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Add the plugin to the map
			plugins[operationID] = p
		}
	}

	// Start the server
	http.HandleFunc("/", handleRequest(spec, plugins))
	http.ListenAndServe(":8080", nil)
}
