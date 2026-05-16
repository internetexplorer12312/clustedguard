package leveldb

import (
	"context"
	"fmt"

	"ClusterGuard/internal/domain"

	"github.com/syndtr/goleveldb/leveldb"
)

// ClusterRepository implements domain.ClusterRepository.
type ClusterRepository struct {
	store *Store
}

func NewClusterRepository(store *Store) *ClusterRepository {
	return &ClusterRepository{store: store}
}

func clusterKey(id string) string { return prefixCluster + id }

func (r *ClusterRepository) Save(ctx context.Context, cluster *domain.Cluster) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if cluster.ID == "" {
		return fmt.Errorf("cluster id is required")
	}

	ids, err := r.store.loadIndex(keyClusterIdx)
	if err != nil {
		return err
	}
	ids = addToIndex(ids, cluster.ID)
	if err := r.store.saveIndex(keyClusterIdx, ids); err != nil {
		return err
	}
	return r.store.putJSON(clusterKey(cluster.ID), cluster)
}

func (r *ClusterRepository) GetByID(ctx context.Context, id string) (*domain.Cluster, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var cluster domain.Cluster
	if err := r.store.getJSON(clusterKey(id), &cluster); err != nil {
		if err == leveldb.ErrNotFound {
			return nil, fmt.Errorf("cluster %s not found", id)
		}
		return nil, err
	}
	return &cluster, nil
}

func (r *ClusterRepository) List(ctx context.Context) ([]*domain.Cluster, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	ids, err := r.store.loadIndex(keyClusterIdx)
	if err != nil {
		return nil, err
	}

	clusters := make([]*domain.Cluster, 0, len(ids))
	for _, id := range ids {
		c, err := r.GetByID(ctx, id)
		if err != nil {
			continue
		}
		clusters = append(clusters, c)
	}
	return clusters, nil
}

func (r *ClusterRepository) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	ids, err := r.store.loadIndex(keyClusterIdx)
	if err != nil {
		return err
	}
	ids = removeFromIndex(ids, id)
	if err := r.store.saveIndex(keyClusterIdx, ids); err != nil {
		return err
	}
	return r.store.delete(clusterKey(id))
}

var _ domain.ClusterRepository = (*ClusterRepository)(nil)
