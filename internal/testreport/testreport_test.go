package testreport

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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
