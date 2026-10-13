package probe

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/durablefile"
	"github.com/gibbonmi/bench/internal/durablefile/durablefiletest"
	"github.com/gibbonmi/bench/internal/probe/probetest"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// The entries most set rows run: two swaps, the mutant bytes each swap writes over clamp.go,
// and an omission that clamp.go does not match.
var (
	guardSwap    = []string{"--swap", "n < 0", "--with", "n > 0"}
	returnSwap   = []string{"--swap", "return n", "--with", "return -n"}
	guardMutant  = strings.Replace(clampSource, "n < 0", "n > 0", 1)
	returnMutant = strings.Replace(clampSource, "return n", "return -n", 1)
	missEntry    = []string{"--omit", "absent"}
)

// bite and pass are the canned mutated runs: one fails TestClampNegative, and one passes.
// alreadyRed is a baseline that fails TestClampNegative before any write.
var (
	bite       = probetest.Start{Events: cannedFailure("caught"), Exit: 1}
	pass       = probetest.Start{Events: cannedPass}
	alreadyRed = probetest.Start{Events: cannedFailure("already red"), Exit: 1}
)

// runBinaryCases are the two run-binary sources of a set: an inherited binary builds nothing,
// and an owned one builds once.
var runBinaryCases = []struct {
	name   string
	owned  bool
	builds int
}{{"inherited", false, 0}, {"owned", true, 1}}

// clampSelection is the selection that most set rows run.
var clampSelection = []string{"--package", "./", "--run", "^TestClampNegative$"}

// mutantsFile writes entries as one mutants file and answers its path.
func mutantsFile(t *testing.T, entries ...[]string) string {
	t.Helper()
	body, err := json.Marshal(entries)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(scratchDir(t), "mutants.json")
	writeFixtureFile(t, path, string(body), 0o644)
	return path
}

// setArgs is the set over clamp.go that reads path, under selection or clampSelection.
func setArgs(path string, selection ...string) []string {
	if len(selection) == 0 {
		selection = clampSelection
	}
	return slices.Concat([]string{"clamp.go", "--mutants", path}, selection)
}

// installSet installs the sequenced canned `go`: a passing baseline, then mutated in order.
func installSet(t *testing.T, f *fixture, mutated ...probetest.Start) {
	t.Helper()
	probetest.InstallSequence(t, f.marker, pass, mutated...)
}

