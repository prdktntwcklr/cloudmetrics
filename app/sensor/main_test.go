package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestNewSensor(t *testing.T) {
	expectedId := 42

	sensor := NewSensor(42, "localhost")
	if sensor.Id != expectedId {
		t.Errorf("Expected sensor id %v, got %v",
			expectedId, sensor.Id)
		}
}

type MockHTTPClient struct {
    MockPost func(url, contentType string, body io.Reader) (*http.Response, error)
}

func (m *MockHTTPClient) Post(url, contentType string, body io.Reader) (*http.Response, error) {
    return m.MockPost(url, contentType, body)
}

func TestSensor_Tick(t *testing.T) {
	expectedSensorID := 123

    called := false
    mockClient := &MockHTTPClient{
        MockPost: func(url, contentType string, body io.Reader) (*http.Response, error) {
            called = true
            
            // Assert URL and Content-Type
            if url != "http://fake-api.local/reading" {
                t.Errorf("expected URL 'http://fake-api.local/reading', got '%s'", url)
            }
			if contentType != "application/json" {
				t.Errorf("expected content type 'application/json', got '%s'", contentType)
			}

			// Verify Payload
			var actualPayload Reading
			decoder := json.NewDecoder(body)
			if err := decoder.Decode(&actualPayload); err != nil {
				t.Fatalf("failed to decode request body JSON: %v", err)
			}

			if actualPayload.SensorID != expectedSensorID {
				t.Errorf("expected JSON SensorID to be '%v', got '%v'", expectedSensorID, actualPayload.SensorID)
			}

            // Return a mock 200 OK response
            return &http.Response{
                StatusCode: http.StatusOK,
                Status:     "200 OK",
                Body:       io.NopCloser(bytes.NewBufferString(`{"status":"success"}`)),
            }, nil
        },
    }

    sensor := &Sensor{
        Id:        123,
        Current:   25.0,
        TargetURL: "http://fake-api.local/reading",
        Client:    mockClient,
    }

    sensor.Tick()

    if !called {
        t.Error("expected HTTP client Post to be called, but it wasn't")
    }
}
