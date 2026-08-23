package worktree

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/usage"
)

// newReclaimPool binds a BENCH_HOME under the test's own temporary directory and creates
// the pool parent inside it.
// It chdirs into a fresh repository and returns the pool path.
// Every fixture here goes through it.
// The package's TestMain reds on residue under the shared private home.
// A test that reached the operator's pool would be caught.
// Binding per test is what keeps it from having to be.
func newReclaimPool(t *testing.T) (pool, root string) {
	t.Helper()
	root = newWorktreeRepo(t)
	home := filepath.Join(root, ".bench-home")
	bindEnv(t, "BENCH_HOME", home)
	pool = poolKeysDirAt(home)
	mustMkdirAll(t, pool, 0o700)
	chdir(t, root)
	return pool, root
}

// plantDeadChild writes one pool child whose `.git` pointer names a repository that was
// never created.
// That is the shape a deleted source repository leaves behind, and the only one the
// predicate may act on.
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
	return reclaimTableRows(t, out, "pool_reclaim[")
}

// reclaimTableRows returns the raw data rows indented under the named table header. Every
// table this package reads is located through here, so "rows live indented under their
// header" is stated once.
func reclaimTableRows(t *testing.T, out, header string) []string {
	t.Helper()
	var rows []string
	inTable := false
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, header) {
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

// reclaimAggregate parses the single row under an aggregate header into its four unquoted
// counts-and-fingerprint fields. The assertion reads the parsed fields rather than a
// substring, which keeps it independent of whether TOON quoted a cell.
// A fingerprint that reads as a number arrives quoted.
// A substring match on the bare digits would red on a counting failure that never happened.
func reclaimAggregate(t *testing.T, out, header string) [4]string {
	t.Helper()
	rows := reclaimTableRows(t, out, header)
	requireTest(t, len(rows) == 1, "want exactly one %s row, got %v: %q", header, rows, out)
	fields := strings.SplitN(rows[0], ",", 4)
	requireTest(t, len(fields) == 4, "aggregate row %q does not carry four fields: %q", rows[0], out)
	var parsed [4]string
	for i, field := range fields {
		parsed[i] = unquoteCell(field)
	}
	return parsed
}

// requireReclaimAggregate asserts the aggregate row's counts and fingerprint.
func requireReclaimAggregate(t *testing.T, out, header, keys, first, retained, fingerprint string) {
	t.Helper()
	got := reclaimAggregate(t, out, header)
	want := [4]string{keys, first, retained, fingerprint}
	requireTest(t, got == want, "%s row = %v, want %v: %q", header, got, want, out)
}

func mustReclaim(t *testing.T, args ...string) (string, int) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := ReclaimCommand(args, &stdout, &stderr)
	requireTest(t, stderr.Len() == 0, "reclaim wrote to stderr: %q", stderr.String())
	return stdout.String(), code
}

// [PL1][PL5] The plan's whole job is to separate the keys a deleted repository left behind
// from the ones still holding work. A plan that names every key, or none, is useless.
// A plan that names the live key is the destructive bug this command must never have.
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
	requireTest(t, strings.Contains(out, "pool_reclaim_aggregate[1]{keys,reclaimable,retained,fingerprint}:"), "plan printed no aggregate header: %q", out)
	requireReclaimAggregate(t, out, "pool_reclaim_aggregate[", "3", "2", "1", reclaimFingerprint(t, out))
}

// [PL2][SH1][SH2][SH3] A retained key the operator expected to be reclaimed has to say
// what protected it. One shared reason string for every retention would leave the operator
// unable to tell a correct retention from a predicate bug.
// Each protecting shape therefore gets its own reason, including the key that mixes a live
// and a dead pointer.
// That key is retained whole rather than half-reclaimed.
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

// [SH6] A session that has acquired its pool directory but not yet checked anything out
// holds an empty key.
// The empty-key clause would otherwise take that key out from under it.
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
	markProof(t, "reclaim/journey/registration")
}

