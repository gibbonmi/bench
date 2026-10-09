package main

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"

	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/usage"
)

// TestNativeRecordCLIJourney drives the native record forms through the production
// dispatcher in a linked worktree, so the primary-checkout guard stays active.
func TestNativeRecordCLIJourney(t *testing.T) {
	f, primary := recordtest.NewLinked(t, 2, recordtest.Delegate, recordtest.Undispatched)
	slug, err := rr.Slug(recordtest.Spec)
	if err != nil {
		t.Fatal(err)
	}
	spec := filepath.Join(f.Root, filepath.FromSlash(recordtest.Spec))
	read := func() []byte {
		data, err := os.ReadFile(spec)
		if err != nil {
			t.Fatal(err)
		}
		return data
	}

	// TestHelpInventoryIsComplete pins the complete grammar in root help. The form help
	// and the family help must advertise that same grammar.
	form := runAXICommandAt(t, f.Root, []string{"record", "assignment", "--help"})
	grammar, ok := strings.CutPrefix(strings.TrimSuffix(form.stdout, "\n"), "usage: bench record assignment ")
	root := runAXICommandAt(t, f.Root, []string{"help"})
	if form.code != 0 || !ok || !strings.Contains(root.stdout, "\n  bench record [--in <label|primary>] assignment "+grammar+"  ") {
		t.Fatalf("record assignment --help = exit %d, stdout %q; want exit 0 and the grammar that root help advertises", form.code, form.stdout)
	}
	if family := runAXICommandAt(t, f.Root, []string{"record", "--help"}); family.code != 0 || !slices.Contains(strings.Split(family.stdout, "\n"), "usage: bench record assignment "+grammar) {
		t.Fatalf("record --help = exit %d, stdout %q; want exit 0 with the assignment form", family.code, family.stdout)
	}

	execution := f.Plan.Execution
	dispatch := func(ticket string) []string {
		return []string{"record", "assignment", slug, "--source", "HEAD", "--ticket", ticket, "--performer", "native-author-1", "--assignment", "native-dispatch-1",
			"--model", "opus", "--effort", "high", "--ref", "native:dispatch-1", "--run", execution.RunID, "--orchestrator", execution.OrchestratorSession,
			"--author-limit", strconv.Itoa(execution.AuthorLimit)}
	}
	if result := runAXICommandAt(t, primary, dispatch("1.md")); result.code != 1 || !strings.Contains(result.stdout+result.stderr, usage.PrimaryCheckoutRefusal()) {
		t.Fatalf("primary checkout = exit %d, stdout %q, stderr %q; want exit 1 and the primary-checkout refusal", result.code, result.stdout, result.stderr)
	}
	source := f.Tip()
	if result := runAXICommandAt(t, f.Root, dispatch("1.md")); result.code != 0 {
		t.Fatalf("first dispatch = exit %d, stdout %q, stderr %q; want exit 0", result.code, result.stdout, result.stderr)
	}
	f.Commit("dispatch 1.md")
	f.Reload()
	want := []rr.Assignment{{Session: "native-author-1", Assignment: "native-dispatch-1", Model: "opus", Effort: "high", Source: source, NativeRef: "native:dispatch-1"}}
	if got := f.Plan.Execution.Assignments; !reflect.DeepEqual(got["1.md"], want) || len(got["2.md"]) != 0 {
		t.Fatalf("assignments = %+v, want 1.md %+v and 2.md empty", got, want)
	}

	// The same session cannot also author 2.md, and the full-plan validator refuses it
	// before the spec changes.
	before := read()
	if result := runAXICommandAt(t, f.Root, dispatch("2.md")); result.code != 1 || !strings.Contains(result.stdout, "already authors ticket 1.md") {
		t.Fatalf("cross-ticket dispatch = exit %d, stdout %q; want exit 1 naming the session's ticket", result.code, result.stdout)
	}
	if !bytes.Equal(read(), before) {
		t.Fatal("cross-ticket refusal changed the spec")
	}
}
