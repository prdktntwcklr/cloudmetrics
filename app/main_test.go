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
