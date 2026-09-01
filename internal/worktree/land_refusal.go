// Landing terminal receipts and refusal rendering: completion lines, next-step pointers, and refusal exits.
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

// conflictRepairPrefix names the hand repair a composition conflict demands, up to the
// commit that records the resolution: merge the incoming commit into the worktree, then
// commit it. The landing and `bench worktree merge` refuse the same conflict, so both
// name this one repair, and the landing appends its own review-and-re-run tail. Neither
// verb composes the repair itself, so the merge step is raw Git and the value says so.
func conflictRepairPrefix(incoming, assignment, path string) string {
	merge := "git -C " + sanitize.ShellQuote(path) + " merge " + conflictCommitArg(incoming)
	if !lineSafe(path) {
		merge = "bench worktree exec " + assignment + " -- git merge " + conflictCommitArg(incoming)
	}
	return merge + " (bench worktree merge refuses this conflict; resolve it by hand); then bench commit"
}

// conflictCommitArg renders the commit a repair line names, and the placeholder that
// stands in when the value is not line-safe.
func conflictCommitArg(commit string) string {
	if lineSafe(commit) {
		return sanitize.ShellQuote(commit)
	}
	return "<full-destination-commit>"
}

// conflictContinuePrefix names the hand repair a source worktree that already holds a
// pending merge demands. The merge is open, so a second merge is not the repair: the
// operator resolves the conflicted paths and finishes the merge, which records the
// resolution itself and needs no separate commit.
func conflictContinuePrefix(assignment, path string) string {
	if !lineSafe(path) {
		return "bench worktree exec " + assignment + " -- git merge --continue (resolve the conflicted paths of the merge in progress first)"
	}
	return "git -C " + sanitize.ShellQuote(path) + " merge --continue (resolve the conflicted paths of the merge in progress first)"
}

// sourceMergePending reports whether the source worktree holds a pending merge, which
// MERGE_HEAD in the worktree's Git directory records. An unreadable Git directory leaves
// the state undecided and answers false, because the commit-and-review route is correct
// under both states while the continuation route is correct under one.
func sourceMergePending(source string) bool {
	dir, err := git.Output("-C", source, "rev-parse", "--absolute-git-dir")
	if err != nil || dir == "" {
		return false
	}
	if _, err := os.ReadFile(filepath.Join(dir, "MERGE_HEAD")); err != nil {
		return false
	}
	return true
}

// landingConflictRefusal is the one constructor the conflict face travels through. The
// route is not an argument, so no call site composes it: the constructor reads the source
// worktree's merge state and the route builder branches on the answer.
func landingConflictRefusal(conflict landing.ConflictError, destination, assignment, specArg, path, source string) refusalError {
	return refusalError{refusal{
		detail: conflict.Error(),
		paths:  conflict.Paths,
		next:   landingConflictNext(destination, assignment, specArg, path, sourceMergePending(source)),
	}}
}

