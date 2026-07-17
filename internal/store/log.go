package store

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
)

// DB is a Store backed by an append-only file on disk.
type DB struct {
	*Store
	file *os.File
}

// Open loads the Op log at path (creating it if missing) and returns a
// persistent DB that appends new ops to the same file.
func Open(path string) (*DB, error) {
	if dir := filepath.Dir(path); dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}
	ops, err := readOps(f)
	if err != nil {
		f.Close()
		return nil, err
	}
	db := &DB{Store: New(), file: f}
	for _, op := range ops {
		db.Store.Apply(op) // replay into memory without re-writing to disk
	}
	return db, nil
}

// Apply writes the op to disk, then applies it in memory.
func (db *DB) Apply(op Op) error {
	line, err := json.Marshal(op)
	if err != nil {
		return err
	}
	if _, err := db.file.Write(append(line, '\n')); err != nil {
		return err
	}
	if err := db.file.Sync(); err != nil {
		return err
	}
	db.Store.Apply(op)
	return nil
}

// Close closes the underlying file.
func (db *DB) Close() error { return db.file.Close() }

func readOps(f *os.File) ([]Op, error) {
	var ops []Op
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if len(sc.Bytes()) == 0 {
			continue
		}
		var op Op
		if err := json.Unmarshal(sc.Bytes(), &op); err != nil {
			return nil, err
		}
		ops = append(ops, op)
	}
	return ops, sc.Err()
}
