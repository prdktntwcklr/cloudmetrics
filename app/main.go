package main

import (
	"net/http"
	"time"
	"math/rand"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ambientTemp = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "cm_ambient_temp_celsius",
		Help: "Current ambient temperature of the simulated environment.",
	})
	readingsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "cm_readings_total",
		Help: "Total number of sensor readings processed.",
	})
)

func main() {
	go func() {
		currentTemp := 25.0
		for {
			// Simulate a Random Walk
			currentTemp += (rand.Float64() - 0.5)
			ambientTemp.Set(currentTemp)
			readingsTotal.Inc()
			
			time.Sleep(2 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())

	println("CloudMetrics server starting on :8080...")
	http.ListenAndServe(":8080", nil)
}
