package worktree

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/toon"

	"github.com/gibbonmi/bench/internal/usage"
)

var resetGrammar = usage.Grammar{
	Cmd: usage.CommandName(usage.WorktreeReset), Help: "usage: " + usage.WorktreeReset,
	MinArgs: 1, MaxArgs: 1,
	Flags: []usage.Flag{
		{Name: "--to", HasValue: true, NoEmptyValue: true},
		{Name: "--restore", HasValue: true, NoEmptyValue: true},
		{Name: "--apply", HasValue: true, NoEmptyValue: true},
	},
}

// The reset refusal face names. The shared registry declares each face, and a raising site
// names the face it raises.
const (
	faceResetCheckoutConflicted = "reset-checkout-conflicted"
	faceResetPlanStale          = "reset-plan-stale"
	faceResetTreeMissing        = "reset-tree-missing"
	faceResetHandback           = "reset-handback"
)

// ResetCommand plans or applies a recoverable return to an assignment checkpoint.
func ResetCommand(root, home string, args []string, stdout, stderr io.Writer) int {
	return resetWith(defaultJoins(), newAmbient(home, stderr), root, args, stdout, stderr)
}

func resetWith(j joins, a ambient, root string, args []string, stdout, stderr io.Writer) int {
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
	refuse := func(err error) int {
		return landRefusalError(stdout, resetFaceRoute(err, resetRerun(parsed.Flags, parsed.Positionals[0])))
	}
	for _, flag := range []string{"--to", "--restore"} {
		if !lineSafe(parsed.Flags[flag]) {
			return refuse(errors.New(flag + " contains control characters"))
		}
	}
	plan, err := planReset(root, parsed.Positionals[0], parsed.Flags["--to"], parsed.Flags["--restore"])
	if err != nil {
		return refuse(err)
	}
	if fingerprint, apply := parsed.Flags["--apply"]; apply {
		return applyReset(j, root, a.home, plan, fingerprint, stdout, refuse)
	}
	next := ""
	if plan.action != "none" {
		next = "," + refusalroute.NextField + "=" + resetPlanCommand(plan) + " --apply " + plan.fingerprint
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
		table, err := toon.Table(resetPathsTable, []string{"path", "status"}, rows)
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

// resetRerun is the caller's own reset plan, with the values it passed. A value that is not
// line-safe prints its placeholder. An apply re-runs as its plan, because the plan prints the
// apply that the checkout takes once the cause is clear.
func resetRerun(flags map[string]string, target string) string {
	flag, placeholder := "--to", "<commit>"
	if flags["--restore"] != "" {
		flag, placeholder = "--restore", "<ref>"
	}
	return resetCommand(flag, landingRerunArg(flags[flag], placeholder), landingRerunArg(target, "<target>"))
}

// resetFaceRoute attaches the caller's own re-run to a reset refusal. A refusal that names a
// registered face keeps it, and every other cause hands back to the reviewer under its own
// sentence.
func resetFaceRoute(err error, rerun string) error {
	if raised := raisedRefusal(err); raised.face == "" {
		return landingFaceRefusal(faceResetHandback, raised, rerun, "")
	}
	return err
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
	plan.tip, err = git.ResolveCommit(root, selected.Branch)
	if err != nil {
		return resetPlan{}, err
	}
	if envelope != "" {
		plan.mode = "restore"
		plan.manifest, err = readResetRestore(root, selected, envelope)
		checkpoint = plan.manifest.Base
	} else {
		checkpoint, err = git.ResolveCommit(selected.Worktree, checkpoint, "--quiet", "--end-of-options")
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
		return resetPlan{}, landingFaceRefusal(faceResetCheckoutConflicted, refusal{values: map[string]string{refusalroute.FactAssignmentID: selected.ID}}, "", "")
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
// the move writes them into the checkout. The collision rule is TreePaths.Collides.
func ignoredCollisions(path string, trees []string) ([]string, error) {
	ignored, err := ignoredListing(path)
	if err != nil {
		return nil, fmt.Errorf("read ignored inventory: %w", err)
	}
	materialized, err := treePaths(path, trees...)
	if err != nil {
		return nil, fmt.Errorf("read materialized tree: %w", err)
	}
	var collisions []string
	for record := range bytes.SplitSeq(ignored, []byte{0}) {
		if len(record) > 0 && materialized.Collides(string(record)) {
			collisions = append(collisions, string(record))
		}
	}
	return collisions, nil
}
