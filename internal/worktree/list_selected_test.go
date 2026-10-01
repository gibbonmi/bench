package worktree

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

func selectedFixture(t *testing.T, states ...intent.AssignmentState) (string, []intent.Assignment) {
	t.Helper()
	root := newWorktreeRepo(t)
	home := filepath.Join(t.TempDir(), "bench-home")
	var assignments []intent.Assignment
	if len(states) == 0 {
		states = []intent.AssignmentState{intent.StateActive, intent.StateComplete, intent.StateCleanupPending, intent.StateRecovered}
	}
	for _, state := range states {
		created := mustCreate(t, root, home, "request-"+string(state), "chosen "+string(state))
		assignment := created.Assignment
		assignment.State = state
		if state == intent.StateRecovered {
			assignment.Recovery = []intent.Recovery{{Ref: intent.RecoveryRefPrefix(assignment.OwnerID, assignment.ID) + "1", Root: assignment.Start, Payloads: []string{assignment.Start}}}
		}
		mustNoError(t, intent.PutAssignment(root, assignment))
		assignments = append(assignments, assignment)
	}
	return root, assignments
}

// selectedBaseline is one captured bare response. The tables hold the captured facts, and
// the TOON owner renders them, so the expectation cannot drift from the encoder's quoting.
type selectedBaseline struct {
	worktreeListResponse
	Tables []struct {
		Name   string
		Fields []string
		Rows   [][]any
	}
}

func (b selectedBaseline) render(t *testing.T, a intent.Assignment) worktreeListResponse {
	t.Helper()
	response := b.worktreeListResponse
	expand := strings.NewReplacer("{{PATH}}", a.Worktree, "{{ID}}", a.ID)
	for _, table := range b.Tables {
		for _, row := range table.Rows {
			for i, cell := range row {
				if value, ok := cell.(string); ok {
					row[i] = expand.Replace(value)
				}
			}
		}
		block, err := toon.TableTyped(table.Name, table.Fields, table.Rows)
		mustNoError(t, err)
		response.Stdout += block
	}
	return response
}

func TestSelectedWorktreesPreserveDefault(t *testing.T) {
	t.Parallel()
	data, err := os.ReadFile("testdata/selected-defaults.json")
	mustNoError(t, err)
	var baseline map[string]selectedBaseline
	mustNoError(t, json.Unmarshal(data, &baseline))
	for name, captured := range baseline {
		t.Run(name, func(t *testing.T) {
			root := ""
			var a intent.Assignment
			if name == "empty" {
				root = newWorktreeRepo(t)
			} else if name != "no-repository" {
				var assignments []intent.Assignment
				root, assignments = selectedFixture(t, intent.AssignmentState(name))
				a = assignments[0]
				if a.State == intent.StateComplete {
					// A numeric-looking identity makes a quoting change observable.
					mustNoError(t, intent.DeleteAssignment(root, a.ID))
					a.ID = strings.Repeat("0", 32)
					a.Branch = intent.AssignmentBranchRef(a.OwnerID, a.ID)
					gitOutput(t, root, "update-ref", a.Branch, a.Start)
					mustNoError(t, intent.PutAssignment(root, a))
				}
			}
			r := runVerb(t, verbList, verbCall{root: root})
			got := worktreeListResponse{Stdout: r.stdout, Exit: r.exit}
			if want := captured.render(t, a); !reflect.DeepEqual(got, want) {
				t.Fatalf("bare output = %#v, want stored pre-change baseline %#v", got, want)
			}
		})
	}
}

// selectedTable is the table block that the selected list renders its rows in.
const selectedTable = "worktrees"

func TestSelectedWorktreeFacts(t *testing.T) {
	t.Parallel()
	root, assignments := selectedFixture(t)
	args := []string{"--view", "paths"}
	var want []any
	for _, a := range assignments {
		args = append(args, "--target", a.ID)
		want = append(want, map[string]any{"target": a.ID, "id": a.ID, "path": a.Worktree, "state": string(a.State), "error": ""})
	}
	r := runVerb(t, verbList, verbCall{root: root, args: args})
	if r.exit != 0 {
		t.Fatalf("selected facts exit=%d output=%q", r.exit, r.stdout)
	}
	if got := r.mustRows(t, selectedTable); !reflect.DeepEqual(got, want) {
		t.Fatalf("selected facts=%#v, want %#v", got, want)
	}
}

