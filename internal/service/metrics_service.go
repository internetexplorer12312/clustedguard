package service

import (
	"context"
	"time"

	"ClusterGuard/internal/domain"
)

// MetricsService collects and stores agent metrics.
type MetricsService struct {
	servers *ServerService
	fetcher domain.MetricsFetcher
	repo    domain.MetricsRepository
}

func NewMetricsService(
	servers *ServerService,
	fetcher domain.MetricsFetcher,
	repo domain.MetricsRepository,
) *MetricsService {
	return &MetricsService{servers: servers, fetcher: fetcher, repo: repo}
}

func (m *MetricsService) Collect(ctx context.Context, serverID string) (*domain.MetricsSample, error) {
	server, err := m.servers.Get(ctx, serverID)
	if err != nil {
		return nil, err
	}
	sample, err := m.fetcher.Fetch(ctx, server.Host, server.AgentPort, server.AgentToken)
	if err != nil {
		return nil, err
	}
	sample.ServerID = server.ID
	if sample.Timestamp == 0 {
		sample.Timestamp = time.Now().Unix()
	}

	_, err = m.servers.ApplyMetrics(ctx, server.ID, sample)
	if err != nil {
		return nil, err
	}
	_ = m.repo.Append(ctx, sample)
	return sample, nil
}

func (m *MetricsService) History(ctx context.Context, serverID string, limit int) ([]*domain.MetricsSample, error) {
	return m.repo.ListByServer(ctx, serverID, limit)
}

func (m *MetricsService) DeleteHistory(ctx context.Context, serverID string) error {
	return m.repo.DeleteByServer(ctx, serverID)
}
