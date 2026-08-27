package gate

import (
	"bytes"
	"context"
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

func containsEntry(env []string, want string) bool {
	for _, entry := range env {
		if entry == want {
			return true
		}
	}
	return false
}
