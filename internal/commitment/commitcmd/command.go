// Package commitcmd adapts commitment operations to the Bench CLI.
package commitcmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/commitment"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// HelpRow is one public root-help row for a commitment operation.
type HelpRow struct{ Suffix, Description string }

type form struct {
	name, suffix, description string
	flags                     []usage.Flag
}

var forms = []form{
	{name: "show", suffix: " show", description: "show the current delivery commitment"},
	{name: "inventory", suffix: " inventory", description: "list roadmap obligations, staged deliverables, and run identities"},
	{name: "plan", suffix: " plan --input <file>", description: "validate an exact commitment transition", flags: []usage.Flag{{Name: "--input", HasValue: true, NoEmptyValue: true, Required: true}}},
	{name: "approve", suffix: " approve --plan <id> --decision <reference> --delayed <ids-or-none> --removed <ids-or-none>", description: "approve and stage one exact commitment transition", flags: []usage.Flag{
		{Name: "--plan", HasValue: true, NoEmptyValue: true, Required: true},
		{Name: "--decision", HasValue: true, NoEmptyValue: true, Required: true},
		{Name: "--delayed", HasValue: true, NoEmptyValue: true, Required: true},
		{Name: "--removed", HasValue: true, NoEmptyValue: true, Required: true},
	}},
}

// HelpRows returns the command family's public forms.
func HelpRows() []HelpRow {
	rows := make([]HelpRow, len(forms))
	for index, form := range forms {
		rows[index] = HelpRow{Suffix: form.suffix, Description: form.description}
	}
	return rows
}

// Command runs one commitment operation in root.
func Command(root string, args []string) (string, int) {
	if len(args) == 1 && (args[0] == "help" || args[0] == "--help" || args[0] == "-h") {
		return usageText(), 0
	}
	if len(args) == 0 {
		return usageText(), 2
	}
	selected, ok := selectForm(args[0])
	if !ok {
		return toon.Usage("bench commitment", args[0]) + "\n", 2
	}
	grammar := usage.Grammar{Cmd: "bench commitment " + selected.name, Help: "usage: bench commitment" + selected.suffix, Flags: selected.flags, MaxArgs: 0}
	parsed, line, code := usage.Parse(grammar, args[1:])
	if line != "" {
		return line + "\n", code
	}
	if root == "" {
		return toon.NotInRepo() + "\n", 1
	}
	store := commitrepo.Store{Root: root}
	switch selected.name {
	case "show":
		return show(store)
	case "inventory":
		return inventory(store)
	case "plan":
		return plan(store, parsed.Flags["--input"])
	case "approve":
		return approve(store, parsed.Flags)
	default:
		panic("unreachable commitment form")
	}
}

func selectForm(name string) (form, bool) {
	for _, form := range forms {
		if form.name == name {
			return form, true
		}
	}
	return form{}, false
}

func usageText() string {
	var lines []string
	for index, form := range forms {
		prefix := "       bench commitment"
		if index == 0 {
			prefix = "usage: bench commitment"
		}
		lines = append(lines, prefix+form.suffix)
	}
	return strings.Join(lines, "\n") + "\n"
}

func show(store commitrepo.Store) (string, int) {
	policy, exists, err := store.Policy()
	if err != nil {
		return refusal("show", err)
	}
	if !exists {
		out, _ := toon.Table("commitment", []string{"state", "active_milestone", "next_outcome"}, [][]string{{"adoption-required", "", ""}})
		return out + "\n", 0
	}
	projection := commitment.Selection(policy)
	rows := make([][]string, 0, len(projection.Outcomes))
	for _, outcome := range projection.Outcomes {
		rows = append(rows, []string{"active", projection.Milestone, outcome})
	}
	out, err := toon.Table("commitment", []string{"state", "active_milestone", "outcome"}, rows)
	if err != nil {
		return refusal("show", err)
	}
	return out + "\n", 0
}

