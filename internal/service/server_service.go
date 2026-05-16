package service

import (
	"context"
	"fmt"
	"time"

	"ClusterGuard/internal/domain"

	"github.com/google/uuid"
)

// ServerService handles server CRUD (Single Responsibility).
type ServerService struct {
	repo domain.ServerRepository
}

func NewServerService(repo domain.ServerRepository) *ServerService {
	return &ServerService{repo: repo}
}

type ServerInput struct {
	ID            string
	Name          string
	Host          string
	Port          int
	Role          domain.ServerRole
	Tags          []string
	CheckType     string
	CheckPath     string
	ClusterID     string
	Notes         string
	UseAgent      bool
	AgentPort     int
	AgentToken    string
	CpuThreshold  float64
	MemThreshold  float64
	DiskThreshold float64
}

func (s *ServerService) Create(ctx context.Context, in ServerInput) (*domain.Server, error) {
	server := &domain.Server{
		ID:            uuid.NewString(),
		Name:          in.Name,
		Host:          in.Host,
		Port:          in.Port,
		Role:          defaultRole(in.Role),
		Status:        domain.StatusUnknown,
		Tags:          in.Tags,
		CheckPath:     in.CheckPath,
		ClusterID:     in.ClusterID,
		Notes:         in.Notes,
		UseAgent:      true,
		AgentPort:     in.AgentPort,
		AgentToken:    in.AgentToken,
		CpuThreshold:  defaultThreshold(in.CpuThreshold, 90),
		MemThreshold:  defaultThreshold(in.MemThreshold, 90),
		DiskThreshold: defaultThreshold(in.DiskThreshold, 90),
	}
	applyAgentDefaults(server)
	if err := s.repo.Save(ctx, server); err != nil {
		return nil, err
	}
	return server, nil
}

func (s *ServerService) Update(ctx context.Context, in ServerInput) (*domain.Server, error) {
	if in.ID == "" {
		return nil, fmt.Errorf("server id is required")
	}
	existing, err := s.repo.GetByID(ctx, in.ID)
	if err != nil {
		return nil, err
	}

	existing.Name = in.Name
	existing.Host = in.Host
	if in.Port > 0 {
		existing.Port = in.Port
	}
	if in.Role != "" {
		existing.Role = in.Role
	}
	if in.Tags != nil {
		existing.Tags = in.Tags
	}
	existing.ClusterID = in.ClusterID
	existing.Notes = in.Notes
	existing.UseAgent = true
	if in.AgentPort > 0 {
		existing.AgentPort = in.AgentPort
	}
	existing.AgentToken = in.AgentToken
	applyAgentDefaults(existing)
	if in.CpuThreshold > 0 {
		existing.CpuThreshold = in.CpuThreshold
	}
	if in.MemThreshold > 0 {
		existing.MemThreshold = in.MemThreshold
	}
	if in.DiskThreshold > 0 {
		existing.DiskThreshold = in.DiskThreshold
	}

	if err := s.repo.Save(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *ServerService) Get(ctx context.Context, id string) (*domain.Server, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ServerService) List(ctx context.Context) ([]*domain.Server, error) {
	all, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, srv := range all {
		needsMigrate := !srv.UseAgent || srv.CheckType != "agent"
		applyAgentDefaults(srv)
		if needsMigrate {
			_ = s.repo.Save(ctx, srv)
		}
	}
	return all, nil
}

func (s *ServerService) ListByCluster(ctx context.Context, clusterID string) ([]*domain.Server, error) {
	all, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	var filtered []*domain.Server
	for _, srv := range all {
		if srv.ClusterID == clusterID {
			filtered = append(filtered, srv)
		}
	}
	return filtered, nil
}

func (s *ServerService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func defaultRole(r domain.ServerRole) domain.ServerRole {
	if r == "" {
		return domain.RoleWorker
	}
	return r
}

func applyAgentDefaults(server *domain.Server) {
	server.UseAgent = true
	server.CheckType = "agent"
	if server.AgentPort == 0 {
		server.AgentPort = 9100
	}
	if server.Port == 0 {
		server.Port = server.AgentPort
	}
}

func defaultThreshold(v, def float64) float64 {
	if v <= 0 {
		return def
	}
	return v
}

// ApplyMetrics saves latest resource usage on the server record.
func (s *ServerService) ApplyMetrics(ctx context.Context, id string, sample *domain.MetricsSample) (*domain.Server, error) {
	server, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	server.CpuPercent = sample.CPUPercent
	server.MemPercent = sample.MemPercent
	server.DiskPercent = sample.DiskPercent
	server.MemAvailBytes = sample.MemAvailBytes
	server.DiskFreeBytes = sample.DiskFreeBytes
	server.LastCheck = sample.Timestamp
	if err := s.repo.Save(ctx, server); err != nil {
		return nil, err
	}
	return server, nil
}

// ApplyHealthResult updates server status from a health probe.
func (s *ServerService) ApplyHealthResult(ctx context.Context, id string, online bool, latencyMs int64) (*domain.Server, error) {
	server, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	server.LastCheck = time.Now().Unix()
	server.LatencyMs = latencyMs
	if online {
		if latencyMs > 500 {
			server.Status = domain.StatusDegraded
		} else {
			server.Status = domain.StatusOnline
		}
	} else {
		server.Status = domain.StatusOffline
	}

	if err := s.repo.Save(ctx, server); err != nil {
		return nil, err
	}
	return server, nil
}
