package spec

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
)

// landingSpec is the spec body of slug s with one status.
func landingSpec(status string) string { return "# Spec\n\nStatus: " + status + "\n" }

// landingCommit writes the spec of slug s with body in root and commits it with git args.
func landingCommit(t *testing.T, root, body string, args ...string) {
	t.Helper()
	path := filepath.Join(root, "specs", "s", "spec.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	gittest.Output(t, root, "add", "-A")
	gittest.Output(t, root, append([]string{"commit", "-q"}, args...)...)
}

func TestLandingIsTheStatusTransitionOfTheDefaultBranch(t *testing.T) {
	for _, tc := range []struct {
		name string
		// build makes the history on main from a staged spec and returns the landing L.
		build func(t *testing.T, root string) string
	}{
		{name: "later implemented edit", build: func(t *testing.T, root string) string {
			gittest.Output(t, root, "checkout", "-q", "-b", "build")
			gittest.Output(t, root, "commit", "-q", "--allow-empty", "-m", "change")
			gittest.Output(t, root, "checkout", "-q", "main")
			gittest.Output(t, root, "merge", "-q", "--no-ff", "--no-commit", "build")
			landingCommit(t, root, landingSpec(StatusImplemented), "-m", "land build")
			landing := gittest.Output(t, root, "rev-parse", "HEAD")
			landingCommit(t, root, landingSpec(StatusImplemented)+"More.\n", "-m", "edit the implemented spec")
			return landing
		}},
		{name: "flip on the branch", build: func(t *testing.T, root string) string {
			gittest.Output(t, root, "checkout", "-q", "-b", "build")
			landingCommit(t, root, landingSpec(StatusImplemented), "-m", "flip on the branch")
			gittest.Output(t, root, "checkout", "-q", "main")
			gittest.Output(t, root, "merge", "-q", "--no-ff", "-m", "land build", "build")
			return gittest.Output(t, root, "rev-parse", "HEAD")
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := gittest.RepoOnBranch(t, "main")
			landingCommit(t, root, landingSpec("staged"), "-m", "stage")
			want := tc.build(t, root)
			got, found, err := Landing(root, "s")
			if err != nil || !found || got != want {
				t.Fatalf("Landing = %q, %v, %v; want the landing %q", got, found, err, want)
			}
		})
	}
}
