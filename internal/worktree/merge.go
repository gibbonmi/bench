// The worktree merge verb: it composes one owned commit into one active owned worktree,
// so a retained worktree reaches a new base and folds a sibling's tip through Bench.
package worktree

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/usage"
)

var mergeGrammar = usage.Grammar{
	Cmd:     "bench worktree merge",
	Help:    "usage: " + usage.WorktreeMerge,
	MinArgs: 1,
	MaxArgs: 1,
	Flags: []usage.Flag{
		{Name: "--from", HasValue: true, NoEmptyValue: true, Required: true},
	},
}

// reconcileMergeCheckout moves the target checkout onto the commit the merge published.
// It is the bare reset the boundary the verb reports stands on: the pre-check and the
// fingerprint recheck already proved the checkout holds nothing the reset would discard.
func reconcileMergeCheckout(path, tip string) error {
	_, err := exec.Command("git", "-C", path, "reset", "--merge", tip).CombinedOutput()
	return err
}

// MergeCommand composes one commit into one active owned worktree by merge. The incoming
// commit is a commit in the default branch's history or a sibling assignment's branch
// tip; no other object reaches the lane.
func MergeCommand(root, home string, args []string, stdout, stderr io.Writer) int {
	return mergeWith(defaultJoins(), root, home, args, stdout, stderr)
}

// mergeWith is MergeCommand with the seam set resolved explicitly at the caller's
// boundary.
func mergeWith(j joins, root, _ string, args []string, stdout, stderr io.Writer) int {
	parsed, line, code := usage.Parse(mergeGrammar, args)
	if line != "" {
		fmt.Fprintln(stderr, line)
		return code
	}
	assignments, err := intent.Assignments(root)
	if err != nil {
		return landRefusal(stdout, "assignment ledger is unreadable")
	}
	target, err := mergeTarget(root, assignments, parsed.Positionals[0])
	if err != nil {
		return landRefusalError(stdout, err)
	}
	spelling := parsed.Flags["--from"]
	incoming, err := mergeIncoming(root, assignments, target, spelling)
	if err != nil {
		return landRefusalError(stdout, err)
	}
	previous, err := mergeTargetTip(root, target)
	if err != nil {
		return landRefusalError(stdout, err)
	}
	fingerprint, err := landing.CheckoutFingerprint(target.Worktree)
	if err != nil {
		return landRefusal(stdout, "merge target checkout fingerprint is unreadable")
	}
	owner, err := mergeOwner(j, target.Worktree, previous)
	if err != nil {
		return landRefusal(stdout, err.Error())
	}
	result, err := owner.Merge(context.Background(), landing.MergeRequest{
		Root: root, Branch: target.Branch, PreviousTip: previous, Incoming: incoming,
		Worktree: target.Worktree, Fingerprint: fingerprint,
		Subject: mergeSubject(spelling, incoming, target.Label),
		Stdout:  stdout, Stderr: stderr,
	})
	if err != nil {
		// A composition conflict carries its paths typed, so it renders through the
		// refusal record's path table rather than as a bare sentence.
		var conflict landing.ConflictError
		if errors.As(err, &conflict) {
			return landRefusalError(stdout, refusalError{refusal{detail: conflict.Error(), paths: conflict.Paths}})
		}
		return landRefusalError(stdout, err)
	}
	if len(result.Resolved) > 0 {
		fmt.Fprintf(stderr, "merge composition{resolved=%s}\n", sanitize.Controls(strings.Join(result.Resolved, ",")))
	}
	record := fmt.Sprintf("merged{worktree=%s,from=%s,kind=%s,previous_tip=%s,tip=%s,tree=%s",
		target.ID, incoming, result.Kind, result.PreviousTip, result.Tip, result.Tree)
	// A `current` target changed nothing, so there is nothing for the checkout to catch
	// up with; only a published tip needs the reconcile.
	if result.Kind != landing.MergeKindCurrent {
		if err := j.mergeReconcile(target.Worktree, result.Tip); err != nil {
			fmt.Fprintln(stdout, record+",next="+sanitize.Controls(mergeReconcileNext(target, result.Tip))+"}")
			return 3
		}
	}
	fmt.Fprintln(stdout, record+"}")
	return 0
}

