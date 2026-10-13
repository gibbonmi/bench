package probe

import (
	"errors"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/probe/probetest"
	"github.com/gibbonmi/bench/internal/responsebound"
)

// The exits below are authored apart from the producer: a biting probe exits 0 while its
// mutated run exits 1, and a probe whose mutated run never started reports no focused exit.
// A producer that copies its command exit into the focused cell reds the biting row.

// readAnswer reads one probe answer and requires a complete result.
func readAnswer(t *testing.T, out string) Result {
	t.Helper()
	result, err := ReadResult([]byte(out))
	if err != nil {
		t.Fatalf("ReadResult(%q) = %v, want a complete result", out, err)
	}
	return result
}

// TestProbeReportsNormalizedExecutionExit covers NR23 and NR24: the focused exit is the mutated run's normalized exit, separate from the
// command exit, and it is null on every path where the mutated run never started.
func TestProbeReportsNormalizedExecutionExit(t *testing.T) {
	one := 1
	for _, tc := range []struct {
		name              string
		baseline, mutated probetest.Start
		args              []string
		command           int
		focused           *int
		verdict           string
		restored          bool
	}{
		{"bit", probetest.Start{Events: cannedPass}, probetest.Start{Events: cannedFailure("caught"), Exit: 1}, probeArgs(), 0, &one, "bit", true},
		{"red baseline", probetest.Start{Events: cannedFailure("red"), Exit: 1}, probetest.Start{Events: cannedPass}, probeArgs(), 1, nil, "invalid", false},
		{"substring miss", probetest.Start{Events: cannedPass}, probetest.Start{Events: cannedPass}, probeArgs("clamp.go", "--omit", "absent", "--package", "./"), 1, nil, "invalid", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t)
			installStubStarts(t, f, tc.baseline, tc.mutated)
			out, code := runProbe(t, tc.args...)
			got := readAnswer(t, out)
			if code != tc.command || got.CommandExit != tc.command || got.Verdict != tc.verdict || got.Restored != tc.restored {
				t.Fatalf("probe = (%d, %+v), want command exit %d, verdict %s, restored %v\n%s", code, got, tc.command, tc.verdict, tc.restored, out)
			}
			if (got.FocusedExit == nil) != (tc.focused == nil) || got.FocusedExit != nil && *got.FocusedExit != *tc.focused {
				t.Fatalf("focused exit = %v, want %v\n%s", got.FocusedExit, tc.focused, out)
			}
			if strings.Count(out, "probe_execution[") != 1 || !strings.Contains(out, "selection[1]{form,target,run,baseline,ran}:\n") {
				t.Fatalf("stdout = %q, want the selection row and one execution row", out)
			}
		})
	}
}

// A failed restore reports the mutated run's exit 1 beside command exit 2, and the
// trailing count covers the preserved row and the report that follow the envelope.
func TestProbeCountsTheRestoreFailureSuffix(t *testing.T) {
	f := newFixture(t)
	installStubGo(t, f, cannedFailure("caught"), 1, "chmod 0500 .\n")
	t.Cleanup(func() { _ = os.Chmod(f.root, 0o755) })
	out, code := runProbe(t, probeArgs()...)
	got := readAnswer(t, out)
	if code != 2 || got.CommandExit != 2 || got.FocusedExit == nil || *got.FocusedExit != 1 || got.Verdict != "restore-failed" || got.Restored {
		t.Fatalf("restore failure = (%d, %+v), want command 2, focused 1, unrestored\n%s", code, got, out)
	}
}

// bitingAnswer is the complete stdout of one biting probe in a new fixture.
func bitingAnswer(t *testing.T) string {
	t.Helper()
	f := newFixture(t)
	installStubGo(t, f, cannedFailure("caught"), 1, "")
	out, code := runProbe(t, probeArgs()...)
	if code != 0 {
		t.Fatalf("probe = (%q, %d), want a bite", out, code)
	}
	return out
}

// TestProbeResultSchemaVersion covers PS54 and PS55: the schema token names the verdict-row
// field set, so ReadResult reads a v2 answer and refuses a v1 answer by its schema rather
// than by a field error. The tokens and the v1 field set are authored apart from the producer.
func TestProbeResultSchemaVersion(t *testing.T) {
	out := bitingAnswer(t)
	if !strings.Contains(out, "probe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,") {
		t.Fatalf("stdout = %q, want the schema cell bench-probe-result-v2", out)
	}
	readAnswer(t, out)
	v1 := strings.Replace(strings.Replace(out, ",first_failure", "", 1), ",TestClampNegative,yes\n", ",yes\n", 1)
	v1 = strings.Replace(v1, "bench-probe-result-v2,", "bench-probe-result-v1,", 1)
	want := `the schema "bench-probe-result-v1" is unsupported`
	if _, err := ReadResult([]byte(v1)); err == nil || !strings.Contains(err.Error(), want) {
		t.Fatalf("ReadResult(v1) = %v, want %q", err, want)
	}
}

