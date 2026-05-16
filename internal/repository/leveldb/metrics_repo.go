package leveldb

import (
	"context"
	"sort"

	"ClusterGuard/internal/domain"
)

const (
	prefixMetrics = "metrics:"
	maxPerServer  = 120
)

// MetricsRepository implements domain.MetricsRepository.
type MetricsRepository struct {
	store *Store
}

func NewMetricsRepository(store *Store) *MetricsRepository {
	return &MetricsRepository{store: store}
}

func metricsKey(serverID string) string { return prefixMetrics + serverID }

func (r *MetricsRepository) Append(ctx context.Context, sample *domain.MetricsSample) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	var samples []*domain.MetricsSample
	_ = r.store.getJSON(metricsKey(sample.ServerID), &samples)

	samples = append(samples, sample)
	if len(samples) > maxPerServer {
		samples = samples[len(samples)-maxPerServer:]
	}
	return r.store.putJSON(metricsKey(sample.ServerID), samples)
}

func (r *MetricsRepository) ListByServer(ctx context.Context, serverID string, limit int) ([]*domain.MetricsSample, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var samples []*domain.MetricsSample
	if err := r.store.getJSON(metricsKey(serverID), &samples); err != nil {
		return []*domain.MetricsSample{}, nil
	}
	sort.Slice(samples, func(i, j int) bool {
		return samples[i].Timestamp < samples[j].Timestamp
	})
	if limit > 0 && len(samples) > limit {
		samples = samples[len(samples)-limit:]
	}
	return samples, nil
}

func (r *MetricsRepository) DeleteByServer(ctx context.Context, serverID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return r.store.delete(metricsKey(serverID))
}

var _ domain.MetricsRepository = (*MetricsRepository)(nil)
