package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func BenchmarkApplyWithSync(b *testing.B) {
	path := filepath.Join(b.TempDir(), "bench.log")
	db, err := Open(path)
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()
	op := Op{Kind: OpAdd, TaskID: "1", Text: "buy milk", At: time.Now()}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := db.Apply(op); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWriteNoSync(b *testing.B) {
	path := filepath.Join(b.TempDir(), "nosync.log")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o644)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	line := []byte(`{"Kind":0,"TaskID":"1","Text":"buy milk"}` + "\n")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := f.Write(line); err != nil {
			b.Fatal(err)
		}
	}
}
