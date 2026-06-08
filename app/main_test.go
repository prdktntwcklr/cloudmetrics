package main

import (
	"net/http"
	"net/http/httptest"
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
