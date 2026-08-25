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
		{Name: "--assignment", HasValue: true, NoEmptyValue: true, Required: true},
		{Name: "--request", HasValue: true, NoEmptyValue: true, Required: true},
		{Name: "--base", HasValue: true, NoEmptyValue: true, Required: true},
		{Name: "--source-tip", HasValue: true, NoEmptyValue: true, Required: true},
	},
}

// unlockWorktree and lockWorktree are the ownership-lock effects the reauthorize refresh
// performs. The joins value carries them, together with reauthorizeBeforeCAS, which lets
// the operation test model a ledger winner after the lock refresh but before the
// expected-old comparison.
func unlockWorktree(root, path string) error {
	_, err := exec.Command("git", "-C", root, "worktree", "unlock", path).CombinedOutput()
	return err
}

func lockWorktree(root, path, reason string) error {
	_, err := exec.Command("git", "-C", root, "worktree", "lock", "--reason", reason, path).CombinedOutput()
	return err
}

// ReauthorizeCommand replaces a retained assignment request after exact identity proofs.
func ReauthorizeCommand(root, home string, args []string, stdout, stderr io.Writer) int {
	return reauthorizeWith(defaultJoins(), root, home, args, stdout, stderr)
}

// reauthorizeWith is ReauthorizeCommand with the seam set resolved explicitly at the
// caller's boundary.
func reauthorizeWith(j joins, root, _ string, args []string, stdout, stderr io.Writer) int {
	parsed, line, code := usage.Parse(reauthorizeGrammar, args)
	if line != "" {
		fmt.Fprintln(stderr, line)
		return code
	}
	id, request := parsed.Flags["--assignment"], parsed.Flags["--request"]
	base, tip := parsed.Flags["--base"], parsed.Flags["--source-tip"]
	path, err := canonicalPath(resolveVerbOperand(root, parsed.Positionals[0]))
	if err != nil {
		fmt.Fprintln(stderr, "bench worktree reauthorize: worktree path is not canonical")
		return 1
	}
	resolvedTip, err := git.Output("-C", root, "rev-parse", "--verify", tip+"^{commit}")
	if err != nil {
		fmt.Fprintln(stderr, "bench worktree reauthorize: source tip is not a commit")
		return 1
	}
	approvedBase := ""
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
			return refusalError{refusal{detail: "recorded start is not an ancestor of source tip", wanted: a.Start}}
		}
		rangeTip, kind, hint := diff.ResolveSourceRange(root, base, resolvedTip)
		if kind == "--base is not an ancestor" {
			return refusalError{refusal{detail: "review base is not an ancestor of source tip", wanted: a.Start}}
		}
		if kind != "" {
			return errors.New(kind + ": " + hint)
		}
		approvedBase = rangeTip.Base
		return nil
	}, func(old, next intent.Assignment) (func(), error) {
		return refreshReauthorizeLock(j, root, path, old, next)
	}, j.reauthorizeBeforeCAS)
	if err != nil {
		var typed refusalError
		if errors.As(err, &typed) {
			fmt.Fprintln(stderr, "bench worktree reauthorize: "+typed.Error())
			return 1
		}
		fmt.Fprintln(stderr, "bench worktree reauthorize: "+sanitize.Controls(err.Error()))
		return 1
	}
	fmt.Fprintf(stdout, "reauthorized{assignment=%s,recorded_start=%s,approved_base=%s,source_tip=%s,state=%s}\n", old.ID, old.Start, approvedBase, resolvedTip, old.State)
	return 0
}

func refreshReauthorizeLock(j joins, root, path string, old, next intent.Assignment) (func(), error) {
	if err := j.reauthorizeUnlock(root, path); err != nil {
		return nil, errors.New("refresh ownership lock")
	}
	if err := j.reauthorizeLock(root, path, lockReason(next)); err != nil {
		if restoreErr := j.reauthorizeLock(root, path, lockReason(old)); restoreErr != nil {
			return nil, errors.Join(errors.New("refresh ownership lock"), restoreErr)
		}
		return nil, errors.New("refresh ownership lock")
	}
	return func() {
		if j.reauthorizeUnlock(root, path) == nil {
			_ = j.reauthorizeLock(root, path, lockReason(old))
		}
	}, nil
}
