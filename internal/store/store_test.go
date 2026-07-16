package store

import (
	"reflect"
	"testing"
	"time"
)

func TestMaterialize(t *testing.T) {
	t0 := time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC)
	t1 := t0.Add(time.Hour)

	tests := []struct {
		name string
		ops  []Op
		want map[string]Task
	}{
		{
			name: "add creates a task",
			ops:  []Op{{Kind: OpAdd, TaskID: "1", Text: "buy milk", At: t0}},
			want: map[string]Task{
				"1": {ID: "1", Text: "buy milk", Created: t0, Updated: t0},
			},
		},
		{
			name: "complete sets done and bumps Updated",
			ops: []Op{
				{Kind: OpAdd, TaskID: "1", Text: "buy milk", At: t0},
				{Kind: OpComplete, TaskID: "1", At: t1},
			},
			want: map[string]Task{
				"1": {ID: "1", Text: "buy milk", Done: true, Created: t0, Updated: t1},
			},
		},
		{
			name: "reopen clears done",
			ops: []Op{
				{Kind: OpAdd, TaskID: "1", Text: "buy milk", At: t0},
				{Kind: OpComplete, TaskID: "1", At: t0},
				{Kind: OpReopen, TaskID: "1", At: t1},
			},
			want: map[string]Task{
				"1": {ID: "1", Text: "buy milk", Done: false, Created: t0, Updated: t1},
			},
		},
		{
			name: "archive then restore",
			ops: []Op{
				{Kind: OpAdd, TaskID: "1", Text: "buy milk", At: t0},
				{Kind: OpArchive, TaskID: "1", At: t0},
				{Kind: OpRestore, TaskID: "1", At: t1},
			},
			want: map[string]Task{
				"1": {ID: "1", Text: "buy milk", Archived: false, Created: t0, Updated: t1},
			},
		},
		{
			name: "op for unknown task is ignored",
			ops: []Op{
				{Kind: OpComplete, TaskID: "ghost", At: t0},
			},
			want: map[string]Task{},
		},
		{
			name: "two independent tasks",
			ops: []Op{
				{Kind: OpAdd, TaskID: "1", Text: "buy milk", At: t0},
				{Kind: OpAdd, TaskID: "2", Text: "call mom", At: t0},
				{Kind: OpComplete, TaskID: "1", At: t1},
			},
			want: map[string]Task{
				"1": {ID: "1", Text: "buy milk", Done: true, Created: t0, Updated: t1},
				"2": {ID: "2", Text: "call mom", Created: t0, Updated: t0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := materialize(tt.ops)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("materialize() =\n  %v\nwant\n  %v", got, tt.want)
			}
		})
	}
}

func TestStoreApplyThenTasks(t *testing.T) {
	t0 := time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC)
	s := New()
	s.Apply(Op{Kind: OpAdd, TaskID: "1", Text: "ship day 3", At: t0})
	s.Apply(Op{Kind: OpComplete, TaskID: "1", At: t0})

	tasks := s.Tasks()
	if got := tasks["1"]; !got.Done || got.Text != "ship day 3" {
		t.Fatalf("unexpected task: %+v", got)
	}
}
