package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	
	dto "github.com/prometheus/client_model/go"
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

func TestIngestReadingHandler_Success(t *testing.T) {
	app := &App{}

	sensorID := 5
	sensorStr := "5"
	expectedTemp := 24.5

	reading := Reading{
		SensorID:    sensorID,
		Temperature: expectedTemp,
	}
	jsonPayload, _ := json.Marshal(reading)

	req := httptest.NewRequest("POST", "/api/readings", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	app.IngestReadingHandler(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		t.Errorf("Expected status 202 Accepted, got %v", resp.Status)
	}

	var gaugeMetric dto.Metric
	err := ambientTempVec.WithLabelValues(sensorStr).Write(&gaugeMetric)
	if err != nil {
		t.Fatalf("Failed to read gauge metric: %v", err)
	}

	actualTemp := gaugeMetric.GetGauge().GetValue()
	if actualTemp != expectedTemp {
		t.Errorf("Expected gauge temperature to be %f, got %f", expectedTemp, actualTemp)
	}

	var counterMetric dto.Metric
	err = readingsTotalVec.WithLabelValues(sensorStr).Write(&counterMetric)
	if err != nil {
		t.Fatalf("Failed to read counter metric: %v", err)
	}

	actualCount := counterMetric.GetCounter().GetValue()
	if actualCount < 1 {
		t.Errorf("Expected counter to be incremented, got %f", actualCount)
	}
}

func TestIngestReadingHandler_BadRequest(t *testing.T) {
	app := &App{}

	req := httptest.NewRequest("POST", "/api/readings", bytes.NewBufferString("not-json"))
	w := httptest.NewRecorder()

	app.IngestReadingHandler(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400 Bad Request for bad JSON, got %v", resp.Status)
	}
}
