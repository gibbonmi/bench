// Landing terminal receipts and refusal rendering: completion lines, next-step pointers, and refusal exits.
package worktree

import (
	"errors"
	"fmt"
	"io"

	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/worktree/landingpolicy"
)

func landedIncomplete(stdout io.Writer, result landing.ReviewedResult, specArg, path, assignment, step string, records int) int {
	next := landingResumeNext(result, specArg, path, assignment)
	outcome := landingpolicy.Terminal(landingpolicy.TerminalFacts{FailedStep: step, Active: true})
	fmt.Fprintf(stdout, "landed{source_base=%s,source_tip=%s,destination_base=%s,published_commit=%s,tree=%s,worktree=%s,next=%s,census=%d}\n", result.SourceBase, result.SourceTip, result.DestinationBase, result.Commit, result.Tree, outcome.WorktreeState, sanitize.Controls(next), records)
	return outcome.ExitCode
}

// landedComplete renders the terminal landed record for a landing whose
// follow-up steps all completed, in this run (active) or a prior one.
func landedComplete(stdout io.Writer, result landing.ReviewedResult, active bool, records int) int {
	outcome := landingpolicy.Terminal(landingpolicy.TerminalFacts{Active: active})
	fmt.Fprintf(stdout, "landed{source_base=%s,source_tip=%s,destination_base=%s,published_commit=%s,tree=%s,worktree=%s,census=%d}\n", result.SourceBase, result.SourceTip, result.DestinationBase, result.Commit, result.Tree, outcome.WorktreeState, records)
	return outcome.ExitCode
}

// landingSlug is the spec slug a landing argument names, and the empty slug for the
// spec-less landing that named no argument.
func landingSlug(arg string) string {
	if arg == "" {
		return ""
	}
	return spec.LiveSpecSlug(arg)
}

func landingResumeNext(result landing.ReviewedResult, specArg, path, assignment string) string {
	values := []string{result.Commit, result.SourceBase, result.SourceTip, specArg}
	for _, value := range values {
		if !lineSafe(value) {
			return "bench worktree exec " + assignment + " -- bench worktree land --resume <full-published-commit> --request <request> --base <full-review-base> --source-tip <full-source-tip> --spec <spec> ."
		}
	}
	// A spec-less landing resumes spec-less, so the resume command it names carries no
	// --spec at all rather than an empty value the grammar refuses.
	specFlag := ""
	if specArg != "" {
		specFlag = " --spec " + sanitize.ShellQuote(specArg)
	}
	command := "bench worktree land --resume " + sanitize.ShellQuote(result.Commit) + " --request <request> --base " + sanitize.ShellQuote(result.SourceBase) + " --source-tip " + sanitize.ShellQuote(result.SourceTip) + specFlag
	return atSourceWorktree(command, path, assignment)
}

// landingConflictNext names the source repair a conflict outside the rule table
// demands, in the order the operator runs it: merge the destination into the source
// worktree, commit the repair, review the new range, and re-run the landing with the
// repaired tip. No Bench verb moves a retained worktree onto the destination yet, so
// the merge step is raw Git and the value says so.
func landingConflictNext(destination, assignment, specArg, path string) string {
	destinationArg := "<full-destination-commit>"
	if lineSafe(destination) {
		destinationArg = sanitize.ShellQuote(destination)
	}
	// A spec-less landing re-runs spec-less, so it names no --spec at all rather than an
	// empty value the grammar refuses.
	specFlag := ""
	if specArg != "" {
		specFlag = " --spec <spec>"
		if lineSafe(specArg) {
			specFlag = " --spec " + sanitize.ShellQuote(specArg)
		}
	}
	merge := "git -C " + sanitize.ShellQuote(path) + " merge " + destinationArg
	if !lineSafe(path) {
		merge = "bench worktree exec " + assignment + " -- git merge " + destinationArg
	}
	rerun := atSourceWorktree("bench worktree land --request <request> --base "+destinationArg+" --source-tip <repaired-source-tip>"+specFlag+" -m <message>", path, assignment)
	return merge + " (no Bench verb moves a retained worktree onto the destination yet); then bench commit; then /bench-review-implementation; then " + rerun
}

// atSourceWorktree addresses the source worktree a command's trailing positional
// names. A path that is not line-safe takes the pointer form every next= uses: the
// assignment id addresses the worktree that the unpasteable path cannot.
func atSourceWorktree(command, path, assignment string) string {
	if lineSafe(path) {
		return command + " " + sanitize.ShellQuote(path)
	}
	return "bench worktree exec " + assignment + " -- " + command + " ."
}

func landRefusal(stdout io.Writer, detail string) int {
	fmt.Fprintln(stdout, "refused{detail="+sanitize.Controls(detail)+"}")
	return 1
}

func landRefusalError(stdout io.Writer, err error) int {
	var typed refusalError
	if !errors.As(err, &typed) {
		return landRefusal(stdout, err.Error())
	}
	fmt.Fprintln(stdout, "refused{"+typed.fields()+"}")
	fmt.Fprint(stdout, typed.table())
	return 1
}
