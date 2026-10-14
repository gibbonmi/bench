package recordcmd

import (
	"errors"
	"fmt"

	"github.com/gibbonmi/bench/internal/probe"
	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/usage"
)

// probeFlags are the three flags of a planned probe.
var probeFlags = []string{"--probe-outcome", "--probe-exit-code", "--probe-restore"}

// probeResult names a complete `bench probe` answer that supplies the excerpt and the three
// probe values in place of --excerpt and probeFlags.
const probeResult = "--probe-result"

// mutantFlag names the row of a mutant-set answer that probeResult imports.
const mutantFlag = "--mutant"

var verificationForm = form{name: reviewrecord.FormVerification, description: "append one planned verification result with its computed digests; " + probeSources(),
	flags: []flag{{"--chunk", "<id>", true, once}, {"--source", "<commit>", false, once}, {"--final", "", false, once},
		{"--requirement", "<id>", true, once}, {"--id", "<id>", true, once}, {"--performer", "<session>", true, once},
		{"--model", "<model>", true, once}, {"--effort", "<effort>", true, once}, {"--exit-code", "<n>", false, once},
		{"--ref", "<ref>", true, once}, {"--excerpt", "<file>", true, once}, {probeFlags[0], "<verdict>", true, once},
		{probeFlags[1], "<n>", false, once}, {probeFlags[2], "pass|fail", false, once}, {probeResult, "<file>", false, once},
		{mutantFlag, "<n>", false, once}},
	layout: []string{"(", "--chunk", "[", "--source", "]", "|", "--final", "--source", ")", "--requirement", "--id", "--performer",
		"--model", "--effort", "--exit-code", "--ref", "(", "--excerpt", "[", probeFlags[0], probeFlags[1], probeFlags[2], "]", "|", probeResult, "[", mutantFlag, "]", ")"},
	valid: verificationValid, run: verification}

// probeSources is the verification form's note on where each copied probe value comes
// from. Each field name comes from the probe owner, so the help names the field it renders.
func probeSources() string {
	return "for a probe, --exit-code is the bench probe command exit, checked against " + probe.FieldSource("command_exit_code") + " on import, " +
		probeFlags[0] + " is bench probe " + probe.FieldSource("verdict") + ", " +
		probeFlags[1] + " is bench probe " + probe.FieldSource("focused_exit_code") + ", the normalized focused-execution exit, " +
		probeFlags[2] + " is bench probe " + probe.FieldSource("restored") + " through the producer owner's restoration projection, and " +
		probeResult + " reads them from one complete bench probe answer; with " + mutantFlag + ", " + probeResult +
		" reads row n of one complete mutant-set answer, --exit-code is checked against " + probe.SetFieldSource("exit_code") +
		", and the probe values are " + probe.SetFieldSource("verdict") + ", " + probe.SetFieldSource("focused_exit_code") + ", and " + probe.SetFieldSource("restored")
}

// verificationValid closes the value sets of the verification form. Command calls it
// only for a set of flags that the layout admits.
func verificationValid(f form, parsed usage.Result) bool {
	_, probe := parsed.Flags[probeFlags[0]]
	_, probeExitCode := integer(parsed.Flags, probeFlags[1])
	return exitsValid(parsed) && (!probe || probeExitCode && chosen(f, parsed.Flags, probeFlags[2]))
}

// exitsValid closes the values that a probe import reads: an integer --exit-code, and a
// positive decimal row when the call names mutantFlag.
func exitsValid(parsed usage.Result) bool {
	_, exitCode := integer(parsed.Flags, "--exit-code")
	return exitCode && mutantValid(parsed)
}

// mutantValid closes the mutantFlag value: a positive decimal row when the call names it.
func mutantValid(parsed usage.Result) bool {
	_, mutant := parsed.Flags[mutantFlag]
	row, decimal := integer(parsed.Flags, mutantFlag)
	return !mutant || decimal && row > 0
}

