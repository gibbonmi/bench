package recordcmd_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/brokermanifest"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/refusalroute/routetest"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordcmd"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// declaredBoundBuild is the bound build whose spec at the build checkout's HEAD declares
// declaration. The primary checkout declares Bench build inputs, and its broker manifest
// binds an installed candidate that is not this test binary.
func declaredBoundBuild(t *testing.T, declaration string) *recordtest.Fixture {
	t.Helper()
	f, primary := boundBuild(t)
	spec := string(bytesOf(t, specFile(f)))
	f.Write(recordtest.Spec, strings.Replace(spec, "Status: staged\n", "Status: staged\n"+declaration, 1))
	f.Commit("declare a persisted field")
	commitmenttest.Write(t, primary, freshness.BuildInputsManifest, "cmd\n")
	commitmenttest.Write(t, primary, "dist/bench", "installed candidate\n")
	dir := filepath.Join(primary, "bin")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, _, err := brokermanifest.Write(dir, freshness.PublishedExecutable(primary), "0.0.0"); err != nil {
		t.Fatal(err)
	}
	return f
}

// TestAssignmentPersistedFields is CB31 and CB61. A declared dispatch on an executable that
// is not the installed candidate refuses before its record read and its gap close, so the
// spec and the ledger keep their bytes. The route ends with the caller's own call in its
// worktree, with the caller's flags in the form's layout order, and names no gate. A
// malformed declaration refuses with its own error and no route.
func TestAssignmentPersistedFields(t *testing.T) {
	t.Run("uninstalled", func(t *testing.T) {
		f := declaredBoundBuild(t, "Persisted fields: intent-ledger:commitment.fixture\n")
		args := assignArgs(t, *f.Plan.Execution, "1.md")
		commitmenttest.RequireLedgerKept(t, f.Root, func() { refuseAssignment(t, f, 1, "candidate-not-installed", args) })
		out, _ := recordcmd.Command(f.Root, args)
		next, printed := routetest.Next(out)
		steps := refusalroute.Steps(next)
		want := "bench record --in 'build' assignment 'example' --source 'HEAD' --ticket '1.md' --performer 'fixture-author-1.md' --assignment 'dispatch-1.md'" +
			" --model 'unknown' --effort 'unknown' --ref 'native:dispatch-1.md' --run 'fixture-run' --orchestrator 'fixture-orchestrator' --author-limit '2'"
		if !printed || steps[len(steps)-1] != want || strings.Contains(next, "bench gate") {
			t.Fatalf("assignment = %q, want one next= route that ends with %q and names no gate", out, want)
		}
	})
	t.Run("malformed declaration", func(t *testing.T) {
		f := declaredBoundBuild(t, "Persisted fields: gate-evidence:fixture\n")
		out, code := recordcmd.Command(f.Root, assignArgs(t, *f.Plan.Execution, "1.md"))
		if _, printed := routetest.Next(out); code != 1 || printed || !strings.Contains(out, "persisted-fields-invalid") {
			t.Fatalf("assignment = (%d, %q), want the persisted-fields-invalid error and no route", code, out)
		}
	})
}
