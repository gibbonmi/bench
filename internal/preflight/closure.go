package preflight

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/tickets"
)

type closureKind string

const (
	fixtureClosure  closureKind = "fixture"
	registryClosure closureKind = "registry"
	anchorClosure   closureKind = "anchor"
)

// closureFamily is every closure kind in row order. Decide renders one row for each
// kind, the proposal lists each kind's missing files, and the proposal tolerates
// each kind's red, so a new kind joins all three from this one table.
var closureFamily = []struct {
	kind  closureKind
	check string
	lead  string
	verb  string
}{
	{fixtureClosure, "fixture-closure", "Writes: entry names a fixture-pinned path without naming the fixture: ", "is pinned by"},
	{registryClosure, "registry-closure", "Writes: entry names a bound package without naming every bound file: ", "requires"},
	{anchorClosure, "anchor-closure", "Writes: entry names an anchored guidance path without naming every anchor registry file: ", "is anchored by"},
}

type closureRequirement struct {
	ticket string
	entry  string
	path   string
	kind   closureKind
}

// closureFacts is the gathered requirement map one closure kind grades.
func closureFacts(f Facts, kind closureKind) map[string][]string {
	switch kind {
	case fixtureClosure:
		return f.WritesFixturePins
	case anchorClosure:
		return f.WritesAnchorFiles
	}
	return f.WritesBoundFiles
}

// missingClosures derives omitted closure facts for verdicts and proposals.
func missingClosures(f Facts, kind closureKind) []closureRequirement {
	var missing []closureRequirement
	seen := map[string]bool{}
	required := closureFacts(f, kind)
	for _, ticket := range f.Tickets {
		owned := ownedPaths(ticket)
		for _, entry := range ticket.Writes {
			for _, path := range required[entry] {
				if pathCovered(path, owned) {
					continue
				}
				key := ticket.Name + "\x00" + entry + "\x00" + path + "\x00" + string(kind)
				if !seen[key] {
					seen[key] = true
					missing = append(missing, closureRequirement{ticket.Name, entry, path, kind})
				}
			}
		}
	}
	return missing
}

// closureRows grades each closure kind into its ticket-gated row. A ticket that
// writes a pinned, bound, or anchored path and omits the file that holds its proof
// finds that file mid-build and pays a repair round.
func closureRows(f Facts) []CheckResult {
	rows := make([]CheckResult, 0, len(closureFamily))
	for _, closure := range closureFamily {
		rows = append(rows, ticketRow(f, closure.check, func(f Facts) CheckResult {
			if missing := closureMessages(missingClosures(f, closure.kind)); len(missing) > 0 {
				return red(closure.check, closure.lead+strings.Join(missing, ", "))
			}
			return green(closure.check)
		}))
	}
	return rows
}

// closureCheckNames is the check name of every closure kind, in row order.
func closureCheckNames() []string {
	names := make([]string, len(closureFamily))
	for i, closure := range closureFamily {
		names[i] = closure.check
	}
	return names
}

func closureMessages(requirements []closureRequirement) []string {
	result := make([]string, len(requirements))
	for i, requirement := range requirements {
		for _, closure := range closureFamily {
			if closure.kind == requirement.kind {
				result[i] = requirement.ticket + ": " + requirement.entry + " " + closure.verb + " " + requirement.path
			}
		}
	}
	return result
}

// probe records the tree facts of one `Writes:` entry. The probe is the gatherer's
// whole I/O contribution to the ownership and closure rows; the policy over the
// resulting facts belongs to Decide.
func (facts *ticketFacts) probe(root, entry string, pins, anchors map[string][]string) {
	path, _ := splitWritesEntry(entry)
	facts.writes[entry] = treeHolds(root, path)
	if pinning := pins[path]; len(pinning) > 0 {
		facts.pins[entry] = pinning
	}
	if bound := tickets.BoundFiles(path); len(bound) > 0 {
		facts.bound[entry] = bound
	}
	if files := anchorFiles(anchors, path); len(files) > 0 {
		facts.anchors[entry] = files
	}
	if systemTagged(root, path) {
		facts.systemTag[entry] = true
	}
}

// anchorFiles is every anchor registry file whose literal names path, or a path under
// it at a `/` boundary, sorted. A directory entry therefore takes the closure of each
// anchored file it holds.
func anchorFiles(anchors map[string][]string, path string) []string {
	if path == "" {
		return nil
	}
	seen := map[string]bool{}
	var files []string
	for literal, holders := range anchors {
		if !pathCovered(literal, []string{path}) {
			continue
		}
		for _, file := range holders {
			if !seen[file] {
				seen[file] = true
				files = append(files, file)
			}
		}
	}
	sort.Strings(files)
	return files
}

// treeHolds reports whether one `Writes:` path names something in the tree.
// Any file type answers yes: the row grades whether the path exists, not
// what sits at it. The lstat never follows a link, so a dangling symlink
// reads as the absent path it is.
func treeHolds(root, path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Lstat(filepath.Join(root, filepath.FromSlash(path)))
	return err == nil
}

// ownedPaths is every tree path one ticket declares, with the (new) marker
// stripped. The closures grade their required names against this one set, so no
// two of them can disagree about what a ticket already owns.
func ownedPaths(ticket tickets.Ticket) []string {
	paths := make([]string, 0, len(ticket.Writes))
	for _, entry := range ticket.Writes {
		path, _ := splitWritesEntry(entry)
		paths = append(paths, path)
	}
	return paths
}

// pathCovered reports whether one required path is already named by the ticket,
// either exactly or through a directory entry that contains it.
func pathCovered(required string, owned []string) bool {
	for _, path := range owned {
		if required == path || strings.HasPrefix(required, path+"/") {
			return true
		}
	}
	return false
}
