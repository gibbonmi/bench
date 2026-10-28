package recordcmd

import (
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/benchhome"
	benchcommit "github.com/gibbonmi/bench/internal/commit"
	"github.com/gibbonmi/bench/internal/probe"
	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/responsebound"
	"github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/reviewrecord/carry"
	"github.com/gibbonmi/bench/internal/usage"
)

var ticketForm = form{name: reviewrecord.FormTicket, description: "write the entries that one ticket owes from its staged excerpts, a carried entry for each unstaged requirement that can carry, " +
	"and an entry from the retained bench probe response for the one requirement left, with --mutant naming its row, in one record write; " +
	"the first ticket of a chunk passes --base, the main tip for the first chunk and the recorded predecessor tip after it; " + refusalroute.NextField + "= prints the record commit; " + writeOrder(reviewrecord.TicketSteps()),
	flags: []flag{{"--chunk", "<id>", true, once}, {"--ticket", "<basename>", true, once}, {"--id-prefix", "<prefix>", true, once}, {"--performer", "<session>", true, once},
		{"--model", "<model>", true, once}, {"--effort", "<effort>", true, once}, {"--ref", "<ref>", true, once}, {"--base", "<commit>", false, once}, {"--tip", "<commit>", false, once},
		{mutantFlag, "<n>", false, once}},
	layout: []string{"--chunk", "--ticket", "--id-prefix", "--performer", "--model", "--effort", "--ref", "[", "--base", "]", "[", "--tip", "]", "[", mutantFlag, "]"},
	valid:  mutantRow, run: ticket}

// ticket writes the entries that one ticket owes from the stages of its owed requirements and
// the carried entries of the unstaged ones, through reviewrecord.RecordTicket at --tip, HEAD by
// default. Before the write it refuses, in order, what reviewrecord.TicketRequirements refuses
// for the plan at the tip and a planned seam cell in a row that the Covers: line of the ticket
// names. After the write it removes the stages that it used, and after its table it prints the
// record commit with the slug and the chunk base on one next= line.
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
	home, slug := benchhome.Dir(), parsed.Positionals[0]
	runs := ownedRuns{home: home, root: root, spec: spec, chunk: call.Chunk, prefix: values["--id-prefix"], owed: owed,
		evidence: reviewrecord.Evidence{Performer: values["--performer"], Model: values["--model"], Effort: values["--effort"], NativeRef: reviewrecord.NativeRef{Ref: values["--ref"]}}}
	runs.mutant, _ = integer(values, mutantFlag)
	if runs.record, err = reviewrecord.RecordPath(spec); err != nil {
		return f.refuse(err.Error())
	}
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
	table, code := row("ticket", []string{"chunk", "ticket", "amendment", "freeze", "base", "tip", "verifications", "carried"},
		[]string{call.Chunk, call.Ticket, amendment, freeze, written.Chunk.Base, written.Chunk.Tip, strings.Join(ids, ","), carried})
	if code != 0 {
		return table, code
	}
	chain := benchcommit.ChainArgs(slug, written.Chunk.Base)
	next := refusalroute.ChainCommit(map[string]string{refusalroute.FactChain: strings.Join(chain, " "), refusalroute.FactRecord: runs.record})
	return table + refusalroute.NextField + "=" + next + "\n", 0
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
// staged one, with the Bench home, the root and the spec of the record, its record path, the
// --mutant row of the retained response, and the evidence fields that every entry shares. calls
// fills carried with the IDs of the carried requirements, in owed order.
type ownedRuns struct {
	home, root, spec, chunk, record, prefix string
	mutant                                  int
	evidence                                reviewrecord.Evidence
	owed                                    []reviewrecord.Requirement
	stages                                  map[string]responsebound.Staged
	carried                                 []string
}

// calls answers the chunk verification call of each owed requirement in owed order: from its
// stage, with the run commit as the source, or else from the carry candidate that carry.Decide
// answers against entry, the chunk entry. When exactly one owed requirement has neither a stage
// nor a carry and the retained call is `bench probe`, stageRetained stages that response for it.
// The count maps the response to one requirement only; no check compares the retained call with
// that requirement's planned command. It refuses a stage whose run reviewrecord.RunBinding does
// not find current for entry, then a stage whose run had a dirty path other than the record
// path, each as a repair of the excerpt form that names the requirement. It then refuses, as one
// repair of the excerpt form, the unstaged requirements that cannot carry.
func (r *ownedRuns) calls(entry reviewrecord.Chunk) ([]reviewrecord.VerificationCall, error) {
	restage := func(id, format string, args ...any) error {
		return reviewrecord.Repair{Form: reviewrecord.FormExcerpt, Chunk: r.chunk, Requirement: id, Err: fmt.Errorf("requirement %s: "+format, append([]any{id}, args...)...)}
	}
	candidates, unstaged, reasons := map[string]reviewrecord.Verification{}, []reviewrecord.Requirement{}, ""
	for _, requirement := range r.owed {
		if _, ok := r.stages[requirement.ID]; ok {
			continue
		}
		candidate, reason, err := carry.Decide(r.root, r.spec, entry, requirement)
		if err != nil {
			return nil, err
		}
		if reason != "" {
			unstaged, reasons = append(unstaged, requirement), reasons+fmt.Sprintf("; carry: %s (%s)", requirement.ID, reason)
			continue
		}
		candidates[requirement.ID] = candidate
	}
	var retained responsebound.Retained
	if len(unstaged) == 1 {
		var err error
		if retained, err = responsebound.ReadRetained(r.home, r.root); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
	}
	if len(unstaged) == 1 && retained.Verb == probe.Verb {
		staged, out := stageRetained(unstaged[0], retained, r.mutant)
		if out != "" {
			return nil, fmt.Errorf("requirement %s: %s", unstaged[0].ID, out)
		}
		r.stages[unstaged[0].ID], unstaged = staged, nil
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
	if len(unstaged) != 0 {
		missing := []string{}
		for _, requirement := range unstaged {
			missing = append(missing, requirement.ID)
		}
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
	imported, out := importedProbe(staged.Output, exit, staged.Mutant)
	return exit, imported, out
}
