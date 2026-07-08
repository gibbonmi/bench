package worktree

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
)

func CleanCommand(args []string, stdout, stderr io.Writer) int {
	return cleanCommand(args, stdout, stderr)
}

func cleanCommand(args []string, stdout, stderr io.Writer) int {
	if len(args) != 0 {
		fmt.Fprint(stderr, WorktreeUsage())
		return 2
	}
	root, err := git.Root()
	if err != nil {
		fmt.Fprintln(stderr, toon.NotInRepo())
		return 1
	}
	// Two independent phases run after the in-repo guard: the branch sweep and the
	// out-of-pool worktree removal. The command's exit is the higher severity of the
	// two, so a swept branch never masks a refused worktree and vice versa.
	sweepExit := sweepDelegateBranches(root, stdout, stderr)
	worktreeExit := cleanOutOfPoolWorktrees(root, stdout, stderr)
	if sweepExit > worktreeExit {
		return sweepExit
	}
	return worktreeExit
}

// sweepDelegateBranches deletes every landed worktree-* scratch orphan that no live worktree
// holds and reports each on stdout; an orphan carrying unique commits is left in place with a
// hand-inspect note, since a scratch name alone is not proof its work landed. Landedness is
// LandedInDefault's proof against the default branch (ancestry or patch containment), never a
// name compare, and the delete is forced so it does not depend on the repo root's HEAD — git's
// own merged-check for a plain delete is HEAD-relative and would refuse a branch merged into
// the default branch but not into HEAD.
//
// Before any orphan is classified the resolved default branch must resolve to a commit; when it does
// not, the sweep refuses loudly on stderr, deletes nothing, and returns 1 — the false-empty guard,
// so an unresolvable default never yields a silent all-clean sweep. The guard is reached only once at
// least one orphan exists: with no orphan there is no mergedness to compute and so no clean report to
// falsify, and gating it there keeps `bench worktree clean` usable for worktree removal in a repo
// whose default branch happens not to resolve. Returns 0 otherwise, whether it deleted, kept, or
// found no orphans.
func sweepDelegateBranches(root string, stdout, stderr io.Writer) int {
	orphans := OrphanedDelegateBranches(root)
	if len(orphans) == 0 {
		return 0
	}
	def, ok := ResolvedDefault(root)
	if !ok {
		fmt.Fprintf(stderr, "error: bench worktree clean cannot resolve the default branch (%s) to a commit; deleting no branches\n", def)
		return 1
	}
	for _, branch := range orphans {
		landed, byContent := LandedInDefault(root, branch, def)
		if !landed {
			fmt.Fprintf(stdout, "kept branch %s (unique commits — inspect or delete by hand)\n", branch)
			continue
		}
		if out, err := exec.Command("git", "-C", root, "branch", "-D", branch).CombinedOutput(); err != nil {
			fmt.Fprintf(stderr, "error: could not delete branch %s: %s\n", branch, strings.TrimSpace(string(out)))
			continue
		}
		if byContent {
			fmt.Fprintf(stdout, "deleted branch %s (landed by content)\n", branch)
			continue
		}
		fmt.Fprintf(stdout, "deleted branch %s\n", branch)
	}
	return 0
}

// cleanOutOfPoolWorktrees removes every out-of-pool worktree without asking: committed work
// lives on the worktree's branch, so removing the checkout destroys nothing git cannot
// recover or revert. A dirty worktree is salvaged first — its uncommitted changes are
// committed onto its own branch — and then removed; the branch survives under the sweep's
// merged/unmerged rule. Only a dirty *detached* worktree is refused: there is no branch to
// hold the salvage, so deleting it would genuinely lose the changes.
func cleanOutOfPoolWorktrees(root string, stdout, stderr io.Writer) int {
	var refused []string
	removed := 0
	for _, wt := range RegisteredWorktrees(root) {
		if wt.Class != ClassOutOfPool {
			continue
		}
		if _, err := os.Stat(wt.Path); os.IsNotExist(err) {
			continue
		}
		if !isClean(wt.Path) {
			branch, err := git.Output("-C", wt.Path, "symbolic-ref", "--quiet", "--short", "HEAD")
			if err != nil || strings.TrimSpace(branch) == "" {
				refused = append(refused, wt.Path)
				continue
			}
			if out, err := exec.Command("git", "-C", wt.Path, "add", "-A").CombinedOutput(); err != nil {
				fmt.Fprintf(stderr, "error: could not salvage %s: %s\n", wt.Path, strings.TrimSpace(string(out)))
				refused = append(refused, wt.Path)
				continue
			}
			// A fixed committer identity: the salvage is machine-made and must not fail in a
			// worktree whose repo never configured user.name/user.email.
			if out, err := exec.Command("git", "-C", wt.Path, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-q", "-m", "wip: salvaged by bench worktree clean").CombinedOutput(); err != nil {
				fmt.Fprintf(stderr, "error: could not salvage %s: %s\n", wt.Path, strings.TrimSpace(string(out)))
				refused = append(refused, wt.Path)
				continue
			}
			fmt.Fprintf(stdout, "salvaged uncommitted changes in %s onto %s\n", wt.Path, strings.TrimSpace(branch))
		}
		if out, err := exec.Command("git", "-C", root, "worktree", "remove", wt.Path).CombinedOutput(); err != nil {
			refused = append(refused, wt.Path)
			if len(out) > 0 {
				fmt.Fprint(stderr, string(out))
			}
			continue
		}
		removed++
		fmt.Fprintf(stdout, "removed %s\n", wt.Path)
	}
	_ = exec.Command("git", "-C", root, "worktree", "prune").Run()
	if len(refused) > 0 {
		printRefused(stdout, refused)
		return 1
	}
	if removed == 0 {
		fmt.Fprintln(stdout, "bench worktree clean: nothing to clean")
	}
	return 0
}

func printRefused(stdout io.Writer, paths []string) {
	fmt.Fprintln(stdout, "refused:")
	for _, path := range paths {
		fmt.Fprintf(stdout, "  %s\n", path)
	}
}

func WorktreeUsage() string {
	return "usage: bench worktree\n       bench worktree clean\n"
}
