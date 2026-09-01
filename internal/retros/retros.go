// Package retros reads pending implementation retrospective capture artifacts.
package retros

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
)

// Directory is the repository-relative retrospective capture directory.
const Directory = "capture/retros"

var retrospectiveSlug = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)

var requiredHeadings = []string{
	"## Outcome",
	"## Gate-stage timings",
	"## Ticket-versus-spec-slice and delegate performance",
	"## Coordinator catches",
	"## Repair attribution",
	"## Agent-experience improvements",
	"### Bench CLI",
	"### Skills",
	"### Process",
}

// ValidSlug reports whether value can name one retrospective file below Directory.
func ValidSlug(value string) bool { return retrospectiveSlug.MatchString(value) }

// Path returns the repository-relative path for one retrospective slug.
func Path(value string) string { return filepath.ToSlash(filepath.Join(Directory, value+".md")) }

// Parse validates the canonical implementation-retrospective headings and every
// improvement destination marker. The writer calls it before it opens the destination.
func Parse(content []byte) error {
	lines := strings.Split(string(content), "\n")
	next := 0
	for _, line := range lines {
		line = strings.TrimSuffix(line, "\r")
		if next == len(requiredHeadings) || line != requiredHeadings[next] {
			continue
		}
		next++
	}
	if next != len(requiredHeadings) {
		return fmt.Errorf("missing or out-of-order heading %q", requiredHeadings[next])
	}
	for _, item := range Recommendations(content) {
		if !item.FeedsMarked() {
			return fmt.Errorf("improvement item on line %d carries no destination marker", item.Line)
		}
	}
	return nil
}

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

// MissingMarkerDiagnostic opens the diagnostic an improvement item without a well-formed
// destination marker raises. The repair is in the message because the retro author reads
// it in a gate report, away from the template that states the grammar.
const MissingMarkerDiagnostic = "improvement item carries no destination marker"

// ValidateImprovementMarkers grades every pending retrospective for the destination marker
// its improvement items carry, and returns one diagnostic per fault.
//
// An absent or an empty capture directory is the ordinary no-pending posture and stays
// quiet. A directory or a file the classifier refuses is a diagnostic rather than a skip:
// an unreadable retro that graded green would be an unmarked retro nobody sees.
func ValidateImprovementMarkers(root string) []string {
	facts := Facts(root)
	if len(facts.Entries) == 0 {
		if facts.State.Failed() {
			return []string{fmt.Sprintf("%s: retro capture directory is %s (%s)", Directory, facts.State, facts.Reason)}
		}
		return nil
	}
	var diags []string
	for _, entry := range facts.Entries {
		if entry.State.Failed() {
			diags = append(diags, fmt.Sprintf("%s: retro file is %s (%s)", entry.Path, entry.State, entry.Reason))
			continue
		}
		for _, item := range Recommendations(entry.Body) {
			if item.FeedsMarked() {
				continue
			}
			diags = append(diags, fmt.Sprintf(
				"%s:%d: %s; end the item with one line reading 'Feeds: FT<n>', 'Feeds: new', or 'Feeds: none'",
				entry.Path, item.Line, MissingMarkerDiagnostic))
		}
	}
	return diags
}
