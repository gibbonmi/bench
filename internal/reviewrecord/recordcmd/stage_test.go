package recordcmd_test

import (
	"errors"
	"io/fs"
	"path/filepath"
	"reflect"
	"slices"
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

// setAnswer answers a mutant-set answer with one row, and that row bit. The probe runs in a scratch
// fixture and binds its own Bench home, so a test calls it before recordtest.Staging.
func setAnswer(t *testing.T) string {
	t.Helper()
	produce, _ := setProducer(t, linked(t, 1))
	data, _ := produce([][]string{bitingEntry}, bite)
	return data
}

// RT10: a retained `bench probe` answer of a probe requirement stages its bytes, exit, commit,
// and dirty paths, and a mutant-set answer stages with its row.
func TestExcerptStagesTheRetainedResponse(t *testing.T) {
	set := setAnswer(t)
	for _, tc := range []struct {
		name, requirement, verb string
		args                    []string
		output                  string
		extra                   []string
		mutant                  int
	}{
		{"probe", "tests", "probe", []string{"clamp.go"}, answer("bit", "yes", 0, "1", report), nil, 0},
		{"quoted command", "quoted", "probe", []string{"clamp.go", "--omit", "a|b"}, answer("bit", "yes", 0, "1", report), nil, 0},
		{"probe row", "tests", "probe", []string{"--mutants", "m.json"}, set, []string{"--mutant", "1"}, 1},
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
	first, second := answer("bit", "yes", 0, "1", report), answer("silent", "yes", 1, "0", report)
	f.Retain(home, "probe", nil, 0, first)
	succeed(t, f, stageArgs(t, "quoted"))
	f.Retain(home, "probe", nil, 1, second)
	succeed(t, f, stageArgs(t, "quoted"))
	if got, err := staged(t, f, home, "quoted"); err != nil || string(got.Output) != second || got.Exit != 1 {
		t.Fatalf("staged = (%q, exit %d, %v), want the second response", got.Output, got.Exit, err)
	}
}

// RT11, RT12, RT13, RT14, RT15, RT17: each refusal stages nothing and writes no worktree file. A
// plain chunk requirement refuses at the plan grade, so the excerpt form stages only probes.
func TestExcerptRefusals(t *testing.T) {
	bit, set := answer("bit", "yes", 0, "1", report), setAnswer(t)
	for _, tc := range []struct {
		name, requirement string
		retain            func(f *recordtest.Fixture, home string)
		extra             []string
		want              string
	}{
		{"no retained response", "quoted", nil, nil, "this worktree has no retained response; run bench probe clamp.go --omit 'a|b' --run 'TestA|TestB' first"},
		{"plain chunk requirement", "quoted", func(f *recordtest.Fixture, home string) {
			f.Plan.Chunks[0].Verification[1].Probe = ""
			f.WritePlan()
			f.Commit("plan a plain chunk requirement")
			f.Retain(home, "probe", nil, 0, bit)
		}, nil, "chunk 1: verification quoted has no probe"},
		{"probe without bench probe", "tests", func(f *recordtest.Fixture, home string) {
			f.Retain(home, "test", []string{"--package", "./internal/example"}, 0, answer("bit", "yes", 0, "1", report))
		}, nil, "requirement tests plans a probe; the retained call is bench test, not bench probe"},
		{"mutant set without --mutant", "tests", func(f *recordtest.Fixture, home string) { f.Retain(home, "probe", nil, 0, set) }, nil,
			"the retained response is a mutant-set answer; name one row with --mutant <n>"},
		{"unplanned requirement", "typo", func(f *recordtest.Fixture, home string) { f.Retain(home, "probe", nil, 0, bit) }, nil,
			"requirement typo is not planned; the planned requirements are tests, quoted"},
		{"over the excerpt limit", "quoted", func(f *recordtest.Fixture, home string) {
			f.Retain(home, "probe", nil, 0, strings.Repeat("x", int(bounds.ControlRecordLimit)+1))
		}, nil, "the retained response is unreadable: read limit exceeded"},
		{"empty", "quoted", func(f *recordtest.Fixture, home string) { f.Retain(home, "probe", nil, 0, "") }, nil,
			"the retained response is empty"},
		{"invalid UTF-8", "quoted", func(f *recordtest.Fixture, home string) { f.Retain(home, "probe", nil, 0, "\xff\n") }, nil,
			"the retained response is malformed"},
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

// RT59: a requirement ID of ../x stages one new file in the assignment scope beside the
// retained response, and nowhere else. No character of the ID forms the file name.
func TestExcerptKeysAHostileRequirementInsideTheScope(t *testing.T) {
	f, home := recordtest.Staging(t)
	f.Plan.Chunks[0].Verification[1].ID = "../x"
	f.RewritePlan()
	f.Retain(home, "probe", nil, 0, answer("bit", "yes", 0, "1", report))
	before := responseboundtest.StoreFiles(t, home)
	if len(before) != 1 {
		t.Fatalf("store files before the stage = %q, want the retained response alone", before)
	}
	succeed(t, f, stageArgs(t, "../x"))
	after := responseboundtest.StoreFiles(t, home)
	added := slices.DeleteFunc(slices.Clone(after), func(file string) bool { return file == before[0] })
	if len(after) != 2 || len(added) != 1 || filepath.Dir(added[0]) != filepath.Dir(before[0]) || strings.ContainsAny(filepath.Base(added[0]), "./x") {
		t.Fatalf("store files after the stage = %q, want one new file beside %s whose name holds no ID character", after, before[0])
	}
	requireNoWorktreeWrite(t, f)
}

// A chunk requirement whose planned command runs the gate is a plain requirement, so the ticket
// form refuses at the plan grade and writes no record.
func TestTicketRefusesAPlainGateRequirement(t *testing.T) {
	f, home := recordtest.Staging(t)
	base, _ := advance(f, "chunk 1")
	stageOwed(f, home, "1", "1.md")
	gate := &f.Plan.Chunks[0].Verification[1]
	gate.Command, gate.Probe = "bench gate --checkpoint "+recordtest.Spec+" --chunk 1", ""
	f.WritePlan()
	f.Commit("plan a plain gate requirement")
	out, code := recordcmd.Command(f.Root, ticketArgs(t, "1", "1.md", "g", "--base", base))
	if want := "chunk 1: verification quoted has no probe"; code != 1 || !strings.Contains(out, want) {
		t.Fatalf("ticket = exit %d, output %q; want exit 1 naming %q", code, out, want)
	}
	if record := recordBytes(t, f); record != nil {
		t.Fatalf("the refusal wrote the record:\n%s", record)
	}
}

// The excerpt form refuses a --mutant value that is not a positive decimal row with its usage.
func TestExcerptRefusesAnInvalidMutantRow(t *testing.T) {
	f, _ := recordtest.Staging(t)
	for _, row := range []string{"x", "0"} {
		out, code := recordcmd.Command(f.Root, stageArgs(t, "tests", "--mutant", row))
		if code != 2 || !strings.HasPrefix(out, "usage: bench record excerpt") {
			t.Errorf("--mutant %s = exit %d, output %q; want exit 2 with a usage: bench record excerpt line", row, code, out)
		}
	}
}
