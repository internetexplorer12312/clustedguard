package agentclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"ClusterGuard/internal/domain"
	"ClusterGuard/pkg/agentapi"
)

// HTTPFetcher implements domain.MetricsFetcher.
type HTTPFetcher struct {
	Client  *http.Client
	Timeout time.Duration
}

func NewHTTPFetcher(timeout time.Duration) *HTTPFetcher {
	if timeout <= 0 {
		timeout = 8 * time.Second
	}
	return &HTTPFetcher{
		Client:  &http.Client{Timeout: timeout},
		Timeout: timeout,
	}
}

func (f *HTTPFetcher) Fetch(ctx context.Context, host string, port int, token string) (*domain.MetricsSample, error) {
	if port == 0 {
		port = 9100
	}
	url := fmt.Sprintf("http://%s:%d/metrics", host, port)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("X-ClusterGuard-Token", token)
	}

	resp, err := f.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("agent returned status %d", resp.StatusCode)
	}

	var m agentapi.Metrics
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}

	return &domain.MetricsSample{
		Timestamp:     m.Timestamp,
		CPUPercent:    m.CPUPercent,
		MemPercent:    m.MemoryUsedPercent,
		DiskPercent:   m.DiskUsedPercent,
		MemAvailBytes: m.MemoryAvailableBytes,
		DiskFreeBytes: m.DiskFreeBytes,
	}, nil
}

var _ domain.MetricsFetcher = (*HTTPFetcher)(nil)
