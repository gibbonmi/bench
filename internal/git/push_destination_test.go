package git

import (
	"path/filepath"
	"strings"
	"testing"
)

// newTopicRepo initialises a repo whose checked-out branch is topic and whose topic
// branch tracks origin/main. The remote-tracking ref is written by hand, so the test
// needs no second repository and no network.
func newTopicRepo(t *testing.T) string {
	t.Helper()
	root := newRepo(t)
	runGit(t, root, "branch", "-M", "main")
	head := strings.TrimSpace(runGit(t, root, "rev-parse", "HEAD"))
	runGit(t, root, "checkout", "-b", "topic")
	runGit(t, root, "remote", "add", "origin", "https://example.invalid/r.git")
	runGit(t, root, "update-ref", "refs/remotes/origin/main", head)
	runGit(t, root, "config", "branch.topic.remote", "origin")
	runGit(t, root, "config", "branch.topic.merge", "refs/heads/main")
	return root
}

// TestCheckedOutName pins the helper the guard and the destination both compose: a named
// branch answers, and a detached head answers no branch rather than the literal "HEAD"
// that CheckedOutBranch reports for that state.
func TestCheckedOutName(t *testing.T) {
	tests := []struct {
		name     string
		detach   bool
		wantName string
		wantOK   bool
	}{
		{name: "named branch", wantName: "topic", wantOK: true},
		{name: "detached HEAD", detach: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := newTopicRepo(t)
			if tt.detach {
				runGit(t, root, "checkout", "--detach", "HEAD")
			}
			got, ok := CheckedOutName(root)
			if got != tt.wantName || ok != tt.wantOK {
				t.Fatalf("CheckedOutName = (%q, %v), want (%q, %v)", got, ok, tt.wantName, tt.wantOK)
			}
		})
	}
}

// TestBarePushDestinationCheckedOutModes pins row PG37: with push.default unset, and
// under the two modes that push the current branch, the destination is the checked-out
// branch.
func TestBarePushDestinationCheckedOutModes(t *testing.T) {
	for _, mode := range []string{"", "simple", "current"} {
		t.Run("mode="+mode, func(t *testing.T) {
			root := newTopicRepo(t)
			if mode != "" {
				runGit(t, root, "config", "push.default", mode)
			}
			got, ok := BarePushDestination(root)
			if !ok || got != "topic" {
				t.Fatalf("BarePushDestination = (%q, %v), want (\"topic\", true)", got, ok)
			}
		})
	}
}

// TestBarePushDestinationUpstreamModes pins the upstream half of row PG38: under
// upstream and under tracking, the destination is the upstream branch's short name on
// the remote, not the checked-out name.
func TestBarePushDestinationUpstreamModes(t *testing.T) {
	for _, mode := range []string{"upstream", "tracking"} {
		t.Run("mode="+mode, func(t *testing.T) {
			root := newTopicRepo(t)
			runGit(t, root, "config", "push.default", mode)
			got, ok := BarePushDestination(root)
			if !ok || got != "main" {
				t.Fatalf("BarePushDestination = (%q, %v), want (\"main\", true)", got, ok)
			}
		})
	}
}

// TestBarePushDestinationNoDestination pins the refusal half of row PG38 and the
// out-of-repository acceptance row: every state with no single destination reports
// no destination.
func TestBarePushDestinationNoDestination(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T) string
	}{
		{
			name: "matching",
			setup: func(t *testing.T) string {
				root := newTopicRepo(t)
				runGit(t, root, "config", "push.default", "matching")
				return root
			},
		},
		{
			name: "nothing",
			setup: func(t *testing.T) string {
				root := newTopicRepo(t)
				runGit(t, root, "config", "push.default", "nothing")
				return root
			},
		},
		{
			name: "detached HEAD",
			setup: func(t *testing.T) string {
				root := newTopicRepo(t)
				runGit(t, root, "checkout", "--detach", "HEAD")
				return root
			},
		},
		{
			name: "upstream mode with no upstream",
			setup: func(t *testing.T) string {
				root := newRepo(t)
				runGit(t, root, "config", "push.default", "upstream")
				return root
			},
		},
		{
			name:  "outside a repository",
			setup: func(t *testing.T) string { return t.TempDir() },
		},
		{
			name: "missing directory",
			setup: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "absent")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := tt.setup(t)
			got, ok := BarePushDestination(root)
			if ok || got != "" {
				t.Fatalf("BarePushDestination = (%q, %v), want (\"\", false)", got, ok)
			}
		})
	}
}