// [PL4] "Nothing to do" is an answer. A clean pool, and a home that never leased a
// worktree at all, both have to produce a definitive zero-row table and exit zero.
// Silence or an error there would turn a successful absence into a failure an operator chases.
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

// [PL3] Acting on the plan must need no invention.
// The printed apply invocation carries the fingerprint the apply will demand, and that
// value is the one the plan itself derived.
func TestReclaimCommandPrintsTheApplyInvocationCarryingTheFingerprint(t *testing.T) {
	pool, root := newReclaimPool(t)
	plantDeadChild(t, pool, "dead-key", "wt")

	out, code := mustReclaim(t)
	requireTest(t, code == 0, "reclaim code=%d out=%q", code, out)
	plan, err := planPoolReclaim(root, filepath.Dir(pool))
	mustNoError(t, err)
	requireTest(t, len(plan.fingerprint) == 64, "plan fingerprint = %q, want a sha256 digest", plan.fingerprint)
	requireTest(t, strings.Contains(out, "bench worktree reclaim --apply "+plan.fingerprint),
		"plan did not print the apply invocation carrying %s: %q", plan.fingerprint, out)
}

// The exclusion needs a current repository to name, so the command requires one like every
// other `bench worktree` subcommand. Refusing with the shared structured error is what lets
// a caller tell "you are in the wrong directory" from "your pool is clean".
func TestReclaimCommandRefusesOutsideARepository(t *testing.T) {
	bindEnv(t, "BENCH_HOME", filepath.Join(t.TempDir(), "home"))
	chdir(t, t.TempDir())
	out, code := mustReclaim(t)
	requireTest(t, code == 1 && strings.Contains(out, "not in a git repository"),
		"reclaim outside a repository code=%d out=%q", code, out)
}

// reclaimFingerprint reads the fingerprint out of the apply invocation the plan printed.
// Feeding the apply the value an operator would copy is what makes the handshake a tested
// round trip rather than two independently asserted strings.
func reclaimFingerprint(t *testing.T, out string) string {
	t.Helper()
	match := regexp.MustCompile(`bench worktree reclaim --apply ([0-9a-f]{64})`).FindStringSubmatch(out)
	requireTest(t, match != nil, "plan printed no apply invocation: %q", out)
	return match[1]
}

// [AP1][PL3][SH7] The apply is the whole point of the plan and the only destructive step
// this command has. It must remove exactly what the plan named, leave every other key
// present, and count what it did.
// It must also aim every removal at a direct child of the pool.
func TestReclaimApplyRemovesExactlyThePlannedKeys(t *testing.T) {
	pool, _ := newReclaimPool(t)
	plantDeadChild(t, pool, "dead-key", "wt")
	mustMkdirAll(t, filepath.Join(pool, "empty-key"), 0o700)
	plantLiveChild(t, pool, "live-key", "wt")
	liveListing := poolListing(t, filepath.Join(pool, "live-key"))

	plan, planCode := mustReclaim(t)
	requireTest(t, planCode == 0, "plan code=%d out=%q", planCode, plan)
	fingerprint := reclaimFingerprint(t, plan)

	out, code := mustReclaim(t, "--apply", fingerprint)
	requireTest(t, code == 0, "apply code=%d out=%q", code, out)
	keys, verdicts := reclaimVerdicts(t, out)
	requireTest(t, len(keys) == 3, "apply keys = %v, want one row per key", keys)
	for key, want := range map[string]string{"dead-key": poolVerdictRemoved, "empty-key": poolVerdictRemoved, "live-key": poolVerdictRetained} {
		requireTest(t, verdicts[key][0] == want, "key %s verdict = %q, want %q (%s)", key, verdicts[key][0], want, verdicts[key][1])
	}
	requireTest(t, strings.Contains(out, "pool_reclaim_applied[1]{keys,removed,retained,fingerprint}:"), "apply printed no applied aggregate header: %q", out)
	requireReclaimAggregate(t, out, "pool_reclaim_applied[", "3", "2", "1", fingerprint)
	for _, key := range []string{"dead-key", "empty-key"} {
		target := filepath.Join(pool, key)
		requireTest(t, filepath.Dir(target) == pool, "removed %s whose parent is not the pool", target)
		_, err := os.Lstat(target)
		requireTest(t, os.IsNotExist(err), "planned key %s survived the apply: %v", key, err)
	}
	requireTest(t, poolListing(t, filepath.Join(pool, "live-key")) == liveListing,
		"the retained key changed across the apply:\nbefore\n%s\nafter\n%s", liveListing, poolListing(t, filepath.Join(pool, "live-key")))
	markProof(t, "reclaim/journey/deletion")
}

