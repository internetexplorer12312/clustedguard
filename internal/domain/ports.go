package domain

import "context"

// ServerRepository persists and retrieves server entities.
type ServerRepository interface {
	Save(ctx context.Context, server *Server) error
	GetByID(ctx context.Context, id string) (*Server, error)
	List(ctx context.Context) ([]*Server, error)
	Delete(ctx context.Context, id string) error
}

// ClusterRepository persists and retrieves cluster entities.
type ClusterRepository interface {
	Save(ctx context.Context, cluster *Cluster) error
	GetByID(ctx context.Context, id string) (*Cluster, error)
	List(ctx context.Context) ([]*Cluster, error)
	Delete(ctx context.Context, id string) error
}

// HealthChecker probes a single node (Strategy — Open/Closed).
type HealthChecker interface {
	Type() string
	Check(ctx context.Context, host string, port int, path string) (online bool, latencyMs int64, err error)
}

// HealthCheckerRegistry resolves checkers by type (Interface Segregation).
type HealthCheckerRegistry interface {
	Get(checkType string) (HealthChecker, bool)
}

// MetricsRepository stores time-series samples for charts.
type MetricsRepository interface {
	Append(ctx context.Context, sample *MetricsSample) error
	ListByServer(ctx context.Context, serverID string, limit int) ([]*MetricsSample, error)
	DeleteByServer(ctx context.Context, serverID string) error
}

// AlertRepository persists alert records.
type AlertRepository interface {
	Save(ctx context.Context, alert *Alert) error
	List(ctx context.Context, limit int) ([]*Alert, error)
	MarkRead(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}

// MetricsFetcher retrieves metrics from a remote agent (DIP).
type MetricsFetcher interface {
	Fetch(ctx context.Context, host string, port int, token string) (*MetricsSample, error)
}

// AlertNotifier delivers alerts to the UI (DIP).
type AlertNotifier interface {
	Notify(alert *Alert)
}
