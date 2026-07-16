package store

import (
	"crypto/rand"
	"encoding/hex"
)

// NewID returns a random hex id for a new task.
func NewID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
