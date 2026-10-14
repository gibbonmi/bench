package preflight

import (
	"slices"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/preflight/ownership"
	"github.com/gibbonmi/bench/internal/tickets"
	"github.com/gibbonmi/bench/internal/toon"
)

// proposalToleratedChecks names every check whose red does not refuse a proposal,
// because the proposal itself reports the repair that the red requires.
var proposalToleratedChecks = append(closureCheckNames(), "fence-writes")

func proposeWritesCommand(root, mode, slug, base, sourceTip, name string, args []string) (string, int) {
	return preparedCommand(root, mode, slug, base, sourceTip, name, args)
}

func proposalSourceCheck(root string, facts Facts, selected *tickets.Entry) string {
	_, failure := loadBuildSources(root, facts, selected, buildSourcePolicy())
	return failure
}

func renderWritesProposal(root string, f Facts, name string) (string, int) {
	if refusal, _ := preparationCheckoutRefusal(root, f, "proposal"); refusal != "" {
		return refusal, 1
	}
	verdict := Decide(f)
	for _, check := range verdict.Checks {
		if check.Verdict == verdictRed && !slices.Contains(proposalToleratedChecks, check.Check) {
			return chargeVerdictRefusal(verdict), 1
		}
	}
	entry, selected, detail, next := preparationTicket(root, f, name)
	if detail != "" {
		return chargeRefusal("ticket", detail, next), 1
	}
	if failure := proposalSourceCheck(root, f, entry); failure != "" {
		return chargeRefusal("source", failure, "restore the named canonical source and rerun the exact proposal"), 1
	}
	rows := proposalRows(f, *selected)
	tables := []struct {
		name   string
		fields []string
		rows   [][]string
	}{
		{"writes_proposal", []string{"path", "source", "fence"}, rows},
		{"ordering", []string{"ticket", "other", "required"}, proposalEdges(f.Tickets, *selected, rows)},
		{"closure", []string{"ticket", "path", "requires"}, closureTable(f, selected.Name)},
		{"closure_evidence", []string{"ticket", "path", "kind", "source", "line", "reason"}, evidenceTable(f, selected.Name)},
		{"source", []string{"base", "tip"}, [][]string{{f.SourceBase, f.SourceTip}}},
		{"premises", []string{"id", "state", "source", "reason"}, premiseTable(f)},
	}
	var out strings.Builder
	for _, table := range tables {
		rendered, err := toon.Table(table.name, table.fields, table.rows)
		if err != nil {
			return toon.RenderError(err) + "\n", 1
		}
		out.WriteString(rendered)
	}
	return out.String(), 0
}

// selectedRequirements is every requirement of the selected ticket, owned or not, in
// closure-family order.
func selectedRequirements(f Facts, name string) []ownership.Requirement {
	var selected []ownership.Requirement
	for _, closure := range closureFamily {
		for _, requirement := range f.Ownership.Requirements {
			if requirement.Ticket == name && requirement.Kind == closure.kind {
				selected = append(selected, requirement)
			}
		}
	}
	return selected
}

// closureTable is the selected ticket's transitive closure: each originating write entry
// beside each path it requires.
func closureTable(f Facts, name string) [][]string {
	var rows [][]string
	for _, requirement := range selectedRequirements(f, name) {
		row := []string{requirement.Ticket, requirement.Entry, requirement.Path}
		if !slices.ContainsFunc(rows, func(seen []string) bool { return slices.Equal(seen, row) }) {
			rows = append(rows, row)
		}
	}
	return rows
}

