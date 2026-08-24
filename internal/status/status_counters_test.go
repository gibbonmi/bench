// Tests for the retirement, orphaned-pickup, tickets-only residue, and roadmap reconcile counts.
package status

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/roadmap"
	"github.com/gibbonmi/bench/internal/roadmap/roadmaptest"
)

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

// ticketsOnlyRepo commits files into a repository parked on its default branch, the
// branch the retirement row requires. This lets a residue fixture be ranked against
// the housekeeping rows it joins.
func ticketsOnlyRepo(t *testing.T, files map[string]string) string {
	t.Helper()
	root := initRepo(t)
	files["tracked.txt"] = "base\n"
	for name, body := range files {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-m", "fixture")
	gitRun(t, root, "branch", "-M", "main")
	return root
}

// H05: the residue row carries the count and the command that closes one. It ranks
// below the two housekeeping rows it joins. A count of residue never displaces a
// retirement or an orphaned pickup inside the five-row budget.
func TestTicketsOnlyResidueRowCountsAndRanksBelowItsBand(t *testing.T) {
	root := ticketsOnlyRepo(t, map[string]string{
		"specs/landed-ticket/tickets/one.md": "ticket\n",
		"specs/second-ticket/tickets/two.md": "ticket\n",
		"specs/merged/spec.md":               "Status: implemented\n",
		"reviews/orphan.md":                  "pickup\n",
	})

	signals := Signals(root)
	position := func(name, detail string) int {
		t.Helper()
		for i, s := range signals {
			if s.Name == name && strings.Contains(s.Detail, detail) {
				return i
			}
		}
		t.Fatalf("no %s row containing %q in %#v", name, detail, signals)
		return -1
	}
	residue := position("specs", "tickets-only")
	if got, want := signals[residue], testSignal(11, "specs", "2 tickets-only spec folders", "bench commit --spec <slug>"); got != want {
		t.Fatalf("residue row = %#v, want %#v", got, want)
	}
	retirement, orphaned := position("specs", "awaiting retirement"), position("reviews", "orphaned review")
	if residue < retirement || residue < orphaned {
		t.Fatalf("residue row at %d outranks retirement %d / orphaned pickup %d", residue, retirement, orphaned)
	}
}

// H06 — a specs tree holding no tickets-only folder renders no residue row at all. A
// spec-backed folder is not residue, so the fixture also pins that the row reads the
// tickets-only predicate rather than counting every child of specs/.
func TestNoTicketsOnlyFolderRendersNoResidueRow(t *testing.T) {
	root := ticketsOnlyRepo(t, map[string]string{"specs/spec-backed/spec.md": "Status: staged\n"})

	for _, s := range Signals(root) {
		if strings.Contains(s.Detail, "tickets-only") {
			t.Fatalf("clean specs tree produced residue row %#v", s)
		}
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

// TestRoadmapReconcileCountsFromRowFile covers story 20 (PR17): the split board keeps
// a row's spec path in roadmap/FT7.md's body, not ROADMAP.md's index line. So the
// reconcile scan must classify a spec named only there. A scan of ROADMAP.md alone
// counts zero.
func TestRoadmapReconcileCountsFromRowFile(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "specs", "merged"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "specs", "merged", "spec.md"), []byte("Status: implemented\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	const heading = "**FT7 (LOW) — x.**"
	body := heading + "\nMerged into specs/merged/spec.md.\n"
	roadmaptest.WriteSplitBoard(t, root, heading+"\n", map[string]string{"FT7.md": body})
	merged, dangling, state := roadmapReconcileCounts(root)
	if merged != 1 || dangling != 0 || state.Failed() {
		t.Fatalf("roadmapReconcileCounts = (%d, %d, %s), want (1, 0, parsed)", merged, dangling, state)
	}
}
