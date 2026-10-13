package recordcmd_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/probe"
	"github.com/gibbonmi/bench/internal/probe/probetest"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordcmd"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// checkpoint adds one independent review for each axis to the chunk 1 record, commits it,
// and answers the chunk 1 checkpoint.
func checkpoint(t *testing.T, f *recordtest.Fixture) error {
	t.Helper()
	record := read(t, f)
	chunk := &record.Chunks[0]
	for _, axis := range rr.Axes() {
		chunk.Reviews = append(chunk.Reviews, rr.Review{Evidence: f.Evidence("review-"+axis, chunk.SourceDigest, "independent-review"), Axis: axis, Base: chunk.Base, Tip: chunk.Tip})
	}
	path := recordFile(t, f)
	document, err := rr.Render(bytesOf(t, path), record)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, document, 0o644); err != nil {
		t.Fatal(err)
	}
	f.Commit("record chunk 1")
	return rr.Check(f.Root, f.Tree(), f.Tip(), recordtest.Spec, "1", false)
}

// The probe answers below are authored apart from the producer, in the envelope that the
// spec fixes: the probe row, the selection row, and the execution row whose last cell
// counts the bytes after it.

// treeLead is a canonical identity row that tree-scoped stdout prints first.
const treeLead = "tree[1]{target,head,dirty}:\n  alpha,none,true\n"

// answer is one probe answer: verdict, restored cell, command exit, focused cell, and the
// suffix after the execution row.
func answer(verdict, restored string, command int, focused, suffix string) string {
	return "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  " + verdict + ",clamp.go,swap,failed,1,TestClampNegative," + restored + "\n" +
		"selection[1]{form,target,run,baseline,ran}:\n  package,./,^TestClampNegative$,passed,1\n" +
		"probe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2," +
		strconv.Itoa(command) + "," + focused + "," + strconv.Itoa(len(suffix)) + "\n" + suffix
}

const report = "packages[1]{package,outcome}:\n  probefixture,failed\n"

// importing edits a verification call to read its probe from the answer in data.
func importing(t *testing.T, data string, exit int) map[string]string {
	t.Helper()
	return map[string]string{"--excerpt": drop, "--probe-outcome": drop, "--probe-exit-code": drop, "--probe-restore": drop,
		"--probe-result": excerptFile(t, "probe.out", data), "--exit-code": strconv.Itoa(exit)}
}

// TestVerificationImportsProbeResult covers NR25, NR26, NR29, and NR30: an import fills the planned probe from the answer and retains the
// exact file bytes as its excerpt. A silent or unrestored answer records as it reads, and
// the checkpoint still refuses it.
func TestVerificationImportsProbeResult(t *testing.T) {
	for _, tc := range []struct {
		name, data string
		exit       int
		want       rr.Probe
		refusal    string
	}{
		{"bit with tree lead", treeLead + answer("bit", "yes", 0, "1", report), 0, rr.Probe{Outcome: "bit", ExitCode: 1, Restore: "pass"}, ""},
		{"silent", answer("silent", "yes", 1, "0", report), 1, rr.Probe{Outcome: "silent", ExitCode: 0, Restore: "pass"}, "verification tests is incomplete, failed"},
		{"restore failed", answer("restore-failed", "no", 2, "1", "preserved[1]{path,reason}:\n  /h/clamp.go,denied\n"+report), 2,
			rr.Probe{Outcome: "restore-failed", ExitCode: 1, Restore: "fail"}, "verification tests is incomplete, failed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := recorded(t, 1)
			verify(t, f, importing(t, tc.data, tc.exit))
			got := onlyResult(t, f)
			tc.want.Mutation, tc.want.NativeRef = f.Plan.Chunks[0].Verification[0].Probe, got.NativeRef
			if got.Probe == nil || !reflect.DeepEqual(*got.Probe, tc.want) || got.NativeRef.Ref != "fixture:terminal-result" {
				t.Fatalf("probe = %+v, want %+v", got.Probe, tc.want)
			}
			if got.NativeRef.Excerpt != tc.data || got.NativeRef.Digest != rr.Digest([]byte(tc.data)) {
				t.Fatalf("excerpt = %q (%s), want the exact file bytes %q", got.NativeRef.Excerpt, got.NativeRef.Digest, tc.data)
			}
			verify(t, f, unprobed)
			if err := checkpoint(t, f); tc.refusal == "" && err != nil || tc.refusal != "" && (err == nil || !strings.Contains(err.Error(), tc.refusal)) {
				t.Fatalf("checkpoint = %v, want refusal %q", err, tc.refusal)
			}
		})
	}
}

