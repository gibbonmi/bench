// The fast-lane authority: the worktree commit's oracle, beside the gate authority.
package authorization

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gittest"
)

// laneFixture is a repository with one commit, which is what a lane run needs to
// materialize the graded tree beside HEAD.
func laneFixture(t *testing.T) (root, tree string) {
	t.Helper()
	root = gittest.RepoOnBranch(t, "main")
	if out, err := benchgit.Raw("-C", root, "commit", "-q", "--allow-empty", "-m", "base"); err != nil {
		t.Fatalf("commit the fixture base: %v\n%s", err, out)
	}
	tree, err := benchgit.Output("-C", root, "rev-parse", "HEAD^{tree}")
	if err != nil {
		t.Fatal(err)
	}
	return root, tree
}

// A lane whose checks all pass attributes a lane pass, states the outcome with its check
// names, and carries no evidence: a lane pass authorizes one commit and nothing else.
func TestLaneAuthorityAttributesAPass(t *testing.T) {
	root, tree := laneFixture(t)
	var stdout, stderr bytes.Buffer
	authority := LaneAuthority{Checks: []gate.Phase{
		{Name: "first", Argv: []string{"true"}},
		{Name: "second", Argv: []string{"true"}},
	}}

	got := authority.Authorize(context.Background(), root, tree, &stdout, &stderr)
	if got.Kind != LanePass {
		t.Fatalf("kind = %q, want %q; stderr=%q", got.Kind, LanePass, stderr.String())
	}
	if got.Evidence != "" {
		t.Errorf("evidence = %q, want none: a lane pass proves nothing a landing may reuse", got.Evidence)
	}
	if !strings.Contains(stdout.String(), "lane{outcome=pass,checks=first,second}") {
		t.Errorf("stdout = %q, want the pass record naming both checks in order", stdout.String())
	}
	if strings.Contains(stdout.String(), "green") {
		t.Errorf("stdout = %q, want no green token", stdout.String())
	}
}

// A failing check attributes a lane fail and states the check with its first output line,
// so the caller refuses without re-reading the stream.
func TestLaneAuthorityAttributesAFailWithTheFirstDiagnostic(t *testing.T) {
	root, tree := laneFixture(t)
	var stdout, stderr bytes.Buffer
	authority := LaneAuthority{Checks: []gate.Phase{
		{Name: "unit", Argv: []string{"sh", "-c", "echo the first line; echo a later line; exit 2"}},
	}}

	got := authority.Authorize(context.Background(), root, tree, &stdout, &stderr)
	if got.Kind != LaneFail {
		t.Fatalf("kind = %q, want %q", got.Kind, LaneFail)
	}
	if !strings.Contains(stdout.String(), "lane{outcome=fail,check=unit}") {
		t.Errorf("stdout = %q, want the fail record naming the check", stdout.String())
	}
	if !strings.Contains(stdout.String(), "the first line") {
		t.Errorf("stdout = %q, want the check's first output line", stdout.String())
	}
}

// A lane that declares no check cannot decide anything, so the run errors and the
// authority attributes a fail rather than an unearned pass.
func TestLaneAuthorityRefusesAnUnrunnableLane(t *testing.T) {
	root, tree := laneFixture(t)
	var stdout, stderr bytes.Buffer

	got := LaneAuthority{}.Authorize(context.Background(), root, tree, &stdout, &stderr)
	if got.Kind != LaneFail {
		t.Fatalf("kind = %q, want %q", got.Kind, LaneFail)
	}
	if !strings.Contains(stderr.String(), "declares no checks") {
		t.Errorf("stderr = %q, want the run's own diagnostic", stderr.String())
	}
	if strings.Contains(stdout.String(), "outcome=") {
		t.Errorf("stdout = %q, want no outcome record for a lane that never ran", stdout.String())
	}
}