// landingConflictNext names the source repair a conflict outside the rule table
// demands, in the order the operator runs it. A source worktree with no pending merge
// merges the destination in, commits the repair, reviews the new range, and re-runs the
// landing with the repaired tip. A source worktree that holds a pending merge finishes
// that merge instead, so its route names the continuation and no second merge.
func landingConflictNext(destination, assignment, specArg, path string, pending bool) string {
	// A spec-less landing re-runs spec-less, so it names no --spec at all rather than an
	// empty value the grammar refuses.
	specFlag := ""
	if specArg != "" {
		specFlag = " --spec <spec>"
		if lineSafe(specArg) {
			specFlag = " --spec " + sanitize.ShellQuote(specArg)
		}
	}
	rerun := atSourceWorktree("bench worktree land --request <request> --base "+conflictCommitArg(destination)+" --source-tip <repaired-source-tip>"+specFlag+" -m <message>", path, assignment)
	repair := conflictRepairPrefix(destination, assignment, path)
	if pending {
		repair = conflictContinuePrefix(assignment, path)
	}
	return repair + "; then /bench-review-implementation; then " + rerun
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

// The landing refusal face names. An operator reads a face's sentence in the record, so
// the name and the registry entry are the same fact.
const (
	faceDestinationNotClean = "destination-not-clean"
	faceDestinationResidue  = "destination-residue"
	faceSourceNotClean      = "source-not-clean"
	faceSourceNotFenced     = "source-not-fenced"
	faceSourceTipMismatch   = "source-tip-mismatch"
	// The resume path refuses on its own destination and marker state, so those two
	// faces are declared here beside the first run's.
	faceResumeDestinationResidue = "resume-destination-residue"
	faceResumeMarker             = "resume-marker"
)

// destinationCleanRepair is the repair a destination that carries uncommitted work
// demands. The first run and the resume refuse on the same destination state, so both
// faces name this one repair rather than each spelling it.
func destinationCleanRepair(rerun string) string {
	return "commit the destination's uncommitted work, or discard it; then " + rerun
}

// destinationResidueRepair is the repair a destination that carries undeclared ignored
// residue demands. The operator has two routes out of the state, so the line names both:
// the declaration file that adopts the paths, and the removal that discards them.
func destinationResidueRepair(rerun string, paths []string) string {
	return "declare the refusal_paths entries in .bench/build-outputs.json, or remove them from the landing checkout with " +
		residueRemovalCommand(paths) + "; then " + rerun
}

// residueRemovalCommand is the exact removal the residue repair names. Git owns no
// removal of an ignored path, so the command is plain `rm`, and the value it names is the
// destination-relative path the refusal's own table lists. A path the operator could not
// paste back takes the placeholder, which the table beside the route resolves.
func residueRemovalCommand(paths []string) string {
	arguments := make([]string, 0, len(paths))
	for _, path := range paths {
		if !lineSafe(path) {
			return "rm -rf <refusal_paths entries>"
		}
		arguments = append(arguments, sanitize.ShellQuote(path))
	}
	if len(arguments) == 0 {
		return "rm -rf <refusal_paths entries>"
	}
	return "rm -rf " + strings.Join(arguments, " ")
}

// landingRefusalFace is one refusal the landing's preflight prints. detail is the
// sentence, and repair composes the face's own repair ahead of the caller's re-run, from
// the refusal's own paths where the repair names them. The registry test walks the slice
// and drives one producing fixture per entry, so a face added without a fixture, or with
// an empty repair, turns the gate red.
type landingRefusalFace struct {
	name   string
	detail string
	repair func(rerun string, raised refusal) string
}

// route is the bare reading of a face's repair, which the proofs that pin a repair with no
// observed values of its own read.
func (f landingRefusalFace) route(rerun string) string { return f.repair(rerun, refusal{}) }

// pathless adapts a repair that reads none of the refusal's observed values.
func pathless(build func(rerun string) string) func(rerun string, raised refusal) string {
	return func(rerun string, _ refusal) string { return build(rerun) }
}

// withPaths adapts a repair that reads the refusal's own path table.
func withPaths(build func(rerun string, paths []string) string) func(rerun string, raised refusal) string {
	return func(rerun string, raised refusal) string { return build(rerun, raised.paths) }
}

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

// landingRefusalFaces is the declared registry. It is the authoritative inventory of the
// landing's refusal faces, so a later ticket adds its face here rather than composing a
// route at the proof that fails.
var landingRefusalFaces = []landingRefusalFace{
	{
		name:   faceDestinationNotClean,
		detail: "landing destination is not clean",
		repair: pathless(destinationCleanRepair),
	},
	{
		name:   faceDestinationResidue,
		detail: "landing destination has undeclared ignored residue",
		repair: withPaths(destinationResidueRepair),
	},
	{
		// The mismatch refusal already carries the tip the tree holds beside the one the
		// caller named, so the route is the caller's own command re-pointed at the tip that
		// works. The operator then has one exact next command for a moved tip.
		name:   faceSourceTipMismatch,
		detail: sourceTipMismatchDetail,
		repair: retargetSourceTip,
	},
	{
		name:   faceSourceNotClean,
		detail: "reviewed source is not clean",
		repair: pathless(func(rerun string) string {
			return "commit the reviewed source's uncommitted work, or discard it; then " + rerun
		}),
	},
	{
		name:   faceSourceNotFenced,
		detail: "reviewed source range or ownership fence is invalid",
		repair: pathless(func(rerun string) string {
			return "take the refusal_paths entries out of the reviewed range, or declare them under the spec's ## Ownership fences; then " + rerun
		}),
	},
	{
		// The residue policy owns the sentence this face prints, so the entry declares
		// none and the constructor carries the observed one.
		name:   faceResumeDestinationResidue,
		repair: pathless(destinationCleanRepair),
	},
	{
		name:   faceResumeMarker,
		detail: landingpolicy.MarkerRefusalDetail,
		repair: pathless(func(rerun string) string {
			return "run bench gate in the landing checkout to record its green marker; then " + rerun
		}),
	},
}

func landingRefusalFaceByName(name string) landingRefusalFace {
	for _, face := range landingRefusalFaces {
		if face.name == name {
			return face
		}
	}
	// A name outside the registry is a programming fault in this package, not an operator
	// condition. It surfaces as a refusal rather than a panic, because the landing must not
	// abort a caller's session over its own bookkeeping.
	return landingRefusalFace{
		name:   name,
		detail: "landing refusal face " + name + " is unregistered",
		repair: pathless(func(rerun string) string { return rerun }),
	}
}

// landingFaceByDetail finds the registered face a refusal's own sentence names. A face
// that declares no sentence of its own matches nothing here, because its sentence comes
// from the policy at the refusal rather than from the registry.
func landingFaceByDetail(detail string) (landingRefusalFace, bool) {
	for _, face := range landingRefusalFaces {
		if face.detail != "" && face.detail == detail {
			return face, true
		}
	}
	return landingRefusalFace{}, false
}

// landingFaceRefusal is the one constructor a registered face travels through. rerun is a
// required argument, so no site can print a face's repair without the caller's own re-run
// behind it. detail carries the observed sentence for a face whose sentence a policy
// owns, and a face that declares its own sentence prints that one. The `refusal` struct
// keeps its optional next field for the verbs outside this registry's reach.
func landingFaceRefusal(name, detail, rerun string, paths []string) refusalError {
	return landingFaceRefusalOf(name, refusal{detail: detail, paths: paths}, rerun)
}

// landingFaceRefusalOf renders one registered face over the refusal a proof already
// raised, so the identities that refusal observed reach the route the face composes.
func landingFaceRefusalOf(name string, raised refusal, rerun string) refusalError {
	face := landingRefusalFaceByName(name)
	if face.detail != "" {
		raised.detail = face.detail
	}
	raised.next = face.repair(rerun, raised)
	return refusalError{raised}
}

// laterProofsSkipped is the sentence a refusal from a short-circuited proof group carries.
// The preflight runs its proofs in groups, and a fault stops the later proofs of its own
// group, so the operator who repairs this one fault must expect another refusal from the
// same group. A refusal from a group that ran to its end carries no such sentence.
const laterProofsSkipped = "later proofs in this group did not run"

// skippedProofs states the short-circuit ahead of the route, so the route still ends with
// the caller's own re-run. A refusal that names no route gains no sentence, because the
// sentence qualifies a repair rather than standing as one.
func skippedProofs(raised refusalError, shortCircuited bool) refusalError {
	if !shortCircuited || raised.next == "" {
		return raised
	}
	raised.next = laterProofsSkipped + "; " + raised.next
	return raised
}

// landingFaceRoute attaches the caller's own re-run to a preflight refusal whose sentence
// names a registered face. The route reads the flag values the caller passed, and the
// preflight assembler is the one place that holds them, so the attachment happens there
// rather than at the proof that failed. shortCircuited states whether this fault stopped
// the later proofs of its own group, which the assembler knows and the proof does not. A
// refusal outside the registry travels unchanged unless it already carries a route, which
// the sentence then qualifies.
func landingFaceRoute(err error, rerun string, shortCircuited bool) error {
	raised := refusal{detail: err.Error()}
	var typed refusalError
	if errors.As(err, &typed) {
		raised = typed.refusal
	}
	if face, ok := landingFaceByDetail(raised.detail); ok {
		return skippedProofs(landingFaceRefusalOf(face.name, raised, rerun), shortCircuited)
	}
	if !shortCircuited || raised.next == "" {
		return err
	}
	return skippedProofs(refusalError{raised}, true)
}

// landingRerun is the caller's own re-run of the landing, with the flag values it passed.
// Every landing-preflight route ends with it, so a repair does not cost the operator its
// flags. An assignment that has not resolved yet has no id for the pointer form to
// address, so the re-run names the operator's own worktree path instead.
func landingRerun(request, base, tip, specArg, path, assignment string) string {
	command := "bench worktree land --request " + landingRerunArg(request, "<request>") +
		" --base " + landingRerunArg(base, "<full-review-base>") +
		landingSourceTipFlag(tip)
	if specArg != "" {
		command += " --spec " + landingRerunArg(specArg, "<spec>")
	}
	command += " -m <message>"
	if assignment != "" {
		return atSourceWorktree(command, path, assignment)
	}
	if lineSafe(path) {
		return command + " " + sanitize.ShellQuote(path)
	}
	return command + " <worktree-path>"
}

// landingSourceTipFlag is the one rendering of the re-run's --source-tip argument. The
// mismatch face swaps this exact text, so the composition and the swap read the same fact.
func landingSourceTipFlag(tip string) string {
	return " --source-tip " + landingRerunArg(tip, "<full-source-tip>")
}

// landingRerunArg renders one flag value the re-run repeats, and the placeholder that
// stands in for a value the operator could not paste back.
func landingRerunArg(value, placeholder string) string {
	if value == "" || !lineSafe(value) {
		return placeholder
	}
	return sanitize.ShellQuote(value)
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
