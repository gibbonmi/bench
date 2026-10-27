// Stack re-run composition: the re-run that ends the route of each stack-raised refusal.
package worktree

import (
	"slices"

	"github.com/gibbonmi/bench/internal/refusalroute"
)

// stackRepairedTip is the slot of a stack override whose source a repair moves. The
// operator fills the tip that the repair committed.
const stackRepairedTip = "<repaired-tip>"

// stackRoute is what the re-run of a stack refusal reads: the caller's stack file, the read
// members and overrides, and the worktree path and assignment of each member. A member
// whose path did not resolve keeps its file value, and a member whose assignment did not
// resolve has an empty assignment.
type stackRoute struct {
	file        string
	input       StackInput
	paths       []string
	assignments []string
}

// stackRerun is what one refusal tells the composer. refused is the index of the refused
// member, or -1 when the refusal names no member. repairs states that the route commits in
// the refused member's source. read is the tip that the tree holds for a source-tip
// mismatch. removed holds the members that the route removes from the stack file, and reset
// holds the members whose source_tip the route sets in it.
type stackRerun struct {
	refused        int
	repairs        bool
	read           string
	removed, reset []int
}

// stackUnattributed is the re-run fact of a refusal that names no member.
var stackUnattributed = stackRerun{refused: -1}

// stackRefusedBy is the re-run fact of the refusal err that member i raised. The face of err
// decides whether the route repairs the source, and a source-tip mismatch carries the tip
// that the tree holds.
func stackRefusedBy(i int, err error) stackRerun {
	face := landingFaceName(err)
	fact := stackRerun{refused: i, repairs: repairsSource(face)}
	if _, read, ok := sourceTips(raisedRefusal(err)); ok && face == faceSourceTipMismatch {
		fact.read = read
	}
	return fact
}

// stackOverride is one `--source-tip` value of a re-run. A repaired override prints the
// repaired-tip slot instead of a commit.
type stackOverride struct {
	request, commit string
	repaired        bool
}

func (o stackOverride) arg() string {
	if o.repaired {
		return refusalroute.Arg("request", o.request) + "=" + stackRepairedTip
	}
	return refusalroute.Arg("request", o.request) + "=" + refusalroute.Arg("commit", o.commit)
}

// rerun composes the re-run of one stack refusal by the stack's rules, in order. It starts
// from the overrides of the refused call and drops the override of each removed and each
// reset member. It sets the repaired slot of a refused member whose source the route
// repairs, or the read tip of a refused member whose tip mismatched. Two or more remaining
// members re-run the stack, and one remaining member re-runs its single landing.
func (s stackRoute) rerun(fact stackRerun) string {
	overrides := make([]stackOverride, 0, len(s.input.Overrides)+1)
	for _, o := range s.input.Overrides {
		overrides = append(overrides, stackOverride{request: o.Request, commit: o.Commit})
	}
	for _, i := range slices.Concat(fact.removed, fact.reset) {
		overrides = slices.DeleteFunc(overrides, func(o stackOverride) bool { return o.request == s.input.Members[i].Request })
	}
	if fact.refused >= 0 && (fact.repairs || fact.read != "") {
		set := stackOverride{request: s.input.Members[fact.refused].Request, commit: fact.read, repaired: fact.repairs}
		if at := slices.IndexFunc(overrides, func(o stackOverride) bool { return o.request == set.request }); at >= 0 {
			overrides[at] = set
		} else {
			overrides = append(overrides, set)
		}
	}
	var remaining []int
	for i := range s.input.Members {
		if !slices.Contains(fact.removed, i) {
			remaining = append(remaining, i)
		}
	}
	if len(remaining) == 1 {
		return s.singleRerun(remaining[0], overrides)
	}
	command := "bench worktree land --stack " + refusalroute.Arg("file", s.file)
	for _, o := range overrides {
		command += " --source-tip " + o.arg()
	}
	return command
}

// singleRerun is the single landing of member i, at the override that remains for it or
// else at its file tip.
func (s stackRoute) singleRerun(i int, overrides []stackOverride) string {
	m := s.input.Members[i]
	tipFlag := landingSourceTipFlag(m.SourceTip)
	if at := slices.IndexFunc(overrides, func(o stackOverride) bool { return o.request == m.Request }); at >= 0 {
		tipFlag = landingSourceTipFlag(overrides[at].commit)
		if overrides[at].repaired {
			tipFlag = repairedSourceTipFlag
		}
	}
	return landingRerunAt(m.Request, m.Base, tipFlag, m.Spec, s.paths[i], s.assignments[i])
}
