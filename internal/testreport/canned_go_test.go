package testreport

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// cannedSet is one recorded `go test -json` stream and the exit its `go` gives.
// The four sets are the outcome kinds a completed child can reach: a failing test,
// a passing test, a package that fails to build, and a package with no test.
type cannedSet struct {
	name   string
	events []string
	exit   int
}

func cannedSets() []cannedSet {
	return []cannedSet{
		{
			name: "failing",
			events: []string{
				`{"Action":"run","Package":"canned","Test":"TestCanned"}`,
				`{"Action":"output","Package":"canned","Test":"TestCanned","Output":"    canned_test.go:9: boom\n"}`,
				`{"Action":"fail","Package":"canned","Test":"TestCanned","Elapsed":0.01}`,
				`{"Action":"fail","Package":"canned","Elapsed":0.5}`,
			},
			exit: 1,
		},
		{
			name: "passing",
			events: []string{
				`{"Action":"run","Package":"canned","Test":"TestCanned"}`,
				`{"Action":"pass","Package":"canned","Test":"TestCanned","Elapsed":0.01}`,
				`{"Action":"pass","Package":"canned","Elapsed":0.25}`,
			},
			exit: 0,
		},
		{
			name: "build-fail",
			events: []string{
				`{"Action":"build-output","ImportPath":"canned","Output":"./canned.go:3:1: syntax error: unexpected }\n"}`,
				`{"Action":"build-fail","ImportPath":"canned"}`,
			},
			exit: 1,
		},
		{
			name: "no-run",
			events: []string{
				`{"Action":"output","Package":"canned","Output":"?   \tcanned\t[no test files]\n"}`,
				`{"Action":"skip","Package":"canned","Elapsed":0}`,
			},
			exit: 0,
		},
	}
}

// installCannedGo puts a stub `go` on PATH that replays one canned stream and exits with
// the set's own code. The build-fail and failing sets are only complete with the nonzero
// exit, which is why this stub carries a code where writeCheckGo does not.
func installCannedGo(t *testing.T, set cannedSet) {
	t.Helper()
	dir := t.TempDir()
	body := "#!/usr/bin/env bash\n"
	for _, event := range set.events {
		body += "printf '%s\\n' " + sanitize.ShellQuote(event) + "\n"
	}
	body += "exit " + strconv.Itoa(set.exit) + "\n"
	writeGoStub(t, filepath.Join(dir, "go"), body)
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

// installCannedSelection installs a factory whose build writes a canned run binary. The
// answered flag turns true once a run selects that binary.
func installCannedSelection(t *testing.T) *bool {
	t.Helper()
	selected := new(bool)
	installTestSelectionFactory(t, runbinary.Factory{
		TempRoot: t.TempDir(),
		Build: func(_ context.Context, _, output string) error {
			*selected = true
			return os.WriteFile(output, []byte("selected"), 0o755)
		},
		Verify: func(string, string) error { return nil },
	})
	return selected
}

// cannedSelections are the two selection forms the goldens cover.
func cannedSelections() []struct {
	name string
	args []string
} {
	return []struct {
		name string
		args []string
	}{
		{"package", []string{"--package", "./..."}},
		{"run", []string{"--package", "./...", "--run", "^TestCanned$"}},
	}
}
