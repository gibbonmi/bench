package status

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/roadmap"
)

func TestTimeoutGateIsDistinctHighestSeveritySignal(t *testing.T) {
	rows := appendGateInfo(nil, GateInfo{Present: true, State: "ready", Status: "timeout"}, t.TempDir())
	if len(rows) != 1 || rows[0].detail != "timeout" || rows[0].sev != 0 || !strings.Contains(rows[0].action, "hang") {
		t.Fatalf("timeout rows = %#v", rows)
	}
}

// retirementCount reads specs/*/spec.md through the shared predicate. The seam is the
// directory: write spec files, then assert the count. The predicate cases (fence, CRLF,
// trailing whitespace, wrong status) are the load-bearing behaviors.
func TestRetirementCount(t *testing.T) {
	cases := []struct {
		name  string
		files map[string]string
		want  int
	}{
		{
			name:  "unfenced implemented marker is counted",
			files: map[string]string{"a/spec.md": "# spec\nStatus: implemented\n"},
			want:  1,
		},
		{
			name:  "staged status is not counted",
			files: map[string]string{"a/spec.md": "# spec\nStatus: staged\n"},
			want:  0,
		},
		{
			name:  "no status line is not counted",
			files: map[string]string{"a/spec.md": "# spec\njust prose\n"},
			want:  0,
		},
		{
			name:  "implemented inside a code fence is not counted",
			files: map[string]string{"a/spec.md": "# spec\n```\nStatus: implemented\n```\n"},
			want:  0,
		},
		{
			name:  "CRLF line endings still match",
			files: map[string]string{"a/spec.md": "# spec\r\nStatus: implemented\r\n"},
			want:  1,
		},
		{
			name:  "trailing tabs and spaces after the value still match",
			files: map[string]string{"a/spec.md": "Status: implemented\t \n"},
			want:  1,
		},
		{
			name: "multiple specs counted correctly",
			files: map[string]string{
				"a/spec.md": "Status: implemented\n",
				"b/spec.md": "Status: staged\n",
				"c/spec.md": "Status:\timplemented\n",
			},
			want: 2,
		},
		{
			name:  "hidden spec file is ignored",
			files: map[string]string{".hidden/spec.md": "Status: implemented\n"},
			want:  0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			for name, body := range tc.files {
				p := filepath.Join(root, "specs", name)
				if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			if got := retirementCount(root); got != tc.want {
				t.Errorf("retirementCount = %d, want %d", got, tc.want)
			}
		})
	}
}

// Absent specs/ directory → 0, never a panic on the missing path.
func TestRetirementCountNoSpecsDir(t *testing.T) {
	if got := retirementCount(t.TempDir()); got != 0 {
		t.Errorf("retirementCount with no specs/ = %d, want 0", got)
	}
}

