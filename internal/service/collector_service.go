package service

import (
	"context"
	"log"
	"time"
)

// CollectorService periodically polls agents and runs health checks.
type CollectorService struct {
	servers  *ServerService
	monitor  *MonitorService
	metrics  *MetricsService
	alerts   *AlertService
	interval time.Duration
}

func NewCollectorService(
	servers *ServerService,
	monitor *MonitorService,
	metrics *MetricsService,
	alerts *AlertService,
	interval time.Duration,
) *CollectorService {
	if interval <= 0 {
		interval = 30 * time.Second
	}
	return &CollectorService{
		servers:  servers,
		monitor:  monitor,
		metrics:  metrics,
		alerts:   alerts,
		interval: interval,
	}
}

func (c *CollectorService) Run(ctx context.Context) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	c.tick(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.tick(ctx)
		}
	}
}

func (c *CollectorService) tick(ctx context.Context) {
	list, err := c.servers.List(ctx)
	if err != nil {
		return
	}
	for _, srv := range list {
		if _, err := c.monitor.CheckServer(ctx, srv.ID); err != nil {
			log.Printf("health check %s: %v", srv.Name, err)
		}
		sample, err := c.metrics.Collect(ctx, srv.ID)
		if err != nil {
			log.Printf("metrics %s: %v", srv.Name, err)
			continue
		}
		if sample == nil {
			continue
		}
		updated, err := c.servers.Get(ctx, srv.ID)
		if err != nil {
			continue
		}
		if _, err := c.alerts.Evaluate(ctx, updated); err != nil {
			log.Printf("alerts %s: %v", srv.Name, err)
		}
	}
}

func (c *CollectorService) CollectNow(ctx context.Context, serverID string) error {
	if _, err := c.metrics.Collect(ctx, serverID); err != nil {
		return err
	}
	srv, err := c.servers.Get(ctx, serverID)
	if err != nil {
		return err
	}
	_, err = c.alerts.Evaluate(ctx, srv)
	return err
}
