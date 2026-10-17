package preflight

import (
	"slices"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/tickets"
)

// fenceWritesCheck grades the spec fence against the union of the ticket `Writes:`
// paths. A fence entry no ticket owns grants write authority nobody planned, and a
// backticked prose token in the fence section parses as exactly such an entry. A
// `Writes:` path the fence omits fails at the first commit that touches it.
//
// Both sides take the one `Writes:` split, so a trailing slash or the (new) marker
// never reads as a difference. Both sides also drop the review pickup, which the
// coverage grammar requires in every fence and no ticket owns, and every path the
// implicit authorizing entries cover.
func fenceWritesCheck(f Facts) CheckResult {
	pickup, _ := reviewrecord.RecordPath(f.SpecPath)
	var writes []string
	for _, ticket := range f.Tickets {
		writes = append(writes, ticket.Writes...)
	}
	fence, owned := unionSide(f, f.FenceEntries, pickup), unionSide(f, writes, pickup)
	var sides []string
	if only := missingFrom(fence, owned); len(only) > 0 {
		sides = append(sides, "fence only: "+strings.Join(only, ", "))
	}
	if only := missingFrom(owned, fence); len(only) > 0 {
		sides = append(sides, "Writes only: "+strings.Join(only, ", "))
	}
	if len(sides) > 0 {
		return red("fence-writes", "spec fence and ticket Writes: union differ: "+strings.Join(sides, "; "))
	}
	return green("fence-writes")
}

// unionSide is one side of the comparison: each entry split to its tree path, less
// the review pickup and every path an implicit entry authorizes.
func unionSide(f Facts, entries []string, pickup string) map[string]bool {
	side := map[string]bool{}
	implicit := implicitEntries(f)
	for _, entry := range entries {
		path, _ := tickets.WritesPath(entry)
		if path == "" || path == pickup || fenceAuthorizes(path, implicit) {
			continue
		}
		side[path] = true
	}
	return side
}

// missingFrom is every path of one side that the other side lacks, sorted.
func missingFrom(side, other map[string]bool) []string {
	var only []string
	for path := range side {
		if !other[path] {
			only = append(only, path)
		}
	}
	sort.Strings(only)
	return only
}

// implicitEntries is every entry that authorizes a path without a fence line: the
// phase-owned capture folder, and the active spec's own folder, which authorizes an
// in-range amendment of the spec. The entries are derived rather than carried as their
// own facts, so the printed spec path and the authorized folder cannot disagree.
func implicitEntries(f Facts) []string {
	entries := []string{captureEntry}
	if folder := specFolder(f.SpecPath); folder != "" {
		entries = append(entries, folder)
	}
	return entries
}

// writesResolveCheck grades each `Writes:` entry against the tree. An entry without a
// marker must name a path the tree holds, or an exact committed deletion of the explicit
// source range. A (deleted) entry must name a tracked file at the review base, whatever
// the tip holds, so no marker gives authority over a path that never existed. Either
// entry also resolves through an exact deletion of the build's landing range, which an
// earlier chunk made before the review base. An entry that the ticket grammar refuses is
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
