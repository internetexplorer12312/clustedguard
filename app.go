package main

import (
	"context"
	"fmt"

	"ClusterGuard/internal/app"
	"ClusterGuard/internal/domain"
	"ClusterGuard/internal/service"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App is the Wails binding layer (adapter — thin facade over services).
type App struct {
	ctx       context.Context
	cancel    context.CancelFunc
	container *app.Container
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	a.ctx = ctx
	a.cancel = cancel

	container, err := app.NewContainer(a)
	if err != nil {
		fmt.Println("failed to init container:", err)
		return
	}
	a.container = container
	go container.CollectorService.Run(ctx)
}

func (a *App) shutdown(context.Context) {
	if a.cancel != nil {
		a.cancel()
	}
	if a.container != nil {
		_ = a.container.Close()
	}
}

// Notify implements domain.AlertNotifier — pushes alerts to the UI.
func (a *App) Notify(alert *domain.Alert) {
	if a.ctx == nil {
		return
	}
	dto := toAlertDTO(alert)
	runtime.EventsEmit(a.ctx, "alert", dto)
}

func (a *App) requireContainer() (*app.Container, error) {
	if a.container == nil {
		return nil, fmt.Errorf("application not initialized")
	}
	return a.container, nil
}

// --- DTOs ---

type ServerDTO struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Host          string   `json:"host"`
	Port          int      `json:"port"`
	Role          string   `json:"role"`
	Status        string   `json:"status"`
	Tags          []string `json:"tags"`
	CheckType     string   `json:"checkType"`
	CheckPath     string   `json:"checkPath"`
	LastCheck     int64    `json:"lastCheck"`
	LatencyMs     int64    `json:"latencyMs"`
	ClusterID     string   `json:"clusterId"`
	Notes         string   `json:"notes"`
	UseAgent      bool     `json:"useAgent"`
	AgentPort     int      `json:"agentPort"`
	AgentToken    string   `json:"agentToken"`
	CpuThreshold  float64  `json:"cpuThreshold"`
	MemThreshold  float64  `json:"memThreshold"`
	DiskThreshold float64  `json:"diskThreshold"`
	CpuPercent    float64  `json:"cpuPercent"`
	MemPercent    float64  `json:"memPercent"`
	DiskPercent   float64  `json:"diskPercent"`
	MemAvailBytes uint64   `json:"memAvailBytes"`
	DiskFreeBytes uint64   `json:"diskFreeBytes"`
}

type ClusterDTO struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ServerIDs   []string `json:"serverIds"`
	CreatedAt   int64    `json:"createdAt"`
	UpdatedAt   int64    `json:"updatedAt"`
}

type ClusterSummaryDTO struct {
	ClusterDTO
	TotalServers int `json:"totalServers"`
	OnlineCount  int `json:"onlineCount"`
	OfflineCount int `json:"offlineCount"`
}

type ServerInputDTO struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Host          string   `json:"host"`
	Port          int      `json:"port"`
	Role          string   `json:"role"`
	Tags          []string `json:"tags"`
	CheckType     string   `json:"checkType"`
	CheckPath     string   `json:"checkPath"`
	ClusterID     string   `json:"clusterId"`
	Notes         string   `json:"notes"`
	UseAgent      bool     `json:"useAgent"`
	AgentPort     int      `json:"agentPort"`
	AgentToken    string   `json:"agentToken"`
	CpuThreshold  float64  `json:"cpuThreshold"`
	MemThreshold  float64  `json:"memThreshold"`
	DiskThreshold float64  `json:"diskThreshold"`
}

type ClusterInputDTO struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ServerIDs   []string `json:"serverIds"`
}

type DashboardStatsDTO struct {
	TotalServers  int `json:"totalServers"`
	OnlineServers int `json:"onlineServers"`
	TotalClusters int `json:"totalClusters"`
	UnreadAlerts  int `json:"unreadAlerts"`
}

type MetricsSampleDTO struct {
	ServerID      string  `json:"serverId"`
	Timestamp     int64   `json:"timestamp"`
	CpuPercent    float64 `json:"cpuPercent"`
	MemPercent    float64 `json:"memPercent"`
	DiskPercent   float64 `json:"diskPercent"`
	MemAvailBytes uint64  `json:"memAvailBytes"`
	DiskFreeBytes uint64  `json:"diskFreeBytes"`
}

type AlertDTO struct {
	ID         string  `json:"id"`
	ServerID   string  `json:"serverId"`
	ServerName string  `json:"serverName"`
	Kind       string  `json:"kind"`
	Value      float64 `json:"value"`
	Threshold  float64 `json:"threshold"`
	Message    string  `json:"message"`
	CreatedAt  int64   `json:"createdAt"`
	Read       bool    `json:"read"`
}

