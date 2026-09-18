package preflight

import (
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/reviewrecord"
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
		path, _ := splitWritesEntry(entry)
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
