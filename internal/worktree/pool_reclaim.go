package worktree

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
	"github.com/gibbonmi/bench/internal/worktree/reclaimpolicy"
)

// poolReclaimFields is the plan's row schema: which key, what the predicate decided, and
// the one fact that decided it.
var poolReclaimFields = []string{"key", "verdict", "reason"}

// The verdict words, the classification tables, and every reclaim decision live in the
// pure policy child; this parent owns the filesystem reads, the enumeration, the
// deletion, and the rendering, translating each key's state into typed facts once.
type poolKeyVerdict = reclaimpolicy.KeyVerdict

const (
	poolVerdictReclaim  = reclaimpolicy.VerdictReclaim
	poolVerdictRetain   = reclaimpolicy.VerdictRetain
	poolVerdictRemoved  = reclaimpolicy.VerdictRemoved
	poolVerdictRetained = reclaimpolicy.VerdictRetained
)

// poolReclaimPlan is one reading of the pool: every key it holds, in name order, and the
// fingerprint over the reclaimable subset.
type poolReclaimPlan struct {
	verdicts    []poolKeyVerdict
	fingerprint string
}

func (p poolReclaimPlan) reclaimableCount() int { return reclaimpolicy.ReclaimableCount(p.verdicts) }

// poolKeysDir is the pool parent this command reads. Nothing else in the tree enumerates
// it: every other reclamation path is anchored at a repository root. That is exactly why
// a key whose repository was deleted is unreachable from all of them.
func poolKeysDirAt(home string) string { return filepath.Join(home, "worktrees") }

// planPoolReclaim classifies every key under the pool parent. An absent pool parent is
// the zero-row answer rather than an error. A home that has never leased a worktree is a
// clean pool, not a broken one. The home is the caller's explicit boundary resolution.
func planPoolReclaim(root, home string) (poolReclaimPlan, error) {
	current := filepath.Base(poolAt(home, canonicalRoot(root)))
	pool := poolKeysDirAt(home)
	entries, err := os.ReadDir(pool)
	if err != nil {
		if os.IsNotExist(err) {
			return poolReclaimPlan{fingerprint: fingerprintPoolReclaim(nil)}, nil
		}
		return poolReclaimPlan{}, err
	}
	plan := poolReclaimPlan{verdicts: make([]poolKeyVerdict, 0, len(entries))}
	for _, entry := range entries {
		name := entry.Name()
		// The current repository's key never reaches a filesystem question: the policy
		// protects it on the Current fact alone, so no facts are gathered for it.
		if name == current {
			plan.verdicts = append(plan.verdicts, reclaimpolicy.ClassifyKey(reclaimpolicy.KeyFacts{Name: name, Current: true}))
			continue
		}
		plan.verdicts = append(plan.verdicts, classifyPoolKey(filepath.Join(pool, name), name))
	}
	plan.fingerprint = fingerprintPoolReclaim(plan.verdicts)
	return plan, nil
}

// classifyPoolKey translates one key's filesystem state and hands the policy the verdict.
func classifyPoolKey(path, name string) poolKeyVerdict {
	return reclaimpolicy.ClassifyKey(gatherPoolKeyFacts(path, name))
}

// gatherPoolKeyFacts is the reclaim fact adapter: it performs every filesystem read for
// one key exactly once and stops reading at the first disqualifying shape. Lstat is used
// throughout, so a symlink where a key, a child, or a `.git` belongs is recorded
// unfollowed, and only a regular pointer file is ever opened — a FIFO there is a shape
// fact, never a blocking read.
func gatherPoolKeyFacts(path, name string) reclaimpolicy.KeyFacts {
	facts := reclaimpolicy.KeyFacts{Name: name}
	facts.Shape, facts.ShapeErr = lstatShape(path)
	if facts.Shape != reclaimpolicy.ShapeDir {
		return facts
	}
	children, err := os.ReadDir(path)
	if err != nil {
		facts.ListErr = err.Error()
		return facts
	}
	for _, child := range children {
		facts.Children = append(facts.Children, gatherPoolChildFacts(filepath.Join(path, child.Name()), child.Name()))
	}
	return facts
}

// gatherPoolChildFacts translates one top-level entry: its shape, its `.git` pointer's
// shape and body, and — for a parseable absolute target — the one Lstat probe whose
// verdict separates a proven absence from an open question.
func gatherPoolChildFacts(path, name string) reclaimpolicy.ChildFacts {
	child := reclaimpolicy.ChildFacts{Name: name}
	child.Shape, child.ShapeErr = lstatShape(path)
	if child.Shape != reclaimpolicy.ShapeDir {
		return child
	}
	pointer := filepath.Join(path, ".git")
	child.Pointer.Shape, child.Pointer.ShapeErr = lstatShape(pointer)
	if child.Pointer.Shape != reclaimpolicy.ShapeRegular {
		return child
	}
	body, err := os.ReadFile(pointer)
	if err != nil {
		child.Pointer.ReadErr = err.Error()
		return child
	}
	child.Pointer.Body = string(body)
	target, ok := reclaimpolicy.GitdirTarget(child.Pointer.Body)
	if !ok || !filepath.IsAbs(target) {
		return child
	}
	if _, err := os.Lstat(target); err == nil {
		child.Pointer.TargetExistence = reclaimpolicy.ExistencePresent
	} else if os.IsNotExist(err) {
		child.Pointer.TargetExistence = reclaimpolicy.ExistenceAbsent
	} else {
		child.Pointer.TargetErr = err.Error()
	}
	return child
}

