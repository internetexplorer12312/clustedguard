package app

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"ClusterGuard/internal/agentclient"
	"ClusterGuard/internal/domain"
	leveldbrepo "ClusterGuard/internal/repository/leveldb"
	"ClusterGuard/internal/service"
)

// Container wires dependencies (Dependency Injection root).
type Container struct {
	ServerService    *service.ServerService
	ClusterService   *service.ClusterService
	MonitorService   *service.MonitorService
	MetricsService   *service.MetricsService
	AlertService     *service.AlertService
	CollectorService *service.CollectorService
	store            *leveldbrepo.Store
}

func NewContainer(notifier domain.AlertNotifier) (*Container, error) {
	dataDir, err := dataDirectory()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}

	store, err := leveldbrepo.Open(filepath.Join(dataDir, "db"))
	if err != nil {
		return nil, err
	}

	serverRepo := leveldbrepo.NewServerRepository(store)
	clusterRepo := leveldbrepo.NewClusterRepository(store)
	metricsRepo := leveldbrepo.NewMetricsRepository(store)
	alertRepo := leveldbrepo.NewAlertRepository(store)

	serverSvc := service.NewServerService(serverRepo)
	clusterSvc := service.NewClusterService(clusterRepo)
	monitorSvc := service.NewMonitorService(serverSvc)
	fetcher := agentclient.NewHTTPFetcher(8 * time.Second)
	metricsSvc := service.NewMetricsService(serverSvc, fetcher, metricsRepo)
	alertSvc := service.NewAlertService(alertRepo, notifier)
	collectorSvc := service.NewCollectorService(serverSvc, monitorSvc, metricsSvc, alertSvc, 30*time.Second)

	return &Container{
		ServerService:    serverSvc,
		ClusterService:   clusterSvc,
		MonitorService:   monitorSvc,
		MetricsService:   metricsSvc,
		AlertService:     alertSvc,
		CollectorService: collectorSvc,
		store:            store,
	}, nil
}

func (c *Container) Close() error {
	if c.store != nil {
		return c.store.Close()
	}
	return nil
}

func dataDirectory() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		home, herr := os.UserHomeDir()
		if herr != nil {
			return "", err
		}
		return filepath.Join(home, ".clusterguard", "data"), nil
	}
	return filepath.Join(configDir, "clusterguard", "data"), nil
}
