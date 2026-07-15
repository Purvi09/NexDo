// Command NexDo runs the NexDo local-first to-do engine.
//
// Day 1: it does nothing useful yet — it just starts an HTTP server so we
// have a running process to build on. Everything (the storage engine, the
// API, the UI) will hang off this server over the coming days.
package main

import (
	"fmt"
	"log"
	"net/http"
)

// addr is where the engine listens. localhost only — this is a local-first
// app, so by default nothing is exposed to the network.
const addr = "localhost:7777"

func main() {
	// A ServeMux is Go's HTTP request router: it maps URL patterns to handlers.
	mux := http.NewServeMux()

	// "GET /" matches GET requests to the root path. The method prefix is a
	// Go 1.22+ feature so we don't have to check r.Method by hand.
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "NexDo is running 🌱")
	})

	log.Printf("NexDo listening on http://%s", addr)

	// ListenAndServe blocks forever, serving requests, until it errors.
	// It only returns on failure, so if we get past it, something went wrong.
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
