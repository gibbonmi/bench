package commitment_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitcmd"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/intent"
)

func TestCommitmentDeliverableTypes(t *testing.T) {
	for _, tc := range []struct {
		name, path, file string
		folder, allowed  bool
	}{
		{"spec", "specs/light/spec.md", "specs/light/spec.md", false, true},
		{"literal-md-folder", "specs/release.md/spec.md", "specs/release.md/spec.md", false, true},
		{"tickets-only", "specs/light", "specs/light/tickets/01-deliver.md", true, true},
		{"source-file", "src/main.go", "src/main.go", false, false},
		{"unstaged-spec", "specs/light/spec.md", "specs/light/spec.md", false, false},
		{"folder-with-spec", "specs/light", "specs/light/spec.md", true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := deliverableBody
			if tc.name == "unstaged-spec" {
				body = "# Unstaged delivery\n"
			}
			root := boundDeliverable(t, tc.path, map[string]string{tc.file: body}, tc.folder)
			planOut, planCode := commitcmd.Command(root, []string{"plan", "--input", filepath.Join(root, ".bench", "commitment.json")})
			if (planCode == 0) != tc.allowed {
				t.Errorf("deliverable plan=(%s,%d), allowed=%v", planOut, planCode, tc.allowed)
			}
			first := commitmenttest.Assignment(t, root, "first")
			ledger, err := intent.Address(root)
			if err != nil {
				t.Fatal(err)
			}
			before, err := os.ReadFile(ledger)
			if err != nil {
				t.Fatal(err)
			}
			out, code := commitcmd.Command(first, []string{"start", "--outcome", "A", "--request", "first", "--deliverable", tc.path})
			after, err := os.ReadFile(ledger)
			if err != nil {
				t.Fatal(err)
			}
			if tc.allowed {
				if code != 0 {
					t.Fatalf("approved deliverable=(%s,%d)", out, code)
				}
			} else if code != 1 || string(before) != string(after) {
				t.Fatalf("invalid deliverable=(%s,%d), unchanged=%v", out, code, string(before) == string(after))
			}
		})
	}
}

func boundDeliverable(t *testing.T, path string, files map[string]string, folder bool) string {
	t.Helper()
	p := policy([]commitment.Milestone{milestone("M1", "A")}, "M1")
	root := commitmenttest.Repo(t, p)
	for file, body := range files {
		commitmenttest.Write(t, root, file, body)
	}
	commitmenttest.Commit(t, root, "publish deliverable")
	identity := commitment.Identity([]byte(files[path]))
	if folder {
		identity = "git-tree:" + gittest.Output(t, root, "rev-parse", "HEAD:"+path)
	}
	p.Milestones[0].Outcomes[0].Deliverables = []commitment.DeliveryBinding{{Source: commitment.SourceBinding{ID: "A.delivery", Path: path, Identity: identity}}}
	commitmenttest.WritePolicy(t, root, p)
	commitmenttest.Commit(t, root, "approve binding")
	return root
}
