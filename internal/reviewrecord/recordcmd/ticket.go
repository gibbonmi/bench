package recordcmd

import (
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/responsebound"
	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/usage"
)

var ticketForm = form{name: reviewrecord.FormTicket, description: "write the entries that one ticket owes from its staged excerpts in one record write; " + writeOrder(reviewrecord.TicketSteps()),
	flags: []flag{{"--chunk", "<id>", true, once}, {"--ticket", "<basename>", true, once}, {"--id-prefix", "<prefix>", true, once}, {"--performer", "<session>", true, once},
		{"--model", "<model>", true, once}, {"--effort", "<effort>", true, once}, {"--ref", "<ref>", true, once}, {"--base", "<commit>", false, once}, {"--tip", "<commit>", false, once}},
	layout: []string{"--chunk", "--ticket", "--id-prefix", "--performer", "--model", "--effort", "--ref", "[", "--base", "]", "[", "--tip", "]"},
	run:    ticket}

// ticket writes the entries that one ticket owes from the stages of its owed requirements,
// through reviewrecord.RecordTicket at --tip, HEAD by default. Before the write it refuses, in
// order, a ticket or a plan that owes no requirement, a planned seam cell for the last ticket,
// and an owed requirement with no stage. After the write it removes the stages that it used.
func ticket(f form, root, spec string, parsed usage.Result) (string, int) {
	values := maps.Clone(parsed.Flags)
	if _, ok := values["--tip"]; !ok {
		values["--tip"] = "HEAD"
	}
	call := reviewrecord.TicketCall{Chunk: values["--chunk"], Ticket: values["--ticket"]}
	var out string
	var code int
	if call.Tip, out, code = commit(f, root, values, "--tip"); out != "" {
		return out, code
	}
	if _, ok := values["--base"]; ok {
		if call.Base, out, code = commit(f, root, values, "--base"); out != "" {
			return out, code
		}
	}
	_, plan, err := reviewrecord.ReadSpecAtCommit(root, spec, call.Tip)
	if err != nil {
		return f.refuse(err.Error())
	}
	owed, last, err := reviewrecord.TicketRequirements(plan, call.Tip, call.Chunk, call.Ticket)
	if err != nil {
		return f.refuse(err.Error())
	}
	if call.Ticket == last {
		if out, code := refusePlanned(f, root, spec, call.Tip, chunkRows(call.Chunk)); out != "" {
			return out, code
		}
	}
	runs := ownedRuns{chunk: call.Chunk, prefix: values["--id-prefix"], owed: owed,
		evidence: reviewrecord.Evidence{Performer: values["--performer"], Model: values["--model"], Effort: values["--effort"], NativeRef: reviewrecord.NativeRef{Ref: values["--ref"]}}}
	if runs.record, err = reviewrecord.RecordPath(spec); err != nil {
		return f.refuse(err.Error())
	}
	home, slug := benchhome.Dir(), parsed.Positionals[0]
	if runs.stages, err = ownedStages(home, root, slug, call.Chunk, owed); err != nil {
		return f.refuse(cause(spec, err))
	}
	call.Results = runs.calls
	written, err := reviewrecord.RecordTicket(root, spec, call)
	if err != nil {
		return f.refuse(cause(spec, err))
	}
	for _, requirement := range owed {
		if err := responsebound.Unstage(home, root, slug, call.Chunk, requirement.ID); err != nil {
			return f.refuse("the record is written; " + err.Error())
		}
	}
	amendment, freeze, ids := "none", "none", []string{}
	if written.Amendment != nil {
		amendment = written.Amendment.To
	}
	if written.Action != "" {
		freeze = written.Action
	}
	for _, entry := range written.Verification {
		ids = append(ids, entry.ID)
	}
	return row("ticket", []string{"chunk", "ticket", "amendment", "freeze", "base", "tip", "verifications"},
		[]string{call.Chunk, call.Ticket, amendment, freeze, written.Chunk.Base, written.Chunk.Tip, strings.Join(ids, ",")})
}