func TestSelectedWorktreePartialFailure(t *testing.T) {
	t.Parallel()
	root, assignments := selectedFixture(t)
	assignments[1].Label = assignments[0].Label
	mustNoError(t, intent.PutAssignment(root, assignments[1]))
	r := runVerb(t, verbList, verbCall{root: root, args: []string{"--view", "paths", "--target", "absent", "--target", assignments[2].ID, "--target", assignments[0].Label}})
	if r.exit != 1 {
		t.Fatalf("partial failure exit=%d output=%q", r.exit, r.stdout)
	}
	rows := r.mustRows(t, selectedTable)
	if len(rows) != 3 {
		t.Fatalf("partial results=%#v, want all three operands", rows)
	}
	ledger, err := intent.Assignments(root)
	mustNoError(t, err)
	var ambiguous ambiguousTargetError
	_, selectErr := selectAssignment(ledger, assignments[0].Label)
	if !errors.As(selectErr, &ambiguous) {
		t.Fatalf("shared label selects %v, want an ambiguity", selectErr)
	}
	for i, want := range map[int]string{0: errTargetUnassigned.Error(), 2: ambiguous.Error()} {
		if row := rows[i].(map[string]any); row["error"] != want || row["id"] != "" {
			t.Fatalf("failed result=%#v, want error %q", row, want)
		}
	}
	if row := rows[1].(map[string]any); row["id"] != assignments[2].ID || row["error"] != "" {
		t.Fatalf("retained success=%#v", row)
	}
}

func TestSelectedWorktreeAliases(t *testing.T) {
	t.Parallel()
	root, assignments := selectedFixture(t)
	first, last := assignments[2], assignments[0]
	args := []string{"--view", "paths"}
	for _, target := range []string{first.Worktree, "missing", first.Label, first.ID, first.ID[:8], "missing", last.Label} {
		args = append(args, "--target", target)
	}
	r := runVerb(t, verbList, verbCall{root: root, args: args})
	if r.exit != 1 {
		t.Fatalf("alias selection exit=%d output=%q", r.exit, r.stdout)
	}
	rows := r.mustRows(t, selectedTable)
	if len(rows) != 3 {
		t.Fatalf("alias results=%#v, want two identities and one distinct failure", rows)
	}
	if row := rows[0].(map[string]any); row["id"] != first.ID || row["target"] != first.Worktree {
		t.Fatalf("first alias=%#v", row)
	}
	if row := rows[1].(map[string]any); row["target"] != "missing" {
		t.Fatalf("first failure=%#v", row)
	}
	if row := rows[2].(map[string]any); row["id"] != last.ID || row["target"] != last.Label {
		t.Fatalf("last identity=%#v", row)
	}
}

func TestSelectedWorktreeHostileTarget(t *testing.T) {
	t.Parallel()
	root, assignments := selectedFixture(t)
	for _, target := range []string{"bad\x1btarget", "tab\ttarget", "line\ntarget", "return\rtarget", "nul\x00target", "--unknown", "$(touch sentinel)", "a space"} {
		t.Run(target, func(t *testing.T) {
			r := runVerb(t, verbList, verbCall{root: root, args: []string{"--view", "paths", "--target", target, "--target", assignments[0].Label}})
			if r.exit != 1 {
				t.Fatalf("hostile target exit=%d output=%q", r.exit, r.stdout)
			}
			rows := r.mustRows(t, selectedTable)
			if len(rows) != 2 {
				t.Fatalf("hostile results=%#v", rows)
			}
			bad := rows[0].(map[string]any)
			if bad["error"] == "" || bad["id"] != "" {
				t.Fatalf("unsafe result=%#v", bad)
			}
			if strings.ContainsAny(target, "\x1b\t\n\r\x00") && (bad["target"] != sanitize.TargetPointer(1) || len(bad["error"].(string)) > 80) {
				t.Fatalf("unsafe ordinal result=%#v", bad)
			}
			if rows[1].(map[string]any)["id"] != assignments[0].ID {
				t.Fatalf("valid result lost: %#v", rows)
			}
		})
	}
	args := []string{"--view", "paths"}
	for _, target := range []string{assignments[0].Label, "bad\x1btarget", "bad\x1btarget", "line\ntarget", assignments[1].ID, "last\rtarget"} {
		args = append(args, "--target", target)
	}
	r := runVerb(t, verbList, verbCall{root: root, args: args})
	if r.exit != 1 {
		t.Fatalf("position matrix exit=%d output=%q", r.exit, r.stdout)
	}
	rows := r.mustRows(t, selectedTable)
	wantTargets := []string{assignments[0].Label, sanitize.TargetPointer(2), sanitize.TargetPointer(4), assignments[1].ID, sanitize.TargetPointer(6)}
	refused := map[int]bool{1: true, 2: true, 4: true}
	if len(rows) != len(wantTargets) {
		t.Fatalf("position results=%#v", rows)
	}
	for i, target := range wantTargets {
		row := rows[i].(map[string]any)
		if row["target"] != target {
			t.Fatalf("position %d = %#v, want %q", i, row, target)
		}
		if refused[i] && (row["id"] != "" || row["error"] == "" || len(row["error"].(string)) > 80) {
			t.Fatalf("position refusal=%#v", row)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "sentinel")); !os.IsNotExist(err) {
		t.Fatalf("operand executed: %v", err)
	}
}

