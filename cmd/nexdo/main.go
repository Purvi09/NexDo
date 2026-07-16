// Command nexdo runs the NexDo local-first to-do engine.
package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/Purvi09/NexDo/internal/store"
)

const addr = "localhost:7777"

func main() {
	db, err := store.Open("data/nexdo.log")
	if err != nil {
		log.Fatalf("open store: %v", err)
	}
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		tasks := sortedTasks(db.Tasks())
		fmt.Fprintf(w, "NexDo — %d task(s)\n\n", len(tasks))
		for _, t := range tasks {
			box := "[ ]"
			if t.Done {
				box = "[x]"
			}
			fmt.Fprintf(w, "%s %s\n", box, t.Text)
		}
	})

	// Temporary demo endpoint so persistence is visible from the browser.
	// A proper API + UI comes later.
	mux.HandleFunc("GET /add", func(w http.ResponseWriter, r *http.Request) {
		text := r.URL.Query().Get("text")
		if text == "" {
			http.Error(w, "missing ?text=", http.StatusBadRequest)
			return
		}
		op := store.Op{Kind: store.OpAdd, TaskID: store.NewID(), Text: text, At: time.Now()}
		if err := db.Apply(op); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	log.Printf("NexDo listening on http://%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

func sortedTasks(m map[string]store.Task) []store.Task {
	out := make([]store.Task, 0, len(m))
	for _, t := range m {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Created.Equal(out[j].Created) {
			return out[i].ID < out[j].ID
		}
		return out[i].Created.Before(out[j].Created)
	})
	return out
}
