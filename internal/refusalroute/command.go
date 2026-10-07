package refusalroute

import (
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// grammar is the declared argument shape usage.Parse enforces for `bench recovery`: the
// matrix takes no argument.
var grammar = usage.Grammar{
	Cmd:  "bench recovery",
	Help: "usage: bench recovery",
}

// RecoveryCommand implements `bench recovery`, the recovery matrix. It reads only the
// compiled registry: one row for each registered face, in registry order. The route cell
// renders the face's route over no facts, so each slot prints its placeholder.
func RecoveryCommand(args []string) (string, int) {
	if _, line, code := usage.Parse(grammar, args); line != "" {
		return line + "\n", code
	}
	rows := make([][]string, 0, len(inventory))
	for _, face := range inventory {
		rows = append(rows, []string{string(face.Verb), face.Name, string(face.Authority), face.Render(Facts{})})
	}
	out, err := toon.Table("recovery", []string{"verb", "face", "authority", "route"}, rows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	return out, 0
}
