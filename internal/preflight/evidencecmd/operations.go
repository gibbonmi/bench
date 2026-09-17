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
	KindVerifyEvidence
	KindCurrentEvidence
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
	flagSource   = "--source"
	flagVerify   = "--verify"
	flagCurrent  = "--check-current"
	modeReview   = "review"
	ModeBuild    = "build"
	modeEvidence = "evidence"
)

// flagSpec is one registered preflight flag. A switch has no placeholder; a valued flag
// names the placeholder its usage text shows.
type flagSpec struct {
	name, placeholder string
}

// flagTable is the one flag registry. The argument grammar, the selector switches, and
// every usage line derive from it.
var flagTable = []flagSpec{
	{FlagBase, "<commit>"},
	{FlagTip, "<commit>"},
	{flagCharge, ""},
	{flagPropose, ""},
	{FlagTicket, "<basename>"},
	{FlagFull, ""},
	{flagQuota, "<n>"},
	{flagCursor, "<cursor>"},
	{flagSource, "<source-id>"},
	{flagVerify, ""},
	{flagCurrent, ""},
}

// modeOperands names the positional operand each mode takes.
var modeOperands = map[string]string{modeReview: "<slug>", ModeBuild: "<slug>", modeEvidence: "<id>"}

// Operation is one implemented public preflight form. The registry below is the one source
// for the argument grammar, the preflight help, and the root help rows.
type Operation struct {
	Mode string
	// selectors are the switch flags that choose this form; required and optional list
	// every other flag the form accepts.
	selectors, required, optional []string
	description                   string
	Kind                          Kind
	// Bounded forms obey the shared response bound for every response they produce.
	Bounded bool
}

var operations = []Operation{
	{Mode: modeReview, optional: []string{FlagBase, FlagTip}, Kind: KindVerdict,
		description: "review-entry checks that a spec's artifacts agree with the tree, one verdict row per check"},
	{Mode: modeReview, selectors: []string{flagCharge}, required: []string{FlagBase, FlagTip}, Kind: KindLegacyCharge,
		description: "legacy review charge that names every omitted source"},
	{Mode: modeReview, selectors: []string{flagCharge, FlagFull}, required: []string{FlagBase, FlagTip}, Kind: KindLegacyCharge,
		description: "legacy review charge that inlines every source"},
	{Mode: ModeBuild, optional: []string{FlagBase, FlagTip}, Kind: KindVerdict,
		description: "build-entry checks that a spec's artifacts agree with the tree, one verdict row per check"},
	{Mode: ModeBuild, selectors: []string{flagCharge}, required: []string{FlagTicket, FlagBase, FlagTip}, optional: []string{flagQuota}, Kind: KindPrepareEvidence, Bounded: true,
		description: "prepare one immutable build evidence artifact and print its bounded orientation"},
	{Mode: ModeBuild, selectors: []string{flagCharge, FlagFull}, required: []string{FlagTicket, FlagBase, FlagTip}, Kind: KindLegacyCharge,
		description: "legacy build charge that inlines every source"},
	{Mode: ModeBuild, selectors: []string{flagPropose}, required: []string{FlagTicket, FlagBase, FlagTip}, Kind: KindProposal,
		description: "propose one ticket's Writes: entries from the pinned source"},
	{Mode: modeEvidence, optional: []string{flagCursor}, Kind: KindReadEvidence, Bounded: true,
		description: "print one bounded fragment of a prepared evidence artifact and its exact successor"},
	{Mode: modeEvidence, selectors: []string{flagSource}, optional: []string{flagCursor}, Kind: KindReadEvidence, Bounded: true,
		description: "print one bounded fragment of one declared source stream and its exact successor"},
	{Mode: modeEvidence, selectors: []string{flagVerify}, Kind: KindVerifyEvidence, Bounded: true,
		description: "verify every stored page and source digest of a prepared evidence artifact"},
	{Mode: modeEvidence, selectors: []string{flagCurrent}, Kind: KindCurrentEvidence, Bounded: true,
		description: "bind a prepared evidence artifact to the current assignment and source pair"},
}

// selectorFlags lists every flag that chooses a registered form, in flag registry order.
// The operations below are the one source of that set, so a flag is a selector exactly
// when some form selects on it.
var selectorFlags = operationSelectors()

func operationSelectors() []string {
	var names []string
	for _, flag := range flagTable {
		for _, op := range operations {
			if contains(op.selectors, flag.name) {
				names = append(names, flag.name)
				break
			}
		}
	}
	return names
}

// Grammar is the preflight argument grammar the operation registry accepts.
var Grammar = usage.Grammar{
	Cmd:     "bench preflight",
	Help:    operationUsage(),
	Flags:   grammarFlags(),
	MinArgs: 2,
	MaxArgs: 2,
}

func grammarFlags() []usage.Flag {
	flags := make([]usage.Flag, len(flagTable))
	for i, flag := range flagTable {
		valued := flag.placeholder != ""
		flags[i] = usage.Flag{Name: flag.name, HasValue: valued, NoEmptyValue: valued}
	}
	return flags
}

// usageLine renders one form after `bench preflight`: the mode and its operand, the first
// selector, the required flags, the remaining selectors, then the bracketed optional flags.
func (op Operation) usageLine() string {
	terms := []string{op.Mode, modeOperands[op.Mode]}
	if len(op.selectors) > 0 {
		terms = append(terms, flagTerm(op.selectors[0]))
	}
	for _, name := range op.required {
		terms = append(terms, flagTerm(name))
	}
	for _, name := range op.selectors[min(1, len(op.selectors)):] {
		terms = append(terms, flagTerm(name))
	}
	for _, name := range op.optional {
		terms = append(terms, "["+flagTerm(name)+"]")
	}
	return strings.Join(terms, " ")
}

// flagTerm is a flag followed by its registered placeholder, if it takes a value.
func flagTerm(name string) string {
	for _, flag := range flagTable {
		if flag.name == name && flag.placeholder != "" {
			return name + " " + flag.placeholder
		}
	}
	return name
}

func operationUsage() string {
	var b strings.Builder
	for i, op := range operations {
		prefix := "       bench preflight "
		if i == 0 {
			prefix = "usage: bench preflight "
		}
		b.WriteString(prefix + op.usageLine() + "\n")
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
		rows[i] = HelpRow{Suffix: " " + op.usageLine(), Description: op.description}
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
	names := make([]string, len(flagTable))
	for i, flag := range flagTable {
		names[i] = flag.name
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
