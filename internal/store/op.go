package store

import "time"

// OpKind is the type of change an Op records.
type OpKind uint8

const (
	OpAdd OpKind = iota
	OpComplete
	OpReopen
	OpArchive
	OpRestore
)

func (k OpKind) String() string {
	switch k {
	case OpAdd:
		return "add"
	case OpComplete:
		return "complete"
	case OpReopen:
		return "reopen"
	case OpArchive:
		return "archive"
	case OpRestore:
		return "restore"
	default:
		return "unknown"
	}
}

// Op is a single append-only change to the task list. The log of Ops is the
// source of truth; Tasks are derived by replaying it.
type Op struct {
	Kind   OpKind
	TaskID string
	Text   string // set only for OpAdd
	At     time.Time
}
