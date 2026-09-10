package roadmap

import (
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/otelrecord"
	"github.com/gibbonmi/bench/internal/retros"
	"github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/tickets"
	"github.com/gibbonmi/bench/internal/toon"
)

// unknownFact is what the scaffold states where the tree holds no derived fact. The
// delegate that fills the draft reads the word and supplies the fact; an invented value
// would read as measured evidence.
const unknownFact = "unknown"

// retroScaffold prints the retrospective draft for slug and writes nothing. The capture
// call owns the exclusive create, so a scaffold that wrote the file would make that call
// refuse on its own draft.
func retroScaffold(slug string) (string, int) {
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	return scaffoldBody(benchhome.Dir(), root, slug), 0
}

// scaffoldBody renders the draft from the retrospective parser's own heading list, so
// every draft the scaffold prints is a draft the capture call accepts.
func scaffoldBody(home, root, slug string) string {
	var out strings.Builder
	for _, heading := range retros.RequiredHeadings() {
		out.WriteString(heading + "\n\n")
		if section := scaffoldSection(home, root, slug, heading); section != "" {
			out.WriteString(section + "\n\n")
		}
	}
	return strings.TrimRight(out.String(), "\n") + "\n"
}

// scaffoldSection returns the derived body of one heading. Every other heading stays
// empty, because the delegate writes it from the run it just watched.
func scaffoldSection(home, root, slug, heading string) string {
	switch heading {
	case retros.TimingsHeading:
		return timingsSection(home, root)
	case retros.RepairHeading:
		return repairSection(root, slug)
	}
	return ""
}

// timingsSection lists the gate stages of the landing's own trace. It names the
// published commit and the trace, so a reader verifies which landing was selected. A
// record that is absent, is unreadable, or names no landing states unknown instead: a
// missing measurement must not block the draft.
func timingsSection(home, root string) string {
	landing, ok := otelrecord.NewestLanding(home, root)
	if !ok {
		return unknownFact
	}
	lines := []string{"- landing: commit " + orUnknown(landing.Commit) + ", trace " + orUnknown(landing.TraceID)}
	for _, stage := range landing.Stages {
		lines = append(lines, "- "+stage.Name+": "+strconv.FormatInt(stage.Elapsed().Milliseconds(), 10)+" ms")
	}
	return strings.Join(lines, "\n")
}

// repairSection prints the repair-attribution table with one row per ticket of the
// slug. The rounds cell and the cause cell stay unknown, because a repair count and a
// cause are read from the build, not from the tree.
func repairSection(root, slug string) string {
	rows := []string{"| ticket | rounds | causes |", "|---|---|---|"}
	for _, ticket := range scaffoldTickets(root, slug) {
		rows = append(rows, "| "+ticket+" | "+unknownFact+" | "+unknownFact+" |")
	}
	return strings.Join(rows, "\n")
}

// scaffoldTickets returns the slug's ticket basenames in name order. An absent, an
// unreadable, and a refused tickets directory each give the one unknown row, so a
// retired folder still leaves the delegate the table's shape.
func scaffoldTickets(root, slug string) []string {
	dir := filepath.Join(root, filepath.FromSlash(filepath.Dir(spec.LiveSpecPath(slug))), "tickets")
	classified := bounds.ClassifyDirNoFollow(dir)
	if classified.State != bounds.StateParsed {
		return []string{unknownFact}
	}
	entries, _, refusal := tickets.EnumerateNoFollow(dir, classified.Entries)
	if refusal != nil || len(entries) == 0 {
		return []string{unknownFact}
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name)
	}
	sort.Strings(names)
	return names
}

func orUnknown(value string) string {
	if value == "" {
		return unknownFact
	}
	return value
}
