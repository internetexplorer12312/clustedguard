package service

import (
	"context"
	"sync"

	"ClusterGuard/internal/agentclient"
	"ClusterGuard/internal/domain"
)

// MonitorService runs health checks via the ClusterGuard agent.
type MonitorService struct {
	servers *ServerService
}

func NewMonitorService(servers *ServerService) *MonitorService {
	return &MonitorService{servers: servers}
}

// CheckServer probes a single server and persists the result.
func (m *MonitorService) CheckServer(ctx context.Context, id string) (*domain.Server, error) {
	server, err := m.servers.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return m.probe(ctx, server)
}

// CheckAll probes every registered server concurrently.
func (m *MonitorService) CheckAll(ctx context.Context) ([]*domain.Server, error) {
	servers, err := m.servers.List(ctx)
	if err != nil {
		return nil, err
	}

	var (
		mu     sync.Mutex
		wg     sync.WaitGroup
		result []*domain.Server
	)

	for _, srv := range servers {
		wg.Add(1)
		go func(s *domain.Server) {
			defer wg.Done()
			updated, err := m.probe(ctx, s)
			if err != nil {
				return
			}
			mu.Lock()
			result = append(result, updated)
			mu.Unlock()
		}(srv)
	}

	wg.Wait()
	return result, nil
}

// CheckCluster probes all servers assigned to a cluster.
func (m *MonitorService) CheckCluster(ctx context.Context, clusterID string) ([]*domain.Server, error) {
	servers, err := m.servers.ListByCluster(ctx, clusterID)
	if err != nil {
		return nil, err
	}

	var result []*domain.Server
	for _, srv := range servers {
		updated, err := m.probe(ctx, srv)
		if err != nil {
			continue
		}
		result = append(result, updated)
	}
	return result, nil
}

func (m *MonitorService) probe(ctx context.Context, server *domain.Server) (*domain.Server, error) {
	online, latency, err := agentclient.CheckHealth(ctx, server.Host, server.AgentPort, server.AgentToken)
	if err != nil {
		return m.servers.ApplyHealthResult(ctx, server.ID, false, latency)
	}
	return m.servers.ApplyHealthResult(ctx, server.ID, online, latency)
}
