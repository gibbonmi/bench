package testreport

import (
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gocache"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// TestDecodeSkipsRunnerLines pins that decode classifies runner lines through testlines.
func TestDecodeSkipsRunnerLines(t *testing.T) {
	stream := strings.Join([]string{
		`{"Action":"output","Package":"p","Test":"TestA","Output":"=== RUN   TestA\n"}`,
		`{"Action":"output","Package":"p","Test":"TestA","Output":"    a_test.go:9: boom\n"}`,
		`{"Action":"output","Package":"p","Test":"TestA","Output":"--- FAIL: TestA (0.00s)\n"}`,
		`{"Action":"output","Package":"p","Output":"ok  \tp\t0.10s\n"}`,
		`{"Action":"output","Package":"p","Output":"package diagnostic\n"}`,
	}, "\n")

	report, err := decode(strings.NewReader(stream))
	if err != nil {
		t.Fatal(err)
	}
	test := report.test("p", "TestA")
	if test.first != "a_test.go:9: boom" || test.last != "a_test.go:9: boom" {
		t.Errorf("test diagnostic = %q/%q, want the non-runner line only", test.first, test.last)
	}
	if report.packageLog["p"] != "package diagnostic" {
		t.Errorf("packageLog = %q, want the non-runner line only", report.packageLog["p"])
	}
}

func TestPackagePattern(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "internal", "usage"), 0o755); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		operand string
		want    string
	}{
		{"internal/usage", "./internal/usage"},
		{"internal/usage/...", "./internal/usage/..."},
		{"./internal/usage", "./internal/usage"},
		{"../x", "../x"},
		{"github.com/gibbonmi/bench/internal/usage", "github.com/gibbonmi/bench/internal/usage"},
		{"does/not/exist", "does/not/exist"},
		{"", "./..."},
	}
	for _, c := range cases {
		if got := packagePattern(root, c.operand); got != c.want {
			t.Errorf("packagePattern(%q) = %q, want %q", c.operand, got, c.want)
		}
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

// T04: the focused `bench test` argv carries -trimpath and -count=1, so a focused run
// warms the same cache entries the gate reads instead of writing a path-keyed set.
func TestFocusedTestArgvCarriesTrimPathAndCountOne(t *testing.T) {
	want := []string{"go", "test", "-trimpath", "-count=1", "-json", "./internal/gate"}
	if got := focusedTestArgv("./internal/gate"); !reflect.DeepEqual(got, want) {
		t.Fatalf("focused test argv = %#v, want %#v", got, want)
	}
}

// The clean probe's protocol. The answer file carries the verb's exit code and its own
// output, and the binary entry lets the module under test re-execute this test binary. A
// child with no answer entry is inert, so the probe row is a no-op in an ordinary run.
const (
	cacheCleanProbeEnv       = "BENCH_TEST_CACHE_CLEAN_PROBE"
	cacheCleanProbeBinaryEnv = "BENCH_TEST_CACHE_CLEAN_PROBE_BIN"
)

// TestCacheCleanProbe is the second process the holder row drives. It runs
// `bench cache clean` and records the verb's own answer, because a POSIX record lock is
// owned per process: a clean inside the holder's process could never contend with it.
func TestCacheCleanProbe(t *testing.T) {
	answerPath := os.Getenv(cacheCleanProbeEnv)
	if answerPath == "" {
		return
	}
	answer, code := gocache.Command([]string{"clean"})
	if err := os.WriteFile(answerPath, []byte(strconv.Itoa(code)+"\n"+answer), 0o600); err != nil {
		t.Fatal(err)
	}
}

// L02: a `bench test` run holds the shared cache lock across its go test child, so
// `bench cache clean` exits 1 while that child is compiling. The probe runs from inside
// the child's own test, which is the one point inside the run's span a second process can
// observe.
func TestFocusedRunHoldsTheCacheLockAcrossItsGoTestChild(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	answerPath := filepath.Join(t.TempDir(), "clean-answer")
	t.Setenv(cacheCleanProbeEnv, answerPath)
	binary, err := filepath.Abs(os.Args[0])
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv(cacheCleanProbeBinaryEnv, binary)

	selected := filepath.Join(t.TempDir(), "bench")
	if err := os.WriteFile(selected, []byte("selected"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv(runbinary.Env, selected)
	installTestSelectionFactory(t, runbinary.Factory{Verify: func(string, string) error { return nil }})

	if output, code := Command(cleanProbeModule(t), nil); code != 0 {
		t.Fatalf("Command = %d\n%s", code, output)
	}
	answer := readTestReportFile(t, answerPath)
	exit, rest, _ := strings.Cut(answer, "\n")
	if exit != "1" || !strings.HasPrefix(rest, "error: cache in use — ") {
		t.Fatalf("clean beside the focused run = exit %s, %q; want the cache-in-use refusal at exit 1", exit, rest)
	}
}

// cleanProbeModule writes the one-package module the focused run compiles. Its single test
// re-executes this test binary's probe row, so the clean it attempts runs while the parent
// still holds the lock.
func cleanProbeModule(t *testing.T) string {
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
