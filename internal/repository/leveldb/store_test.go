package leveldb

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStore_PutGetDelete_OnDisk(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "db")

	s, err := Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}

	type sample struct {
		Name string `json:"name"`
	}
	if err := s.putJSON("test:key", &sample{Name: "web-01"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}

	// Повторное открытие — данные на диске (LevelDB)
	s2, err := Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer s2.Close()

	var got sample
	if err := s2.getJSON("test:key", &got); err != nil {
		t.Fatal(err)
	}
	if got.Name != "web-01" {
		t.Fatalf("want web-01, got %s", got.Name)
	}

	if err := s2.delete("test:key"); err != nil {
		t.Fatal(err)
	}
}

func TestStore_DataDirectoryContainsLevelDBFiles(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "db")

	s, err := Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.putJSON("ping", "pong"); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}

	entries, err := os.ReadDir(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("expected LevelDB files in data directory")
	}
}

func TestStore_ListByPrefix(t *testing.T) {
	dir := t.TempDir()
	s, err := Open(filepath.Join(dir, "db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	_ = s.putJSON(prefixServer+"a", map[string]string{"id": "a"})
	_ = s.putJSON(prefixServer+"b", map[string]string{"id": "b"})

	items, err := s.listByPrefix(prefixServer)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("want 2 items, got %d", len(items))
	}
}
