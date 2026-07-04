package status

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
)

// retirementCount reads specs/*.md through the awk-ported predicate. The seam is the
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
			files: map[string]string{"a.md": "# spec\nStatus: implemented\n"},
			want:  1,
		},
		{
			name:  "staged status is not counted",
			files: map[string]string{"a.md": "# spec\nStatus: staged\n"},
			want:  0,
		},
		{
			name:  "no status line is not counted",
			files: map[string]string{"a.md": "# spec\njust prose\n"},
			want:  0,
		},
		{
			name:  "implemented inside a code fence is not counted",
			files: map[string]string{"a.md": "# spec\n```\nStatus: implemented\n```\n"},
			want:  0,
		},
		{
			name:  "CRLF line endings still match",
			files: map[string]string{"a.md": "# spec\r\nStatus: implemented\r\n"},
			want:  1,
		},
		{
			name:  "trailing tabs and spaces after the value still match",
			files: map[string]string{"a.md": "Status: implemented\t \n"},
			want:  1,
		},
		{
			name: "multiple specs counted correctly",
			files: map[string]string{
				"a.md": "Status: implemented\n",
				"b.md": "Status: staged\n",
				"c.md": "Status:\timplemented\n",
			},
			want: 2,
		},
		{
			name:  "hidden spec file is ignored",
			files: map[string]string{".hidden.md": "Status: implemented\n"},
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

// short guards the [:7] tree-prefix slice against a short or "none" hash.
func TestShort(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"", ""},
		{"none", "none"},
		{"abcdef", "abcdef"},
		{"0123456789", "0123456"},
	} {
		if got := short(tc.in); got != tc.want {
			t.Errorf("short(%q) = %q, want %q", tc.in, got, tc.want)
		}
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

	out := render(root)
	if out != "bench: clean — nothing pending\n" {
		t.Errorf("clean render = %q", out)
	}
}

// A dirty tree leads with the git action; the parked-idea footer is present but never leads.
func TestRenderDirtyLeadsGitAndFooterTrails(t *testing.T) {
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
	// Park an idea so the footer is non-empty.
	if err := os.WriteFile(filepath.Join(root, "ROADMAP.md"), []byte("- 2026-07-03  an idea\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	out := render(root)
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) < 2 {
		t.Fatalf("expected at least a lead + footer, got:\n%s", out)
	}
	if !strings.HasPrefix(lines[0], "▶ commit on green / push  (git)") {
		t.Errorf("lead line = %q, want git action lead", lines[0])
	}
	last := lines[len(lines)-1]
	if !strings.HasSuffix(last, "idea(s) parked — bench roadmap") {
		t.Errorf("footer line = %q, want parked-idea footer last", last)
	}
	if strings.HasPrefix(last, "▶") {
		t.Errorf("footer must never lead: %q", last)
	}
}

func TestRenderSurfacesOrphanedWorktreeBranch(t *testing.T) {
	root := initRepo(t)
	if err := os.WriteFile(filepath.Join(root, "f.txt"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "base")
	gitRun(t, root, "branch", "worktree-agent-orphan")

	out := render(root)
	if !strings.Contains(out, "orphaned worktree branch") {
		t.Fatalf("status did not surface orphaned worktree branch:\n%s", out)
	}
	if !strings.Contains(out, "delete scratch branch") {
		t.Fatalf("status did not recommend scratch branch cleanup:\n%s", out)
	}
}

// Command rejects an unknown argument with a usage line and exit 2, and prints usage on -h.
func TestCommandArgs(t *testing.T) {
	if r, c := Command([]string{"--bogus"}); c != 2 || !strings.Contains(r, "usage:") {
		t.Errorf("unknown arg: report %q exit %d", r, c)
	}
	if r, c := Command([]string{"-h"}); c != 0 || !strings.Contains(r, "usage: bench status") {
		t.Errorf("help: report %q exit %d", r, c)
	}
}
