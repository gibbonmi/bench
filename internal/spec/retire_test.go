package spec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/usage"
)

// retireRepo builds the minimal repository bench spec retire reads and deletes from: a
// folder spec committed at HEAD with body, plus any extraFiles (repo-relative path to
// content) committed alongside it — used to place a roadmap/FT<n>.md detail file. It
// returns a linked worktree of that repository, because retire refuses the primary
// checkout; retirePrimary returns the primary checkout instead.
func retireRepo(t *testing.T, slug, body string, extraFiles map[string]string) (root string) {
	t.Helper()
	return retireWorktree(t, retirePrimary(t, slug, body, extraFiles))
}

// retirePrimary builds the repository itself and returns its primary checkout.
func retirePrimary(t *testing.T, slug, body string, extraFiles map[string]string) (root string) {
	t.Helper()
	root = gittest.RepoOnBranch(t, "main")
	writeFolderSpec(t, root, slug, body)
	for path, content := range extraFiles {
		full := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	gittest.MustRun(t, root, "add", "-A")
	gittest.Commit(t, root, "base")
	return root
}

// retireWorktree adds a linked worktree of primary and returns its resolved toplevel —
// the checkout shape a Bench phase runs a retire from.
func retireWorktree(t *testing.T, primary string) string {
	t.Helper()
	root := filepath.Join(t.TempDir(), "linked")
	gittest.MustRun(t, primary, "worktree", "add", "-q", "-b", "topic", root)
	return gittest.Output(t, root, "rev-parse", "--show-toplevel")
}

// runRetire runs retireCommand from root, the way the CLI does: RepoBase resolves
// through the process cwd.
func runRetire(t *testing.T, root, arg string) (string, int) {
	t.Helper()
	return runSpec(t, root, "retire", arg)
}

// runSpec runs one spec subcommand from dir, so the process cwd anchors the repository.
func runSpec(t *testing.T, dir string, args ...string) (string, int) {
	t.Helper()
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldwd)
	out, code, _ := Command(args)
	return out, code
}