func toServerDTO(s *domain.Server) ServerDTO {
	return ServerDTO{
		ID:            s.ID,
		Name:          s.Name,
		Host:          s.Host,
		Port:          s.Port,
		Role:          string(s.Role),
		Status:        string(s.Status),
		Tags:          s.Tags,
		CheckType:     s.CheckType,
		CheckPath:     s.CheckPath,
		LastCheck:     s.LastCheck,
		LatencyMs:     s.LatencyMs,
		ClusterID:     s.ClusterID,
		Notes:         s.Notes,
		UseAgent:      s.UseAgent,
		AgentPort:     s.AgentPort,
		AgentToken:    s.AgentToken,
		CpuThreshold:  s.CpuThreshold,
		MemThreshold:  s.MemThreshold,
		DiskThreshold: s.DiskThreshold,
		CpuPercent:    s.CpuPercent,
		MemPercent:    s.MemPercent,
		DiskPercent:   s.DiskPercent,
		MemAvailBytes: s.MemAvailBytes,
		DiskFreeBytes: s.DiskFreeBytes,
	}
}

func toClusterDTO(c *domain.Cluster) ClusterDTO {
	return ClusterDTO{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		ServerIDs:   c.ServerIDs,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func toServerInput(in ServerInputDTO) service.ServerInput {
	return service.ServerInput{
		ID:            in.ID,
		Name:          in.Name,
		Host:          in.Host,
		Port:          in.Port,
		Role:          domain.ServerRole(in.Role),
		Tags:          in.Tags,
		CheckType:     in.CheckType,
		CheckPath:     in.CheckPath,
		ClusterID:     in.ClusterID,
		Notes:         in.Notes,
		UseAgent:      true,
		AgentPort:     in.AgentPort,
		AgentToken:    in.AgentToken,
		CpuThreshold:  in.CpuThreshold,
		MemThreshold:  in.MemThreshold,
		DiskThreshold: in.DiskThreshold,
	}
}

func toClusterInput(in ClusterInputDTO) service.ClusterInput {
	return service.ClusterInput{
		ID:          in.ID,
		Name:        in.Name,
		Description: in.Description,
		ServerIDs:   in.ServerIDs,
	}
}

func toMetricsSampleDTO(m *domain.MetricsSample) MetricsSampleDTO {
	return MetricsSampleDTO{
		ServerID:      m.ServerID,
		Timestamp:     m.Timestamp,
		CpuPercent:    m.CPUPercent,
		MemPercent:    m.MemPercent,
		DiskPercent:   m.DiskPercent,
		MemAvailBytes: m.MemAvailBytes,
		DiskFreeBytes: m.DiskFreeBytes,
	}
}

func toAlertDTO(a *domain.Alert) AlertDTO {
	return AlertDTO{
		ID:         a.ID,
		ServerID:   a.ServerID,
		ServerName: a.ServerName,
		Kind:       string(a.Kind),
		Value:      a.Value,
		Threshold:  a.Threshold,
		Message:    a.Message,
		CreatedAt:  a.CreatedAt,
		Read:       a.Read,
	}
}

// --- Server API ---

func (a *App) ListServers() ([]ServerDTO, error) {
	c, err := a.requireContainer()
	if err != nil {
		return nil, err
	}
	servers, err := c.ServerService.List(a.ctx)
	if err != nil {
		return nil, err
	}
	out := make([]ServerDTO, len(servers))
	for i, s := range servers {
		out[i] = toServerDTO(s)
	}
	return out, nil
}

func (a *App) CreateServer(in ServerInputDTO) (ServerDTO, error) {
	c, err := a.requireContainer()
	if err != nil {
		return ServerDTO{}, err
	}
	s, err := c.ServerService.Create(a.ctx, toServerInput(in))
	if err != nil {
		return ServerDTO{}, err
	}
	return toServerDTO(s), nil
}

func (a *App) UpdateServer(in ServerInputDTO) (ServerDTO, error) {
	c, err := a.requireContainer()
	if err != nil {
		return ServerDTO{}, err
	}
	s, err := c.ServerService.Update(a.ctx, toServerInput(in))
	if err != nil {
		return ServerDTO{}, err
	}
	return toServerDTO(s), nil
}

func (a *App) DeleteServer(id string) error {
	c, err := a.requireContainer()
	if err != nil {
		return err
	}
	_ = c.MetricsService.DeleteHistory(a.ctx, id)
	return c.ServerService.Delete(a.ctx, id)
}

func (a *App) CheckServer(id string) (ServerDTO, error) {
	c, err := a.requireContainer()
	if err != nil {
		return ServerDTO{}, err
	}
	s, err := c.MonitorService.CheckServer(a.ctx, id)
	if err != nil {
		return ServerDTO{}, err
	}
	_ = c.CollectorService.CollectNow(a.ctx, id)
	s, _ = c.ServerService.Get(a.ctx, id)
	return toServerDTO(s), nil
}

func (a *App) CheckAllServers() ([]ServerDTO, error) {
	c, err := a.requireContainer()
	if err != nil {
		return nil, err
	}
	servers, err := c.MonitorService.CheckAll(a.ctx)
	if err != nil {
		return nil, err
	}
	for _, s := range servers {
		_ = c.CollectorService.CollectNow(a.ctx, s.ID)
	}
	all, _ := c.ServerService.List(a.ctx)
	out := make([]ServerDTO, len(all))
	for i, s := range all {
		out[i] = toServerDTO(s)
	}
	return out, nil
}

func (a *App) GetMetricsHistory(serverID string, limit int) ([]MetricsSampleDTO, error) {
	c, err := a.requireContainer()
	if err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 60
	}
	samples, err := c.MetricsService.History(a.ctx, serverID, limit)
	if err != nil {
		return nil, err
	}
	out := make([]MetricsSampleDTO, len(samples))
	for i, s := range samples {
		out[i] = toMetricsSampleDTO(s)
	}
	return out, nil
}

func (a *App) CollectServerMetrics(serverID string) (ServerDTO, error) {
	c, err := a.requireContainer()
	if err != nil {
		return ServerDTO{}, err
	}
	if err := c.CollectorService.CollectNow(a.ctx, serverID); err != nil {
		return ServerDTO{}, err
	}
	s, err := c.ServerService.Get(a.ctx, serverID)
	if err != nil {
		return ServerDTO{}, err
	}
	return toServerDTO(s), nil
}

// --- Cluster API ---

func (a *App) ListClusters() ([]ClusterSummaryDTO, error) {
	c, err := a.requireContainer()
	if err != nil {
		return nil, err
	}
	clusters, err := c.ClusterService.List(a.ctx)
	if err != nil {
		return nil, err
	}
	allServers, _ := c.ServerService.List(a.ctx)

	out := make([]ClusterSummaryDTO, len(clusters))
	for i, cl := range clusters {
		summary := ClusterSummaryDTO{ClusterDTO: toClusterDTO(cl)}
		for _, srv := range allServers {
			if srv.ClusterID == cl.ID {
				summary.TotalServers++
				switch srv.Status {
				case domain.StatusOnline, domain.StatusDegraded:
					summary.OnlineCount++
				case domain.StatusOffline:
					summary.OfflineCount++
				}
			}
		}
		out[i] = summary
	}
	return out, nil
}

func (a *App) CreateCluster(in ClusterInputDTO) (ClusterDTO, error) {
	c, err := a.requireContainer()
	if err != nil {
		return ClusterDTO{}, err
	}
	cl, err := c.ClusterService.Create(a.ctx, toClusterInput(in))
	if err != nil {
		return ClusterDTO{}, err
	}
	return toClusterDTO(cl), nil
}

func (a *App) UpdateCluster(in ClusterInputDTO) (ClusterDTO, error) {
	c, err := a.requireContainer()
	if err != nil {
		return ClusterDTO{}, err
	}
	cl, err := c.ClusterService.Update(a.ctx, toClusterInput(in))
	if err != nil {
		return ClusterDTO{}, err
	}
	return toClusterDTO(cl), nil
}

func (a *App) DeleteCluster(id string) error {
	c, err := a.requireContainer()
	if err != nil {
		return err
	}
	return c.ClusterService.Delete(a.ctx, id)
}

func (a *App) CheckCluster(id string) ([]ServerDTO, error) {
	c, err := a.requireContainer()
	if err != nil {
		return nil, err
	}
	servers, err := c.MonitorService.CheckCluster(a.ctx, id)
	if err != nil {
		return nil, err
	}
	for _, s := range servers {
		_ = c.CollectorService.CollectNow(a.ctx, s.ID)
	}
	out := make([]ServerDTO, len(servers))
	for i, s := range servers {
		updated, _ := c.ServerService.Get(a.ctx, s.ID)
		if updated != nil {
			s = updated
		}
		out[i] = toServerDTO(s)
	}
	return out, nil
}

// --- Alerts ---

func (a *App) ListAlerts(limit int) ([]AlertDTO, error) {
	c, err := a.requireContainer()
	if err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 50
	}
	alerts, err := c.AlertService.List(a.ctx, limit)
	if err != nil {
		return nil, err
	}
	out := make([]AlertDTO, len(alerts))
	for i, al := range alerts {
		out[i] = toAlertDTO(al)
	}
	return out, nil
}

