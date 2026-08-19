package worktree

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
)

// newReclaimPool binds a BENCH_HOME under the test's own temporary directory, creates the
// pool parent inside it, chdirs into a fresh repository, and returns the pool path. Every
// fixture here goes through it: the package's TestMain reds on residue under the shared
// private home, so a test that reached the operator's pool would be caught, and binding
// per test is what keeps it from having to be.
func newReclaimPool(t *testing.T) (pool, root string) {
	t.Helper()
	root = newWorktreeRepo(t)
	t.Setenv("BENCH_HOME", filepath.Join(root, ".bench-home"))
	pool = poolKeysDir()
	mustMkdirAll(t, pool, 0o700)
	t.Chdir(root)
	return pool, root
}

// plantDeadChild writes one pool child whose `.git` pointer names a repository that was
// never created — the shape a deleted source repository leaves behind, and the only one
// the predicate may act on.
func plantDeadChild(t *testing.T, pool, key, child string) string {
	t.Helper()
	target := filepath.Join(t.TempDir(), "deleted-"+key+"-"+child, ".git", "worktrees", child)
	plantChild(t, pool, key, child, target)
	return target
}

// plantLiveChild writes one pool child whose `.git` pointer names a directory that exists.
func plantLiveChild(t *testing.T, pool, key, child string) {
	t.Helper()
	target := filepath.Join(t.TempDir(), "live-"+key+"-"+child)
	mustMkdirAll(t, target, 0o700)
	plantChild(t, pool, key, child, target)
}

func plantChild(t *testing.T, pool, key, child, target string) {
	t.Helper()
	dir := filepath.Join(pool, key, child)
	mustMkdirAll(t, dir, 0o700)
	mustWrite(t, filepath.Join(dir, ".git"), []byte("gitdir: "+target+"\n"), 0o644)
}

// poolListing is the sorted recursive listing a "removed nothing" assertion compares. It
// records every path's mode, so a plan that truncated a file rather than unlinking it
// would still show up.
func poolListing(t *testing.T, pool string) string {
	t.Helper()
	var lines []string
	mustNoError(t, filepath.WalkDir(pool, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, statErr := os.Lstat(path)
		if statErr != nil {
			return statErr
		}
		lines = append(lines, path+" "+info.Mode().String())
		return nil
	}))
	sort.Strings(lines)
	return strings.Join(lines, "\n")
}

// reclaimVerdicts parses the plan table into key order and, per key, its verdict and
// reason. TOON quotes any cell it must, so each field is unquoted before comparison.
func reclaimVerdicts(t *testing.T, out string) ([]string, map[string][2]string) {
	t.Helper()
	keys := make([]string, 0, 8)
	verdicts := make(map[string][2]string, 8)
	for _, row := range reclaimRows(t, out) {
		fields := strings.SplitN(row, ",", 3)
		requireTest(t, len(fields) == 3, "plan row %q does not carry three fields", row)
		key, verdict, reason := unquoteCell(fields[0]), unquoteCell(fields[1]), unquoteCell(fields[2])
		keys = append(keys, key)
		verdicts[key] = [2]string{verdict, reason}
	}
	return keys, verdicts
}

func unquoteCell(field string) string {
	trimmed := strings.TrimPrefix(strings.TrimSuffix(field, `"`), `"`)
	return strings.ReplaceAll(trimmed, `\"`, `"`)
}

// reclaimRows returns the plan table's raw data rows.
func reclaimRows(t *testing.T, out string) []string {
	t.Helper()
	var rows []string
	inTable := false
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, "pool_reclaim[") {
			inTable = true
			continue
		}
		if !strings.HasPrefix(line, "  ") {
			inTable = false
			continue
		}
		if inTable {
			rows = append(rows, strings.TrimSpace(line))
		}
	}
	return rows
}

func mustReclaim(t *testing.T, args ...string) (string, int) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := ReclaimCommand(args, &stdout, &stderr)
	requireTest(t, stderr.Len() == 0, "reclaim wrote to stderr: %q", stderr.String())
	return stdout.String(), code
}

// [PL1][PL5] The plan's whole job is to separate the keys a deleted repository left behind
// from the ones still holding work. A plan that names every key, or none, is useless; a
// plan that names the live key is the destructive bug this command must never have.
func TestReclaimCommandPlansOnlyTheProvablyDeadKeys(t *testing.T) {
	pool, _ := newReclaimPool(t)
	plantDeadChild(t, pool, "dead-key", "wt")
	mustMkdirAll(t, filepath.Join(pool, "empty-key"), 0o700)
	plantLiveChild(t, pool, "live-key", "wt")

	out, code := mustReclaim(t)
	requireTest(t, code == 0, "reclaim code=%d out=%q", code, out)
	keys, verdicts := reclaimVerdicts(t, out)
	requireTest(t, len(keys) == 3, "reclaim keys = %v, want one row per key", keys)
	for key, want := range map[string]string{"dead-key": poolVerdictReclaim, "empty-key": poolVerdictReclaim, "live-key": poolVerdictRetain} {
		requireTest(t, verdicts[key][0] == want, "key %s verdict = %q, want %q", key, verdicts[key][0], want)
	}
	requireTest(t, strings.Contains(out, "pool_reclaim_aggregate[1]{keys,reclaimable,retained,fingerprint}:") && strings.Contains(out, "\n  3,2,1,"),
		"reclaim aggregate did not count 3 keys, 2 reclaimable, 1 retained: %q", out)
}

