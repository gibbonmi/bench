package worktree

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// poolReclaimFingerprintVersion binds the fingerprint to this plan's derivation. A later
// change to what the plan carries changes the digest rather than silently reusing a value
// an apply would still accept.
const poolReclaimFingerprintVersion = "bench-pool-reclaim/v1"

// poolReclaimFields is the plan's row schema: which key, what the predicate decided, and
// the one fact that decided it.
var poolReclaimFields = []string{"key", "verdict", "reason"}

const (
	poolVerdictReclaim = "reclaim"
	poolVerdictRetain  = "retain"
)

// poolKeyVerdict is one key's classification. targets holds the gitdir: pointers that
// proved the key dead, in child order; it is empty for every retained key and for an
// empty one, and it feeds the fingerprint so a key whose children changed identity
// invalidates a plan that counted them.
type poolKeyVerdict struct {
	key     string
	verdict string
	reason  string
	targets []string
}

func (v poolKeyVerdict) reclaimable() bool { return v.verdict == poolVerdictReclaim }

// poolReclaimPlan is one reading of the pool: every key it holds, in name order, and the
// fingerprint over the reclaimable subset.
type poolReclaimPlan struct {
	verdicts    []poolKeyVerdict
	fingerprint string
}

func (p poolReclaimPlan) reclaimableCount() int {
	count := 0
	for _, verdict := range p.verdicts {
		if verdict.reclaimable() {
			count++
		}
	}
	return count
}

// poolKeysDir is the pool parent this command reads. Nothing else in the tree enumerates
// it: every other reclamation path is anchored at a repository root, which is exactly why
// a key whose repository was deleted is unreachable from all of them.
func poolKeysDir() string { return filepath.Join(benchHome(), "worktrees") }

// planPoolReclaim classifies every key under the pool parent. An absent pool parent is
// the zero-row answer rather than an error — a home that has never leased a worktree is a
// clean pool, not a broken one.
func planPoolReclaim(root string) (poolReclaimPlan, error) {
	current := filepath.Base(Pool(canonicalRoot(root)))
	pool := poolKeysDir()
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
		// The current repository's key is excluded before the predicate runs. A session
		// between acquiring its pool directory and its first checkout holds an empty key,
		// which the empty-key clause would otherwise take out from under it.
		if name == current {
			plan.verdicts = append(plan.verdicts, poolKeyVerdict{key: name, verdict: poolVerdictRetain, reason: "key belongs to the current repository"})
			continue
		}
		plan.verdicts = append(plan.verdicts, classifyPoolKey(filepath.Join(pool, name), name))
	}
	plan.fingerprint = fingerprintPoolReclaim(plan.verdicts)
	return plan, nil
}

// classifyPoolKey is the one reclaimability predicate in the tree. A key is reclaimable
// when it holds nothing at top level, or when every top-level entry is a real directory
// holding a regular `.git` file whose `gitdir:` target is provably absent. Everything else
// retains and says what protected it.
//
// Absence is proven only by os.IsNotExist. Every other error leaves existence unknown, and
// unknown retains — treating a permission failure as absence is the one direction that
// destroys work. Lstat throughout: a symlink where a key, a child, or a `.git` belongs is
// retained unfollowed, so the pool keeps bounding what the apply can ever remove.
func classifyPoolKey(path, name string) poolKeyVerdict {
	retain := func(format string, args ...any) poolKeyVerdict {
		return poolKeyVerdict{key: name, verdict: poolVerdictRetain, reason: fmt.Sprintf(format, args...)}
	}
	info, err := os.Lstat(path)
	switch {
	case err != nil:
		return retain("key cannot be read: %v", err)
	case info.Mode()&os.ModeSymlink != 0:
		return retain("key is a symlink")
	case !info.IsDir():
		return retain("key is not a directory")
	}
	children, err := os.ReadDir(path)
	if err != nil {
		return retain("key contents cannot be listed: %v", err)
	}
	if len(children) == 0 {
		return poolKeyVerdict{key: name, verdict: poolVerdictReclaim, reason: "key holds nothing"}
	}
	targets := make([]string, 0, len(children))
	for _, child := range children {
		target, verdict := classifyPoolChild(filepath.Join(path, child.Name()), child.Name())
		if verdict != "" {
			return retain("%s", verdict)
		}
		targets = append(targets, target)
	}
	return poolKeyVerdict{key: name, verdict: poolVerdictReclaim, reason: "every child points at an absent repository", targets: targets}
}

