package worktree

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

const unclaimedAssignmentFingerprintVersion = "bench-unclaimed-assignment-branches/v1"

type unclaimedAssignmentBranch struct{ ref, oid, reason string }
type unclaimedAssignmentSet struct {
	rows        []unclaimedAssignmentBranch
	fingerprint string
}

// unclaimedBranchReason reports whether ref sits in a Bench-created branch namespace and
// names the reason its row carries. Bench creates both namespaces, so Bench retires both.
func unclaimedBranchReason(ref string) (string, bool) {
	switch {
	case strings.HasPrefix(ref, intent.AssignmentBranchPrefix()):
		return "unclaimed assignment branch", true
	case strings.HasPrefix(ref, intent.ShiftBranchPrefix()):
		return "shift residue branch", true
	default:
		return "", false
	}
}

// planUnclaimedAssignmentSet selects the unclaimed assignment and shift branches. It
// excludes every assignment record, checked-out ref, and the default branch.
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
		reason, owned := unclaimedBranchReason(ref)
		if !owned || protected[ref] {
			continue
		}
		oid, err := git.Output("-C", root, "rev-parse", "--verify", ref+"^{commit}")
		if err != nil {
			return unclaimedAssignmentSet{}, fmt.Errorf("git assignment branch identity %s: %w", ref, err)
		}
		set.rows = append(set.rows, unclaimedAssignmentBranch{ref, oid, reason})
	}
	sort.Slice(set.rows, func(i, j int) bool { return set.rows[i].ref < set.rows[j].ref })
	if len(set.rows) == 0 {
		return set, nil
	}
	parts := [][]byte{[]byte(unclaimedAssignmentFingerprintVersion), []byte("discard-branch=true"), []byte(fmt.Sprintf("unclaimed=%t", options.Unclaimed))}
	for _, row := range set.rows {
		parts = append(parts, []byte(row.ref), []byte(row.oid), []byte(row.reason))
	}
	set.fingerprint = fingerprintParts(parts...)
	return set, nil
}

func renderUnclaimedAssignmentSet(stdout io.Writer, set unclaimedAssignmentSet) error {
	plans := make([]CleanupPlan, 0, len(set.rows))
	for _, row := range set.rows {
		plans = append(plans, CleanupPlan{Target: row.ref, Action: ActionDiscardRemove, Tracked: "unclaimed", ignoredSummary: "none", Recovery: "none", Fingerprint: set.fingerprint, Reason: row.reason})
	}
	return renderCleanups(stdout, plans)
}

func staleUnclaimedPlans(set unclaimedAssignmentSet) []CleanupPlan {
	return []CleanupPlan{{Target: "unknown", Action: ActionError, Tracked: "unclaimed", ignoredSummary: "none", Recovery: "none", Fingerprint: set.fingerprint, Reason: errStaleFingerprint.Error()}}
}

// unclaimedOptions is what this mode's options are when no caller parsed any. The status
// reader asks for the selection outside the clean grammar, so it has no invocation to read
// them from. A caller that parsed an invocation passes its own options instead.
func unclaimedOptions() CleanupOptions {
	return CleanupOptions{DiscardBranch: true, Unclaimed: true}
}

// unclaimedReplan is this mode's own re-plan command, beside the landed selector's. It takes
// the caller's options so a later grammar change cannot leave the rendered command behind.
// Today the two are pinned equal: valid() requires --discard-branch and refuses the other
// two modifiers under --unclaimed, so every valid invocation carries unclaimedOptions().
func unclaimedReplan(options CleanupOptions) []axi.InvocationArgument {
	return cleanArguments(options, "--unclaimed")
}

// applyUnclaimedAssignmentSet deletes each planned branch at the exact object the plan
// named. It reports the outcome rows alone. A stale refusal carries no rows, because the
// refusal row is this command surface's own spelling and the caller renders it; a row
// returned here would be a second derivation the caller discards.
func applyUnclaimedAssignmentSet(j joins, root string, set unclaimedAssignmentSet, options CleanupOptions) ([]CleanupPlan, error) {
	// The window that refusal exists for: a branch can enter or leave the namespace between
	// the caller's plan read and this re-plan. The boundary is nil in production; it lets a
	// test stand in that window, as the three checkout modes already do.
	if err := hit(j.cleanupBoundary, StepApplyLocked); err != nil {
		return nil, err
	}
	current, err := planUnclaimedAssignmentSet(root, options)
	if err != nil {
		return nil, err
	}
	if current.fingerprint != set.fingerprint || len(current.rows) != len(set.rows) {
		return nil, errStaleFingerprint
	}
	plans := make([]CleanupPlan, 0, len(set.rows))
	for i, planned := range set.rows {
		if current.rows[i] != planned {
			return nil, errStaleFingerprint
		}
		if err := git.DeleteBranchExact(root, planned.ref, planned.oid); err != nil {
			return append(plans, CleanupPlan{Target: planned.ref, Action: ActionError, Tracked: "unclaimed", ignoredSummary: "none", Recovery: "none", Fingerprint: set.fingerprint, Reason: err.Error()}), err
		}
		plans = append(plans, CleanupPlan{Target: planned.ref, Action: ActionRemoved, Tracked: "unclaimed", ignoredSummary: "none", Recovery: "none", Fingerprint: set.fingerprint})
	}
	return plans, nil
}

// UnclaimedAssignmentBranchRefs gives status the same sorted assignment-and-shift
// selection used by clean.
func UnclaimedAssignmentBranchRefs(root string) ([]string, error) {
	set, err := planUnclaimedAssignmentSet(root, unclaimedOptions())
	if err != nil {
		return nil, err
	}
	refs := make([]string, len(set.rows))
	for i, row := range set.rows {
		refs[i] = row.ref
	}
	return refs, nil
}
