package worktree

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

// refClass is the class the unclaimed plan gives one unrecorded Bench-namespace branch.
type refClass string

const (
	// classLanded passes one of the four landed proofs.
	classLanded refClass = "landed"
	// classSubsumed has a tip that equals, or is a strict ancestor of, a holder's tip.
	classSubsumed refClass = "subsumed"
	// classUnique is every other ref: it carries content main lacks and nothing holds it.
	classUnique refClass = "unique"
)

// The detail spellings of a classified row. A unique row ends with uniqueRetainedDetail.
const (
	classField           = "class="
	holderField          = " holder="
	uniqueRetainedDetail = "retained: content main lacks"
)

// errFaultedUnclaimedRef refuses a plan that holds an error row.
var errFaultedUnclaimedRef = errors.New("an unclaimed ref has an error row; the plan offers no apply")

// refVerdict is what the class function decides for one ref. A fault names a symref or a ref
// whose tip is not a commit; that ref has no class and no holder.
type refVerdict struct {
	oid, fault string
	class      refClass
	holder     string
}

// action is the cleanup action a classified row carries.
func (v refVerdict) action() CleanupAction {
	switch {
	case v.fault != "":
		return ActionError
	case v.class == classUnique:
		return ActionRetain
	default:
		return ActionDiscardRemove
	}
}

// detail is the row's detail cell: the class, the holder of a subsumed row, then the removal
// reason of a removing row or the retained text of a unique row.
func (v refVerdict) detail(reason string) string {
	if v.fault != "" {
		return v.fault
	}
	text := classField + string(v.class)
	if v.class == classSubsumed {
		text += holderField + v.holder
	}
	if v.class == classUnique {
		return text + "; " + uniqueRetainedDetail
	}
	return text + "; " + reason
}

// holderTip is one ref that can hold a subsumed ref, at its resolved tip.
type holderTip struct{ ref, oid string }

// classifyUnclaimedRefs gives each of the sorted unrecorded refs one class and one holder.
// It reads refs only. It calls the shared landed proof for every ref, so this sweep and the
// landing prune cannot disagree. A landed recorded branch holds nothing, because its later
// retirement would leave a ref beneath it with no handle. A symref faults, because a delete
// through it removes its target.
func classifyUnclaimedRefs(root string, assignments []intent.Assignment, protected map[string]bool, defaultBranch string, refs []string) ([]refVerdict, error) {
	verdicts := make([]refVerdict, len(refs))
	var recorded []holderTip
	for _, assignment := range assignments {
		if assignment.State != intent.StateActive && assignment.State != intent.StateCleanupPending {
			continue
		}
		oid, err := git.Output("-C", root, "rev-parse", "--verify", "--quiet", assignment.Branch+"^{commit}")
		if err != nil {
			continue
		}
		landed, _, err := git.LandedInDefault(root, assignment.Branch, defaultBranch)
		if err != nil {
			return nil, fmt.Errorf("git landedness %s: %w", assignment.Branch, err)
		}
		if !landed {
			recorded = append(recorded, holderTip{assignment.Branch, oid})
		}
	}
	sort.Slice(recorded, func(i, j int) bool { return recorded[i].ref < recorded[j].ref })
	landed := make([]bool, len(refs))
	for i, ref := range refs {
		if protected[ref] {
			return nil, fmt.Errorf("protected branch %s reached the unclaimed classes", ref)
		}
		if target, err := git.Output("-C", root, "symbolic-ref", "--quiet", ref); err == nil {
			verdicts[i].fault = fmt.Sprintf("%s is a symref to %s", ref, target)
			continue
		}
		oid, err := git.Output("-C", root, "rev-parse", "--verify", "--quiet", ref+"^{commit}")
		if err != nil {
			kind, _ := git.Output("-C", root, "cat-file", "-t", ref)
			verdicts[i].fault = fmt.Sprintf("%s names a %s object that is not a commit", ref, kind)
			continue
		}
		verdicts[i].oid = oid
		if landed[i], _, err = git.LandedInDefault(root, ref, defaultBranch); err != nil {
			return nil, fmt.Errorf("git landedness %s: %w", ref, err)
		}
	}
	var open []int
	for i, ref := range refs {
		switch {
		case verdicts[i].fault != "":
		case landed[i]:
			verdicts[i].class = classLanded
		default:
			if holder, ok := firstReaching(root, recorded, ref, verdicts[i].oid); ok {
				verdicts[i].class, verdicts[i].holder = classSubsumed, holder
				continue
			}
			open = append(open, i)
		}
	}
	roots, err := uniqueRoots(root, refs, verdicts, open)
	if err != nil {
		return nil, err
	}
	for _, i := range open {
		if holder, ok := firstReaching(root, roots, refs[i], verdicts[i].oid); ok {
			verdicts[i].class, verdicts[i].holder = classSubsumed, holder
			continue
		}
		verdicts[i].class = classUnique
	}
	return verdicts, nil
}

// uniqueRoots names the unique roots among the refs no other holder reaches: each tip no
// other open tip descends from, held by the lexically first open ref at that tip.
func uniqueRoots(root string, refs []string, verdicts []refVerdict, open []int) ([]holderTip, error) {
	if len(open) == 0 {
		return nil, nil
	}
	seen := map[string]bool{}
	args := []string{"-C", root, "merge-base", "--independent"}
	for _, i := range open {
		if !seen[verdicts[i].oid] {
			seen[verdicts[i].oid] = true
			args = append(args, verdicts[i].oid)
		}
	}
	out, err := git.Output(args...)
	if err != nil {
		return nil, fmt.Errorf("git independent unclaimed tips: %w", err)
	}
	maximal := map[string]bool{}
	for _, oid := range strings.Split(out, "\n") {
		maximal[oid] = true
	}
	var roots []holderTip
	for _, i := range open {
		if maximal[verdicts[i].oid] {
			roots = append(roots, holderTip{refs[i], verdicts[i].oid})
			delete(maximal, verdicts[i].oid)
		}
	}
	return roots, nil
}

// firstReaching returns the first holder, other than the ref itself, whose tip equals oid or
// descends from it. An ancestry check that cannot run reads as no reach, so the ref stays
// unique and the sweep retains it.
func firstReaching(root string, holders []holderTip, ref, oid string) (string, bool) {
	for _, holder := range holders {
		if holder.ref == ref {
			continue
		}
		if holder.oid == oid || git.OK("-C", root, "merge-base", "--is-ancestor", oid, holder.oid) {
			return holder.ref, true
		}
	}
	return "", false
}
