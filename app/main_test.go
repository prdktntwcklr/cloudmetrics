package main

import (
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
