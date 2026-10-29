// Stack resume flow: the proofs that a published merge is the publication of a stack
// file's prefix, the post-publication work that the stack did not finish, and the terminal
// record that a stack and its resume print.
package worktree

import (
	"errors"
	"fmt"
	"io"
	"slices"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/usage"
	"github.com/gibbonmi/bench/internal/worktree/landingpolicy"
)

// stackResumeGrammar is the stack grammar with the published merge to resume.
var stackResumeGrammar = usage.Grammar{
	Cmd:   usage.CommandName(usage.WorktreeLandStackResume),
	Help:  "usage: " + usage.WorktreeLandStackResume,
	Flags: slices.Concat([]usage.Flag{{Name: "--resume", HasValue: true, NoEmptyValue: true, Required: true}}, stackLandGrammar.Flags),
}

// resumeStack is the stack resume form of the land leaf. Every proof runs before the first
// step, so a refused resume writes nothing. It finishes the steps after the publication
// once for the published merge, releases each published member that is still active, and
// leaves each later member active.
func resumeStack(j joins, a ambient, root string, args []string, stdout, stderr io.Writer) int {
	parsed, line, code := usage.Parse(stackResumeGrammar, args)
	if line != "" {
		fmt.Fprintln(stderr, line)
		return code
	}
	// A member path that does not resolve reaches the assignment proofs as given, and they
	// refuse it.
	route, _, err := readStackRoute(a, root, parsed)
	if err != nil {
		return stackUsageRefusal(stderr, stackResumeGrammar, err)
	}
	published := expandIdentity(root, parsed.Flags["--resume"])
	rerun := route.resume(published)
	members := route.input.Members
	bases, tips := make([]string, len(members)), make([]string, len(members))
	for i, m := range members {
		bases[i], tips[i] = expandIdentity(root, m.Base), expandIdentity(root, m.EffectiveTip)
	}
	destination, branch, marker, err := landingDestinationIdentity(root)
	if err != nil {
		return landRefusalError(stdout, landingFaceRoute(err, rerun, false))
	}
	k, err := stackPublishedCount(root, published, tips)
	if err != nil {
		return landRefusalError(stdout, landingFaceRoute(err, rerun, false))
	}
	// The walk proves each published member merge as one resume proves its merge, from the
	// top merge down. The first parent of each merge is the merge of the member before it,
	// and the last first parent is the destination base of the stack.
	merge := published
	for i := k - 1; i >= 0; i-- {
		_, _, firstParent, _, err := resumePublished(j, root, destination, merge, bases[i], tips[i], members[i].Spec)
		if err != nil {
			return landRefusalError(stdout, stackMemberRefusal(i, "", landingFaceRoute(err, rerun, false)))
		}
		merge = firstParent
	}
	stackBase := merge
	var active []int
	for i := range k {
		_, live, skipped, err := resumeAssignment(j, a.run, root, route.paths[i], members[i].Request, tips[i], bases[i], landingSlug(members[i].Spec))
		if err == nil && !live {
			_, err = terminalResumeReceipt(a.run, root, route.paths[i], members[i].Request, tips[i])
		}
		if err != nil {
			return landRefusalError(stdout, stackMemberRefusal(i, "", landingFaceRoute(err, rerun, skipped)))
		}
		if live {
			active = append(active, i)
		}
	}
	end := stackEnd{destinationBase: stackBase, published: published, landed: k}
	refused, step := resumePublishedSteps(j, a.run, root, destination, branch, marker, published, stackBase, rerun)
	if refused != nil {
		return landRefusalError(stdout, refused)
	}
	for _, i := range active {
		if step == "" && !releaseLanded(j, a, root, members[i].Request, route.paths[i], stderr) {
			step = "release"
		}
	}
	if step == "" {
		step = runLandedEffects(j, a, root, stackBase, false, stdout, stderr)
	}
	return printStackEnd(stdout, route, end.after(step))
}

// stackEnd is the terminal record of a stack and its exit code. step names the step after
// the publication that stopped, empty when none stopped.
type stackEnd struct {
	destinationBase, published, worktree, step string
	landed, exit                               int
}

// after is end with the step after the publication that stopped, empty when none stopped,
// and the worktree state and exit that the step decides.
func (e stackEnd) after(step string) stackEnd {
	outcome := landingpolicy.Terminal(landingpolicy.TerminalFacts{FailedStep: step, Active: true})
	e.step, e.worktree, e.exit = step, outcome.WorktreeState, outcome.ExitCode
	return e
}

// printStackEnd prints the land-incomplete record of a stopped step, whose route is the
// stack resume of the published merge with each override of the call, then the stack
// record. It answers the exit of end.
func printStackEnd(stdout io.Writer, route stackRoute, end stackEnd) int {
	if end.step != "" {
		stopped := refusal{detail: "landing stack stopped after publication at the " + end.step + " step"}
		landRefusalError(stdout, landingFaceRefusal(faceLandIncomplete, stopped, route.resume(end.published), ""))
	}
	fmt.Fprintf(stdout, "stack{destination_base=%s,published_commit=%s,landed=%d,worktree=%s}\n", end.destinationBase, end.published, end.landed, end.worktree)
	return end.exit
}

// stackPublishedCount is the published count of a stack: the 1-based index of the one
// member whose effective tip, among tips, is the second parent of published.
func stackPublishedCount(root, published string, tips []string) (int, error) {
	second := secondParent(root, published)
	var matches []int
	for i, tip := range tips {
		if second != "" && tip == second {
			matches = append(matches, i+1)
		}
	}
	if len(matches) != 1 {
		return 0, errStackPublishedCount
	}
	return matches[0], nil
}

// errStackPublishedCount refuses a published commit that gives no single published count.
var errStackPublishedCount = errors.New("the second parent of the published commit is not the tip of exactly one stack member")

// secondParent is the second parent of commit, or the empty identity when it has none.
func secondParent(root, commit string) string {
	parent, err := git.ResolveCommit(root, commit+"^2", "--quiet")
	if err != nil {
		return ""
	}
	return parent
}
