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

type flag struct{ name, placeholder string }

type form struct {
	name, description string
	flags             []flag
}

var forms = []form{
	{name: "start", description: "claim an eligible committed outcome", flags: []flag{{"--outcome", "<id>"}, {"--request", "<request>"}, {"--deliverable", "<path>"}}},
	{name: "block", description: "record an outcome blocker", flags: []flag{{"--outcome", "<id>"}, {"--reason", "<text>"}}},
	{name: "unblock", description: "clear an outcome blocker", flags: []flag{{"--outcome", "<id>"}}},
	{name: "show", description: "show the current delivery commitment"},
	{name: "inventory", description: "list roadmap obligations, staged deliverables, and run identities"},
	{name: "plan", description: "validate an exact commitment transition", flags: []flag{{"--input", "<file>"}}},
	{name: "approve", description: "approve and stage one exact commitment transition", flags: []flag{
		{"--plan", "<id>"}, {"--decision", "<reference>"}, {"--delayed", "<ids-or-none>"}, {"--removed", "<ids-or-none>"},
	}},
	{name: "verify", description: "verify milestone criterion evidence and record a completion receipt", flags: []flag{{"--milestone", "<id>"}, {"--evidence", "<file>"}}},
}

// suffix is the form's help grammar after the family name: every flag shows its placeholder.
func (f form) suffix() string { return f.suffixWith(nil) }

// suffixWith is the form's grammar after the family name, with each flag that values names
// carrying that value in place of its placeholder.
func (f form) suffixWith(values map[string]string) string {
	terms := []string{"", f.name}
	for _, flag := range f.flags {
		value, known := values[flag.name]
		if !known {
			value = flag.placeholder
		}
		terms = append(terms, flag.name, value)
	}
	return strings.Join(terms, " ")
}

// command is the complete command form that a next action names.
func (f form) command() string { return "bench commitment" + f.suffix() }

func (f form) grammar() usage.Grammar {
	g := usage.Grammar{Cmd: "bench commitment " + f.name, Help: "usage: bench commitment" + f.suffix(), MaxArgs: 0}
	for _, flag := range f.flags {
		g.Flags = append(g.Flags, usage.Flag{Name: flag.name, HasValue: flag.placeholder != "", NoEmptyValue: flag.placeholder != "", Required: true})
	}
	return g
}

