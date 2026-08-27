package worktree

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

const unclaimedAssignmentFingerprintVersion = "bench-unclaimed-assignment-branches/v1"

type unclaimedAssignmentBranch struct{ ref, oid string }
type unclaimedAssignmentSet struct {
	rows        []unclaimedAssignmentBranch
	fingerprint string
}

// planUnclaimedAssignmentSet excludes every assignment record and checked-out ref.
func planUnclaimedAssignmentSet(root string, options CleanupOptions) (unclaimedAssignmentSet, error) {
	assignments, err := intent.Assignments(root)
	if err != nil {
		return unclaimedAssignmentSet{}, err
	}
	protected := make(map[string]bool, len(assignments))
	for _, assignment := range assignments {
		protected[assignment.Branch] = true
	}
	checkouts, err := git.Worktrees(root)
	if err != nil {
		return unclaimedAssignmentSet{}, err
	}
	for _, checkout := range checkouts {
		if checkout.BranchRef != "" {
			protected[checkout.BranchRef] = true
		}
	}
	defaultBranch, ok := git.ResolvedDefault(root)
	if !ok {
		return unclaimedAssignmentSet{}, fmt.Errorf("git repository has no resolvable default branch")
	}
	protected["refs/heads/"+defaultBranch] = true
	branches, err := git.LocalBranches(root)
	if err != nil {
		return unclaimedAssignmentSet{}, fmt.Errorf("git local branches: %w", err)
	}
	set := unclaimedAssignmentSet{}
	for _, branch := range branches {
		ref := "refs/heads/" + branch
		if !strings.HasPrefix(ref, intent.AssignmentBranchPrefix()) || protected[ref] {
			continue
		}
		oid, err := git.Output("-C", root, "rev-parse", "--verify", ref+"^{commit}")
		if err != nil {
			return unclaimedAssignmentSet{}, fmt.Errorf("git assignment branch identity %s: %w", ref, err)
		}
		set.rows = append(set.rows, unclaimedAssignmentBranch{ref, oid})
	}
	sort.Slice(set.rows, func(i, j int) bool { return set.rows[i].ref < set.rows[j].ref })
	if len(set.rows) == 0 {
		return set, nil
	}
	parts := [][]byte{[]byte(unclaimedAssignmentFingerprintVersion), []byte("discard-branch=true"), []byte(fmt.Sprintf("unclaimed=%t", options.Unclaimed))}
	for _, row := range set.rows {
		parts = append(parts, []byte(row.ref), []byte(row.oid))
	}
	set.fingerprint = fingerprintParts(parts...)
	return set, nil
}

func renderUnclaimedAssignmentSet(stdout io.Writer, set unclaimedAssignmentSet) error {
	plans := make([]CleanupPlan, 0, len(set.rows))
	for _, row := range set.rows {
		plans = append(plans, CleanupPlan{Target: row.ref, Action: ActionDiscardRemove, Tracked: "unclaimed", ignoredSummary: "none", Recovery: "none", Fingerprint: set.fingerprint, Reason: "unclaimed assignment branch"})
	}
	return renderCleanups(stdout, plans)
}

func staleUnclaimedPlans(set unclaimedAssignmentSet) []CleanupPlan {
	return []CleanupPlan{{Target: "unknown", Action: ActionError, Tracked: "unclaimed", ignoredSummary: "none", Recovery: "none", Fingerprint: set.fingerprint, Reason: errStaleFingerprint.Error()}}
}

func applyUnclaimedAssignmentSet(root string, set unclaimedAssignmentSet) ([]CleanupPlan, error) {
	current, err := planUnclaimedAssignmentSet(root, CleanupOptions{DiscardBranch: true, Unclaimed: true})
	if err != nil {
		return nil, err
	}
	if current.fingerprint != set.fingerprint || len(current.rows) != len(set.rows) {
		return staleUnclaimedPlans(set), errStaleFingerprint
	}
	plans := make([]CleanupPlan, 0, len(set.rows))
	for i, planned := range set.rows {
		if current.rows[i] != planned {
			return staleUnclaimedPlans(set), errStaleFingerprint
		}
		if err := git.DeleteBranchExact(root, planned.ref, planned.oid); err != nil {
			return append(plans, CleanupPlan{Target: planned.ref, Action: ActionError, Tracked: "unclaimed", ignoredSummary: "none", Recovery: "none", Fingerprint: set.fingerprint, Reason: err.Error()}), err
		}
		plans = append(plans, CleanupPlan{Target: planned.ref, Action: ActionRemoved, Tracked: "unclaimed", ignoredSummary: "none", Recovery: "none", Fingerprint: set.fingerprint})
	}
	return plans, nil
}

// UnclaimedAssignmentBranchRefs gives status the same sorted selection used by clean.
func UnclaimedAssignmentBranchRefs(root string) ([]string, error) {
	set, err := planUnclaimedAssignmentSet(root, CleanupOptions{DiscardBranch: true, Unclaimed: true})
	if err != nil {
		return nil, err
	}
	refs := make([]string, len(set.rows))
	for i, row := range set.rows {
		refs[i] = row.ref
	}
	return refs, nil
}
