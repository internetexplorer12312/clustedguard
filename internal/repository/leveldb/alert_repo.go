package leveldb

import (
	"context"
	"fmt"
	"sort"

	"ClusterGuard/internal/domain"

	"github.com/syndtr/goleveldb/leveldb"
)

const (
	prefixAlert = "alert:"
	keyAlertIdx = "index:alerts"
	maxAlerts   = 200
)

// AlertRepository implements domain.AlertRepository.
type AlertRepository struct {
	store *Store
}

func NewAlertRepository(store *Store) *AlertRepository {
	return &AlertRepository{store: store}
}

func alertKey(id string) string { return prefixAlert + id }

func (r *AlertRepository) Save(ctx context.Context, alert *domain.Alert) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if alert.ID == "" {
		return fmt.Errorf("alert id required")
	}

	ids, err := r.store.loadIndex(keyAlertIdx)
	if err != nil {
		return err
	}
	ids = addToIndex(ids, alert.ID)
	if len(ids) > maxAlerts {
		ids = ids[len(ids)-maxAlerts:]
	}
	if err := r.store.saveIndex(keyAlertIdx, ids); err != nil {
		return err
	}
	return r.store.putJSON(alertKey(alert.ID), alert)
}

func (r *AlertRepository) List(ctx context.Context, limit int) ([]*domain.Alert, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	ids, err := r.store.loadIndex(keyAlertIdx)
	if err != nil {
		return nil, err
	}

	var alerts []*domain.Alert
	for _, id := range ids {
		var a domain.Alert
		if err := r.store.getJSON(alertKey(id), &a); err != nil {
			continue
		}
		alerts = append(alerts, &a)
	}
	sort.Slice(alerts, func(i, j int) bool {
		return alerts[i].CreatedAt > alerts[j].CreatedAt
	})
	if limit > 0 && len(alerts) > limit {
		alerts = alerts[:limit]
	}
	return alerts, nil
}

func (r *AlertRepository) MarkRead(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	var a domain.Alert
	if err := r.store.getJSON(alertKey(id), &a); err != nil {
		if err == leveldb.ErrNotFound {
			return fmt.Errorf("alert not found")
		}
		return err
	}
	a.Read = true
	return r.store.putJSON(alertKey(id), &a)
}

func (r *AlertRepository) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	ids, err := r.store.loadIndex(keyAlertIdx)
	if err != nil {
		return err
	}
	ids = removeFromIndex(ids, id)
	if err := r.store.saveIndex(keyAlertIdx, ids); err != nil {
		return err
	}
	return r.store.delete(alertKey(id))
}

var _ domain.AlertRepository = (*AlertRepository)(nil)
