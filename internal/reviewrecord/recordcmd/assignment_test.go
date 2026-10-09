package recordcmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordcmd"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/toon"
)

// undispatched is a delegated fixture in a linked worktree whose count tickets have empty
// assignment histories. Each option follows Delegate.
func undispatched(t *testing.T, count int, options ...recordtest.Option) *recordtest.Fixture {
	t.Helper()
	f, _ := recordtest.NewLinked(t, count, append([]recordtest.Option{recordtest.Delegate, recordtest.Undispatched}, options...)...)
	return f
}

// assignArgs is the first-dispatch argv for ticket at HEAD that asserts execution.
func assignArgs(t *testing.T, execution rr.Execution, ticket string, edits ...map[string]string) []string {
	t.Helper()
	return formArgs(t, "assignment", map[string]string{"--source": "HEAD", "--ticket": ticket, "--performer": recordtest.Author(ticket),
		"--assignment": "dispatch-" + ticket, "--model": "unknown", "--effort": "unknown", "--ref": "native:dispatch-" + ticket,
		"--run": execution.RunID, "--orchestrator": execution.OrchestratorSession, "--author-limit": strconv.Itoa(execution.AuthorLimit)}, edits...)
}

// dispatch runs the assignment form for ticket over the plan of f and requires exit 0.
func dispatch(t *testing.T, f *recordtest.Fixture, ticket string, edits ...map[string]string) string {
	t.Helper()
	return succeed(t, f, assignArgs(t, *f.Plan.Execution, ticket, edits...))
}

func specFile(f *recordtest.Fixture) string {
	return filepath.Join(f.Root, filepath.FromSlash(recordtest.Spec))
}

// entries lists the names in dir.
func entries(t *testing.T, dir string) []string {
	t.Helper()
	list, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	names := []string{}
	for _, entry := range list {
		names = append(names, entry.Name())
	}
	return names
}

// refuseAssignment runs args over the spec of f and requires exit code, an output that
// holds want, unchanged spec bytes, and no new file beside the spec.
func refuseAssignment(t *testing.T, f *recordtest.Fixture, code int, want string, args []string) {
	t.Helper()
	before, names := bytesOf(t, specFile(f)), entries(t, filepath.Dir(specFile(f)))
	out, got := recordcmd.Command(f.Root, args)
	if got != code || !strings.Contains(out, want) {
		t.Fatalf("%s = exit %d, output %q; want exit %d naming %q", args[0], got, out, code, want)
	}
	if after := bytesOf(t, specFile(f)); !bytes.Equal(before, after) {
		t.Fatalf("refusal changed the spec")
	}
	if after := entries(t, filepath.Dir(specFile(f))); !reflect.DeepEqual(after, names) {
		t.Fatalf("spec directory = %q, want %q", after, names)
	}
}

// committedPlan commits the working tree of f and returns the plan that Reload reads at
// the new HEAD.
func committedPlan(t *testing.T, f *recordtest.Fixture, message string) rr.Plan {
	t.Helper()
	f.Commit(message)
	f.Reload()
	return f.Plan
}

// editTicket commits the prepared ticket name of f after edit rewrites it.
func editTicket(t *testing.T, f *recordtest.Fixture, name string, edit func(string) string) {
	t.Helper()
	path := f.TicketPath(name)
	f.Write(path, edit(string(bytesOf(t, filepath.Join(f.Root, filepath.FromSlash(path))))))
	f.Commit("edit ticket " + name)
}

func TestAssignmentFirstDispatch(t *testing.T) {
	f := undispatched(t, 2)
	source, before := f.Tip(), *f.Plan.Execution
	out := dispatch(t, f, "1.md")
	plan := committedPlan(t, f, "dispatch 1.md")
	want := []rr.Assignment{{Session: recordtest.Author("1.md"), Assignment: "dispatch-1.md", Model: "unknown", Effort: "unknown", Source: source, NativeRef: "native:dispatch-1.md"}}
	if got := plan.Execution.Assignments["1.md"]; !reflect.DeepEqual(got, want) {
		t.Fatalf("1.md history = %+v, want %+v", got, want)
	}
	if got := plan.Execution.Assignments["2.md"]; len(got) != 0 {
		t.Fatalf("2.md history = %+v, want it still empty", got)
	}
	after := *plan.Execution
	before.Assignments, after.Assignments = nil, nil
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("execution = %+v, want the source declaration %+v", after, before)
	}
	row, err := toon.Table("assignment", []string{"ticket", "session", "occurrence", "predecessor", "source"}, [][]string{{"1.md", recordtest.Author("1.md"), "1", "", source}})
	if err != nil {
		t.Fatal(err)
	}
	if out != row {
		t.Fatalf("output = %q, want %q", out, row)
	}
}

