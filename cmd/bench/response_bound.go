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
	// form limits the flag to the form position, so an operand equal to the flag does not
	// exempt the call.
	form bool
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

// boundExemptForm exempts only a call whose form position holds form.
func boundExemptForm(form, reason string) boundDisposition {
	return boundDisposition{exempt: reason, flag: form, form: true}
}

// call answers the disposition of one call with args. An exemption limited to a flag
// holds only when the call carries that flag. An exemption limited to a form holds only
// when the form position holds that form.
func (d boundDisposition) call(args []string) boundDisposition {
	switch {
	case d.flag == "":
		return d
	case d.form && formArgument(args) == d.flag, !d.form && slices.Contains(args, d.flag):
		return d
	}
	return boundResponse
}

// formArgument answers the form of one call, or the empty string for a call with no form.
// The form is the first argument after a leading tree-target pair, as the dispatcher reads it.
func formArgument(args []string) string {
	if leadsWithTreeTarget(args) {
		args = args[min(2, len(args)):]
	}
	if len(args) == 0 {
		return ""
	}
	return args[0]
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
