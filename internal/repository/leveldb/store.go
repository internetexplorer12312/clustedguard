package leveldb

import (
	"encoding/json"
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

const (
	prefixServer  = "server:"
	prefixCluster = "cluster:"
	keyServerIdx  = "index:servers"
	keyClusterIdx = "index:clusters"
)

// Store wraps LevelDB with JSON serialization helpers.
type Store struct {
	db *leveldb.DB
}

func Open(path string) (*Store, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, fmt.Errorf("open leveldb: %w", err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) putJSON(key string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return s.db.Put([]byte(key), data, nil)
}

func (s *Store) getJSON(key string, v any) error {
	data, err := s.db.Get([]byte(key), nil)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func (s *Store) delete(key string) error {
	return s.db.Delete([]byte(key), nil)
}

func (s *Store) loadIndex(key string) ([]string, error) {
	var ids []string
	err := s.getJSON(key, &ids)
	if err == leveldb.ErrNotFound {
		return []string{}, nil
	}
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (s *Store) saveIndex(key string, ids []string) error {
	return s.putJSON(key, ids)
}

func addToIndex(ids []string, id string) []string {
	for _, existing := range ids {
		if existing == id {
			return ids
		}
	}
	return append(ids, id)
}

func removeFromIndex(ids []string, id string) []string {
	out := make([]string, 0, len(ids))
	for _, existing := range ids {
		if existing != id {
			out = append(out, existing)
		}
	}
	return out
}

func (s *Store) listByPrefix(prefix string) ([][]byte, error) {
	iter := s.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	defer iter.Release()

	var items [][]byte
	for iter.Next() {
		val := make([]byte, len(iter.Value()))
		copy(val, iter.Value())
		items = append(items, val)
	}
	return items, iter.Error()
}
