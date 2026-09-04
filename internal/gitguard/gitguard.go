// Package gitguard is the destructive-git analyzer behind the PreToolUse guard
// (.bench/hooks/block-dangerous-git.sh). It is the single source that both classifies an
// agent's Bash command (deny label or "") and enumerates the deny classes for the guard
// manifest. This is so enforcement and advertisement cannot drift. The shim pipes the
// PreToolUse envelope to `bench guard-git`, which classifies through this package.
//
// Threat model: an honest-mistake layer, not an
// evasion-resistant boundary. Wrapper scanning goes exactly one level deep; the
// backstops for a misaligned agent are the git pre-push hook and pooled-worktree
// isolation, not this analyzer.
package gitguard

import (
	"encoding/json"

	"github.com/gibbonmi/bench/internal/shellcommand"
)

// Checker resolves repository truth for the verdicts that need it (checkout ref-ness,
// forced-creation clobber, and the push destination rule). It is injected so
// classification is unit testable without a real repo; the guard subcommand wires it to
// internal/git. Each of the three push facts reports a name and whether one exists; a
// nil field reports no name, so the push verdict fails closed on an unwired Checker.
type Checker struct {
	RefResolves     func(string) bool
	BranchExists    func(string) bool
	DefaultBranch   func() (string, bool)
	CheckedOut      func() (string, bool)
	BareDestination func() (string, bool)
}

func (c Checker) defaultBranch() (string, bool)   { return callFact(c.DefaultBranch) }
func (c Checker) checkedOut() (string, bool)      { return callFact(c.CheckedOut) }
func (c Checker) bareDestination() (string, bool) { return callFact(c.BareDestination) }

func callFact(fact func() (string, bool)) (string, bool) {
	if fact == nil {
		return "", false
	}
	return fact()
}

// denyTable is the ordered source for every destructive class; classification returns
// its labels in the live block verdict. The advice column is the one source of the
// sentence a refusal appends, and only a class whose fix the agent can type carries one.
var denyTable = []struct{ key, label, advice string }{
	{"push-default", "git push to the default branch", ""},
	{"push-force", "git push --force", ""},
	{"push-delete", "git push --delete", ""},
	{"push-all", "git push --all", ""},
	{"push-mirror", "git push --mirror", ""},
	{"push-tags", "git push --tags", ""},
	{"push-unresolved", "git push with an unresolved destination", "Name the remote and the branch: git push <remote> <branch>."},
	{"reset", "git reset --hard", ""},
	{"clean", "git clean -f", ""},
	{"branch-force", "git branch -f", ""},
	{"branch-delete-safe", "git branch -d", ""},
	{"branch-delete-force", "git branch -D", ""},
	{"checkout", "git checkout path", ""},
	{"switch", "git switch --force", ""},
	{"restore", "git restore path", ""},
	{"rebase", "history rewrite", ""},
	{"filter-branch", "git filter-branch", ""},
	{"amend", "git commit --amend", ""},
	{"update-ref", "git update-ref -d", ""},
	{"tag", "git tag -d", ""},
	{"reflog", "git reflog expire", ""},
	{"worktree", "git worktree remove --force", ""},
	{"stash-drop", "git stash drop", ""},
	{"stash-clear", "git stash clear", ""},
	{"rm-force", "git rm -rf", ""},
}

var denyLabels = func() map[string]string {
	m := make(map[string]string, len(denyTable))
	for _, d := range denyTable {
		m[d.key] = d.label
	}
	return m
}()

// denyAdvice indexes the table's advice column by label, so BlockMessage appends the
// sentence of the class it names without a second copy of it.
var denyAdvice = func() map[string]string {
	m := make(map[string]string, len(denyTable))
	for _, d := range denyTable {
		if d.advice != "" {
			m[d.label] = d.advice
		}
	}
	return m
}()

// CommandFromEnvelope extracts tool_input.command from the PreToolUse JSON envelope on
// stdin. It returns "" when the envelope is malformed or the field is absent. The
// analyzer then classifies "" as allow, so a non-JSON or command-less envelope is
// permitted exactly as the Python shim's `|| true` did.
func CommandFromEnvelope(data []byte) string {
	var e struct {
		ToolInput struct {
			Command string `json:"command"`
		} `json:"tool_input"`
	}
	if err := json.Unmarshal(data, &e); err != nil {
		return ""
	}
	return e.ToolInput.Command
}

// Classify returns the deny label for command, or "" to allow. It tokenizes then scans
// with wrapper recursion enabled; chk supplies ref/branch truth to the two verdicts
// that need it.
func Classify(command string, chk Checker) string {
	return scan(shellcommand.Parse(command), chk, true)
}

// BlockMessage is the actionable refusal the guard returns to the agent for a blocked
// command, naming the deny label and appending that label's advice sentence when the
// table carries one. One source for the message the hook emits.
func BlockMessage(label string) string {
	msg := "BLOCKED: `" + label + "` — you don't have authority over this. " +
		"The merge and any history rewrite are the user's, and discarding work leaves a " +
		"gate verdict answering for something other than what is " +
		"on disk. A failed shift is rolled back by bench, not by you. Stop and hand back."
	if advice := denyAdvice[label]; advice != "" {
		msg += " " + advice
	}
	return msg
}
