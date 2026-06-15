package main

import (
	"fmt"
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
	ambientTempVec = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cm_ambient_temp_celsius",
		Help: "Current ambient temperature of the simulated environment labeled by sensor.",
	}, []string{"sensor_id"})
	readingsTotalVec = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "cm_readings_total",
		Help: "Total number of sensor readings processed labeled by sensor.",
	}, []string{"sensor_id"})
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

type Sensor struct {
	Id       string
	Current  float64
}

func NewSensor(id int) *Sensor {
	return &Sensor{
		Id:       fmt.Sprintf("sensor-%02d", id),
		Current:  25.0,
	}
}

func (s *Sensor) Run() {
	slog.Info("Starting sensor simulation...", "sensor_id", s.Id)

	tempGauge := ambientTempVec.WithLabelValues(s.Id)
	counter := readingsTotalVec.WithLabelValues(s.Id)

	for {
		delta := (rand.Float64() - 0.5)
		s.Current += delta
		tempGauge.Set(s.Current)
		counter.Inc()
		time.Sleep(2 * time.Second)
	}
}

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
	
	for i := 1; i <= 10; i++ {
		sensor := NewSensor(i)
		go sensor.Run()
	}

	slog.Info("Starting Cloudmetrics App", "version", app.GitCommit)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", app.IndexHandler)
	mux.HandleFunc("GET /healthz", HealthzHandler)
	mux.Handle("GET /metrics", promhttp.Handler())
	mux.HandleFunc("GET /styles.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "styles.css")
	})
	
	loggingMux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
