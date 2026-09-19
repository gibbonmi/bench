// Landing re-run rendering: the caller's own landing command, and the one rendering of
// each flag a refusal route swaps.
package worktree

import "github.com/gibbonmi/bench/internal/sanitize"

// landingRerun is the caller's own re-run of the landing, with the flag values it passed.
// Every landing-preflight route ends with it, so a repair does not cost the operator its
// flags. An assignment that has not resolved yet has no id for the pointer form to
// address, so the re-run names the operator's own worktree path instead.
func landingRerun(request, base, tip, specArg, path, assignment string) string {
	command := "bench worktree land --request " + landingRerunArg(request, "<request>") +
		landingBaseFlag(base) +
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

// landingBaseFlag is the one rendering of the re-run's --base argument. The fence face
// swaps this exact text when the source folded a later default-branch commit.
func landingBaseFlag(base string) string {
	return " --base " + landingRerunArg(base, "<full-review-base>")
}

// landingRerunArg renders one flag value the re-run repeats, and the placeholder that
// stands in for a value the operator could not paste back.
func landingRerunArg(value, placeholder string) string {
	if value == "" || !lineSafe(value) {
		return placeholder
	}
	return sanitize.ShellQuote(value)
}
