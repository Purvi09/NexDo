package store

// Store is an in-memory append-only log of Ops.
type Store struct {
	ops []Op
}

// New returns an empty Store.
func New() *Store { return &Store{} }

// Apply appends an op to the log.
func (s *Store) Apply(op Op) { s.ops = append(s.ops, op) }

// Tasks materializes the current state by replaying the log.
func (s *Store) Tasks() map[string]Task { return materialize(s.ops) }

func materialize(ops []Op) map[string]Task {
	tasks := make(map[string]Task)
	for _, op := range ops {
		switch op.Kind {
		case OpAdd:
			tasks[op.TaskID] = Task{
				ID:      op.TaskID,
				Text:    op.Text,
				Created: op.At,
				Updated: op.At,
			}
		case OpComplete:
			mutate(tasks, op, func(t *Task) { t.Done = true })
		case OpReopen:
			mutate(tasks, op, func(t *Task) { t.Done = false })
		case OpArchive:
			mutate(tasks, op, func(t *Task) { t.Archived = true })
		case OpRestore:
			mutate(tasks, op, func(t *Task) { t.Archived = false })
		}
	}
	return tasks
}

func mutate(tasks map[string]Task, op Op, fn func(*Task)) {
	t, ok := tasks[op.TaskID]
	if !ok {
		return
	}
	fn(&t)
	t.Updated = op.At
	tasks[op.TaskID] = t
}
