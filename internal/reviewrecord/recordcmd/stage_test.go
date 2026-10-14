package recordcmd_test

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/fs"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/responsebound"
	"github.com/gibbonmi/bench/internal/responsebound/responseboundtest"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordcmd"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// stageArgs is the excerpt argv for requirement of chunk 1, with the extra flags after it.
func stageArgs(t *testing.T, requirement string, extra ...string) []string {
	t.Helper()
	return append([]string{"excerpt", slug(t), "--chunk", "1", "--requirement", requirement}, extra...)
}

// staged answers the staged excerpt of requirement of chunk 1.
func staged(t *testing.T, f *recordtest.Fixture, home, requirement string) (responsebound.Staged, error) {
	t.Helper()
	return responsebound.ReadStaged(home, f.Root, slug(t), "1", requirement)
}

// requireNoWorktreeWrite requires a clean worktree, so the form wrote no review record and no
// worktree file.
func requireNoWorktreeWrite(t *testing.T, f *recordtest.Fixture) {
	t.Helper()
	if status := f.Git("status", "--porcelain", "--untracked-files=all"); status != "" {
		t.Fatalf("worktree status = %q, want clean", status)
	}
}

// setAnswer is a one-row mutant-set answer whose row bit, authored apart from the producer in
// the envelope that the probe owner reads.
const setAnswer = "mutants[1]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n" +
	"  1,bit,swap,failed,1,TestClampNegative,1,1,0,null,yes\n" +
	"selection[1]{form,target,run,baseline,builds}:\n  package,./,^TestClampNegative$,passed,0\n" +
	"probe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"

// RT10, RT57: a retained call that equals the planned command stages its bytes, exit, commit,
// and dirty paths, and the planned single-quoted operand folds as the shell folds it. A probe
// requirement stages a `bench probe` answer, and a mutant-set answer stages with its row.
func TestExcerptStagesTheRetainedResponse(t *testing.T) {
	for _, tc := range []struct {
		name, requirement, verb string
		args                    []string
		output                  string
		extra                   []string
		mutant                  int
	}{
		{"quoted operand", "quoted", "test", []string{"--run", "TestA|TestB"}, "ok\n", nil, 0},
		{"probe", "tests", "probe", []string{"clamp.go"}, answer("bit", "yes", 0, "1", report), nil, 0},
		{"probe row", "tests", "probe", []string{"--mutants", "m.json"}, setAnswer, []string{"--mutant", "1"}, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, home := recordtest.Staging(t)
			f.Retain(home, tc.verb, tc.args, 3, tc.output)
			retained, err := responsebound.ReadRetained(home, f.Root)
			if err != nil {
				t.Fatal(err)
			}
			succeed(t, f, stageArgs(t, tc.requirement, tc.extra...))
			got, err := staged(t, f, home, tc.requirement)
			want := responsebound.Staged{Retained: retained, Mutant: tc.mutant}
			if err != nil || !reflect.DeepEqual(got, want) || string(got.Output) != tc.output || got.Head != f.Tip() {
				t.Fatalf("staged = (%+v, %v), want %+v", got, err, want)
			}
			requireNoWorktreeWrite(t, f)
		})
	}
}

// RT16: a second stage of one requirement replaces the first staged bytes.
func TestExcerptReplacesTheFirstStage(t *testing.T) {
	f, home := recordtest.Staging(t)
	args := []string{"--run", "TestA|TestB"}
	f.Retain(home, "test", args, 0, "first\n")
	succeed(t, f, stageArgs(t, "quoted"))
	f.Retain(home, "test", args, 1, "second\n")
	succeed(t, f, stageArgs(t, "quoted"))
	if got, err := staged(t, f, home, "quoted"); err != nil || string(got.Output) != "second\n" || got.Exit != 1 {
		t.Fatalf("staged = (%q, exit %d, %v), want the second response", got.Output, got.Exit, err)
	}
}

