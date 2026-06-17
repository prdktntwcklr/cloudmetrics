// +build smoke

package main

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestHealthzSmoke(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "cloudmetrics:smoke-test",
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
