// Landing re-run rendering: the caller's own landing command, and the one rendering of
// each flag a refusal route swaps.
package worktree

import "github.com/gibbonmi/bench/internal/refusalroute"

// repairedSourceTipFlag is the re-run's --source-tip argument after a repair that commits
// in the source. The commit moves the tip the caller named, so the operator fills the tip
// that its repair committed.
const repairedSourceTipFlag = " --source-tip <repaired-source-tip>"

// landingRerun is the caller's own re-run of the landing, with the flag values it passed.
// Every landing-preflight route ends with it, so a repair does not cost the operator its
// flags.
func landingRerun(request, base, tip, specArg, path, assignment string) string {
	return landingRerunAt(request, base, landingSourceTipFlag(tip), specArg, path, assignment)
}

// landingRerunAt is landingRerun with its --source-tip argument already rendered, so a
// route whose repair moves the tip can name the repaired one.
func landingRerunAt(request, base, tipFlag, specArg, path, assignment string) string {
	command := "bench worktree land --request " + refusalroute.Arg("request", request) +
		landingBaseFlag(base) +
		tipFlag
	if specArg != "" {
		command += " --spec " + refusalroute.Arg("spec", specArg)
	}
	return atSourceWorktree(command+" -m <message>", path, assignment)
}

// landingSourceTipFlag is the one rendering of the re-run's --source-tip argument. The
// mismatch face swaps this exact text, so the composition and the swap read the same fact.
func landingSourceTipFlag(tip string) string {
	return " --source-tip " + refusalroute.Arg("full-source-tip", tip)
}

// landingBaseFlag is the one rendering of the re-run's --base argument, so a proof can
// find the caller's base in a route.
func landingBaseFlag(base string) string {
	return " --base " + refusalroute.Arg("full-review-base", base)
}