func TestAssignmentExecutionIdentity(t *testing.T) {
	f := undispatched(t, 1)
	for name, edit := range map[string]map[string]string{
		"run":          {"--run": "other-run"},
		"orchestrator": {"--orchestrator": "other-orchestrator"},
		"author limit": {"--author-limit": strconv.Itoa(f.Plan.Execution.AuthorLimit + 1)},
		"review mode":  {"--review-mode": "unified"},
	} {
		t.Run(name, func(t *testing.T) {
			refuseAssignment(t, f, 1, name, assignArgs(t, *f.Plan.Execution, "1.md", edit))
		})
	}
	for _, limit := range []string{"0", "-1"} {
		t.Run("non-positive author limit "+limit, func(t *testing.T) {
			refuseAssignment(t, f, 2, "usage: bench record assignment", assignArgs(t, *f.Plan.Execution, "1.md", map[string]string{"--author-limit": limit}))
		})
	}
	t.Run("omitted review mode asserts the default", func(t *testing.T) {
		unified := undispatched(t, 1, func(plan *rr.Plan) { plan.Execution.ReviewMode = "unified" })
		refuseAssignment(t, unified, 1, "review mode", assignArgs(t, *unified.Plan.Execution, "1.md"))
		dispatch(t, unified, "1.md", map[string]string{"--review-mode": "unified"})
	})
}

func TestAssignmentWholePlanValidation(t *testing.T) {
	t.Run("invalid non-target ticket", func(t *testing.T) {
		f := undispatched(t, 2)
		editTicket(t, f, "2.md", func(text string) string { return strings.Replace(text, "Writes: source.txt\n", "", 1) })
		refuseAssignment(t, f, 1, "2.md", assignArgs(t, *f.Plan.Execution, "1.md"))
	})
	t.Run("missing dependency", func(t *testing.T) {
		f := undispatched(t, 2)
		editTicket(t, f, "2.md", func(text string) string { return strings.Replace(text, "Blocked by: 1.md", "Blocked by: 9.md", 1) })
		refuseAssignment(t, f, 1, "9.md", assignArgs(t, *f.Plan.Execution, "1.md"))
	})
	t.Run("cyclic dependency", func(t *testing.T) {
		f, _ := recordtest.NewLinked(t, 2, oneChunk, recordtest.Delegate, recordtest.Undispatched)
		editTicket(t, f, "1.md", func(text string) string { return strings.Replace(text, "Blocked by: none", "Blocked by: 2.md", 1) })
		refuseAssignment(t, f, 1, "cycle", assignArgs(t, *f.Plan.Execution, "1.md"))
	})
	t.Run("session of another ticket", func(t *testing.T) {
		f := undispatched(t, 2)
		dispatch(t, f, "1.md")
		f.Commit("dispatch 1.md")
		refuseAssignment(t, f, 1, "already authors ticket 1.md", assignArgs(t, *f.Plan.Execution, "2.md", map[string]string{"--performer": recordtest.Author("1.md")}))
	})
	t.Run("version 1 plan", func(t *testing.T) {
		f, _ := recordtest.NewLinked(t, 1)
		refuseAssignment(t, f, 1, "version 2", assignArgs(t, rr.Execution{RunID: "run", OrchestratorSession: "orchestrator", AuthorLimit: 1}, "1.md"))
	})
}

func TestAssignmentSourceRefusals(t *testing.T) {
	t.Run("source older than HEAD", func(t *testing.T) {
		f := undispatched(t, 1)
		old, _ := advance(f, "later")
		refuseAssignment(t, f, 1, "HEAD", assignArgs(t, *f.Plan.Execution, "1.md", map[string]string{"--source": old}))
	})
	t.Run("changed working spec", func(t *testing.T) {
		f := undispatched(t, 1)
		f.Write(recordtest.Spec, specText(f)+"\nAn uncommitted line.\n")
		refuseAssignment(t, f, 1, recordtest.Spec, assignArgs(t, *f.Plan.Execution, "1.md"))
	})
	t.Run("changed planned ticket", func(t *testing.T) {
		f := undispatched(t, 2)
		path := f.TicketPath("2.md")
		f.Write(path, string(bytesOf(t, filepath.Join(f.Root, filepath.FromSlash(path))))+"\nAn uncommitted line.\n")
		refuseAssignment(t, f, 1, path, assignArgs(t, *f.Plan.Execution, "1.md"))
	})
}