// RT11, RT12, RT13, RT14, RT15, RT17: each refusal stages nothing and writes no worktree file.
func TestExcerptRefusals(t *testing.T) {
	quoted := []string{"--run", "TestA|TestB"}
	for _, tc := range []struct {
		name, requirement string
		retain            func(f *recordtest.Fixture, home string)
		extra             []string
		want              string
	}{
		{"no retained response", "quoted", nil, nil, "this worktree has no retained response; run bench test --run 'TestA|TestB' first"},
		{"differing arguments", "quoted", func(f *recordtest.Fixture, home string) {
			f.Retain(home, "test", []string{"--run", "TestA"}, 0, "ok\n")
		}, nil,
			"the retained call bench test --run TestA differs from the planned command bench test --run 'TestA|TestB'"},
		{"probe without bench probe", "tests", func(f *recordtest.Fixture, home string) {
			f.Retain(home, "test", []string{"--package", "./internal/example"}, 0, answer("bit", "yes", 0, "1", report))
		}, nil, "requirement tests plans a probe; the retained call is bench test, not bench probe"},
		{"mutant set without --mutant", "tests", func(f *recordtest.Fixture, home string) { f.Retain(home, "probe", nil, 0, setAnswer) }, nil,
			"the retained response is a mutant-set answer; name one row with --mutant <n>"},
		{"unplanned requirement", "typo", func(f *recordtest.Fixture, home string) { f.Retain(home, "test", quoted, 0, "ok\n") }, nil,
			"requirement typo is not planned; the planned requirements are tests, quoted"},
		{"over the excerpt limit", "quoted", func(f *recordtest.Fixture, home string) {
			f.Retain(home, "test", quoted, 0, strings.Repeat("x", int(bounds.ControlRecordLimit)+1))
		}, nil, "the retained response is unreadable: read limit exceeded"},
		{"mutant without a probe", "quoted", func(f *recordtest.Fixture, home string) { f.Retain(home, "test", quoted, 0, "ok\n") }, []string{"--mutant", "1"},
			"--mutant needs a probe requirement; requirement quoted plans no probe"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, home := recordtest.Staging(t)
			if tc.retain != nil {
				tc.retain(f, home)
			}
			out, code := recordcmd.Command(f.Root, stageArgs(t, tc.requirement, tc.extra...))
			if code != 1 || !strings.Contains(out, tc.want) {
				t.Fatalf("excerpt = exit %d, output %q; want exit 1 naming %q", code, out, tc.want)
			}
			if got, err := staged(t, f, home, tc.requirement); !errors.Is(err, fs.ErrNotExist) {
				t.Fatalf("staged after the refusal = (%+v, %v), want none", got, err)
			}
			requireNoWorktreeWrite(t, f)
		})
	}
}

// RT58: a planned command that is not one lexed `bench` command of words alone refuses, even
// when the retained call equals its words after the first.
func TestExcerptRefusesAPlannedCommandThatIsNotOneBenchCommand(t *testing.T) {
	for name, command := range map[string]string{
		"foreign executable": "go test ./...",
		"pipe":               "bench test | tail -5",
		"unlexed":            "bench test 'open",
	} {
		t.Run(name, func(t *testing.T) {
			f, home := recordtest.Staging(t)
			f.Plan.Chunks[0].Verification[1].Command = command
			f.RewritePlan()
			words := strings.Fields(command)
			f.Retain(home, words[1], words[2:], 0, "ok\n")
			out, code := recordcmd.Command(f.Root, stageArgs(t, "quoted"))
			if want := "is not one lexed bench command of words alone"; code != 1 || !strings.Contains(out, want) {
				t.Fatalf("excerpt = exit %d, output %q; want exit 1 naming %q", code, out, want)
			}
		})
	}
}

// RT59: a requirement ID of ../x stages one file, named by the digest of the slug, the chunk
// ID, and the requirement ID, in the assignment scope beside the retained response, and
// nowhere else.
func TestExcerptKeysAHostileRequirementInsideTheScope(t *testing.T) {
	f, home := recordtest.Staging(t)
	f.Plan.Chunks[0].Verification[1].ID = "../x"
	f.RewritePlan()
	f.Retain(home, "test", []string{"--run", "TestA|TestB"}, 0, "ok\n")
	before := responseboundtest.StoreFiles(t, home)
	if len(before) != 1 {
		t.Fatalf("store files before the stage = %q, want the retained response alone", before)
	}
	succeed(t, f, stageArgs(t, "../x"))
	sum := sha256.Sum256([]byte(slug(t) + "\x00" + "1" + "\x00" + "../x"))
	want := []string{filepath.Join(filepath.Dir(before[0]), hex.EncodeToString(sum[:])), before[0]}
	if got := responseboundtest.StoreFiles(t, home); !reflect.DeepEqual(got, want) {
		t.Fatalf("store files after the stage = %q, want %q", got, want)
	}
	requireNoWorktreeWrite(t, f)
}
