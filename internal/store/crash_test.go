package store

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCrashRecovery(t *testing.T) {
	t0 := time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC)
	path := filepath.Join(t.TempDir(), "crash.log")

	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	must(t, db.Apply(Op{Kind: OpAdd, TaskID: "1", Text: "safe task", At: t0}))
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	// Simulate a crash mid-write: a full header claiming a 50-byte payload,
	// but only 4 payload bytes actually reached the disk.
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	var hdr [8]byte
	binary.BigEndian.PutUint32(hdr[0:4], 50)
	f.Write(hdr[:])
	f.Write([]byte("part"))
	f.Close()

	// Reopen: the good record survives, the torn one is dropped.
	db2, err := Open(path)
	if err != nil {
		t.Fatalf("recovery failed: %v", err)
	}
	if n := len(db2.Tasks()); n != 1 || db2.Tasks()["1"].Text != "safe task" {
		t.Fatalf("recovery lost data: %d tasks", n)
	}
	must(t, db2.Apply(Op{Kind: OpAdd, TaskID: "2", Text: "after recovery", At: t0}))
	if err := db2.Close(); err != nil {
		t.Fatal(err)
	}

	// Reopen once more: torn bytes stay gone, both tasks present.
	db3, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db3.Close()
	if n := len(db3.Tasks()); n != 2 {
		t.Fatalf("recovery not durable: want 2 tasks, got %d", n)
	}
}
