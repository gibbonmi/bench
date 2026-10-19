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

// writesResolveCheck grades each `Writes:` entry against the tree. An entry without a
// marker must name a path the tree holds, or an exact committed deletion of the explicit
// source range. A (deleted) entry must name a tracked file at the review base, whatever
// the tip holds, so no marker gives authority over a path that never existed. Either
// entry also resolves through an exact deletion of the build's landing range, from the
// landing base to the source tip. An entry that the ticket grammar refuses is
// tickets-parse's red.
func writesResolveCheck(f Facts) CheckResult {
	var unresolved, unborn []string
	seen := map[string]bool{}
	for _, ticket := range f.Tickets {
		for _, entry := range ticket.Writes {
			write, err := tickets.ParseWrite(entry)
			named := ticket.Name + ": " + entry
			switch {
			case err != nil || write.Kind == tickets.WriteNew || seen[named] || slices.Contains(f.BuildDeletedPaths, write.Path):
			case write.Kind == tickets.WriteDeleted:
				if !f.DeletedBefore[entry] {
					unborn = append(unborn, named)
				}
			case !f.WritesPathExists[entry] && !slices.Contains(f.DeletedPaths, write.Path):
				unresolved = append(unresolved, named)
			}
			seen[named] = true
		}
	}
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
