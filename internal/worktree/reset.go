package worktree

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/toon"

	"github.com/gibbonmi/bench/internal/usage"
)

var resetGrammar = usage.Grammar{
	Cmd: "bench worktree reset", Help: "usage: " + usage.WorktreeReset,
	MinArgs: 1, MaxArgs: 1,
	Flags: []usage.Flag{
		{Name: "--to", HasValue: true, NoEmptyValue: true, Required: true},
		{Name: "--apply", HasValue: true, NoEmptyValue: true},
	},
}

// ResetCommand plans or applies a recoverable return to an assignment checkpoint.
func ResetCommand(root, home string, args []string, stdout, stderr io.Writer) int {
	return resetWith(defaultJoins(), root, home, args, stdout, stderr)
}

func resetWith(j joins, root, home string, args []string, stdout, stderr io.Writer) int {
	parsed, line, code := usage.Parse(resetGrammar, args)
	if line != "" {
		if code == 0 {
			fmt.Fprintln(stdout, line)
		} else {
			fmt.Fprintln(stderr, line)
		}
		return code
	}
	if !lineSafe(parsed.Flags["--to"]) {
		return landRefusal(stdout, "--to contains control characters")
	}
	plan, err := planReset(root, parsed.Positionals[0], parsed.Flags["--to"])
	if err != nil {
		return landRefusalError(stdout, err)
	}
	if fingerprint, apply := parsed.Flags["--apply"]; apply {
		return applyReset(j, root, home, plan, fingerprint, stdout)
	}
	next := ""
	if plan.action != "none" {
		next = ",next=" + resetPlanCommand(plan) + " --apply " + plan.fingerprint
	}
	fmt.Fprintf(stdout, "reset_plan{worktree=%s,mode=reset,action=%s,checkpoint=%s,head=%s,ref=%s,tip=%s,tracked=%s,lock=%s,preserve=%s%s,fingerprint=%s}\n",
		plan.assignment.ID, plan.action, plan.checkpoint, plan.head, plan.ref, plan.tip, plan.tracked, plan.lock, plan.preserve, next, plan.fingerprint)
	var rows [][]string
	for _, entry := range plan.paths {
		if entry.Status != "" {
			rows = append(rows, []string{sanitize.Controls(entry.Path), entry.Status})
		}
	}
	if len(rows) > 0 {
		table, err := toon.Table("reset_paths", []string{"path", "status"}, rows)
		if err != nil {
			return landRefusalError(stdout, err)
		}
		fmt.Fprint(stdout, table)
	}
	return 0
}

func resetPlanCommand(plan resetPlan) string {
	return "bench worktree reset --to " + plan.checkpoint + " " + plan.assignment.ID
}

type resetPlan struct {
	assignment                                   intent.Assignment
	checkpoint, head, ref, tip                   string
	action, tracked, lock, preserve, fingerprint string
	status                                       []byte
	paths                                        []git.PorcelainEntry
	registrationLocked                           bool
}

