// Package commitcmd adapts commitment operations to the Bench CLI.
package commitcmd

import (
	"fmt"
	"maps"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/commitment"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/refusalroute"
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
	store, c := commitrepo.Store{Root: root}, call{form: selected, flags: parsed.Flags}
	switch selected.name {
	case "start", "block", "unblock":
		return admission(store, c)
	case "show":
		return show(store, c)
	case "inventory":
		return inventory(store, c)
	case "plan":
		return plan(store, c)
	case "approve":
		return approve(store, c)
	case "verify":
		return verify(store, c)
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

func show(store commitrepo.Store, c call) (string, int) {
	policy, exists, err := store.Policy()
	if err != nil {
		return c.refuse(refusalroute.CommitmentHandback, err, nil)
	}
	outlook, err := store.Outlook()
	if err != nil {
		return c.refuse(refusalroute.CommitmentHandback, err, nil)
	}
	next, err := outlookTables(withCommand(outlook))
	if err != nil {
		return c.refuse(refusalroute.CommitmentHandback, err, nil)
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
		return c.refuse(refusalroute.CommitmentHandback, err, nil)
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

func inventory(store commitrepo.Store, c call) (string, int) {
	items, err := store.Inventory()
	if err != nil {
		return c.refuse(refusalroute.CommitmentHandback, err, nil)
	}
	rows := make([][]string, 0, len(items))
	for _, item := range items {
		rows = append(rows, []string{item.Kind, item.ID, item.State, item.Identity, item.Scope})
	}
	out, err := toon.Table("commitment_inventory", []string{"kind", "id", "state", "identity", "scope"}, rows)
	if err != nil {
		return c.refuse(refusalroute.CommitmentHandback, err, nil)
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

// plan validates and records the proposal in the input file. A refusal that raises no face
// of its own names the input, which the caller corrects.
func plan(store commitrepo.Store, c call) (string, int) {
	data, err := readInput(c.flags["--input"])
	if err != nil {
		return c.refuse(refusalroute.CommitmentPlanInput, err, nil)
	}
	result, err := store.Plan(data)
	if err != nil {
		return c.refuse(refusalroute.CommitmentPlanInput, err, nil)
	}
	planTable, err := toon.Table("commitment_plan", []string{"id", "predecessor", "proposal"}, [][]string{{result.ID, result.Predecessor, result.ProposalIdentity}})
	if err != nil {
		return c.refuse(refusalroute.CommitmentHandback, err, nil)
	}
	var rows [][]string
	for _, item := range effectRows(result.Effects) {
		rows = append(rows, item)
	}
	effects, err := toon.Table("effects", []string{"kind", "outcome"}, rows)
	if err != nil {
		return c.refuse(refusalroute.CommitmentHandback, err, nil)
	}
	return planTable + "\n" + effects + "\n", 0
}

// approve stages the exact transition of one plan. The approval is a commitment change, so a
// refusal that raises no face of its own is the reviewer's decision.
func approve(store commitrepo.Store, c call) (string, int) {
	primary, err := git.IsPrimaryCheckout(store.Root)
	if err != nil {
		return c.refuse(refusalroute.CommitmentHandback, err, nil)
	}
	if _, owned := intent.AssignmentForWorktree(store.Root); primary || !owned {
		return c.refuse(refusalroute.CommitmentNeedsAssignment, fmt.Errorf("an owned planning worktree is required"), nil)
	}

	delayed, ok := operandSet(c.flags["--delayed"])
	if !ok {
		return c.refuse(refusalroute.CommitmentDecision, fmt.Errorf("invalid delayed operands %q", c.flags["--delayed"]), nil)
	}
	removed, ok := operandSet(c.flags["--removed"])
	if !ok {
		return c.refuse(refusalroute.CommitmentDecision, fmt.Errorf("invalid removed operands %q", c.flags["--removed"]), nil)
	}
	changed, err := store.Approve(c.flags["--plan"], c.flags["--decision"], delayed, removed)
	if err != nil {
		return c.refuse(refusalroute.CommitmentDecision, err, nil)
	}
	out, err := toon.Table("commitment_approval", []string{"plan", "decision", "changed"}, [][]string{{c.flags["--plan"], c.flags["--decision"], fmt.Sprint(changed)}})
	if err != nil {
		return c.refuse(refusalroute.CommitmentHandback, err, nil)
	}
	return out + "\n", 0
}

// verify records a milestone verification receipt. A refusal names verify of the same
// milestone, because the next action is corrected evidence; success names the completion
// proposal.
func verify(store commitrepo.Store, c call) (string, int) {
	milestone := map[string]string{refusalroute.FactMilestone: c.flags["--milestone"]}
	data, err := readInput(c.flags["--evidence"])
	if err != nil {
		return c.refuse(refusalroute.CommitmentVerifyEvidence, err, milestone)
	}
	verification, err := store.Verify(c.flags["--milestone"], data)
	if err != nil {
		return c.refuse(refusalroute.CommitmentVerifyEvidence, err, milestone)
	}
	out, err := toon.Table("commitment_verification", []string{"id", "milestone", "revision"}, [][]string{{verification.ID, verification.Milestone, verification.Revision}})
	if err != nil {
		return c.refuse(refusalroute.CommitmentHandback, err, nil)
	}
	planForm, _ := selectForm("plan")
	return out + "\n" + nextTable(planForm.command()), 0
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

// call is one parsed commitment form: the form and the flag values that the caller passed.
type call struct {
	form  form
	flags map[string]string
}

// rerun is the caller's own command, each flag value rendered by refusalroute.Arg under the
// flag's placeholder.
func (c call) rerun() string {
	values := map[string]string{}
	for _, flag := range c.form.flags {
		values[flag.name] = refusalroute.Arg(strings.Trim(flag.placeholder, "<>"), c.flags[flag.name])
	}
	return "bench commitment" + c.form.suffixWith(values)
}

// refuse prints a refusal of the call: the cause, and in the next table the route of the face
// that cause raised, or else of face. values are the facts that the call observed, and the
// call's own re-run joins them.
func (c call) refuse(face string, cause error, values map[string]string) (string, int) {
	facts := map[string]string{refusalroute.FactRerun: c.rerun()}
	maps.Copy(facts, values)
	refusal := refusalroute.Printed(face, cause, facts)
	return toon.Errorf("bench commitment "+c.form.name+" refused", refusal.Sentence) + "\n" + nextTable(refusal.Route), 1
}

// nextTable is the next table that names one command: a refusal's route, or the command that
// follows a success. The registry owns the table's label and its column.
func nextTable(command string) string {
	table, _ := toon.Table(refusalroute.NextField, []string{refusalroute.NextColumn}, [][]string{{command}})
	return table + "\n"
}
