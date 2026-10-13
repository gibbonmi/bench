package probe

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/gibbonmi/bench/internal/jsonfile"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// maxMutants caps the entries of one mutant set, so one call stays bounded.
const maxMutants = 64

// entryGrammar parses one mutant-set entry: the shared mutation flags and --expect, with no
// positional. So an entry means what the same mutation argv means in the single form.
var entryGrammar = usage.Grammar{
	Cmd:   grammar.Cmd,
	Help:  grammar.Help,
	Flags: slices.Concat(mutationFlags, []usage.Flag{expectFlag}),
}

// setEntry is one resolved entry: its mutation, its mutant over the start bytes, and its
// --expect witness, which is empty when it names none.
type setEntry struct {
	input  mutation
	mutant mutationResult
	expect string
}

// readMutantSet reads the --mutants file of parsed and resolves each entry over the subject's
// start bytes. Every refusal answers before the selection is prepared and before any run.
func readMutantSet(parsed usage.Result, subject subject) ([]setEntry, string, int) {
	path := parsed.Flags[mutantsFlag.Name]
	body, refused := readMutationFile(path)
	if refused != nil {
		return nil, refused.line(), 1
	}
	var list [][]string
	if err := jsonfile.DecodeDocument([]byte(body), &list); err != nil {
		return nil, setMalformed(path, err.Error()), 1
	}
	// A JSON null decodes to a nil list or a nil entry, and [] decodes to an empty one.
	if list == nil {
		return nil, setMalformed(path, "the list is null"), 1
	}
	if n := slices.IndexFunc(list, func(entry []string) bool { return entry == nil }); n >= 0 {
		return nil, setMalformed(path, "entry "+strconv.Itoa(n+1)+" is null"), 1
	}
	if len(list) == 0 {
		return nil, toon.Errorf("probe mutant set empty", path+" holds no entry") + "\n", 1
	}
	if len(list) > maxMutants {
		return nil, toon.Errorf("probe mutant set too large", fmt.Sprintf("%s holds %d entries, the cap is %d", path, len(list), maxMutants)) + "\n", 1
	}
	entries := make([]setEntry, len(list))
	for i, argv := range list {
		entry, line, code := resolveEntry(i+1, argv, parsed, subject)
		if line != "" {
			return nil, line, code
		}
		entries[i] = entry
	}
	return entries, "", 0
}

// resolveEntry grades entry n through the mutation half of the form matrix and the cross-half
// rule against the top-level selection. It then resolves the entry's mutant over the start
// bytes.
func resolveEntry(n int, argv []string, top usage.Result, subject subject) (setEntry, string, int) {
	parsed, line, _ := usage.Parse(entryGrammar, argv)
	if line != "" || !mutationFormFits(parsed) || !expectFitsSelection(parsed, top) {
		line, code := usageLine()
		return setEntry{}, toon.Errorf("probe mutant set entry refused", fmt.Sprintf("entry %d does not match the mutation grammar", n)) + "\n" + line, code
	}
	input, refused := mutationInput(parsed)
	var mutant mutationResult
	if refused == nil {
		mutant, refused = mutate(subject.start, input, subject.display)
	}
	if refused != nil {
		refused.detail = fmt.Sprintf("entry %d: %s", n, refused.detail)
		return setEntry{}, refused.line(), 1
	}
	return setEntry{input: input, mutant: mutant, expect: parsed.Flags[expectFlag.Name]}, "", 0
}

func setMalformed(path, reason string) string {
	return toon.Errorf("probe mutant set malformed", path+": "+reason) + "\n"
}
