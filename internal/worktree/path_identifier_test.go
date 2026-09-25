package worktree

import (
	"bytes"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/intent"
)

// TestListActiveRowsUseTargetSlot shows that the active-row actions print once, however
// many active rows the list holds. The id cell fills the `<target>` slot, so no help row
// repeats an id.
func TestListActiveRowsUseTargetSlot(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	for _, name := range []string{"one", "two", "three"} {
		creation := mustCreate(t, root, home, "slot-"+name, "slot "+name)
		// A commit keeps the assignment unlanded, so the list adds no batch-clean row.
		commitInWorktree(t, creation.Path, name+".txt", name+"\n", name)
	}
	out, code := ListCommand(root, home, nil)
	if code != 0 {
		t.Fatalf("ListCommand = (%d, %q), want exit 0", code, out)
	}
	document, err := axitest.DecodeDocument(out)
	if err != nil {
		t.Fatalf("decode list: %v\n%s", err, out)
	}
	rows, err := document.Rows("worktrees")
	if err != nil || len(rows) != 3 {
		t.Fatalf("worktrees rows = %#v, %v, want three rows", rows, err)
	}
	var ids []string
	for _, raw := range rows {
		row, _ := raw.(map[string]any)
		if row["state"] != string(intent.StateActive) || row["tree"] != "present" {
			t.Fatalf("row = %#v, want an active row with a present tree", row)
		}
		id, _ := row["id"].(string)
		ids = append(ids, id)
	}
	actions, err := document.HelpActions()
	if err != nil {
		t.Fatal(err)
	}
	for _, action := range actions {
		for _, id := range ids {
			if strings.Contains(action.Cmd, id) {
				t.Fatalf("help row %q names the active id %q:\n%s", action.Cmd, id, out)
			}
		}
	}
	want := "help[2]{cmd,why}:\n" + activeHelpRows
	if !strings.HasSuffix(out, want) {
		t.Fatalf("ListCommand = %q, want the help block %q", out, want)
	}
}

// The list advertises `bench worktree path <target>` once, and the id cell of an active
// row fills the slot. So the id cell has to be the address the resolver accepts. The id
// fills the slot rather than the label, because ids are unique and labels can collide.
// It is the only address an agent has for a worktree it did not create.
func TestListPathActionRunsAsAdvertised(t *testing.T) {
	root, creation, home := newOwnedAssignment(t, "advertised")
	chdir(t, root)
	listed, code := ListCommand(root, home, nil)
	if code != 0 {
		t.Fatalf("list code=%d out=%q", code, listed)
	}
	document, err := axitest.DecodeDocument(listed)
	if err != nil {
		t.Fatalf("decode list: %v\n%s", err, listed)
	}
	actions, err := document.HelpActions()
	if err != nil {
		t.Fatal(err)
	}
	if !slices.ContainsFunc(actions, func(action axitest.HelpAction) bool { return action.Cmd == "bench worktree path <target>" }) {
		t.Fatalf("list advertised no target-slot path action:\n%s", listed)
	}
	rows, err := document.Rows("worktrees")
	if err != nil || len(rows) != 1 {
		t.Fatalf("worktrees rows = %#v, %v, want one row", rows, err)
	}
	row, _ := rows[0].(map[string]any)
	id, _ := row["id"].(string)
	if id != creation.Assignment.ID {
		t.Fatalf("id cell = %q, want the assignment id %q", id, creation.Assignment.ID)
	}
	var stdout, stderr bytes.Buffer
	if code := PathCommand(root, home, []string{id}, &stdout, &stderr); code != 0 {
		t.Fatalf("id cell %q exited %d: %s", id, code, stderr.String())
	}
}

// [PB34] The path serves the file tools, and a delegate that pastes it into a shell step
// leaves the worktree boundary. The note says so on stderr, because stdout stays the one
// path line a `$(...)` capture reads. The note quotes the target as typed, so the
// suggested command is the one the caller can run.
func TestPathNotesTheFileToolRouteOnStderr(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "filetools")
	target := creation.Assignment.Label
	var stdout, stderr bytes.Buffer
	if code := PathCommand(root, home, []string{target}, &stdout, &stderr); code != 0 {
		t.Fatalf("path exited %d: %s", code, stderr.String())
	}
	printed := strings.TrimSuffix(stdout.String(), "\n")
	if strings.Contains(printed, "\n") || !filepath.IsAbs(printed) {
		t.Fatalf("stdout = %q, want one absolute path line alone", stdout.String())
	}
	want := "note: the path serves the file tools; run a shell step through bench worktree exec " + target + " -- <command>\n"
	if stderr.String() != want {
		t.Fatalf("stderr = %q, want %q", stderr.String(), want)
	}
}

// The label stays a valid address; accepting the id widens the grammar, it does not
// replace it.
func TestPathResolvesTheLabelAndTheIdAlike(t *testing.T) {
	t.Parallel()
	root, creation, home := newOwnedAssignment(t, "both")
	for _, target := range []string{creation.Assignment.ID, creation.Assignment.Label} {
		var byTarget, stderr bytes.Buffer
		if code := PathCommand(root, home, []string{target}, &byTarget, &stderr); code != 0 {
			t.Fatalf("target %q exited %d: %s", target, code, stderr.String())
		}
		if byTarget.Len() == 0 {
			t.Fatalf("target %q resolved to no path", target)
		}
	}
}