// classifyPoolChild decides one top-level entry. It returns the child's dangling gitdir:
// target, or the reason that entry protects the whole key — a key mixing one live and one
// dead pointer is retained whole, because partial reclamation would amputate the live half.
func classifyPoolChild(path, name string) (target, retain string) {
	info, err := os.Lstat(path)
	switch {
	case err != nil:
		return "", fmt.Sprintf("entry %s cannot be read: %v", name, err)
	case info.Mode()&os.ModeSymlink != 0:
		return "", fmt.Sprintf("entry %s is a symlink", name)
	case !info.IsDir():
		return "", fmt.Sprintf("entry %s is not a directory", name)
	}
	pointer := filepath.Join(path, ".git")
	pointerInfo, err := os.Lstat(pointer)
	switch {
	case err != nil && os.IsNotExist(err):
		return "", fmt.Sprintf("child %s holds no .git entry", name)
	case err != nil:
		return "", fmt.Sprintf("child %s .git cannot be read: %v", name, err)
	case pointerInfo.Mode()&os.ModeSymlink != 0:
		return "", fmt.Sprintf("child %s .git is a symlink", name)
	case pointerInfo.IsDir():
		return "", fmt.Sprintf("child %s .git is a repository directory", name)
	case !pointerInfo.Mode().IsRegular():
		return "", fmt.Sprintf("child %s .git is not a regular file", name)
	}
	body, err := os.ReadFile(pointer)
	if err != nil {
		return "", fmt.Sprintf("child %s .git cannot be read: %v", name, err)
	}
	target, ok := gitdirTarget(string(body))
	if !ok {
		return "", fmt.Sprintf("child %s .git carries no gitdir: target", name)
	}
	if _, err := os.Lstat(target); err == nil {
		return "", fmt.Sprintf("child %s gitdir: target exists", name)
	} else if !os.IsNotExist(err) {
		return "", fmt.Sprintf("child %s gitdir: target cannot be read: %v", name, err)
	}
	return target, ""
}

// gitdirTarget reads the pointer a git worktree's `.git` file carries. A file with no
// `gitdir:` line, or one whose value is blank, reports no target: an unparseable pointer
// is never proof that anything is absent.
func gitdirTarget(body string) (string, bool) {
	for line := range strings.SplitSeq(body, "\n") {
		rest, found := strings.CutPrefix(strings.TrimSpace(line), "gitdir:")
		if !found {
			continue
		}
		if target := strings.TrimSpace(rest); target != "" {
			return target, true
		}
	}
	return "", false
}

// fingerprintPoolReclaim digests exactly what an apply would remove: the reclaimable key
// names in pool order and each one's child gitdir: targets. A change elsewhere in the pool
// leaves it alone, so an operator is not sent back to re-plan by a key the plan did not
// name; a change to a target inside a counted key does move it.
func fingerprintPoolReclaim(verdicts []poolKeyVerdict) string {
	parts := [][]byte{[]byte(poolReclaimFingerprintVersion)}
	for _, verdict := range verdicts {
		if !verdict.reclaimable() {
			continue
		}
		parts = append(parts, []byte("key"), []byte(verdict.key), []byte(strconv.Itoa(len(verdict.targets))))
		for _, target := range verdict.targets {
			parts = append(parts, []byte("target"), []byte(target))
		}
	}
	return fingerprintParts(parts...)
}

// ReclaimCommand plans the reclaimable keys in `$BENCH_HOME/worktrees` and removes
// nothing. It runs inside a repository like every other `bench worktree` subcommand,
// because the current repository's key is the one key it may never name.
func ReclaimCommand(args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 {
		fmt.Fprintln(stdout, toon.Usage(usage.WorktreeReclaim, args[0]))
		return 2
	}
	root, err := git.Root()
	if err != nil {
		fmt.Fprintln(stdout, toon.NotInRepo())
		return 1
	}
	plan, err := planPoolReclaim(root)
	if err != nil {
		fmt.Fprintln(stdout, toon.Errorf("cannot read the worktree pool", "make "+poolKeysDir()+" readable and retry"))
		return 1
	}
	out, err := renderPoolReclaim(plan)
	if err != nil {
		fmt.Fprintln(stdout, toon.RenderError(err))
		return 1
	}
	fmt.Fprint(stdout, out)
	return 0
}

// renderPoolReclaim projects the plan: one row per key, an aggregate carrying the
// fingerprint, and the exact apply invocation that acts on it. The apply command is
// advertised only when there is something to apply, so a clean pool reads as an answer
// rather than an invitation.
func renderPoolReclaim(plan poolReclaimPlan) (string, error) {
	rows := make([][]string, 0, len(plan.verdicts))
	for _, verdict := range plan.verdicts {
		rows = append(rows, []string{verdict.key, verdict.verdict, verdict.reason})
	}
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
