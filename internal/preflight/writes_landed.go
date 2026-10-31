package preflight

import (
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/tickets"
)

// landedDeletions maps each plain `Writes:` entry whose path landed work deleted to the
// full identity of its deleting commit. It considers a plain entry that the checkout does
// not hold, as exists reports. Only build mode fills the map, so a review keeps every such
// entry in writes-resolve. The gatherer reads no deletion range here: the writes-resolve
// and writes-landed precedence alone decides which answer an entry gets.
func landedDeletions(root, mode string, parsed []tickets.Ticket, exists map[string]bool) map[string]string {
	if mode != modeBuild {
		return nil
	}
	landed := map[string]string{}
	tip, resolved := "", false
	for _, ticket := range parsed {
		for _, entry := range ticket.Writes {
			write, err := tickets.ParseWrite(entry)
			if err != nil || write.Kind != tickets.WriteExisting || exists[entry] {
				continue
			}
			if !resolved {
				if tip, resolved = defaultBranchTip(root); !resolved {
					return landed
				}
			}
			if commit, ok := deletingCommit(root, tip, write.Path); ok {
				landed[entry] = commit
			}
		}
	}
	return landed
}

// defaultBranchTip resolves the default branch through git.ResolvedDefault and its tip
// through git.ResolveCommit. Either failure answers false.
func defaultBranchTip(root string) (string, bool) {
	branch, ok := git.ResolvedDefault(root)
	if !ok {
		return "", false
	}
	tip, err := git.ResolveCommit(root, branch)
	return tip, err == nil
}

// deletingCommit answers the newest first-parent commit of tip that deleted path, when tip
// lacks the path. The first-parent diff of each merge counts, so a landing merge and a
// merged side-branch deletion both answer the merge. A rename counts as a deletion of its
// old path. The literal pathspec keeps a glob character or a space from matching a
// sibling. A failed query or an empty answer is no deletion.
func deletingCommit(root, tip, path string) (string, bool) {
	pathspec := ":(literal)" + path
	listing, err := git.Output("-C", root, "ls-tree", tip, "--", pathspec)
	if err != nil || listing != "" {
		return "", false
	}
	commit, err := git.Output("-C", root, "log", "-1", "--format=%H", "--diff-filter=D", "--no-renames", "--first-parent", "-m", tip, "--", pathspec)
	commit, _, _ = strings.Cut(commit, "\n")
	return commit, err == nil && commit != ""
}