func inventory(store commitrepo.Store) (string, int) {
	items, err := store.Inventory()
	if err != nil {
		return refusal("inventory", err)
	}
	rows := make([][]string, 0, len(items))
	for _, item := range items {
		rows = append(rows, []string{item.Kind, item.ID, item.State, item.Identity})
	}
	out, err := toon.Table("commitment_inventory", []string{"kind", "id", "state", "identity"}, rows)
	if err != nil {
		return refusal("inventory", err)
	}
	return out + "\n", 0
}

func plan(store commitrepo.Store, path string) (string, int) {
	classified := bounds.ClassifyNoFollow(filepath.Clean(path))
	if classified.State != bounds.StateParsed {
		return refusal("plan", fmt.Errorf("input %q is %s: %s", path, classified.State, classified.Reason))
	}
	result, err := store.Plan(classified.Data)
	if err != nil {
		return refusal("plan", err)
	}
	planTable, err := toon.Table("commitment_plan", []string{"id", "predecessor", "proposal"}, [][]string{{result.ID, result.Predecessor, result.ProposalIdentity}})
	if err != nil {
		return refusal("plan", err)
	}
	var rows [][]string
	for _, item := range effectRows(result.Effects) {
		rows = append(rows, item)
	}
	effects, err := toon.Table("effects", []string{"kind", "outcome"}, rows)
	if err != nil {
		return refusal("plan", err)
	}
	return planTable + "\n" + effects + "\n", 0
}

func approve(store commitrepo.Store, flags map[string]string) (string, int) {
	primary, err := git.IsPrimaryCheckout(store.Root)
	if err != nil {
		return refusal("approve", err)
	}
	if _, owned := intent.AssignmentForWorktree(store.Root); primary || !owned {
		return refusalNext("approve", fmt.Errorf("an owned planning worktree is required"), "bench worktree create --request <request> --label <label>")
	}

	delayed, ok := operandSet(flags["--delayed"])
	if !ok {
		return refusal("approve", fmt.Errorf("invalid delayed operands %q", flags["--delayed"]))
	}
	removed, ok := operandSet(flags["--removed"])
	if !ok {
		return refusal("approve", fmt.Errorf("invalid removed operands %q", flags["--removed"]))
	}
	changed, err := store.Approve(flags["--plan"], flags["--decision"], delayed, removed)
	if err != nil {
		return refusal("approve", err)
	}
	out, err := toon.Table("commitment_approval", []string{"plan", "decision", "changed"}, [][]string{{flags["--plan"], flags["--decision"], fmt.Sprint(changed)}})
	if err != nil {
		return refusal("approve", err)
	}
	return out + "\n", 0
}

func operandSet(value string) ([]string, bool) {
	if value == "none" {
		return []string{}, true
	}
	parts := strings.Split(value, ",")
	for _, part := range parts {
		if strings.TrimSpace(part) == "" || part != strings.TrimSpace(part) {
			return nil, false
		}
	}
	return parts, true
}

func effectRows(effects commitment.Effects) [][]string {
	var rows [][]string
	for _, outcomes := range []struct {
		name string
		ids  []string
	}{{"added", effects.Added}, {"reordered", effects.Reordered}, {"delayed", effects.Delayed}, {"removed", effects.Removed}, {"activated", effects.Activated}, {"switched", effects.Switched}, {"parallel-authorized", effects.ParallelAuthorized}} {
		for _, outcome := range outcomes.ids {
			rows = append(rows, []string{outcomes.name, outcome})
		}
	}
	return rows
}

func refusal(operation string, err error) (string, int) {
	form, _ := selectForm("plan")
	return refusalNext(operation, err, "bench commitment"+form.suffix)
}

func refusalNext(operation string, err error, next string) (string, int) {
	actions, _ := toon.Table("next", []string{"command"}, [][]string{{next}})
	return toon.Errorf("bench commitment "+operation+" refused", err.Error()) + "\n" + actions + "\n", 1
}