func TestAssignmentReplacementHistory(t *testing.T) {
	f := undispatched(t, 1)
	dispatch(t, f, "1.md")
	first := committedPlan(t, f, "dispatch 1.md").Execution.Assignments["1.md"]
	replacement := map[string]string{"--performer": "replacement-author", "--assignment": "dispatch-2", "--ref": "native:dispatch-2",
		"--predecessor": recordtest.Author("1.md"), "--trigger": "terminal-failure", "--stopped": "native:stopped-1", "--preserved": "refs/bench/reset/1"}
	for name, refusal := range map[string]struct {
		code  int
		want  string
		edits []map[string]string
	}{
		"no replacement fields":       {1, "predecessor", nil},
		"another predecessor":         {1, "predecessor", []map[string]string{replacement, {"--predecessor": "someone-else"}}},
		"predecessor as performer":    {1, "distinct preceding author", []map[string]string{replacement, {"--performer": recordtest.Author("1.md")}}},
		"missing stopped evidence":    {2, "usage: bench record assignment", []map[string]string{replacement, {"--stopped": drop}}},
		"missing preserved work":      {2, "usage: bench record assignment", []map[string]string{replacement, {"--preserved": drop}}},
		"no progress, no reassessing": {1, "reassessment", []map[string]string{replacement, {"--trigger": "no-progress"}}},
		"reassessing another trigger": {1, "reassessment", []map[string]string{replacement, {"--reassessment": "native:reassessed"}}},
	} {
		t.Run(name, func(t *testing.T) {
			refuseAssignment(t, f, refusal.code, refusal.want, assignArgs(t, *f.Plan.Execution, "1.md", append([]map[string]string{{"--performer": "replacement-author"}}, refusal.edits...)...))
		})
	}
	source := f.Tip()
	dispatch(t, f, "1.md", replacement, map[string]string{"--trigger": "no-progress", "--reassessment": "native:reassessed"})
	history := committedPlan(t, f, "replace the author of 1.md").Execution.Assignments["1.md"]
	want := rr.Assignment{Session: "replacement-author", Assignment: "dispatch-2", Model: "unknown", Effort: "unknown", Source: source, NativeRef: "native:dispatch-2",
		Predecessor: recordtest.Author("1.md"), Trigger: "no-progress", Stopped: "native:stopped-1", Preserved: "refs/bench/reset/1", Reassessment: "native:reassessed"}
	if len(history) != 2 || !reflect.DeepEqual(history[:1], first) || !reflect.DeepEqual(history[1], want) {
		t.Fatalf("1.md history = %+v, want %+v followed by %+v", history, first, want)
	}
}

func TestAssignmentPreservesDocument(t *testing.T) {
	f := undispatched(t, 1)
	const fence, closing = "```bench-completion-plan\n", "\n```\n"
	f.Write(recordtest.Spec, "Préface — ✓\r\nA CRLF line.\r\n"+specText(f)+"\r\n## Après\r\n\r\nTrailing prose ✓")
	f.Commit("frame the plan")
	frame := func(document []byte) (string, string, string) {
		head, rest, ok := strings.Cut(string(document), fence)
		payload, tail, closed := strings.Cut(rest, closing)
		if !ok || !closed {
			t.Fatalf("spec %q holds no plan fence", document)
		}
		return head, payload, tail
	}
	head, payload, tail := frame(bytesOf(t, specFile(f)))
	source := f.Plan
	dispatch(t, f, "1.md")
	gotHead, gotPayload, gotTail := frame(bytesOf(t, specFile(f)))
	if gotHead != head || gotTail != tail || gotPayload == payload {
		t.Fatalf("spec frame = %q / %q, want %q / %q around a changed payload", gotHead, gotTail, head, tail)
	}
	plan := committedPlan(t, f, "dispatch 1.md")
	if !reflect.DeepEqual(plan.Chunks, source.Chunks) || !reflect.DeepEqual(plan.FinalVerification, source.FinalVerification) || len(plan.Execution.Assignments["1.md"]) != 1 {
		t.Fatalf("plan = %+v, want the source chunks and final verification with one 1.md assignment", plan)
	}
}

func TestAssignmentWriteFailure(t *testing.T) {
	f := undispatched(t, 1)
	dir := filepath.Dir(specFile(f))
	if err := os.Chmod(dir, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })
	probe := filepath.Join(dir, "probe")
	if err := os.WriteFile(probe, nil, 0o644); err == nil {
		_ = os.Remove(probe)
		capability.Capability(t, capability.Privilege, "the test process writes into a mode 0555 directory, so a failed temporary write is unobservable")
	}
	refuseAssignment(t, f, 1, "", assignArgs(t, *f.Plan.Execution, "1.md"))
}
