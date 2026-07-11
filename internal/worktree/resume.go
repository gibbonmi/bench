package worktree

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/toon"
)

type ResumeResult struct {
	Worktrees, Branches   int
	Dirty, Locked, Leased int
	OpenIntent            int
}

type classifiedCandidate struct {
	registered Registered
	clean      bool
}

// ConservativeCleanup removes only clean, unlocked, out-of-pool worktrees and
// orphan scratch refs proven landed. It never stages, commits, or forces a worktree.
func ConservativeCleanup(root string) (ResumeResult, error) {
	registered, err := ClassifyRegisteredWorktrees(root)
	if err != nil {
		return ResumeResult{}, fmt.Errorf("git worktree list failed: %w", err)
	}
	result := ResumeResult{}
	var candidates []classifiedCandidate
	for _, wt := range registered {
		switch wt.Class {
		case ClassPoolLease:
			result.Leased++
		case ClassOutOfPool:
			if wt.Locked {
				result.Locked++
				continue
			}
			if _, err := os.Stat(wt.Path); os.IsNotExist(err) {
				continue
			}
			status, err := git.Raw("-C", wt.Path, "status", "--porcelain=v1", "-z", "--no-renames", "--ignored=matching")
			if err != nil {
				return ResumeResult{}, fmt.Errorf("classify worktree %s: %w", wt.Path, err)
			}
			clean := len(status) == 0
			if !clean {
				result.Dirty++
			}
			candidates = append(candidates, classifiedCandidate{registered: wt, clean: clean})
		}
	}
	for _, candidate := range candidates {
		if !candidate.clean {
			continue
		}
		if out, err := exec.Command("git", "-C", root, "worktree", "remove", candidate.registered.Path).CombinedOutput(); err != nil {
			return result, fmt.Errorf("remove worktree %s: %s", candidate.registered.Path, strings.TrimSpace(string(out)))
		}
		result.Worktrees++
	}
	_ = exec.Command("git", "-C", root, "worktree", "prune").Run()

	orchards, err := OrphanedDelegateBranches(root)
	if err != nil {
		return result, fmt.Errorf("classify orphan branches: %w", err)
	}
	if len(orchards) > 0 {
		def, ok := git.ResolvedDefault(root)
		if !ok {
			return result, fmt.Errorf("cannot resolve default branch; deleting no unproven branches")
		}
		for _, branch := range orchards {
			landed, _, err := git.LandedInDefault(root, branch, def)
			if err != nil {
				return result, fmt.Errorf("classify landed branch %s: %w", branch, err)
			}
			if !landed {
				continue
			}
			if out, err := exec.Command("git", "-C", root, "branch", "-D", branch).CombinedOutput(); err != nil {
				return result, fmt.Errorf("delete landed branch %s: %s", branch, strings.TrimSpace(string(out)))
			}
			result.Branches++
		}
	}
	return result, nil
}

func ResumeCleanCommand(args []string, stdout, stderr io.Writer) int {
	if len(args) != 0 {
		fmt.Fprintln(stderr, "usage: bench resume-clean")
		return 2
	}
	root, err := git.Root()
	if err != nil {
		fmt.Fprintln(stderr, toon.NotInRepo())
		return 1
	}
	result, err := ConservativeCleanup(root)
	if err != nil {
		fmt.Fprintf(stderr, "bench resume-clean: %v\n", err)
		return 1
	}
	if err := intent.Compact(root); err != nil {
		fmt.Fprintf(stderr, "bench resume-clean: intent refresh: %v\n", err)
		return 1
	}
	live, err := intent.Snapshot(root)
	if err != nil {
		fmt.Fprintf(stderr, "bench resume-clean: intent snapshot: %v\n", err)
		return 1
	}
	result.OpenIntent = len(live)
	if result.Worktrees+result.Branches+result.Dirty+result.Locked+result.Leased+result.OpenIntent == 0 {
		return 0
	}
	fmt.Fprintf(stdout, "bench resume: cleaned %d worktree(s), %d landed branch(es); kept %d dirty, %d locked, %d leased; %d open intent(s)\n",
		result.Worktrees, result.Branches, result.Dirty, result.Locked, result.Leased, result.OpenIntent)
	return 0
}
