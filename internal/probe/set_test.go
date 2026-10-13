package probe

import (
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
	"github.com/gibbonmi/bench/internal/probe/probetest"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/testreport"
	"github.com/gibbonmi/bench/internal/toon"
)

// Each entry of a set prints its own row, and two biting entries take three Go starts.
func TestMutantSetGradesEachEntry(t *testing.T) {
	f := newFixture(t)
	installSet(t, f, bite, bite)
	out, code := runProbe(t, setArgs(mutantsFile(t, guardSwap, returnSwap))...)
	requireRow(t, out, code, 0, setRows(t, biteRow(1, "swap"), biteRow(2, "swap")))
	requireStarts(t, f, 3)
	requireSubjectBytes(t, f, clampSource)
}

// An entry means what the same argv means alone: the set's run sees the single form's bytes.
func TestMutantSetEntryMatchesTheSingleForm(t *testing.T) {
	f := newFixture(t)
	replacement := filepath.Join(scratchDir(t), "replacement")
	writeFixtureFile(t, replacement, "n >= 0", 0o644)
	entries := [][]string{
		{"--swap", "return", "--with", "return -", "--occurrence", "2"},
		{"--swap", "n < 0", "--with-file", replacement},
	}
	copies := scratchDir(t)
	hooks := make([]probetest.Start, len(entries))
	for k, entry := range entries {
		copied := filepath.Join(copies, strconv.Itoa(k))
		installStubStarts(t, f, pass, probetest.Start{Hook: "cp " + sanitize.ShellQuote(f.subject) + " " + sanitize.ShellQuote(copied) + "\n", Events: bite.Events, Exit: 1})
		if out, code := runProbe(t, slices.Concat([]string{"clamp.go"}, entry, clampSelection)...); code != 0 {
			t.Fatalf("single form of entry %d = (%q, %d), want a bite", k+1, out, code)
		}
		single, err := os.ReadFile(copied)
		if err != nil || string(single) == clampSource {
			t.Fatalf("single form of entry %d wrote (%q, %v), want a mutant", k+1, single, err)
		}
		hooks[k] = probetest.Start{Hook: mutantHook(t, f, string(single)), Events: bite.Events, Exit: 1}
	}
	installSet(t, f, hooks...)
	out, code := runProbe(t, setArgs(mutantsFile(t, entries...))...)
	requireRow(t, out, code, 0, setRows(t, biteRow(1, "swap"), biteRow(2, "swap")))
}

// Each entry grades its own --expect, so one missed witness changes only its own row.
func TestMutantSetExpectIsPerEntry(t *testing.T) {
	f := newFixture(t)
	installSet(t, f, bite, bite)
	path := mutantsFile(t, append(slices.Clone(guardSwap), "--expect", "TestClampPositive"), append(slices.Clone(returnSwap), "--expect", "TestClampNegative"))
	out, code := runProbe(t, setArgs(path)...)
	missed := []any{1, "invalid", "swap", "expect-missed", 1, "TestClampNegative", 1, 1, 1, nil, "yes"}
	requireRow(t, out, code, 1, setRows(t, missed, biteRow(2, "swap")))
}

// An entry's repeated swaps make one combined mutant, as the same argv does alone.
func TestMutantSetEntryTakesACombinedMutant(t *testing.T) {
	f := newFixture(t)
	both := strings.Replace(guardMutant, "return n", "return -n", 1)
	installSet(t, f, probetest.Start{Hook: mutantHook(t, f, both), Events: bite.Events, Exit: 1})
	out, code := runProbe(t, setArgs(mutantsFile(t, slices.Concat(guardSwap, returnSwap)))...)
	requireRow(t, out, code, 0, setRows(t, biteRow(1, "swap+swap")))
}

