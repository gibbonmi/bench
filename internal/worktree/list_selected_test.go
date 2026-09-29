package worktree

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/toon"
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
			out, code := ListCommand(root, "", nil)
			got := worktreeListResponse{Stdout: out, Exit: code}
			if want := captured.render(t, a); !reflect.DeepEqual(got, want) {
				t.Fatalf("bare output = %#v, want stored pre-change baseline %#v", got, want)
			}
		})
	}
}

func selectedRows(t *testing.T, out string) []any {
	t.Helper()
	document, err := axitest.DecodeDocument(out)
	mustNoError(t, err)
	rows, err := document.Rows("worktrees")
	mustNoError(t, err)
	return rows
}

func TestSelectedWorktreeFacts(t *testing.T) {
	t.Parallel()
	root, assignments := selectedFixture(t)
	args := []string{"--view", "paths"}
	var want []any
	for _, a := range assignments {
		args = append(args, "--target", a.ID)
		want = append(want, map[string]any{"target": a.ID, "id": a.ID, "path": a.Worktree, "state": string(a.State), "error": ""})
	}
	out, code := ListCommand(root, "", args)
	if code != 0 {
		t.Fatalf("selected facts exit=%d output=%q", code, out)
	}
	if got := selectedRows(t, out); !reflect.DeepEqual(got, want) {
		t.Fatalf("selected facts=%#v, want %#v", got, want)
	}
}

