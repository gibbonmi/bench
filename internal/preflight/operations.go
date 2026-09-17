package preflight

import (
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// Operation kinds. Each registered operation dispatches to exactly one kind.
const (
	opVerdict = iota
	opLegacyCharge
	opProposal
	opPrepareEvidence
	opReadEvidence
)

// Flag spellings the operations share.
const (
	flagBase     = "--base"
	flagTip      = "--source-tip"
	flagCharge   = "--charge"
	flagPropose  = "--propose-writes"
	flagTicket   = "--ticket"
	flagFull     = "--full"
	flagQuota    = "--max-store-bytes"
	flagCursor   = "--cursor"
	modeReview   = "review"
	modeEvidence = "evidence"
)

// operation is one implemented public preflight form. The registry below is the one source
// for the argument grammar, the preflight help, and the root help rows.
type operation struct {
	mode string
	// selectors are the switch flags that choose this form; required and optional list
	// every other flag the form accepts.
	selectors, required, optional []string
	usage, description            string
	kind                          int
	// bounded forms obey the shared response bound for every response they produce.
	bounded bool
}

var operations = []operation{
	{mode: modeReview, optional: []string{flagBase, flagTip}, kind: opVerdict,
		usage: "review <slug> [--base <commit>] [--source-tip <commit>]", description: "review-entry checks that a spec's artifacts agree with the tree, one verdict row per check"},
	{mode: modeReview, selectors: []string{flagCharge}, required: []string{flagBase, flagTip}, kind: opLegacyCharge,
		usage: "review <slug> --charge --base <commit> --source-tip <commit>", description: "legacy review charge that names every omitted source"},
	{mode: modeReview, selectors: []string{flagCharge, flagFull}, required: []string{flagBase, flagTip}, kind: opLegacyCharge,
		usage: "review <slug> --charge --base <commit> --source-tip <commit> --full", description: "legacy review charge that inlines every source"},
	{mode: modeBuild, optional: []string{flagBase, flagTip}, kind: opVerdict,
		usage: "build <slug> [--base <commit>] [--source-tip <commit>]", description: "build-entry checks that a spec's artifacts agree with the tree, one verdict row per check"},
	{mode: modeBuild, selectors: []string{flagCharge}, required: []string{flagTicket, flagBase, flagTip}, optional: []string{flagQuota}, kind: opPrepareEvidence, bounded: true,
		usage: "build <slug> --charge --ticket <basename> --base <commit> --source-tip <commit> [--max-store-bytes <n>]", description: "prepare one immutable build evidence artifact and print its bounded orientation"},
	{mode: modeBuild, selectors: []string{flagCharge, flagFull}, required: []string{flagTicket, flagBase, flagTip}, kind: opLegacyCharge,
		usage: "build <slug> --charge --ticket <basename> --base <commit> --source-tip <commit> --full", description: "legacy build charge that inlines every source"},
	{mode: modeBuild, selectors: []string{flagPropose}, required: []string{flagTicket, flagBase, flagTip}, kind: opProposal,
		usage: "build <slug> --propose-writes --ticket <basename> --base <commit> --source-tip <commit>", description: "propose one ticket's Writes: entries from the pinned source"},
	{mode: modeEvidence, optional: []string{flagCursor}, kind: opReadEvidence, bounded: true,
		usage: "evidence <id> [--cursor <cursor>]", description: "print one bounded fragment of a prepared evidence artifact and its exact successor"},
}

// selectorFlags lists every switch flag in declaration order.
var selectorFlags = []string{flagCharge, flagPropose, flagFull}

var grammar = usage.Grammar{
	Cmd:  "bench preflight",
	Help: operationUsage(),
	Flags: []usage.Flag{
		{Name: flagBase, HasValue: true, NoEmptyValue: true},
		{Name: flagTip, HasValue: true, NoEmptyValue: true},
		{Name: flagCharge},
		{Name: flagPropose},
		{Name: flagTicket, HasValue: true, NoEmptyValue: true},
		{Name: flagFull},
		{Name: flagQuota, HasValue: true, NoEmptyValue: true},
		{Name: flagCursor, HasValue: true, NoEmptyValue: true},
	},
	MinArgs: 2,
	MaxArgs: 2,
}

func operationUsage() string {
	var b strings.Builder
	for i, op := range operations {
		prefix := "       bench preflight "
		if i == 0 {
			prefix = "usage: bench preflight "
		}
		b.WriteString(prefix + op.usage + "\n")
	}
	return b.String()
}

// HelpRow is one public help row: the text after `bench preflight` and its description.
type HelpRow struct {
	Suffix, Description string
}

// HelpRows projects the operation registry for the root help inventory.
func HelpRows() []HelpRow {
	rows := make([]HelpRow, len(operations))
	for i, op := range operations {
		rows[i] = HelpRow{Suffix: " " + op.usage, Description: op.description}
	}
	return rows
}

// selectOperation returns the one registered form that accepts mode and the given flags, or
// the usage line that explains why none does.
func selectOperation(mode string, flags map[string]string) (operation, string) {
	present := func(name string) bool { _, ok := flags[name]; return ok }
	var selectors []string
	for _, name := range selectorFlags {
		if present(name) {
			selectors = append(selectors, name)
		}
	}
	var sameMode, otherMode *operation
	modeKnown := false
	for i := range operations {
		op := &operations[i]
		if op.mode == mode {
			modeKnown = true
		}
		if !sameSet(op.selectors, selectors) {
			continue
		}
		if op.mode == mode {
			sameMode = op
		} else if otherMode == nil {
			otherMode = op
		}
	}
	switch {
	case !modeKnown:
		return operation{}, toon.Usage(grammar.Cmd, mode)
	case sameMode != nil:
		for _, name := range sameMode.required {
			if !present(name) {
				return operation{}, requirementLine(*sameMode)
			}
		}
		allowed := append(append(append([]string{}, sameMode.selectors...), sameMode.required...), sameMode.optional...)
		for _, name := range grammarFlagNames() {
			if present(name) && !contains(allowed, name) {
				return operation{}, toon.Usage(grammar.Cmd, name)
			}
		}
		return *sameMode, ""
	case otherMode != nil:
		return operation{}, requirementLine(*otherMode)
	}
	for _, op := range operations {
		if subset(selectors, op.selectors) {
			missing := without(op.selectors, selectors...)
			return operation{}, toon.Usage(grammar.Cmd, strings.Join(selectors, " and ")+" requires "+strings.Join(missing, " and "))
		}
	}
	return operation{}, toon.Usage(grammar.Cmd, strings.Join(selectors, " and ")+" cannot be combined")
}

func requirementLine(op operation) string {
	needs := append([]string{op.mode}, op.required...)
	list := strings.Join(needs[:len(needs)-1], ", ") + ", and " + needs[len(needs)-1]
	if len(needs) == 2 {
		list = needs[0] + " and " + needs[1]
	}
	return toon.Usage(grammar.Cmd, strings.Join(op.selectors, " and ")+" requires "+list)
}

func grammarFlagNames() []string {
	names := make([]string, len(grammar.Flags))
	for i, flag := range grammar.Flags {
		names[i] = flag.Name
	}
	return names
}

func sameSet(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	x, y := append([]string{}, a...), append([]string{}, b...)
	sort.Strings(x)
	sort.Strings(y)
	return strings.Join(x, "\x00") == strings.Join(y, "\x00")
}

func subset(inner, outer []string) bool {
	for _, value := range inner {
		if !contains(outer, value) {
			return false
		}
	}
	return true
}

func without(values []string, drop ...string) []string {
	var kept []string
	for _, value := range values {
		if !contains(drop, value) {
			kept = append(kept, value)
		}
	}
	return kept
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
