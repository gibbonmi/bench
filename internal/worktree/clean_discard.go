// The explicit discard of one unrecorded branch: the fallback that resolves a target the
// ledger does not hold, the row that target plans, and the transaction that writes the
// dated discarded ref before it deletes the branch.

package worktree

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

// The two windows of the discard transaction. The first precedes the discarded ref write and
// the second follows it, so a test can stand in each one. The tokens sit here rather than in
// ownership.go, which is over its size budget.
const (
	StepDiscardedRefWrite     LifecycleStep = "discarded-ref-write"
	StepDiscardedBranchDelete LifecycleStep = "discarded-branch-delete"
)

// branchSegment is the last path segment of a branch ref. For an assignment branch it is the
// assignment id.
func branchSegment(ref string) string { return ref[strings.LastIndex(ref, "/")+1:] }

// discardSelector is the --target operand that selects one unclaimed branch: the id segment
// of an assignment branch, or the branch path of a shift branch, which has no id segment.
func discardSelector(ref string) string {
	if strings.HasPrefix(ref, intent.AssignmentBranchPrefix()) && intent.ValidIdentity(branchSegment(ref)) {
		return branchSegment(ref)
	}
	return strings.TrimPrefix(ref, "refs/heads/")
}

// discardTargetCommand is the explicit discard command for one unclaimed branch. A unique row
// of the unclaimed plan ends with it, and so does an explicit row planned without the flag.
func discardTargetCommand(ref string) string { return DiscardTargetCommand(discardSelector(ref)) }

// DiscardTargetCommand is the one spelling of the explicit discard command for a --target
// selector: an assignment id, or the branch path of a shift branch. The retire listing names
// each superseded assignment with it.
func DiscardTargetCommand(selector string) string {
	return cleanCommand(CleanupOptions{DiscardBranch: true}, "--target", selector)
}

// Unique is how many rows of the counted plan carry the unique class: the rows that only an
// explicit discard removes.
func (counts UnclaimedRefCounts) Unique() int {
	for _, class := range counts.Classes {
		if class.Class == string(classUnique) {
			return class.Count
		}
	}
	return 0
}

// classDetail is refVerdict.detail with suffix in the place of the retained text a unique row
// ends with there. Every other row already ends with the suffix it is given.
func classDetail(v refVerdict, suffix string) string {
	if v.class != classUnique {
		return v.detail(suffix)
	}
	return strings.TrimSuffix(v.detail(suffix), uniqueRetainedDetail) + suffix
}

// discardPrefixRef reports whether target names a branch in a Bench-created namespace, with
// or without refs/heads/, and returns the full ref. The ledger declares both namespaces.
func discardPrefixRef(target string) (string, bool) {
	for _, prefix := range []string{intent.AssignmentBranchPrefix(), intent.ShiftBranchPrefix()} {
		if strings.HasPrefix(target, prefix) {
			return target, true
		}
		if strings.HasPrefix(target, strings.TrimPrefix(prefix, "refs/heads/")) {
			return "refs/heads/" + target, true
		}
	}
	return "", false
}

// explicitTarget is what one operand resolves to: a recorded assignment, or one branch of the
// unclaimed set when branch is set.
type explicitTarget struct {
	assignment intent.Assignment
	branch     *unclaimedAssignmentBranch
}

// discardCandidates is the unclaimed set the fallback selects from. It is the set the bulk
// planner selects, planned once per selection and only when an operand needs it.
type discardCandidates struct {
	root    string
	planned bool
	set     unclaimedAssignmentSet
	err     error
}

func (c *discardCandidates) rows() ([]unclaimedAssignmentBranch, error) {
	if !c.planned {
		c.set, c.err = planUnclaimedAssignmentSet(c.root, unclaimedOptions())
		c.planned = true
	}
	return c.set.rows, c.err
}

// resolve routes a namespace-prefix operand to the fallback before the record path, whose
// relative-path check would refuse it. Any other operand takes the record path, and only an
// id that resolves no record falls back to the unclaimed branches that end with it.
func (c *discardCandidates) resolve(target string) (explicitTarget, error) {
	if ref, ok := discardPrefixRef(target); ok {
		return c.byRef(ref)
	}
	assignment, err := resolveAssignmentIn(c.root, target, resumeActiveState)
	if !errors.Is(err, errTargetUnassigned) || !intent.ValidIdentity(target) {
		return explicitTarget{assignment: assignment}, err
	}
	rows, err := c.rows()
	if err != nil {
		return explicitTarget{}, err
	}
	var matched []unclaimedAssignmentBranch
	for _, row := range rows {
		if strings.HasPrefix(row.ref, intent.AssignmentBranchPrefix()) && branchSegment(row.ref) == target {
			matched = append(matched, row)
		}
	}
	if len(matched) == 0 {
		return explicitTarget{}, errTargetUnassigned
	}
	if len(matched) > 1 {
		refs := make([]string, len(matched))
		for i, row := range matched {
			refs[i] = row.ref
		}
		return explicitTarget{}, fmt.Errorf("target matches more than one unclaimed branch: %s", strings.Join(refs, " and "))
	}
	return explicitTarget{branch: &matched[0]}, nil
}

