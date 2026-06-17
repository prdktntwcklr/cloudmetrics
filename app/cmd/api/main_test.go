package main

import (
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	
	"github.com/prometheus/client_golang/prometheus/testutil"
)

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
	t.Parallel()

	expectedString := "test-sha-12345"

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		t.Fatalf("Failed to parse index.html template: %v", err)
	}

	app := &App{
		GitCommit: expectedString,
		Templates: tmpl,
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	app.IndexHandler(w, req)
	
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

func TestNewSensor(t *testing.T) {
	expectedId := "sensor-42"

	sensor := NewSensor(42)
	if sensor.Id != expectedId {
		t.Errorf("Expected sensor id %s, got %s",
			expectedId, sensor.Id)
		}
}

func TestSensorRun(t *testing.T) {
	sensorId := "sensor-42"
	sensor := NewSensor(42)
	go sensor.Run()

	time.Sleep(50 * time.Millisecond)

	counterCount := testutil.ToFloat64(readingsTotalVec.WithLabelValues(sensorId))
	if counterCount < 1 {
		t.Errorf("Expected readings counter to be >= 1, got %v", counterCount)
	}

	currentTemp := testutil.ToFloat64(ambientTempVec.WithLabelValues(sensorId))
	if currentTemp < 24.0 || currentTemp > 26.0 {
		t.Errorf("Expected initial temperature to be roughly between 24 and 26, got %v", currentTemp)
	}
}