// TestVerificationImportRefusals covers NR28 and the other import refusals: an answer that the probe owner refuses, an exit flag that
// disagrees with the command exit, or an answer with no focused exit refuses before any
// record byte changes. Naming an import beside the scalar flags is a usage error.
func TestVerificationImportRefusals(t *testing.T) {
	bit := answer("bit", "yes", 0, "1", report)
	f := recorded(t, 1)
	for name, tc := range map[string]struct {
		data, want string
		exit       int
	}{
		"duplicate envelope":  {strings.Replace(bit, ","+strconv.Itoa(len(report))+"\n", ","+strconv.Itoa(len(bit))+"\n", 1) + bit, "repeats", 0},
		"substituted schema":  {strings.Replace(bit, "bench-probe-result-v2", "bench-probe-result-v9", 1), "schema", 0},
		"older envelope":      {strings.Join(strings.SplitAfter(bit, "\n")[:4], "") + report, "probe_execution", 0},
		"string exit":         {strings.Replace(bit, "-v2,0,1,", "-v2,0,\"1\",", 1), "malformed", 0},
		"count mismatch":      {bit + "\n", "trailing_bytes", 0},
		"exit disagreement":   {bit, "--exit-code 1 differs", 1},
		"no mutated run":      {answer("invalid", "untouched", 1, "null", ""), "null", 1},
		"projected response":  {bit + "spilled{lines=20,bytes=900,omitted_lines=11,cut_lines=0,path=/outside}\n", "incomplete response", 0},
		"malformed tree lead": {"tree[1]{target,head}:\n  alpha,none\n" + bit, "identity row", 0},
		"reordered envelope": {strings.Replace(bit, "{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,",
			"{schema,focused_exit_code,command_exit_code,trailing_bytes}:\n  bench-probe-result-v2,1,0,", 1), "not canonical", 0},
	} {
		t.Run(name, func(t *testing.T) {
			refuseVerify(t, f, tc.want, importing(t, tc.data, tc.exit))
		})
	}
	combined := importing(t, bit, 0)
	combined["--probe-outcome"], combined["--probe-exit-code"], combined["--probe-restore"] = "bit", "1", "pass"
	if out, code := recordcmd.Command(f.Root, verifyArgs(t, combined)); code != 2 || !strings.HasPrefix(out, "usage: bench record verification") {
		t.Fatalf("scalar and import together = exit %d, output %q; want the usage line and exit 2", code, out)
	}
}

// The mutant-set answers below come from probe runs over the clamp repository. The silent
// row reads exit 1 and focused exit 0, and the biting row reads exit 0 and focused exit 1. The
// missed entry matches no clamp bytes.
var (
	silentEntry = []string{"--swap", "return n", "--with", "return -n"}
	bitingEntry = []string{"--swap", "n < 0", "--with", "n > 0"}
	missedEntry = []string{"--omit", "absent"}
	bite        = probetest.ClampFailure("caught")
)

// producer runs `bench probe --mutants` over entries, with their mutated runs in order, and
// answers its complete stdout and exit.
type producer func(entries [][]string, mutated ...probetest.Start) (string, int)

// setProducer makes the root of f the clamp repository and answers its producer and its Bench home.
func setProducer(t *testing.T, f *recordtest.Fixture) (producer, string) {
	t.Helper()
	_, home := probetest.Repository(t, f.Root, freshness.Publish)
	marker := filepath.Join(t.TempDir(), "go-started")
	return func(entries [][]string, mutated ...probetest.Start) (string, int) {
		t.Helper()
		body, err := json.Marshal(entries)
		if err != nil {
			t.Fatal(err)
		}
		probetest.InstallSequence(t, marker, probetest.ClampPass, mutated...)
		return probe.Command(slices.Concat([]string{probetest.ClampPath, "--mutants", excerptFile(t, "mutants.json", string(body))}, probetest.ClampSelection))
	}, home
}

// silentThenBiting is the two-row answer of the silent entry and then the biting entry.
func silentThenBiting(produce producer) (string, int) {
	return produce([][]string{silentEntry, bitingEntry}, probetest.ClampPass, bite)
}

// importingRow edits a verification call to import the row named by row from the answer in data.
func importingRow(t *testing.T, data string, exit int, row string) map[string]string {
	t.Helper()
	edits := importing(t, data, exit)
	edits["--mutant"] = row
	return edits
}

// --mutant 2 imports row 2 alone: its verdict, its focused exit, and its restore.
func TestVerificationImportsOneMutantRow(t *testing.T) {
	f := recorded(t, 1)
	produce, _ := setProducer(t, f)
	data, _ := silentThenBiting(produce)
	requireImportedRow(t, f, data, 0, "2", rr.Probe{Outcome: "bit", ExitCode: 1, Restore: "pass"})
}

// --mutant 1 of a row whose restore failed records that outcome and a failed restore. The
// mutated run removes the preserved copy, so the restore cannot put the start bytes back.
func TestVerificationImportsARestoreFailedRow(t *testing.T) {
	f := recorded(t, 1)
	produce, home := setProducer(t, f)
	damaging := bite
	damaging.Hook = "find " + sanitize.ShellQuote(home) + " -name clamp.go -delete\n"
	data, _ := produce([][]string{bitingEntry}, damaging)
	requireImportedRow(t, f, data, 2, "1", rr.Probe{Outcome: "restore-failed", ExitCode: 1, Restore: "fail"})
}

