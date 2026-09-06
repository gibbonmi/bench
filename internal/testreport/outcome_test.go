package testreport

import (
	"context"
	"errors"
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

// installCannedGo puts a stub `go` on PATH that replays one canned stream. The prior
// art is writeCheckGo in check_test.go; this stub adds the exit code, because the
// build-fail and failing sets are only complete with the nonzero exit.
func installCannedGo(t *testing.T, set cannedSet) {
	t.Helper()
	dir := t.TempDir()
	body := "#!/usr/bin/env bash\n"
	for _, event := range set.events {
		body += "printf '%s\\n' " + sanitize.ShellQuote(event) + "\n"
	}
	body += "exit " + strconv.Itoa(set.exit) + "\n"
	if err := os.WriteFile(filepath.Join(dir, "go"), []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func installCannedSelection(t *testing.T) {
	t.Helper()
	installTestSelectionFactory(t, runbinary.Factory{
		TempRoot: t.TempDir(),
		Build: func(_ context.Context, _, output string) error {
			return os.WriteFile(output, []byte("selected"), 0o755)
		},
		Verify: func(string, string) error { return nil },
	})
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

// TestExecuteClassifiesTheOutcome grades the six kinds a focused run can reach. The
// four completed kinds run through the canned streams, `refused` runs through a
// selection factory whose build fails, and `interrupted` runs through a stub `go`
// that signals the owner and then parks. (Coverage row PB27.)
func TestExecuteClassifiesTheOutcome(t *testing.T) {
	for _, tc := range []struct {
		set   string
		want  Outcome
		count int
	}{
		{set: "failing", want: Outcome{Kind: OutcomeFailed, FailedTests: 1}},
		{set: "passing", want: Outcome{Kind: OutcomePassed}},
		{set: "build-fail", want: Outcome{Kind: OutcomeBuildFailed}},
		{set: "no-run", want: Outcome{Kind: OutcomeNoTestRun}},
	} {
		t.Run(tc.set, func(t *testing.T) {
			installCannedGo(t, cannedSetNamed(t, tc.set))
			installCannedSelection(t)

			outcome := executeSelection(t, t.TempDir(), []string{"--package", "./..."})
			if outcome != tc.want {
				t.Fatalf("outcome = %+v, want %+v", outcome, tc.want)
			}
		})
	}

	t.Run("no-test-run through the run pattern", func(t *testing.T) {
		installCannedGo(t, cannedSetNamed(t, "no-run"))
		installCannedSelection(t)

		outcome := executeSelection(t, t.TempDir(), []string{"--package", "./...", "--run", "^TestCanned$"})
		if want := (Outcome{Kind: OutcomeNoTestRun}); outcome != want {
			t.Fatalf("outcome = %+v, want %+v", outcome, want)
		}
	})

	t.Run("refused", func(t *testing.T) {
		installCannedGo(t, cannedSetNamed(t, "passing"))
		installTestSelectionFactory(t, runbinary.Factory{
			TempRoot: t.TempDir(),
			Build:    func(context.Context, string, string) error { return errors.New("no Bench executable") },
			Verify:   func(string, string) error { return nil },
		})

		outcome := executeSelection(t, t.TempDir(), []string{"--package", "./..."})
		if want := (Outcome{Kind: OutcomeRefused}); outcome != want {
			t.Fatalf("outcome = %+v, want %+v", outcome, want)
		}
	})

	t.Run("interrupted", func(t *testing.T) {
		installSignallingGo(t)
		installCannedSelection(t)

		outcome := executeSelection(t, t.TempDir(), []string{"--package", "./..."})
		if want := (Outcome{Kind: OutcomeInterrupted}); outcome != want {
			t.Fatalf("outcome = %+v, want %+v", outcome, want)
		}
	})
}

// TestCommandIsThePrepareExecuteProjection grades that `Command` is the two steps and
// nothing else, so no second render path can drift from the pair. (Coverage row PB28.)
func TestCommandIsThePrepareExecuteProjection(t *testing.T) {
	for _, set := range cannedSets() {
		for _, selection := range cannedSelections() {
			t.Run(set.name+"/"+selection.name, func(t *testing.T) {
				installCannedGo(t, set)
				installCannedSelection(t)
				root := t.TempDir()

				commandOutput, commandCode := Command(root, selection.args)
				request, line, code := Prepare(root, selection.args)
				if line != "" {
					t.Fatalf("Prepare refused the selection: (%d, %q)", code, line)
				}
				_, stepOutput, stepCode := Execute(root, request)
				if commandOutput != stepOutput || commandCode != stepCode {
					t.Fatalf("Command = (%d, %q), want Prepare then Execute (%d, %q)", commandCode, commandOutput, stepCode, stepOutput)
				}
			})
		}
	}
}

func executeSelection(t *testing.T, root string, args []string) Outcome {
	t.Helper()
	request, line, code := Prepare(root, args)
	if line != "" {
		t.Fatalf("Prepare refused the selection: (%d, %q)", code, line)
	}
	outcome, _, _ := Execute(root, request)
	return outcome
}

func cannedSetNamed(t *testing.T, name string) cannedSet {
	t.Helper()
	for _, set := range cannedSets() {
		if set.name == name {
			return set
		}
	}
	t.Fatalf("no canned set named %q", name)
	return cannedSet{}
}

// installSignallingGo puts a stub `go` on PATH that interrupts the owner and then
// parks. The signal reaches the owner only after Execute installed its cancel
// handler, because the stub is the child that Execute started.
func installSignallingGo(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	body := "#!/usr/bin/env bash\nkill -INT \"$PPID\"\nsleep " + strconv.FormatInt(parkSeconds(), 10) + "\n"
	if err := os.WriteFile(filepath.Join(dir, "go"), []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

// baseGoldens holds the exact stdout and exit that `Command` gave for each canned set
// in each selection form. Every literal was recorded from the base commit
// 963aa880f9365b6a34211e49fd067ef421bff399, before the Prepare/Execute split, so a
// rendering or exit change anywhere in the split reds this table rather than passing
// as a fresh expectation. (Coverage row PB28.)
var baseGoldens = map[string]struct {
	output string
	code   int
}{
	"failing/package": {
		output: "packages[1]{package,status,elapsed_ms}:\n  canned,fail,500\nfailures[1]{package,test,line}:\n  canned,TestCanned,\"canned_test.go:9: boom\"\nskips[0]{package,test,reason}:\n",
		code:   1,
	},
	"failing/run": {
		output: "packages[1]{package,status,elapsed_ms}:\n  canned,fail,500\nfailures[1]{package,test,line}:\n  canned,TestCanned,\"canned_test.go:9: boom\"\nskips[0]{package,test,reason}:\n",
		code:   1,
	},
	"passing/package": {
		output: "packages[1]{package,status,elapsed_ms}:\n  canned,pass,250\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n",
		code:   0,
	},
	"passing/run": {
		output: "packages[1]{package,status,elapsed_ms}:\n  canned,pass,250\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n",
		code:   0,
	},
	"build-fail/package": {
		output: "packages[1]{package,status,elapsed_ms}:\n  canned,fail,0\nfailures[1]{package,test,line}:\n  canned,\"\",\"./canned.go:3:1: syntax error: unexpected }\"\nskips[0]{package,test,reason}:\n",
		code:   1,
	},
	"build-fail/run": {
		output: "error: go test reported no test runs — run pattern matched no tests\n",
		code:   1,
	},
	"no-run/package": {
		output: "packages[1]{package,status,elapsed_ms}:\n  canned,no-tests,0\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:\n",
		code:   0,
	},
	"no-run/run": {
		output: "error: go test reported no test runs — run pattern matched no tests\n",
		code:   1,
	},
}

// TestCommandKeepsItsBaseOutput holds `bench test`'s contract across the split: the
// caller reads the same bytes and the same exit it read at the base commit.
func TestCommandKeepsItsBaseOutput(t *testing.T) {
	for _, set := range cannedSets() {
		for _, selection := range cannedSelections() {
			key := set.name + "/" + selection.name
			golden, ok := baseGoldens[key]
			if !ok {
				t.Fatalf("no recorded golden for %s", key)
			}
			t.Run(key, func(t *testing.T) {
				installCannedGo(t, set)
				installCannedSelection(t)

				output, code := Command(t.TempDir(), selection.args)
				if output != golden.output || code != golden.code {
					t.Fatalf("Command = (%d, %q), want the base commit's (%d, %q)", code, output, golden.code, golden.output)
				}
			})
		}
	}
}
