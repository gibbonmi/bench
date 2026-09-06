package main

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
)

type processAttachment string

const (
	attachmentDirect processAttachment = "direct"
	attachmentSystem processAttachment = "system"
	attachmentShip   processAttachment = "ship"
)

type commandDisposition struct {
	Name       string
	Attachment processAttachment
}

type commandHandler func(Command, []string) int

type commandKind uint8

const (
	commandStandard commandKind = iota
	commandHelp
)

type inventoryVisibility uint8

const (
	inventoryUnclassified inventoryVisibility = iota
	inventoryPublic
	inventoryInternal
)

type helpRow struct {
	Order       int
	Prefix      string
	Suffix      string
	Gap         int
	Description string
}

type commandInventory struct {
	Visibility inventoryVisibility
	Rows       []helpRow
}

func publicInventory(rows ...helpRow) commandInventory {
	return commandInventory{Visibility: inventoryPublic, Rows: rows}
}

var internalInventory = commandInventory{Visibility: inventoryInternal}

type commandAXIDisposition struct {
	root      bool
	children  []string
	exemption string
}

var axiApprovedRoot = commandAXIDisposition{root: true}

var commandAliases = map[string]string{
	"--help": "help",
	"-h":     "help",
}

func axiApprovedChildren(children ...string) commandAXIDisposition {
	return commandAXIDisposition{children: children}
}

func axiExempt(reason string) commandAXIDisposition {
	return commandAXIDisposition{exemption: reason}
}

const (
	axiReasonOperational = "operational surface has its own output contract"
	axiReasonMutation    = "state-changing surface is outside the read-only AXI query set"
	axiReasonPlumbing    = "internal plumbing is outside the public AXI query set"
	axiReasonRelease     = "release surface has its own ship-tier contract"
)

type commandDefinition struct {
	Name        string
	Attachment  processAttachment
	AXI         commandAXIDisposition
	Inventory   commandInventory
	WrapperOnly bool
	// Hook marks a hook plumbing verb: a verb a harness hook shim pipes its envelope to.
	// The set is the record's one source for which dispatches open a hook span.
	Hook bool
	Kind commandKind
	Run  commandHandler
}

// leafRootNeed states how one leaf of a nested command family receives the repository
// root. The dispatcher reads this value once per call and resolves the root at that one
// site, so the family has a single producer of the not-in-repo refusal.
type leafRootNeed uint8

const (
	// rootNone gives the leaf no repository context: it resolves what it needs itself.
	rootNone leafRootNeed = iota
	// rootRequired refuses outside a repository before the leaf runs.
	rootRequired
	// rootBoundary passes the empty string outside a repository, so the leaf still
	// answers its grammar there.
	rootBoundary
)

// commandLeaf is one row of a nested command family. The row is the single declaration
// of the leaf: its name, the grammar its help form answers, how it takes the root, and
// the handler that runs it.
type commandLeaf struct {
	Name string
	// Grammar is the grammar `<family> <leaf> --help` answers at exit 0. It is empty
	// for a leaf that takes `--help` as an operand or refuses it, and for a leaf with
	// no grammar constant.
	Grammar string
	Root    leafRootNeed
	Run     func(c Command, root string, args []string) int
}

// dispatchLeafFamily answers the forms that belong to the family rather than to any
// leaf, then routes an explicit leaf. A bare call, a help call, and an unknown leaf are
// all answered before the root is resolved, so a typo creates no worktree and acquires
// no assignment. The help form matches on exactly one argument: a second argument is an
// unknown-argument refusal, not a help request.
func dispatchLeafFamily(c Command, family, familyUsage string, leaves []commandLeaf, args []string) int {
	if len(args) == 0 {
		fmt.Fprint(c.Stdout, familyUsage)
		return 2
	}
	if len(args) == 1 && (args[0] == "help" || args[0] == "--help" || args[0] == "-h") {
		fmt.Fprint(c.Stdout, familyUsage)
		return 0
	}
	for _, leaf := range leaves {
		if leaf.Name != args[0] {
			continue
		}
		root, ok := leafRoot(c, leaf.Root)
		if !ok {
			return 1
		}
		return leaf.Run(c, root, args[1:])
	}
	fmt.Fprintln(c.Stderr, toon.Usage(family, args[0]))
	return 2
}

