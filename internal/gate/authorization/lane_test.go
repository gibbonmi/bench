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
	return root, laneRev(t, root, "HEAD^{tree}")
}

func laneRev(t *testing.T, root, revision string) string {
	t.Helper()
	resolved, err := benchgit.Output("-C", root, "rev-parse", revision)
	if err != nil {
		t.Fatal(err)
	}
	return resolved
}

// A lane whose checks all pass attributes a lane pass, states the outcome with its check
// names, and carries no evidence: a lane pass authorizes one commit and nothing else.
func TestLaneAuthorityAttributesAPass(t *testing.T) {
	root, tree := laneFixture(t)
	base := laneRev(t, root, "HEAD^{commit}")
	var stdout, stderr bytes.Buffer
	authority := LaneAuthority{Base: base, Checks: []gate.Phase{
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
	base := laneRev(t, root, "HEAD^{commit}")
	var stdout, stderr bytes.Buffer
	authority := LaneAuthority{Base: base, Checks: []gate.Phase{
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
	base := laneRev(t, root, "HEAD^{commit}")
	var stdout, stderr bytes.Buffer

	got := LaneAuthority{Base: base}.Authorize(context.Background(), root, tree, &stdout, &stderr)
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

// PL4: the declared base makes the authority derive the prose placeholder itself: the
// changed Markdown the composed tree holds, and nothing else. An unchanged file is not
// prose this commit brought, and a deleted file is prose the tree no longer has.
func TestLaneAuthorityDerivesTheProseSubjectFromTheBase(t *testing.T) {
	root := gittest.RepoOnBranch(t, "main")
	commitFiles(t, root, "base", map[string]string{"kept.md": "kept\n", "changed.md": "before\n", "gone.md": "gone\n"})
	base := laneRev(t, root, "HEAD^{commit}")
	if err := os.Remove(filepath.Join(root, "gone.md")); err != nil {
		t.Fatal(err)
	}
	commitFiles(t, root, "graded", map[string]string{"changed.md": "after\n", "added.md": "added\n", "notes.txt": "prose is not this\n"})
	tree := laneRev(t, root, "HEAD^{tree}")
	argv := filepath.Join(t.TempDir(), "prose-argv")
	var stdout, stderr bytes.Buffer
	authority := LaneAuthority{
		Base: base,
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
		t.Fatalf("prose argv = %q, want the changed Markdown the composed tree holds", got)
	}
}

// commitFiles writes and commits one set of files, so a row states the two trees its
// derivation is measured between.
func commitFiles(t *testing.T, root, message string, files map[string]string) {
	t.Helper()
	for name, body := range files {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
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

// PL5, WM36: the prose subject is framed with `-z` and split on NUL, so a path with a
// space and a byte above ASCII reaches the prose check as one argument of its own bytes.
// A newline split under the default `core.quotepath` hands the check a C-quoted name that
// matches no file, and a field split halves the name at the space.
func TestLaneAuthorityCarriesANonASCIIProsePathVerbatim(t *testing.T) {
	root := gittest.RepoOnBranch(t, "main")
	commitFiles(t, root, "base", map[string]string{"kept.md": "kept\n"})
	base := laneRev(t, root, "HEAD^{commit}")
	commitFiles(t, root, "graded", map[string]string{"café notes.md": "prose\n"})
	tree := laneRev(t, root, "HEAD^{tree}")
	argv := filepath.Join(t.TempDir(), "prose-argv")
	var stdout, stderr bytes.Buffer
	authority := LaneAuthority{
		Base: base,
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
	if got := string(recorded); got != "café notes.md\n" {
		t.Fatalf("prose argv = %q, want the added path's own bytes as one argument", got)
	}
}

// TestLaneAuthorityNamesTheSelectedChecksAndClasses is PL20. A selective lane's pass
// line names the checks that ran and the classes that selected them, so a reader learns
// why each check ran and why the others did not.
func TestLaneAuthorityNamesTheSelectedChecksAndClasses(t *testing.T) {
	root := gittest.RepoOnBranch(t, "main")
	// The two directories the embed derivation walks. The kit's own checkout carries
	// them, and a tree that omits one reports an absent directory rather than an empty
	// embed list.
	commitFiles(t, root, "base", map[string]string{
		"kept.md": "kept\n", "cmd/bench/main.go": "package main\n\nfunc main() {}\n",
		"internal/x/x.go": "package x\n",
	})
	base := laneRev(t, root, "HEAD^{commit}")
	commitFiles(t, root, "graded", map[string]string{"note.md": "note\n"})
	tree := laneRev(t, root, "HEAD^{tree}")
	var stdout, stderr bytes.Buffer
	authority := LaneAuthority{Base: base, Selective: true, Checks: []gate.Phase{
		{Name: "gofmt", Argv: []string{"true"}},
		{Name: "prose", Argv: []string{"true"}},
		{Name: "vet", Argv: []string{"true"}},
	}}

	got := authority.Authorize(context.Background(), root, tree, &stdout, &stderr)
	if got.Kind != LanePass {
		t.Fatalf("kind = %q, want %q; stderr=%q", got.Kind, LanePass, stderr.String())
	}
	if !strings.Contains(stdout.String(), "lane{outcome=pass,checks=prose,classes=markdown}") {
		t.Errorf("stdout = %q, want the selected checks and their classes", stdout.String())
	}
}
