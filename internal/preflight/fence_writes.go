package preflight

import (
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/tickets"
)

// ticketFence is the ticket half of the ownership fence: the tree path that each `Writes:`
// entry of every ticket names, once, in ticket and entry order. Each entry takes the one
// `Writes:` split, so a trailing slash or a marker never becomes a path segment. The
// charge evidence fence block and the authorizing entries both read this one answer.
func ticketFence(f Facts) []string {
	var fence []string
	for _, ticket := range f.Tickets {
		for _, entry := range ticket.Writes {
			if path, _ := tickets.WritesPath(entry); !slices.Contains(fence, path) {
				fence = append(fence, path)
			}
		}
	}
	return fence
}

// implicitEntries is every entry that authorizes a path without a ticket entry: the
// phase-owned capture folder, the active spec's own folder, which authorizes an in-range
// amendment of the spec, and the exact review pickup of that spec, which no ticket names.
// The entries are derived rather than carried as their own facts, so the printed spec
// path and the authorized paths cannot disagree.
func implicitEntries(f Facts) []string {
	entries := []string{captureEntry}
	if folder := specFolder(f.SpecPath); folder != "" {
		entries = append(entries, folder)
	}
	if pickup, err := reviewrecord.RecordPath(f.SpecPath); err == nil {
		entries = append(entries, pickup)
	}
	return entries
}

// writesAnswer is the red answer that the writes-resolve and writes-landed precedence gives
// one `Writes:` entry. An entry outside both rows, or green in both, gets no answer.
type writesAnswer uint8

const (
	writesUnborn writesAnswer = iota
	writesLanded
	writesUnresolved
)

// walkWrites gives each `Writes:` entry of every ticket its one answer, in ticket and entry
// order, and visits each entry that has a red answer. Both writes rows read this one walk.
// An entry that the ticket grammar refuses is tickets-parse's red, and a (new) entry or a
// repeated entry gets no answer. Either kind of the other entries also resolves through an
// exact deletion of the build's landing range, from the landing base to the source tip. A
// (deleted) entry must name a tracked file at the review base, whatever the tip holds, so
// no marker gives authority over a path that never existed. A plain entry resolves through
// the checkout or through an exact committed deletion of the explicit source range, and
// that resolution wins over a landed deletion. A plain entry with a landed deletion is a
// writes-landed red, and every other plain entry is a writes-resolve red.
func walkWrites(f Facts, visit func(answer writesAnswer, named, entry string)) {
	seen := map[string]bool{}
	for _, ticket := range f.Tickets {
		for _, entry := range ticket.Writes {
			write, err := tickets.ParseWrite(entry)
			named := ticket.Name + ": " + entry
			switch {
			case err != nil || write.Kind == tickets.WriteNew || seen[named] || slices.Contains(f.BuildDeletedPaths, write.Path):
			case write.Kind == tickets.WriteDeleted:
				if !f.DeletedBefore[entry] {
					visit(writesUnborn, named, entry)
				}
			case f.WritesPathExists[entry] || slices.Contains(f.DeletedPaths, write.Path):
			case hasLandedDeletion(f, entry):
				visit(writesLanded, named, entry)
			default:
				visit(writesUnresolved, named, entry)
			}
			seen[named] = true
		}
	}
}

// hasLandedDeletion reports whether the gatherer recorded a deleting commit for entry.
func hasLandedDeletion(f Facts, entry string) bool {
	_, landed := f.LandedDeletions[entry]
	return landed
}

// writesResolveCheck reds each entry that walkWrites answers for writes-resolve: a plain
// entry that names no tree path and no deletion, and a (deleted) entry that names no
// tracked file at the review base.
func writesResolveCheck(f Facts) CheckResult {
	var unresolved, unborn []string
	walkWrites(f, func(answer writesAnswer, named, _ string) {
		switch answer {
		case writesUnresolved:
			unresolved = append(unresolved, named)
		case writesUnborn:
			unborn = append(unborn, named)
		}
	})
	var details []string
	if len(unresolved) > 0 {
		details = append(details, "Writes: entry names no tree path and carries no (new) marker: "+strings.Join(unresolved, ", "))
	}
	if len(unborn) > 0 {
		details = append(details, "(deleted) entry names no tracked file at the review base: "+strings.Join(unborn, ", "))
	}
	if len(details) > 0 {
		return red("writes-resolve", strings.Join(details, "; "))
	}
	return green("writes-resolve")
}

// writesLandedCheck reds each entry that walkWrites answers for writes-landed, and names
// the deleting commit of each, so the staleness pass can cite the landed work.
func writesLandedCheck(f Facts) CheckResult {
	var landed []string
	walkWrites(f, func(answer writesAnswer, named, entry string) {
		if answer == writesLanded {
			landed = append(landed, named+" ("+f.LandedDeletions[entry]+")")
		}
	})
	if len(landed) > 0 {
		return red("writes-landed", "Writes: entry names a path that a landed commit deleted: "+strings.Join(landed, ", "))
	}
	return green("writes-landed")
}