// leafRoot is the family's one root producer. A required root that cannot resolve prints
// the not-in-repo line here, so no leaf repeats that refusal.
func leafRoot(c Command, need leafRootNeed) (string, bool) {
	switch need {
	case rootRequired:
		root, err := git.Root()
		if err != nil {
			fmt.Fprintln(c.Stderr, toon.NotInRepo())
			return "", false
		}
		return root, true
	case rootBoundary:
		return boundaryRoot(), true
	default:
		return "", true
	}
}

// Command is the in-process production entry for ordinary command behavior.
type Command struct {
	Stdin          io.Reader
	Stdout, Stderr io.Writer
	Executable     string
	Observe        io.Writer
}

// Run executes args through the production command dispatcher.
func (c Command) Run(args []string) int {
	if c.Stdin == nil {
		c.Stdin = emptyReader{}
	}
	if c.Stdout == nil {
		c.Stdout = io.Discard
	}
	if c.Stderr == nil {
		c.Stderr = io.Discard
	}
	if len(args) == 0 {
		args = []string{"status", "--route"}
	}
	name := args[0]
	if canonical, ok := commandAliases[name]; ok {
		name = canonical
	}
	definition, ok := commandByName(name)
	if !ok {
		fmt.Fprintf(c.Stderr, "bench: unknown subcommand: %q\n", args[0])
		return 2
	}
	if c.Observe != nil {
		fmt.Fprintln(c.Observe, commandImplementationID(definition))
	}
	if definition.Kind == commandHelp {
		return helpCommand(c, args[1:])
	}
	// The hook plumbing verbs record here, at their one shared dispatch, so no adapter
	// package opens a span of its own.
	if definition.Hook {
		finishSpan := beginHookSpan(definition.Name)
		exit := definition.Run(c, args[1:])
		finishSpan(exit)
		return exit
	}
	return definition.Run(c, args[1:])
}

func commandImplementationID(definition commandDefinition) string {
	return "command-registry:" + definition.Name
}

func commandByName(name string) (commandDefinition, bool) {
	for _, definition := range commandRegistry {
		if definition.Name == name && !definition.WrapperOnly {
			return definition, true
		}
	}
	return commandDefinition{}, false
}

const helpInventoryTitle = "bench — Pocock pipeline meets Kun Chen substrate, gated by your invariants."

func renderCommandHelp() string {
	type orderedRow struct {
		order int
		text  string
	}
	var rows []orderedRow
	for _, definition := range commandRegistry {
		if definition.Inventory.Visibility != inventoryPublic {
			continue
		}
		for _, row := range definition.Inventory.Rows {
			prefix := row.Prefix
			if prefix == "" {
				prefix = "bench"
			}
			command := prefix + " " + definition.Name + row.Suffix
			text := fmt.Sprintf("  %-25s  %s", command, row.Description)
			if row.Gap > 0 {
				text = "  " + command + strings.Repeat(" ", row.Gap) + row.Description
			}
			rows = append(rows, orderedRow{order: row.Order, text: text})
		}
	}
	sort.SliceStable(rows, func(i, j int) bool { return rows[i].order < rows[j].order })
	lines := make([]string, 1, len(rows)+1)
	lines[0] = helpInventoryTitle
	for _, row := range rows {
		lines = append(lines, row.text)
	}
	return strings.Join(lines, "\n") + "\n"
}

func commandDispositions() []commandDisposition {
	var dispositions []commandDisposition
	for _, definition := range commandRegistry {
		if definition.Attachment == "" {
			continue
		}
		dispositions = append(dispositions, commandDisposition{Name: definition.Name, Attachment: definition.Attachment})
	}
	return dispositions
}

type emptyReader struct{}

func (emptyReader) Read([]byte) (int, error) { return 0, io.EOF }