// executionLine answers the index of the execution row in lines, in either answer form.
func executionLine(t *testing.T, lines []string) int {
	t.Helper()
	header := slices.IndexFunc(lines, func(line string) bool {
		return strings.HasPrefix(line, resultTables[2].name+"[") || strings.HasPrefix(line, setTables[2].name+"[")
	})
	if header < 0 {
		t.Fatalf("answer %q has no execution row", lines)
	}
	return header + 1
}

// declared rewrites the trailing_bytes cell of out's execution row to n.
func declared(t *testing.T, out string, n int) string {
	t.Helper()
	lines := strings.SplitAfter(out, "\n")
	row := executionLine(t, lines)
	cells := strings.Split(strings.TrimSuffix(lines[row], "\n"), ",")
	cells[len(cells)-1] = strconv.Itoa(n)
	lines[row] = strings.Join(cells, ",") + "\n"
	return strings.Join(lines, "")
}

// TestProbeResultEnvelopeRefusals covers NR27, NR44, and NR45: every envelope that is incomplete, noncanonical, repeated, mistyped,
// of another schema, or miscounted refuses. The tree rows are authored apart from Row.
func TestProbeResultEnvelopeRefusals(t *testing.T) {
	out := bitingAnswer(t)
	lines := strings.SplitAfter(out, "\n")
	envelope, suffix := strings.Join(lines[:6], ""), strings.Join(lines[6:], "")
	trailing := len(suffix)
	lead := "tree[1]{target,head,dirty}:\n  primary,none,false\n"
	if got := readAnswer(t, lead+out); got.CommandExit != 0 || *got.FocusedExit != 1 {
		t.Fatalf("tree-prefixed result = %+v, want command 0 and focused 1", got)
	}
	for name, data := range map[string]string{
		"older envelope":         strings.Join(lines[:4], "") + suffix,
		"substituted schema":     strings.Replace(out, resultSchema+",", "bench-probe-result-v0,", 1),
		"absent focused field":   strings.Replace(strings.Replace(out, ",focused_exit_code", "", 1), ",0,1,", ",0,", 1),
		"string focused exit":    strings.Replace(out, ",0,1,"+strconv.Itoa(trailing)+"\n", ",0,\"1\","+strconv.Itoa(trailing)+"\n", 1),
		"fractional exit":        strings.Replace(out, ",0,1,", ",0,1.5,", 1),
		"repeated envelope":      declared(t, envelope, len(out)) + out,
		"missing envelope":       suffix,
		"truncated envelope":     out[:len(lines[0])+len(lines[1])+4],
		"truncated suffix":       out[:len(out)-1],
		"extra suffix byte":      out + "x",
		"count one over":         declared(t, out, trailing+1),
		"count one under":        declared(t, out, trailing-1),
		"negative count":         declared(t, envelope, -1),
		"malformed tree row":     "tree[1]{target,head}:\n  primary,none\n" + out,
		"repeated tree row":      lead + lead + out,
		"unterminated tree row":  strings.TrimSuffix(lead, "\n"),
		"spill frame in suffix":  declared(t, out, trailing+len("spill-failed{reason=r}\n")) + "spill-failed{reason=r}\n",
		"noncanonical row space": strings.Replace(out, "  bit,", "   bit,", 1),
		// The two rows below decode cleanly under the strict decoder, so only the re-render
		// check refuses them.
		"reordered execution fields": strings.Replace(strings.Replace(out, "{schema,command_exit_code,focused_exit_code,", "{schema,focused_exit_code,command_exit_code,", 1),
			resultSchema+",0,1,", resultSchema+",1,0,", 1),
		"quoted bare cell": strings.Replace(out, "  bit,", "  \"bit\",", 1),
	} {
		if result, err := ReadResult([]byte(data)); err == nil {
			t.Errorf("%s: ReadResult = %+v, want a refusal", name, result)
		}
	}
	if _, err := ReadResult([]byte(out + "spilled{lines=1,bytes=1,omitted_lines=0,cut_lines=0,path=/p}\n")); !errors.Is(err, responsebound.ErrIncomplete) {
		t.Fatalf("projected result = %v, want responsebound.ErrIncomplete", err)
	}
}

