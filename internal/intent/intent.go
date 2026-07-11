// Package intent persists machine-written work intent in the repository's shared
// git common directory. It is the only package that knows the ledger's address,
// schema, locking, atomic replacement, and proof-based lifecycle.
package intent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"
	"unicode"

	"github.com/gibbonmi/bench/internal/git"
)

const (
	Schema   = 1
	Filename = "bench-intent.json"
)

type Kind string

const (
	KindShift       Kind = "shift"
	KindWorktree    Kind = "worktree"
	KindClaudeAgent Kind = "claude-agent"
)

type Entry struct {
	Key       string    `json:"key"`
	Kind      Kind      `json:"kind"`
	Objective string    `json:"objective"`
	CreatedAt time.Time `json:"created_at"`
	Worktree  string    `json:"worktree,omitempty"`
	Branch    string    `json:"branch,omitempty"`
}

type Ledger struct {
	Schema  int     `json:"schema"`
	Entries []Entry `json:"entries"`
}

// Address resolves the ledger through git's absolute common-directory query, so
// the primary checkout and every linked worktree share one file.
func Address(root string) (string, error) {
	common, err := git.Output("-C", root, "rev-parse", "--path-format=absolute", "--git-common-dir")
	if err != nil || common == "" {
		return "", fmt.Errorf("resolve git common directory: %w", err)
	}
	return filepath.Join(common, Filename), nil
}

func Read(root string) (Ledger, error) {
	path, err := Address(root)
	if err != nil {
		return Ledger{}, err
	}
	return readPath(path)
}

func readPath(path string) (Ledger, error) {
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return Ledger{Schema: Schema, Entries: []Entry{}}, nil
	}
	if err != nil {
		return Ledger{}, fmt.Errorf("read intent ledger: %w", err)
	}
	if info.Mode().Perm()&0o444 == 0 {
		return Ledger{}, errors.New("read intent ledger: file is unreadable")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Ledger{}, fmt.Errorf("read intent ledger: %w", err)
	}
	if len(data) == 0 {
		return Ledger{}, errors.New("read intent ledger: present file is empty")
	}
	var ledger Ledger
	dec := json.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&ledger); err != nil {
		return Ledger{}, fmt.Errorf("read intent ledger: malformed JSON: %w", err)
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return Ledger{}, errors.New("read intent ledger: trailing JSON value")
	}
	if ledger.Schema != Schema {
		return Ledger{}, fmt.Errorf("read intent ledger: unsupported schema %d", ledger.Schema)
	}
	if ledger.Entries == nil {
		ledger.Entries = []Entry{}
	}
	seen := map[string]bool{}
	for _, entry := range ledger.Entries {
		if err := validEntry(entry); err != nil {
			return Ledger{}, fmt.Errorf("read intent ledger: %w", err)
		}
		if seen[entry.Key] {
			return Ledger{}, fmt.Errorf("read intent ledger: duplicate key %q", entry.Key)
		}
		seen[entry.Key] = true
	}
	return ledger, nil
}

func validEntry(entry Entry) error {
	if entry.Key == "" || entry.Objective == "" || entry.CreatedAt.IsZero() {
		return errors.New("entry requires key, objective, and creation time")
	}
	switch entry.Kind {
	case KindShift, KindWorktree, KindClaudeAgent:
		return nil
	default:
		return fmt.Errorf("entry %q has unknown writer kind %q", entry.Key, entry.Kind)
	}
}

// NewEntry creates the stable process/time key bench-owned writers persist before
// acquire. Callers retain the returned value and enrich that same key later.
func NewEntry(kind Kind, objective string) Entry {
	now := time.Now().UTC()
	return Entry{
		Key:       fmt.Sprintf("%s-%d-%d", kind, os.Getpid(), now.UnixNano()),
		Kind:      kind,
		Objective: objective,
		CreatedAt: now,
	}
}

// Upsert inserts or enriches one stable writer key. An identical upsert does not
// replace the file, preserving exact bytes and mtime.
func Upsert(root string, entry Entry) error {
	if err := validEntry(entry); err != nil {
		return err
	}
	path, err := Address(root)
	if err != nil {
		return err
	}
	release, err := acquire(path + ".lock")
	if err != nil {
		return err
	}
	defer release()
	ledger, err := readPath(path)
	if err != nil {
		return err
	}
	changed := true
	for i := range ledger.Entries {
		if ledger.Entries[i].Key != entry.Key {
			continue
		}
		if entry.CreatedAt.IsZero() {
			entry.CreatedAt = ledger.Entries[i].CreatedAt
		}
		if ledger.Entries[i] == entry {
			return nil
		}
		ledger.Entries[i] = entry
		changed = true
		goto write
	}
	ledger.Entries = append(ledger.Entries, entry)

write:
	if !changed {
		return nil
	}
	return writePath(path, ledger)
}