// An unmatched entry reads substring-miss while the other entries run.
func TestMutantSetMissIsPerEntry(t *testing.T) {
	f := newFixture(t)
	installSet(t, f, bite)
	out, code := runProbe(t, setArgs(mutantsFile(t, missEntry, guardSwap))...)
	requireRow(t, out, code, 1, setRows(t, missRow(1), biteRow(2, "swap")))
	requireStarts(t, f, 2)
}

// A set whose every entry misses answers before the gate lock: no lock, no run, no build.
func TestMutantSetOfMissesStartsNoRun(t *testing.T) {
	f := newFixture(t)
	installSet(t, f)
	out, code := runProbe(t, setArgs(mutantsFile(t, missEntry, missEntry))...)
	want := setRows(t, missRow(1), missRow(2)) + setSelection(t, "", 0) + setExecution(t, 1, 0)
	requireRefused(t, f, out, code, want, 1)
}

// One baseline serves the whole set, so three entries take four Go starts.
func TestMutantSetSharesOneBaseline(t *testing.T) {
	f := newFixture(t)
	installSet(t, f, bite, bite, bite)
	out, code := runProbe(t, setArgs(mutantsFile(t, guardSwap, returnSwap, []string{"--omit", "return 0"}))...)
	if code != 0 {
		t.Fatalf("exit = %d, want 0\n%s", code, out)
	}
	requireStarts(t, f, 4)
}

// A failed baseline gives every row its kind before any write, and its report follows. The
// builds cell counts the baseline's own build.
func TestMutantSetBaselineEndsTheSet(t *testing.T) {
	for _, tc := range runBinaryCases {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t)
			if tc.owned {
				buildEachRun(t, f, "clamp.go")
			}
			probetest.InstallSequence(t, f.marker, alreadyRed, bite, bite)
			out, code := runProbe(t, setArgs(mutantsFile(t, guardSwap, returnSwap))...)
			requireRow(t, out, code, 1, setRows(t, redRow(1), redRow(2))+setSelection(t, "failed", tc.builds))
			if !strings.Contains(out, "already red") {
				t.Fatalf("stdout = %q, want the baseline's own report", out)
			}
			requireTrailingBytes(t, out, 1)
			requireStarts(t, f, 1)
			requireSubjectBytes(t, f, clampSource)
			requireHomeEmpty(t, f)
		})
	}
}

// Each mutant runs alone: its write replaces the start bytes, never an earlier mutant, and its
// run sees only its own bytes.
func TestMutantSetRunsEachMutantAlone(t *testing.T) {
	f := newFixture(t)
	installSet(t, f,
		probetest.Start{Hook: mutantHook(t, f, guardMutant), Events: bite.Events, Exit: 1},
		probetest.Start{Hook: mutantHook(t, f, returnMutant), Events: bite.Events, Exit: 1})
	native := durablefile.Native()
	replacer := native
	var replaced []string
	replacer.Rename = func(oldpath, newpath string) error {
		body, err := os.ReadFile(newpath)
		if err != nil {
			t.Fatal(err)
		}
		replaced = append(replaced, string(body))
		return native.Rename(oldpath, newpath)
	}
	out, code := runSetWith(t, f, replacer, setArgs(mutantsFile(t, guardSwap, returnSwap)))
	requireRow(t, out, code, 0, setRows(t, biteRow(1, "swap"), biteRow(2, "swap")))
	if want := []string{clampSource, guardMutant, clampSource, returnMutant}; !slices.Equal(replaced, want) {
		t.Fatalf("each replacement replaced %q, want %q", replaced, want)
	}
}

// The subject keeps its start bytes and mode after the set, and no copy stays in the home.
func TestMutantSetRestoresAfterEachMutant(t *testing.T) {
	f := newFixture(t)
	if err := os.Chmod(f.subject, 0o640); err != nil {
		t.Fatal(err)
	}
	installSet(t, f, bite, bite)
	if out, code := runProbe(t, setArgs(mutantsFile(t, guardSwap, returnSwap))...); code != 0 {
		t.Fatalf("exit = %d, want 0\n%s", code, out)
	}
	requireSubjectBytes(t, f, clampSource)
	if got := permissions(t, f.subject); got != 0o640 {
		t.Fatalf("subject mode = %o, want 640", got)
	}
	requireHomeEmpty(t, f)
}

