package store

import "time"

// Task is the derived state of a to-do item, materialized from the Op log.
type Task struct {
	ID       string
	Text     string
	Done     bool
	Archived bool
	Created  time.Time
	Updated  time.Time
}
