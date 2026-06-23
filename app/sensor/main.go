package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// Injected from the Dockerfile
var Version = "unknown"

type HTTPClient interface {
    Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

type Sensor struct {
    Id        int
    Current   float64
    TargetURL string
    Client    HTTPClient
}

func NewSensor(id int, targetURL string) *Sensor {
	return &Sensor{
		Id:        id,
		Current:   25.0,
		TargetURL: targetURL,
		Client:    &http.Client{Timeout: 5 * time.Second},
	}
}

type Reading struct {
	SensorID    int     `json:"sensor_id"`
	Temperature float64 `json:"temperature"`
}

func (s *Sensor) Run(ctx context.Context) {
    slog.Info("Starting sensor simulation...", "sensor_id", s.Id)
    baseDelay := 2 * time.Second

    for {
        select {
        case <-ctx.Done():
            slog.Info("Stopping sensor simulation...", "sensor_id", s.Id)
            return
        default:
            s.Tick()
            
            jitter := time.Duration((rand.Float64() - 0.5) * float64(500*time.Millisecond))
            time.Sleep(baseDelay + jitter)
        }
    }
}

func (s *Sensor) Tick() {
    delta := (rand.Float64() - 0.5)
    s.Current += delta

    payload := Reading{
        SensorID:    s.Id,
        Temperature: s.Current,
    }
    
    jsonData, err := json.Marshal(payload)
    if err != nil {
        slog.Error("Failed to marshal JSON", "error", err)
        return
    }

	if s.Client == nil {
		slog.Error("HTTP client is not initialized", "sensor_id", s.Id)
		return
	}

    resp, err := s.Client.Post(s.TargetURL, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        slog.Error("Failed to send reading to API", "sensor_id", s.Id, "error", err)
        return
    }
    defer resp.Body.Close()
    
    slog.Debug("Successfully sent reading", "sensor_id", s.Id, "status", resp.Status)
}

func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug, 
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)

    apiURL := os.Getenv("API_URL")
    if apiURL == "" {
        apiURL = "http://localhost:8080/api/readings"
    }

	slog.Info("Starting Cloudmetrics Sensors", "target_api", apiURL, "version", Version)

	for i := 1; i <= 10; i++ {
		sensor := NewSensor(i, apiURL)
		go sensor.Run(context.Background())
	}

	// Keep the main routine alive indefinitely
	select {}
}