// lstatShape translates one Lstat into the policy's typed entry shape. Absence is proven
// only by os.IsNotExist; every other error is an unreadable shape carrying its text.
func lstatShape(path string) (reclaimpolicy.EntryShape, string) {
	info, err := os.Lstat(path)
	switch {
	case err != nil && os.IsNotExist(err):
		return reclaimpolicy.ShapeMissing, err.Error()
	case err != nil:
		return reclaimpolicy.ShapeUnreadable, err.Error()
	case info.Mode()&os.ModeSymlink != 0:
		return reclaimpolicy.ShapeSymlink, ""
	case info.IsDir():
		return reclaimpolicy.ShapeDir, ""
	case info.Mode().IsRegular():
		return reclaimpolicy.ShapeRegular, ""
	default:
		return reclaimpolicy.ShapeOther, ""
	}
}

// fingerprintPoolReclaim digests the policy-derived plan material with the package's
// shared fingerprint owner, so what a plan covers has one source and how bytes digest
// has another.
func fingerprintPoolReclaim(verdicts []poolKeyVerdict) string {
	return fingerprintParts(reclaimpolicy.FingerprintMaterial(verdicts)...)
}

// ReclaimCommand plans the reclaimable keys in `$BENCH_HOME/worktrees`, and with
// `--apply <fingerprint>` removes the ones that plan named. It runs inside a repository
// like every other `bench worktree` subcommand, because the current repository's key is
// the one key it may never name.
func ReclaimCommand(root string, args []string, stdout, stderr io.Writer) int {
	applying, fingerprint, code := parseReclaimArgs(args, stdout)
	if code != 0 {
		return code
	}
	if !inRepository(root) {
		fmt.Fprintln(stdout, toon.NotInRepo())
		return 1
	}
	home := benchHome()
	plan, err := planPoolReclaim(root, home)
	if err != nil {
		fmt.Fprintln(stdout, toon.Errorf("cannot read the worktree pool", "make "+poolKeysDirAt(home)+" readable and retry"))
		return 1
	}
	if !applying {
		out, err := renderPoolReclaim(plan)
		if err != nil {
			fmt.Fprintln(stdout, toon.RenderError(err))
			return 1
		}
		fmt.Fprint(stdout, out)
		return 0
	}
	// The plan just re-read the pool. A supplied fingerprint that no longer matches it
	// means the pool moved since the operator read the plan. Nothing is removed on the
	// strength of that stale reading.
	if reclaimpolicy.PlanDrift(fingerprint, plan.fingerprint) {
		out, err := renderPoolReclaimStale()
		if err != nil {
			fmt.Fprintln(stdout, toon.RenderError(err))
			return 1
		}
		fmt.Fprint(stdout, out)
		return 1
	}
	applied := applyPoolReclaim(plan, home)
	out, err := renderPoolReclaimApplied(applied, plan.fingerprint)
	if err != nil {
		fmt.Fprintln(stdout, toon.RenderError(err))
		return 1
	}
	fmt.Fprint(stdout, out)
	// The policy owns the exit verdict: a key the plan named and the apply could not
	// remove leaves the operator's intent unsatisfied, and the exit code is what a
	// script reads. The rows say which and why.
	if reclaimpolicy.ApplyIncomplete(plan.verdicts, applied) {
		return 1
	}
	return 0
}

// parseReclaimArgs accepts the bare plan and exactly one `--apply <fingerprint>`. A
// missing, empty, or malformed value is a usage refusal rather than a weaker apply. A
// parser that read an absent value as "apply everything" would turn a fumbled flag
// into a destructive run.
func parseReclaimArgs(args []string, stdout io.Writer) (applying bool, fingerprint string, code int) {
	for i := 0; i < len(args); i++ {
		if args[i] != "--apply" {
			fmt.Fprintln(stdout, toon.Usage(usage.WorktreeReclaim, args[i]))
			return false, "", 2
		}
		if applying {
			fmt.Fprintln(stdout, toon.Usage(usage.WorktreeReclaim, args[i]))
			return false, "", 2
		}
		if i+1 >= len(args) {
			fmt.Fprintln(stdout, toon.MissingArg(usage.WorktreeReclaim, "fingerprint"))
			return false, "", 2
		}
		i++
		if !wellFormedFingerprint(args[i]) {
			fmt.Fprintln(stdout, toon.Usage(usage.WorktreeReclaim, args[i]))
			return false, "", 2
		}
		applying, fingerprint = true, args[i]
	}
	return applying, fingerprint, 0
}

