package recordcmd

import (
	"errors"
	"strings"

	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/usage"
)

// probeFlags are the three flags of a planned probe.
var probeFlags = []string{"--probe-outcome", "--probe-exit-code", "--probe-restore"}

var verificationForm = form{name: "verification", description: "append one planned verification result with its computed digests",
	flags: []flag{{"--chunk", "<id>", true, once}, {"--source", "<commit>", false, once}, {"--final", "", false, once},
		{"--requirement", "<id>", true, once}, {"--id", "<id>", true, once}, {"--performer", "<session>", true, once},
		{"--model", "<model>", true, once}, {"--effort", "<effort>", true, once}, {"--exit-code", "<n>", false, once},
		{"--ref", "<ref>", true, once}, {"--excerpt", "<file>", true, once}, {probeFlags[0], "<verdict>", true, once},
		{probeFlags[1], "<n>", false, once}, {probeFlags[2], "pass|fail", false, once}},
	layout: []string{"(", "--chunk", "[", "--source", "]", "|", "--final", "--source", ")", "--requirement", "--id", "--performer",
		"--model", "--effort", "--exit-code", "--ref", "--excerpt", "[", probeFlags[0], probeFlags[1], probeFlags[2], "]"},
	valid: verificationValid, run: verification}

// verificationValid closes the value sets of the verification form. Command calls it
// only for a set of flags that the layout admits.
func verificationValid(f form, parsed usage.Result) bool {
	_, exitCode := integer(parsed.Flags, "--exit-code")
	_, probe := parsed.Flags[probeFlags[0]]
	_, probeExitCode := integer(parsed.Flags, probeFlags[1])
	return exitCode && (!probe || probeExitCode && chosen(f, parsed.Flags, probeFlags[2]))
}

// cause is the refusal cause of err. A missing chunk entry names the chunk form, and a
// probe mismatch names the probe flags.
func cause(err error) string {
	switch {
	case errors.Is(err, reviewrecord.ErrNoChunkEntry):
		return err.Error() + "; write it with " + command + " chunk"
	case errors.Is(err, reviewrecord.ErrProbe):
		return err.Error() + "; the probe flags are " + strings.Join(probeFlags, ", ")
	}
	return err.Error()
}

func verification(f form, root, spec string, parsed usage.Result) (string, int) {
	values := parsed.Flags
	ref, out, code := native(values)
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
	if _, ok := values[probeFlags[0]]; ok {
		exitCode, _ := integer(values, probeFlags[1])
		call.Probe = &reviewrecord.Probe{Outcome: values[probeFlags[0]], ExitCode: exitCode, Restore: values[probeFlags[2]]}
	}
	entry, err := reviewrecord.RecordVerification(root, spec, call)
	if err != nil {
		return f.refuse(cause(err))
	}
	list := "chunk"
	if call.Chunk == "" {
		list = "completion"
	}
	return row("verification", []string{"list", "chunk", "id", "requirement", "role", "outcome", "source_digest", "excerpt_digest"},
		[]string{list, call.Chunk, entry.ID, entry.Requirement, entry.Role, entry.Outcome, entry.SourceDigest, entry.NativeRef.Digest})
}
