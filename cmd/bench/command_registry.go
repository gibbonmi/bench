package main

import (
	"fmt"
	"io"
	"sort"
	"strings"
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
	Kind        commandKind
	Run         commandHandler
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
