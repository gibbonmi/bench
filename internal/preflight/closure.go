package preflight

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/preflight/ownership"
	"github.com/gibbonmi/bench/internal/tickets"
)

// closureFamily is every closure kind in row order. Decide renders one row for each
// kind, the proposal lists each kind's missing files, and the proposal tolerates
// each kind's red, so a new kind joins all three from this one table.
var closureFamily = []struct {
	kind  ownership.Kind
	check string
	lead  string
	verb  string
}{
	{ownership.KindFixture, "fixture-closure", "Writes: entry names a fixture-pinned path without naming the fixture: ", "is pinned by"},
	{ownership.KindRegistry, "registry-closure", "Writes: entry names a bound package without naming every bound file: ", "requires"},
	{ownership.KindAnchor, "anchor-closure", "Writes: entry names an anchored guidance path without naming every anchor registry file: ", "is anchored by"},
}

// missingClosures is every requirement of one kind that its ticket does not own, in
// ticket, entry, and path order. Verdicts and proposals both read it.
func missingClosures(f Facts, kind ownership.Kind) []ownership.Requirement {
	owned := map[string][]string{}
	for _, ticket := range f.Tickets {
		owned[ticket.Name] = append(owned[ticket.Name], ownership.OwnedPaths(ticket)...)
	}
	var missing []ownership.Requirement
	for _, requirement := range f.Ownership.Requirements {
		if requirement.Kind == kind && !ownership.Owned(owned[requirement.Ticket], requirement.Path) {
			missing = append(missing, requirement)
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

func closureMessages(requirements []ownership.Requirement) []string {
	result := make([]string, len(requirements))
	for i, requirement := range requirements {
		for _, closure := range closureFamily {
			if closure.kind == requirement.Kind {
				result[i] = requirement.Ticket + ": " + requirement.Entry + " " + closure.verb + " " + requirement.Path
			}
		}
	}
	return result
}

// probe records the tree facts of one `Writes:` entry: whether its path exists and
// whether it joins the system suite. The closure itself is the ownership collector's; the
// policy over both belongs to Decide. A (new) test file the tree does not hold yet takes
// its system tag from its directory, so build preflight grades the kit pin before the
// build creates the file.
func (facts *ticketFacts) probe(root, entry string) {
	path, isNew := tickets.WritesPath(entry)
	facts.writes[entry] = treeHolds(root, path)
	if systemTagged(root, path) || isNew && !facts.writes[entry] && newSystemTest(root, path) {
		facts.systemTag[entry] = true
	}
}

// collectOwnership reads the canonical closure of the parsed tickets against the source
// that the gatherer resolved, and the spec's planning-facts inventory. An owner that does
// not read refuses the gather, so no partial closure reaches Decide.
func collectOwnership(root string, content []byte, source ownership.Source, parsed []tickets.Ticket) (ownership.Snapshot, *BootstrapFailure) {
	snapshot, err := ownership.Collect(root, source, ownership.ParseManifest(content), parsed)
	var refusal *ownership.Error
	if errors.As(err, &refusal) {
		return ownership.Snapshot{}, &BootstrapFailure{refusal.Stage + " refused", refusal.Error()}
	}
	if err != nil {
		return ownership.Snapshot{}, &BootstrapFailure{"ownership not readable", err.Error()}
	}
	return snapshot, nil
}

// ownershipSource is the collector's view of the resolved source range and the verified
// kit context. The seal verifier already graded the published binary against root's
// sources, so only a verified seal names a kit source identity.
func ownershipSource(root, specPath string, source diff.SourceRange, changed []string, sealPresent, sealVerified bool) ownership.Source {
	paths := make([]ownership.PathStatus, 0, len(changed))
	for _, path := range changed {
		status := ownership.StatusChanged
		if slices.Contains(source.DeletedPaths, path) {
			status = ownership.StatusDeleted
		}
		paths = append(paths, ownership.PathStatus{Path: path, Status: status})
	}
	result := ownership.Source{Base: source.Base, Tip: source.Tip, SpecPath: specPath, Paths: paths}
	if sealPresent && sealVerified {
		if sources, _, err := freshness.SealDigests(freshness.PublishedExecutable(root)); err == nil {
			result.Kit.Identity = sources
		}
	}
	return result
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
