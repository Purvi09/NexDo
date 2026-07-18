package store

import (
	"sort"
	"time"
)

// Add records a new task and returns its id.
func (db *DB) Add(text string) (string, error) {
	id := NewID()
	return id, db.Apply(Op{Kind: OpAdd, TaskID: id, Text: text, At: time.Now()})
}

func (db *DB) Complete(id string) error {
	return db.Apply(Op{Kind: OpComplete, TaskID: id, At: time.Now()})
}

func (db *DB) Reopen(id string) error {
	return db.Apply(Op{Kind: OpReopen, TaskID: id, At: time.Now()})
}

// Archive moves a task to the archive ("let go"). It's a non-destructive op.
func (db *DB) Archive(id string) error {
	return db.Apply(Op{Kind: OpArchive, TaskID: id, At: time.Now()})
}

// Restore brings a task back from the archive.
func (db *DB) Restore(id string) error {
	return db.Apply(Op{Kind: OpRestore, TaskID: id, At: time.Now()})
}

// List returns active (non-archived) tasks, oldest first.
func (db *DB) List() []Task { return filterSorted(db.Tasks(), false) }

// ListArchived returns archived tasks, oldest first.
func (db *DB) ListArchived() []Task { return filterSorted(db.Tasks(), true) }

func filterSorted(m map[string]Task, archived bool) []Task {
	out := make([]Task, 0)
	for _, t := range m {
		if t.Archived == archived {
			out = append(out, t)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Created.Equal(out[j].Created) {
			return out[i].ID < out[j].ID
		}
		return out[i].Created.Before(out[j].Created)
	})
	return out
}
