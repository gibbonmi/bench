package tickets

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/toon"
)

// Ext is the one extension the grammar parses. Any other entry under a tickets
// directory is an asset the grammar ignores rather than grades.
const Ext = ".md"

// Entry is one enumerated ticket file: the path a diagnostic names it by, its
// basename identity, and its bytes.
type Entry struct {
	// Rel is the slash-spelled path from the tickets directory.
	Rel string
	// Name is the basename, which is the identity a blocker edge resolves against.
	Name string
	Data []byte
}

// Refusal is the one fail-closed answer an enumeration gives when an entry will
// not read. Kind names the class and Hint names the entry and the reason.
type Refusal struct {
	Kind, Hint string
}

// Enumerate walks one already-classified tickets directory and returns the
// ticket files below it. Every venue that grades ticket grammar reads through
// this one seam, so the enumeration semantics cannot drift between them.
//
// The walk classifies each entry before it tests the extension, so a special
// file, a dangling symlink, and an unreadable entry are refused whatever they
// are named. The lstat-first refusal holds at every depth. A path spec-TOON
// cannot render is refused here, before any verdict row or diagnostic can carry
// it. Only Ext files parse as tickets.
//
// A basename that appears at two depths is a duplicate identity, because a
// blocker edge resolves by basename alone. The second copy is named in the
// returned diagnostics and dropped, so the returned set holds one entry per
// name.
func Enumerate(base string, entries []fs.DirEntry) ([]Entry, []string, *Refusal) {
	found, refusal := scan(base, base, entries)
	if refusal != nil {
		return nil, nil, refusal
	}
	unique := make([]Entry, 0, len(found))
	var diagnostics []string
	seen := make(map[string]string, len(found))
	for _, entry := range found {
		if first, duplicate := seen[entry.Name]; duplicate {
			diagnostics = append(diagnostics, fmt.Sprintf("duplicate ticket basename %s at %s and %s", entry.Name, first, entry.Rel))
			continue
		}
		seen[entry.Name] = entry.Rel
		unique = append(unique, entry)
	}
	return unique, diagnostics, nil
}

// scan collects one directory level and recurses with the same classification.
func scan(base, dir string, entries []fs.DirEntry) ([]Entry, *Refusal) {
	var files []Entry
	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			below := bounds.ClassifyDir(path)
			switch below.State {
			case bounds.StateEmpty:
				continue
			case bounds.StateParsed:
				nested, refusal := scan(base, path, below.Entries)
				if refusal != nil {
					return nil, refusal
				}
				files = append(files, nested...)
			default:
				return nil, &Refusal{"tickets directory not readable", path + " is " + string(below.State) + ": " + below.Reason}
			}
			continue
		}
		rel := filepath.ToSlash(relTo(base, path))
		if !toon.Representable(rel) {
			return nil, &Refusal{"ticket path not representable", fmt.Sprintf("ticket %q contains a byte spec-TOON cannot represent", rel)}
		}
		classified := bounds.Classify(path, bounds.ControlRecordLimit)
		unreadable := &Refusal{"ticket file not readable", path + " is " + string(classified.State) + ": " + classified.Reason}
		switch classified.State {
		case bounds.StateWrongType, bounds.StateUnreadable, bounds.StateAbsent:
			return nil, unreadable
		}
		if filepath.Ext(entry.Name()) != Ext {
			continue
		}
		switch classified.State {
		case bounds.StateParsed, bounds.StateEmpty:
			files = append(files, Entry{Rel: rel, Name: entry.Name(), Data: classified.Data})
		default:
			return nil, unreadable
		}
	}
	return files, nil
}

func relTo(base, path string) string {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return path
	}
	return rel
}

// TagOf is the uppercase alphabetic prefix of one row ID — the tag a spec
// declares and the tag the citation grammar scopes its comparison to. A
// digit-leading ID answers the empty tag, which is the tickets-only posture:
// the Covers tag rule stands down rather than grading every token as foreign.
func TagOf(rowID string) string {
	i := 0
	for i < len(rowID) && rowID[i] >= 'A' && rowID[i] <= 'Z' {
		i++
	}
	return rowID[:i]
}

// UnrepresentableValue reports the first declared field value spec-TOON cannot
// render, with the field that declared it. A blocker basename and a `Writes:`
// entry both reach a rendered detail cell, so both are refused before a verdict
// renders rather than at the render itself.
func UnrepresentableValue(ticket Ticket) (field, value string, found bool) {
	for _, blocker := range ticket.Blockers {
		if !toon.Representable(blocker) {
			return fieldBlockedBy, blocker, true
		}
	}
	for _, entry := range ticket.Writes {
		if !toon.Representable(entry) {
			return fieldWrites, entry, true
		}
	}
	return "", "", false
}