func (a *App) MarkAlertRead(id string) error {
	c, err := a.requireContainer()
	if err != nil {
		return err
	}
	return c.AlertService.MarkRead(a.ctx, id)
}

func (a *App) DeleteAlert(id string) error {
	c, err := a.requireContainer()
	if err != nil {
		return err
	}
	return c.AlertService.Delete(a.ctx, id)
}

func (a *App) GetDashboardStats() (DashboardStatsDTO, error) {
	c, err := a.requireContainer()
	if err != nil {
		return DashboardStatsDTO{}, err
	}
	servers, err := c.ServerService.List(a.ctx)
	if err != nil {
		return DashboardStatsDTO{}, err
	}
	clusters, err := c.ClusterService.List(a.ctx)
	if err != nil {
		return DashboardStatsDTO{}, err
	}
	alerts, _ := c.AlertService.List(a.ctx, 100)

	stats := DashboardStatsDTO{
		TotalServers:  len(servers),
		TotalClusters: len(clusters),
	}
	for _, s := range servers {
		if s.Status == domain.StatusOnline || s.Status == domain.StatusDegraded {
			stats.OnlineServers++
		}
	}
	for _, al := range alerts {
		if !al.Read {
			stats.UnreadAlerts++
		}
	}
	return stats, nil
}
