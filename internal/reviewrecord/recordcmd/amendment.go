package recordcmd

import (
	"slices"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/usage"
)

var amendmentForm = form{name: "amendment", description: "record the plan-digest change at a source commit",
	flags: []flag{{"--source", "<commit>", false, once}, {"--map", "<old>=<new>[,<new>...]", true, repeated}}, valid: amendmentValid, run: amendment}

// changes reads the --map values of one amendment call into each old chunk ID and its
// new chunk IDs. ok is false for an empty chunk ID on either side and for an old chunk ID
// that an earlier value names.
func changes(parsed usage.Result) (map[string][]string, bool) {
	result := map[string][]string{}
	for _, value := range parsed.Repeated["--map"] {
		// A value with no "=" has an empty new side, so the empty chunk ID rule refuses it.
		old, targets, _ := strings.Cut(value, "=")
		ids := strings.Split(targets, ",")
		if _, named := result[old]; old == "" || named || slices.Contains(ids, "") {
			return nil, false
		}
		result[old] = ids
	}
	return result, true
}

// amendmentValid closes the --map value grammar.
func amendmentValid(_ form, parsed usage.Result) bool {
	_, ok := changes(parsed)
	return ok
}

func amendment(f form, root, spec string, parsed usage.Result) (string, int) {
	source, out, code := commit(f, root, parsed.Flags, "--source")
	if out != "" {
		return out, code
	}
	mapped, _ := changes(parsed)
	entry, err := reviewrecord.RecordAmendment(root, spec, source, mapped)
	if err != nil {
		return f.refuse(cause(err))
	}
	return row("amendment", []string{"from", "to", "chunks"}, []string{entry.From, entry.To, strconv.Itoa(len(entry.ChunkIDs))})
}