// The writer quotes an aggregate cell that would otherwise read as a number, so a
// fingerprint of all digits arrives quoted. Which spelling a run gets is luck: an
// all-digit digest is roughly one in a thousand.
// Both are therefore exercised here directly rather than waited for.
func TestReclaimAggregateReadsAFingerprintInEitherSpelling(t *testing.T) {
	digits := strings.Repeat("0123456789", 6) + "0123"
	hex := strings.Repeat("0123456789abcdef", 4)
	for name, out := range map[string]string{
		"quoted":   "pool_reclaim_applied[1]{keys,removed,retained,fingerprint}:\n  3,2,1,\"" + digits + "\"\n",
		"unquoted": "pool_reclaim_applied[1]{keys,removed,retained,fingerprint}:\n  3,2,1," + hex + "\n",
	} {
		fingerprint := hex
		if name == "quoted" {
			fingerprint = digits
		}
		t.Run(name, func(t *testing.T) {
			requireReclaimAggregate(t, out, "pool_reclaim_applied[", "3", "2", "1", fingerprint)
		})
	}
}

// [AP2] A fingerprint the pool no longer matches is a reading that has stopped being true.
// Removing on the strength of it is the failure the handshake exists to prevent, so the
// apply refuses, names the re-plan, and touches nothing.
func TestReclaimApplyRefusesAFingerprintThePoolNoLongerMatches(t *testing.T) {
	pool, _ := newReclaimPool(t)
	plantDeadChild(t, pool, "dead-key", "wt")

	plan, planCode := mustReclaim(t)
	requireTest(t, planCode == 0, "plan code=%d out=%q", planCode, plan)
	fingerprint := reclaimFingerprint(t, plan)

	plantDeadChild(t, pool, "arrived-later", "wt")
	before := poolListing(t, pool)

	out, code := mustReclaim(t, "--apply", fingerprint)
	requireTest(t, code == 1, "stale apply code=%d out=%q", code, out)
	requireTest(t, strings.Contains(out, "worktree pool reclaim plan is stale") && strings.Contains(out, "bench worktree reclaim,"),
		"stale refusal did not name the re-plan command: %q", out)
	requireTest(t, poolListing(t, pool) == before, "a stale apply changed the pool:\nbefore\n%s\nafter\n%s", before, poolListing(t, pool))
}

