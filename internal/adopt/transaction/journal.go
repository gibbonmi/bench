package transaction

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gibbonmi/bench/internal/bounds"
)

type entry struct {
	Span        *Span  `json:"span,omitempty"`
	Destination string `json:"destination"`
	Before      image  `json:"before"`
	After       image  `json:"after"`
}

type recordState string

const (
	statePrepared        recordState = "prepared"
	stateApplied         recordState = "applied"
	stateUndoing         recordState = "undoing"
	stateUndone          recordState = "undone"
	privateDirectoryMode os.FileMode = 0o700
	privateRecordMode    os.FileMode = 0o600
)

type record struct {
	Version int         `json:"version"`
	State   recordState `json:"state"`
	Entries []entry     `json:"entries"`
}

func (s Store) prepare(entries []entry) (string, record, error) {
	r := record{Version: 1, State: statePrepared, Entries: entries}
	if err := ensurePrivateDirectory(s.Directory, s.syncDirectory); err != nil {
		return "", r, err
	}
	var nonce [16]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return "", r, err
	}
	id := hex.EncodeToString(nonce[:])
	dir := filepath.Join(s.Directory, id)
	if err := os.Mkdir(dir, privateDirectoryMode); err != nil {
		return "", r, err
	}
	for i, e := range entries {
		if e.Before.Kind != "file" {
			continue
		}
		backup := image{Kind: "file", Mode: privateRecordMode, Data: backupData(e)}
		if err := s.publish(filepath.Join(dir, strconv.Itoa(i)+".before"), backup); err != nil {
			return id, r, err
		}
	}
	if err := s.save(id, r); err != nil {
		return id, r, err
	}
	if err := s.syncDirectory(s.Directory); err != nil {
		return id, r, err
	}
	return id, r, nil
}

func (s Store) save(id string, r record) error {
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}
	if int64(len(data)) > bounds.ControlRecordLimit {
		return fmt.Errorf("repair metadata exceeds control-record limit")
	}
	return s.publish(filepath.Join(s.Directory, id, "record.json"), image{Kind: "file", Mode: privateRecordMode, Data: data})
}

func (s Store) load(id string) (record, error) {
	var r record
	decoded, err := hex.DecodeString(id)
	if err != nil || len(decoded) != 16 || hex.EncodeToString(decoded) != id {
		return r, fmt.Errorf("invalid repair identifier")
	}
	dir := filepath.Join(s.Directory, id)
	if err := bounds.RefuseLinks(dir); err != nil {
		return r, err
	}
	for _, path := range []string{s.Directory, dir} {
		if err := requireMode(path, os.ModeDir|privateDirectoryMode); err != nil {
			return r, err
		}
	}
	metadata := filepath.Join(dir, "record.json")
	if err := requireMode(metadata, privateRecordMode); err != nil {
		return r, err
	}
	got := bounds.ClassifyNoFollow(metadata)
	if got.State != bounds.StateParsed {
		return r, fmt.Errorf("repair record %s: %s: %s", id, got.State, got.Reason)
	}
	if err := json.Unmarshal(got.Data, &r); err != nil {
		return r, fmt.Errorf("malformed repair record: %w", err)
	}
	if r.Version != 1 {
		return r, fmt.Errorf("unsupported repair record version")
	}
	switch r.State {
	case statePrepared, stateApplied, stateUndoing, stateUndone:
	default:
		return r, fmt.Errorf("unsupported repair record state")
	}
	for i := range r.Entries {
		e := &r.Entries[i]
		if !filepath.IsAbs(e.Destination) {
			return r, fmt.Errorf("invalid repair destination")
		}
		if e.Before.Kind != "file" {
			continue
		}
		path := filepath.Join(dir, strconv.Itoa(i)+".before")
		if err := bounds.RefuseLinks(path); err != nil {
			return r, err
		}
		if err := requireMode(path, privateRecordMode); err != nil {
			return r, err
		}
		backup, err := inspect(path)
		if err != nil {
			return r, err
		}
		digest := e.Before.Digest
		if e.Span != nil {
			digest = e.Span.Digest
		}
		if backup.Kind != "file" || backup.Digest != digest {
			return r, fmt.Errorf("repair backup mismatch: %s", path)
		}
		e.Before.Data = backup.Data
	}
	return r, nil
}

func requireMode(path string, want os.FileMode) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if info.Mode() != want {
		return fmt.Errorf("repair path has unsafe mode %v; want %v: %s", info.Mode(), want, path)
	}
	return nil
}

// EnsurePrivateDirectory creates a private namespace or refuses its existing unsafe mode.
func EnsurePrivateDirectory(path string) error {
	return ensurePrivateDirectory(path, SyncDirectory)
}

func ensurePrivateDirectory(path string, persist func(string) error) error {
	if err := bounds.RefuseLinks(path); err != nil {
		return err
	}
	if err := ensureDirectory(path, privateDirectoryMode, persist); err != nil {
		return err
	}
	return requireMode(path, os.ModeDir|privateDirectoryMode)
}
