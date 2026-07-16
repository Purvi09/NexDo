package store

import (
	"path/filepath"
	"testing"
	"time"
)

func TestPersistenceRoundTrip(t *testing.T) {
	t0 := time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC)
	path := filepath.Join(t.TempDir(), "test.log")

	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	must(t, db.Apply(Op{Kind: OpAdd, TaskID: "1", Text: "survive a restart", At: t0}))
	must(t, db.Apply(Op{Kind: OpComplete, TaskID: "1", At: t0}))
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	// Reopen from the same file — the state must still be there.
	db2, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db2.Close()

	got := db2.Tasks()["1"]
	if got.Text != "survive a restart" || !got.Done {
		t.Fatalf("state did not persist: %+v", got)
	}
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