// ownedStages reads the stage of each owed requirement of chunk, in owed order. One refusal
// names every owed requirement with no stage, as a repair of the excerpt form.
func ownedStages(home, root, slug, chunk string, owed []reviewrecord.Requirement) ([]responsebound.Staged, error) {
	stages, missing := []responsebound.Staged{}, []string{}
	for _, requirement := range owed {
		staged, err := responsebound.ReadStaged(home, root, slug, chunk, requirement.ID)
		if errors.Is(err, fs.ErrNotExist) {
			missing = append(missing, requirement.ID)
			continue
		}
		if err != nil {
			return nil, err
		}
		stages = append(stages, staged)
	}
	if len(missing) != 0 {
		return nil, reviewrecord.Repair{Form: reviewrecord.FormExcerpt, Chunk: chunk, Err: fmt.Errorf("these owed requirements have no staged excerpt: %s", strings.Join(missing, ", "))}
	}
	return stages, nil
}

// ownedRuns is the stage of each owed requirement of one ticket call, in owed order, with the
// record path of the spec and the evidence fields that every entry shares.
type ownedRuns struct {
	chunk, record, prefix string
	evidence              reviewrecord.Evidence
	owed                  []reviewrecord.Requirement
	stages                []responsebound.Staged
}

// calls answers the chunk verification call of each owed requirement at source from its stage.
// It refuses a stage whose run commit is not source, then a stage whose run had a dirty path
// other than the record path, each as a repair of the excerpt form that names the requirement.
// Each evidence ID joins the prefix and the requirement ID with a hyphen.
func (r ownedRuns) calls(source string) ([]reviewrecord.VerificationCall, error) {
	restage := func(i int, format string, args ...any) error {
		id := r.owed[i].ID
		return reviewrecord.Repair{Form: reviewrecord.FormExcerpt, Chunk: r.chunk, Requirement: id, Err: fmt.Errorf("requirement %s: "+format, append([]any{id}, args...)...)}
	}
	for i, staged := range r.stages {
		if staged.Head != source {
			return nil, restage(i, "the staged run is at %s, not the chunk source %s", staged.Head, source)
		}
	}
	for i, staged := range r.stages {
		if dirty := slices.DeleteFunc(slices.Clone(staged.Dirty), func(path string) bool { return path == r.record }); len(dirty) != 0 {
			return nil, restage(i, "the staged run had the dirty path %s", dirty[0])
		}
	}
	calls := []reviewrecord.VerificationCall{}
	for i, staged := range r.stages {
		requirement := r.owed[i]
		call := reviewrecord.VerificationCall{Chunk: r.chunk, Source: source, Requirement: requirement.ID, ExitCode: staged.Exit, Evidence: r.evidence}
		call.Evidence.ID, call.Evidence.NativeRef = r.prefix+"-"+requirement.ID, nativeOf(r.evidence.NativeRef.Ref, staged.Output)
		if requirement.Probe != "" {
			var out string
			if call.ExitCode, call.Probe, out = stagedProbe(staged); out != "" {
				return nil, fmt.Errorf("requirement %s: %s", requirement.ID, out)
			}
		}
		calls = append(calls, call)
	}
	return calls, nil
}

// stagedProbe reads the probe of a staged answer. A single-form answer binds the staged exit,
// and a mutant-set row binds its own exit, so the entry records the exit that the
// verification form checks for the same answer.
func stagedProbe(staged responsebound.Staged) (int, *reviewrecord.Probe, string) {
	exit := staged.Exit
	if staged.Mutant != 0 {
		answer, out := readImported(staged.Output, staged.Mutant, retainedSource)
		if out != "" {
			return 0, nil, out
		}
		exit = answer.exit
	}
	probe, out := importedProbe(staged.Output, exit, staged.Mutant)
	return exit, probe, out
}