// setAnswer is the complete stdout of a set over clamp.go whose entries bite, miss, and bite,
// with its expected rows.
func setAnswer(t *testing.T) (string, [][]any) {
	t.Helper()
	f := newFixture(t)
	installSet(t, f, bite, bite)
	out, code := runProbe(t, setArgs(mutantsFile(t, guardSwap, missEntry, returnSwap))...)
	rows := [][]any{biteRow(1, "swap"), missRow(2), biteRow(3, "swap")}
	requireRow(t, out, code, 1, setRows(t, rows...))
	return out, rows
}

// mutantResult derives what ReadMutant answers for one expected mutants row of an answer whose
// command exit is command.
func mutantResult(row []any, command int) MutantResult {
	cell := func(field string) any { return row[slices.Index(mutantFieldsOfSet, field)] }
	result := MutantResult{Verdict: cell("verdict").(string), Restored: cell("restored") == restoredCell, RowExit: cell("exit_code").(int), CommandExit: command}
	if focused, ok := cell("focused_exit_code").(int); ok {
		result.FocusedExit = &focused
	}
	return result
}

// A single-form reader refuses a mutant-set answer with the sentinel that names the set.
func TestReadResultRefusesAMutantSet(t *testing.T) {
	out, _ := setAnswer(t)
	if result, err := ReadResult([]byte(out)); !errors.Is(err, ErrMutantSet) {
		t.Fatalf("ReadResult(set) = (%+v, %v), want ErrMutantSet", result, err)
	}
}

// ReadMutant answers each row of a complete answer whose trailing_bytes counts every report. An
// answer cut or extended by one byte refuses, and so does a row number past the rows.
func TestReadMutantRequiresTheCompleteAnswer(t *testing.T) {
	out, rows := setAnswer(t)
	requireTrailingBytes(t, out, 1)
	for i, row := range rows {
		if got, err := ReadMutant([]byte(out), i+1); err != nil || !reflect.DeepEqual(got, mutantResult(row, 1)) {
			t.Errorf("ReadMutant(row %d) = (%+v, %v), want %+v", i+1, got, err, mutantResult(row, 1))
		}
	}
	for name, data := range map[string]string{"cut": out[:len(out)-1], "extended": out + "x"} {
		if result, err := ReadMutant([]byte(data), 1); err == nil {
			t.Errorf("%s answer: ReadMutant = %+v, want a refusal", name, result)
		}
	}
	for _, n := range []int{0, len(rows) + 1} {
		var missing *RowRangeError
		if _, err := ReadMutant([]byte(out), n); !errors.As(err, &missing) || missing.Rows != len(rows) {
			t.Errorf("ReadMutant(row %d) = %v, want a RowRangeError of %d rows", n, err, len(rows))
		}
	}
}

// ReadMutant keeps the single form's rules: a spill frame, a non-canonical row, a repeated
// mutants table, and another schema refuse before any row returns, and a single-form answer
// refuses with the sentinel that names it.
func TestReadMutantSharesTheReaderRules(t *testing.T) {
	out, _ := setAnswer(t)
	lines := strings.SplitAfter(out, "\n")
	trailing := len(strings.Join(lines[executionLine(t, lines)+1:], ""))
	frame := "spill-failed{reason=r}\n"
	for name, data := range map[string]string{
		"spill frame":            declared(t, out, trailing+len(frame)) + frame,
		"quoted row cell":        strings.Replace(out, ",bit,", ",\"bit\",", 1),
		"repeated mutants table": declared(t, out, trailing+len(out)) + out,
		"substituted schema":     strings.Replace(out, setSchema+",", "bench-probe-set-v0,", 1),
	} {
		if result, err := ReadMutant([]byte(data), 1); err == nil {
			t.Errorf("%s: ReadMutant = %+v, want a refusal", name, result)
		}
	}
	if _, err := ReadMutant([]byte(out+"spilled{lines=1,bytes=1,omitted_lines=0,cut_lines=0,path=/p}\n"), 1); !errors.Is(err, responsebound.ErrIncomplete) {
		t.Fatalf("projected set = %v, want responsebound.ErrIncomplete", err)
	}
	if _, err := ReadMutant([]byte(bitingAnswer(t)), 1); !errors.Is(err, ErrSingleForm) {
		t.Fatalf("ReadMutant(single form) = %v, want ErrSingleForm", err)
	}
}