// evidenceTable is every source reason of the selected ticket's closure, each naming its
// canonical owner and exact location, then each BENCH_KIT obligation a required system
// test carries.
func evidenceTable(f Facts, name string) [][]string {
	var rows [][]string
	add := func(row []string) {
		if !slices.ContainsFunc(rows, func(seen []string) bool { return slices.Equal(seen, row) }) {
			rows = append(rows, row)
		}
	}
	for _, requirement := range selectedRequirements(f, name) {
		for _, reason := range requirement.Provenance {
			add([]string{requirement.Ticket, requirement.Path, string(requirement.Kind), reason.Owner, strconv.Itoa(reason.Line), reason.Reason})
		}
	}
	for _, obligation := range f.KitObligations {
		if obligation.Ticket == name {
			add([]string{obligation.Ticket, obligation.Path, "kit-pin", obligation.Owner, strconv.Itoa(obligation.Line), "system-tagged test requires the ticket to state " + kitEnvMarker})
		}
	}
	return rows
}

// premiseTable is every premise row of the inventory, then each diagnostic as an unknown
// row at the source tip, so an incomplete collection never prints as complete proof.
func premiseTable(f Facts) [][]string {
	var rows [][]string
	for _, premise := range f.Ownership.Premises {
		rows = append(rows, []string{premise.ID, string(premise.State), premise.Source, premise.Reason})
	}
	for _, diagnostic := range f.Ownership.Diagnostics {
		rows = append(rows, []string{diagnostic.ID, string(ownership.StateUnknown), f.SourceTip, diagnostic.Reason})
	}
	return rows
}

func proposalRows(f Facts, selected tickets.Ticket) [][]string {
	seen := map[string]bool{}
	var rows [][]string
	for _, closure := range closureFamily {
		if !closure.writes {
			continue
		}
		for _, requirement := range missingClosures(f, closure.kind) {
			if requirement.Ticket != selected.Name || seen[requirement.Path] {
				continue
			}
			seen[requirement.Path] = true
			rows = append(rows, []string{requirement.Path, proposalSource(requirement), proposalFence(requirement.Path, f.FenceEntries)})
		}
	}
	return rows
}

// proposalSource names the fact that requires one proposed path.
func proposalSource(requirement ownership.Requirement) string {
	switch requirement.Kind {
	case ownership.KindFixture:
		return "fixture " + requirement.Path
	case ownership.KindAnchor:
		return "anchor " + requirement.Entry
	case ownership.KindCaller:
		return "caller " + requirement.Entry
	}
	return "registry " + tickets.BindingRegistry
}

func proposalFence(path string, entries []string) string {
	if fenceAuthorizes(path, entries) {
		return "covered"
	}
	return "spec fence expansion required"
}

func proposalEdges(all []tickets.Ticket, selected tickets.Ticket, rows [][]string) [][]string {
	proposed := append(append([]string{}, ownership.OwnedPaths(selected)...), pathsFromRows(rows)...)
	var edges [][]string
	for _, other := range all {
		if other.Name == selected.Name || ordered(selected, other, all) || !overlaps(proposed, ownership.OwnedPaths(other)) {
			continue
		}
		edges = append(edges, []string{selected.Name, other.Name, "approved ordering edge required"})
	}
	return edges
}

func pathsFromRows(rows [][]string) []string {
	paths := make([]string, len(rows))
	for i, row := range rows {
		paths[i] = row[0]
	}
	return paths
}

func overlaps(left, right []string) bool {
	for _, a := range left {
		for _, b := range right {
			if tickets.Covers(b, a) || tickets.Covers(a, b) {
				return true
			}
		}
	}
	return false
}

func ordered(left, right tickets.Ticket, all []tickets.Ticket) bool {
	return reaches(left.Name, right.Name, all, map[string]bool{}) || reaches(right.Name, left.Name, all, map[string]bool{})
}

func reaches(from, target string, all []tickets.Ticket, seen map[string]bool) bool {
	if seen[from] {
		return false
	}
	seen[from] = true
	for _, ticket := range all {
		if ticket.Name != from {
			continue
		}
		for _, blocker := range ticket.Blockers {
			if blocker == target || reaches(blocker, target, all, seen) {
				return true
			}
		}
	}
	return false
}