// mergeSubject derives the one message the verb publishes, so no `-m` exists and the log
// reads one way. The spelling is the operand as typed, because that is what the operator
// will search the log for.
func mergeSubject(spelling, incoming, label string) string {
	return "merge: compose " + spelling + " " + incoming[:8] + " into " + label
}

// mergeReconcileNext names the repair the exit-3 boundary leaves the operator. A path
// that is not line-safe takes the pointer form every next= uses.
func mergeReconcileNext(target intent.Assignment, tip string) string {
	if !lineSafe(target.Worktree) {
		return "bench worktree exec " + target.ID + " -- git reset --merge " + tip
	}
	return "git -C " + sanitize.ShellQuote(target.Worktree) + " reset --merge " + tip
}

// mergeOnAssignmentBranch proves one checkout is on its assignment's branch. It runs
// before the identity bundle reads the same registration, because the bundle's
// registration sentence reports the fact without naming the ref the operator restores.
func mergeOnAssignmentBranch(a intent.Assignment, detail string) error {
	branch, err := git.Output("-C", a.Worktree, "symbolic-ref", "--quiet", "HEAD")
	if err != nil || branch != a.Branch {
		return refusalError{refusal{detail: detail, observed: branch, wanted: a.Branch}}
	}
	return nil
}

// mergeTarget resolves the operand to the assignment record the verb needs: the id for
// the record, the label for the subject, and the branch for the compare-and-swap. The
// identity bundle is this verb's whole authority, so it is proved before anything moves.
func mergeTarget(root string, assignments []intent.Assignment, operand string) (intent.Assignment, error) {
	if !lineSafe(operand) {
		return intent.Assignment{}, refusalError{refusal{detail: "target contains control characters"}}
	}
	selected, err := selectAssignment(assignments, operand)
	if err != nil {
		return intent.Assignment{}, refusalError{refusal{detail: err.Error()}}
	}
	if err := mergeOnAssignmentBranch(selected, "merge target is not on its assignment branch"); err != nil {
		return intent.Assignment{}, err
	}
	if err := identityBundleRefusal(root, selected.Worktree, selected, landingActiveState); err != nil {
		return intent.Assignment{}, err
	}
	return selected, nil
}

// mergeTargetTip proves the moving identity the compare-and-swap depends on: the
// assignment branch's tip is the checkout's HEAD. It then refuses a dirty target,
// because the reconcile after publication resets it.
func mergeTargetTip(root string, a intent.Assignment) (string, error) {
	head, err := git.Output("-C", a.Worktree, "rev-parse", "HEAD^{commit}")
	if err != nil {
		return "", refusalError{refusal{detail: "merge target checkout has no commit"}}
	}
	tip, err := git.Output("-C", root, "rev-parse", "--verify", a.Branch+"^{commit}")
	if err != nil || tip != head {
		return "", refusalError{refusal{detail: "merge target branch tip is not the checkout HEAD", observed: head, wanted: tip}}
	}
	if err := mergeCheckoutClean(a.Worktree, "merge target checkout is not clean", ""); err != nil {
		return "", err
	}
	return tip, nil
}

// mergeCheckoutClean refuses any status entry, untracked included. Ignored residue is
// excluded, because a worktree's build output is not uncommitted work.
func mergeCheckoutClean(path, detail, next string) error {
	raw, err := git.Raw("-C", path, "status", "--porcelain=v1", "-z", "--no-renames", "--untracked-files=all")
	if err != nil {
		return refusalError{refusal{detail: "checkout status is unreadable"}}
	}
	entries, err := git.ParsePorcelainZStrict(raw)
	if err != nil {
		return refusalError{refusal{detail: "checkout status is unreadable"}}
	}
	if len(entries) == 0 {
		return nil
	}
	paths := make([]string, 0, len(entries))
	for _, entry := range entries {
		paths = append(paths, entry.Path)
	}
	return refusalError{refusal{detail: detail, next: next, paths: paths}}
}

