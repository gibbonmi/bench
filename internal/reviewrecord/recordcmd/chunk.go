package recordcmd

import (
	"strconv"

	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/usage"
)

var chunkForm = form{name: reviewrecord.FormChunk, description: "write one chunk's frozen pair, digests, and acceptance rows into reviews/<slug>.md",
	flags: []flag{{"--chunk", "<id>", true, once}, {"--base", "<commit>", false, once}, {"--tip", "<commit>", false, once}}, run: chunk}

func chunk(f form, root, spec string, parsed usage.Result) (string, int) {
	values := parsed.Flags
	base, out, code := commit(f, root, values, "--base")
	if out != "" {
		return out, code
	}
	tip, out, code := commit(f, root, values, "--tip")
	if out != "" {
		return out, code
	}
	if out, code := refusePlanned(f, root, spec, tip, chunkRows(values["--chunk"])); out != "" {
		return out, code
	}
	action, entry, err := reviewrecord.RecordChunk(root, spec, values["--chunk"], base, tip)
	if err != nil {
		return f.refuse(cause(spec, err))
	}
	return row("chunk", []string{"id", "action", "base", "tip", "source_digest", "plan_digest", "rows"},
		[]string{entry.ID, action, entry.Base, entry.Tip, entry.SourceDigest, entry.PlanDigest, strconv.Itoa(len(entry.AcceptanceRows))})
}
