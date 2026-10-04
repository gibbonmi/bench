package status

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/intent"
)

// TestAppendDeliveryRendersTheOutlook holds the delivery row to the shared outlook, so a
// staged spec never becomes the board's delivery action. Before adoption a waiting staged
// spec shows the adoption row, and a repository with no staged spec shows no row. Each
// expected row is written by hand.
func TestAppendDeliveryRendersTheOutlook(t *testing.T) {
	root := initRepo(t)
	path := filepath.Join(root, "specs", "staged", "spec.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("Status: staged\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	blocked := []intent.OutcomeBlocker{{Outcome: "A", Reason: "vendor fix"}}
	adoption := commitment.Outlook{State: "adoption-required", Operation: "plan", Command: "bench commitment plan --input <file>"}
	for _, tc := range []struct {
		name    string
		root    string
		outlook commitment.Outlook
		want    []row
	}{
		{
			name:    "eligible successor",
			outlook: commitment.Outlook{State: "eligible", Milestone: "M1", Next: "B", Blocked: blocked, Operation: "start", Command: "bench commitment start --outcome B --request <request> --deliverable <path>"},
			want:    []row{{4, "commitment", "eligible B in M1; blocked A", commandActionWithArgument(commitmentAction, "start --outcome B --request <request> --deliverable <path>")}},
		},
		{
			name:    "active outcome",
			outlook: commitment.Outlook{State: "active", Milestone: "M1", Active: []string{"A"}},
			want:    []row{{4, "commitment", "active A in M1", advisoryAction("")}},
		},
		{
			name:    "active outcome without a staged spec",
			root:    initRepo(t),
			outlook: commitment.Outlook{State: "active", Milestone: "M1", Active: []string{"A"}},
			want:    []row{{4, "commitment", "active A in M1", advisoryAction("")}},
		},
		{
			name:    "all blocked",
			outlook: commitment.Outlook{State: "all-blocked", Milestone: "M1", Blocked: blocked, Operation: "plan", Command: "bench commitment plan --input <file>"},
			want:    []row{{4, "commitment", "all-blocked in M1; blocked A", commandActionWithArgument(commitmentAction, "plan --input <file>")}},
		},
		{
			name:    "adoption required",
			outlook: adoption,
			want:    []row{{4, "commitment", "adoption-required", commandActionWithArgument(commitmentAction, "plan --input <file>")}},
		},
		{
			name:    "adoption required without a staged spec",
			root:    initRepo(t),
			outlook: adoption,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			caseRoot := root
			if tc.root != "" {
				caseRoot = tc.root
			}
			got := appendDeliveryOutlook(nil, caseRoot, tc.outlook)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("rows = %#v, want %#v", got, tc.want)
			}
			if len(got) != 0 && got[0].signal == "commitment" && got[0].action.render() != tc.outlook.Command {
				t.Fatalf("action = %q, want the outlook command %q", got[0].action.render(), tc.outlook.Command)
			}
		})
	}
}
