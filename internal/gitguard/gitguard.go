// Package gitguard is the destructive-git analyzer behind the PreToolUse guard
// (.bench/hooks/block-dangerous-git.sh) — the single source that both classifies an
// agent's Bash command (deny label or "") and enumerates the deny classes for the guard
// manifest, so enforcement and advertisement cannot drift. The shim pipes the PreToolUse
// envelope to `bench guard-git`, which classifies through this package.
//
// Threat model: an honest-mistake layer, not an
// evasion-resistant boundary. Wrapper scanning goes exactly one level deep; the
// backstops for a misaligned agent are the git pre-push hook and pooled-worktree
// isolation, not this analyzer.
package gitguard

import (
	"encoding/json"
	"strings"
)

// Checker resolves repository truth for the two verdicts that need it (checkout
// ref-ness and forced-creation clobber). It is injected so classification is unit
// testable without a real repo; the guard subcommand wires it to internal/git.
type Checker struct {
	RefResolves  func(string) bool
	BranchExists func(string) bool
}

// denyTable is the ordered source for every destructive class: classify returns from it
// (via denyLabels) and DescribeClasses enumerates it, so the block surface and its
// advertised manifest share one definition.
var denyTable = []struct{ key, label string }{
	{"push", "git push"},
	{"reset", "git reset --hard"},
	{"clean", "git clean -f"},
	{"branch-force", "git branch -f"},
	{"branch-delete", "git branch -D"},
	{"checkout", "git checkout path"},
	{"switch", "git switch --force"},
	{"restore", "git restore path"},
	{"rebase", "history rewrite"},
	{"stash", "git stash drop"},
	{"amend", "git commit --amend"},
	{"update-ref", "git update-ref -d"},
	{"tag", "git tag -d"},
	{"reflog", "git reflog expire"},
	{"worktree", "git worktree remove --force"},
}

var denyLabels = func() map[string]string {
	m := make(map[string]string, len(denyTable))
	for _, d := range denyTable {
		m[d.key] = d.label
	}
	return m
}()

// CommandFromEnvelope extracts tool_input.command from the PreToolUse JSON envelope on
// stdin, returning "" when the envelope is malformed or the field is absent — the
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
	return scan(tokenize(command), chk, true)
}

// BlockMessage is the actionable refusal the guard returns to the agent for a blocked
// command, naming the deny label. One source for the message the hook emits.
func BlockMessage(label string) string {
	return "BLOCKED: `" + label + "` — you don't have authority over this. " +
		"The merge and any history rewrite are the user's; a failed shift is rolled back " +
		"by bench, not by you. Stop and hand back."
}

// DescribeClasses is the comma-joined deny surface for the guard manifest
// (`--describe-classes`): the unique labels in table order.
func DescribeClasses() string {
	var seen []string
	for _, d := range denyTable {
		if !containsStr(seen, d.label) {
			seen = append(seen, d.label)
		}
	}
	return strings.Join(seen, ", ")
}

func containsStr(xs []string, s string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}
