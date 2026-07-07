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
