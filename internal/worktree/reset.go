package worktree

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

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
		{Name: "--to", HasValue: true, NoEmptyValue: true},
		{Name: "--restore", HasValue: true, NoEmptyValue: true},
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
	if (parsed.Flags["--to"] == "") == (parsed.Flags["--restore"] == "") {
		fmt.Fprintln(stderr, resetGrammar.Help)
		return 2
	}
	for _, flag := range []string{"--to", "--restore"} {
		if !lineSafe(parsed.Flags[flag]) {
			return landRefusal(stdout, flag+" contains control characters")
		}
	}
	plan, err := planReset(root, parsed.Positionals[0], parsed.Flags["--to"], parsed.Flags["--restore"])
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
	if plan.envelope != "" {
		next += ",envelope=" + plan.envelope
	}
	fmt.Fprintf(stdout, "reset_plan{worktree=%s,mode=%s,action=%s,checkpoint=%s,head=%s,ref=%s,tip=%s,tracked=%s,lock=%s,preserve=%s%s,fingerprint=%s}\n",
		plan.assignment.ID, plan.mode, plan.action, plan.checkpoint, plan.head, plan.ref, plan.tip, plan.tracked, plan.lock, plan.preserve, next, plan.fingerprint)
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

// resetCommand owns the layout of every reset command the verb prints, so the plan's
// apply, the record's restore, and the exit-3 next cell never drift apart.
func resetCommand(flag, operand, id string) string {
	return "bench worktree reset " + flag + " " + operand + " " + id
}

func resetPlanCommand(plan resetPlan) string {
	if plan.envelope != "" {
		return resetCommand("--restore", plan.envelope, plan.assignment.ID)
	}
	return resetCommand("--to", plan.checkpoint, plan.assignment.ID)
}

type resetPlan struct {
	assignment                                   intent.Assignment
	checkpoint, head, ref, tip                   string
	action, tracked, lock, preserve, fingerprint string
	status                                       []byte
	paths                                        []git.PorcelainEntry
	registrationLocked                           bool
	mode, envelope                               string
	manifest                                     recoveryManifest
}

func planReset(root, operand, checkpoint, envelope string) (resetPlan, error) {
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
	plan := resetPlan{assignment: selected, checkpoint: checkpoint, mode: "reset", envelope: envelope, action: "reset", tracked: "dirty", lock: "ok", preserve: "envelope"}
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
	if envelope != "" {
		plan.mode = "restore"
		plan.manifest, err = readResetRestore(root, selected, envelope)
		checkpoint = plan.manifest.Base
	} else {
		checkpoint, err = git.Output("-C", selected.Worktree, "rev-parse", "--verify", "--quiet", "--end-of-options", checkpoint+"^{commit}")
	}
	if err != nil && envelope != "" {
		return resetPlan{}, err
	}
	if err != nil || checkpoint == "" {
		return resetPlan{}, errors.New("checkpoint is not a commit")
	}
	plan.checkpoint = checkpoint
	if envelope == "" {
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
	plan.status, err = checkoutStatus(selected.Worktree)
	if err != nil {
		return resetPlan{}, err
	}
	plan.paths, err = git.ParsePorcelainZStrict(plan.status)
	if err != nil {
		return resetPlan{}, err
	}
	index, err := rawIndexEntries(selected.Worktree)
	if err != nil {
		return resetPlan{}, err
	}
	_, conflicted, err := parseIndexEntries(index)
	if err != nil {
		return resetPlan{}, err
	}
	if conflicted {
		return resetPlan{}, refusalError{refusal{detail: "checkout is conflicted", next: "bench worktree clean " + selected.ID}}
	}
	hidden, err := hiddenIndexFlags(selected.Worktree)
	if err != nil {
		return resetPlan{}, err
	}
	if len(hidden) > 0 {
		return resetPlan{}, refusalError{refusal{detail: "index carries hidden flags", paths: hidden}}
	}
	materialized := []string{checkpoint}
	if envelope != "" {
		materialized = []string{plan.manifest.Tip, plan.manifest.Base, plan.manifest.Layers["working"]}
	}
	collisions, err := ignoredCollisions(selected.Worktree, materialized)
	if err != nil {
		return resetPlan{}, err
	}
	if len(collisions) > 0 {
		return resetPlan{}, refusalError{refusal{detail: "ignored content would be overwritten", paths: collisions}}
	}
	if len(plan.status) == 0 {
		plan.tracked = "clean"
		targetTip := checkpoint
		if envelope != "" {
			targetTip = plan.manifest.Tip
		}
		if plan.head == checkpoint && plan.tip == targetTip {
			plan.preserve = "none"
		}
		if envelope == "" && plan.head == checkpoint && plan.ref == selected.Branch && plan.lock == "ok" {
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
			[]byte(selected.ID), []byte(plan.mode), []byte(checkpoint), []byte(plan.head), []byte(plan.ref),
			[]byte(plan.tip), []byte(evidence.registration.LockReason), plan.status, index, []byte(content), []byte(envelope))
	}
	return plan, err
}

// hiddenIndexFlags lists every index entry marked assume-unchanged or skip-worktree.
// Such an entry hides its edit from the status, the content identity, and the staged
// listing, so no envelope could carry it.
func hiddenIndexFlags(path string) ([]string, error) {
	listing, err := git.Raw("--no-optional-locks", "-C", path, "ls-files", "-v", "-z")
	if err != nil {
		return nil, fmt.Errorf("read index flags: %w", err)
	}
	var hidden []string
	for record := range bytes.SplitSeq(listing, []byte{0}) {
		if len(record) < 2 {
			continue
		}
		if tag := record[0]; tag == 'h' || tag == 'S' || tag == 's' {
			hidden = append(hidden, string(record[2:]))
		}
	}
	return hidden, nil
}

// ignoredCollisions lists every ignored path that the given trees would overwrite when
// the move writes them into the checkout. A tracked file whose path is a parent
// directory of the ignored path collides, because the move replaces that directory with
// the file. A tracked file below the ignored path collides too, because the move
// replaces the ignored file with a directory.
func ignoredCollisions(path string, trees []string) ([]string, error) {
	ignored, err := ignoredListing(path)
	if err != nil {
		return nil, fmt.Errorf("read ignored inventory: %w", err)
	}
	tracked, above := map[string]bool{}, map[string]bool{}
	for _, tree := range trees {
		names, err := git.Raw("-C", path, "ls-tree", "-r", "-z", "--name-only", tree)
		if err != nil {
			return nil, fmt.Errorf("read materialized tree: %w", err)
		}
		for name := range bytes.SplitSeq(names, []byte{0}) {
			if len(name) == 0 {
				continue
			}
			tracked[string(name)] = true
			for dir := string(name); strings.Contains(dir, "/"); {
				dir = dir[:strings.LastIndexByte(dir, '/')]
				above[dir] = true
			}
		}
	}
	var collisions []string
	for record := range bytes.SplitSeq(ignored, []byte{0}) {
		if len(record) == 0 {
			continue
		}
		if above[string(record)] {
			collisions = append(collisions, string(record))
			continue
		}
		for candidate := string(record); ; candidate = candidate[:strings.LastIndexByte(candidate, '/')] {
			if tracked[candidate] {
				collisions = append(collisions, string(record))
				break
			}
			if !strings.Contains(candidate, "/") {
				break
			}
		}
	}
	return collisions, nil
}
