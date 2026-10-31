package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/responsebound/responseboundtest"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/usage"
	"github.com/gibbonmi/bench/internal/worktree"
)

// putSelectedAssignment records one complete assignment whose tree is absent, so the
// selected view answers a non-active record that the active-path resolver refuses.
func putSelectedAssignment(t *testing.T, root, id, owner string) intent.Assignment {
	t.Helper()
	assignment := intent.Assignment{
		Schema: intent.AssignmentRecordSchema, ID: id, OwnerID: owner,
		Request: intent.RequestDigest("selected-route-" + id), RequestToken: "selected-route-" + id,
		Label: "selected route " + id[:1], Start: gittest.Output(t, root, "rev-parse", "HEAD"),
		Branch: intent.AssignmentBranchRef(owner, id), Worktree: filepath.Join(root, "missing tree "+id[:1]), State: intent.StateComplete,
	}
	if err := intent.PutAssignment(root, assignment); err != nil {
		t.Fatal(err)
	}
	return assignment
}

func TestSelectedWorktreeCommandRoutes(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	assignment := putSelectedAssignment(t, root, strings.Repeat("a", 32), strings.Repeat("b", 32))
	for _, executable := range []string{"bench", filepath.Join(root, ".bench", "bin", "bench")} {
		for _, target := range []string{"bad\x1btarget", "space target", "--target", "$(touch sentinel)"} {
			result := runAXICommandAsAt(t, filepath.Join(root, "nested", "deep"), executable, []string{"worktree", "list", "--view", "paths", "--target", target, "--target", assignment.ID})
			if result.code != 1 || result.stderr != "" {
				t.Fatalf("selected route=%#v", result)
			}
			document, err := axitest.DecodeDocument(result.stdout)
			if err != nil {
				t.Fatal(err)
			}
			rows, err := document.Rows("worktrees")
			if err != nil {
				t.Fatal(err)
			}
			if len(rows) != 2 {
				t.Fatalf("selected rows=%#v", rows)
			}
			row := rows[1].(map[string]any)
			if row["id"] != assignment.ID || row["path"] != assignment.Worktree || row["state"] != "complete" || row["error"] != "" {
				t.Fatalf("non-active route=%#v", row)
			}
		}
	}
}

func TestSelectedWorktreeHelpDiscovery(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	result := runAXICommandAt(t, root, []string{"worktree", "--help"})
	if result.code != 0 || !strings.Contains(result.stdout, usage.WorktreeListPaths) {
		t.Fatalf("worktree help=%#v", result)
	}
	result = runAXICommandAt(t, root, []string{"worktree", "list", "--view", "paths", "--help"})
	if result.code != 0 || result.stdout != "usage: "+usage.WorktreeListPaths+"\n" {
		t.Fatalf("selected help=%#v", result)
	}
}

// TestSelectedWorktreeWithinResponseBound reads the printed response, not the spill file.
// Three resolved targets give 1 table header, 3 rows, 1 help header, and 1 help row: 6 lines.
func TestSelectedWorktreeWithinResponseBound(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	owner := strings.Repeat("b", 32)
	argv := []string{"worktree", "list", "--view", "paths"}
	for _, id := range []string{strings.Repeat("a", 32), strings.Repeat("c", 32), strings.Repeat("d", 32)} {
		argv = append(argv, "--target", putSelectedAssignment(t, root, id, owner).ID)
	}
	putSelectedAssignment(t, root, strings.Repeat("e", 32), owner)
	t.Chdir(root)
	var stdout, stderr bytes.Buffer
	code := Command{Stdout: &stdout, Stderr: &stderr, Executable: "bench"}.Run(argv)
	printed := stdout.String()
	if code != 0 || stderr.Len() != 0 {
		t.Fatalf("selected response exit=%d stdout=%q stderr=%q", code, printed, stderr.String())
	}
	if _, spilled := responseboundtest.Find(printed); spilled {
		t.Fatalf("selected response spilled: %q", printed)
	}
	if lines := strings.Count(printed, "\n"); lines != 6 || !strings.HasSuffix(printed, "\n") {
		t.Fatalf("selected response has %d lines, want exactly 6: %q", lines, printed)
	}
}