// requireImportedRow imports row of the answer in data with the explicit exit, and requires the
// one recorded probe want against the chunk's planned mutation, with the exact answer bytes.
func requireImportedRow(t *testing.T, f *recordtest.Fixture, data string, exit int, row string, want rr.Probe) {
	t.Helper()
	verify(t, f, importingRow(t, data, exit, row))
	got := onlyResult(t, f)
	want.Mutation, want.NativeRef = f.Plan.Chunks[0].Verification[0].Probe, got.NativeRef
	if got.Probe == nil || !reflect.DeepEqual(*got.Probe, want) || got.NativeRef.Excerpt != data {
		t.Fatalf("probe = %+v with excerpt %q, want %+v and the exact answer bytes", got.Probe, got.NativeRef.Excerpt, want)
	}
}

// The explicit exit is checked against the row's exit_code, never the set's command exit.
func TestVerificationChecksTheRowExit(t *testing.T) {
	f := recorded(t, 1)
	produce, _ := setProducer(t, f)
	data, command := silentThenBiting(produce)
	if command != 1 {
		t.Fatalf("command exit = %d, want 1, apart from row 2's exit", command)
	}
	refuseVerify(t, f, "— --exit-code 1 differs from the imported mutants.exit_code 0\n", importingRow(t, data, 1, "2"))
}

// A mutant-set answer without --mutant refuses by name and routes to --mutant.
func TestVerificationRefusesAMutantSetWithoutARow(t *testing.T) {
	f := recorded(t, 1)
	produce, _ := setProducer(t, f)
	data, command := silentThenBiting(produce)
	refuseVerify(t, f, "— --probe-result is a mutant-set answer; name one row with --mutant <n>\n", importing(t, data, command))
}

// --mutant refuses a single-form answer, a row past the answer, and a row whose mutated run never
// started. A bad row number, or --mutant without --probe-result, prints the form usage.
func TestVerificationMutantRefusals(t *testing.T) {
	f := recorded(t, 1)
	produce, _ := setProducer(t, f)
	set, _ := silentThenBiting(produce)
	missed, _ := produce([][]string{missedEntry, bitingEntry}, bite)
	if row := strings.Split(missed, "\n")[1]; !strings.Contains(row, ",substring-miss,") || !strings.HasSuffix(row, ",untouched") {
		t.Fatalf("missed row = %q, want the producer's substring-miss row", row)
	}
	for name, tc := range map[string]struct {
		edits map[string]string
		want  string
	}{
		"single-form answer":  {importingRow(t, answer("bit", "yes", 0, "1", report), 0, "1"), "--mutant needs a mutant-set answer; --probe-result holds a bench-probe-result-v2 answer"},
		"row past the answer": {importingRow(t, set, 1, "3"), "--mutant 3 names no row; the answer holds 2 rows"},
		"null focused exit":   {importingRow(t, missed, 1, "1"), "the imported mutants.focused_exit_code is null, because no mutated run started"},
	} {
		t.Run(name, func(t *testing.T) {
			refuseVerify(t, f, "— "+tc.want+"\n", tc.edits)
		})
	}
	for name, edits := range map[string]map[string]string{
		"zero row":           importingRow(t, set, 0, "0"),
		"negative row":       importingRow(t, set, 0, "-1"),
		"non-decimal row":    importingRow(t, set, 0, "x"),
		"row with no import": {"--mutant": "1"},
	} {
		if out, code := recordcmd.Command(f.Root, verifyArgs(t, edits)); code != 2 || !strings.HasPrefix(out, "usage: bench record verification") {
			t.Errorf("%s = exit %d, output %q; want the usage line and exit 2", name, code, out)
		}
	}
}

// TestVerificationHelpFieldSources covers NR32: the help names the producer verb and field of each copied probe value, so the
// command exit and the focused-execution exit cannot be confused.
func TestVerificationHelpFieldSources(t *testing.T) {
	var description string
	for _, row := range recordcmd.HelpRows() {
		if strings.HasPrefix(row.Suffix, " verification ") {
			description = row.Description
		}
	}
	for _, source := range []string{
		"--exit-code is the bench probe command exit, checked against probe_execution.command_exit_code on import",
		"--probe-outcome is bench probe probe.verdict",
		"--probe-exit-code is bench probe probe_execution.focused_exit_code, the normalized focused-execution exit",
		"--probe-restore is bench probe probe.restored through the producer owner's restoration projection",
	} {
		if !strings.Contains(description, source) {
			t.Errorf("verification help = %q, want %q", description, source)
		}
	}
}