func planReset(root, operand, checkpoint string) (resetPlan, error) {
	assignments, err := intent.Assignments(root)
	if err != nil {
		return resetPlan{}, errors.New("assignment ledger is unreadable")
	}
	selected, err := selectAssignment(assignments, operand)
	if err != nil {
		return resetPlan{}, err
	}
	if !landingActiveState(selected.State) {
		return resetPlan{}, componentRefusal(componentAssignmentState, selected.ID, string(selected.State), string(intent.StateActive))
	}
	if _, err := os.Stat(selected.Worktree); os.IsNotExist(err) {
		return resetPlan{}, missingTreeRefusal(root, selected)
	}
	evidence, err := ownerMarkerRefusal(root, selected.Worktree, selected)
	if err != nil {
		return resetPlan{}, err
	}
	plan := resetPlan{assignment: selected, checkpoint: checkpoint, action: "reset", tracked: "dirty", lock: "ok", preserve: "envelope"}
	plan.registrationLocked = evidence.registration.Locked
	plan.head, err = git.Output("-C", selected.Worktree, "rev-parse", "HEAD^{commit}")
	if err != nil {
		return resetPlan{}, err
	}
	plan.ref, _ = git.Output("-C", selected.Worktree, "symbolic-ref", "--quiet", "HEAD")
	if plan.ref == "" {
		plan.ref = "detached"
	}
	plan.tip, err = git.Output("-C", root, "rev-parse", "--verify", selected.Branch+"^{commit}")
	if err != nil {
		return resetPlan{}, err
	}
	checkpoint, err = git.Output("-C", root, "rev-parse", "--verify", "--quiet", "--end-of-options", checkpoint+"^{commit}")
	if err != nil || checkpoint == "" {
		return resetPlan{}, errors.New("checkpoint is not a commit")
	}
	plan.checkpoint = checkpoint
	fromStart, err := authorization.IsAncestor(root, selected.Start, checkpoint)
	if err != nil {
		return resetPlan{}, err
	}
	toTip, err := authorization.IsAncestor(root, checkpoint, plan.tip)
	if err != nil {
		return resetPlan{}, err
	}
	if !fromStart || !toTip {
		return resetPlan{}, refusalError{refusal{detail: "checkpoint is outside the assignment history", wanted: selected.Start + ".." + plan.tip}}
	}
	nested, err := classifyNestedState(selected.Worktree)
	if err != nil || nested == nestedUnknown {
		return resetPlan{}, errors.New("nested repository state is unknown")
	}
	if nested == nestedDirty {
		return resetPlan{}, errors.New("nested repository is dirty")
	}
	if nested == nestedEmbeddedClean || nested == nestedEmbeddedDirty {
		return resetPlan{}, errors.New("embedded repository is retained")
	}
	lease, err := LeaseFile(selected.Worktree)
	if err != nil {
		return resetPlan{}, err
	}
	if ProbeLease(lease) == LeaseLive {
		body, err := os.ReadFile(lease)
		if err != nil {
			return resetPlan{}, err
		}
		pid, ok := leaseOwnerPID(body)
		if !ok || pid != os.Getpid() {
			return resetPlan{}, errors.New("assignment has a live lease")
		}
	}
	if !evidence.registration.Locked || evidence.registration.LockReason != lockReason(selected) {
		plan.lock = "repair"
	}
	plan.status, err = resetStatus(selected.Worktree)
	if err != nil {
		return resetPlan{}, err
	}
	plan.paths, err = git.ParsePorcelainZStrict(plan.status)
	if err != nil {
		return resetPlan{}, err
	}
	_, conflicted, err := readIndexEntries(selected.Worktree)
	if err != nil {
		return resetPlan{}, err
	}
	if conflicted {
		return resetPlan{}, refusalError{refusal{detail: "checkout is conflicted", next: "bench worktree clean " + selected.ID}}
	}
	if len(plan.status) == 0 {
		plan.tracked = "clean"
		if plan.head == checkpoint && plan.tip == checkpoint {
			plan.preserve = "none"
		}
		if plan.head == checkpoint && plan.ref == selected.Branch && plan.lock == "ok" {
			plan.action, plan.fingerprint = "none", "none"
		}
	}
	if plan.action != "none" {
		content, err := explicitContentIdentity(selected.Worktree)
		if err != nil {
			return resetPlan{}, err
		}
		common, err := git.CommonDir(root)
		if err != nil {
			return resetPlan{}, err
		}
		plan.fingerprint = fingerprintParts([]byte("bench-reset/v1"), []byte(common), []byte(selected.OwnerID),
			[]byte(selected.ID), []byte("reset"), []byte(checkpoint), []byte(plan.head), []byte(plan.ref),
			[]byte(plan.tip), []byte(evidence.registration.LockReason), plan.status, []byte(content), []byte(""))
	}
	return plan, err
}

func resetStatus(path string) ([]byte, error) {
	return git.Raw("--no-optional-locks", "-C", path, "status", "--porcelain=v1", "-z", "--untracked-files=all", "--ignore-submodules=none")
}
