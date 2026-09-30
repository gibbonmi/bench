package treetarget

import (
	"strings"

	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/toon"
)

// Operand is the grammar operand of the tree-target flag. The missing-value refusal and
// the root help row both read it here.
const Operand = "<label|primary>"

// targetValue splits the arguments after the tree-target flag into the target value and the
// verb's own arguments. It answers the exit-2 usage line of command when the arguments name
// no value. A missing or empty value never reaches the lookup, and a value that starts with
// a dash is a misplaced flag, not a label.
func targetValue(command string, args []string) (value string, rest []string, usage string) {
	if len(args) == 0 || args[0] == "" {
		return "", nil, toon.MissingArg(command, Operand)
	}
	if strings.HasPrefix(args[0], "-") {
		return "", nil, valueUsage(command, args[0])
	}
	return args[0], args[1:], ""
}

// valueUsage is the exit-2 usage line of a value that names no tree target. The value is
// caller text that can carry a control character, so it always prints escaped.
func valueUsage(command, value string) string {
	return toon.Usage(command, sanitize.Controls(value))
}