// A failed restore after mutant 1 keeps the copy and stops the set at exit 2.
func TestMutantSetStopsOnARestoreFailure(t *testing.T) {
	f := newFixture(t)
	installSet(t, f, bite, bite)
	router := subjectFaults()
	router.Fail[restoreReplacement] = map[string]error{"rename": os.ErrPermission}
	out, code := runSetWith(t, f, router.Replacer(), setArgs(mutantsFile(t, guardSwap, returnSwap)))
	requireRow(t, out, code, 2, setRows(t, restoreFailedRow(1), untouchedRow(2, "set-stopped")))
	requirePreserved(t, f, out)
	requireTrailingBytes(t, out, 2)
	requireStarts(t, f, 2)
}

// A failed write of mutant 1 stops the set. An unpublished failure leaves the subject
// untouched, and a published one restores it. A failed restore of a published failure keeps
// the copy at exit 2.
func TestMutantSetStopsOnAMutationFailure(t *testing.T) {
	published := map[string]error{"directory-sync": os.ErrPermission}
	failed := func(verdict string, exit int, restored string) []any {
		return []any{1, verdict, "swap", "mutation-failed", 0, nil, 0, nil, exit, nil, restored}
	}
	for _, tc := range []struct {
		name    string
		fail    map[durablefiletest.Fault]map[string]error
		row     []any
		command int
	}{
		{"rename", map[durablefiletest.Fault]map[string]error{mutationReplacement: {"rename": os.ErrPermission}}, failed("invalid", 1, "untouched"), 1},
		{"directory-sync", map[durablefiletest.Fault]map[string]error{mutationReplacement: published}, failed("invalid", 1, "yes"), 1},
		{"failed restore", map[durablefiletest.Fault]map[string]error{mutationReplacement: published, restoreReplacement: {"rename": os.ErrPermission}}, failed("restore-failed", 2, "no"), 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t)
			installSet(t, f, bite, bite)
			router := subjectFaults()
			router.Fail = tc.fail
			out, code := runSetWith(t, f, router.Replacer(), setArgs(mutantsFile(t, guardSwap, returnSwap)))
			requireRow(t, out, code, tc.command, setRows(t, tc.row, untouchedRow(2, "set-stopped")))
			requireStarts(t, f, 1)
			if tc.command == 2 {
				requirePreserved(t, f, out)
				return
			}
			requireSubjectBytes(t, f, clampSource)
			requireHomeEmpty(t, f)
		})
	}
}

// A failed preservation stops the set before mutant 1's write, and the answer prints no
// failure text.
func TestMutantSetStopsOnAPreserveFailure(t *testing.T) {
	f := newFixture(t)
	writeFixtureFile(t, filepath.Join(f.home, "probe"), "not a directory\n", 0o644)
	installSet(t, f, bite, bite)
	out, code := runProbe(t, setArgs(mutantsFile(t, guardSwap, returnSwap))...)
	want := setRows(t, untouchedRow(1, "preserve-failed"), untouchedRow(2, "set-stopped")) + setSelection(t, "passed", 0) + setExecution(t, 1, 0)
	if out != want || code != 1 {
		t.Fatalf("stdout = (%q, %d), want (%q, 1)", out, code, want)
	}
	requireStarts(t, f, 1)
	requireSubjectBytes(t, f, clampSource)
}

// A mutant outside the run binary's closure reuses the baseline's build.
func TestMutantSetReusesTheBuildOutsideTheClosure(t *testing.T) {
	f := newFixture(t)
	records := buildEachRun(t, f, "clamp.go")
	installSet(t, f, bite, bite)
	out, code := runProbe(t, setArgs(mutantsFile(t, guardSwap, returnSwap))...)
	requireRow(t, out, code, 0, setRows(t, biteRow(1, "swap"), biteRow(2, "swap")))
	if got := buildRecords(t, records); !slices.Equal(got, []string{clampSource}) {
		t.Fatalf("build script reads = %q, want one build over the start bytes", got)
	}
}

