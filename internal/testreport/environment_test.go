package testreport

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/gocache"
	"github.com/gibbonmi/bench/internal/gocache/cleanprobe"
	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// A focused run of the kit's tests starts git with no auto-maintenance. A background
// prune from that maintenance can remove a worktree admin entry a test fixture planted.
func TestKitRunEnvironmentGitStartsNoAutoMaintenance(t *testing.T) {
	kit := t.TempDir()
	t.Setenv("BENCH_KIT", kit)
	t.Setenv("TMPDIR", t.TempDir())
	cache, err := gocache.Dir(os.Environ())
	if err != nil {
		t.Fatal(err)
	}
	argv := gittest.KitRunProbe(t, os.Environ())
	for i := range argv {
		argv[i] = sanitize.ShellQuote(argv[i])
	}
	marker := filepath.Join(t.TempDir(), "home")
	body := "#!/bin/sh\nset -eu\n" + strings.Join(argv, " ") + "\n" +
		"test \"$GOCACHE\" = " + sanitize.ShellQuote(cache) + "\n" +
		"printf '%s' \"$HOME\" > " + sanitize.ShellQuote(marker) + "\n" +
		"printf '%s\\n' '{\"Action\":\"pass\",\"Package\":\"kitprobe\"}'\n"
	goDir := t.TempDir()
	writeGoStub(t, filepath.Join(goDir, "go"), body)
	t.Setenv("PATH", goDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	installCannedSelection(t)
	if out, code := Command(kit, nil); code != 0 {
		t.Fatalf("kit run probe exit = %d: %s", code, out)
	}
	home, err := os.ReadFile(marker)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Dir(string(home))); !os.IsNotExist(err) {
		t.Fatalf("private run remains after child exit: %v", err)
	}
}

func TestLinkedRunKeepsOperatorEnvironment(t *testing.T) {
	root, kit := t.TempDir(), t.TempDir()
	t.Setenv("BENCH_KIT", kit)
	marker := filepath.Join(t.TempDir(), "environment")
	goDir := t.TempDir()
	writeCheckGo(t, filepath.Join(goDir, "go"), marker)
	t.Setenv("PATH", goDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	installCannedSelection(t)
	t.Setenv("BENCH_KIT", kit)
	if out, code := Command(root, nil); code != 0 {
		t.Fatalf("linked run exit = %d: %s", code, out)
	}
	child := readTestReportFile(t, marker)
	for _, key := range []string{"HOME", "TMPDIR", "GIT_CONFIG_GLOBAL", "GIT_CONFIG_NOSYSTEM"} {
		if value, present := os.LookupEnv(key); present && !strings.Contains(child, key+"="+value+"\n") {
			t.Errorf("linked child changed %s", key)
		}
	}
}

// Stubbed test output must not replace the Go settings that the run owner resolves.
func writeGoStub(t *testing.T, path, source string) {
	t.Helper()
	realGo, err := exec.LookPath("go")
	if err != nil {
		t.Fatal(err)
	}
	header, body, ok := strings.Cut(source, "\n")
	if !ok {
		t.Fatal("Go stub needs an interpreter line")
	}
	source = header + "\nif [ \"$1\" = env ]; then exec " + sanitize.ShellQuote(realGo) + " \"$@\"; fi\n" + body
	if err := os.WriteFile(path, []byte(source), 0o755); err != nil {
		t.Fatal(err)
	}
}

// C10: the `bench test` child carries the Bench build cache entry, so a focused run
// warms the archives a later gate reads.
func TestTestEnvironmentCarriesTheBenchBuildCache(t *testing.T) {
	child, err := testEnvironment([]string{"HOME=/home/agent", "GOCACHE=/ambient/cache"}, "/selected/bench")
	if err != nil {
		t.Fatal(err)
	}
	want, err := gocache.Dir([]string{"HOME=/home/agent"})
	if err != nil {
		t.Fatal(err)
	}
	entries := []string{}
	for _, entry := range child {
		if strings.HasPrefix(entry, gocache.Env+"=") {
			entries = append(entries, entry)
		}
	}
	if len(entries) != 1 || entries[0] != gocache.Env+"="+want {
		t.Fatalf("cache entries = %#v, want exactly %s=%s", entries, gocache.Env, want)
	}
}

func TestTestEnvironmentRefusesWithoutAnAbsoluteHome(t *testing.T) {
	child, err := testEnvironment([]string{"PATH=/usr/bin"}, "/selected/bench")
	if err == nil {
		t.Fatalf("testEnvironment = %#v, want an error", child)
	}
	if !strings.Contains(err.Error(), "HOME") {
		t.Fatalf("error = %q, want it to name HOME", err)
	}
}

// cacheCleanProbeBinaryEnv names this test binary, so the module under test can
// re-execute it. The answer entry's name and wire format live in cleanprobe.
const cacheCleanProbeBinaryEnv = "BENCH_TEST_CACHE_CLEAN_PROBE_BIN"

// TestCacheCleanProbe is the second process the holder row drives. The shared body runs
// `bench cache clean` and records the verb's own answer.
func TestCacheCleanProbe(t *testing.T) { cleanprobe.Answer(t) }

// L02: a `bench test` run holds the shared cache lock across its go test child, so
// `bench cache clean` exits 1 while that child is compiling. The probe runs from inside
// the child's own test, which is the one point inside the run's span a second process can
// observe.
func TestFocusedRunHoldsTheCacheLockAcrossItsGoTestChild(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	answerPath := filepath.Join(t.TempDir(), "clean-answer")
	t.Setenv(cleanprobe.Env, answerPath)
	binary, err := filepath.Abs(os.Args[0])
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv(cacheCleanProbeBinaryEnv, binary)

	selected := filepath.Join(t.TempDir(), "bench")
	if err := os.WriteFile(selected, []byte("selected"), 0o755); err != nil {
		t.Fatal(err)
	}
	installTestSelectionFactory(t, runbinary.Factory{Verify: func(string, string) error { return nil }})
	t.Setenv(runbinary.Env, selected)

	if output, code := Command(cleanProbeModule(t, home), nil); code != 0 {
		t.Fatalf("Command = %d\n%s", code, output)
	}
	cleanprobe.Require(t, readTestReportFile(t, answerPath))
}

// cleanProbeModule writes the one-package module the focused run compiles. Its single test
// re-executes this test binary's probe row, so the clean it attempts runs while the parent
// still holds the lock.
func cleanProbeModule(t *testing.T, home string) string {
	t.Helper()
	root := t.TempDir()
	source := `package cleanprobe

import (
	"os"
	"os/exec"
	"testing"
)

func TestCleanProbe(t *testing.T) {
	probe := exec.Command(os.Getenv("` + cacheCleanProbeBinaryEnv + `"), "-test.run=^TestCacheCleanProbe$")
	// The competing clean targets the operator cache, not the private test home.
	probe.Env = append(os.Environ(), "HOME=" + ` + strconv.Quote(home) + `)
	if output, err := probe.CombinedOutput(); err != nil {
		t.Fatalf("clean probe: %v\n%s", err, output)
	}
}
`
	for name, body := range map[string]string{
		"go.mod":             "module cleanprobe\n\ngo 1.25\n",
		"cleanprobe_test.go": source,
	} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}
