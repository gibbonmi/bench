package worktree

import (
	"errors"
	"fmt"
	"io"
	"os/exec"

	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/usage"
)

var reauthorizeGrammar = usage.Grammar{
	Cmd:     "bench worktree reauthorize",
	Help:    "usage: " + usage.WorktreeReauthorize,
	MinArgs: 1,
	MaxArgs: 1,
	Flags: []usage.Flag{
		{Name: "--assignment", HasValue: true, NoEmptyValue: true},
		{Name: "--request", HasValue: true, NoEmptyValue: true},
		{Name: "--base", HasValue: true, NoEmptyValue: true},
		{Name: "--source-tip", HasValue: true, NoEmptyValue: true},
	},
}

var reauthorizeUnlock = func(root, path string) error {
	_, err := exec.Command("git", "-C", root, "worktree", "unlock", path).CombinedOutput()
	return err
}

var reauthorizeLock = func(root, path, reason string) error {
	_, err := exec.Command("git", "-C", root, "worktree", "lock", "--reason", reason, path).CombinedOutput()
	return err
}

// reauthorizeBeforeCAS lets the operation test model a ledger winner after the lock
// refresh but before the expected-old comparison.
var reauthorizeBeforeCAS func(*intent.Assignment)

// ReauthorizeCommand replaces a retained assignment request after exact identity proofs.
func ReauthorizeCommand(root string, args []string, stdout, stderr io.Writer) int {
	parsed, line, code := usage.Parse(reauthorizeGrammar, args)
	if line != "" {
		fmt.Fprintln(stderr, line)
		return code
	}
	if len(parsed.Flags) != 4 {
		fmt.Fprintln(stderr, reauthorizeGrammar.Help)
		return 2
	}
	id, request := parsed.Flags["--assignment"], parsed.Flags["--request"]
	base, tip := parsed.Flags["--base"], parsed.Flags["--source-tip"]
	path, err := canonicalPath(parsed.Positionals[0])
	if err != nil {
		fmt.Fprintln(stderr, "bench worktree reauthorize: worktree path is not canonical")
		return 1
	}
	resolvedTip, err := git.Output("-C", root, "rev-parse", "--verify", tip+"^{commit}")
	if err != nil {
		fmt.Fprintln(stderr, "bench worktree reauthorize: source tip is not a commit")
		return 1
	}
	rangeTip, kind, hint := diff.ResolveSourceRange(root, base, resolvedTip)
	if kind != "" {
		fmt.Fprintf(stderr, "bench worktree reauthorize: %s: %s\n", sanitize.Controls(kind), sanitize.Controls(hint))
		return 1
	}
	old, err := intent.ReauthorizeAssignment(root, id, request, func(a intent.Assignment) error {
		if a.State != intent.StateActive || a.Worktree != path {
			return errors.New("assignment identity mismatch")
		}
		evidence, err := validateOwnerMarker(root, path)
		if err != nil {
			return err
		}
		if evidence.marker.OwnerID != a.OwnerID || evidence.marker.Path != a.Worktree || evidence.registration.BranchRef != a.Branch || !evidence.registration.Locked || evidence.registration.LockReason != lockReason(a) {
			return errors.New("owner marker or branch mismatch")
		}
		branch, err := git.Output("-C", path, "symbolic-ref", "--quiet", "HEAD")
		if err != nil || branch != a.Branch {
			return errors.New("assignment branch is not checked out")
		}
		head, err := git.Output("-C", path, "rev-parse", "HEAD")
		if err != nil || head != resolvedTip {
			return errors.New("worktree tip mismatch")
		}
		branchTip, err := git.Output("-C", root, "rev-parse", "--verify", a.Branch+"^{commit}")
		if err != nil || branchTip != resolvedTip {
			return errors.New("assignment branch tip mismatch")
		}
		start, startKind, _ := diff.ResolveSourceRange(root, a.Start, resolvedTip)
		if startKind != "" || start.Base != a.Start {
			return errors.New("recorded start is not an ancestor of source tip")
		}
		return nil
	}, func(old, next intent.Assignment) (func(), error) {
		return refreshReauthorizeLock(root, path, old, next)
	}, reauthorizeBeforeCAS)
	if err != nil {
		fmt.Fprintln(stderr, "bench worktree reauthorize: "+sanitize.Controls(err.Error()))
		return 1
	}
	fmt.Fprintf(stdout, "reauthorized{assignment=%s,recorded_start=%s,approved_base=%s,source_tip=%s,state=%s}\n", old.ID, old.Start, rangeTip.Base, resolvedTip, old.State)
	return 0
}

func refreshReauthorizeLock(root, path string, old, next intent.Assignment) (func(), error) {
	if err := reauthorizeUnlock(root, path); err != nil {
		return nil, errors.New("refresh ownership lock")
	}
	if err := reauthorizeLock(root, path, lockReason(next)); err != nil {
		if restoreErr := reauthorizeLock(root, path, lockReason(old)); restoreErr != nil {
			return nil, errors.Join(errors.New("refresh ownership lock"), restoreErr)
		}
		return nil, errors.New("refresh ownership lock")
	}
	return func() {
		if reauthorizeUnlock(root, path) == nil {
			_ = reauthorizeLock(root, path, lockReason(old))
		}
	}, nil
}
