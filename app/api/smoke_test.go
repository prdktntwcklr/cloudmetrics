// +build smoke

package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestHealthzSmoke(t *testing.T) {
	ctx := context.Background()

	imageName := os.Getenv("TEST_IMAGE_NAME")
	if imageName == "" {
		t.Fatalf("Smoke test failed: TEST_IMAGE_NAME environment variable is not set.")
	}

	req := testcontainers.ContainerRequest{
		Image:        imageName,
		ExposedPorts: []string{"8080/tcp"},
		WaitingFor: wait.ForHTTP("/healthz").
			WithPort("8080/tcp").
			WithStatusCodeMatcher(func(status int) bool {
				return status == http.StatusOK
			}).
			WithStartupTimeout(15 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Could not start container: %v", err)
	}

	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("Could not terminate container: %v", err)
		}
	}()

	t.Log("SUCCESS: Container started and /healthz responded with 200 OK")
}

func TestIngestReadingSmoke(t *testing.T) {
	ctx := context.Background()

	imageName := os.Getenv("TEST_IMAGE_NAME")
	if imageName == "" {
		t.Fatalf("Smoke test failed: TEST_IMAGE_NAME environment variable is not set.")
	}

	req := testcontainers.ContainerRequest{
		Image:        imageName,
		ExposedPorts: []string{"8080/tcp"},
		WaitingFor: wait.ForHTTP("/healthz").
			WithPort("8080/tcp").
			WithStatusCodeMatcher(func(status int) bool {
				return status == http.StatusOK
			}).
			WithStartupTimeout(15 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Could not start container: %v", err)
	}

	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("Could not terminate container: %v", err)
		}
	}()

	mappedPort, err := container.MappedPort(ctx, "8080/tcp")
	if err != nil {
		t.Fatalf("Failed to get mapped port: %v", err)
	}

	apiURL := fmt.Sprintf("http://localhost:%s/api/readings", mappedPort.Port())
	payload := []byte(`{"sensor_id": 99, "temperature": 27.8}`)

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Failed to send POST request to container: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		t.Errorf("Expected smoke test status 202 Accepted, got %v", resp.Status)
	} else {
		t.Log("SUCCESS: Sent mock sensor data and container responded with 202 Accepted")
	}
}
