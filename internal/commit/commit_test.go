package commit

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/specbuild"
)

type commitTestGate struct{}

func (commitTestGate) Bootstrap(context.Context, string, string, string, string) error { return nil }
func (commitTestGate) AdvanceMarker(context.Context, string, string, string, string) error {
	return nil
}

func TestCommandRefusesSpecCommitWhileBuildActive(t *testing.T) {
	root := t.TempDir()
	gitCommit := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	gitCommit("init", "-q")
	if err := os.MkdirAll(filepath.Join(root, "specs", "active"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "specs", "active", "spec.md"), []byte("# Active\n\nStatus: staged\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "tracked.txt"), []byte("base\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitCommit("add", ".")
	gitCommit("-c", "user.name=Test", "-c", "user.email=test@example.com", "commit", "-qm", "base")
	if _, err := specbuild.New(root, commitTestGate{}, nil).Start(context.Background(), "active"); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "tracked.txt"), []byte("changed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitCommit("add", "tracked.txt")
	before := strings.TrimSpace(string(runGit(t, root, "rev-parse", "HEAD")))
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldwd)
	var stdout, stderr bytes.Buffer
	code := Command([]string{"-m", "commit", "--spec", "active", "tracked.txt"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("exit = %d, want 1; stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "bench spec build promote active") {
		t.Fatalf("stderr = %q, want promote remedy", stderr.String())
	}
	if got := string(mustReadFile(t, filepath.Join(root, "specs", "active", "spec.md"))); !strings.Contains(got, "Status: staged") {
		t.Fatalf("spec changed: %q", got)
	}
	if got := strings.TrimSpace(string(runGit(t, root, "rev-parse", "HEAD"))); got != before {
		t.Fatalf("HEAD = %s, want %s", got, before)
	}
}

func runGit(t *testing.T, root string, args ...string) []byte {
	t.Helper()
	out, err := exec.Command("git", append([]string{"-C", root}, args...)...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return out
}
func mustReadFile(t *testing.T, path string) []byte {
	t.Helper()
	out, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(fmt.Errorf("read %s: %w", path, err))
	}
	return out
}

// The integrated command behavior (block-check, gate, flip, stage, commit, exit codes)
// is gate-observed through the CLI in internal/contract/runtime; this unit test pins only
// the pure argument parser, whose branch table the black-box tests exercise only by
// outcome. The rendered reason text belongs to the shared grammar and its toon renderer,
// so the misuse cases assert the whole line rather than a local phrasing.
func TestParseArgs(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		wantMsg  string
		wantSpec string
		wantPath string // first positional path; "" means none expected
		wantN    int    // number of positional paths
		wantErr  string // substring of usageErr; "" means no error
	}{
		{name: "message and one path", args: []string{"-m", "msg", "a.txt"}, wantMsg: "msg", wantPath: "a.txt", wantN: 1},
		{name: "spec flag consumes its value", args: []string{"-m", "m", "--spec", "feature", "a.txt"}, wantMsg: "m", wantSpec: "feature", wantPath: "a.txt", wantN: 1},
		{name: "paths on both sides of flags", args: []string{"a.txt", "-m", "m", "b.txt"}, wantMsg: "m", wantPath: "a.txt", wantN: 2},
		{name: "-- makes a leading-dash token a path", args: []string{"-m", "m", "--", "-weird.txt"}, wantMsg: "m", wantPath: "-weird.txt", wantN: 1},
		{name: "no message", args: []string{"a.txt"}, wantErr: "-m <msg> is required"},
		{name: "no paths", args: []string{"-m", "m"}, wantErr: "at least one <path> is required"},
		{name: "unknown flag", args: []string{"-m", "m", "--nope", "a.txt"}, wantErr: "usage: bench commit (unknown argument: --nope)"},
		{name: "dangling -m", args: []string{"-m"}, wantErr: "usage: bench commit (missing argument: -m)"},
		{name: "dangling --spec", args: []string{"-m", "m", "--spec"}, wantErr: "usage: bench commit (missing argument: --spec)"},
		// The two shapes an unset or blank shell variable produces. An empty path is
		// caught by the shared grammar (it would otherwise resolve to the cwd and
		// widen the commit); an empty or blank message is caught here, so it is
		// reported like the missing -m rather than by a raw `git commit -m ""`.
		{name: "empty path", args: []string{"-m", "m", ""}, wantErr: `usage: bench commit (unknown argument: "")`},
		{name: "empty message", args: []string{"-m", "", "a.txt"}, wantErr: "-m <msg> must not be empty"},
		{name: "blank message", args: []string{"-m", "   \t\n", "a.txt"}, wantErr: "-m <msg> must not be empty"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msg, specSlug, paths, help, usageErr := parseArgs(tc.args)
			if help != "" {
				t.Fatalf("unexpected help %q", help)
			}
			if tc.wantErr != "" {
				if usageErr == "" || !strings.Contains(usageErr, tc.wantErr) {
					t.Fatalf("usageErr = %q, want substring %q", usageErr, tc.wantErr)
				}
				return
			}
			if usageErr != "" {
				t.Fatalf("unexpected usageErr %q", usageErr)
			}
			if msg != tc.wantMsg || specSlug != tc.wantSpec || len(paths) != tc.wantN {
				t.Errorf("got msg=%q spec=%q paths=%v, want msg=%q spec=%q n=%d", msg, specSlug, paths, tc.wantMsg, tc.wantSpec, tc.wantN)
			}
			if len(paths) > 0 && paths[0] != tc.wantPath {
				t.Errorf("first path = %q, want %q", paths[0], tc.wantPath)
			}
		})
	}
}

// Help is a success the caller prints, never a misuse: all three spellings return the
// declared help text with no usage error.
func TestParseArgsHelpIsSuccess(t *testing.T) {
	for _, spelling := range []string{"help", "--help", "-h"} {
		t.Run(spelling, func(t *testing.T) {
			_, _, _, help, usageErr := parseArgs([]string{spelling})
			if usageErr != "" {
				t.Fatalf("usageErr = %q, want none", usageErr)
			}
			if !strings.HasPrefix(help, "usage: bench commit") {
				t.Fatalf("help = %q, want the declared help text", help)
			}
		})
	}
}