// The builds cell counts the builds that the set owned: none over an inherited binary.
func TestMutantSetReportsTheBuildCount(t *testing.T) {
	for _, tc := range runBinaryCases {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t)
			if tc.owned {
				buildEachRun(t, f, "clamp.go")
			}
			installSet(t, f, bite, bite)
			out, _ := runProbe(t, setArgs(mutantsFile(t, guardSwap, returnSwap))...)
			if want := setSelection(t, "passed", tc.builds); !strings.Contains(out, want) {
				t.Fatalf("stdout = %q, want the selection row %q", out, want)
			}
		})
	}
}

// A prose set selects no run binary, so it builds nothing.
func TestMutantSetOfProseBuildsNothing(t *testing.T) {
	f := newFixture(t)
	trackProseNote(t, f)
	records := buildEachRun(t, f, "note.md")
	path := mutantsFile(t, []string{"--swap", shortLine, "--with", longSentence}, []string{"--swap", "# Note", "--with", "# A note"})
	out, _ := runProbe(t, "note.md", "--mutants", path, "--check", "prose")
	selection, err := toon.TableTyped("selection", setSelectionFields, [][]any{{"check", "prose", testreport.AllTests, "passed", 0}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, selection) {
		t.Fatalf("stdout = %q, want the selection row %q", out, selection)
	}
	if got := buildRecords(t, records); len(got) != 0 {
		t.Fatalf("builds = %d, want none", len(got))
	}
}

// The answer opens with the designed mutants header, and one selection row follows the rows.
func TestMutantSetRowsKeepTheirSchema(t *testing.T) {
	f := newFixture(t)
	installSet(t, f, bite, bite)
	out, _ := runProbe(t, setArgs(mutantsFile(t, guardSwap, returnSwap))...)
	lines := strings.Split(out, "\n")
	if len(lines) < 7 {
		t.Fatalf("stdout = %q, want the three tables", out)
	}
	want := []string{
		"mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:",
		"selection[1]{form,target,run,baseline,builds}:",
		"probe_set_execution[1]{schema,command_exit_code,trailing_bytes}:",
	}
	if got := []string{lines[0], lines[3], lines[5]}; !slices.Equal(got, want) {
		t.Fatalf("table headers = %q, want %q", got, want)
	}
	if got, want := lines[4], "  package,./,^TestClampNegative$,passed,0"; got != want {
		t.Fatalf("selection row = %q, want %q", got, want)
	}
}

// Each row's exit follows the verdict map. The command exit is 0 when every row bit, 2 when a
// restore failed, and 1 otherwise.
func TestMutantSetCommandExit(t *testing.T) {
	silent := []any{2, "silent", "swap", "passed", 0, nil, 1, 0, 1, nil, "yes"}
	for _, tc := range []struct {
		name    string
		second  probetest.Start
		fault   map[durablefiletest.Fault]map[string]error
		row     []any
		command int
	}{
		{"bit and silent", pass, nil, silent, 1},
		{"two bites", bite, nil, biteRow(2, "swap"), 0},
		{"bit and restore-failed", bite, map[durablefiletest.Fault]map[string]error{{Destination: "subject", Ordinal: 4}: {"rename": os.ErrPermission}}, restoreFailedRow(2), 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t)
			installSet(t, f, bite, tc.second)
			router := subjectFaults()
			if tc.fault != nil {
				router.Fail = tc.fault
			}
			out, code := runSetWith(t, f, router.Replacer(), setArgs(mutantsFile(t, guardSwap, returnSwap)))
			requireRow(t, out, code, tc.command, setRows(t, biteRow(1, "swap"), tc.row))
		})
	}
}

