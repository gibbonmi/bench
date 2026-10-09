package recordcmd

import (
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

var verificationForm = form{name: reviewrecord.FormVerification, description: "append one planned verification result with its computed digests; " + probeSources(),
	flags: []flag{{"--chunk", "<id>", true, once}, {"--source", "<commit>", false, once}, {"--final", "", false, once},
		{"--requirement", "<id>", true, once}, {"--id", "<id>", true, once}, {"--performer", "<session>", true, once},
		{"--model", "<model>", true, once}, {"--effort", "<effort>", true, once}, {"--exit-code", "<n>", false, once},
		{"--ref", "<ref>", true, once}, {"--excerpt", "<file>", true, once}, {probeFlags[0], "<verdict>", true, once},
		{probeFlags[1], "<n>", false, once}, {probeFlags[2], "pass|fail", false, once}, {probeResult, "<file>", false, once}},
	layout: []string{"(", "--chunk", "[", "--source", "]", "|", "--final", "--source", ")", "--requirement", "--id", "--performer",
		"--model", "--effort", "--exit-code", "--ref", "(", "--excerpt", "[", probeFlags[0], probeFlags[1], probeFlags[2], "]", "|", probeResult, ")"},
	valid: verificationValid, run: verification}

// probeSources is the verification form's note on where each copied probe value comes
// from. Each field name comes from the probe owner, so the help names the field it renders.
func probeSources() string {
	return "for a probe, --exit-code is the bench probe command exit, checked against " + probe.FieldSource("command_exit_code") + " on import, " +
		probeFlags[0] + " is bench probe " + probe.FieldSource("verdict") + ", " +
		probeFlags[1] + " is bench probe " + probe.FieldSource("focused_exit_code") + ", the normalized focused-execution exit, " +
		probeFlags[2] + " is bench probe " + probe.FieldSource("restored") + " through the probe owner's restoration projection, and " +
		probeResult + " reads them from one complete bench probe answer"
}

// verificationValid closes the value sets of the verification form. Command calls it
// only for a set of flags that the layout admits.
func verificationValid(f form, parsed usage.Result) bool {
	_, exitCode := integer(parsed.Flags, "--exit-code")
	_, probe := parsed.Flags[probeFlags[0]]
	_, probeExitCode := integer(parsed.Flags, probeFlags[1])
	return exitCode && (!probe || probeExitCode && chosen(f, parsed.Flags, probeFlags[2]))
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
		if call.Probe, out = importedProbe([]byte(ref.Excerpt), call.ExitCode); out != "" {
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

// importedProbe reads the probe values from one complete probe answer through the probe
// owner. The explicit exit must equal the answer's command exit, and an answer whose mutated
// run never started has no focused exit to record. A silent or unrestored answer records
// as it reads, so the checkpoint still refuses it.
func importedProbe(data []byte, exitCode int) (*reviewrecord.Probe, string) {
	result, err := probe.ReadResult(data)
	if err != nil {
		return nil, probeResult + " " + err.Error()
	}
	if result.CommandExit != exitCode {
		return nil, fmt.Sprintf("--exit-code %d differs from the imported %s %d", exitCode, probe.FieldSource("command_exit_code"), result.CommandExit)
	}
	if result.FocusedExit == nil {
		return nil, fmt.Sprintf("the imported %s is null, because no mutated run started", probe.FieldSource("focused_exit_code"))
	}
	restore := "fail"
	if result.Restored {
		restore = "pass"
	}
	return &reviewrecord.Probe{Outcome: result.Verdict, ExitCode: *result.FocusedExit, Restore: restore}, ""
}
