package spec

import (
	"strings"

	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// BuildInvocation is the parsed, lifecycle-neutral shape of one spec build command.
type BuildInvocation struct {
	Operation, Slug string
	Flags           map[string]string
}

// buildOperationOrder is the one declared lifecycle order for spec build operations.
// buildOperations and the empty-argument diagnostic both derive from it, so neither can
// list an operation the other doesn't, and the diagnostic reads in this order rather than
// Go's randomized map order.
var buildOperationOrder = []string{
	"start", "assign", "checkpoint", "integrate", "review", "status", "promote", "abandon", "reclaim",
}

var buildOperations = func() map[string]usage.Grammar {
	flags := map[string][]usage.Flag{
		"assign":     {buildValueFlag("--ticket"), buildValueFlag("--request"), buildValueFlag("--refresh")},
		"checkpoint": {buildValueFlag("--assignment"), buildValueFlag("--evidence")},
		"integrate":  {buildValueFlag("--assignment")},
		"review":     {buildValueFlag("--evidence")},
		"status":     {{Name: "--full"}},
		"abandon":    {buildValueFlag("--apply")},
		"reclaim":    {buildValueFlag("--apply")},
	}
	m := make(map[string]usage.Grammar, len(buildOperationOrder))
	for _, operation := range buildOperationOrder {
		m[operation] = buildGrammar(operation, flags[operation])
	}
	return m
}()

// buildOperationsHelp is the `|`-separated operation list the empty-argument diagnostic
// prints, derived from buildOperationOrder so it can never drift from the grammar table.
var buildOperationsHelp = strings.Join(buildOperationOrder, "|")

// applySuffix is the help tail every plan-then-apply operation shares.
const applySuffix = " [--apply <fingerprint>]"

func buildValueFlag(name string) usage.Flag {
	return usage.Flag{Name: name, HasValue: true, NoEmptyValue: true}
}

func buildGrammar(operation string, flags []usage.Flag) usage.Grammar {
	help := "usage: bench spec build " + operation + " <slug>"
	suffixes := map[string]string{
		"assign": " --ticket <ticket> --request <id> [--refresh <receipt>]", "checkpoint": " --assignment <id> --evidence <receipt>",
		"integrate": " --assignment <id>", "review": " --evidence <receipt>", "status": " [--full]",
		"abandon": applySuffix, "reclaim": applySuffix,
	}
	return usage.Grammar{Cmd: "bench spec build " + operation, Help: help + suffixes[operation], Flags: flags, MinArgs: 1, MaxArgs: 1}
}

// ParseBuild parses exactly the operations the grammar table declares, without interpreting
// flag values or text after `--` as lifecycle routing.
func ParseBuild(args []string) (BuildInvocation, string, int) {
	if len(args) == 0 {
		return BuildInvocation{}, toon.MissingArg("bench spec build", buildOperationsHelp) + "\n", 2
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