// byRef resolves one full branch ref. A recorded branch plans through its record. A branch
// outside the unclaimed set refuses with the fact that keeps it out of that set.
func (c *discardCandidates) byRef(ref string) (explicitTarget, error) {
	assignments, err := intent.Assignments(c.root)
	if err != nil {
		return explicitTarget{}, err
	}
	for _, assignment := range assignments {
		if assignment.Branch == ref {
			selected, err := resolveAssignmentIn(c.root, assignment.ID, resumeActiveState)
			return explicitTarget{assignment: selected}, err
		}
	}
	rows, err := c.rows()
	if err != nil {
		return explicitTarget{}, err
	}
	for _, row := range rows {
		if row.ref == ref {
			return explicitTarget{branch: &row}, nil
		}
	}
	checkouts, err := git.Worktrees(c.root)
	if err != nil {
		return explicitTarget{}, err
	}
	for _, checkout := range checkouts {
		if checkout.BranchRef == ref {
			return explicitTarget{}, fmt.Errorf("branch %s is checked out at %s", ref, checkout.Path)
		}
	}
	if defaultBranch, ok := git.ResolvedDefault(c.root); ok && "refs/heads/"+defaultBranch == ref {
		return explicitTarget{}, fmt.Errorf("branch %s is the default branch", ref)
	}
	return explicitTarget{}, errTargetUnassigned
}

// unrecordedCleanupRow is one selected unclaimed branch: its class verdict, the discarded ref
// its apply writes first, and the row the plan prints. A row with no discarded ref writes none.
type unrecordedCleanupRow struct {
	unclaimedAssignmentBranch
	discarded string
	plan      CleanupPlan
}

// planUnrecordedRow plans one unclaimed branch. Without --discard-branch every class retains
// and names the flag. With it every class removes, and a unique branch plans its discarded
// ref on the day now names.
func planUnrecordedRow(branch unclaimedAssignmentBranch, now time.Time, options CleanupOptions) (unrecordedCleanupRow, error) {
	if branch.fault != "" {
		return unrecordedCleanupRow{}, errors.New(branch.fault)
	}
	row := unrecordedCleanupRow{unclaimedAssignmentBranch: branch, plan: branch.plan("")}
	row.plan.Action, row.plan.Reason = ActionRetain, classDetail(branch.refVerdict, discardTargetCommand(branch.ref))
	if !options.DiscardBranch {
		return row, nil
	}
	row.plan.Action, row.plan.Reason = ActionDiscardRemove, classDetail(branch.refVerdict, branch.reason)
	if branch.class == classUnique {
		row.discarded = intent.DiscardedRef(now, branch.ref)
		row.plan.Recovery = row.discarded
	}
	return row, nil
}

// requalifyUnrecordedRow re-plans one unclaimed member against the repository as it stands
// now. A member that left the unclaimed set, or whose verdict or discarded ref changed, is
// stale.
func requalifyUnrecordedRow(j joins, root string, planned unrecordedCleanupRow, options CleanupOptions) (unrecordedCleanupRow, error) {
	set, err := planUnclaimedAssignmentSet(root, unclaimedOptions())
	if err != nil {
		return planned, err
	}
	for _, branch := range set.rows {
		if branch.ref != planned.ref {
			continue
		}
		current, err := planUnrecordedRow(branch, j.now(), options)
		if err != nil {
			return planned, err
		}
		current.plan.Fingerprint = planned.plan.Fingerprint
		if current.unclaimedAssignmentBranch != planned.unclaimedAssignmentBranch || current.discarded != planned.discarded {
			return current, errStaleFingerprint
		}
		return current, nil
	}
	return planned, errStaleFingerprint
}

// applyUnrecordedRows runs each unrecorded member after every recorded one. Each removing
// member re-plans just before its own transaction, because an earlier member's removal can
// take away the holder a subsumed member relied on.
func applyUnrecordedRows(j joins, root string, set explicitCleanupSet, options CleanupOptions, plans []CleanupPlan) ([]CleanupPlan, error) {
	members := set.plans()
	for k, planned := range set.unrecorded {
		unreached := members[len(set.rows)+k+1:]
		if !planned.plan.Action.Removes() {
			plans = append(plans, planned.plan)
			continue
		}
		current, err := requalifyUnrecordedRow(j, root, planned, options)
		if err != nil {
			return notAttemptedPlans(append(plans, requalifiedOutcome(planned.plan, current.plan, err)), unreached, notAttemptedDetail), err
		}
		if err := discardUnrecordedBranch(j, root, current); err != nil {
			return notAttemptedPlans(append(plans, faultedPlan(current.plan, CleanupPlan{}, err)), unreached, notAttemptedDetail), err
		}
		removed := current.plan
		removed.Action, removed.Reason = ActionRemoved, ""
		plans = append(plans, removed)
	}
	return plans, nil
}

// discardUnrecordedBranch writes the discarded ref, then deletes the branch at the exact tip
// the plan named. A branch that moved after the write survives the delete, and the discarded
// ref keeps the tip the plan named.
func discardUnrecordedBranch(j joins, root string, row unrecordedCleanupRow) error {
	if err := writeDiscardedRef(j, root, row); err != nil {
		return err
	}
	return git.DeleteBranchExact(root, row.ref, row.oid)
}

// writeDiscardedRef writes the row's discarded ref at the row's tip. A ref already at that tip
// is the write a stopped apply made, so a retry skips it. A ref at another tip is a handle
// this write must not replace, so the apply refuses and keeps the branch. The zero old value
// makes a ref that appears after the read fail the write.
func writeDiscardedRef(j joins, root string, row unrecordedCleanupRow) error {
	if row.discarded == "" {
		return nil
	}
	if err := hit(j.cleanupBoundary, StepDiscardedRefWrite); err != nil {
		return err
	}
	existing, _ := git.Output("-C", root, "show-ref", "--verify", "--hash", row.discarded)
	switch existing {
	case row.oid:
	case "":
		if _, err := git.Output("-C", root, "update-ref", row.discarded, row.oid, strings.Repeat("0", len(row.oid))); err != nil {
			return fmt.Errorf("write discarded ref %s: %w", row.discarded, err)
		}
	default:
		return fmt.Errorf("discarded ref %s already names %s; the branch stays", row.discarded, existing)
	}
	return hit(j.cleanupBoundary, StepDiscardedBranchDelete)
}
