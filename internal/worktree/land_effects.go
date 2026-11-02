// Landing effects: the work a landing owns after publication and release, each with its
// own recorded result and its own resume path.
package worktree

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/gibbonmi/bench/internal/brokermanifest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// An effect's result word. A landing reports one of these for every effect it owns, so a
// reader tells a skipped effect from a finished one without a second read of the tree.
const (
	effectComplete      = "complete"
	effectFailed        = "failed"
	effectPending       = "pending"
	effectNotApplicable = "not-applicable"
	effectHeld          = "held"
)

// refreshEffect and cleanupEffect name the two effects in the effects row and in the
// incomplete record's step cell, so `incomplete:refresh` and `incomplete:cleanup` cannot
// disagree with the row cell beside them.
const (
	refreshEffect = "refresh"
	cleanupEffect = "cleanup"
)

// effectRow is one effect's line in the landing's effects table.
type effectRow struct {
	effect string
	result string
}

// landedAfterEffects runs the landing's post-release effects at root, prints their table,
// and renders the landing's terminal record. It is the one place the effects and the
// record meet, so the first run and the resume state the same thing about the same tree.
// A failed effect reaches the existing incomplete render under its own step name, which
// carries the resume command every other incomplete step carries.
//
// brokerChanged states that the landed diff changes the promotion broker's build inputs.
// The broker notice prints after the effects row, so it reads the refresh result. A resume
// has no source worktree to read, so it passes false.
func landedAfterEffects(j joins, a ambient, root string, result landing.ReviewedResult, specArg, path, assignment string, active, brokerChanged bool, records int, stdout, stderr io.Writer) int {
	if failed := runLandedEffects(j, a, root, result.DestinationBase, brokerChanged, stdout, stderr); failed != "" {
		return landedIncomplete(stdout, result, specArg, path, assignment, failed, records)
	}
	return landedComplete(stdout, result, active, records)
}

// runLandedEffects runs the two landing effects once at root for a publication over
// destinationBase, prints their table and the broker notice, and answers the failed effect.
// It answers the empty name when no effect failed. A first run, a resume, and a stack all
// run the effects through it.
func runLandedEffects(j joins, a ambient, root, destinationBase string, brokerChanged bool, stdout, stderr io.Writer) string {
	refresh := refreshBroker(j, a.run, root, stderr)
	// A failed effect stops every later effect, so the cleanup never starts and never
	// touches a checkout. It reports pending, which is the word for an effect that has
	// not run rather than one that ran and did nothing.
	cleanup := effectPending
	if refresh != effectFailed {
		cleanup = cleanLandedSiblings(j, a, root, destinationBase, stderr)
	}
	rows := []effectRow{{effect: refreshEffect, result: refresh}, {effect: cleanupEffect, result: cleanup}}
	fmt.Fprintf(stdout, "effects[%d]{effect,result}:\n", len(rows))
	for _, row := range rows {
		fmt.Fprintf(stdout, "  %s,%s\n", row.effect, row.result)
	}
	if brokerChanged {
		if notice := brokerChangeNotice(a.kit, root, refresh); notice != "" {
			fmt.Fprintln(stderr, notice)
		}
	}
	for _, row := range rows {
		if row.result == effectFailed {
			return row.effect
		}
	}
	return ""
}

// publishedSteps runs the steps that follow a publication at root: the green marker
// advance, then the steps after the marker. It answers the name of the step that failed,
// or the empty name. A first run and a stack both run these steps once for the commit that
// they published.
func publishedSteps(j joins, run *git.Run, root, branch string, result landing.ReviewedResult, priorMarker string) string {
	if err := advanceLandedMarker(root, branch, result, priorMarker); err != nil {
		return "marker"
	}
	return stepsAfterMarker(j, run, root, result.Commit, result.Commit, result.DestinationBase)
}

// stepsAfterMarker runs the steps that follow the marker step for the publication of
// published at root, in order: the commitment reconcile, the destination checkout reconcile
// to destination, and the branch prune. destination is the destination tip now. It answers
// the name of the step that failed, or the empty name. Every landing and every resume runs
// these steps through it.
func stepsAfterMarker(j joins, run *git.Run, root, destination, published, destinationBase string) string {
	if err := j.reconcileCommitment(commitrepo.Store{Root: root, Run: run}); err != nil {
		return "commitment"
	}
	if err := reconcileLandingDestination(j, root, destination, published, destinationBase); err != nil {
		return "reconcile"
	}
	if _, err := intent.In(run, root).PruneUnclaimedLandedBranches(); err != nil {
		return "prune"
	}
	return ""
}