// [AP3] The fingerprint proves the plan as a whole is current; it cannot see a change made
// inside a key it already counted. Only the re-check immediately before the RemoveAll can.
// That is why a key that goes live in that window survives while its neighbours are still
// removed. The seam is the apply itself, because a plan whose fingerprint still matches is
// exactly the state the fingerprint cannot distinguish.
func TestReclaimApplyRetainsAKeyThatWentLiveAfterThePlan(t *testing.T) {
	pool, root := newReclaimPool(t)
	plantDeadChild(t, pool, "still-dead", "wt")
	mustMkdirAll(t, filepath.Join(pool, "went-live"), 0o700)

	plan, err := planPoolReclaim(root, filepath.Dir(pool))
	mustNoError(t, err)
	requireTest(t, plan.reclaimableCount() == 2, "plan = %#v, want both keys reclaimable", plan.verdicts)

	plantLiveChild(t, pool, "went-live", "wt")

	applied := applyPoolReclaim(plan, filepath.Dir(pool))
	verdicts := make(map[string]poolKeyVerdict, len(applied))
	for _, verdict := range applied {
		verdicts[verdict.Key] = verdict
	}
	requireTest(t, verdicts["went-live"].Verdict == poolVerdictRetained && strings.Contains(verdicts["went-live"].Reason, "stopped qualifying"),
		"went-live = %q/%q, want a retain naming the re-check", verdicts["went-live"].Verdict, verdicts["went-live"].Reason)
	_, statErr := os.Lstat(filepath.Join(pool, "went-live", "wt", ".git"))
	requireTest(t, statErr == nil, "the key that went live was removed anyway: %v", statErr)
	requireTest(t, verdicts["still-dead"].Verdict == poolVerdictRemoved, "still-dead = %q/%q, want it removed", verdicts["still-dead"].Verdict, verdicts["still-dead"].Reason)
	_, statErr = os.Lstat(filepath.Join(pool, "still-dead"))
	requireTest(t, os.IsNotExist(statErr), "still-dead survived: %v", statErr)
	markProof(t, "reclaim/journey/process-liveness")
}

// [AP4] A repeat invocation has to be safe to run. An apply over a pool holding nothing
// reclaimable is a successful no-op, not an error an operator has to interpret.
func TestReclaimApplyOverNothingToReclaimIsASuccessfulNoOp(t *testing.T) {
	pool, root := newReclaimPool(t)
	plantLiveChild(t, pool, "live-key", "wt")
	before := poolListing(t, pool)

	plan, err := planPoolReclaim(root, filepath.Dir(pool))
	mustNoError(t, err)
	requireTest(t, plan.reclaimableCount() == 0, "plan = %#v, want nothing reclaimable", plan.verdicts)

	out, code := mustReclaim(t, "--apply", plan.fingerprint)
	requireTest(t, code == 0, "no-op apply code=%d out=%q", code, out)
	requireReclaimAggregate(t, out, "pool_reclaim_applied[", "1", "0", "1", plan.fingerprint)
	requireTest(t, poolListing(t, pool) == before, "a no-op apply changed the pool:\nbefore\n%s\nafter\n%s", before, poolListing(t, pool))
}

// [AP5] A fumbled flag must never become a destructive run. An absent, empty, or malformed
// value — and a repeated flag — are all usage refusals that read nothing and remove nothing.
func TestReclaimApplyRefusesAFumbledFingerprint(t *testing.T) {
	pool, root := newReclaimPool(t)
	plantDeadChild(t, pool, "dead-key", "wt")
	plan, err := planPoolReclaim(root, filepath.Dir(pool))
	mustNoError(t, err)
	before := poolListing(t, pool)

	for _, args := range [][]string{
		{"--apply"},
		{"--apply", ""},
		{"--apply", strings.Repeat("a", 63)},
		{"--apply", strings.Repeat("a", 65)},
		{"--apply", strings.Repeat("g", 64)},
		{"--apply", strings.ToUpper(plan.fingerprint)},
		{"--apply", plan.fingerprint, "--apply", plan.fingerprint},
		{"--apply", plan.fingerprint, "extra"},
		{"--force"},
	} {
		out, code := mustReclaim(t, args...)
		requireTest(t, code == 2 && strings.Contains(out, usage.WorktreeReclaim),
			"args=%q code=%d out=%q, want a usage refusal", args, code, out)
		requireTest(t, poolListing(t, pool) == before, "args=%q changed the pool:\nbefore\n%s\nafter\n%s", args, before, poolListing(t, pool))
	}
}