// TestRetireNextLineNamesTheBoardRemainder covers FC1-FC4 and FC6, and the last-line
// and NBSP edges: the next: line's exact text as a function of the spec's Roadmap:
// value and whether roadmap/FT<n>.md exists on disk.
func TestRetireNextLineNamesTheBoardRemainder(t *testing.T) {
	generic := "next: promote durable content, remove the ROADMAP row FT<n> and its roadmap/FT<n>.md detail file, commit as `spec-retire: s`\n"
	cases := []struct {
		name       string
		body       string
		extraFiles map[string]string
		want       string
	}{
		{
			name:       "FC1: FT7 with existing detail file names both",
			body:       "Status: implemented\nRoadmap: FT7\n",
			extraFiles: map[string]string{"roadmap/FT7.md": "row\n"},
			want:       "next: promote durable content, remove the ROADMAP row FT7 and its roadmap/FT7.md detail file, commit as `spec-retire: s`\n",
		},
		{
			name: "FC2: no Roadmap line prints the generic line",
			body: "Status: implemented\n",
			want: generic,
		},
		{
			name: "FC3: FT7 with no detail file names FT7 and no detail path",
			body: "Status: implemented\nRoadmap: FT7\n",
			want: "next: promote durable content, remove the ROADMAP row FT7, commit as `spec-retire: s`\n",
		},
		{
			name: "FC4: lowercase ft7 prints the generic line",
			body: "Status: implemented\nRoadmap: ft7\n",
			want: generic,
		},
		{
			name: "FC4: internal space FT 7 prints the generic line",
			body: "Status: implemented\nRoadmap: FT 7\n",
			want: generic,
		},
		{
			name: "FC4: empty Roadmap value prints the generic line",
			body: "Status: implemented\nRoadmap:\n",
			want: generic,
		},
		{
			name:       "FC6: fenced Roadmap line is not read",
			body:       "Status: implemented\n\n```\nRoadmap: FT7\n```\n",
			extraFiles: map[string]string{"roadmap/FT7.md": "row\n"},
			want:       generic,
		},
		{
			name:       "edge: Roadmap line is the last line with no trailing newline",
			body:       "Status: implemented\nRoadmap: FT7",
			extraFiles: map[string]string{"roadmap/FT7.md": "row\n"},
			want:       "next: promote durable content, remove the ROADMAP row FT7 and its roadmap/FT7.md detail file, commit as `spec-retire: s`\n",
		},
		{
			name:       "edge: NBSP around the id trims like a space",
			body:       "Status: implemented\nRoadmap:  FT7 \n",
			extraFiles: map[string]string{"roadmap/FT7.md": "row\n"},
			want:       "next: promote durable content, remove the ROADMAP row FT7 and its roadmap/FT7.md detail file, commit as `spec-retire: s`\n",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := retireRepo(t, "s", tc.body, tc.extraFiles)
			out, code := runRetire(t, root, "s")
			if code != 0 {
				t.Fatalf("code = %d, want 0; out = %q", code, out)
			}
			lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
			got := lines[len(lines)-1] + "\n"
			if got != tc.want {
				t.Errorf("next line = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestRetireStagedRefusesAndDeletesNothing pins FC5: a staged spec still refuses at
// exit 1, and the spec file survives untouched.
func TestRetireStagedRefusesAndDeletesNothing(t *testing.T) {
	root := retireRepo(t, "s", "Status: staged\nRoadmap: FT7\n", nil)
	specPath := filepath.Join(root, "specs", "s", "spec.md")
	before, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatal(err)
	}
	out, code := runRetire(t, root, "s")
	if code != 1 {
		t.Fatalf("code = %d, want 1; out = %q", code, out)
	}
	after, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatalf("spec file was deleted: %v", err)
	}
	if string(after) != string(before) {
		t.Errorf("spec file content changed on refusal")
	}
}

// TestRetireDeletesTheFolderAndExitsZero pins FA4: the surviving verb still removes a
// merged-implemented spec's whole folder at exit 0.
func TestRetireDeletesTheFolderAndExitsZero(t *testing.T) {
	root := retireRepo(t, "s", "Status: implemented\nRoadmap: FT7\n", nil)
	out, code := runRetire(t, root, "s")
	if code != 0 {
		t.Fatalf("code = %d, want 0; out = %q", code, out)
	}
	if _, err := os.Stat(filepath.Join(root, "specs", "s")); !os.IsNotExist(err) {
		t.Errorf("specs/s still present: err = %v", err)
	}
}

// TestCommitmentRetireClosedRow pins retirement after a verified closure. That delivery
// closed its row in its own publication, so the next: line schedules no board cleanup,
// and the verb still deletes the folder. A HEAD policy that records no delivery of this
// spec keeps the board-remainder line.
func TestCommitmentRetireClosedRow(t *testing.T) {
	closed := "next: promote durable content, commit as `spec-retire: s`\n"
	open := "next: promote durable content, remove the ROADMAP row FT7, commit as `spec-retire: s`\n"
	for _, tc := range []struct {
		name, delivered, want string
		uncommitted           bool
	}{
		{name: "delivered", delivered: "specs/s/spec.md", want: closed},
		{name: "undelivered", want: open},
		{name: "other-spec", delivered: "specs/other/spec.md", want: open},
		{name: "working-tree-only", delivered: "specs/s/spec.md", uncommitted: true, want: open},
	} {
		t.Run(tc.name, func(t *testing.T) {
			policy := map[string]string{commitment.PolicyPath: retirePolicy(t, tc.delivered)}
			committed := policy
			if tc.uncommitted {
				committed = nil
			}
			root := retireRepo(t, "s", "Status: implemented\nRoadmap: FT7\n", committed)
			if tc.uncommitted {
				if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(commitment.PolicyPath)), []byte(policy[commitment.PolicyPath]), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			out, code := runRetire(t, root, "s")
			lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
			if code != 0 || lines[len(lines)-1]+"\n" != tc.want {
				t.Fatalf("retire = (%d, %q), want exit 0 and the last line %q", code, out, tc.want)
			}
			if _, err := os.Stat(filepath.Join(root, "specs", "s")); !os.IsNotExist(err) {
				t.Errorf("specs/s still present: err = %v", err)
			}
		})
	}
}

// retirePolicy is a valid commitment policy whose outcome A owns FT7 and approves
// specs/s/spec.md and specs/other/spec.md. It records the delivery of delivered, if any.
func retirePolicy(t *testing.T, delivered string) string {
	t.Helper()
	binding := func(id, path string) commitment.DeliveryBinding {
		return commitment.DeliveryBinding{Source: commitment.SourceBinding{ID: id, Path: path, Identity: "sha256:" + id}, Obligations: []string{"FT7"}}
	}
	policy := commitment.Policy{Version: 1, ActiveMilestone: "M", Milestones: []commitment.Milestone{{ID: "M", Outcomes: []commitment.Outcome{{
		ID: "A", Criteria: []commitment.Criterion{{ID: "done", Text: "The row is delivered."}},
		Sources:      []commitment.SourceBinding{{ID: "FT7", Path: "roadmap/FT7.md", Identity: "sha256:FT7"}},
		Deliverables: []commitment.DeliveryBinding{binding("s", "specs/s/spec.md"), binding("other", "specs/other/spec.md")},
	}}}}}
	var err error
	if delivered != "" {
		if policy, _, err = commitment.Deliver(policy, delivered, "source", "evidence"); err != nil {
			t.Fatal(err)
		}
	}
	data, err := commitment.Bytes(policy)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// TestRetireOnPrimaryCheckoutRefusesAndDeletesNothing pins the primary-checkout refusal:
// retire prints the shared refusal with its own retire follow-on step, so the operator
// reads one route, and it leaves the merged-implemented spec on disk.
func TestRetireOnPrimaryCheckoutRefusesAndDeletesNothing(t *testing.T) {
	root := retirePrimary(t, "s", "Status: implemented\nRoadmap: FT7\n", nil)
	out, code := runRetire(t, root, "s")
	if code != 1 {
		t.Fatalf("code = %d, want 1; out = %q", code, out)
	}
	if shared := usage.PrimaryCheckoutRefusal(); !strings.HasPrefix(out, shared) {
		t.Errorf("out = %q, want the shared refusal as its prefix (%q)", out, shared)
	}
	if out == usage.PrimaryCheckoutRefusal()+"\n" {
		t.Errorf("out = %q names only the worktree creation, want a retire follow-on step too", out)
	}
	if !strings.Contains(out, "bench spec retire") {
		t.Errorf("out = %q, want it to name bench spec retire as the follow-on step", out)
	}
	if _, err := os.Stat(filepath.Join(root, "specs", "s", "spec.md")); err != nil {
		t.Errorf("spec was deleted on refusal: %v", err)
	}
}
