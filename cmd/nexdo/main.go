// Command nexdo runs the NexDo local-first to-do engine.
package main

import (
	"fmt"
	"log"
	"net/http"
)

const addr = "localhost:7777"

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "NexDo is running 🌱")
	})

	log.Printf("NexDo listening on http://%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
