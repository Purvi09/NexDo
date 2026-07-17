package store

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
)

const headerSize = 8 // 4-byte length + 4-byte CRC32

// DB is a Store backed by an append-only file on disk. Each record is framed
// as: 4-byte payload length | 4-byte CRC32 | JSON payload.
type DB struct {
	*Store
	file *os.File
}

// Open loads the Op log at path (creating it if missing), recovers from any
// torn trailing record left by a crash, and returns a DB that appends new ops.
func Open(path string) (*DB, error) {
	if dir := filepath.Dir(path); dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}
	ops, good, err := readOps(f)
	if err != nil {
		f.Close()
		return nil, err
	}
	if err := f.Truncate(good); err != nil { // drop torn trailing bytes
		f.Close()
		return nil, err
	}
	db := &DB{Store: New(), file: f}
	for _, op := range ops {
		db.Store.Apply(op)
	}
	return db, nil
}

// Apply frames the op (length + checksum + payload), writes and fsyncs it,
// then applies it in memory.
func (db *DB) Apply(op Op) error {
	payload, err := json.Marshal(op)
	if err != nil {
		return err
	}
	rec := make([]byte, headerSize+len(payload))
	binary.BigEndian.PutUint32(rec[0:4], uint32(len(payload)))
	binary.BigEndian.PutUint32(rec[4:8], crc32.ChecksumIEEE(payload))
	copy(rec[headerSize:], payload)
	if _, err := db.file.Write(rec); err != nil {
		return err
	}
	if err := db.file.Sync(); err != nil {
		return err
	}
	db.Store.Apply(op)
	return nil
}

// Close closes the underlying file.
func (db *DB) Close() error { return db.file.Close() }

// readOps reads every intact record and returns them plus the byte offset of
// the last good record's end. A torn or corrupt final record (a crash
// mid-write) stops the read without failing — that offset is where recovery
// truncates.
func readOps(f *os.File) ([]Op, int64, error) {
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, 0, err
	}
	r := bufio.NewReader(f)
	var ops []Op
	var good int64
	for {
		var hdr [headerSize]byte
		if _, err := io.ReadFull(r, hdr[:]); err != nil {
			break // clean EOF or torn header
		}
		length := binary.BigEndian.Uint32(hdr[0:4])
		crc := binary.BigEndian.Uint32(hdr[4:8])
		payload := make([]byte, length)
		if _, err := io.ReadFull(r, payload); err != nil {
			break // torn payload
		}
		if crc32.ChecksumIEEE(payload) != crc {
			break // corrupt record
		}
		var op Op
		if err := json.Unmarshal(payload, &op); err != nil {
			break
		}
		ops = append(ops, op)
		good += headerSize + int64(length)
	}
	return ops, good, nil
}
