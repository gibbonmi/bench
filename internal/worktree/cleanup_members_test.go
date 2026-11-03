// The shared cleanup traversal on its own: its visit order, its pass-through of members a
// mode does not visit, and its stop at the first failed visit. Each visit is injected and
// records a call trace, so the expectation states every call the traversal makes.
package worktree

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// tracedTraversal is a traversal over named members. A name that starts with "keep" is not
// visited, a name in fail fails its visit, and every predicate, visit, and projection call
// appends to calls.
func tracedTraversal(calls *[]string, fail map[string]bool) cleanupTraversal[string] {
	return cleanupTraversal[string]{
		visits: func(member string) bool { return !strings.HasPrefix(member, "keep") },
		retained: func(member string) CleanupPlan {
			*calls = append(*calls, "retain "+member)
			return CleanupPlan{Target: member, Action: ActionRetain}
		},
		visit: func(member string) (CleanupPlan, error) {
			*calls = append(*calls, "visit "+member)
			if fail[member] {
				return CleanupPlan{Target: member, Action: ActionError}, errors.New("fault on " + member)
			}
			return CleanupPlan{Target: member, Action: ActionRemoved}, nil
		},
		stopped: func(plans []CleanupPlan, failed int) []CleanupPlan {
			*calls = append(*calls, fmt.Sprintf("stopped at %d", failed))
			return append(plans, CleanupPlan{Target: "suffix", Action: ActionNotAttempted})
		},
	}
}

// traversalRows is the target and action of each row, in order.
func traversalRows(plans []CleanupPlan) []string {
	rows := make([]string, 0, len(plans))
	for _, plan := range plans {
		rows = append(rows, plan.Target+" "+string(plan.Action))
	}
	return rows
}

// TestCleanupTraversalOrder is R02. Every member is visited once, in the order given, after
// the rows the caller already holds.
func TestCleanupTraversalOrder(t *testing.T) {
	t.Parallel()
	var calls []string
	plans, err := tracedTraversal(&calls, nil).run([]string{"c", "a", "b"}, []CleanupPlan{{Target: "held", Action: ActionRemoved}})
	mustNoError(t, err)
	if want := []string{"visit c", "visit a", "visit b"}; !reflect.DeepEqual(calls, want) {
		t.Fatalf("calls = %q, want %q", calls, want)
	}
	if got, want := traversalRows(plans), []string{"held removed", "c removed", "a removed", "b removed"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("rows = %q, want %q", got, want)
	}
}

// TestCleanupTraversalRetainedVisits is R05. A member the predicate does not select reaches
// no visit and passes through as its retained row, in its own position.
func TestCleanupTraversalRetainedVisits(t *testing.T) {
	t.Parallel()
	var calls []string
	plans, err := tracedTraversal(&calls, nil).run([]string{"keep-x", "a", "keep-y"}, nil)
	mustNoError(t, err)
	if want := []string{"retain keep-x", "visit a", "retain keep-y"}; !reflect.DeepEqual(calls, want) {
		t.Fatalf("calls = %q, want %q", calls, want)
	}
	if got, want := traversalRows(plans), []string{"keep-x retain", "a removed", "keep-y retain"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("rows = %q, want %q", got, want)
	}
}

// TestCleanupTraversalStopsAtFailedVisit is R14 and the traversal half of R11. The first
// failed visit is the last call before the projection: no later member is visited or passed
// through, the failed member's own outcome is the last row ahead of the projection, and its
// error is the result.
func TestCleanupTraversalStopsAtFailedVisit(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		members []string
		calls   []string
		rows    []string
	}{
		{"first", []string{"a", "b", "keep-x"}, []string{"visit a", "stopped at 0"}, []string{"a error", "suffix not-attempted"}},
		{"later", []string{"keep-x", "b", "a", "c"}, []string{"retain keep-x", "visit b", "visit a", "stopped at 2"}, []string{"keep-x retain", "b removed", "a error", "suffix not-attempted"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var calls []string
			plans, err := tracedTraversal(&calls, map[string]bool{"a": true}).run(tc.members, nil)
			if err == nil || err.Error() != "fault on a" {
				t.Fatalf("error = %v, want the failed visit's own error", err)
			}
			if !reflect.DeepEqual(calls, tc.calls) {
				t.Fatalf("calls = %q, want %q", calls, tc.calls)
			}
			if got := traversalRows(plans); !reflect.DeepEqual(got, tc.rows) {
				t.Fatalf("rows = %q, want %q", got, tc.rows)
			}
		})
	}
}