// retireSpecs commits count empty retire commits for each slug.
func retireSpecs(t *testing.T, root string, count int, slugs ...string) {
	t.Helper()
	for _, slug := range slugs {
		for i := 0; i < count; i++ {
			gittest.Commit(t, root, "spec-retire: "+slug, "--allow-empty")
		}
	}
}

func TestSelectedHistoryHelpDiscovery(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	result := runAXICommandAt(t, root, []string{"help"})
	if definition, _ := commandByName("spec"); result.code != 0 || !strings.Contains(result.stdout, "  bench spec"+definition.treeTargetInsertion()+strings.TrimPrefix(spec.SelectedHistoryUsage, "bench spec")+"  ") {
		t.Fatalf("selected history discovery=%#v", result)
	}
	result = runAXICommandAt(t, root, []string{"spec", "history", "--spec", "chosen", "--help"})
	if result.code != 0 || result.stdout != "usage: "+spec.SelectedHistoryUsage+"\n" {
		t.Fatalf("selected history help=%#v", result)
	}
}

func TestSelectedHistoryCommandRoutes(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	retireSpecs(t, root, 2, "chosen")
	for _, executable := range []string{"bench", filepath.Join(root, ".bench", "bin", "bench")} {
		for _, target := range []string{"bad\x1btarget", "a space", "--spec", "$(touch sentinel)", "quote'name", "café"} {
			result := runAXICommandAsAt(t, filepath.Join(root, "nested", "deep"), executable, []string{"spec", "history", "--spec", target, "--spec", "chosen", "--limit", "1"})
			unsafe := target == "bad\x1btarget"
			wantExit := 0
			if unsafe {
				wantExit = 1
			}
			if result.code != wantExit || result.stderr != "" {
				t.Fatalf("selected route=%#v", result)
			}
			document, err := axitest.DecodeDocument(result.stdout)
			if err != nil {
				t.Fatal(err)
			}
			rows, err := document.Rows("histories")
			if err != nil || len(rows) != 2 {
				t.Fatalf("summary=%#v error=%v", rows, err)
			}
			first := rows[0].(map[string]any)
			if unsafe {
				if first["target"] != sanitize.TargetPointer(1) || first["error"] == "" {
					t.Fatalf("unsafe route=%#v", first)
				}
			} else if first["target"] != target || first["error"] != "" || fmt.Sprint(first["total_events"]) != "0" {
				t.Fatalf("empty route=%#v", first)
			}
			row := rows[1].(map[string]any)
			if row["slug"] != "chosen" || row["error"] != "" || fmt.Sprint(row["total_events"]) != "2" || fmt.Sprint(row["omitted_events"]) != "1" {
				t.Fatalf("retained successful route=%#v", row)
			}
			events, err := document.Rows("history")
			if err != nil || len(events) != 1 || events[0].(map[string]any)["slug"] != "chosen" {
				t.Fatalf("selected events=%#v error=%v", events, err)
			}
		}
	}
	if _, err := os.Stat(filepath.Join(root, "nested", "deep", "sentinel")); !os.IsNotExist(err) {
		t.Fatalf("shell-shaped target executed: %v", err)
	}
}

// TestSelectedHistoryWithinResponseBound reads the printed response, not the spill file.
// Two 3-event specs at limit 2 give 1 summary header, 2 summary rows, 1 event header, and
// 2 × 2 event rows: 8 lines.
func TestSelectedHistoryWithinResponseBound(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	retireSpecs(t, root, 3, "first", "second")
	t.Chdir(root)
	var stdout, stderr bytes.Buffer
	code := Command{Stdout: &stdout, Stderr: &stderr, Executable: "bench"}.Run([]string{"spec", "history", "--spec", "first", "--spec", "second", "--limit", "2"})
	printed := stdout.String()
	if code != 0 || stderr.Len() != 0 {
		t.Fatalf("selected response exit=%d stdout=%q stderr=%q", code, printed, stderr.String())
	}
	if _, spilled := responseboundtest.Find(printed); spilled {
		t.Fatalf("selected response spilled: %q", printed)
	}
	if lines := strings.Count(withoutTreeRow(printed), "\n"); lines != 8 || !strings.HasSuffix(printed, "\n") {
		t.Fatalf("selected response has %d lines, want exactly 8: %q", lines, printed)
	}
}

// relatedID is the assignment identity of related fixture row i.
func relatedID(i int) string { return fmt.Sprintf("%032x", i) }

