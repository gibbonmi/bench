package recordcmd

import (
	"strconv"

	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/usage"
)

var completionForm = form{name: "completion", description: "validate and write final completion evidence for a source commit",
	flags: []flag{{"--source", "<commit>", false, once}}, run: completion}

func completion(f form, root, spec string, parsed usage.Result) (string, int) {
	source, out, code := commit(f, root, parsed.Flags, "--source")
	if out != "" {
		return out, code
	}
	entry, err := reviewrecord.RecordCompletion(root, spec, source)
	if err != nil {
		return f.refuse(cause(err))
	}
	return row("completion", []string{"state", "source_digest", "performer", "rows", "verification"},
		[]string{entry.State, entry.SourceDigest, entry.Performer, strconv.Itoa(len(entry.Reconciliation)), strconv.Itoa(len(entry.Verification))})
}
