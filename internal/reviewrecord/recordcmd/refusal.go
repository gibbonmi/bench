package recordcmd

import (
	"errors"
	"strings"

	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/reviewrecord"
)

// cause is the refusal cause of err for the record of spec. A repair names the usage line
// of the form that writes the missing entry, with the slug of spec and each value that the
// refusing owner knows in place of its placeholder. A probe mismatch names the probe flags.
func cause(spec string, err error) string {
	var repair reviewrecord.Repair
	if errors.As(err, &repair) {
		if form := named(repair.Form); form != nil {
			return err.Error() + "; write it with " + command + form.suffix(known(spec, repair))
		}
	}
	if errors.Is(err, reviewrecord.ErrProbe) {
		return err.Error() + "; the probe flags are " + strings.Join(probeFlags, ", ")
	}
	return err.Error()
}

// known maps the slug of spec onto the slug operand and each value that repair knows onto
// the flag that reads it. A spec path with no slug keeps the placeholder. A form shows only
// the flags it declares, so a flag that the repair form lacks drops out there.
func known(spec string, repair reviewrecord.Repair) map[string]string {
	slug, _ := reviewrecord.Slug(spec)
	return map[string]string{refusalroute.FactSlug: slug, "--chunk": repair.Chunk, "--fold": repair.Fold, "--requirement": repair.Requirement, "--ticket": repair.Ticket,
		"--axis": repair.Axis, "--source": repair.Source}
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