// putRelatedAssignment records assignment id with the given state and tree.
func putRelatedAssignment(t *testing.T, root, id string, state intent.AssignmentState, tree string) intent.Assignment {
	t.Helper()
	owner := strings.Repeat("b", 32)
	assignment := intent.Assignment{
		Schema: intent.AssignmentRecordSchema, ID: id, OwnerID: owner,
		Request: intent.RequestDigest("related-" + id), RequestToken: "related-" + id,
		Label: "related " + id[28:], Start: gittest.Output(t, root, "rev-parse", "HEAD"),
		Branch: intent.AssignmentBranchRef(owner, id), Worktree: tree, State: state,
	}
	if err := intent.PutAssignment(root, assignment); err != nil {
		t.Fatal(err)
	}
	return assignment
}

// relatedFixture records forty complete assignments with absent trees and one active
// assignment in a real linked worktree, all bound to specs/x/spec.md, and binds one more
// assignment that the ledger does not hold yet. The projection spans several pages.
func relatedFixture(t *testing.T) (string, intent.Assignment) {
	t.Helper()
	root := newAXIEnvelopeRepo(t)
	writeAXIFixture(t, filepath.Join(root, ".git", "info", "exclude"), "*.ignored\n")
	var bindings []intent.DeliveryBinding
	for i := 1; i <= 40; i++ {
		bindings = append(bindings, intent.DeliveryBinding{Assignment: putRelatedAssignment(t, root, relatedID(i), intent.StateComplete, filepath.Join(root, "absent", relatedID(i))).ID})
	}
	linked := filepath.Join(t.TempDir(), "linked")
	branch := intent.AssignmentBranchRef(strings.Repeat("b", 32), relatedID(100))
	gittest.MustRun(t, root, "worktree", "add", "-q", "-b", strings.TrimPrefix(branch, "refs/heads/"), linked)
	active := putRelatedAssignment(t, root, relatedID(100), intent.StateActive, linked)
	bindings = append(bindings, intent.DeliveryBinding{Assignment: active.ID}, intent.DeliveryBinding{Assignment: relatedID(200)})
	commitmenttest.EditCommitment(t, root, func(state *intent.CommitmentState) {
		for _, binding := range bindings {
			binding.Request, binding.Milestone, binding.Outcome, binding.Deliverable, binding.Identity = intent.RequestDigest("related-"+binding.Assignment), "fixture", "related-"+binding.Assignment, "specs/x/spec.md", "fixture"
			state.Bindings = append(state.Bindings, binding)
		}
	})
	return root, active
}

// relatedPage runs one related page from root and answers its source token, its
// continuation argv, and its completion.
func relatedPage(t *testing.T, root string, argv []string) (string, []string, bool) {
	t.Helper()
	result := runAXICommandAt(t, root, argv)
	document, err := axitest.DecodeDocument(result.stdout)
	if result.code != 0 || err != nil {
		t.Fatalf("%q = %#v, %v", argv, result, err)
	}
	rows, err := document.Rows("page")
	if err != nil {
		t.Fatal(err)
	}
	page := rows[0].(map[string]any)
	next, _ := axitest.RecoverHelpCommandArgv(result.stdout)
	if page["complete"] == true {
		return page["source"].(string), nil, true
	}
	return page["source"].(string), next[1:], false
}