// applyPoolReclaim removes the keys the plan named and reports what it did per key. The
// keys the plan retained are reported untouched, so the destructive step leaves evidence
// covering the whole pool, not only the part it acted on.
func applyPoolReclaim(plan poolReclaimPlan, home string) []poolKeyVerdict {
	applied := make([]poolKeyVerdict, 0, len(plan.verdicts))
	for _, verdict := range plan.verdicts {
		if !verdict.Reclaimable() {
			applied = append(applied, reclaimpolicy.RetainedOnApply(verdict))
			continue
		}
		applied = append(applied, removePoolKey(home, verdict.Key))
	}
	return applied
}

// removePoolKey is the only place in the tree that deletes pool bytes. The caller has
// already matched the fingerprint, which speaks for the plan as a whole. The re-check
// here speaks for this one key at the instant of removal. The parent assertion bounds
// the target to a direct child of the pool.
func removePoolKey(home, key string) poolKeyVerdict {
	pool := poolKeysDirAt(home)
	target, refusal, ok := reclaimpolicy.RemovalBounds(pool, key)
	if !ok {
		return refusal
	}
	if refusal, ok := reclaimpolicy.RemovalRequalified(key, classifyPoolKey(target, key)); !ok {
		return refusal
	}
	removeErr := ""
	if err := os.RemoveAll(target); err != nil {
		removeErr = err.Error()
	}
	return reclaimpolicy.RemovalOutcome(key, removeErr)
}

// poolRows projects verdicts into the one row shape poolReclaimFields names, so the plan
// and the apply cannot drift into two shapes for one schema.
func poolRows(verdicts []poolKeyVerdict) [][]string {
	rows := make([][]string, 0, len(verdicts))
	for _, verdict := range verdicts {
		rows = append(rows, []string{verdict.Key, verdict.Verdict, verdict.Reason})
	}
	return rows
}

// renderPoolReclaimApplied projects the apply in the plan's row shape, with an aggregate
// naming what was removed rather than what could be.
func renderPoolReclaimApplied(applied []poolKeyVerdict, fingerprint string) (string, error) {
	removed := 0
	for _, verdict := range applied {
		if verdict.Verdict == poolVerdictRemoved {
			removed++
		}
	}
	table, err := toon.Table("pool_reclaim", poolReclaimFields, poolRows(applied))
	if err != nil {
		return "", err
	}
	aggregate, err := toon.TableTyped("pool_reclaim_applied", []string{"keys", "removed", "retained", "fingerprint"},
		[][]any{{len(applied), removed, len(applied) - removed, fingerprint}})
	if err != nil {
		return "", err
	}
	help, err := axi.RenderHelp(nil)
	if err != nil {
		return "", err
	}
	return table + aggregate + help, nil
}

// renderPoolReclaimStale is the drift refusal. It removes nothing and names the re-plan
// through the same action renderer the plan advertises its apply with. The command an
// operator is sent back to cannot drift from the command that printed the fingerprint.
func renderPoolReclaimStale() (string, error) {
	help, err := axi.RenderHelp([]axi.Action{axi.ExecutableInvocation("re-plan the pool and apply the fingerprint it prints",
		axi.KnownArgument("worktree"), axi.KnownArgument("reclaim"))})
	if err != nil {
		return "", err
	}
	return toon.Errorf("worktree pool reclaim plan is stale",
		"the pool changed since that plan; nothing was removed") + "\n" + help, nil
}

// renderPoolReclaim projects the plan: one row per key, an aggregate carrying the
// fingerprint, and the exact apply invocation that acts on it. The apply command is
// advertised only when there is something to apply, so a clean pool reads as an answer
// rather than an invitation.
func renderPoolReclaim(plan poolReclaimPlan) (string, error) {
	rows := poolRows(plan.verdicts)
	table, err := toon.Table("pool_reclaim", poolReclaimFields, rows)
	if err != nil {
		return "", err
	}
	reclaimable := plan.reclaimableCount()
	aggregate, err := toon.TableTyped("pool_reclaim_aggregate", []string{"keys", "reclaimable", "retained", "fingerprint"},
		[][]any{{len(plan.verdicts), reclaimable, len(plan.verdicts) - reclaimable, plan.fingerprint}})
	if err != nil {
		return "", err
	}
	var actions []axi.Action
	if reclaimable > 0 {
		actions = append(actions, axi.ExecutableInvocation("remove the planned pool keys",
			axi.KnownArgument("worktree"), axi.KnownArgument("reclaim"), axi.KnownArgument("--apply"), axi.KnownArgument(plan.fingerprint)))
	}
	help, err := axi.RenderHelp(actions)
	if err != nil {
		return "", err
	}
	return table + aggregate + help, nil
}