// [PL2][SH1][SH2][SH3] A retained key the operator expected to be reclaimed has to say
// what protected it. One shared reason string for every retention would leave the operator
// unable to tell a correct retention from a predicate bug, so each protecting shape gets
// its own — including the key that mixes a live and a dead pointer, which is retained
// whole rather than half-reclaimed.
func TestReclaimCommandNamesWhatProtectedEachRetainedKey(t *testing.T) {
	pool, _ := newReclaimPool(t)
	plantLiveChild(t, pool, "live-pointer", "wt")
	mustMkdirAll(t, filepath.Join(pool, "repo-directory", "wt", ".git"), 0o700)
	mustMkdirAll(t, filepath.Join(pool, "no-gitdir", "wt"), 0o700)
	mustWrite(t, filepath.Join(pool, "no-gitdir", "wt", ".git"), []byte("ref: refs/heads/main\n"), 0o644)
	mustMkdirAll(t, filepath.Join(pool, "stray-entry"), 0o700)
	mustWrite(t, filepath.Join(pool, "stray-entry", "notes.txt"), []byte("mine\n"), 0o644)
	plantDeadChild(t, pool, "mixed", "a-dead")
	plantLiveChild(t, pool, "mixed", "b-live")
	mustMkdirAll(t, filepath.Join(pool, "no-git-entry", "wt"), 0o700)

	out, code := mustReclaim(t)
	requireTest(t, code == 0, "reclaim code=%d out=%q", code, out)
	keys, verdicts := reclaimVerdicts(t, out)
	reasons := make(map[string]string, len(keys))
	for _, key := range keys {
		verdict, reason := verdicts[key][0], verdicts[key][1]
		requireTest(t, verdict == poolVerdictRetain, "key %s verdict = %q, want retain (%s)", key, verdict, reason)
		for other, seen := range reasons {
			requireTest(t, seen != reason, "keys %s and %s share the reason %q", other, key, reason)
		}
		reasons[key] = reason
	}
	requireTest(t, len(reasons) == 6, "retained keys = %v, want all six", reasons)
	for key, want := range map[string]string{
		"live-pointer":   "gitdir: target exists",
		"repo-directory": ".git is a repository directory",
		"no-gitdir":      ".git carries no gitdir: target",
		"stray-entry":    "entry notes.txt is not a directory",
		"mixed":          "b-live gitdir: target exists",
		"no-git-entry":   "holds no .git entry",
	} {
		requireTest(t, strings.Contains(reasons[key], want), "key %s reason = %q, want it to name %q", key, reasons[key], want)
	}
}

// [SH4] Unknown is not absence. A child whose `.git` cannot be read, and a `gitdir:` target
// whose stat fails for any reason other than absence, both leave existence undecided —
// counting either as proof of deadness is the one direction that destroys work.
func TestPoolKeyPredicateRetainsAKeyWhoseExistenceItCannotProve(t *testing.T) {
	if os.Geteuid() == 0 {
		capability.Capability(t, capability.Privilege, "root bypasses directory permissions; cannot deny the stat that leaves existence unknown")
	}
	pool, _ := newReclaimPool(t)

	plantDeadChild(t, pool, "unreadable-child", "wt")
	sealed := filepath.Join(pool, "unreadable-child", "wt")
	t.Cleanup(func() { _ = os.Chmod(sealed, 0o700) })
	mustChmod(t, sealed, 0o000)

	denied := filepath.Join(t.TempDir(), "denied")
	mustMkdirAll(t, denied, 0o700)
	plantChild(t, pool, "unstattable-target", "wt", filepath.Join(denied, "gone", ".git"))
	t.Cleanup(func() { _ = os.Chmod(denied, 0o700) })
	mustChmod(t, denied, 0o000)

	for _, key := range []string{"unreadable-child", "unstattable-target"} {
		verdict := classifyPoolKey(filepath.Join(pool, key), key)
		requireTest(t, verdict.verdict == poolVerdictRetain && strings.Contains(verdict.reason, "cannot be read"),
			"key %s = %q/%q, want a retain naming the unreadable step", key, verdict.verdict, verdict.reason)
	}
}

