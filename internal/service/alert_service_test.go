package service

import (
	"context"
	"testing"

	"ClusterGuard/internal/domain"
)

type memAlertRepo struct {
	saved []*domain.Alert
}

func (m *memAlertRepo) Save(_ context.Context, a *domain.Alert) error {
	m.saved = append(m.saved, a)
	return nil
}
func (m *memAlertRepo) List(_ context.Context, _ int) ([]*domain.Alert, error) { return m.saved, nil }
func (m *memAlertRepo) MarkRead(_ context.Context, _ string) error               { return nil }
func (m *memAlertRepo) Delete(_ context.Context, _ string) error                 { return nil }

func TestAlertService_NoAlertBelowThreshold(t *testing.T) {
	repo := &memAlertRepo{}
	svc := NewAlertService(repo, nil)
	srv := &domain.Server{
		ID: "s1", Name: "web-01", CpuPercent: 50, CpuThreshold: 90,
		MemPercent: 40, MemThreshold: 85, DiskPercent: 30, DiskThreshold: 90,
	}
	got, err := svc.Evaluate(context.Background(), srv)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 || len(repo.saved) != 0 {
		t.Fatalf("expected no alerts, got %d saved %d", len(got), len(repo.saved))
	}
}

func TestAlertService_CreatesAlertWhenOverThreshold(t *testing.T) {
	repo := &memAlertRepo{}
	svc := NewAlertService(repo, nil)
	srv := &domain.Server{
		ID: "s1", Name: "web-01", CpuPercent: 92.3, CpuThreshold: 90,
		MemPercent: 40, MemThreshold: 85, DiskPercent: 30, DiskThreshold: 90,
	}
	_, err := svc.Evaluate(context.Background(), srv)
	if err != nil {
		t.Fatal(err)
	}
	if len(repo.saved) != 1 {
		t.Fatalf("expected 1 saved alert, got %d", len(repo.saved))
	}
	if repo.saved[0].Kind != domain.AlertCPU {
		t.Fatalf("expected CPU alert, got %s", repo.saved[0].Kind)
	}
}
