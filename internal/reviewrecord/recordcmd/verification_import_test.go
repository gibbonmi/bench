package recordcmd_test

import (
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordcmd"
	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
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