// mergeIncoming resolves `--from` in the two lookups the bootstrap authority allows: a
// sibling assignment's branch tip, and a commit in the default branch's history. A value
// both lookups answer is ambiguous, because a first-match resolver would merge whichever
// lookup happened to run first.
func mergeIncoming(root string, assignments []intent.Assignment, target intent.Assignment, from string) (string, error) {
	if !lineSafe(from) {
		return "", refusalError{refusal{detail: "--from contains control characters"}}
	}
	sibling, siblingID, hasSibling, err := mergeSiblingTip(root, assignments, target, from)
	if err != nil {
		return "", err
	}
	commit, resolved, owned := mergeDefaultBranchCommit(root, from)
	switch {
	case hasSibling && owned:
		// The assignment id, not its tip, is the unambiguous respelling of the label, so
		// the operator can retype either lookup exactly.
		return "", refusalError{refusal{detail: "--from names both an assignment and a commit", observed: from, wanted: siblingID + " or " + commit}}
	case hasSibling:
		return sibling, nil
	case owned:
		return commit, nil
	case resolved:
		return "", refusalError{refusal{detail: "--from is outside the default branch's history and is no sibling tip", observed: commit}}
	}
	return "", refusalError{refusal{detail: "--from names no assignment and no commit", observed: from}}
}

// mergeSiblingTip answers the assignment lookup. A sibling contributes its committed
// branch tip alone: `bench commit` stays the one snapshot composer, so a detached or
// dirty sibling refuses rather than have its uncommitted work silently dropped.
func mergeSiblingTip(root string, assignments []intent.Assignment, target intent.Assignment, from string) (tip, id string, ok bool, err error) {
	selected, err := selectAssignment(assignments, from)
	if err != nil {
		return "", "", false, nil
	}
	if selected.ID == target.ID {
		return "", "", false, refusalError{refusal{detail: "--from resolves to the target itself", observed: selected.ID}}
	}
	if err := mergeOnAssignmentBranch(selected, "sibling is not on its assignment branch"); err != nil {
		return "", "", false, err
	}
	if err := identityBundleRefusal(root, selected.Worktree, selected, landingActiveState); err != nil {
		return "", "", false, err
	}
	tip, err = git.Output("-C", root, "rev-parse", "--verify", selected.Branch+"^{commit}")
	if err != nil {
		return "", "", false, refusalError{refusal{detail: "sibling assignment branch has no commit", observed: selected.Branch}}
	}
	head, err := git.Output("-C", selected.Worktree, "rev-parse", "HEAD^{commit}")
	if err != nil || head != tip {
		return "", "", false, refusalError{refusal{detail: "sibling checkout is not at its branch tip", observed: head, wanted: tip}}
	}
	if err := mergeCheckoutClean(selected.Worktree, "sibling checkout is not clean", "bench worktree exec "+selected.ID+" -- bench commit"); err != nil {
		return "", "", false, err
	}
	return tip, selected.ID, true, nil
}

// mergeDefaultBranchCommit answers the commit lookup. The two results are separate
// facts: a spelling Git can peel, and a commit the default branch's history contains.
// Only the second authorizes the lane to execute the tree.
func mergeDefaultBranchCommit(root, from string) (commit string, resolved, owned bool) {
	commit, err := git.Output("-C", root, "rev-parse", "--verify", "--quiet", from+"^{commit}")
	if err != nil || commit == "" {
		return "", false, false
	}
	branch, ok := git.ResolvedDefault(root)
	if !ok {
		return commit, true, false
	}
	tip, err := git.Output("-C", root, "rev-parse", "--verify", branch+"^{commit}")
	if err != nil {
		return commit, true, false
	}
	contains, err := authorization.IsAncestor(root, commit, tip)
	return commit, true, err == nil && contains
}

// mergeOwner resolves the authority the composed tree is graded under, the way
// `bench commit` resolves it, but for the target worktree. A root with no declared lane
// keeps the whole-project gate.
func mergeOwner(j joins, target, previous string) (landing.Owner, error) {
	checks, kit, err := j.mergeLane(target)
	if err != nil {
		return landing.Owner{}, err
	}
	if checks == nil {
		return landing.New(), nil
	}
	// The previous tip, not a path list, is what the authority needs: it grades the
	// composed tree it already holds, so the verb never composes a second time to learn
	// which prose the merge brought in.
	return landing.NewLane(authorization.LaneAuthority{Checks: checks, Kit: kit, PreviousTip: previous}), nil
}
