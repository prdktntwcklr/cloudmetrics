// +build integration

package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthzIntegration(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(HealthzHandler))
    defer server.Close()

    resp, err := http.Get(server.URL + "/healthz")
    if err != nil {
        t.Fatalf("Failed to make request: %v", err)
    }
    defer resp.Body.Close()

	expectedStatus := http.StatusOK
    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status %v, got %v", expectedStatus, resp.StatusCode)
    }

    body, _ := io.ReadAll(resp.Body)
    expectedBody := `{"status":"OK"}`
    if string(body) != expectedBody {
        t.Errorf("Expected body %s, got %s", expectedBody, string(body))
    }
}
