package testreport

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/runbinary"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/toon"
)

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
		{set: "failing", want: Outcome{Kind: OutcomeFailed, FailedTests: 1, FailedNames: []string{"TestCanned"}, Ran: 1}},
		{set: "passing", want: Outcome{Kind: OutcomePassed, Ran: 1}},
		{set: "build-fail", want: Outcome{Kind: OutcomeBuildFailed, BuildLine: "./canned.go:3:1: syntax error: unexpected }"}},
		{set: "no-run", want: Outcome{Kind: OutcomeNoTestRun}},
	} {
		t.Run(tc.set, func(t *testing.T) {
			installCannedGo(t, cannedSetNamed(t, tc.set))
			installCannedSelection(t)

			outcome := executeSelection(t, t.TempDir(), []string{"--package", "./..."})
			if !sameOutcome(outcome, tc.want) {
				t.Fatalf("outcome = %+v, want %+v", outcome, tc.want)
			}
		})
	}

	t.Run("build failure through the run pattern", func(t *testing.T) {
		installCannedGo(t, cannedSetNamed(t, "build-fail"))
		installCannedSelection(t)
		outcome := executeSelection(t, t.TempDir(), []string{"--package", "./...", "--run", "^TestCanned$"})
		if want := (Outcome{Kind: OutcomeBuildFailed, BuildLine: "./canned.go:3:1: syntax error: unexpected }"}); !sameOutcome(outcome, want) {
			t.Fatalf("outcome = %+v, want %+v", outcome, want)
		}
	})

	t.Run("no-test-run through the run pattern", func(t *testing.T) {
		installCannedGo(t, cannedSetNamed(t, "no-run"))
		installCannedSelection(t)

		outcome := executeSelection(t, t.TempDir(), []string{"--package", "./...", "--run", "^TestCanned$"})
		if want := (Outcome{Kind: OutcomeNoTestRun}); !sameOutcome(outcome, want) {
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
		if want := (Outcome{Kind: OutcomeRefused}); !sameOutcome(outcome, want) {
			t.Fatalf("outcome = %+v, want %+v", outcome, want)
		}
	})

	t.Run("interrupted", func(t *testing.T) {
		installSignallingGo(t)
		installCannedSelection(t)

		outcome := executeSelection(t, t.TempDir(), []string{"--package", "./..."})
		if want := (Outcome{Kind: OutcomeInterrupted}); !sameOutcome(outcome, want) {
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

// sameOutcome compares two outcomes field by field, because the FailedNames slice makes
// Outcome incomparable with ==.
func sameOutcome(got, want Outcome) bool {
	return got.Kind == want.Kind && got.FailedTests == want.FailedTests && got.Ran == want.Ran && slices.Equal(got.FailedNames, want.FailedNames) && got.BuildLine == want.BuildLine
}

// TestOnlyAScriptFailureIsBuildFailed covers PS45, PS46, PS47, PS51, and PS52. A run-binary
// build script that exits 1 reads build-failed with the refusal bytes and exit that `bench
// test` prints, and BuildLine names the first output line with a Go position, or else the
// first non-empty line. A cancelled build and a build with no Go on PATH stay refused.
func TestOnlyAScriptFailureIsBuildFailed(t *testing.T) {
	failing := func(output string) string {
		return "printf '%s' " + sanitize.ShellQuote(output) + "\nexit 1\n"
	}
	for _, tc := range []struct {
		name, output, message, line string
	}{
		{"a header before a .go: line", "# fixture\n./broken.go:3:1: syntax error\n", "# fixture\n./broken.go:3:1: syntax error", "./broken.go:3:1: syntax error"},
		{"no .go: line", "\n  go-build.sh: no toolchain\nsecond line\n", "go-build.sh: no toolchain\nsecond line", "go-build.sh: no toolchain"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			outcome, out, code := executeOverBuildScript(t, failing(tc.output), nil)
			want := toon.Errorf("Bench executable selection failed", "build private Bench executable: exit status 1: "+tc.message) + "\n"
			if out != want || code != 1 {
				t.Fatalf("Execute = (%d, %q), want the selection refusal (1, %q)", code, out, want)
			}
			if want := (Outcome{Kind: OutcomeBuildFailed, BuildLine: tc.line}); !sameOutcome(outcome, want) {
				t.Fatalf("outcome = %+v, want %+v", outcome, want)
			}
		})
	}
	t.Run("a cancelled build", func(t *testing.T) {
		cancelled := func(ctx context.Context, source, output string) error {
			ctx, cancel := context.WithCancel(ctx)
			cancel()
			return runbinary.Build(ctx, source, output)
		}
		outcome, out, _ := executeOverBuildScript(t, "exec sleep 30\n", cancelled)
		if want := (Outcome{Kind: OutcomeRefused}); !sameOutcome(outcome, want) {
			t.Fatalf("outcome = %+v, want %+v\n%s", outcome, want, out)
		}
	})
	t.Run("no Go on PATH", func(t *testing.T) {
		t.Setenv("PATH", t.TempDir())
		outcome, out, _ := executeOverBuildScript(t, failing("./broken.go:3:1: syntax error\n"), nil)
		if want := (Outcome{Kind: OutcomeRefused}); !sameOutcome(outcome, want) {
			t.Fatalf("outcome = %+v, want %+v\n%s", outcome, want, out)
		}
	})
}

// executeOverBuildScript runs one package selection whose run binary builds through a
// fixture scripts/go-build.sh with the body script. A nil build selects the canonical
// builder.
func executeOverBuildScript(t *testing.T, script string, build runbinary.Builder) (Outcome, string, int) {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "scripts", "go-build.sh"), []byte("#!/usr/bin/env bash\n"+script), 0o755); err != nil {
		t.Fatal(err)
	}
	installTestSelectionFactory(t, runbinary.Factory{TempRoot: t.TempDir(), Build: build, Verify: func(string, string) error { return nil }})
	request, line, code := Prepare(root, []string{"--package", "./..."})
	if line != "" {
		t.Fatalf("Prepare refused the selection: (%d, %q)", code, line)
	}
	return Execute(root, request)
}

// TestOutcomeNamesTheFailedTests covers PS56: FailedNames lists the failed tests in the
// order that the failures table prints them, package first and then test, whatever the
// event order. A prose run counts its findings and names no test.
func TestOutcomeNamesTheFailedTests(t *testing.T) {
	set := cannedSet{name: t.Name(), exit: 1}
	for _, failed := range [][2]string{{"z", "TestA"}, {"a", "TestC"}, {"a", "TestB"}} {
		set.events = append(set.events,
			`{"Action":"run","Package":"`+failed[0]+`","Test":"`+failed[1]+`"}`,
			`{"Action":"output","Package":"`+failed[0]+`","Test":"`+failed[1]+`","Output":"    boom\n"}`,
			`{"Action":"fail","Package":"`+failed[0]+`","Test":"`+failed[1]+`","Elapsed":0}`)
	}
	set.events = append(set.events, `{"Action":"fail","Package":"z","Elapsed":0.1}`, `{"Action":"fail","Package":"a","Elapsed":0.1}`)
	installCannedGo(t, set)
	installCannedSelection(t)
	outcome := executeSelection(t, t.TempDir(), []string{"--package", "./..."})
	if want := (Outcome{Kind: OutcomeFailed, FailedTests: 3, FailedNames: []string{"TestB", "TestC", "TestA"}, Ran: 3}); !sameOutcome(outcome, want) {
		t.Fatalf("outcome = %+v, want %+v", outcome, want)
	}

	root := t.TempDir()
	writeProseCheckFile(t, root, ".bench/prose-exclusions", "")
	writeProseCheckFile(t, root, "long.md", strings.Repeat("word ", 27)+".\n")
	prose := executeSelection(t, root, []string{"--check", ProseCheckName})
	if prose.Kind != OutcomeFailed || prose.FailedTests == 0 || len(prose.FailedNames) != 0 {
		t.Fatalf("prose outcome = %+v, want a failed run with findings and no failed names", prose)
	}
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
	writeGoStub(t, filepath.Join(dir, "go"), body)
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

// baseGoldens holds the exact stdout and exit `Command` owes for each canned set in each
// selection form. Every literal is independent of the Prepare/Execute pair, so a rendering
// or exit change anywhere in the pair reds this table rather than passing as a fresh
// expectation. (Coverage row PB28.)
var baseGoldens = map[string]struct {
	output string
	code   int
}{
	"failing/package": {
		output: "packages[1]{package,status,elapsed_ms,tests_run}:\n  canned,fail,500,1\nfailures[1]{package,test,line,lines}:\n  canned,TestCanned,\"canned_test.go:9: boom\",1\nskips[0]{package,test,reason}:\n",
		code:   1,
	},
	"failing/run": {
		output: "packages[1]{package,status,elapsed_ms,tests_run}:\n  canned,fail,500,1\nfailures[1]{package,test,line,lines}:\n  canned,TestCanned,\"canned_test.go:9: boom\",1\nskips[0]{package,test,reason}:\n",
		code:   1,
	},
	"passing/package": {
		output: "packages[1]{package,status,elapsed_ms,tests_run}:\n  canned,pass,250,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n",
		code:   0,
	},
	"passing/run": {
		output: "packages[1]{package,status,elapsed_ms,tests_run}:\n  canned,pass,250,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n",
		code:   0,
	},
	"build-fail/package": {
		output: "packages[1]{package,status,elapsed_ms,tests_run}:\n  canned,fail,0,0\nfailures[1]{package,test,line,lines}:\n  canned,\"\",\"./canned.go:3:1: syntax error: unexpected }\",1\nskips[0]{package,test,reason}:\n",
		code:   1,
	},
	"build-fail/run": {
		output: "packages[1]{package,status,elapsed_ms,tests_run}:\n  canned,fail,0,0\nfailures[1]{package,test,line,lines}:\n  canned,\"\",\"./canned.go:3:1: syntax error: unexpected }\",1\nskips[0]{package,test,reason}:\n",
		code:   1,
	},
	"no-run/package": {
		output: "packages[1]{package,status,elapsed_ms,tests_run}:\n  canned,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n",
		code:   0,
	},
	"no-run/run": {
		output: "error: go test reported no test runs — run pattern matched no tests\n",
		code:   1,
	},
}

// TestCommandKeepsItsBaseOutput pins the command output and exit for each selection.
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
					t.Fatalf("Command = (%d, %q), want the command contract (%d, %q)", code, output, golden.code, golden.output)
				}
			})
		}
	}
}
