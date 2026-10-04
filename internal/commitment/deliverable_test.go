package commitment_test

import (
	"os"
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
		{"tickets-only", "specs/light", "specs/light/tickets/01-deliver.md", true, true},
		{"source-file", "src/main.go", "src/main.go", false, false},
		{"unstaged-spec", "specs/light/spec.md", "specs/light/spec.md", false, false},
		{"folder-with-spec", "specs/light", "specs/light/spec.md", true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			p := policy([]commitment.Milestone{milestone("M1", "A")}, "M1")
			root := commitmenttest.Repo(t, p)
			body := deliverableBody
			if tc.name == "unstaged-spec" {
				body = "# Unstaged delivery\n"
			}
			commitmenttest.Write(t, root, tc.file, body)
			commitmenttest.Commit(t, root, "publish deliverable")
			identity := commitment.Identity([]byte(body))
			if tc.folder {
				identity = "git-tree:" + gittest.Output(t, root, "rev-parse", "HEAD:"+tc.path)
			}
			p.Milestones[0].Outcomes[0].Deliverables = []commitment.DeliveryBinding{{Source: commitment.SourceBinding{ID: "A.delivery", Path: tc.path, Identity: identity}}}
			commitmenttest.WritePolicy(t, root, p)
			commitmenttest.Commit(t, root, "approve binding")
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
