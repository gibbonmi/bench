// Package intent persists machine-written work intent in the repository's shared
// git common directory. It is the only package that knows the ledger's address,
// schema, locking, atomic replacement, and proof-based lifecycle.
package intent

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/jsonfile"
)

const (
	LegacySchema = 1
	Schema       = 2
	Filename     = "bench-intent.json"
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
	Schema          int              `json:"schema"`
	Entries         []Entry          `json:"entries"`
	Assignments     []Assignment     `json:"assignments,omitempty"`
	CleanupReceipts []CleanupReceipt `json:"cleanup_receipts,omitempty"`
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
	if err := jsonfile.Decode(data, &ledger); err != nil {
		return Ledger{}, fmt.Errorf("read intent ledger: malformed JSON: %w", err)
	}
	if ledger.Schema != LegacySchema && ledger.Schema != Schema {
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
	if ledger.Assignments == nil {
		ledger.Assignments = []Assignment{}
	}
	if ledger.Schema == LegacySchema && (len(ledger.Assignments) != 0 || len(ledger.CleanupReceipts) != 0) {
		return Ledger{}, errors.New("read intent ledger: legacy schema cannot authorize lifecycle records")
	}
	assignmentIDs := map[string]bool{}
	requests := map[string]bool{}
	for _, assignment := range ledger.Assignments {
		if err := ValidateAssignment(assignment); err != nil {
			return Ledger{}, fmt.Errorf("read intent ledger: %w", err)
		}
		if assignmentIDs[assignment.ID] || requests[assignment.Request] {
			return Ledger{}, errors.New("read intent ledger: duplicate assignment identity")
		}
		assignmentIDs[assignment.ID] = true
		requests[assignment.Request] = true
	}
	if err := validateCleanupReceipts(ledger.CleanupReceipts); err != nil {
		return Ledger{}, fmt.Errorf("read intent ledger: %w", err)
	}
	return ledger, nil
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
		entry.CreatedAt = ledger.Entries[i].CreatedAt
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
	if ledger.Assignments == nil {
		ledger.Assignments = []Assignment{}
	}
	if ledger.CleanupReceipts == nil {
		ledger.CleanupReceipts = []CleanupReceipt{}
	}
	sort.Slice(ledger.Entries, func(i, j int) bool { return ledger.Entries[i].Key < ledger.Entries[j].Key })
	sort.Slice(ledger.Assignments, func(i, j int) bool { return ledger.Assignments[i].ID < ledger.Assignments[j].ID })
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
			if landed, _, err := git.LandedInDefault(root, entry.Branch, def); err == nil && landed {
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
	worktrees, err := git.Worktrees(root)
	if err != nil {
		return false, fmt.Errorf("classify intent worktrees: %w", err)
	}
	for _, worktree := range worktrees {
		if strings.HasPrefix(filepath.Base(worktree.Path), "worktree-agent-") {
			return true, nil
		}
	}
	branches, err := git.LocalBranches(root)
	if err != nil {
		return false, fmt.Errorf("classify intent branches: %w", err)
	}
	for _, branch := range branches {
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
