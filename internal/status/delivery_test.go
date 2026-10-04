package status

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/intent"
)

// TestAppendDeliveryRendersTheOutlook holds the delivery row to the shared outlook. A
// published commitment replaces the staged-spec row, so a staged spec outside the
// commitment never becomes the board's delivery action. Before adoption the staged-spec
// row stays. Each expected row is written by hand.
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
	for _, tc := range []struct {
		name    string
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
			name:    "all blocked",
			outlook: commitment.Outlook{State: "all-blocked", Milestone: "M1", Blocked: blocked, Operation: "plan", Command: "bench commitment plan --input <file>"},
			want:    []row{{4, "commitment", "all-blocked in M1; blocked A", commandActionWithArgument(commitmentAction, "plan --input <file>")}},
		},
		{
			name:    "adoption required",
			outlook: commitment.Outlook{State: "adoption-required", Operation: "plan", Command: "bench commitment plan --input <file>"},
			want:    []row{{4, "specs", "1 staged spec(s)", commandActionWithArgument(implementSpecPhaseAction, "specs/staged/spec.md")}},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := appendDeliveryOutlook(nil, root, tc.outlook)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("rows = %#v, want %#v", got, tc.want)
			}
			if got[0].signal == "commitment" && got[0].action.render() != tc.outlook.Command {
				t.Fatalf("action = %q, want the outlook command %q", got[0].action.render(), tc.outlook.Command)
			}
		})
	}
}