func TestOrphanedPickupCount(t *testing.T) {
	root := t.TempDir()
	for name, body := range map[string]string{
		"reviews/paired.md":    "pickup\n",
		"reviews/orphaned.md":  "pickup\n",
		"specs/paired/spec.md": "Status: staged\n",
	} {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if got := orphanedPickupCount(root); got != 1 {
		t.Fatalf("orphanedPickupCount = %d, want 1", got)
	}
}

func TestRoadmapReconcileCounts(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "specs", "merged"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "specs", "merged", "spec.md"), []byte("Status: implemented\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	content := "specs/merged/spec.md specs/retired/spec.md specs/<slug>/spec.md\n```\nspecs/fenced/spec.md\n```\n"
	if err := os.WriteFile(filepath.Join(root, roadmap.RoadmapFile), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	merged, dangling, state := roadmapReconcileCounts(root)
	if merged != 1 || dangling != 1 || state.Failed() {
		t.Fatalf("roadmapReconcileCounts = (%d, %d, %s), want (1, 1, parsed)", merged, dangling, state)
	}
}

// short guards the [:7] tree-prefix slice against a short or "none" hash.
func TestShort(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"", ""},
		{"none", "none"},
		{"abcdef", "abcdef"},
		{"0123456789", "0123456"},
	} {
		if got := Short(tc.in); got != tc.want {
			t.Errorf("Short(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestStaleGateDetailActionCurrentTreeNoneFailsClosed(t *testing.T) {
	detail, action := staleGateDetailAction(t.TempDir(), "0123456789abcdef", "none")
	if detail != "stale (gated tree 0123456, work tree none)" {
		t.Fatalf("detail = %q, want strong stale detail", detail)
	}
	if action != "re-run the gate" {
		t.Fatalf("action = %q, want re-run the gate", action)
	}
}

func initRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	gitRun(t, root, "init")
	gitRun(t, root, "config", "user.email", "t@example.com")
	gitRun(t, root, "config", "user.name", "t")
	return root
}

func gitRun(t *testing.T, root string, args ...string) string {
	t.Helper()
	out, err := git.Output(append([]string{"-C", root}, args...)...)
	if err != nil {
		t.Fatalf("git %v: %v (%s)", args, err, out)
	}
	return out
}

// A clean, committed tree with no signals renders the clean message and nothing else.
func TestRenderClean(t *testing.T) {
	root := initRepo(t)
	if err := os.WriteFile(filepath.Join(root, "f.txt"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")

	out := render(root, false)
	if out != "bench: clean — nothing pending\n" {
		t.Errorf("clean render = %q", out)
	}
}

// A dirty tree leads with the git action; the capture-drain row is present but outranked.
func TestRenderDirtyLeadsGitOverDrainRow(t *testing.T) {
	root := initRepo(t)
	if err := os.WriteFile(filepath.Join(root, "f.txt"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")
	// Make the tree dirty (uncommitted change) so the git signal fires.
	if err := os.WriteFile(filepath.Join(root, "f.txt"), []byte("y\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Park an idea so the drain signal fires.
	if err := os.WriteFile(filepath.Join(root, "IDEAS.md"), []byte("- 2026-07-03  an idea\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	out := render(root, false)
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if !strings.HasPrefix(lines[0], "▶ commit on green  (git)") {
		t.Errorf("lead line = %q, want git action lead", lines[0])
	}
	if !strings.Contains(out, "1 idea(s), 0 open learning(s)") || !strings.Contains(out, "/bench-what-next") {
		t.Errorf("drain row missing from:\n%s", out)
	}
}

// A working roadmap alone is not pending capture: no drain row, the board stays clean.
func TestRenderWorkingRoadmapAloneIsClean(t *testing.T) {
	root := initRepo(t)
	content := "# Roadmap\n\n## Recommended sequence\n\n1. Shape next item - /bench-shape-idea\n"
	if err := os.WriteFile(filepath.Join(root, "ROADMAP.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")

	if out := render(root, false); out != "bench: clean — nothing pending\n" {
		t.Errorf("roadmap-only render = %q, want clean board", out)
	}
}

func TestAppendWorktreeIgnoresUnownedBranchPrefix(t *testing.T) {
	root := initRepo(t)
	if err := os.WriteFile(filepath.Join(root, "f.txt"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")
	gitRun(t, root, "branch", "worktree-agent-orphan")

	if rows := appendWorktree(nil, root); len(rows) != 0 {
		t.Fatalf("branch prefix created worktree status without ownership evidence: %#v", rows)
	}
}

// appendWorktree must surface a `git worktree list` failure as a visible row, never as
// silence that a reader mistakes for "no worktree signals" — the false-empty class FT29
// swept. The failure is induced deterministically (the FT29 gitOpError style: break the
// git query itself, here by revoking read access to .git) rather than a PATH-shimmed git.
func TestAppendWorktreeSurfacesClassifyFailure(t *testing.T) {
	root := initRepo(t)
	if err := os.WriteFile(filepath.Join(root, "f.txt"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")

	gitDir := filepath.Join(root, ".git")
	if err := os.Chmod(gitDir, 0o000); err != nil {
		t.Fatalf("chmod .git unreadable: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(gitDir, 0o755) })

	rows := appendWorktree(nil, root)
	if len(rows) == 0 {
		t.Fatal("appendWorktree dropped the classify failure instead of surfacing a row")
	}
	if !strings.Contains(rows[0].detail, "worktree list failed") {
		t.Errorf("row detail = %q, want it to name the git worktree-list failure", rows[0].detail)
	}
}

// Command rejects an unknown argument with a usage line and exit 2, prints usage on -h,
// and accepts --all as the one added token — while --all plus junk and near-misses stay
// usage errors so a typo never silently prints the default board.
func TestCommandArgs(t *testing.T) {
	if r, c := Command([]string{"--bogus"}); c != 2 || !strings.Contains(r, "usage:") {
		t.Errorf("unknown arg: report %q exit %d", r, c)
	}
	if r, c := Command([]string{"-h"}); c != 0 || !strings.Contains(r, "usage: bench status") {
		t.Errorf("help: report %q exit %d", r, c)
	}
	if r, c := Command([]string{"-h"}); !strings.Contains(r, "[--all]") {
		t.Errorf("help usage should advertise [--all], got %q exit %d", r, c)
	}
	if r, c := Command([]string{"--all"}); c != 0 {
		t.Errorf("--all should be accepted with exit 0, got report %q exit %d", r, c)
	}
	for _, bad := range [][]string{{"--all", "extra"}, {"--allx"}, {"-a"}} {
		if r, c := Command(bad); c != 2 || !strings.Contains(r, "usage:") {
			t.Errorf("args %q: report %q exit %d, want usage exit 2", bad, r, c)
		}
	}
}
