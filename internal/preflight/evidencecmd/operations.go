// Package evidencecmd owns the public preflight operation registry and the bounded charge
// evidence commands: build evidence publication, stateless evidence reads, and the shared
// response bound. Package preflight dispatches through it and supplies the preparation
// attempts; this package does not import package preflight.
package evidencecmd

import (
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// Kind names the handler one registered operation dispatches to.
type Kind int

// Operation kinds. Each registered operation dispatches to exactly one kind.
const (
	KindVerdict Kind = iota
	KindLegacyCharge
	KindProposal
	KindPrepareEvidence
	KindReadEvidence
)

// Flag spellings and modes the operations share.
const (
	FlagBase     = "--base"
	FlagTip      = "--source-tip"
	flagCharge   = "--charge"
	flagPropose  = "--propose-writes"
	FlagTicket   = "--ticket"
	FlagFull     = "--full"
	flagQuota    = "--max-store-bytes"
	flagCursor   = "--cursor"
	modeReview   = "review"
	ModeBuild    = "build"
	modeEvidence = "evidence"
)

// Operation is one implemented public preflight form. The registry below is the one source
// for the argument grammar, the preflight help, and the root help rows.
type Operation struct {
	Mode string
	// selectors are the switch flags that choose this form; required and optional list
	// every other flag the form accepts.
	selectors, required, optional []string
	usage, description            string
	Kind                          Kind
	// Bounded forms obey the shared response bound for every response they produce.
	Bounded bool
}

var operations = []Operation{
	{Mode: modeReview, optional: []string{FlagBase, FlagTip}, Kind: KindVerdict,
		usage: "review <slug> [--base <commit>] [--source-tip <commit>]", description: "review-entry checks that a spec's artifacts agree with the tree, one verdict row per check"},
	{Mode: modeReview, selectors: []string{flagCharge}, required: []string{FlagBase, FlagTip}, Kind: KindLegacyCharge,
		usage: "review <slug> --charge --base <commit> --source-tip <commit>", description: "legacy review charge that names every omitted source"},
	{Mode: modeReview, selectors: []string{flagCharge, FlagFull}, required: []string{FlagBase, FlagTip}, Kind: KindLegacyCharge,
		usage: "review <slug> --charge --base <commit> --source-tip <commit> --full", description: "legacy review charge that inlines every source"},
	{Mode: ModeBuild, optional: []string{FlagBase, FlagTip}, Kind: KindVerdict,
		usage: "build <slug> [--base <commit>] [--source-tip <commit>]", description: "build-entry checks that a spec's artifacts agree with the tree, one verdict row per check"},
	{Mode: ModeBuild, selectors: []string{flagCharge}, required: []string{FlagTicket, FlagBase, FlagTip}, optional: []string{flagQuota}, Kind: KindPrepareEvidence, Bounded: true,
		usage: "build <slug> --charge --ticket <basename> --base <commit> --source-tip <commit> [--max-store-bytes <n>]", description: "prepare one immutable build evidence artifact and print its bounded orientation"},
	{Mode: ModeBuild, selectors: []string{flagCharge, FlagFull}, required: []string{FlagTicket, FlagBase, FlagTip}, Kind: KindLegacyCharge,
		usage: "build <slug> --charge --ticket <basename> --base <commit> --source-tip <commit> --full", description: "legacy build charge that inlines every source"},
	{Mode: ModeBuild, selectors: []string{flagPropose}, required: []string{FlagTicket, FlagBase, FlagTip}, Kind: KindProposal,
		usage: "build <slug> --propose-writes --ticket <basename> --base <commit> --source-tip <commit>", description: "propose one ticket's Writes: entries from the pinned source"},
	{Mode: modeEvidence, optional: []string{flagCursor}, Kind: KindReadEvidence, Bounded: true,
		usage: "evidence <id> [--cursor <cursor>]", description: "print one bounded fragment of a prepared evidence artifact and its exact successor"},
}

// selectorFlags lists every switch flag in declaration order.
var selectorFlags = []string{flagCharge, flagPropose, FlagFull}

// Grammar is the preflight argument grammar the operation registry accepts.
var Grammar = usage.Grammar{
	Cmd:  "bench preflight",
	Help: operationUsage(),
	Flags: []usage.Flag{
		{Name: FlagBase, HasValue: true, NoEmptyValue: true},
		{Name: FlagTip, HasValue: true, NoEmptyValue: true},
		{Name: flagCharge},
		{Name: flagPropose},
		{Name: FlagTicket, HasValue: true, NoEmptyValue: true},
		{Name: FlagFull},
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

// Select returns the one registered form that accepts mode and the given flags, or the
// usage line that explains why none does.
func Select(mode string, flags map[string]string) (Operation, string) {
	present := func(name string) bool { _, ok := flags[name]; return ok }
	var selectors []string
	for _, name := range selectorFlags {
		if present(name) {
			selectors = append(selectors, name)
		}
	}
	var sameMode, otherMode *Operation
	modeKnown := false
	for i := range operations {
		op := &operations[i]
		if op.Mode == mode {
			modeKnown = true
		}
		if !sameSet(op.selectors, selectors) {
			continue
		}
		if op.Mode == mode {
			sameMode = op
		} else if otherMode == nil {
			otherMode = op
		}
	}
	switch {
	case !modeKnown:
		return Operation{}, toon.Usage(Grammar.Cmd, mode)
	case sameMode != nil:
		for _, name := range sameMode.required {
			if !present(name) {
				return Operation{}, requirementLine(*sameMode)
			}
		}
		allowed := append(append(append([]string{}, sameMode.selectors...), sameMode.required...), sameMode.optional...)
		for _, name := range grammarFlagNames() {
			if present(name) && !contains(allowed, name) {
				return Operation{}, toon.Usage(Grammar.Cmd, name)
			}
		}
		return *sameMode, ""
	case otherMode != nil:
		return Operation{}, requirementLine(*otherMode)
	}
	for _, op := range operations {
		if subset(selectors, op.selectors) {
			missing := without(op.selectors, selectors...)
			return Operation{}, toon.Usage(Grammar.Cmd, strings.Join(selectors, " and ")+" requires "+strings.Join(missing, " and "))
		}
	}
	return Operation{}, toon.Usage(Grammar.Cmd, strings.Join(selectors, " and ")+" cannot be combined")
}

func requirementLine(op Operation) string {
	needs := append([]string{op.Mode}, op.required...)
	list := strings.Join(needs[:len(needs)-1], ", ") + ", and " + needs[len(needs)-1]
	if len(needs) == 2 {
		list = needs[0] + " and " + needs[1]
	}
	return toon.Usage(Grammar.Cmd, strings.Join(op.selectors, " and ")+" requires "+list)
}

func grammarFlagNames() []string {
	names := make([]string, len(Grammar.Flags))
	for i, flag := range Grammar.Flags {
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
