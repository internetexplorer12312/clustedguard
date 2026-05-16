package leveldb

import (
	"context"
	"encoding/json"
	"fmt"

	"ClusterGuard/internal/domain"

	"github.com/syndtr/goleveldb/leveldb"
)

// ServerRepository implements domain.ServerRepository.
type ServerRepository struct {
	store *Store
}

func NewServerRepository(store *Store) *ServerRepository {
	return &ServerRepository{store: store}
}

func serverKey(id string) string { return prefixServer + id }

func (r *ServerRepository) Save(ctx context.Context, server *domain.Server) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if server.ID == "" {
		return fmt.Errorf("server id is required")
	}

	ids, err := r.store.loadIndex(keyServerIdx)
	if err != nil {
		return err
	}
	ids = addToIndex(ids, server.ID)
	if err := r.store.saveIndex(keyServerIdx, ids); err != nil {
		return err
	}
	return r.store.putJSON(serverKey(server.ID), server)
}

func (r *ServerRepository) GetByID(ctx context.Context, id string) (*domain.Server, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var server domain.Server
	if err := r.store.getJSON(serverKey(id), &server); err != nil {
		if err == leveldb.ErrNotFound {
			return nil, fmt.Errorf("server %s not found", id)
		}
		return nil, err
	}
	return &server, nil
}

func (r *ServerRepository) List(ctx context.Context) ([]*domain.Server, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	ids, err := r.store.loadIndex(keyServerIdx)
	if err != nil {
		return nil, err
	}

	servers := make([]*domain.Server, 0, len(ids))
	for _, id := range ids {
		s, err := r.GetByID(ctx, id)
		if err != nil {
			continue
		}
		servers = append(servers, s)
	}
	return servers, nil
}

func (r *ServerRepository) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	ids, err := r.store.loadIndex(keyServerIdx)
	if err != nil {
		return err
	}
	ids = removeFromIndex(ids, id)
	if err := r.store.saveIndex(keyServerIdx, ids); err != nil {
		return err
	}
	return r.store.delete(serverKey(id))
}

func (r *ServerRepository) ListRaw() ([]*domain.Server, error) {
	raw, err := r.store.listByPrefix(prefixServer)
	if err != nil {
		return nil, err
	}
	servers := make([]*domain.Server, 0, len(raw))
	for _, data := range raw {
		var s domain.Server
		if json.Unmarshal(data, &s) == nil {
			servers = append(servers, &s)
		}
	}
	return servers, nil
}

var _ domain.ServerRepository = (*ServerRepository)(nil)
