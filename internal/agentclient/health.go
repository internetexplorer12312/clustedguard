package agentclient

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// CheckHealth probes the ClusterGuard agent /health endpoint.
func CheckHealth(ctx context.Context, host string, port int, token string) (online bool, latencyMs int64, err error) {
	if port == 0 {
		port = 9100
	}
	url := fmt.Sprintf("http://%s:%d/health", host, port)

	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, 0, err
	}
	if token != "" {
		req.Header.Set("X-ClusterGuard-Token", token)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	latencyMs = time.Since(start).Milliseconds()
	if err != nil {
		return false, latencyMs, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, latencyMs, nil
}
