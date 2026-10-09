package recordcmd

import (
	"errors"
	"strings"

	"github.com/gibbonmi/bench/internal/reviewrecord"
)

// cause is the refusal cause of err. A repair names the usage line of the form that writes
// the missing entry, with each value that the refusing owner knows in place of its
// placeholder. A probe mismatch names the probe flags.
func cause(err error) string {
	var repair reviewrecord.Repair
	switch {
	case errors.As(err, &repair) && named(repair.Form) != nil:
		return err.Error() + "; write it with " + command + named(repair.Form).suffix(known(repair))
	case errors.Is(err, reviewrecord.ErrProbe):
		return err.Error() + "; the probe flags are " + strings.Join(probeFlags, ", ")
	}
	return err.Error()
}

// known maps each value that repair knows onto the flag that reads it. A form shows only
// the flags it declares, so a flag that the repair form lacks drops out there.
func known(repair reviewrecord.Repair) map[string]string {
	return map[string]string{"--chunk": repair.Chunk, "--requirement": repair.Requirement, "--ticket": repair.Ticket, "--axis": repair.Axis, "--source": repair.Source}
}

// completionOrder is the completion form's note on the order of the record writes that
// the completion needs, read from the repair owner.
func completionOrder() string {
	steps := []string{}
	for _, step := range reviewrecord.CompletionOrder() {
		if step.Final {
			step.Form = "final " + step.Form
		}
		steps = append(steps, step.Form)
	}
	return "write in the order " + strings.Join(steps, ", ")
}
