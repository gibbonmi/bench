package preflight

import (
	"slices"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/preflight/evidencecmd"
	"github.com/gibbonmi/bench/internal/preflight/ownership"
	"github.com/gibbonmi/bench/internal/tickets"
	"github.com/gibbonmi/bench/internal/toon"
)

// proposalToleratedChecks names every check whose red does not refuse a proposal,
// because the proposal itself reports the repair that the red requires.
var proposalToleratedChecks = append(closureCheckNames(), "fence-writes")

// breakdownToleratedChecks adds two reds of a graph that slicing has not finished to the
// proposal's set: rows-owned, whose uncited rows the breakdown lists, and completion-plan,
// because a valid plan names only ticket files that exist. Only this read-only view
// tolerates them.
var breakdownToleratedChecks = append(slices.Clone(proposalToleratedChecks), "rows-owned", "completion-plan")

func proposeWritesCommand(root string, kind evidencecmd.Kind, mode, slug, base, sourceTip, name string, args []string) (string, int) {
	return preparedCommand(root, kind, mode, slug, base, sourceTip, name, args)
}

func proposalSourceCheck(root string, facts Facts, selected *tickets.Entry) string {
	_, failure := loadBuildSources(root, facts, selected, buildSourcePolicy())
	return failure
}

func renderWritesProposal(root string, f Facts, name string) (string, int) {
	if refusal, _ := preparationCheckoutRefusal(root, f, "proposal"); refusal != "" {
		return refusal, 1
	}
	if refused := untolerated(Decide(f), proposalToleratedChecks); refused.Red {
		return chargeVerdictRefusal(refused), 1
	}
	entry, selected, detail, next := preparationTicket(root, f, name)
	if detail != "" {
		return chargeRefusal("ticket", detail, next), 1
	}
	if failure := proposalSourceCheck(root, f, entry); failure != "" {
		return chargeRefusal("source", failure, "restore the named canonical source and rerun the exact proposal"), 1
	}
	rows := proposalRows(f, *selected)
	return renderTables([]responseTable{
		{"writes_proposal", []string{"path", "source", "fence"}, rows},
		{"ordering", []string{"ticket", "other", "required"}, proposalEdges(f.Tickets, *selected, rows)},
		{"closure", []string{"ticket", "path", "requires"}, closureTable(f, selected.Name)},
		{"closure_evidence", []string{"ticket", "path", "kind", "source", "line", "reason"}, evidenceTable(f, selected.Name)},
		{"source", []string{"base", "tip"}, [][]string{{f.SourceBase, f.SourceTip}}},
		{"premises", []string{"id", "state", "source", "reason"}, premiseTable(f)},
	})
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
		if closure.source == nil {
			continue
		}
		for _, requirement := range missingClosures(f, closure.kind) {
			if requirement.Ticket != selected.Name || seen[requirement.Path] {
				continue
			}
			seen[requirement.Path] = true
			rows = append(rows, []string{requirement.Path, closure.source(requirement), proposalFence(requirement.Path, f.FenceEntries)})
		}
	}
	return rows
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
	for _, other := range ownership.Unordered(all, selected, proposed) {
		edges = append(edges, []string{selected.Name, other, "approved ordering edge required"})
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

// renderBreakdown prints the read-only authoring view of the gathered facts: the source,
// each uncited coverage row with its seam cell, every ticket's closure, the premises, and
// each tolerated red check. It refuses every other red check, and it grants no frontier,
// ordering, or write authority.
func renderBreakdown(f Facts) (string, int) {
	verdict := Decide(f)
	if refused := untolerated(verdict, breakdownToleratedChecks); refused.Red {
		return chargeVerdictRefusal(refused), 1
	}
	var uncited, closure [][]string
	for _, row := range f.Ownership.Uncited {
		uncited = append(uncited, []string{row.ID, row.Story, row.Seam})
	}
	for _, ticket := range f.Tickets {
		closure = append(closure, closureTable(f, ticket.Name)...)
	}
	return renderTables([]responseTable{
		{"source", []string{"base", "tip"}, [][]string{{f.SourceBase, f.SourceTip}}},
		{"uncited", []string{"row", "story", "seam"}, uncited},
		{"closure", []string{"ticket", "path", "requires"}, closure},
		{"premises", []string{"id", "state", "source", "reason"}, premiseTable(f)},
		{"checks", checkFields, redRows(verdict)},
	})
}

// untolerated is verdict with only its red checks that tolerated does not name, so a refusal
// names the first red that refuses and never a red that the response itself reports.
func untolerated(verdict Verdict, tolerated []string) Verdict {
	var refused Verdict
	for _, check := range verdict.Checks {
		if check.Verdict == verdictRed && !slices.Contains(tolerated, check.Check) {
			refused.Checks, refused.Red = append(refused.Checks, check), true
		}
	}
	return refused
}

// responseTable is one TOON table of a proposal or breakdown response.
type responseTable struct {
	name   string
	fields []string
	rows   [][]string
}

// renderTables renders tables in order, or the render error alone.
func renderTables(tables []responseTable) (string, int) {
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
