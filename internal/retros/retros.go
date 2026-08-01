// Package retros reads pending implementation retrospective capture artifacts.
package retros

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
)

// Directory is the repository-relative retrospective capture directory.
const Directory = "capture/retros"

// Fact is one eligible retrospective file and its bounded classification.
type Fact struct {
	Path, Reason string
	State        bounds.FileState
	Body         []byte
}

// Result is the classified directory and every eligible retrospective it names.
type Result struct {
	State   bounds.FileState
	Reason  string
	Entries []Fact
}

// Facts classifies the capture directory before enumerating it, then classifies every
// non-hidden Markdown candidate before reading it. Absence and emptiness are the
// ordinary no-pending posture; every eligible entry remains evidence even if degraded.
func Facts(root string) Result {
	dir := filepath.Join(root, filepath.FromSlash(Directory))
	classified := bounds.ClassifyDir(dir)
	if classified.State == bounds.StateAbsent || classified.State == bounds.StateEmpty {
		return Result{State: bounds.StateParsed}
	}
	if classified.State != bounds.StateParsed {
		return Result{State: classified.State, Reason: classified.Reason}
	}
	entries := make([]Fact, 0, len(classified.Entries))
	for _, entry := range classified.Entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") || !strings.HasSuffix(name, ".md") {
			continue
		}
		path := filepath.Join(dir, name)
		c := bounds.Classify(path, bounds.ControlRecordLimit)
		entries = append(entries, Fact{
			Path:   filepath.ToSlash(filepath.Join(Directory, name)),
			State:  c.State,
			Reason: c.Reason,
			Body:   c.Data,
		})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Path < entries[j].Path })
	result := Result{State: bounds.StateParsed, Entries: entries}
	for _, entry := range entries {
		if entry.State == bounds.StateParsed || entry.State == bounds.StateEmpty {
			continue
		}
		result.State, result.Reason = entry.State, entry.Reason
		break
	}
	return result
}
