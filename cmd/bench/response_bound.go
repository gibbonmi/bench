package main

import "slices"

// boundDisposition is how one public command or leaf takes the response bound: bounded, or
// exempt with a reason. The zero value declares nothing, so a new public entry cannot skip
// the bound silently. A hook or internal plumbing entry declares nothing and stays outside
// the bound.
type boundDisposition struct {
	bounded bool
	// exempt names why the entry prints its complete response.
	exempt string
	// flag limits the exemption to a call that carries this argument. Any other call of
	// the entry is bounded.
	flag string
}

// boundResponse holds the entry's stdout and stderr to the response bound.
var boundResponse = boundDisposition{bounded: true}

// The exempt set is closed. Each member names its reason here, and the registry test
// compares the set with its closed members.
const (
	boundReasonHelp     = "a help form prints the complete inventory or grammar in one call"
	boundReasonArtifact = "stdout is a machine artifact"
	boundReasonTerminal = "the command drives an interactive terminal"
	boundReasonShipTier = "CI runs it, and a spill file on a discarded runner loses the evidence"
	boundReasonWrapper  = "only the wrapper runs it, so it never reaches Command.Run"
)

// boundHelpForm is the disposition of every help form.
var boundHelpForm = boundExempt(boundReasonHelp)

func boundExempt(reason string) boundDisposition { return boundDisposition{exempt: reason} }

// boundExemptWith exempts only a call that carries flag.
func boundExemptWith(flag, reason string) boundDisposition {
	return boundDisposition{exempt: reason, flag: flag}
}

// call answers the disposition of one call with args. An exemption limited to a flag
// holds only when the call carries that flag.
func (d boundDisposition) call(args []string) boundDisposition {
	if d.flag != "" && !slices.Contains(args, d.flag) {
		return boundResponse
	}
	return d
}

// helpArgument reports one of the three spellings of a help request.
func helpArgument(arg string) bool { return arg == "help" || arg == "--help" || arg == "-h" }

// helpForm reports a help form: `bench help`, or a command or a leaf followed by exactly
// one help argument. A longer argument list is no help form, so an exec child's own
// `--help` stays bounded.
func (definition commandDefinition) helpForm(args []string) bool {
	if definition.Kind == commandHelp {
		return len(args) == 0
	}
	if _, ok := leafNamed(definition.Leaves, args); ok {
		args = args[1:]
	}
	return len(args) == 1 && helpArgument(args[0])
}
