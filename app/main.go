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

type TemperatureSetter interface {
	Set(float64)
}

func UpdateTemperature(setter TemperatureSetter, currentTemp *float64) {
	delta := (rand.Float64() - 0.5)
	*currentTemp += delta
	setter.Set(*currentTemp)
}

type CounterIncrementer interface {
	Inc()
}

func IncrementReadings(counter CounterIncrementer) {
	counter.Inc()
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {
	go func() {
		currentTemp := 25.0

		for {
			UpdateTemperature(ambientTemp, &currentTemp)
			IncrementReadings(readingsTotal)
			
			time.Sleep(2 * time.Second)
		}
	}()

	http.HandleFunc("/", IndexHandler)
	http.Handle("/metrics", promhttp.Handler())

	println("CloudMetrics server starting on :8080...")
	http.ListenAndServe(":8080", nil)
}
