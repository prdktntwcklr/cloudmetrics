package main

import (
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func TestIncrementReadings(t *testing.T) {
	counter := prometheus.NewCounter(prometheus.CounterOpts{Name: "test_counter"})
	
	IncrementReadings(counter)

	var metric dto.Metric
	counter.Write(&metric)
	if metric.GetCounter().GetValue() != 1 {
		t.Errorf("Expected 1, got %v", metric.GetCounter().GetValue())
	}
}

func TestUpdateTemperature(t *testing.T) {
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{Name: "test_gauge"})
	temp := 25.0

	UpdateTemperature(gauge, &temp)

	var metric dto.Metric
	gauge.Write(&metric)
	if metric.GetGauge().GetValue() == 25.0 {
		t.Errorf("Expected temperature to change, got %f", metric.GetGauge().GetValue())
	}
}

func TestHealthzHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthzHandler)
	handler.ServeHTTP(rr, req)

	expectedStatus := http.StatusOK
	if rr.Code != expectedStatus {
		t.Errorf("Expected status %v, got %v",
			expectedStatus, rr.Code)
	}

	expectedBody := `{"status":"OK"}`
	if rr.Body.String() != expectedBody {
    t.Errorf("Expected body %s, got %s",
        expectedBody, rr.Body.String())
	}
}

func TestIndexHandler(t *testing.T) {
	// Save original global state and schedule restoration
	oldGitCommit := GitCommit
	oldTemplates := templates

	t.Cleanup(func() {
		GitCommit = oldGitCommit
		templates = oldTemplates
	})
	
	// Override the global GitCommit variable for the test environment
	GitCommit = "test-sha-12345"
	expectedString := "test-sha-12345"

	// Initialize the global template cache just for the test
	// (Normally main() handles this, but go test bypasses main())
	var err error
	templates, err = template.ParseFiles("index.html")
	if err != nil {
		t.Fatalf("Failed to parse index.html template: %v", err)
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	IndexHandler(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	bodyString := string(body)
	if !strings.Contains(bodyString, expectedString) {
		t.Errorf("Expected HTML body to contain '%v', but received:\n%s", expectedString, bodyString)
	}
}
