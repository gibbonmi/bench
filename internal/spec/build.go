package spec

import (
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// BuildInvocation is the parsed, lifecycle-neutral shape of one spec build command.
type BuildInvocation struct {
	Operation, Slug string
	Flags           map[string]string
}

var buildOperations = map[string]usage.Grammar{
	"start":      buildGrammar("start", nil),
	"assign":     buildGrammar("assign", []usage.Flag{buildValueFlag("--ticket"), buildValueFlag("--request")}),
	"checkpoint": buildGrammar("checkpoint", []usage.Flag{buildValueFlag("--assignment"), buildValueFlag("--evidence")}),
	"integrate":  buildGrammar("integrate", []usage.Flag{buildValueFlag("--assignment")}),
	"review":     buildGrammar("review", []usage.Flag{buildValueFlag("--evidence")}),
	"status":     buildGrammar("status", []usage.Flag{{Name: "--full"}}),
	"promote":    buildGrammar("promote", nil),
	"abandon":    buildGrammar("abandon", []usage.Flag{buildValueFlag("--apply")}),
}

func buildValueFlag(name string) usage.Flag {
	return usage.Flag{Name: name, HasValue: true, NoEmptyValue: true}
}

func buildGrammar(operation string, flags []usage.Flag) usage.Grammar {
	help := "usage: bench spec build " + operation + " <slug>"
	suffixes := map[string]string{
		"assign": " --ticket <ticket> --request <id>", "checkpoint": " --assignment <id> --evidence <receipt>",
		"integrate": " --assignment <id>", "review": " --evidence <receipt>", "status": " [--full]", "abandon": " [--apply <fingerprint>]",
	}
	return usage.Grammar{Cmd: "bench spec build " + operation, Help: help + suffixes[operation], Flags: flags, MinArgs: 1, MaxArgs: 1}
}

// ParseBuild parses exactly the eight public operations without interpreting flag values
// or text after `--` as lifecycle routing.
func ParseBuild(args []string) (BuildInvocation, string, int) {
	if len(args) == 0 {
		return BuildInvocation{}, toon.MissingArg("bench spec build", "start|assign|checkpoint|integrate|review|status|promote|abandon") + "\n", 2
	}
	operation := args[0]
	grammar, known := buildOperations[operation]
	if !known {
		return BuildInvocation{}, toon.Usage("bench spec build", operation) + "\n", 2
	}
	parsed, line, code := usage.Parse(grammar, args[1:])
	if line != "" {
		return BuildInvocation{}, line + "\n", code
	}
	if parsed.PositionalsBeforeTerminator != 1 {
		return BuildInvocation{}, toon.Usage(grammar.Cmd, parsed.Positionals[0]) + "\n", 2
	}
	required := map[string][]string{
		"assign": {"--ticket", "--request"}, "checkpoint": {"--assignment", "--evidence"},
		"integrate": {"--assignment"}, "review": {"--evidence"},
	}
	for _, flag := range required[operation] {
		if _, ok := parsed.Flags[flag]; !ok {
			return BuildInvocation{}, toon.MissingArg(grammar.Cmd, flag) + "\n", 2
		}
	}
	return BuildInvocation{Operation: operation, Slug: parsed.Positionals[0], Flags: parsed.Flags}, "", 0
}