// ER48: a related continuation refuses before content when its projection changed since
// the cursor while the commitment data stayed equal. Each case changes the assignment
// population, a rendered field, or an observed tree fact, and with it the rows, the counts,
// or the navigation actions.
func TestRelatedContinuationRefusesChangedProjection(t *testing.T) {
	transact := func(t *testing.T, root string, edit func([]intent.Assignment) []intent.Assignment) {
		if err := intent.Transact(root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
			ledger.Assignments = edit(ledger.Assignments)
			return ledger, true, nil
		}, nil); err != nil {
			t.Fatal(err)
		}
	}
	edit := func(change func(*intent.Assignment)) func(*testing.T, string, intent.Assignment) {
		return func(t *testing.T, root string, active intent.Assignment) {
			change(&active)
			if err := intent.PutAssignment(root, active); err != nil {
				t.Fatal(err)
			}
		}
	}
	for _, tc := range []struct {
		name   string
		change func(*testing.T, string, intent.Assignment)
	}{
		{"a bound assignment joins", func(t *testing.T, root string, _ intent.Assignment) {
			putRelatedAssignment(t, root, relatedID(200), intent.StateComplete, filepath.Join(root, "absent", "joined"))
		}},
		{"an unbound assignment joins", func(t *testing.T, root string, _ intent.Assignment) {
			putRelatedAssignment(t, root, relatedID(300), intent.StateComplete, filepath.Join(root, "absent", "unbound"))
		}},
		{"a bound assignment leaves", func(t *testing.T, root string, _ intent.Assignment) {
			transact(t, root, func(assignments []intent.Assignment) []intent.Assignment { return assignments[1:] })
		}},
		{"a label changes", edit(func(a *intent.Assignment) { a.Label = "renamed" })},
		{"a request changes", edit(func(a *intent.Assignment) { a.RequestToken, a.Request = "renamed", intent.RequestDigest("renamed") })},
		{"a state changes", edit(func(a *intent.Assignment) { a.State = intent.StateCleanupPending })},
		{"a path changes", edit(func(a *intent.Assignment) { a.Worktree += "-moved" })},
		{"a tree disappears", func(t *testing.T, _ string, active intent.Assignment) {
			if err := os.RemoveAll(active.Worktree); err != nil {
				t.Fatal(err)
			}
		}},
		{"a lease appears", func(t *testing.T, _ string, active intent.Assignment) {
			lease, err := worktree.LeaseFile(active.Worktree)
			if err != nil {
				t.Fatal(err)
			}
			writeAXIFixture(t, lease, "not a lease\n")
		}},
		{"the branch stops being landed", func(t *testing.T, _ string, active intent.Assignment) {
			writeAXIFixture(t, filepath.Join(active.Worktree, "own.txt"), "own\n")
			gittest.MustRun(t, active.Worktree, "add", "own.txt")
			gittest.Commit(t, active.Worktree, "own commit")
		}},
		{"an ignored file appears", func(t *testing.T, _ string, active intent.Assignment) {
			writeAXIFixture(t, filepath.Join(active.Worktree, "scratch.ignored"), "ignored\n")
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, active := relatedFixture(t)
			first := []string{"worktree", "list", "--spec", "x"}
			token, next, complete := relatedPage(t, root, first)
			if complete {
				t.Fatal("the related fixture fits one page, want a continuation")
			}
			before, err := intent.Read(root)
			if err != nil {
				t.Fatal(err)
			}
			tc.change(t, root, active)
			after, err := intent.Read(root)
			if err != nil {
				t.Fatal(err)
			}
			if was, is := fmt.Sprintf("%+v", *before.Commitment), fmt.Sprintf("%+v", *after.Commitment); was != is {
				t.Fatalf("commitment changed: %s -> %s", was, is)
			}
			if fresh, _, _ := relatedPage(t, root, first); fresh == token {
				t.Fatal("the change left the related projection equal")
			}
			result := runAXICommandAt(t, root, next)
			if result.code != 1 || !strings.HasPrefix(result.stdout, "error: source changed since the cursor") || strings.Contains(result.stdout, "content[") {
				t.Fatalf("continuation = %#v, want the source-changed refusal before content", result)
			}
		})
	}
}

// The related view refuses the selected-view flags, and its help names its own grammar.
func TestRelatedWorktreeHelpAndRefusals(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	for _, tc := range []struct {
		argv []string
		code int
		want string
	}{
		{[]string{"worktree", "--help"}, 0, usage.WorktreeListRelated + "\n"},
		{[]string{"worktree", "list", "--spec", "x", "--help"}, 0, "usage: " + usage.WorktreeListRelated + "\n"},
		{[]string{"worktree", "list", "--spec", "x", "--view", "paths"}, 2, "--spec"},
		{[]string{"worktree", "list", "--spec", "x", "--target", "a", "--target", "b"}, 2, "--spec"},
		{[]string{"worktree", "list", "--spec", "x", "--offset", "0"}, 2, "usage: " + usage.WorktreeListRelated + "\n"},
		{[]string{"worktree", "list", "--spec", "specs/x"}, 1, "error: specs/x is not a spec path"},
	} {
		if result := runAXICommandAt(t, root, tc.argv); result.code != tc.code || !strings.Contains(result.stdout, tc.want) {
			t.Fatalf("%q = %#v, want exit %d and %q", tc.argv, result, tc.code, tc.want)
		}
	}
}
