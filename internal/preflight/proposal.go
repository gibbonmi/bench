package preflight

import (
	"github.com/gibbonmi/bench/internal/tickets"
	"github.com/gibbonmi/bench/internal/toon"
)

const ticketBindingRegistry = "internal/tickets/registry_data.go"

func proposeWritesCommand(root, mode, slug, base, sourceTip, name string, args []string) (string, int) {
	return preparedCommand(root, mode, slug, base, sourceTip, name, false, args, proposalPreparation)
}

func proposalSourceCheck(root string, facts Facts, selected *tickets.Entry) string {
	_, failure := chargeSources(root, facts.SourceTip, facts.SpecPath, selected)
	return failure
}

func renderWritesProposal(root string, f Facts, name string) (string, int) {
	if refusal := preparationCheckoutRefusal(root, f, "proposal"); refusal != "" {
		return refusal, 1
	}
	verdict := Decide(f)
	for _, check := range verdict.Checks {
		if check.Verdict == verdictRed && check.Check != "fixture-closure" && check.Check != "registry-closure" {
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
	table, err := toon.Table("writes_proposal", []string{"path", "source", "fence"}, rows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	edges, err := toon.Table("ordering", []string{"ticket", "other", "required"}, proposalEdges(f.Tickets, *selected, rows))
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return table + edges, 0
}

func proposalRows(f Facts, selected tickets.Ticket) [][]string {
	seen := map[string]bool{}
	var rows [][]string
	for _, requirement := range append(missingClosures(f, fixtureClosure), missingClosures(f, registryClosure)...) {
		if requirement.ticket != selected.Name || seen[requirement.path] {
			continue
		}
		seen[requirement.path] = true
		source := "registry " + ticketBindingRegistry
		if requirement.kind == fixtureClosure {
			source = "fixture " + requirement.path
		}
		rows = append(rows, []string{requirement.path, source, proposalFence(requirement.path, f.FenceEntries)})
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
	proposed := append(append([]string{}, ownedPaths(selected)...), pathsFromRows(rows)...)
	var edges [][]string
	for _, other := range all {
		if other.Name == selected.Name || ordered(selected, other, all) || !overlaps(proposed, ownedPaths(other)) {
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
			if pathCovered(a, []string{b}) || pathCovered(b, []string{a}) {
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
