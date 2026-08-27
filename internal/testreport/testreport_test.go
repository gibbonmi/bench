package testreport

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

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

func TestFocusedRequestPackageCompatibility(t *testing.T) {
	root := focusedTestModule(t)
	installTestSelectionFactory(t, runbinary.Factory{
		TempRoot: t.TempDir(),
		Build: func(_ context.Context, _, output string) error {
			return os.WriteFile(output, []byte("selected"), 0o755)
		},
		Verify: func(string, string) error { return nil },
	})

	for _, tc := range []struct {
		name string
		args []string
	}{
		{name: "default"},
		{name: "positional", args: []string{"chosen"}},
		{name: "package flag", args: []string{"--package", "chosen"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			output, code := Command(root, tc.args)
			if code != 0 {
				t.Fatalf("Command(%v) = %d\n%s", tc.args, code, output)
			}
			if !strings.Contains(output, "focusedfixture/chosen") {
				t.Fatalf("Command(%v) output = %q, want selected package", tc.args, output)
			}
		})
	}
}

func TestRunPatternReachesGoAsOneArgument(t *testing.T) {
	root := focusedTestModule(t)
	installTestSelectionFactory(t, runbinary.Factory{
		TempRoot: t.TempDir(),
		Build: func(_ context.Context, _, output string) error {
			return os.WriteFile(output, []byte("selected"), 0o755)
		},
		Verify: func(string, string) error { return nil },
	})
	const pattern = "^TestRunArgument$"
	for _, tc := range []struct {
		name string
		args []string
	}{
		{name: "default", args: []string{"--run", pattern}},
		{name: "explicit package", args: []string{"--package", "chosen", "--run", pattern}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			marker := filepath.Join(t.TempDir(), "argv")
			t.Setenv("BENCH_TEST_MARKER", marker)
			output, code := Command(root, tc.args)
			if code != 0 {
				t.Fatalf("Command(%v) = %d\n%s", tc.args, code, output)
			}
			if got := readTestReportFile(t, marker); !strings.Contains(got, "-test.run="+pattern) {
				t.Fatalf("test argv = %q, want %q as one value", got, "-test.run="+pattern)
			}
		})
	}
}

func TestRunPatternRefusesZeroMatches(t *testing.T) {
	root := focusedTestModule(t)
	addFocusedFailureCases(t, root)
	installTestSelectionFactory(t, runbinary.Factory{
		TempRoot: t.TempDir(),
		Build: func(_ context.Context, _, output string) error {
			return os.WriteFile(output, []byte("selected"), 0o755)
		},
		Verify: func(string, string) error { return nil },
	})
	for _, tc := range []struct {
		name     string
		pattern  string
		wantCode int
		want     string
	}{
		{name: "zero matches", pattern: "^Missing$", wantCode: 1, want: "go test reported no test runs"},
		{name: "matched skip", pattern: "^TestMatchedSkip$", wantCode: 0, want: "TestMatchedSkip"},
		{name: "matched failure", pattern: "^TestMatchedFailure$", wantCode: 1, want: "TestMatchedFailure"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			args := []string{"--package", "chosen", "--run", tc.pattern}
			output, code := Command(root, args)
			if code != tc.wantCode {
				t.Fatalf("Command(%v) = %d, want %d\n%s", args, code, tc.wantCode, output)
			}
			if !strings.Contains(output, tc.want) {
				t.Fatalf("Command(%v) output = %q, want %q", args, output, tc.want)
			}
		})
	}
}

func TestExplicitFocusedRunsWriteNoGateOwnedRecords(t *testing.T) {
	root := focusedTestModule(t)
	installTestSelectionFactory(t, runbinary.Factory{
		TempRoot: t.TempDir(),
		Build: func(_ context.Context, _, output string) error {
			return os.WriteFile(output, []byte("selected"), 0o755)
		},
		Verify: func(string, string) error { return nil },
	})
	records := map[string]string{
		".git/bench-last-gate":            "cache",
		".git/bench-last-lane":            "lane",
		".git/bench-gate-owner":           "owner",
		".git/bench-gate-pin":             "pin",
		".git/bench-gate-evidence/record": "evidence",
		".git/refs/bench/green/main":      "green",
	}
	for path, contents := range records {
		file := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(file, []byte(contents), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	for _, args := range [][]string{nil, {"--package", "chosen"}, {"--run", "^TestRunArgument$"}} {
		output, code := Command(root, args)
		if code != 0 {
			t.Fatalf("Command(%v) = %d\n%s", args, code, output)
		}
		for path, want := range records {
			if got := readTestReportFile(t, filepath.Join(root, path)); got != want {
				t.Fatalf("gate-owned record %s = %q, want %q after Command(%v)", path, got, want, args)
			}
		}
	}
}

func focusedTestModule(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for path, source := range map[string]string{
		"go.mod":                "module focusedfixture\n\ngo 1.25\n",
		"root_test.go":          focusedTestSource("focusedfixture"),
		"chosen/chosen_test.go": focusedTestSource("chosen"),
	} {
		file := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(file, []byte(source), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func focusedTestSource(pkg string) string {
	return "package " + pkg + `

import (
	"os"
	"strings"
	"testing"
)

func TestRoot(t *testing.T) {}

func TestChosen(t *testing.T) {}

func TestRunArgument(t *testing.T) {
	if marker := os.Getenv("BENCH_TEST_MARKER"); marker != "" {
		if err := os.WriteFile(marker, []byte(strings.Join(os.Args, "\n")), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}
`
}

func addFocusedFailureCases(t *testing.T, root string) {
	t.Helper()
	const source = `package chosen

import "testing"

func TestMatchedSkip(t *testing.T) { t.Skip("matched skip") }

func TestMatchedFailure(t *testing.T) { t.Fatal("matched failure") }
`
	if err := os.WriteFile(filepath.Join(root, "chosen", "failure_test.go"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
}
