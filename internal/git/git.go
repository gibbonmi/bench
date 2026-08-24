// Package git is the one source of git subprocess invocation for the AXI query
// commands. Every ported parser shells out to git through here. Git stays the source of
// repository truth — root, config, merge-base, diff — exactly as the shell commands did,
// so the invocation form lives in one place.
package git

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
)

// refCheckTimeout is the hook-scoped fail-safe for destructive-git classification,
// unlike the policy-owned worktree discovery bound below. It bounds the ref and branch
// existence probes the destructive-git guard runs per classification — internal/gitguard's
// checkout and forced-creation verdicts. A hung git must never stall a PreToolUse Bash
// hook, so the code bounds each probe at two seconds and resolves it to its caller's
// fail-safe default.
const refCheckTimeout = 2 * time.Second

var worktreeListTimeout = bounds.WorktreeListTimeout

// SetWorktreeListTimeoutForTest installs a test-only discovery bound and restores it.
func SetWorktreeListTimeoutForTest(limit time.Duration) func() {
	previous := worktreeListTimeout
	worktreeListTimeout = limit
	return func() { worktreeListTimeout = previous }
}

// RefResolves reports whether arg names a commit-ish that resolves in the process
// working directory. That directory is the agent's cwd, where the guarded Bash command
// runs, not a fixed root: the check runs `git rev-parse --verify --quiet
// <arg>^{commit}`. On a timeout, or a failure to run git at all, it returns false. The
// function treats an undeterminable target as unresolvable, so a checkout of it fails
// closed and blocks; the fail-toward-blocking default is deliberate.
func RefResolves(arg string) bool {
	exitZero, ran := refCheck(arg + "^{commit}")
	if !ran {
		return false
	}
	return exitZero
}

// BranchExists reports whether refs/heads/<name> exists in the cwd repo. On a timeout or
// a failure to run git, it returns TRUE, the opposite default from RefResolves. Its only
// caller — forced branch or switch creation — must block when it cannot rule out that the
// force would clobber an existing branch.
func BranchExists(name string) bool {
	exitZero, ran := refCheck("refs/heads/" + name)
	if !ran {
		return true
	}
	return exitZero
}

// refCheck runs `git rev-parse --verify --quiet <ref>` under the time bound and reports
// (git exited zero, git ran to a verdict). ran is false when the probe timed out or git
// could not run — the undeterminable branch each caller resolves to its own fail-safe
// default.
func refCheck(ref string) (exitZero, ran bool) {
	ctx, cancel := context.WithTimeout(context.Background(), refCheckTimeout)
	defer cancel()
	err := exec.CommandContext(ctx, "git", "rev-parse", "--verify", "--quiet", ref).Run()
	if ctx.Err() == context.DeadlineExceeded {
		return false, false
	}
	if err == nil {
		return true, true
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return false, true
	}
	return false, false
}

// Root returns the working tree's top-level directory. It returns an error when the cwd
// is not inside a git repository — the `not in a git repository` posture of every
// command.
func Root() (string, error) {
	return Output("rev-parse", "--show-toplevel")
}

// RootAt returns the top-level directory for the repository containing dir. Harness
// events carry an explicit cwd, so their routing must not depend on the hook process's
// ambient working directory.
func RootAt(dir string) (string, error) {
	return Output("-C", dir, "rev-parse", "--show-toplevel")
}

// LocalBranches is the sole local-head enumeration and parsing owner.
func LocalBranches(root string) ([]string, error) {
	out, err := Output("-C", root, "for-each-ref", "--format=%(refname:short)", "refs/heads/")
	if err != nil {
		return nil, err
	}
	if out == "" {
		return []string{}, nil
	}
	return strings.Split(out, "\n"), nil
}

