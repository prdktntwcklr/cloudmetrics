package main

import (
	"html/template"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
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

type App struct {
	GitCommit string
	Templates *template.Template
}

// Injected from the Dockerfile
var Version = "unknown"

func (app *App) IndexHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"GitCommit": app.GitCommit,
	}

	err := app.Templates.Execute(w, data)
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
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug, 
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	app := &App{
        GitCommit: Version, 
        Templates: template.Must(template.ParseFiles("index.html")),
    }
	
	go func() {
		currentTemp := 25.0
		slog.Info("Starting sensor simulation loop...")

		for {
			UpdateTemperature(ambientTemp, &currentTemp)
			IncrementReadings(readingsTotal)
			
			time.Sleep(2 * time.Second)
		}
	}()

	slog.Info("Starting Cloudmetrics App", "version", app.GitCommit)

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.IndexHandler)
	mux.HandleFunc("/healthz", HealthzHandler)
	mux.Handle("/metrics", promhttp.Handler())
	
	loggingMux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Consider excluding "/metrics" from logs to reduce noise in production
        slog.Debug("HTTP request received", 
			"method", r.Method, 
			"path", r.URL.Path, 
			"remote_addr", r.RemoteAddr,
		)
        mux.ServeHTTP(w, r)
    })
	
	port := "8080"
	address := ":" + port

	slog.Info("Starting server", "address", address)
	if err := http.ListenAndServe(address, loggingMux); err != nil {
		slog.Error("Server crashed", "error", err)
		os.Exit(1)
	}
}
