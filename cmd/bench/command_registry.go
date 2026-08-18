package main

import (
	"fmt"
	"io"
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

type commandAXIDisposition struct {
	root      bool
	children  []string
	exemption string
}

var axiApprovedRoot = commandAXIDisposition{root: true}

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
	Name       string
	Attachment processAttachment
	AXI        commandAXIDisposition
	Help       []string
	Run        commandHandler
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
	definition, ok := commandByName(args[0])
	if !ok {
		fmt.Fprintf(c.Stderr, "bench: unknown subcommand: %q\n", args[0])
		return 2
	}
	if c.Observe != nil {
		fmt.Fprintln(c.Observe, commandImplementationID(definition))
	}
	if definition.Help != nil {
		return helpCommand(c, definition, args[1:])
	}
	return definition.Run(c, args[1:])
}

func commandImplementationID(definition commandDefinition) string {
	return "command-registry:" + definition.Name
}

func commandByName(name string) (commandDefinition, bool) {
	for _, definition := range commandRegistry {
		if definition.Name == name {
			return definition, true
		}
	}
	return commandDefinition{}, false
}

func renderCommandHelp(definition commandDefinition) string {
	if len(definition.Help) == 0 {
		return ""
	}
	return strings.Join(definition.Help, "\n") + "\n"
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