func TestSelectedWorktreePartialFailure(t *testing.T) {
	t.Parallel()
	root, assignments := selectedFixture(t)
	assignments[1].Label = assignments[0].Label
	mustNoError(t, intent.PutAssignment(root, assignments[1]))
	out, code := ListCommand(root, "", []string{"--view", "paths", "--target", "absent", "--target", assignments[2].ID, "--target", assignments[0].Label})
	if code != 1 {
		t.Fatalf("partial failure exit=%d output=%q", code, out)
	}
	rows := selectedRows(t, out)
	if len(rows) != 3 {
		t.Fatalf("partial results=%#v, want all three operands", rows)
	}
	for i, fragment := range map[int]string{0: "unassigned", 2: "ambiguous"} {
		row := rows[i].(map[string]any)
		if !strings.Contains(row["error"].(string), fragment) || row["id"] != "" {
			t.Fatalf("failed result=%#v, want %s", row, fragment)
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
	out, code := ListCommand(root, "", args)
	if code != 1 {
		t.Fatalf("alias selection exit=%d output=%q", code, out)
	}
	rows := selectedRows(t, out)
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
			out, code := ListCommand(root, "", []string{"--view", "paths", "--target", target, "--target", assignments[0].Label})
			if code != 1 {
				t.Fatalf("hostile target exit=%d output=%q", code, out)
			}
			rows := selectedRows(t, out)
			if len(rows) != 2 {
				t.Fatalf("hostile results=%#v", rows)
			}
			bad := rows[0].(map[string]any)
			if bad["error"] == "" || bad["id"] != "" {
				t.Fatalf("unsafe result=%#v", bad)
			}
			if strings.ContainsAny(target, "\x1b\t\n\r\x00") && (bad["target"] != "target-1" || len(bad["error"].(string)) > 80) {
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
	out, code := ListCommand(root, "", args)
	if code != 1 {
		t.Fatalf("position matrix exit=%d output=%q", code, out)
	}
	rows := selectedRows(t, out)
	wantTargets := []string{assignments[0].Label, "target-2", "target-4", assignments[1].ID, "target-6"}
	if len(rows) != len(wantTargets) {
		t.Fatalf("position results=%#v", rows)
	}
	for i, target := range wantTargets {
		row := rows[i].(map[string]any)
		if row["target"] != target {
			t.Fatalf("position %d = %#v, want %q", i, row, target)
		}
		if strings.HasPrefix(target, "target-") && (row["id"] != "" || row["error"] == "" || len(row["error"].(string)) > 80) {
			t.Fatalf("position refusal=%#v", row)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "sentinel")); !os.IsNotExist(err) {
		t.Fatalf("operand executed: %v", err)
	}
}

func TestSelectedWorktreeGrammar(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{
		{"--view"}, {"--target"}, {"--view", "paths"}, {"--target", "one"},
		{"--view", "other", "--target", "one"}, {"--view", "paths", "--target", ""},
		{"--view", "paths", "--target", "one", "extra"},
		{"--view", "paths", "--target", "one", "--limit", "1"},
		{"--view", "paths", "--target", "one", "--view", "paths"},
		{"--view", "paths", "--target", "one", "--"},
	} {
		out, code := ListCommand("", "", args)
		if code != 2 || !strings.HasPrefix(out, "usage:") {
			t.Errorf("grammar %q = (%d,%q), want usage before repository lookup", args, code, out)
		}
	}
	for _, help := range []string{"--help", "-h", "help"} {
		out, code := ListCommand("", "", []string{help})
		if code != 0 || out != "usage: bench worktree list\n" {
			t.Errorf("bare help %q = (%d,%q)", help, code, out)
		}
	}
	out, code := ListCommand("", "", []string{"--view", "paths", "--help"})
	if code != 0 || out != selectedWorktreeGrammar.Help+"\n" {
		t.Fatalf("selected help=(%d,%q)", code, out)
	}
}

func TestSelectedWorktreesExcludeOldOutput(t *testing.T) {
	t.Parallel()
	root, assignments := selectedFixture(t)
	out, code := ListCommand(root, "", []string{"--view", "paths", "--target", assignments[0].ID})
	if code != 0 {
		t.Fatalf("selected output=(%d,%q)", code, out)
	}
	document, err := axitest.DecodeDocument(out)
	mustNoError(t, err)
	if !reflect.DeepEqual(document.Blocks, []string{"worktrees", "help"}) {
		t.Fatalf("selected blocks=%q", document.Blocks)
	}
	rows := selectedRows(t, out)
	if len(rows) != 1 || len(rows[0].(map[string]any)) != 5 {
		t.Fatalf("selected projection=%#v", rows)
	}
	for _, a := range assignments[1:] {
		if strings.Contains(out, a.ID) || strings.Contains(out, a.Worktree) {
			t.Fatalf("unrequested worktree survives: %q", out)
		}
	}
}

func TestSelectedWorktreeDetailRoute(t *testing.T) {
	t.Parallel()
	root, assignments := selectedFixture(t)
	for _, target := range []string{assignments[0].ID, "absent"} {
		out, _ := ListCommand(root, "", []string{"--view", "paths", "--target", target})
		document, err := axitest.DecodeDocument(out)
		mustNoError(t, err)
		actions, err := document.HelpActions()
		mustNoError(t, err)
		if len(actions) != 1 || actions[0].Cmd != "bench worktree list" {
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
			out, code := ListCommand(root, "", []string{"--view", "paths", "--target", bad.Label, "--target", good.ID, "--target", bad.ID, "--target", bad.Label})
			wantCode := 1
			first := map[string]any{"target": bad.Label, "id": "", "path": "", "state": "", "error": "assignment path is not representable"}
			if tc.permitted {
				wantCode = 0
				first = map[string]any{"target": bad.Label, "id": bad.ID, "path": bad.Worktree, "state": string(bad.State), "error": ""}
			}
			if code != wantCode {
				t.Fatalf("stored path exit=%d output=%q, want exit %d", code, out, wantCode)
			}
			want := []any{first, map[string]any{"target": good.ID, "id": good.ID, "path": good.Worktree, "state": string(good.State), "error": ""}}
			if rows := selectedRows(t, out); !reflect.DeepEqual(rows, want) {
				t.Fatalf("stored path results=%#v, want %#v", rows, want)
			}
		})
	}
}
