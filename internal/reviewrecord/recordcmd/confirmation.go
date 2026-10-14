package recordcmd

import (
	"strings"

	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/usage"
)

var confirmationForm = form{name: reviewrecord.FormConfirmation, description: "append one coordinator confirmation of a test-only repair to a recorded chunk; " +
	"the reviewed pair and the open findings derive from the current axis results, " + probeResult + " reads the probe as the verification form reads it, and --mutation records the caller's text unverified",
	flags: []flag{{"--chunk", "<id>", true, once}, {"--id", "<id>", true, once}, {"--performer", "<session>", true, once}, {"--model", "<model>", true, once},
		{"--effort", "<effort>", true, once}, {"--exit-code", "<n>", false, once}, {"--ref", "<ref>", true, once}, {probeResult, "<file>", false, once},
		{mutantFlag, "<n>", false, once}, {"--mutation", "<text>", true, once}, {"--finding", "<id>", true, repeated}},
	layout: []string{"--chunk", "--id", "--performer", "--model", "--effort", "--exit-code", "--ref", probeResult, "[", mutantFlag, "]", "--mutation", "--finding"},
	valid:  func(_ form, parsed usage.Result) bool { return exitsValid(parsed) }, run: confirmation}

func confirmation(f form, root, spec string, parsed usage.Result) (string, int) {
	values := parsed.Flags
	ref, out, code := native(values, probeResult)
	if out != "" {
		return out, code
	}
	exitCode, _ := integer(values, "--exit-code")
	mutant, _ := integer(values, mutantFlag)
	probe, out := importedProbe([]byte(ref.Excerpt), exitCode, mutant)
	if out != "" {
		return f.refuse(out)
	}
	probe.Mutation = values["--mutation"]
	call := reviewrecord.ConfirmationCall{Chunk: values["--chunk"], Findings: parsed.Repeated["--finding"], Probe: *probe,
		Evidence: reviewrecord.Evidence{ID: values["--id"], Performer: values["--performer"], Model: values["--model"], Effort: values["--effort"], NativeRef: ref}}
	entry, reviewed, err := reviewrecord.RecordConfirmation(root, spec, call)
	if err != nil {
		return f.refuse(cause(spec, err))
	}
	return row("confirmation", []string{"chunk", "id", "reviewed_tip", "reviewed_digest", "repaired_tip", "source_digest", "supersedes", "excerpt_digest"},
		[]string{call.Chunk, entry.ID, entry.ReviewedTip, reviewed.Digest, entry.RepairedTip, entry.SourceDigest, strings.Join(entry.Supersedes, ","), entry.NativeRef.Digest})
}
