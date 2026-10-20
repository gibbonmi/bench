package recordcmd_test

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/refusalroute/routetest"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordcmd"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// declaredLine is a declaration of one field that no ledger type holds.
const declaredLine = "Persisted fields: intent-ledger:commitment.fixture\n"

// declaredBoundBuild is the bound build whose spec at the build checkout's HEAD declares
// declaration. The primary checkout declares Bench build inputs, and its broker manifest
// binds an installed candidate that is not this test binary.
func declaredBoundBuild(t *testing.T, declaration string) *recordtest.Fixture {
	t.Helper()
	f, primary := boundBuild(t)
	declare(t, f, declaration)
	commitmenttest.Write(t, primary, freshness.BuildInputsManifest, "cmd\n")
	commitmenttest.Write(t, primary, "dist/bench", "installed candidate\n")
	commitmenttest.BindManifest(t, primary, freshness.PublishedExecutable(primary))
	return f
}

// declare commits the spec of f with declaration in its header.
func declare(t *testing.T, f *recordtest.Fixture, declaration string) {
	t.Helper()
	spec := string(bytesOf(t, specFile(f)))
	f.Write(recordtest.Spec, strings.Replace(spec, "Status: staged\n", "Status: staged\n"+declaration, 1))
	f.Commit("declare a persisted field")
}

// refusedRerun runs args in the checkout of f and requires one next= route that names no
// gate. It answers the last step of that route.
func refusedRerun(t *testing.T, f *recordtest.Fixture, args []string) string {
	t.Helper()
	out, _ := recordcmd.Command(f.Root, args)
	next, printed := routetest.Next(out)
	steps := refusalroute.Steps(next)
	if !printed || strings.Contains(next, "bench gate") {
		t.Fatalf("assignment = %q, want one next= route that names no gate", out)
	}
	return steps[len(steps)-1]
}

// TestAssignmentPersistedFields is CB31 and CB61. A declared dispatch on an executable that
// is not the installed candidate refuses before its record read and its gap close, so the
// spec and the ledger keep their bytes. The declaration comes from HEAD, not from the
// working tree. The route ends with the caller's own call in its worktree, with the caller's
// flags in the form's layout order, each value quoted, and names no gate. A malformed
// declaration refuses with its own error and no route, except in a repository without
// Bench build inputs, where the entry reads no declaration.
func TestAssignmentPersistedFields(t *testing.T) {
	t.Run("uninstalled", func(t *testing.T) {
		f := declaredBoundBuild(t, declaredLine)
		args := assignArgs(t, *f.Plan.Execution, "1.md")
		commitmenttest.RequireLedgerKept(t, f.Root, func() { refuseAssignment(t, f, 1, "candidate-not-installed", args) })
		want := "bench record --in 'build' assignment 'example' --source 'HEAD' --ticket '1.md' --performer 'fixture-author-1.md' --assignment 'dispatch-1.md'" +
			" --model 'unknown' --effort 'unknown' --ref 'native:dispatch-1.md' --run 'fixture-run' --orchestrator 'fixture-orchestrator' --author-limit '2'"
		if got := refusedRerun(t, f, args); got != want {
			t.Fatalf("rerun:\n got %q\nwant %q", got, want)
		}
	})
	t.Run("hostile slug and flag value", func(t *testing.T) {
		f := declaredBoundBuild(t, declaredLine)
		f.Write("specs/it's here/spec.md", string(bytesOf(t, specFile(f))))
		f.Commit("copy the declared spec to a hostile slug")
		args := assignArgs(t, *f.Plan.Execution, "1.md", map[string]string{"--ref": "native:it's here"})
		args[1] = "it's here"
		want := `bench record --in 'build' assignment 'it'\''s here' --source 'HEAD' --ticket '1.md' --performer 'fixture-author-1.md' --assignment 'dispatch-1.md'` +
			` --model 'unknown' --effort 'unknown' --ref 'native:it'\''s here' --run 'fixture-run' --orchestrator 'fixture-orchestrator' --author-limit '2'`
		if got := refusedRerun(t, f, args); got != want {
			t.Fatalf("rerun:\n got %q\nwant %q", got, want)
		}
	})
	t.Run("an undeclared working tree", func(t *testing.T) {
		f := declaredBoundBuild(t, declaredLine)
		spec := string(bytesOf(t, specFile(f)))
		f.Write(recordtest.Spec, strings.Replace(spec, declaredLine, "", 1))
		refuseAssignment(t, f, 1, "candidate-not-installed", assignArgs(t, *f.Plan.Execution, "1.md"))
	})
	t.Run("malformed declaration", func(t *testing.T) {
		f := declaredBoundBuild(t, "Persisted fields: gate-evidence:fixture\n")
		out, code := recordcmd.Command(f.Root, assignArgs(t, *f.Plan.Execution, "1.md"))
		if _, printed := routetest.Next(out); code != 1 || printed || !strings.Contains(out, "persisted-fields-invalid") {
			t.Fatalf("assignment = (%d, %q), want the persisted-fields-invalid error and no route", code, out)
		}
	})
	t.Run("malformed declaration without build inputs", func(t *testing.T) {
		f, _ := boundBuild(t)
		declare(t, f, "Persisted fields: gate-evidence:fixture\n")
		if out, code := recordcmd.Command(f.Root, assignArgs(t, *f.Plan.Execution, "1.md")); code != 0 {
			t.Fatalf("assignment = (%d, %q), want a dispatch that reads no declaration", code, out)
		}
	})
}