// DeleteBranchExact removes one full branch ref only while it still has the
// caller-proven OID.
func DeleteBranchExact(root, ref, oid string) error {
	out, err := exec.Command("git", "-C", root, "update-ref", "-d", ref, oid).CombinedOutput()
	if err != nil {
		return fmt.Errorf("update exact branch ref: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// PruneLandedBranches removes local non-default branches whose work the default branch
// already contains. It skips a branch checked out in a registered worktree or named by a
// caller-protected lifecycle record. The exact old OID makes each deletion fail closed if
// the branch moves after classification.
func PruneLandedBranches(root string, protectedBranches []string) (int, error) {
	def, ok := ResolvedDefault(root)
	if !ok {
		return 0, errors.New("git repository has no resolvable default branch")
	}
	worktrees, err := Worktrees(root)
	if err != nil {
		return 0, fmt.Errorf("worktree discovery failed: %w", err)
	}
	checkedOut := map[string]bool{}
	for _, worktree := range worktrees {
		if worktree.Branch != "" {
			checkedOut[worktree.Branch] = true
		}
	}
	for _, branch := range protectedBranches {
		checkedOut[strings.TrimPrefix(branch, "refs/heads/")] = true
	}
	branches, err := LocalBranches(root)
	if err != nil {
		return 0, fmt.Errorf("git local branches: %w", err)
	}
	pruned := 0
	for _, branch := range branches {
		if branch == def || checkedOut[branch] {
			continue
		}
		landed, _, err := LandedInDefault(root, branch, def)
		if err != nil {
			return pruned, fmt.Errorf("git landedness %s: %w", branch, err)
		}
		if !landed {
			continue
		}
		oid, err := Output("-C", root, "rev-parse", "--verify", "refs/heads/"+branch+"^{commit}")
		if err != nil {
			return pruned, fmt.Errorf("git branch identity %s: %w", branch, err)
		}
		if err := DeleteBranchExact(root, "refs/heads/"+branch, oid); err != nil {
			return pruned, fmt.Errorf("delete landed branch %s: %w", branch, err)
		}
		pruned++
	}
	return pruned, nil
}

// GateCacheFile is the filename of the cached gate verdict, written under the git
// directory (never the worktree, so it is never a diff or commit candidate). It is the
// one gate owner resolves and composes with its absolute Git-directory path.
const GateCacheFile = "bench-last-gate"

// LandedInDefault proves a local branch landed by ancestry, patch containment, or
// reverse-applying the branch's cumulative diff to def's tree. The reverse-apply proof
// alone proves a squash-landing, where the branch's commits compose into one commit no
// other proof can see. git cherry cannot prove merge-only content, so the function keeps
// the merge check. It runs that check before either content proof, which keeps
// conservatism for the fourth proof too.
//
// A true verdict is the sole authority every cleanup path deletes a branch on, so
// ambiguity never rounds up. The function reports not landed whenever the reverse-apply
// proof cannot generate, apply cleanly, or represent the diff byte-for-byte —
// reverseAppliesToDefault names the exact forms. The cost of refusing is an orphaned
// branch. The cost of a wrong true is lost work with nothing standing behind it.
func LandedInDefault(root, branch, def string) (landed, byContent bool, err error) {
	ancestor := exec.Command("git", "-C", root, "merge-base", "--is-ancestor", branch, def)
	if err := ancestor.Run(); err == nil {
		return true, false, nil
	} else if exit, ok := err.(*exec.ExitError); !ok || exit.ExitCode() != 1 {
		return false, false, fmt.Errorf("merge-base %s against %s: %w", branch, def, err)
	}
	merges, err := Output("-C", root, "rev-list", "--merges", "--max-count=1", def+".."+branch)
	if err != nil {
		return false, false, fmt.Errorf("list merges on %s: %w", branch, err)
	}
	if merges != "" {
		return false, false, nil
	}
	out, err := Output("-C", root, "cherry", def, branch)
	if err != nil {
		return false, false, fmt.Errorf("compare patches on %s: %w", branch, err)
	}
	for _, line := range strings.Split(out, "\n") {
		if line != "" && !strings.HasPrefix(line, "-") {
			if reverseAppliesToDefault(root, branch, def) {
				return true, true, nil
			}
			return false, false, nil
		}
	}
	return true, true, nil
}

// Output runs `git <args>` and returns stdout with a single trailing newline trimmed.
// err is non-nil on a nonzero exit. Callers use it for single-value reads — a root, a
// config key, a resolved sha — where the trailing newline is noise.
func Output(args ...string) (string, error) {
	var out bytes.Buffer
	cmd := exec.Command("git", args...)
	cmd.Stdout = &out
	err := cmd.Run()
	return strings.TrimRight(out.String(), "\n"), err
}

// OK reports whether `git <args>` exits zero. It discards all output. This is the test
// form — cat-file -e, merge-base --is-ancestor — where only the exit code matters.
func OK(args ...string) bool {
	return exec.Command("git", args...).Run() == nil
}

// Raw runs `git <args>` and returns stdout verbatim, with no trimming, along with the
// exit status. Callers use it for `diff -z` output, whose NUL framing and trailing bytes
// are load-bearing.
func Raw(args ...string) ([]byte, error) {
	var out bytes.Buffer
	cmd := exec.Command("git", args...)
	cmd.Stdout = &out
	err := cmd.Run()
	return out.Bytes(), err
}