func TestSelectedWorktreeGrammar(t *testing.T) {
	t.Parallel()
	// Each refusal is the selected grammar's own line. The bare grammar refuses the first
	// flag as an unknown argument, so a request that misses the selected route fails here.
	selectedHelp := selectedWorktreeGrammar.Help
	for _, tc := range []struct {
		args []string
		want string
	}{
		{[]string{"--view"}, toon.MissingArg(usage.WorktreeList, "--view")},
		{[]string{"--target"}, toon.MissingArg(usage.WorktreeList, "--target")},
		{[]string{"--view", "paths"}, selectedHelp},
		{[]string{"--target", "one"}, selectedHelp},
		{[]string{"--view", "other", "--target", "one"}, selectedHelp},
		{[]string{"--view", "paths", "--target", ""}, toon.Usage(usage.WorktreeList, usage.EmptyFlagValue("--target"))},
		{[]string{"--view", "paths", "--target", "one", "extra"}, toon.Usage(usage.WorktreeList, "extra")},
		{[]string{"--view", "paths", "--target", "one", "--limit", "1"}, toon.Usage(usage.WorktreeList, "--limit")},
		{[]string{"--view", "paths", "--target", "one", "--view", "paths"}, toon.Usage(usage.WorktreeList, "--view")},
		{[]string{"--view", "paths", "--target", "one", "--"}, selectedHelp},
	} {
		r := runVerb(t, verbList, verbCall{args: tc.args})
		if r.exit != 2 || r.stdout != tc.want+"\n" {
			t.Errorf("grammar %q = (%d,%q), want (2,%q) before repository lookup", tc.args, r.exit, r.stdout, tc.want+"\n")
		}
	}
	for _, help := range []string{"--help", "-h", "help"} {
		r := runVerb(t, verbList, verbCall{args: []string{help}})
		if r.exit != 0 || r.stdout != worktreeListGrammar.Help+"\n" {
			t.Errorf("bare help %q = (%d,%q)", help, r.exit, r.stdout)
		}
	}
	r := runVerb(t, verbList, verbCall{args: []string{"--view", "paths", "--help"}})
	if r.exit != 0 || r.stdout != selectedWorktreeGrammar.Help+"\n" {
		t.Fatalf("selected help=(%d,%q)", r.exit, r.stdout)
	}
}

func TestSelectedWorktreesExcludeOldOutput(t *testing.T) {
	t.Parallel()
	root, assignments := selectedFixture(t)
	r := runVerb(t, verbList, verbCall{root: root, args: []string{"--view", "paths", "--target", assignments[0].ID}})
	if r.exit != 0 {
		t.Fatalf("selected output=(%d,%q)", r.exit, r.stdout)
	}
	document, err := axitest.DecodeDocument(r.stdout)
	mustNoError(t, err)
	if !reflect.DeepEqual(document.Blocks, []string{selectedTable, "help"}) {
		t.Fatalf("selected blocks=%q", document.Blocks)
	}
	rows := r.mustRows(t, selectedTable)
	if len(rows) != 1 || len(rows[0].(map[string]any)) != 5 {
		t.Fatalf("selected projection=%#v", rows)
	}
	for _, a := range assignments[1:] {
		if strings.Contains(r.stdout, a.ID) || strings.Contains(r.stdout, a.Worktree) {
			t.Fatalf("unrequested worktree survives: %q", r.stdout)
		}
	}
}

func TestSelectedWorktreeDetailRoute(t *testing.T) {
	t.Parallel()
	root, assignments := selectedFixture(t)
	for _, target := range []string{assignments[0].ID, "absent"} {
		r := runVerb(t, verbList, verbCall{root: root, args: []string{"--view", "paths", "--target", target}})
		document, err := axitest.DecodeDocument(r.stdout)
		mustNoError(t, err)
		actions, err := document.HelpActions()
		mustNoError(t, err)
		if len(actions) != 1 || actions[0].Cmd != usage.WorktreeList {
			t.Fatalf("detail actions=%#v, want complete worktree inventory", actions)
		}
	}
}

func TestSelectedWorktreeHostilePath(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, suffix string
		permitted    bool
	}{{"escape", "\x1b", false}, {"tab", "\t", true}, {"newline", "\n", true}, {"return", "\r", true}} {
		t.Run(tc.name, func(t *testing.T) {
			root, assignments := selectedFixture(t)
			bad, good := assignments[0], assignments[1]
			bad.Worktree += tc.suffix
			mustNoError(t, intent.PutAssignment(root, bad))
			r := runVerb(t, verbList, verbCall{root: root, args: []string{"--view", "paths", "--target", bad.Label, "--target", good.ID, "--target", bad.ID, "--target", bad.Label}})
			wantCode := 1
			first := map[string]any{"target": bad.Label, "id": "", "path": "", "state": "", "error": selectedPathUnrepresentable}
			if tc.permitted {
				wantCode = 0
				first = map[string]any{"target": bad.Label, "id": bad.ID, "path": bad.Worktree, "state": string(bad.State), "error": ""}
			}
			if r.exit != wantCode {
				t.Fatalf("stored path exit=%d output=%q, want exit %d", r.exit, r.stdout, wantCode)
			}
			want := []any{first, map[string]any{"target": good.ID, "id": good.ID, "path": good.Worktree, "state": string(good.State), "error": ""}}
			if rows := r.mustRows(t, selectedTable); !reflect.DeepEqual(rows, want) {
				t.Fatalf("stored path results=%#v, want %#v", rows, want)
			}
		})
	}
}