func verification(f form, root, spec string, parsed usage.Result) (string, int) {
	values := parsed.Flags
	_, imported := values[probeResult]
	excerpt := "--excerpt"
	if imported {
		excerpt = probeResult
	}
	ref, out, code := native(values, excerpt)
	if out != "" {
		return out, code
	}
	call := reviewrecord.VerificationCall{Chunk: values["--chunk"], Requirement: values["--requirement"],
		Evidence: reviewrecord.Evidence{ID: values["--id"], Performer: values["--performer"], Model: values["--model"], Effort: values["--effort"], NativeRef: ref}}
	if _, ok := values["--source"]; ok {
		if call.Source, out, code = commit(f, root, values, "--source"); out != "" {
			return out, code
		}
	}
	call.ExitCode, _ = integer(values, "--exit-code")
	if imported {
		mutant, _ := integer(values, mutantFlag)
		if call.Probe, out = importedProbe([]byte(ref.Excerpt), call.ExitCode, mutant); out != "" {
			return f.refuse(out)
		}
	} else if _, ok := values[probeFlags[0]]; ok {
		exitCode, _ := integer(values, probeFlags[1])
		call.Probe = &reviewrecord.Probe{Outcome: values[probeFlags[0]], ExitCode: exitCode, Restore: values[probeFlags[2]]}
	}
	entry, err := reviewrecord.RecordVerification(root, spec, call)
	if err != nil {
		return f.refuse(cause(spec, err))
	}
	list := "chunk"
	if call.Chunk == "" {
		list = "completion"
	}
	return row("verification", []string{"list", "chunk", "id", "requirement", "role", "outcome", "source_digest", "excerpt_digest"},
		[]string{list, call.Chunk, entry.ID, entry.Requirement, entry.Role, entry.Outcome, entry.SourceDigest, entry.NativeRef.Digest})
}

// importedAnswer is the probe values that one answer form binds, with the producer field of
// each exit.
type importedAnswer struct {
	verdict                 string
	restored                bool
	exit                    int
	focused                 *int
	exitField, focusedField string
}

// readImported reads the probe values through the probe owner: a single-form answer when
// mutant is 0, and row mutant of a mutant-set answer otherwise. The explicit exit binds to the
// command exit of a single-form answer and to the row's own exit of a mutant set. source
// names the answer in a refusal.
func readImported(data []byte, mutant int, source string) (importedAnswer, string) {
	if mutant == 0 {
		result, err := probe.ReadResult(data)
		if err != nil {
			return importedAnswer{}, importRefusal(err, mutant, source)
		}
		return importedAnswer{verdict: result.Verdict, restored: result.Restored, exit: result.CommandExit, focused: result.FocusedExit,
			exitField: probe.FieldSource("command_exit_code"), focusedField: probe.FieldSource("focused_exit_code")}, ""
	}
	result, err := probe.ReadMutant(data, mutant)
	if err != nil {
		return importedAnswer{}, importRefusal(err, mutant, source)
	}
	return importedAnswer{verdict: result.Verdict, restored: result.Restored, exit: result.RowExit, focused: result.FocusedExit,
		exitField: probe.SetFieldSource("exit_code"), focusedField: probe.SetFieldSource("focused_exit_code")}, ""
}

// importRefusal maps a probe reader's refusal of the answer that source names to its record
// line. The probe owner's sentinels carry no flag, so each line here names the flag that
// repairs the call.
func importRefusal(err error, mutant int, source string) string {
	var missing *probe.RowRangeError
	switch {
	case errors.Is(err, probe.ErrMutantSet):
		return source + " is a mutant-set answer; name one row with " + mutantFlag + " <n>"
	case errors.Is(err, probe.ErrSingleForm):
		return mutantFlag + " needs a mutant-set answer; " + source + " holds a " + probe.ResultSchema + " answer"
	case errors.As(err, &missing):
		return fmt.Sprintf("%s %d names no row; the answer holds %d rows", mutantFlag, mutant, missing.Rows)
	}
	return source + " " + err.Error()
}

// importedProbe reads the probe values from one complete probe answer, or from its row mutant
// when mutant is not 0. The explicit exit must equal the exit that the answer binds, and an answer
// whose mutated run never started has no focused exit to record. A silent or unrestored answer
// imports as it reads; the checkpoint or the confirmation writer refuses it.
func importedProbe(data []byte, exitCode, mutant int) (*reviewrecord.Probe, string) {
	answer, out := readImported(data, mutant, probeResult)
	if out != "" {
		return nil, out
	}
	if answer.exit != exitCode {
		return nil, fmt.Sprintf("--exit-code %d differs from the imported %s %d", exitCode, answer.exitField, answer.exit)
	}
	if answer.focused == nil {
		return nil, fmt.Sprintf("the imported %s is null, because no mutated run started", answer.focusedField)
	}
	restore := "fail"
	if answer.restored {
		restore = "pass"
	}
	return &reviewrecord.Probe{Outcome: answer.verdict, ExitCode: *answer.focused, Restore: restore}, ""
}
