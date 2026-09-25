//go:build system

package systemtest

import (
	"os"
	"path/filepath"
	"testing"
)

// BO30: a nested Bench verb bounds its own response in its own process. The child prints
// a 30-line blob, so exec receives the child's 10-line projection and passes it through
// unchanged. The spill line names the child's spill in the assignment's scope, and exec,
// which runs from the primary checkout, writes no spill of its own.
func TestExecNestedBenchBoundsOnce(t *testing.T) {
	fixture, source := execBoundFixture(t)
	systemCommit(t, source.path, "thirty.txt", joinLines(seqLines(1, 30)), "thirty lines")
	if err := owner.observeSelected(); err != nil {
		t.Fatal(err)
	}
	// The dispatch observation line reaches stderr outside every bound, so the child's own
	// line would join its projection. This run turns the observation off.
	environment := append(fixture.environment(fixture.home), "BENCH_COMMAND_OBSERVE=0")
	result := owner.runWithInput(fixture.root, environment, "", owner.selected.path,
		execBoundArgs(source, owner.selected.path, "worktree", "show", source.assignment, "HEAD:thirty.txt")...)
	if result.code != 0 {
		t.Fatalf("exec = (%d, %q, %q), want exit 0", result.code, result.stdout, result.stderr)
	}
	path := spilledPath(t, result.stdout, seqLines(1, 4), seqLines(26, 30))
	if got, want := readSpillFile(t, path), joinLines(seqLines(1, 30)); got != want {
		t.Fatalf("spill file = %q, want the child's complete output %q", got, want)
	}
	scope := filepath.Dir(path)
	if filepath.Base(scope) != source.assignment {
		t.Fatalf("spill path = %q, want the child's spill in the scope of assignment %q", path, source.assignment)
	}
	// The repository's store holds the child's scope alone, so exec spilled nothing.
	scopes, err := os.ReadDir(filepath.Dir(scope))
	if err != nil || len(scopes) != 1 {
		t.Fatalf("store scopes = %v (%v), want only the child's scope", scopes, err)
	}
}
