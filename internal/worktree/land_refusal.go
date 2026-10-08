// Landing terminal receipts and refusal rendering: completion lines, next-step pointers,
// the refused record the worktree verbs print, and refusal exits.
package worktree

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/spec"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/worktree/landingpolicy"
)

func landedIncomplete(stdout io.Writer, result landing.ReviewedResult, specArg, path, assignment, step string, records int) int {
	resume := landingResumeNext(result, specArg, path, assignment)
	next := refusalroute.New(faceLandIncomplete, refusalroute.Facts{Values: map[string]string{refusalroute.FactRerun: resume}}).Route
	outcome := landingpolicy.Terminal(landingpolicy.TerminalFacts{FailedStep: step, Active: true})
	fmt.Fprintf(stdout, "landed{source_base=%s,source_tip=%s,destination_base=%s,published_commit=%s,tree=%s,worktree=%s,%s=%s,census=%d}\n", result.SourceBase, result.SourceTip, result.DestinationBase, result.Commit, result.Tree, outcome.WorktreeState, refusalroute.NextField, sanitize.Controls(next), records)
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
			return atSourceWorktree("bench worktree land --resume <full-published-commit> --request <request> --base <full-review-base> --source-tip <full-source-tip> --spec <spec>", path, assignment)
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

// conflictFacts are the facts the hand merge of a composition conflict reads: the commit
// to merge and the Git command addressed at the conflicted worktree. The landing and
// `bench worktree merge` refuse the same conflict, so both read these. A path that is not
// line-safe takes the assignment pointer form, because no quoting makes a control byte
// pasteable.
func conflictFacts(commit, assignment, path string) map[string]string {
	checkoutGit := "git -C " + sanitize.ShellQuote(path)
	if !lineSafe(path) {
		checkoutGit = "bench worktree exec " + assignment + " -- git"
	}
	return map[string]string{refusalroute.FactCheckoutGit: checkoutGit, refusalroute.FactConflictCommit: commit}
}

// sourceMergePending reports whether the source worktree holds a pending merge, which
// MERGE_HEAD in the worktree's Git directory records. An unreadable Git directory leaves
// the state undecided and answers false, because the commit-and-review route is correct
// under both states while the continuation route is correct under one.
func sourceMergePending(source string) bool {
	dir, err := git.AdminDir(source)
	if err != nil {
		return false
	}
	if _, err := os.ReadFile(filepath.Join(dir, "MERGE_HEAD")); err != nil {
		return false
	}
	return true
}

// landingConflictRefusal is the one constructor the conflict faces travel through. It
// reads the source worktree's merge state, which chooses the face, and composes the re-run
// the repaired source takes. The registered face puts the review step and that re-run
// behind the hand repair, so no call site composes a route.
func landingConflictRefusal(conflict landing.ConflictError, destination, assignment, specArg, path, source string) refusalError {
	face := faceCompositionConflict
	if sourceMergePending(source) {
		face = faceCompositionConflictPending
	}
	// The conflict route never prints the caller's request, so the re-run names its
	// placeholder.
	rerun := landingRerunAt("", destination, repairedSourceTipFlag, specArg, path, assignment)
	raised := refusal{detail: conflict.Error(), paths: conflict.Paths, values: conflictFacts(destination, assignment, path)}
	return landingFaceRefusal(face, raised, rerun, "")
}

// atSourceWorktree addresses the source worktree a command's trailing positional
// names. A path that is not line-safe takes the registry's checkout form: the route looks
// the path up by the assignment id, which an assignment that has not resolved yet leaves
// empty, and the command names the checkout placeholder.
func atSourceWorktree(command, path, assignment string) string {
	if lineSafe(path) {
		return command + " " + sanitize.ShellQuote(path)
	}
	return refusalroute.AtCheckout(command, assignment)
}

// The landing refusal face names. The shared registry declares each face, and a raising
// site names the face it raises.
const (
	faceDestinationNotClean        = "destination-not-clean"
	faceDestinationCollision       = "destination-collision"
	faceSourceNotClean             = "source-not-clean"
	faceSourceNotFenced            = "source-not-fenced"
	faceSourceTipMismatch          = "source-tip-mismatch"
	faceCompositionConflict        = "composition-conflict"
	faceCompositionConflictPending = "composition-conflict-pending"
	faceResumeDestinationResidue   = "resume-destination-residue"
	faceResumeMarker               = "resume-marker"
	faceLandIncomplete             = "land-incomplete"
	faceLandHandback               = "land-handback"
	// The landing raises these two from the composed tree's authorization refusal, by the
	// kind the gate attributed.
	faceLandRed            = "land-red"
	faceLandInfrastructure = "land-infrastructure"
)

// retargetSourceTip re-points the caller's own re-run at the source tip the landing read
// in the tree, so a moved tip leaves exactly one command to run. The refusal carries the
// requested tip and the read one already, and landingSourceTipFlag is the single rendering
// of that flag, so the route comes from an exact swap and not from a re-parse of the
// command. A refusal that names neither tip keeps the caller's command unchanged.
func retargetSourceTip(rerun string, raised refusal) string {
	if raised.observed == "" || raised.wanted == "" {
		return rerun
	}
	return strings.Replace(rerun, landingSourceTipFlag(raised.observed), landingSourceTipFlag(raised.wanted), 1)
}

// landingFaceOf names the registered face a refusal raised. A proof that raised a face
// names it on the refusal. A proof that reads a cause outside this package states the
// face's own sentence and then the cause, so a sentence that opens with a face's declared
// one names that face too. A face whose sentence a policy owns matches no sentence here.
func landingFaceOf(raised refusal) (string, bool) {
	if raised.face != "" {
		return raised.face, true
	}
	for _, face := range refusalroute.Faces(refusalroute.Land) {
		if face.Sentence != "" && (face.Sentence == raised.detail || strings.HasPrefix(raised.detail, face.Sentence+": ")) {
			return face.Name, true
		}
	}
	return "", false
}

// landingFaceRefusal is the one constructor a landing face travels through. raised carries
// what the proof observed: the sentence for a face whose sentence a policy owns, the paths,
// and the values the face's route reads. rerun is a required argument, so no site prints a
// face's repair without the caller's own re-run behind it. preface states a qualifier
// ahead of the route's first step. A refusal that already states a sentence keeps it,
// because that sentence carries the cause the face's declared one drops.
func landingFaceRefusal(name string, raised refusal, rerun, preface string) refusalError {
	if name == faceSourceTipMismatch {
		rerun = retargetSourceTip(rerun, raised)
	}
	values := map[string]string{refusalroute.FactRerun: rerun}
	for slot, value := range raised.values {
		values[slot] = value
	}
	built := refusalroute.New(name, refusalroute.Facts{Sentence: raised.detail, Paths: raised.paths, Preface: preface, Values: values})
	raised.face, raised.detail, raised.next = name, built.Sentence, built.Route
	return refusalError{raised}
}

// laterProofsSkipped is the sentence a refusal from a short-circuited proof group carries.
// The preflight runs its proofs in groups, and a fault stops the later proofs of its own
// group, so the operator who repairs this one fault must expect another refusal from the
// same group. A refusal from a group that ran to its end carries no such sentence.
const laterProofsSkipped = "later proofs in this group did not run"

// landingFaceRoute attaches the caller's own re-run to a landing refusal. The route reads
// the flag values the caller passed, and the verb's assembler is the one place that holds
// them, so the attachment happens there rather than at the proof that failed.
// shortCircuited states whether this fault stopped the later proofs of its own group,
// which the assembler knows and the proof does not; the face states it ahead of the route,
// so the route still ends with the re-run.
//
// A refusal that names a registered face takes that face. A refusal that carries a route
// of its own keeps it, and the skipped-proof sentence qualifies that route. Every other
// cause, an identity component with no recovery command included, has no route of its
// own, so it hands back to the reviewer under its own sentence.
func landingFaceRoute(err error, rerun string, shortCircuited bool) error {
	raised := raisedRefusal(err)
	preface := ""
	if shortCircuited {
		preface = laterProofsSkipped
	}
	if name, ok := landingFaceOf(raised); ok {
		return landingFaceRefusal(name, raised, rerun, preface)
	}
	if raised.next == "" {
		return landingFaceRefusal(faceLandHandback, raised, rerun, preface)
	}
	if !shortCircuited {
		return err
	}
	raised.next = laterProofsSkipped + "; " + raised.next
	return refusalError{raised}
}

// raisedRefusal is the refused record an error carries: the typed refusal, the face that a
// policy raised with the values it observed, or a refusal whose sentence is the error's own.
func raisedRefusal(err error) refusal {
	var typed refusalError
	if errors.As(err, &typed) {
		return typed.refusal
	}
	var raised refusalroute.Raised
	if errors.As(err, &raised) {
		return refusal{detail: err.Error(), face: raised.Name, values: raised.Values}
	}
	return refusal{detail: err.Error()}
}

// landingSourceRoute attaches the caller's own re-run to a first-run source proof's
// refusal. The repair of a source that is not clean or not fenced commits in the source,
// which moves the tip the caller named, so that re-run names the repaired tip.
func landingSourceRoute(err error, request, base, tip, specArg, path, assignment string) error {
	tipFlag := landingSourceTipFlag(tip)
	if name, _ := landingFaceOf(raisedRefusal(err)); name == faceSourceNotClean || name == faceSourceNotFenced {
		tipFlag = repairedSourceTipFlag
	}
	return landingFaceRoute(err, landingRerunAt(request, base, tipFlag, specArg, path, assignment), false)
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

// refusal is one refused record. face names the registered face the refusing proof raised,
// and values holds the facts that proof observed for the face's route, such as the label
// of the assignment that owns the refusing tree. A refusal outside the registry leaves
// both empty.
type refusal struct {
	detail, observed, wanted, next string
	paths                          []string
	face                           string
	values                         map[string]string
}
type refusalError struct{ refusal }

func (r refusal) fields() string {
	fields := []string{"detail=" + sanitize.Controls(r.detail)}
	for _, pair := range [][2]string{{"observed", r.observed}, {"wanted", r.wanted}, {refusalroute.NextField, r.next}} {
		if pair[1] != "" {
			fields = append(fields, pair[0]+"="+sanitize.Controls(pair[1]))
		}
	}
	return strings.Join(fields, ",")
}
func (e refusalError) Error() string {
	text := sanitize.Controls(e.detail)
	if fields := strings.TrimPrefix(e.fields(), "detail="+sanitize.Controls(e.detail)); fields != "" {
		text += "; " + strings.TrimPrefix(fields, ",")
	}
	if table := e.table(); table != "" {
		text += "\n" + table
	}
	return text
}
func (r refusal) table() string {
	if len(r.paths) == 0 {
		return ""
	}
	shown := len(r.paths)
	if shown > ignoredEntryLimit {
		shown = ignoredEntryLimit
	}
	rows := make([][]string, 0, shown)
	for _, path := range r.paths[:shown] {
		rows = append(rows, []string{sanitize.Controls(path)})
	}
	out, err := toon.Table(refusalPathsTable, []string{"path"}, rows)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("paths_total=%d\n%s", len(r.paths), out)
}