func writePath(path string, ledger Ledger) error {
	ledger.Schema = Schema
	if ledger.Entries == nil {
		ledger.Entries = []Entry{}
	}
	sort.Slice(ledger.Entries, func(i, j int) bool { return ledger.Entries[i].Key < ledger.Entries[j].Key })
	data, err := json.Marshal(ledger)
	if err != nil {
		return fmt.Errorf("encode intent ledger: %w", err)
	}
	data = append(data, '\n')
	tmp, err := os.CreateTemp(filepath.Dir(path), ".bench-intent-*")
	if err != nil {
		return fmt.Errorf("create intent ledger temporary file: %w", err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if err := tmp.Chmod(0o600); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return fmt.Errorf("write intent ledger: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return fmt.Errorf("sync intent ledger: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close intent ledger: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("replace intent ledger: %w", err)
	}
	if dir, err := os.Open(filepath.Dir(path)); err == nil {
		_ = dir.Sync()
		_ = dir.Close()
	}
	return nil
}

func acquire(lock string) (func(), error) {
	deadline := time.Now().Add(2 * time.Second)
	for {
		file, err := os.OpenFile(lock, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
		if err == nil {
			owner := fmt.Sprintf("%d %d\n", os.Getpid(), time.Now().Unix())
			if _, err := file.WriteString(owner); err != nil {
				_ = file.Close()
				_ = os.Remove(lock)
				return nil, err
			}
			if err := file.Sync(); err != nil {
				_ = file.Close()
				_ = os.Remove(lock)
				return nil, err
			}
			if err := file.Close(); err != nil {
				_ = os.Remove(lock)
				return nil, err
			}
			return func() { _ = os.Remove(lock) }, nil
		} else if !errors.Is(err, os.ErrExist) {
			return nil, fmt.Errorf("lock intent ledger: %w", err)
		}
		if staleLock(lock) {
			_ = os.Remove(lock)
			continue
		}
		if time.Now().After(deadline) {
			return nil, errors.New("lock intent ledger: timed out")
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func staleLock(lock string) bool {
	data, err := os.ReadFile(lock)
	if err != nil {
		return false
	}
	var pid int
	var created int64
	if _, err := fmt.Sscanf(string(data), "%d %d", &pid, &created); err != nil {
		info, statErr := os.Stat(lock)
		return statErr == nil && time.Since(info.ModTime()) > 100*time.Millisecond
	}
	if time.Since(time.Unix(created, 0)) > 30*time.Second {
		return true
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		return true
	}
	return process.Signal(syscall.Signal(0)) != nil
}

// Snapshot returns proof-live entries without mutating the ledger.
func Snapshot(root string) ([]Entry, error) {
	ledger, err := Read(root)
	if err != nil {
		return nil, err
	}
	return filterLive(root, ledger.Entries)
}

func filterLive(root string, entries []Entry) ([]Entry, error) {
	candidates, err := claudeCandidates(root)
	if err != nil {
		return nil, err
	}
	def, defOK := git.ResolvedDefault(root)
	live := make([]Entry, 0, len(entries))
	for _, entry := range entries {
		if entry.Kind == KindClaudeAgent && entry.Worktree == "" && entry.Branch == "" {
			if candidates {
				live = append(live, entry)
			}
			continue
		}
		if entry.Worktree != "" {
			if _, err := os.Stat(entry.Worktree); errors.Is(err, os.ErrNotExist) {
				continue
			}
		}
		if entry.Branch != "" && defOK {
			if landed, _ := git.LandedInDefault(root, entry.Branch, def); landed {
				continue
			}
		}
		live = append(live, entry)
	}
	sort.Slice(live, func(i, j int) bool {
		if live[i].CreatedAt.Equal(live[j].CreatedAt) {
			return live[i].Key < live[j].Key
		}
		return live[i].CreatedAt.Before(live[j].CreatedAt)
	})
	return live, nil
}

func claudeCandidates(root string) (bool, error) {
	worktrees, err := git.Output("-C", root, "worktree", "list", "--porcelain")
	if err != nil {
		return false, fmt.Errorf("classify intent worktrees: %w", err)
	}
	for _, line := range strings.Split(worktrees, "\n") {
		if strings.HasPrefix(line, "worktree ") && strings.HasPrefix(filepath.Base(strings.TrimPrefix(line, "worktree ")), "worktree-agent-") {
			return true, nil
		}
	}
	branches, err := git.Output("-C", root, "for-each-ref", "--format=%(refname:short)", "refs/heads/")
	if err != nil {
		return false, fmt.Errorf("classify intent branches: %w", err)
	}
	for _, branch := range strings.Split(branches, "\n") {
		if strings.HasPrefix(branch, "worktree-agent-") {
			return true, nil
		}
	}
	return false, nil
}

// Compact atomically removes only entries Snapshot has proven done.
func Compact(root string) error {
	path, err := Address(root)
	if err != nil {
		return err
	}
	release, err := acquire(path + ".lock")
	if err != nil {
		return err
	}
	defer release()
	current, err := readPath(path)
	if err != nil {
		return err
	}
	live, err := filterLive(root, current.Entries)
	if err != nil {
		return err
	}
	liveKeys := map[string]bool{}
	for _, entry := range live {
		liveKeys[entry.Key] = true
	}
	next := current.Entries[:0]
	for _, entry := range current.Entries {
		if liveKeys[entry.Key] {
			next = append(next, entry)
		}
	}
	if len(next) == len(current.Entries) {
		return nil
	}
	current.Entries = next
	return writePath(path, current)
}

// Preview safely encodes terminal context and caps it at 120 Unicode code points.
func Preview(value string) string {
	runes := []rune(value)
	truncated := len(runes) > 120
	if truncated {
		runes = runes[:120]
	}
	var b strings.Builder
	for _, r := range runes {
		switch {
		case r == '\n':
			b.WriteString(`\n`)
		case r == '\r':
			b.WriteString(`\r`)
		case r == '\t':
			b.WriteString(`\t`)
		case unicode.IsControl(r):
			fmt.Fprintf(&b, "\\u%04x", r)
		case r == '\\':
			b.WriteString(`\\`)
		default:
			b.WriteRune(r)
		}
	}
	if truncated {
		fmt.Fprintf(&b, "… (%d bytes)", len(value))
	}
	return b.String()
}
