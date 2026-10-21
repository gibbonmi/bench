package recordcmd

import (
	"errors"
	"strings"

	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/usage"
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

// refuseUninstalled runs the candidate check of a guarded form of a declared build, before
// any record read or ledger write. The declaration comes from the spec at the commit that
// HEAD names in root. It answers no output when the check passes. A refusal that a
// registered face raises prints its route on its own next= line, and the route re-runs the
// caller's own call in the worktree of the active assignment that owns root.
func (f form) refuseUninstalled(root, spec string, parsed usage.Result) (string, int) {
	err := commitrepo.CheckEntry(root, spec)
	var raised refusalroute.Raised
	switch {
	case err == nil:
		return "", 0
	case !errors.As(err, &raised):
		return f.refuse(err.Error())
	}
	owner, _ := intent.AssignmentForWorktree(root)
	slug, _ := reviewrecord.Slug(spec)
	rerun := refusalroute.RecordRerun(map[string]string{refusalroute.FactLabel: owner.Label, refusalroute.FactArguments: f.arguments(slug, parsed)})
	refusal := refusalroute.Printed("", err, map[string]string{refusalroute.FactRerun: rerun})
	out, code := f.refuse(refusal.Sentence)
	return out + refusalroute.NextField + "=" + refusal.Route + "\n", code
}

// completionOrder is the completion form's note on the order of the record writes that
// the completion needs, read from the repair owner.
func completionOrder() string { return writeOrder(reviewrecord.CompletionOrder()) }

// writeOrder is a form's note on the order of the record writes in order.
func writeOrder(order []reviewrecord.Step) string {
	steps := []string{}
	for _, step := range order {
		if step.Final {
			step.Form = "final " + step.Form
		}
		steps = append(steps, step.Form)
	}
	return "write in the order " + strings.Join(steps, ", ")
}
