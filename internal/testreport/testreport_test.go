package testreport

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gocache"
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