// setRows derives the expected mutants table through the encoder the verb renders with.
func setRows(t *testing.T, rows ...[]any) string {
	t.Helper()
	out, err := toon.TableTyped("mutants", mutantFieldsOfSet, rows)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// biteRow is the row of entry n whose run failed TestClampNegative.
func biteRow(n int, mutation string) []any {
	return []any{n, "bit", mutation, "failed", 1, "TestClampNegative", 1, 1, 0, nil, "yes"}
}

// missRow is the row of entry n, a missEntry, which the subject does not match.
func missRow(n int) []any {
	return []any{n, "invalid", "omit", "substring-miss", 0, nil, 0, nil, 1, nil, "untouched"}
}

// untouchedRow is the row of entry n that ended before any write, with cause.
func untouchedRow(n int, cause string) []any {
	return []any{n, "invalid", "swap", cause, 0, nil, 0, nil, 1, nil, "untouched"}
}

// redRow is the row of entry n, a swap, after an alreadyRed baseline.
func redRow(n int) []any {
	return []any{n, "invalid", "swap", "baseline-failed", 1, "TestClampNegative", 0, nil, 1, nil, "untouched"}
}

// restoreFailedRow is the row of entry n, a swap that bit and whose restore failed.
func restoreFailedRow(n int) []any {
	return []any{n, "restore-failed", "swap", "failed", 1, "TestClampNegative", 1, 1, 2, nil, "no"}
}

// requireTrailingBytes asserts the execution row of out, with its trailing_bytes cell derived
// from the bytes that follow that row.
func requireTrailingBytes(t *testing.T, out string, command int) {
	t.Helper()
	header, _, _ := strings.Cut(setExecution(t, command, 0), "\n")
	_, after, found := strings.Cut(out, "\n"+header+"\n")
	row, suffix, _ := strings.Cut(after, "\n")
	if !found {
		t.Fatalf("stdout = %q, want an execution row", out)
	}
	if got, want := header+"\n"+row+"\n", setExecution(t, command, len(suffix)); got != want {
		t.Fatalf("execution row = %q, want %q", got, want)
	}
}

// requirePreserved asserts the preserved table of out, which names the one kept copy, and that
// copy's start bytes.
func requirePreserved(t *testing.T, f *fixture, out string) {
	t.Helper()
	kept := preservedCopy(t, f)
	if !strings.Contains(out, "preserved[1]{path,reason}:\n  "+kept) {
		t.Fatalf("stdout = %q, want the preserved row naming %s", out, kept)
	}
	durablefiletest.RequireBytes(t, kept, clampSource)
}

// setSelection derives the expected selection table of clampSelection.
func setSelection(t *testing.T, baseline string, builds int) string {
	t.Helper()
	out, err := toon.TableTyped("selection", setSelectionFields, [][]any{{"package", "./", "^TestClampNegative$", baseline, builds}})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// setExecution derives the expected execution table.
func setExecution(t *testing.T, command, trailing int) string {
	t.Helper()
	out, err := toon.TableTyped("probe_set_execution", setExecutionFields, [][]any{{setSchema, command, trailing}})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// runSetWith runs the set argv over f with each subject replacement through replacer.
func runSetWith(t *testing.T, f *fixture, replacer durablefile.Replacer, argv []string) (string, int) {
	t.Helper()
	parsed, line, _ := usage.Parse(grammar, argv)
	if line != "" {
		t.Fatalf("parse refused %q: %q", argv, line)
	}
	return runWith(replacer, f.root, parsed)
}

func requireStarts(t *testing.T, f *fixture, want int) {
	t.Helper()
	if got := startCount(t, f); got != want {
		t.Fatalf("Go test starts = %d, want %d", got, want)
	}
}

// entryRefused is the answer to a set whose entry n does not match the mutation grammar.
func entryRefused(n string) string {
	return "error: probe mutant set entry refused — entry " + n + " does not match the mutation grammar\n" + grammar.Help + "\n"
}

// An entry that names a selection token, --mutants, a positional, or no mutation at all
// refuses with its entry number before any run.
func TestMutantSetEntryRefusesSelectionTokens(t *testing.T) {
	f := refusalFixture(t)
	for _, tc := range []struct {
		name  string
		entry []string
	}{
		{"package", append(slices.Clone(guardSwap), "--package", "./")},
		{"check", append(slices.Clone(guardSwap), "--check", "prose")},
		{"run", append(slices.Clone(guardSwap), "--run", "^X$")},
		{"full", append(slices.Clone(guardSwap), "--full")},
		{"mutants", append(slices.Clone(guardSwap), "--mutants", "other.json")},
		{"positional", append([]string{"clamp.go"}, guardSwap...)},
		{"empty", []string{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out, code := runProbe(t, setArgs(mutantsFile(t, returnSwap, tc.entry))...)
			requireRefused(t, f, out, code, entryRefused("2"), 2)
		})
	}
}

// Under a prose selection, an entry's --expect refuses, because a prose run names no test.
func TestMutantSetEntryRefusesProseExpect(t *testing.T) {
	f := refusalFixture(t)
	path := mutantsFile(t, []string{"--omit", "n < 0"}, []string{"--omit", "n < 0", "--expect", "TestClampNegative"})
	out, code := runProbe(t, "clamp.go", "--mutants", path, "--check", "prose")
	requireRefused(t, f, out, code, entryRefused("2"), 2)
}

// --mutants beside any top-level mutation flag or --expect is a usage refusal, so no flag
// joins every entry in silence.
func TestMutantsExcludesTopLevelMutation(t *testing.T) {
	f := refusalFixture(t)
	path := mutantsFile(t, guardSwap)
	for _, extra := range [][]string{
		{"--swap", "n < 0"}, {"--omit", "n < 0"}, {"--unwrap", "f(n)"}, {"--insert-before", "return"},
		{"--omit-file", path}, {"--delete"}, {"--with", "x"}, {"--with-file", path}, {"--occurrence", "1"},
		{"--expect", "TestClampNegative"},
	} {
		t.Run(extra[0], func(t *testing.T) {
			out, code := runProbe(t, slices.Concat([]string{"clamp.go", "--mutants", path}, extra, clampSelection)...)
			requireRefused(t, f, out, code, grammar.Help+"\n", 2)
		})
	}
}

// The mutants file reads through the bounded mutation-file reader, so a hostile path refuses
// with that reader's line before any write.
func TestMutantsFileRefusesHostilePaths(t *testing.T) {
	f := refusalFixture(t)
	dir := scratchDir(t)
	valid := mutantsFile(t, guardSwap)
	paths := map[string]string{
		"absent":         filepath.Join(dir, "absent.json"),
		"a symlink":      filepath.Join(dir, "link.json"),
		"a directory":    dir,
		"a special file": filepath.Join(dir, "fifo"),
		"oversized":      filepath.Join(dir, "oversized.json"),
		"empty":          filepath.Join(dir, "empty.json"),
	}
	if err := os.Symlink(valid, paths["a symlink"]); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(paths["a special file"], 0o600); err != nil {
		t.Fatal(err)
	}
	writeFixtureFile(t, paths["oversized"], strings.Repeat(" ", int(bounds.ControlRecordLimit)+1), 0o644)
	writeFixtureFile(t, paths["empty"], "", 0o644)
	for reason, path := range paths {
		t.Run(reason, func(t *testing.T) {
			out, code := runProbe(t, setArgs(path)...)
			requireRefused(t, f, out, code, "error: probe mutation file unavailable — "+path+" is "+reason+"\n", 1)
		})
	}
}

// A file that is not one strict JSON list of string lists refuses as malformed, and an empty
// list refuses as empty, before any run.
func TestMutantsFileRefusesMalformedJSON(t *testing.T) {
	f := refusalFixture(t)
	for name, body := range map[string]string{
		"object":         `{}`,
		"string entry":   `["--swap"]`,
		"number token":   `[["--swap", 1]]`,
		"null list":      `null`,
		"null entry":     `[["--omit", "n < 0"], null]`,
		"trailing value": `[["--omit", "n < 0"]] []`,
	} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(scratchDir(t), "mutants.json")
			writeFixtureFile(t, path, body, 0o644)
			out, code := runProbe(t, setArgs(path)...)
			prefix := "error: probe mutant set malformed — " + path + ": "
			if !strings.HasPrefix(out, prefix) || strings.Count(out, "\n") != 1 || code != 1 {
				t.Fatalf("refusal = (%q, %d), want one line that starts %q and 1", out, code, prefix)
			}
			requireNoRunChild(t, f)
			requireHomeEmpty(t, f)
		})
	}
	t.Run("empty list", func(t *testing.T) {
		path := filepath.Join(scratchDir(t), "mutants.json")
		writeFixtureFile(t, path, `[]`, 0o644)
		out, code := runProbe(t, setArgs(path)...)
		requireRefused(t, f, out, code, "error: probe mutant set empty — "+path+" holds no entry\n", 1)
	})
}

