package gate

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gocache"
)

// cacheClosureFixture is an outcome fixture whose manifest declares HOME, which is the
// shape every Bench project's own manifest carries. The closure derives its build cache
// entry from that declared value.
func cacheClosureFixture(t *testing.T) string {
	t.Helper()
	root := outcomeFixture(t)
	outcomeWrite(t, root, ".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":["HOME"],"paths":[],"tools":[]}`+"\n", 0o644)
	outcomeCommit(t, root, "declare HOME")
	return root
}

// C06: the closed oracle env carries the entry derived from the closure's own HOME.
func TestClosedOracleEnvCarriesTheBenchBuildCache(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := cacheClosureFixture(t)

	s, err := buildSubject(root)
	if err != nil {
		t.Fatal(err)
	}
	if !s.Closed {
		t.Fatalf("subject open: %s", s.Reason)
	}
	want, err := gocache.Dir([]string{"HOME=" + home})
	if err != nil {
		t.Fatal(err)
	}
	entries := []string{}
	for _, entry := range s.Env {
		if strings.HasPrefix(entry, gocache.Env+"=") {
			entries = append(entries, entry)
		}
	}
	if len(entries) != 1 || entries[0] != gocache.Env+"="+want {
		t.Fatalf("closure cache entries = %#v, want exactly %s=%s", entries, gocache.Env, want)
	}
}

// C07: an ambient XDG_CACHE_HOME is not a closure input, so it moves neither the entry
// nor the identity. A derivation that read it would make an undeclared variable steer
// the gate.
func TestAmbientXDGCacheHomeMovesNeitherEntryNorIdentity(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	root := cacheClosureFixture(t)

	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	first, err := buildSubject(root)
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	second, err := buildSubject(root)
	if err != nil {
		t.Fatal(err)
	}
	if first.Oracle != second.Oracle {
		t.Fatalf("oracle = %s and %s, want one identity", first.Oracle, second.Oracle)
	}
	if strings.Join(first.Env, "\x00") != strings.Join(second.Env, "\x00") {
		t.Fatalf("closure env = %#v and %#v, want one env", first.Env, second.Env)
	}
	want := gocache.Env + "=" + home + "/.cache/bench/go-build"
	if !containsEntry(first.Env, want) {
		t.Fatalf("closure env = %#v, want %s", first.Env, want)
	}
}

// C08: the phase-child base environment carries the Bench entry and no other, so an
// ambient GOCACHE never survives into a gate child.
func TestGateEnvCarriesOnlyTheBenchBuildCache(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv(gocache.Env, "/ambient/cache")

	base, err := gateEnv()
	if err != nil {
		t.Fatal(err)
	}
	want, err := gocache.Dir([]string{"HOME=" + home})
	if err != nil {
		t.Fatal(err)
	}
	entries := []string{}
	for _, entry := range base {
		if strings.HasPrefix(entry, gocache.Env+"=") {
			entries = append(entries, entry)
		}
	}
	if len(entries) != 1 || entries[0] != gocache.Env+"="+want {
		t.Fatalf("gateEnv cache entries = %#v, want exactly %s=%s", entries, gocache.Env, want)
	}
}

// C11: with no absolute HOME the phase reds before the child starts and names HOME.
func TestPhaseRedsWithoutAnAbsoluteHome(t *testing.T) {
	t.Setenv("HOME", "relative/home")
	root := t.TempDir()
	var stdout, stderr bytes.Buffer

	result := runPhase(context.Background(), root, Phase{Name: "test", Argv: []string{"true"}}, &stdout, &stderr)

	if result.Code == 0 || result.Skipped {
		t.Fatalf("phase result = %#v, want a red", result)
	}
	if !strings.Contains(stderr.String(), "HOME") {
		t.Fatalf("stderr = %q, want it to name HOME", stderr.String())
	}
}

// unwritableCacheHome answers a HOME whose derived build cache directory exists and denies
// a write, which is the state that fails a holder's lock-file open. The directory is the
// derivation's own answer rather than a second spelling of it, and its mode is restored
// before the temporary tree is removed.
func unwritableCacheHome(t *testing.T) string {
	t.Helper()
	requireDirectoryWriteDenied(t)
	home := t.TempDir()
	dir, err := gocache.Dir([]string{"HOME=" + home})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	restoreMode(t, dir)
	if err := os.Chmod(dir, 0o500); err != nil {
		t.Fatal(err)
	}
	return home
}

// plannedEvaluation answers one fixed subject at every stage, so a row drives the run
// transaction to a named step without staging a tree that must hash to that subject.
type plannedEvaluation struct{ plan subject }

func (e plannedEvaluation) acceptPre() (subject, error)   { return e.plan, nil }
func (e plannedEvaluation) validatePre() (subject, error) { return e.plan, nil }
func (e plannedEvaluation) capturePost() (subject, error) { return e.plan, nil }

// LQ17: an unwritable cache directory refuses the gate transaction, and the refusal names
// the hold error and the path. A run that graded a tree with the lock unheld would compile
// beside a clean that is removing its archives.
func TestGateTransactionRefusesAnUnheldCache(t *testing.T) {
	home := unwritableCacheHome(t)
	root := outcomeFixture(t)
	plan := subject{
		Tree:       strings.Repeat("a", 40),
		Resolution: Resolution{Kind: GateSh},
		Closed:     true,
		Env:        []string{"HOME=" + home},
	}
	var stdout, stderr bytes.Buffer

	result := executeSubjectWithRunBinary(context.Background(), root, root, &stdout, &stderr, nil, forceRun, plannedEvaluation{plan: plan}, nil, "")

	if result.ActionExit != 1 {
		t.Fatalf("result = %#v, want action exit 1; stderr=%q", result, stderr.String())
	}
	requireCacheRefusal(t, stderr.String(), home)
}

// LQ17: an unwritable cache directory refuses the lane the same way, because one rule
// covers every holder.
func TestLaneRefusesAnUnheldCache(t *testing.T) {
	home := unwritableCacheHome(t)
	t.Setenv("HOME", home)
	root := outcomeFixture(t)
	tree := outcomeGit(t, root, "rev-parse", "HEAD^{tree}")
	var stdout, stderr bytes.Buffer

	result, err := RunLane(context.Background(), LaneRequest{
		Root: root, Tree: tree,
		Checks: []Phase{{Name: "declared", Argv: []string{"true"}}},
		Stdout: &stdout, Stderr: &stderr,
	})

	if err == nil {
		t.Fatalf("lane result = %#v, want a refusal", result)
	}
	requireCacheRefusal(t, stderr.String(), home)
}

// requireCacheRefusal grades one call site's refusal: the shared kind, the hold error, and
// the cache path the operator has to fix.
func requireCacheRefusal(t *testing.T, output, home string) {
	t.Helper()
	dir, err := gocache.Dir([]string{"HOME=" + home})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output, "cache lock unavailable") || !strings.Contains(output, dir) {
		t.Fatalf("output = %q, want the cache refusal naming %q", output, dir)
	}
}

// LQ28: a closure that declares a HOME the derivation refuses fails subject construction.
// Dropping the entry there would hand the oracle a hash no child of the run ever used.
func TestSubjectPropagatesTheCacheApplyRefusal(t *testing.T) {
	root := cacheClosureFixture(t)
	t.Setenv("HOME", "relative/home")

	s, err := buildSubject(root)

	if err == nil {
		t.Fatalf("buildSubject = %#v, want the Apply refusal", s)
	}
	if !strings.Contains(err.Error(), "HOME") {
		t.Fatalf("error = %q, want it to name HOME", err)
	}
}

func containsEntry(env []string, want string) bool {
	for _, entry := range env {
		if entry == want {
			return true
		}
	}
	return false
}