// [SH5] A symlink is never followed. If one at the key, the child, or the `.git` were
// traversed, the pool would stop bounding what an apply can remove and bytes outside it
// would become the subject.
func TestPoolKeyPredicateRetainsSymlinksUnfollowed(t *testing.T) {
	pool, _ := newReclaimPool(t)
	empty := t.TempDir()

	mustNoError(t, os.Symlink(empty, filepath.Join(pool, "symlinked-key")))

	mustMkdirAll(t, filepath.Join(pool, "symlinked-child"), 0o700)
	mustNoError(t, os.Symlink(empty, filepath.Join(pool, "symlinked-child", "wt")))

	dead := plantDeadChild(t, pool, "symlinked-git", "wt")
	pointer := filepath.Join(pool, "symlinked-git", "wt", ".git")
	mustRemove(t, pointer)
	elsewhere := filepath.Join(t.TempDir(), "pointer")
	mustWrite(t, elsewhere, []byte("gitdir: "+dead+"\n"), 0o644)
	mustNoError(t, os.Symlink(elsewhere, pointer))

	for key, want := range map[string]string{
		"symlinked-key":   "key is a symlink",
		"symlinked-child": "entry wt is a symlink",
		"symlinked-git":   "child wt .git is a symlink",
	} {
		verdict := classifyPoolKey(filepath.Join(pool, key), key)
		requireTest(t, verdict.verdict == poolVerdictRetain && verdict.reason == want,
			"key %s = %q/%q, want a retain reading %q", key, verdict.verdict, verdict.reason, want)
	}
}

// [SH6] A session that has acquired its pool directory but not yet checked anything out
// holds an empty key, which the empty-key clause would otherwise take out from under it.
// The current repository's key is excluded before the predicate ever runs.
func TestReclaimCommandRetainsTheCurrentRepositorysEmptyKey(t *testing.T) {
	pool, root := newReclaimPool(t)
	current := filepath.Base(Pool(canonicalRoot(root)))
	mustMkdirAll(t, filepath.Join(pool, current), 0o700)

	out, code := mustReclaim(t)
	requireTest(t, code == 0, "reclaim code=%d out=%q", code, out)
	keys, verdicts := reclaimVerdicts(t, out)
	requireTest(t, len(keys) == 1 && keys[0] == current && verdicts[current][0] == poolVerdictRetain,
		"reclaim rows = %v/%v, want the current repository's empty key retained", keys, verdicts)
}

// [PL4] "Nothing to do" is an answer. A clean pool, and a home that never leased a
// worktree at all, both have to produce a definitive zero-row table and exit zero, because
// silence or an error there turns a successful absence into a failure an operator chases.
func TestReclaimCommandAnswersAnEmptyPoolWithZeroRows(t *testing.T) {
	t.Run("empty pool", func(t *testing.T) {
		newReclaimPool(t)
		out, code := mustReclaim(t)
		requireTest(t, code == 0 && strings.Contains(out, "pool_reclaim[0]{key,verdict,reason}:"),
			"empty pool reclaim code=%d out=%q", code, out)
	})
	t.Run("absent pool directory", func(t *testing.T) {
		pool, _ := newReclaimPool(t)
		mustRemove(t, pool)
		out, code := mustReclaim(t)
		requireTest(t, code == 0 && strings.Contains(out, "pool_reclaim[0]{key,verdict,reason}:"),
			"absent pool reclaim code=%d out=%q", code, out)
	})
}

// [PL5] The bare plan is an inspection. Removing anything is the worst failure this command
// has, so the pool's whole recursive listing must survive a plan over every shape unchanged.
func TestReclaimCommandRemovesNothing(t *testing.T) {
	pool, _ := newReclaimPool(t)
	plantDeadChild(t, pool, "dead-key", "wt")
	mustMkdirAll(t, filepath.Join(pool, "empty-key"), 0o700)
	plantLiveChild(t, pool, "live-key", "wt")
	mustWrite(t, filepath.Join(pool, "stray-key-file"), []byte("x\n"), 0o644)
	before := poolListing(t, pool)

	out, code := mustReclaim(t)
	requireTest(t, code == 0, "reclaim code=%d out=%q", code, out)
	requireTest(t, poolListing(t, pool) == before, "the pool changed across a bare plan:\nbefore\n%s\nafter\n%s", before, poolListing(t, pool))
}

// [PL3] Acting on the plan must need no invention: the printed apply invocation carries the
// fingerprint the apply will demand, and that value is the one the plan itself derived.
func TestReclaimCommandPrintsTheApplyInvocationCarryingTheFingerprint(t *testing.T) {
	pool, root := newReclaimPool(t)
	plantDeadChild(t, pool, "dead-key", "wt")

	out, code := mustReclaim(t)
	requireTest(t, code == 0, "reclaim code=%d out=%q", code, out)
	plan, err := planPoolReclaim(root)
	mustNoError(t, err)
	requireTest(t, len(plan.fingerprint) == 64, "plan fingerprint = %q, want a sha256 digest", plan.fingerprint)
	requireTest(t, strings.Contains(out, "bench worktree reclaim --apply "+plan.fingerprint),
		"plan did not print the apply invocation carrying %s: %q", plan.fingerprint, out)
}

// The exclusion needs a current repository to name, so the command requires one like every
// other `bench worktree` subcommand. Refusing with the shared structured error is what lets
// a caller tell "you are in the wrong directory" from "your pool is clean".
func TestReclaimCommandRefusesOutsideARepository(t *testing.T) {
	t.Setenv("BENCH_HOME", filepath.Join(t.TempDir(), "home"))
	t.Chdir(t.TempDir())
	out, code := mustReclaim(t)
	requireTest(t, code == 1 && strings.Contains(out, "not in a git repository"),
		"reclaim outside a repository code=%d out=%q", code, out)
}
