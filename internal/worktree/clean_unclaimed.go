package worktree

import (
	"fmt"
	"io"
	"slices"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

const unclaimedAssignmentFingerprintVersion = "bench-unclaimed-assignment-branches/v2"

// StepUnlockedReplan is the re-plan this selector runs before it deletes anything. Every
// StepApplyLocked site sits inside the registration lock a checkout holds, and this mode
// locks no checkout: it compares refs, so it detects a concurrent writer instead of
// excluding one. The token stays here, beside its only site, because ownership.go is over
// its size budget.
const StepUnlockedReplan LifecycleStep = "unlocked-replan"

type unclaimedAssignmentBranch struct {
	ref, reason string
	refVerdict
}
type unclaimedAssignmentSet struct {
	rows        []unclaimedAssignmentBranch
	fingerprint string
	options     CleanupOptions
}

// faulted reports whether any row is an error row. Such a set has no fingerprint, so a
// fingerprint apply refuses as stale, and its plan render refuses, so the plan exits 1 and
// --apply-current stops before its apply.
func (set unclaimedAssignmentSet) faulted() bool {
	for _, row := range set.rows {
		if row.fault != "" {
			return true
		}
	}
	return false
}

// plan is the cleanup row one classified branch renders.
func (row unclaimedAssignmentBranch) plan(fingerprint string) CleanupPlan {
	return CleanupPlan{Target: row.ref, Action: row.action(), Tracked: "unclaimed", ignoredSummary: "none", Recovery: "none", Fingerprint: fingerprint, Reason: row.detail(row.reason)}
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

// planUnclaimedAssignmentSet selects the unclaimed assignment and shift branches and
// classifies each one. It excludes every assignment record, checked-out ref, and the
// default branch.
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
	set := unclaimedAssignmentSet{options: options}
	for _, branch := range branches {
		ref := "refs/heads/" + branch
		reason, owned := unclaimedBranchReason(ref)
		if !owned || protected[ref] {
			continue
		}
		set.rows = append(set.rows, unclaimedAssignmentBranch{ref: ref, reason: reason})
	}
	sort.Slice(set.rows, func(i, j int) bool { return set.rows[i].ref < set.rows[j].ref })
	refs := make([]string, len(set.rows))
	for i, row := range set.rows {
		refs[i] = row.ref
	}
	verdicts, err := classifyUnclaimedRefs(root, assignments, protected, defaultBranch, refs)
	if err != nil {
		return unclaimedAssignmentSet{}, err
	}
	parts := [][]byte{[]byte(unclaimedAssignmentFingerprintVersion), []byte("discard-branch=true"), []byte(fmt.Sprintf("unclaimed=%t", options.Unclaimed))}
	for i := range set.rows {
		set.rows[i].refVerdict = verdicts[i]
		row := set.rows[i]
		parts = append(parts, []byte(row.ref), []byte(row.oid), []byte(row.class), []byte(row.holder))
	}
	if len(set.rows) == 0 || set.faulted() {
		return set, nil
	}
	set.fingerprint = fingerprintParts(parts...)
	return set, nil
}

// renderUnclaimedAssignmentSet prints the plan rows, then the apply command when a row
// removes.
func renderUnclaimedAssignmentSet(stdout io.Writer, set unclaimedAssignmentSet) error {
	plans := make([]CleanupPlan, 0, len(set.rows))
	removes := false
	for _, row := range set.rows {
		plans = append(plans, row.plan(set.fingerprint))
		removes = removes || row.action().Removes()
	}
	if err := renderCleanups(stdout, plans); err != nil {
		return err
	}
	if set.faulted() {
		return errFaultedUnclaimedRef
	}
	if !removes {
		return nil
	}
	// The apply command is the re-plan command plus the digest this plan authorizes.
	arguments := append(unclaimedReplan(set.options), axi.KnownArgument("--apply"), axi.KnownArgument(set.fingerprint))
	help, err := axi.RenderHelp([]axi.Action{axi.ExecutableInvocation("apply the unclaimed branch plan", arguments...)})
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(stdout, help)
	return err
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

// unclaimedSelector is the selector that names this mode.
const unclaimedSelector = "--unclaimed"

// unclaimedReplan is this mode's own re-plan command, beside the landed selector's. It takes
// the caller's options so a later grammar change cannot leave the rendered command behind.
// Today the two are pinned equal: valid() requires --discard-branch and refuses the other
// two modifiers under --unclaimed, so every valid invocation carries unclaimedOptions().
func unclaimedReplan(options CleanupOptions) []axi.InvocationArgument {
	return cleanArguments(options, unclaimedSelector)
}

// UnclaimedPlanCommand is the plan-only command that status routes an unclaimed ref to. It
// carries the modifiers and the selector of the re-plan this mode's refusal offers.
func UnclaimedPlanCommand() string {
	words := append([]string{"bench", "worktree", "clean"}, cleanupModifierFlags(unclaimedOptions())...)
	return strings.Join(append(words, unclaimedSelector), " ")
}

// applyUnclaimedAssignmentSet deletes each landed and subsumed branch at the exact object
// the plan named, and it reports a retained row as planned. It reports the outcome rows
// alone. A stale refusal carries no rows, because the refusal row is this command surface's
// own spelling and the caller renders it; a row returned here would be a second derivation
// the caller discards.
func applyUnclaimedAssignmentSet(j joins, root string, set unclaimedAssignmentSet, options CleanupOptions) ([]CleanupPlan, error) {
	// The window that refusal exists for: a branch can enter or leave the namespace between
	// the caller's plan read and this re-plan. The boundary is nil in production; it lets a
	// test stand in that window. This mode holds no lock across the window, so a concurrent
	// writer is refused after the fact rather than excluded, and the step token says so.
	if err := hit(j.cleanupBoundary, StepUnlockedReplan); err != nil {
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
		if !planned.action().Removes() {
			plans = append(plans, planned.plan(set.fingerprint))
			continue
		}
		if err := git.DeleteBranchExact(root, planned.ref, planned.oid); err != nil {
			return append(plans, CleanupPlan{Target: planned.ref, Action: ActionError, Tracked: "unclaimed", ignoredSummary: "none", Recovery: "none", Fingerprint: set.fingerprint, Reason: err.Error()}), err
		}
		plans = append(plans, CleanupPlan{Target: planned.ref, Action: ActionRemoved, Tracked: "unclaimed", ignoredSummary: "none", Recovery: "none", Fingerprint: set.fingerprint})
	}
	return plans, nil
}

// unclaimedClassOrder is the order the class counts are reported in.
var unclaimedClassOrder = []refClass{classLanded, classSubsumed, classUnique}

// UnclaimedClassCount is how many rows of the unclaimed plan carry one class.
type UnclaimedClassCount struct {
	Class string
	Count int
}

// UnclaimedRefCounts is the unclaimed plan by class, in unclaimedClassOrder. Faulted counts
// the error rows, which carry no class and make the plan refuse its apply. Refs holds the
// ref of every row, so a reader that counts branches by another proof can leave them out.
type UnclaimedRefCounts struct {
	Classes []UnclaimedClassCount
	Faulted int
	Refs    map[string]bool
}

// Rows is how many rows the plan holds.
func (counts UnclaimedRefCounts) Rows() int {
	rows := counts.Faulted
	for _, class := range counts.Classes {
		rows += class.Count
	}
	return rows
}

// CountUnclaimedRefs counts the classes of the same plan the unclaimed clean prints, so
// status and the plan cannot disagree. A row with a class outside unclaimedClassOrder is an
// error, not a silent unique count.
func CountUnclaimedRefs(root string) (UnclaimedRefCounts, error) {
	set, err := planUnclaimedAssignmentSet(root, unclaimedOptions())
	if err != nil {
		return UnclaimedRefCounts{}, err
	}
	counts := UnclaimedRefCounts{Classes: make([]UnclaimedClassCount, len(unclaimedClassOrder)), Refs: make(map[string]bool, len(set.rows))}
	for i, class := range unclaimedClassOrder {
		counts.Classes[i].Class = string(class)
	}
	for _, row := range set.rows {
		counts.Refs[row.ref] = true
		if row.fault != "" {
			counts.Faulted++
			continue
		}
		i := slices.Index(unclaimedClassOrder, row.class)
		if i < 0 {
			return UnclaimedRefCounts{}, fmt.Errorf("unclaimed ref %s has unknown class %q", row.ref, row.class)
		}
		counts.Classes[i].Count++
	}
	return counts, nil
}
