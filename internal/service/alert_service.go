package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"ClusterGuard/internal/domain"

	"github.com/google/uuid"
)

// AlertService evaluates thresholds and emits notifications.
type AlertService struct {
	repo           domain.AlertRepository
	notifier       domain.AlertNotifier
	notifyInterval time.Duration
	saveCooldown   time.Duration
	lastNotify     map[string]int64
	lastSaved      map[string]int64
	mu             sync.Mutex
}

func NewAlertService(repo domain.AlertRepository, notifier domain.AlertNotifier) *AlertService {
	return &AlertService{
		repo:           repo,
		notifier:       notifier,
		notifyInterval: 30 * time.Second,
		saveCooldown:   5 * time.Minute,
		lastNotify:     make(map[string]int64),
		lastSaved:      make(map[string]int64),
	}
}

func (a *AlertService) Evaluate(ctx context.Context, server *domain.Server) ([]*domain.Alert, error) {
	if server == nil {
		return nil, nil
	}

	now := time.Now().Unix()
	var created []*domain.Alert

	checks := []struct {
		kind      domain.AlertKind
		value     float64
		threshold float64
	}{
		{domain.AlertCPU, server.CpuPercent, server.CpuThreshold},
		{domain.AlertMemory, server.MemPercent, server.MemThreshold},
		{domain.AlertDisk, server.DiskPercent, server.DiskThreshold},
	}

	for _, c := range checks {
		key := server.ID + ":" + string(c.kind)
		if c.threshold <= 0 || c.value < c.threshold {
			a.clearKey(key)
			continue
		}

		if a.shouldNotify(key, now) {
			a.markNotify(key, now)
			alert := a.buildAlert(server, c.kind, c.value, c.threshold, now)
			if a.notifier != nil {
				a.notifier.Notify(alert)
			}
		}

		if !a.shouldSave(key, now) {
			continue
		}
		a.markSaved(key, now)

		alert := a.buildAlert(server, c.kind, c.value, c.threshold, now)
		alert.ID = uuid.NewString()
		if err := a.repo.Save(ctx, alert); err != nil {
			return created, err
		}
		created = append(created, alert)
	}

	return created, nil
}

func (a *AlertService) clearKey(key string) {
	a.mu.Lock()
	delete(a.lastNotify, key)
	delete(a.lastSaved, key)
	a.mu.Unlock()
}

func (a *AlertService) shouldNotify(key string, now int64) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	last, ok := a.lastNotify[key]
	if !ok {
		return true
	}
	return now-last >= int64(a.notifyInterval.Seconds())
}

func (a *AlertService) shouldSave(key string, now int64) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	last, ok := a.lastSaved[key]
	if !ok {
		return true
	}
	return now-last >= int64(a.saveCooldown.Seconds())
}

func (a *AlertService) markNotify(key string, now int64) {
	a.mu.Lock()
	a.lastNotify[key] = now
	a.mu.Unlock()
}

func (a *AlertService) markSaved(key string, now int64) {
	a.mu.Lock()
	a.lastSaved[key] = now
	a.mu.Unlock()
}

func (a *AlertService) buildAlert(
	server *domain.Server,
	kind domain.AlertKind,
	value, threshold float64,
	now int64,
) *domain.Alert {
	kindRu := map[domain.AlertKind]string{
		domain.AlertCPU:    "ЦП",
		domain.AlertMemory: "ОЗУ",
		domain.AlertDisk:   "Диск",
	}
	label := kindRu[kind]
	if label == "" {
		label = string(kind)
	}
	msg := fmt.Sprintf("%s: %s %.1f%% (порог %.1f%%)", server.Name, label, value, threshold)
	return &domain.Alert{
		ServerID:   server.ID,
		ServerName: server.Name,
		Kind:       kind,
		Value:      value,
		Threshold:  threshold,
		Message:    msg,
		CreatedAt:  now,
		Read:       false,
	}
}

func (a *AlertService) List(ctx context.Context, limit int) ([]*domain.Alert, error) {
	return a.repo.List(ctx, limit)
}

func (a *AlertService) MarkRead(ctx context.Context, id string) error {
	return a.repo.MarkRead(ctx, id)
}

func (a *AlertService) Delete(ctx context.Context, id string) error {
	return a.repo.Delete(ctx, id)
}
