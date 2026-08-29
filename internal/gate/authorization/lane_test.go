// The fast-lane authority: the worktree commit's oracle, beside the gate authority.
package authorization

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
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

// A declared previous tip makes the authority derive the prose placeholder itself: the
// Markdown the graded tree changes against that tip's tree, and nothing else. A caller
// that stated the whole tree would grade unchanged prose.
func TestLaneAuthorityDerivesTheProseSubjectFromThePreviousTip(t *testing.T) {
	root := gittest.RepoOnBranch(t, "main")
	commitFiles(t, root, "base", map[string]string{"kept.md": "kept\n", "changed.md": "before\n"})
	previous, err := benchgit.Output("-C", root, "rev-parse", "HEAD^{commit}")
	if err != nil {
		t.Fatal(err)
	}
	commitFiles(t, root, "graded", map[string]string{"changed.md": "after\n", "added.md": "added\n", "notes.txt": "prose is not this\n"})
	tree, err := benchgit.Output("-C", root, "rev-parse", "HEAD^{tree}")
	if err != nil {
		t.Fatal(err)
	}
	argv := filepath.Join(t.TempDir(), "prose-argv")
	var stdout, stderr bytes.Buffer
	authority := LaneAuthority{
		PreviousTip: previous,
		// A stated list must lose to the derived one, so the row proves which input wins.
		NamedMarkdown: []string{"kept.md"},
		Checks: []gate.Phase{{Name: "prose", Argv: []string{"sh", "-c",
			`for path in "$@"; do printf '%s\n' "$path"; done > ` + argv, "prose", gate.LaneNamedMarkdownToken}}},
	}

	if got := authority.Authorize(context.Background(), root, tree, &stdout, &stderr); got.Kind != LanePass {
		t.Fatalf("kind = %q, want %q; stdout=%q stderr=%q", got.Kind, LanePass, stdout.String(), stderr.String())
	}
	recorded, err := os.ReadFile(argv)
	if err != nil {
		t.Fatalf("the prose check recorded no argv: %v", err)
	}
	if got := string(recorded); got != "added.md\nchanged.md\n" {
		t.Fatalf("prose argv = %q, want the Markdown the graded tree changes against the previous tip", got)
	}
}

// commitFiles writes and commits one set of files, so a row states the two trees its
// derivation is measured between.
func commitFiles(t *testing.T, root, message string, files map[string]string) {
	t.Helper()
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(root, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if out, err := benchgit.Raw("-C", root, "add", "-A"); err != nil {
		t.Fatalf("stage %s: %v\n%s", message, err, out)
	}
	if out, err := benchgit.Raw("-C", root, "commit", "-q", "-m", message); err != nil {
		t.Fatalf("commit %s: %v\n%s", message, err, out)
	}
}