// [SH7] The pool is what bounds this command's blast radius. The removal asserts its
// target's parent rather than trusting the key name it was handed. A name that is not a
// plain direct child leaves the bytes it points at alone.
func TestRemovePoolKeyRefusesATargetOutsideThePool(t *testing.T) {
	pool, _ := newReclaimPool(t)
	decoy := filepath.Join(filepath.Dir(pool), "decoy")
	mustMkdirAll(t, decoy, 0o700)
	mustWrite(t, filepath.Join(decoy, "keep.txt"), []byte("keep\n"), 0o644)

	for _, key := range []string{"", ".", "..", filepath.Join("..", "decoy"), "nested/child", decoy} {
		verdict := removePoolKey(filepath.Dir(pool), key)
		requireTest(t, verdict.Verdict == poolVerdictRetained && strings.Contains(verdict.Reason, "not a direct child of "+pool),
			"key %q = %q/%q, want a retain naming the pool boundary", key, verdict.Verdict, verdict.Reason)
	}
	_, err := os.Lstat(filepath.Join(decoy, "keep.txt"))
	requireTest(t, err == nil, "a removal escaped the pool: %v", err)
	_, err = os.Lstat(pool)
	requireTest(t, err == nil, "the pool itself was removed: %v", err)
}

// [RP1] This is the bug this spec exists for, end to end and through the real command.
// A repository that genuinely existed, genuinely leased a pool key, and was then deleted
// leaves a key no repository-anchored path can reach.
// It must be planned and then actually removed.
func TestReclaimReclaimsTheKeyOfADeletedRepository(t *testing.T) {
	pool, _ := newReclaimPool(t)
	source := newWorktreeRepo(t)
	key := filepath.Base(Pool(canonicalRoot(source)))
	created := mustCreate(t, source, "reclaim-repro", "repro")
	requireTest(t, filepath.Dir(created.Path) == filepath.Join(pool, key), "created worktree %q is not under the source's pool key %q", created.Path, key)
	mustNoError(t, os.RemoveAll(source))

	plan, planCode := mustReclaim(t)
	requireTest(t, planCode == 0, "plan code=%d out=%q", planCode, plan)
	_, verdicts := reclaimVerdicts(t, plan)
	requireTest(t, verdicts[key][0] == poolVerdictReclaim, "deleted repository key %s = %q/%q, want reclaim", key, verdicts[key][0], verdicts[key][1])

	out, code := mustReclaim(t, "--apply", reclaimFingerprint(t, plan))
	requireTest(t, code == 0, "apply code=%d out=%q", code, out)
	_, err := os.Lstat(filepath.Join(pool, key))
	requireTest(t, os.IsNotExist(err), "the deleted repository's key survived the apply: %v", err)
	markProof(t, "reclaim/journey/lease")
}

// [S4] A key the plan named and the apply could not remove leaves the operator's intent
// unsatisfied. The rows say which and why, but a script reads the exit code, so a partial
// apply that exited 0 would be read as complete. A write-denied pool fails the RemoveAll
// without changing any key's classification, so the plan's fingerprint still matches and
// the drift refusal does not fire first.
func TestReclaimApplyExitsNonZeroWhenAPlannedKeySurvives(t *testing.T) {
	pool, _ := newReclaimPool(t)
	plantDeadChild(t, pool, "really-dead", "wt")

	plan, planCode := mustReclaim(t)
	requireTest(t, planCode == 0, "plan code=%d out=%q", planCode, plan)
	fingerprint := reclaimFingerprint(t, plan)

	mustNoError(t, os.Chmod(pool, 0o500))
	t.Cleanup(func() { _ = os.Chmod(pool, 0o700) })

	out, code := mustReclaim(t, "--apply", fingerprint)
	requireTest(t, !strings.Contains(out, "stale"), "the drift refusal fired; this test must reach the apply: %q", out)
	requireTest(t, code == 1, "partial apply code=%d out=%q, want 1", code, out)
	requireTest(t, strings.Contains(out, "removal failed"), "out=%q, want a row naming the failed removal", out)
	_, err := os.Lstat(filepath.Join(pool, "really-dead"))
	requireTest(t, err == nil, "the key was removed despite a write-denied pool: %v", err)
}
