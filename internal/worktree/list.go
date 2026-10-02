package worktree

import (
	"errors"
	"os"
	"strings"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// These constants name each table block that the package renders, so a test reads the
// name that its renderer passes to toon.Table or toon.TableTyped. Both list views render
// listTable.
const (
	listTable                 = "worktrees"
	createTable               = "worktree_create"
	releaseTable              = "worktree_release"
	buildTable                = "worktree_build"
	cleanupTable              = "worktree_cleanup"
	ignoredPathsTable         = "ignored_paths"
	refusalPathsTable         = "refusal_paths"
	resetPathsTable           = "reset_paths"
	poolReclaimTable          = "pool_reclaim"
	poolReclaimAppliedTable   = "pool_reclaim_applied"
	poolReclaimAggregateTable = "pool_reclaim_aggregate"
)

var worktreeListFields = []string{"id", "label", "request", "state", "source", "tree", "lease", "landed", "ignored"}

var worktreeListGrammar = usage.Grammar{Cmd: usage.WorktreeList, Help: "usage: " + usage.WorktreeList, HelpOnlyWhenSole: true, UnquotedEmptyPositional: true}

type listRow struct {
	values []any
	// orphanPath carries a foreign registration's path, and assignmentPath an owned
	// assignment's. A recovery action pastes the path its own verb takes, so the row
	// carries the address the table cell does not.
	orphanPath     string
	assignmentPath string
}

// ListCommand implements the read-only AXI worktree population query.
func ListCommand(root, _ string, args []string) (string, int) {
	if selectsWorktrees(args) {
		return listSelectedWorktrees(root, args)
	}
	j := defaultJoins()
	_, line, code := usage.Parse(worktreeListGrammar, args)
	if line != "" {
		return line + "\n", code
	}
	if !inRepository(root) {
		return toon.NotInRepo() + "\n", 1
	}
	registrations, err := git.Worktrees(root)
	if err != nil {
		var typed git.WorktreeFailure
		if errors.As(err, &typed) {
			return toon.Errorf(typed.Error(), typed.WorktreeAction()) + "\n", 1
		}
		return toon.Errorf("cannot read registered worktrees", "run git worktree list and retry") + "\n", 1
	}
	assignments, err := intent.Assignments(root)
	if err != nil {
		return assignmentsReadRefusal(), 1
	}

	def, defaultResolved := git.ResolvedDefault(root)
	rows := make([]listRow, 0, len(assignments)+len(registrations))
	assignedPaths := make(map[string]bool, len(assignments))
	for _, assignment := range assignments {
		assignedPaths[assignment.Worktree] = true
		rows = append(rows, listRow{values: listAssignmentRow(j, root, assignment, def, defaultResolved), assignmentPath: assignment.Worktree})
	}
	mainRoot := canonicalRoot(root)
	for _, registration := range registrations {
		if samePath(registration.Path, mainRoot) || assignedPaths[registration.Path] {
			continue
		}
		label := registration.Branch
		if label == "" {
			label = registration.Path
		}
		tree := listTree(registration.Path)
		row := listRow{values: []any{"foreign", label, "", "foreign", "foreign", tree, listLease(registration.Path), listLanded(root, registration.Branch, def, defaultResolved)}}
		row.values = append(row.values, listIgnored(j, registration.Path))
		if tree == "missing" {
			row.orphanPath = registration.Path
		}
		rows = append(rows, row)
	}
	values := make([][]any, 0, len(rows))
	for _, row := range rows {
		values = append(values, row.values)
	}
	landed := false
	for _, assignment := range assignments {
		if assignment.State != intent.StateActive {
			continue
		}
		if plan := planLandedAssignment(j, root, assignment, CleanupOptions{}); assignmentLanded(assignment, plan) {
			landed = true
			break
		}
	}
	out, err := toon.TableTyped(listTable, worktreeListFields, values)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	actions := actionsForRows(rows)
	if landed {
		actions = append(actions, recoverMissingTree(true, "", "").action())
	}
	help, err := axi.RenderHelp(actions)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return out + help, 0
}

// assignmentsReadRefusal is the response of each list view when the intent ledger cannot
// be read, so the bare and selected views name one repair.
func assignmentsReadRefusal() string {
	return toon.Errorf("cannot read worktree assignments", "repair the Bench intent ledger and retry") + "\n"
}

// actionsForRows derives the help rows from the row states. An active row with a present
// tree gives a path action and an exec action with a `<target>` slot, because the id cell
// already holds its address. The actions of each such row are identical, so the help
// render collapses them to one pair. The other states keep a row-specific action, because
// their operand appears in no cell.
func actionsForRows(rows []listRow) []axi.Action {
	actions := make([]axi.Action, 0, len(rows))
	for _, row := range rows {
		if len(row.values) < 4 {
			continue
		}
		if row.values[3] == string(intent.StateCleanupPending) {
			request, _ := row.values[listRequestCell].(string)
			if request == "" || row.assignmentPath == "" {
				continue
			}
			if rowUnlanded(row) {
				id, _ := row.values[0].(string)
				actions = append(actions, discardUnlanded(row.assignmentPath, id).action())
				continue
			}
			actions = append(actions, axi.ExecutableInvocation(
				"resume the cleanup-pending assignment",
				axi.KnownArgument("worktree"), axi.KnownArgument("release"),
				axi.KnownArgument("--request"), axi.KnownArgument(request), axi.KnownArgument(row.assignmentPath),
			))
			continue
		}
		if row.values[3] == string(intent.StateActive) {
			id, ok := row.values[0].(string)
			if !ok || id == "" {
				continue
			}
			// A row whose tree cell reads `missing` advertises its recovery verb alone. The
			// path and exec actions both enter the tree, so neither can succeed on it.
			if len(row.values) > listTreeCell && row.values[listTreeCell] == "missing" {
				request, _ := row.values[listRequestCell].(string)
				actions = append(actions, recoverMissingTree(rowLanded(row), request, row.assignmentPath).action())
				continue
			}
			actions = append(actions,
				axi.ExecutableInvocation("inspect an active worktree by its id", axi.KnownArgument("worktree"), axi.KnownArgument("path"), axi.FutureInput("target")),
				axi.ExecutableInvocation("run a command in an active worktree by its id", axi.KnownArgument("worktree"), axi.KnownArgument("exec"), axi.FutureInput("target"), axi.KnownArgument("--"), axi.FutureInput("command")))
			continue
		}
		if row.values[3] == "foreign" && row.orphanPath != "" {
			actions = append(actions, axi.ExecutableInvocation("clean the orphaned worktree", axi.KnownArgument("worktree"), axi.KnownArgument("clean"), axi.KnownArgument(row.orphanPath)))
		}
	}
	return actions
}

// The cells an action reads by position. The field list above declares the order, so a
// reader compares the two in one place.
const (
	listRequestCell = 2
	listTreeCell    = 5
	listLandedCell  = 7
)

// rowLanded reports the landedness the row itself discloses. A cell the query could not
// prove reads `unknown`, and an unproven record keeps the release route, because only a
// proven landing licenses the batch clean.
func rowLanded(row listRow) bool {
	return len(row.values) > listLandedCell && row.values[listLandedCell] == true
}

// rowUnlanded reports a branch the row proves unlanded. An `unknown` cell proves neither
// outcome, so that row keeps the release route.
func rowUnlanded(row listRow) bool {
	return len(row.values) > listLandedCell && row.values[listLandedCell] == false
}

// recoveryRoute is the route out of an assignment record that a retirement cannot finish.
// It renders as a refusal's `next=` line and as a `list` help row, so the two surfaces
// cannot name different verbs.
type recoveryRoute struct {
	words []string
	path  string
	why   string
}

// recoverMissingTree owns the landed rule. A landed assignment leaves with the batch
// clean, which is the one route `list` already advertises for the whole set. Any other
// record needs its own release, so the operator reads the request token that opened it.
func recoverMissingTree(landed bool, request, path string) recoveryRoute {
	if landed {
		return recoveryRoute{words: []string{"worktree", "clean", "--landed"}, why: "clean landed assignments"}
	}
	return recoveryRoute{words: []string{"worktree", "release", "--request", request}, path: path, why: "release the assignment whose worktree tree is missing"}
}

// discardUnlanded is the route out of an assignment whose branch never landed. A release
// retains that tree every time, so its refusal and the `list` row name the clean that plans
// the branch discard; that plan still asks for its fingerprint. A path that no line can
// carry leaves the route to the assignment id.
func discardUnlanded(path, assignment string) recoveryRoute {
	words, why := []string{"worktree", "clean", "--discard-branch"}, "discard the unlanded cleanup-pending assignment"
	if !lineSafe(path) {
		return recoveryRoute{words: append(words, "--target", assignment), why: why}
	}
	return recoveryRoute{words: words, path: path, why: why}
}

// releaseNext is the route a retained release names. An unlanded branch never releases, so
// it names the discard clean; any other retain clears with its cause, then release runs again.
func releaseNext(reason CleanupReason, target, assignment string) string {
	if reason == ReasonUnmerged {
		return discardUnlanded(target, assignment).line()
	}
	if lineSafe(target) {
		return "bench worktree release --request <request> " + sanitize.ShellQuote(target)
	}
	return "bench worktree exec " + assignment + " -- bench worktree release --request <request> ."
}

// line renders the refusal's route. axi owns the quoting here too, so the `next=` line
// and the help row give the operator one pasteable spelling of the same path.
func (r recoveryRoute) line() string {
	line := "bench " + strings.Join(r.words, " ")
	if r.path != "" {
		line += " " + axi.ShellQuote(r.path)
	}
	return line
}

// action renders the same route as a help row. axi owns the quoting there, so the path
// passes through as a known argument.
func (r recoveryRoute) action() axi.Action {
	arguments := make([]axi.InvocationArgument, 0, len(r.words)+1)
	for _, word := range r.words {
		arguments = append(arguments, axi.KnownArgument(word))
	}
	if r.path != "" {
		arguments = append(arguments, axi.KnownArgument(r.path))
	}
	return axi.ExecutableInvocation(r.why, arguments...)
}

func listAssignmentRow(j joins, root string, assignment intent.Assignment, def string, defaultResolved bool) []any {
	return []any{assignment.ID, assignment.Label, assignment.RequestToken, string(assignment.State), "assignment", listTree(assignment.Worktree), listLease(assignment.Worktree), listLanded(root, assignment.Branch, def, defaultResolved), listIgnored(j, assignment.Worktree)}
}

func listTree(path string) string {
	if _, err := os.Stat(path); err != nil {
		return "missing"
	}
	return "present"
}

func listLease(path string) string {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "none"
		}
		return "unknown"
	}
	if listPathHasSpecialGitMetadata(path) {
		return "unknown"
	}
	lease, err := LeaseFile(path)
	if err != nil {
		return "unknown"
	}
	if _, err := os.Stat(lease); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "none"
		}
		return "unknown"
	}
	return string(ProbeLease(lease))
}

func listLanded(root, branch, def string, defaultResolved bool) any {
	if !defaultResolved || branch == "" {
		return "unknown"
	}
	landed, _, err := git.LandedInDefault(root, branch, def)
	if err != nil {
		return "unknown"
	}
	return landed
}

func listIgnored(j joins, path string) any {
	if _, err := os.Stat(path); err != nil {
		return "unknown"
	}
	if listPathHasSpecialGitMetadata(path) {
		return "unknown"
	}
	inventory, _, err := inventoryIgnored(j, path, true)
	if err != nil || inventory.Uncertain {
		return "unknown"
	}
	return inventory.Count
}

func listPathHasSpecialGitMetadata(path string) bool {
	shape, err := ClassifyPathShape(path)
	return err == nil && shape == ShapeSpecialMetadata
}
