package main

import (
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"

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

// TODO: refactor to avoid using global variables
var GitCommit = "unknown"
var templates = template.Must(template.ParseFiles("index.html"))

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"GitCommit": GitCommit,
	}

	err := templates.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"OK"}`))
}

func main() {
	go func() {
		currentTemp := 25.0
		log.Println("Starting sensor simulation loop...")

		for {
			UpdateTemperature(ambientTemp, &currentTemp)
			IncrementReadings(readingsTotal)
			
			time.Sleep(2 * time.Second)
		}
	}()

	log.Printf("Starting Cloudmetrics App Version: %s", GitCommit)

	mux := http.NewServeMux()
	mux.HandleFunc("/", IndexHandler)
	mux.HandleFunc("/healthz", HealthzHandler)
	mux.Handle("/metrics", promhttp.Handler())
	
	loggingMux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Consider excluding "/metrics" from logs to reduce noise in production
        log.Printf("Request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
        mux.ServeHTTP(w, r)
    })
	
	port := "8080"
	address := ":" + port

	log.Printf("Starting server on %s", address)
	log.Fatal(http.ListenAndServe(address, loggingMux))
}
