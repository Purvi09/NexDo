// Command nexdo runs the NexDo local-first to-do engine.
package main

import (
	"fmt"
	"log"
	"net/http"

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
		active := db.List()
		fmt.Fprintf(w, "NexDo — %d active task(s)\n\n", len(active))
		for _, t := range active {
			box := "[ ]"
			if t.Done {
				box = "[x]"
			}
			fmt.Fprintf(w, "%s  %s  %s\n", t.ID, box, t.Text)
		}
		fmt.Fprintf(w, "\n(%d archived — see /archived)\n", len(db.ListArchived()))
	})

	mux.HandleFunc("GET /archived", func(w http.ResponseWriter, r *http.Request) {
		archived := db.ListArchived()
		fmt.Fprintf(w, "NexDo — %d archived task(s)\n\n", len(archived))
		for _, t := range archived {
			fmt.Fprintf(w, "%s  %s\n", t.ID, t.Text)
		}
	})

	// Temporary demo endpoints so the store API is visible from the browser.
	// The real JSON API + UI come in Phase 2.
	mux.HandleFunc("GET /add", demo(func(r *http.Request) error {
		text := r.URL.Query().Get("text")
		if text == "" {
			return fmt.Errorf("missing ?text=")
		}
		_, err := db.Add(text)
		return err
	}))
	mux.HandleFunc("GET /complete", demo(func(r *http.Request) error {
		return db.Complete(r.URL.Query().Get("id"))
	}))
	mux.HandleFunc("GET /archive", demo(func(r *http.Request) error {
		return db.Archive(r.URL.Query().Get("id"))
	}))
	mux.HandleFunc("GET /restore", demo(func(r *http.Request) error {
		return db.Restore(r.URL.Query().Get("id"))
	}))

	log.Printf("NexDo listening on http://%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

func demo(action func(*http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := action(r); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