// HelpRows returns the command family's public forms.
func HelpRows() []HelpRow {
	rows := make([]HelpRow, len(forms))
	for index, form := range forms {
		rows[index] = HelpRow{Suffix: form.suffix(), Description: form.description}
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
	grammar := selected.grammar()
	parsed, line, code := usage.Parse(grammar, args[1:])
	if line != "" {
		return line + "\n", code
	}
	if root == "" {
		return toon.NotInRepo() + "\n", 1
	}
	store := commitrepo.Store{Root: root}
	switch selected.name {
	case "start", "block", "unblock":
		return admission(store, selected.name, parsed.Flags)
	case "show":
		return show(store)
	case "inventory":
		return inventory(store)
	case "plan":
		return plan(store, parsed.Flags["--input"])
	case "approve":
		return approve(store, parsed.Flags)
	case "verify":
		return verify(store, selected, parsed.Flags)
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
		lines = append(lines, prefix+form.suffix())
	}
	return strings.Join(lines, "\n") + "\n"
}

func show(store commitrepo.Store) (string, int) {
	policy, exists, err := store.Policy()
	if err != nil {
		return refusal("show", err)
	}
	outlook, err := store.Outlook()
	if err != nil {
		return refusal("show", err)
	}
	next, err := outlookTables(withCommand(outlook))
	if err != nil {
		return refusal("show", err)
	}
	if !exists {
		out, _ := toon.Table("commitment", []string{"state", "active_milestone", "next_outcome"}, [][]string{{commitment.OutlookAdoptionRequired, "", ""}})
		return out + "\n" + next, 0
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
	return out + "\n" + next, 0
}

// Outlook is the one commitment projection for root, with its next-action command filled
// from the form table. The roadmap, status, and dashboard readers render it. A policy or
// runtime record that does not read projects as unreadable, so a reader stays available
// and `bench commitment show` names the cause.
func Outlook(root string) commitment.Outlook {
	outlook, err := commitrepo.Store{Root: root}.Outlook()
	if err != nil {
		outlook = commitment.Unreadable()
	}
	return withCommand(outlook)
}

// withCommand fills the outlook's next-action command. Each value that the outlook knows
// replaces its placeholder; the request stays a placeholder, because only the worker names it.
func withCommand(outlook commitment.Outlook) commitment.Outlook {
	selected, ok := selectForm(outlook.Operation)
	if !ok {
		return outlook
	}
	values := map[string]string{}
	for flag, value := range map[string]string{"--outcome": outlook.Next, "--milestone": outlook.Milestone, "--deliverable": outlook.Deliverable} {
		if value != "" {
			values[flag] = value
		}
	}
	outlook.Command = "bench commitment" + selected.suffixWith(values)
	return outlook
}

// outlookTables renders the outlook as its TOON outlook table and its blocker table.
func outlookTables(outlook commitment.Outlook) (string, error) {
	next, err := toon.Table(commitment.OutlookTable, commitment.OutlookFields, [][]string{outlook.Cells()})
	if err != nil {
		return "", err
	}
	blockers, err := toon.Table(commitment.BlockerTable, commitment.BlockerFields, outlook.BlockerCells())
	if err != nil {
		return "", err
	}
	return next + "\n" + blockers + "\n", nil
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

// readInput reads one regular input file without following a link.
func readInput(path string) ([]byte, error) {
	classified := bounds.ClassifyNoFollow(filepath.Clean(path))
	if classified.State != bounds.StateParsed {
		return nil, fmt.Errorf("input %q is %s: %s", path, classified.State, classified.Reason)
	}
	return classified.Data, nil
}

func plan(store commitrepo.Store, path string) (string, int) {
	data, err := readInput(path)
	if err != nil {
		return refusal("plan", err)
	}
	result, err := store.Plan(data)
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

// verify records a milestone verification receipt. A refusal names verify itself, because
// the next action is corrected evidence; success names the completion proposal.
func verify(store commitrepo.Store, selected form, flags map[string]string) (string, int) {
	retry := selected.command()
	data, err := readInput(flags["--evidence"])
	if err != nil {
		return refusalNext("verify", err, retry)
	}
	verification, err := store.Verify(flags["--milestone"], data)
	if err != nil {
		return refusalNext("verify", err, retry)
	}
	out, err := toon.Table("commitment_verification", []string{"id", "milestone", "revision"}, [][]string{{verification.ID, verification.Milestone, verification.Revision}})
	if err != nil {
		return refusalNext("verify", err, retry)
	}
	planForm, _ := selectForm("plan")
	actions, err := toon.Table("next", []string{"command"}, [][]string{{planForm.command()}})
	if err != nil {
		return refusalNext("verify", err, retry)
	}
	return out + "\n" + actions + "\n", 0
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
	}{{"added", effects.Added}, {"reordered", effects.Reordered}, {"delayed", effects.Delayed}, {"removed", effects.Removed}, {"activated", effects.Activated}, {"switched", effects.Switched}, {"parallel-authorized", effects.ParallelAuthorized}, {"completed", effects.Completed}} {
		for _, outcome := range outcomes.ids {
			rows = append(rows, []string{outcomes.name, outcome})
		}
	}
	return rows
}

func refusal(operation string, err error) (string, int) {
	form, _ := selectForm("plan")
	return refusalNext(operation, err, form.command())
}

func refusalNext(operation string, err error, next string) (string, int) {
	actions, _ := toon.Table("next", []string{"command"}, [][]string{{next}})
	return toon.Errorf("bench commitment "+operation+" refused", err.Error()) + "\n" + actions + "\n", 1
}
