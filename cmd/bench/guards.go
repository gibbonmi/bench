package main

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/benchguard"
	"github.com/gibbonmi/bench/internal/census"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/gitguard"
	"github.com/gibbonmi/bench/internal/poolkey"
	"github.com/gibbonmi/bench/internal/worktree"
	"github.com/gibbonmi/bench/internal/writeguard"
)

// guardGit is the destructive-git guard subcommand. It reads the PreToolUse envelope on
// stdin, classifies through internal/gitguard, and yields the verdict as an exit code:
// 0 allow, 2 block, with the `BLOCKED:` message on stderr, or 3 a genuine failure to run.
// The deferred recover maps any panic to 3, not Go's default exit-2, so exit 2 means
// only an intentional block, and the shim can trust it.
func guardGit(_ []string, stdin io.Reader, _ io.Writer, stderr io.Writer) (code int) {
	defer func() {
		if r := recover(); r != nil {
			code = 3
		}
	}()
	data, err := io.ReadAll(stdin)
	if err != nil {
		return 3
	}
	command := gitguard.CommandFromEnvelope(data)
	if command == "" {
		return 0
	}
	chk := gitguard.Checker{
		RefResolves:     git.RefResolves,
		BranchExists:    git.BranchExists,
		DefaultBranch:   guardDefaultBranch,
		CheckedOut:      guardCheckedOut,
		BareDestination: guardBareDestination,
	}
	label := gitguard.Classify(command, chk)
	if label == "" {
		return 0
	}
	fmt.Fprintln(stderr, gitguard.BlockMessage(label))
	return 2
}

// guardProbeRoot is the directory the guard's three push facts read. The guarded Bash
// command runs in the agent's cwd, and git.RefResolves and git.BranchExists already probe
// that directory with no root operand, so the push facts name it explicitly to reach the
// same repository.
const guardProbeRoot = "."

// guardDefaultBranch reports the repository's default branch, from the one Go owner of
// that fact. No answer denies the push, so the guard never guesses a protected name.
func guardDefaultBranch() (string, bool) { return git.ResolvedDefault(guardProbeRoot) }

// guardCheckedOut reports the checked-out branch, or no branch, from the one Go owner of
// that mapping. No answer denies a `HEAD` refspec with the unresolved class.
func guardCheckedOut() (string, bool) { return git.CheckedOutName(guardProbeRoot) }

// guardBareDestination reports the branch a bare `git push` targets, from the one Go
// owner of that fact.
func guardBareDestination() (string, bool) { return git.BarePushDestination(guardProbeRoot) }

func guardBenchFollowOn(_ []string, stdin io.Reader, _ io.Writer, stderr io.Writer) (code int) {
	defer func() {
		if recover() != nil {
			code = 3
		}
	}()
	data, err := io.ReadAll(stdin)
	if err != nil {
		return 3
	}
	command, err := benchguard.CommandFromEnvelope(data)
	if err != nil {
		fmt.Fprintln(stderr, "WARNING: block-bench-follow-on: unreadable command field — allowing Bash.")
		return 0
	}
	recordFollowOn(command)
	// The pool denial runs first. A pool reference with a Bench call after it has two
	// faults, and the pool reference is the cause the reader repairs.
	if target := benchguard.PoolReference(command, poolkey.Pools(worktree.Home())); target != "" {
		fmt.Fprintln(stderr, benchguard.PoolReferenceMessage(target))
		return 2
	}
	verdict := benchguard.Classify(command, benchguard.DefaultResolver())
	if !verdict.Blocked {
		return 0
	}
	fmt.Fprintln(stderr, verdict.Message())
	return 2
}

// guardFileWrite is the primary-checkout file-write guard subcommand. It reads the
// PreToolUse envelope on stdin, classifies through internal/writeguard, and yields the
// verdict as an exit code: 0 allow, 2 block with the `BLOCKED:` message on stderr, or 3 a
// genuine failure to run. The deferred recover maps any panic to 3, not Go's default
// exit-2, so exit 2 means only an intentional block, and the shim can trust it.
func guardFileWrite(_ []string, stdin io.Reader, _ io.Writer, stderr io.Writer) (code int) {
	defer func() {
		if recover() != nil {
			code = 3
		}
	}()
	data, err := io.ReadAll(stdin)
	if err != nil {
		return 3
	}
	path, err := writeguard.PathFromEnvelope(data)
	if err != nil {
		fmt.Fprintln(stderr, "WARNING: block-primary-file-write: unreadable file_path field — allowing the write.")
		return 0
	}
	verdict := writeguard.Classify(path, writeguard.Checker{
		RootAt:    git.RootAt,
		IsPrimary: git.IsPrimaryCheckout,
		IsTracked: guardPathTracked,
		IsIgnored: guardPathIgnored,
	})
	if !verdict.Blocked {
		return 0
	}
	fmt.Fprintln(stderr, verdict.Message())
	return 2
}

// guardPathTracked reports whether git tracks path in root. An absolute pathspec is what
// the envelope carries, and git resolves it against the repository itself.
func guardPathTracked(root, path string) bool {
	return git.OK("-C", root, "ls-files", "--error-unmatch", "--", path)
}

// guardPathIgnored reports whether root's ignore rules cover path.
func guardPathIgnored(root, path string) bool {
	return git.OK("-C", root, "check-ignore", "-q", "--", path)
}

// recordFollowOn records a raw call through the exec census. It tests the command
// text for the pool prefix before it resolves any root, so an ordinary call outside
// a Bench worktree spawns no git process. Its own failure is silent and never reaches
// the verdict: this call sits before the verdict, so no later return can skip it.
func recordFollowOn(command string) {
	home := worktree.Home()
	if !strings.Contains(command, poolkey.Pools(home)+string(filepath.Separator)) {
		return
	}
	root, err := git.Root()
	if err != nil {
		return
	}
	_ = census.Record(command, root, home, time.Now())
}