// Each started mutant's report follows the execution row under its mutant number, in entry
// order, and a missed entry prints none, the last entry included. trailing_bytes counts every byte after that row.
func TestMutantSetPrintsEachReport(t *testing.T) {
	f := newFixture(t)
	first, third := cannedFailure("first caught"), cannedFailure("third caught")
	installSet(t, f, probetest.Start{Events: first, Exit: 1}, probetest.Start{Events: third, Exit: 1})
	out, code := runProbe(t, setArgs(mutantsFile(t, guardSwap, missEntry, returnSwap, missEntry))...)
	reports := map[string]string{}
	for _, events := range []string{first, third} {
		installStubStarts(t, f, probetest.Start{Events: events, Exit: 1}, probetest.Start{Events: events, Exit: 1})
		reports[events], _ = testreport.Command(f.root, clampSelection)
	}
	suffix := "mutant_report[1]{mutant}:\n  1\n" + reports[first] + "mutant_report[1]{mutant}:\n  3\n" + reports[third]
	head := setRows(t, biteRow(1, "swap"), missRow(2), biteRow(3, "swap"), missRow(4)) + setSelection(t, "passed", 0) + setExecution(t, 1, len(suffix))
	if want := head + suffix; out != want || code != 1 {
		t.Fatalf("stdout = (%q, %d), want (%q, 1)", out, code, want)
	}
}

// unusedOs is a mutated package whose build output starts with the compiler line.
const unusedOs = `{"Action":"build-output","ImportPath":"probefixture","Output":"clamp.go:3:2: \"os\" imported and not used\n"}
{"Action":"build-fail","ImportPath":"probefixture"}`

// A mutant that fails to build names its first build line in its own row only.
func TestMutantSetNamesABuildFailure(t *testing.T) {
	f := newFixture(t)
	installSet(t, f, probetest.Start{Events: unusedOs, Exit: 1}, bite)
	out, code := runProbe(t, setArgs(mutantsFile(t, guardSwap, returnSwap), "--package", "./")...)
	broken := []any{1, "invalid", "swap", "build-failed", 0, nil, 0, 1, 1, `clamp.go:3:2: "os" imported and not used`, "yes"}
	requireRow(t, out, code, 1, setRows(t, broken, biteRow(2, "swap")))
}

// A set passes every repeated --package value, in argv order, to each of its Go starts.
func TestMutantSetTakesRepeatedPackages(t *testing.T) {
	f := newFixture(t)
	argvs := scratchDir(t)
	record := "printf '%s\\n' \"$@\" > " + sanitize.ShellQuote(argvs) + "/$n\n"
	probetest.InstallSequence(t, f.marker, probetest.Start{Hook: record, Events: cannedPass},
		probetest.Start{Hook: record, Events: bite.Events, Exit: 1}, probetest.Start{Hook: record, Events: bite.Events, Exit: 1})
	selection := []string{"--package", "./clamp", "--package", "./other", "--run", "^TestClampNegative$"}
	out, code := runProbe(t, setArgs(mutantsFile(t, guardSwap, returnSwap), selection...)...)
	if code != 0 {
		t.Fatalf("exit = %d, want 0\n%s", code, out)
	}
	starts, err := os.ReadDir(argvs)
	if err != nil || len(starts) != 3 {
		t.Fatalf("Go test starts = (%d, %v), want 3", len(starts), err)
	}
	tail := []string{"-run", "^TestClampNegative$", "./clamp", "./other"}
	for _, start := range starts {
		body, err := os.ReadFile(filepath.Join(argvs, start.Name()))
		if err != nil {
			t.Fatal(err)
		}
		args := strings.Split(strings.TrimSuffix(string(body), "\n"), "\n")
		if len(args) < len(tail) || !slices.Equal(args[len(args)-len(tail):], tail) {
			t.Fatalf("go argv of start %s = %q, want %q as its last arguments", start.Name(), args, tail)
		}
	}
}