// A set holds at most 64 entries: 65 refuse with the cap, and 64 run.
func TestMutantSetEntryCap(t *testing.T) {
	f := refusalFixture(t)
	misses := func(n int) [][]string {
		entries := make([][]string, n)
		for i := range entries {
			entries[i] = missEntry
		}
		return entries
	}
	over := mutantsFile(t, misses(65)...)
	out, code := runProbe(t, setArgs(over)...)
	requireRefused(t, f, out, code, "error: probe mutant set too large — "+over+" holds 65 entries, the cap is 64\n", 1)
	out, code = runProbe(t, setArgs(mutantsFile(t, misses(64)...))...)
	if !strings.HasPrefix(out, "mutants[64]{") || code != 1 {
		t.Fatalf("stdout = (%q, %d), want 64 rows and 1", out, code)
	}
	requireNoRunChild(t, f)
}

// An entry's mutation refusal keeps its single-form title, and its detail names the entry.
func TestMutantSetEntryRefusalNamesTheEntry(t *testing.T) {
	f := refusalFixture(t)
	absent := filepath.Join(scratchDir(t), "absent")
	for _, tc := range []struct {
		name  string
		entry []string
		want  string
	}{
		{"ambiguous", []string{"--swap", "return", "--with", "x"}, "error: probe mutation ambiguous — entry 2: the old string matches 2 times, want exactly 1\n"},
		{"absent text file", []string{"--swap", "n < 0", "--with-file", absent}, "error: probe mutation file unavailable — entry 2: " + absent + " is absent\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out, code := runProbe(t, setArgs(mutantsFile(t, guardSwap, tc.entry))...)
			requireRefused(t, f, out, code, tc.want, 1)
		})
	}
}