// releaseLanded releases the landed assignment that request and path name. A refused
// release states its diagnostic on stderr and answers false.
func releaseLanded(j joins, a ambient, root, request, path string, stderr io.Writer) bool {
	var diagnostic bytes.Buffer
	if j.releaseLandingAssignment(j, a, root, []string{"--request", request, path}, io.Discard, &diagnostic) == 0 {
		return true
	}
	if diagnostic.Len() > 0 {
		fmt.Fprintln(stderr, sanitize.Controls(strings.TrimSuffix(diagnostic.String(), "\n")))
	}
	return false
}

// refreshBroker republishes root's own Bench executable after the landing published the
// sources it is built from, through freshness.Refresh. It applies only where the
// destination declares Bench build inputs, which is the predicate the landing's broker
// notice reads; a linked installation declares none and stays untouched. Refresh settles
// completion from the tree, so a resume needs no journal.
//
// An installed candidate never verifies against the destination's sources, so a rebuild
// would replace it with an executable that refuses its build's ledger. The refresh therefore
// reads the held-install predicate first, and a held refresh writes nothing.
func refreshBroker(j joins, run *git.Run, root string, stderr io.Writer) string {
	if !freshness.DeclaresBuildInputs(root) {
		return effectNotApplicable
	}
	if j.heldInstall(commitrepo.Store{Root: root, Run: run}) {
		return effectHeld
	}
	// The build runs the destination's own entry point, which is candidate code. It runs
	// here and nowhere else, after the release, so it can reach neither the gate, nor the
	// composition, nor the destination compare-and-swap that published these bytes. The
	// refresh never executes the published executable itself.
	_, buildErr, err := freshness.Refresh(context.Background(), root, j.buildSubject)
	if buildErr != nil {
		fmt.Fprintln(stderr, "landing refresh build failed: "+sanitize.Controls(buildErr.Error()))
	}
	if err != nil {
		fmt.Fprintln(stderr, "landing refresh failed: "+sanitize.Controls(err.Error()))
		return effectFailed
	}
	return effectComplete
}

// readHeldInstall is the held-install predicate at the store's root. It holds while a
// declared build is open and the broker manifest beside the root's wrapper binds the digest
// of the root's published executable. It also holds when the open builds do not read, so the
// refresh never rebuilds over a ledger that it cannot read.
func readHeldInstall(store commitrepo.Store) bool {
	_, declared, err := store.DeclaredBuild()
	if err != nil {
		return true
	}
	return declared && manifestBindsExecutable(store.Root)
}

// manifestBindsExecutable reports whether the broker manifest beside root's wrapper binds
// the digest of root's published executable.
func manifestBindsExecutable(root string) bool {
	fields, err := brokermanifest.Read(manifestPath(root))
	if err != nil {
		return false
	}
	digest, err := brokermanifest.Digest(freshness.PublishedExecutable(root))
	return err == nil && digest == fields[brokermanifest.DigestField]
}

// cleanLandedSiblings retires the sibling worktrees whose work this landing carries. The
// scope is the destination base the landing composed against, so the set holds only the
// assignments this landing landed and never the repository's other landed work.
//
// The landing plans and applies in one process, so it prints no fingerprint and asks for
// no second call. Completion is a predicate over the tree, never a recorded step: a
// resume re-plans the same narrowed set and finds it empty once the set is settled.
//
// A destination base the landing could not resolve names no scope, and a repository-wide
// removal is not what a landing was asked for. So there is nothing this landing owns to
// clean, and the effect completes without a mutation.
func cleanLandedSiblings(j joins, a ambient, root, base string, stderr io.Writer) string {
	if base == "" {
		return effectComplete
	}
	options := CleanupOptions{}
	set, err := planLandedSet(j, a.run, root, options, base)
	if err != nil {
		fmt.Fprintln(stderr, "landing cleanup plan failed: "+sanitize.Controls(err.Error()))
		return effectFailed
	}
	printLandedCleanupRows(stderr, set)
	if len(set.rows) == 0 {
		return effectComplete
	}
	if _, err := applyLandedSet(j, a, root, set, options, base); err != nil {
		fmt.Fprintln(stderr, "landing cleanup failed: "+sanitize.Controls(err.Error()))
		return effectFailed
	}
	return effectComplete
}

// printLandedCleanupRows states the narrowed plan beside the landing's other diagnostic
// evidence, so stdout keeps the effects row and the landed record alone.
//
// A worktree path this line cannot carry verbatim is replaced by its assignment pointer.
// The raw bytes would forge a line of their own, and an escaped path would name a tree
// that does not exist.
func printLandedCleanupRows(stderr io.Writer, set landedCleanupSet) {
	for _, row := range set.rows {
		reason := string(row.plan.ReasonCode)
		if reason == "" {
			reason = "none"
		}
		target := row.assignment.Worktree
		if !lineSafe(target) {
			target = "assignment/" + row.assignment.ID
		}
		fmt.Fprintf(stderr, "landing cleanup{assignment=%s,action=%s,reason=%s,target=%s}\n", row.assignment.ID, row.plan.Action, reason, target)
	}
}
