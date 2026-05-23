package leveldb

import (
	"context"
	"path/filepath"
	"testing"

	"ClusterGuard/internal/domain"
)

func TestServerRepository_PersistsAcrossReopen(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "db")

	save := func() error {
		s, err := Open(path)
		if err != nil {
			return err
		}
		defer s.Close()
		repo := NewServerRepository(s)
		return repo.Save(context.Background(), &domain.Server{
			ID: "srv-1", Name: "web-01", Host: "127.0.0.1", Port: 9100,
			Role: domain.RoleWorker, Status: domain.StatusOnline,
		})
	}
	if err := save(); err != nil {
		t.Fatal(err)
	}

	s2, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s2.Close()

	got, err := NewServerRepository(s2).GetByID(context.Background(), "srv-1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "web-01" || got.Host != "127.0.0.1" {
		t.Fatalf("unexpected server: %+v", got)
	}
}
