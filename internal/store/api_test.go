package store

import (
	"path/filepath"
	"testing"
)

func TestStoreAPI(t *testing.T) {
	path := filepath.Join(t.TempDir(), "api.log")
	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	id, err := db.Add("buy milk")
	must(t, err)
	if len(db.List()) != 1 {
		t.Fatalf("want 1 active task, got %d", len(db.List()))
	}

	must(t, db.Complete(id))
	if !db.List()[0].Done {
		t.Fatal("task should be done after Complete")
	}

	must(t, db.Archive(id))
	if len(db.List()) != 0 {
		t.Fatal("archived task should leave the active list")
	}
	if len(db.ListArchived()) != 1 {
		t.Fatal("archived task should appear in the archive")
	}

	must(t, db.Restore(id))
	if len(db.List()) != 1 {
		t.Fatal("restored task should return to the active list")
	}
	if len(db.ListArchived()) != 0 {
		t.Fatal("archive should be empty after Restore")
	}
}
