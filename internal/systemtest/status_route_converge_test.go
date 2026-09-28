//go:build system

package systemtest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/intent"
)

// TestStatusRouteExecutesUnclaimedCleanup runs the command status routes to. Over a unique
// ref the routed plan keeps the ref; over a landed ref the apply command that plan prints
// removes it.
func TestStatusRouteExecutesUnclaimedCleanup(t *testing.T) {
	for _, tc := range []struct {
		name   string
		landed bool
	}{{name: "unique ref survives"}, {name: "landed ref applies", landed: true}} {
		t.Run(tc.name, func(t *testing.T) {
			repo, err := os.MkdirTemp(owner.root, "status-route [unclaimed]-")
			if err != nil {
				t.Fatal(err)
			}
			git := func(args ...string) {
				t.Helper()
				if result := owner.runAt(repo, nil, "git", args...); result.code != 0 {
					t.Fatalf("git %s = (%d, %q)", strings.Join(args, " "), result.code, result.stderr)
				}
			}
			git("init", "-q", "-b", "main")
			git("config", "user.email", "route@example.test")
			git("config", "user.name", "Route Test")
			git("commit", "--allow-empty", "-qm", "base")
			branch := intent.AssignmentBranchRef(strings.Repeat("a", 32), strings.Repeat("b", 32))
			git("branch", strings.TrimPrefix(branch, "refs/heads/"))
			if !tc.landed {
				git("checkout", "-q", strings.TrimPrefix(branch, "refs/heads/"))
				if err := os.WriteFile(filepath.Join(repo, "orphan.txt"), []byte("orphan\n"), 0o644); err != nil {
					t.Fatal(err)
				}
				git("add", "orphan.txt")
				git("commit", "-qm", "orphan assignment")
				git("checkout", "-q", "main")
			}

			route := owner.runSelected(repo, "status", "--route")
			if route.code != 0 {
				t.Fatalf("bench status --route = (%d, %q, %q)", route.code, route.stdout, route.stderr)
			}
			const cleanupPrefix = "bench worktree clean --discard-branch --unclaimed"
			start := strings.Index(route.stdout, cleanupPrefix)
			if start < 0 {
				t.Fatalf("bench status --route returned %q, want %q", route.stdout, cleanupPrefix)
			}
			command := route.stdout[start:]
			if end := strings.IndexByte(command, '\n'); end >= 0 {
				command = command[:end]
			}
			shim := t.TempDir()
			if err := os.Symlink(owner.selected.path, filepath.Join(shim, "bench")); err != nil {
				t.Fatal(err)
			}
			run := owner.runAt(repo, []string{"PATH=" + shim + string(os.PathListSeparator) + os.Getenv("PATH")}, "bash", "-c", command)
			if run.code != 0 {
				t.Fatalf("routed command %q = (%d, %q, %q)", command, run.code, run.stdout, run.stderr)
			}
			if tc.landed {
				argv, err := axitest.RecoverHelpCommandArgv(run.stdout)
				if err != nil || len(argv) < 2 || argv[0] != "bench" || argv[len(argv)-2] != "--apply" {
					t.Fatalf("routed plan %q printed apply %q (%v), want a --apply <fingerprint> command", run.stdout, argv, err)
				}
				if apply := owner.runSelected(repo, argv[1:]...); apply.code != 0 {
					t.Fatalf("printed apply %q = (%d, %q, %q)", argv, apply.code, apply.stdout, apply.stderr)
				}
			}
			present := owner.runAt(repo, nil, "git", "show-ref", "--verify", "--quiet", branch).code == 0
			if present == tc.landed {
				t.Fatalf("after routed command %q, ref %q present=%t, want %t", command, branch, present, !tc.landed)
			}
		})
	}
}
