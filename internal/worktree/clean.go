package worktree

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/terminal"
	"github.com/gibbonmi/bench/internal/toon"
)

const cleanConfirm = "clean worktrees"

func CleanCommand(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	return cleanCommand(args, stdin, stdout, stderr, terminal.IsTerminal)
}

func cleanCommand(args []string, stdin io.Reader, stdout, stderr io.Writer, isTerminal func(io.Reader) bool) int {
	if len(args) != 0 {
		fmt.Fprint(stderr, WorktreeUsage())
		return 2
	}
	root, err := git.Root()
	if err != nil {
		fmt.Fprintln(stderr, toon.NotInRepo())
		return 1
	}
	// Two independent phases run after the in-repo guard: the branch sweep (no TTY, no
	// confirmation) and the out-of-pool worktree removal. The command's exit is the higher
	// severity of the two, so a swept branch never masks a refused worktree and vice versa.
	sweepExit := sweepDelegateBranches(root, stdout, stderr)
	worktreeExit := cleanOutOfPoolWorktrees(root, stdin, stdout, stderr, isTerminal)
	if sweepExit > worktreeExit {
		return sweepExit
	}
	return worktreeExit
}

// sweepDelegateBranches deletes every fully-merged worktree-* scratch orphan that no live worktree
// holds and reports each on stdout; an orphan carrying unique commits is left in place with a
// hand-inspect note, since a scratch name alone is not proof its work landed. Mergedness is a
// merge-base ancestry test against the default branch, never a name compare, and the delete is
// forced so it does not depend on the repo root's HEAD — git's own merged-check for a plain delete
// is HEAD-relative and would refuse a branch merged into the default branch but not into HEAD.
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
	def := git.DefaultBranch(root)
	if !git.OK("-C", root, "rev-parse", "--verify", "--quiet", def+"^{commit}") {
		fmt.Fprintf(stderr, "error: bench worktree clean cannot resolve the default branch (%s) to a commit; deleting no branches\n", def)
		return 1
	}
	for _, branch := range orphans {
		if !git.OK("-C", root, "merge-base", "--is-ancestor", branch, def) {
			fmt.Fprintf(stdout, "kept branch %s (unique commits — inspect or delete by hand)\n", branch)
			continue
		}
		if out, err := exec.Command("git", "-C", root, "branch", "-D", branch).CombinedOutput(); err != nil {
			fmt.Fprintf(stderr, "error: could not delete branch %s: %s\n", branch, strings.TrimSpace(string(out)))
			continue
		}
		fmt.Fprintf(stdout, "deleted branch %s\n", branch)
	}
	return 0
}

func cleanOutOfPoolWorktrees(root string, stdin io.Reader, stdout, stderr io.Writer, isTerminal func(io.Reader) bool) int {
	var candidates, refused []string
	for _, wt := range RegisteredWorktrees(root) {
		if wt.Class != ClassOutOfPool {
			continue
		}
		if _, err := os.Stat(wt.Path); os.IsNotExist(err) {
			continue
		}
		if !isClean(wt.Path) {
			refused = append(refused, wt.Path)
			continue
		}
		candidates = append(candidates, wt.Path)
	}
	if len(candidates) == 0 {
		_ = exec.Command("git", "-C", root, "worktree", "prune").Run()
		if len(refused) == 0 {
			fmt.Fprintln(stdout, "bench worktree clean: nothing to clean")
			return 0
		}
		printRefused(stdout, refused)
		return 1
	}
	if !isTerminal(stdin) {
		fmt.Fprintln(stderr, "error: bench worktree clean requires an interactive TTY")
		return 1
	}
	fmt.Fprintln(stdout, "bench worktree clean will remove:")
	for _, path := range candidates {
		fmt.Fprintf(stdout, "  %s\n", path)
	}
	if len(refused) > 0 {
		printRefused(stdout, refused)
	}
	fmt.Fprintf(stdout, "Type '%s' to remove clean out-of-pool worktrees: ", cleanConfirm)
	line, err := bufio.NewReader(stdin).ReadString('\n')
	if err != nil && err != io.EOF {
		fmt.Fprintln(stderr, "error: could not read confirmation")
		return 1
	}
	if strings.TrimSpace(line) != cleanConfirm {
		fmt.Fprintln(stderr, "bench worktree clean declined; no worktrees removed")
		return 1
	}
	removed := 0
	for _, path := range candidates {
		cmd := exec.Command("git", "-C", root, "worktree", "remove", path)
		if out, err := cmd.CombinedOutput(); err != nil {
			refused = append(refused, path)
			if len(out) > 0 {
				fmt.Fprint(stderr, string(out))
			}
			continue
		}
		removed++
		fmt.Fprintf(stdout, "removed %s\n", path)
	}
	_ = exec.Command("git", "-C", root, "worktree", "prune").Run()
	if len(refused) > 0 {
		printRefused(stdout, refused)
	}
	if removed == 0 && len(refused) == 0 {
		fmt.Fprintln(stdout, "bench worktree clean: nothing to clean")
	}
	if len(refused) > 0 {
		return 1
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
