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
	"github.com/gibbonmi/bench/internal/reviewrecord/carry"
	"github.com/gibbonmi/bench/internal/usage"
)

var ticketForm = form{name: reviewrecord.FormTicket, description: "write the entries that one ticket owes from its staged excerpts, and a carried entry for each unstaged requirement that can carry, in one record write; " + writeOrder(reviewrecord.TicketSteps()),
	flags: []flag{{"--chunk", "<id>", true, once}, {"--ticket", "<basename>", true, once}, {"--id-prefix", "<prefix>", true, once}, {"--performer", "<session>", true, once},
		{"--model", "<model>", true, once}, {"--effort", "<effort>", true, once}, {"--ref", "<ref>", true, once}, {"--base", "<commit>", false, once}, {"--tip", "<commit>", false, once}},
	layout: []string{"--chunk", "--ticket", "--id-prefix", "--performer", "--model", "--effort", "--ref", "[", "--base", "]", "[", "--tip", "]"},
	run:    ticket}

// ticket writes the entries that one ticket owes from the stages of its owed requirements and
// the carried entries of the unstaged ones, through reviewrecord.RecordTicket at --tip, HEAD by
// default. Before the write it refuses, in order, what reviewrecord.TicketRequirements refuses
// for the plan at the tip and a planned seam cell in a row that the Covers: line of the ticket
// names. After the write it removes the stages that it used.
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
	owed, err := reviewrecord.TicketRequirements(plan, call.Tip, call.Chunk, call.Ticket)
	if err != nil {
		return f.refuse(err.Error())
	}
	if out, code := refusePlanned(f, root, spec, call.Tip, ticketRows(call.Chunk, call.Ticket)); out != "" {
		return out, code
	}
	runs := ownedRuns{root: root, spec: spec, chunk: call.Chunk, prefix: values["--id-prefix"], owed: owed,
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
	amendment, freeze, ids, carried := "none", "none", []string{}, "none"
	if written.Amendment != nil {
		amendment = written.Amendment.To
	}
	if written.Action != "" {
		freeze = written.Action
	}
	for _, entry := range written.Verification {
		ids = append(ids, entry.ID)
	}
	if len(runs.carried) != 0 {
		carried = strings.Join(runs.carried, ",")
	}
	return row("ticket", []string{"chunk", "ticket", "amendment", "freeze", "base", "tip", "verifications", "carried"},
		[]string{call.Chunk, call.Ticket, amendment, freeze, written.Chunk.Base, written.Chunk.Tip, strings.Join(ids, ","), carried})
}

// ownedStages reads the stage of each owed requirement of chunk that has one, by requirement ID.
func ownedStages(home, root, slug, chunk string, owed []reviewrecord.Requirement) (map[string]responsebound.Staged, error) {
	stages := map[string]responsebound.Staged{}
	for _, requirement := range owed {
		staged, err := responsebound.ReadStaged(home, root, slug, chunk, requirement.ID)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, err
		}
		stages[requirement.ID] = staged
	}
	return stages, nil
}

// ownedRuns is the owed requirements of one ticket call, in owed order, and the stage of each
// staged one, with the root and the spec of the record, its record path, and the evidence
// fields that every entry shares. calls fills carried with the IDs of the carried
// requirements, in owed order.
type ownedRuns struct {
	root, spec, chunk, record, prefix string
	evidence                          reviewrecord.Evidence
	owed                              []reviewrecord.Requirement
	stages                            map[string]responsebound.Staged
	carried                           []string
}

// calls answers the chunk verification call of each owed requirement in owed order: from its
// stage, with the run commit as the source, or else from the carry candidate that carry.Decide
// answers against entry, the chunk entry. It refuses a stage whose run reviewrecord.RunBinding
// does not find current for entry, then a stage whose run had a dirty path other than the
// record path, each as a repair of the excerpt form that names the requirement. It then
// refuses, as one repair of the excerpt form, the unstaged requirements that cannot carry.
func (r *ownedRuns) calls(entry reviewrecord.Chunk) ([]reviewrecord.VerificationCall, error) {
	restage := func(id, format string, args ...any) error {
		return reviewrecord.Repair{Form: reviewrecord.FormExcerpt, Chunk: r.chunk, Requirement: id, Err: fmt.Errorf("requirement %s: "+format, append([]any{id}, args...)...)}
	}
	for _, requirement := range r.owed {
		staged, ok := r.stages[requirement.ID]
		if !ok {
			continue
		}
		_, current, err := reviewrecord.RunBinding(r.root, r.spec, entry, staged.Head)
		if err != nil {
			return nil, err
		}
		if !current {
			return nil, restage(requirement.ID, "the staged run at %s is not on the chunk line %s..%s", staged.Head, entry.Base, entry.Tip)
		}
	}
	for _, requirement := range r.owed {
		staged := r.stages[requirement.ID]
		if dirty := slices.DeleteFunc(slices.Clone(staged.Dirty), func(path string) bool { return path == r.record }); len(dirty) != 0 {
			return nil, restage(requirement.ID, "the staged run had the dirty path %s", dirty[0])
		}
	}
	candidates, missing, reasons := map[string]reviewrecord.Verification{}, []string{}, ""
	for _, requirement := range r.owed {
		if _, ok := r.stages[requirement.ID]; ok {
			continue
		}
		candidate, reason, err := carry.Decide(r.root, r.spec, entry, requirement)
		if err != nil {
			return nil, err
		}
		if reason != "" {
			missing, reasons = append(missing, requirement.ID), reasons+fmt.Sprintf("; carry: %s (%s)", requirement.ID, reason)
			continue
		}
		candidates[requirement.ID] = candidate
	}
	if len(missing) != 0 {
		return nil, reviewrecord.Repair{Form: reviewrecord.FormExcerpt, Chunk: r.chunk, Err: fmt.Errorf("these owed requirements have no staged excerpt: %s%s", strings.Join(missing, ", "), reasons)}
	}
	calls := []reviewrecord.VerificationCall{}
	for _, requirement := range r.owed {
		call := reviewrecord.VerificationCall{Chunk: r.chunk, Requirement: requirement.ID, Evidence: r.evidence}
		call.Evidence.ID = reviewrecord.EvidenceID(r.prefix, requirement.ID)
		if candidate, ok := candidates[requirement.ID]; ok {
			// The carried entry states the run of the candidate, and the call states who carried it.
			call.Source, call.ExitCode, call.Probe = candidate.SourceCommit, *candidate.ExitCode, candidate.Probe
			call.Evidence.NativeRef = reviewrecord.NativeRef{Ref: "carried:" + candidate.ID, Digest: candidate.NativeRef.Digest, Excerpt: candidate.NativeRef.Excerpt}
			calls, r.carried = append(calls, call), append(r.carried, requirement.ID)
			continue
		}
		staged := r.stages[requirement.ID]
		call.Source, call.ExitCode, call.Evidence.NativeRef = staged.Head, staged.Exit, nativeOf(r.evidence.NativeRef.Ref, staged.Output)
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
