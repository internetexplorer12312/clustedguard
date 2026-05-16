package health

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"ClusterGuard/internal/domain"
)

// HTTPChecker probes reachability via HTTP GET.
type HTTPChecker struct {
	Timeout time.Duration
	Client  *http.Client
}

func NewHTTPChecker(timeout time.Duration) *HTTPChecker {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &HTTPChecker{
		Timeout: timeout,
		Client:  &http.Client{Timeout: timeout},
	}
}

func (c *HTTPChecker) Type() string { return "http" }

func (c *HTTPChecker) Check(ctx context.Context, host string, port int, path string) (bool, int64, error) {
	if path == "" {
		path = "/"
	}
	url := fmt.Sprintf("http://%s:%d%s", host, port, path)
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, 0, err
	}

	resp, err := c.Client.Do(req)
	latency := time.Since(start).Milliseconds()
	if err != nil {
		return false, latency, nil
	}
	defer resp.Body.Close()

	online := resp.StatusCode >= 200 && resp.StatusCode < 500
	return online, latency, nil
}

var _ domain.HealthChecker = (*HTTPChecker)(nil)
