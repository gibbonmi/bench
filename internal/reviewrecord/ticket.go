package reviewrecord

import (
	"fmt"
	"slices"
	"strings"
)

// TicketSteps is the order of the record writes of one ticket: the completion steps before
// the first review step.
func TicketSteps() []Step {
	steps := CompletionOrder()
	return steps[:slices.IndexFunc(steps, func(step Step) bool { return step.Form == FormReview })]
}

// TicketCall is the caller's part of one ticket record. Chunk and Ticket name the planned
// chunk and ticket, and Tip is the full commit ID whose plan owes the requirements. Base is the
// full commit ID of the chunk base, or empty to derive it. Results answers one chunk
// verification call for each owed requirement against chunk, the chunk entry after the chunk
// step.
type TicketCall struct {
	Chunk, Ticket, Base, Tip string
	Results                  func(chunk Chunk) ([]VerificationCall, error)
}

// TicketRecord is one ticket record as written. Amendment is nil when the plan digest did not
// move. Action is created, added, or updated when the call freezes the chunk, and empty when the
// call leaves the chunk entry as it is. Chunk is the chunk entry after the write.
type TicketRecord struct {
	Amendment    *Amendment
	Action       string
	Chunk        Chunk
	Verification []Verification
}

// TicketRequirements answers the requirements of chunk in plan, the plan at tip, that ticket
// owes, in plan order. A ticket that owes nothing answers an empty list, also when the chunk
// lists no requirement. It refuses a chunk that the plan does not name, a ticket that the chunk
// does not list, and a plan that is not delegated, whose requirements name no ticket.
func TicketRequirements(plan Plan, tip, chunk, ticket string) ([]Requirement, error) {
	planned, err := plannedChunk(plan, tip, chunk)
	if err != nil {
		return nil, err
	}
	if !slices.Contains(planned.Tickets, ticket) {
		return nil, fmt.Errorf("ticket %s is not a ticket of chunk %s in the plan at %s; the chunk tickets are %s", ticket, chunk, tip, strings.Join(planned.Tickets, ", "))
	}
	if !plan.Delegated() {
		return nil, fmt.Errorf("the requirements of chunk %s in the plan at %s name no ticket; the ticket form reads a version 2 plan", chunk, tip)
	}
	owed := []Requirement{}
	for _, requirement := range planned.Verification {
		if requirement.Ticket == ticket {
			owed = append(owed, requirement)
		}
	}
	return owed, nil
}

// EvidenceID is the evidence ID of requirement under prefix, with each trailing hyphen of the
// prefix removed. A requirement ID that already starts with the prefix and a hyphen is the
// evidence ID, so the prefix occurs once. Any other requirement ID follows the prefix and one
// hyphen.
func EvidenceID(prefix, requirement string) string {
	prefix = strings.TrimRight(prefix, "-")
	if strings.HasPrefix(requirement, prefix+"-") {
		return requirement
	}
	return prefix + "-" + requirement
}

// plannedChunk answers chunk id of plan, the plan at source, or refuses a chunk that the plan
// does not name.
func plannedChunk(plan Plan, source, id string) (*PlannedChunk, error) {
	if planned := findChunk(plan, id); planned != nil {
		return planned, nil
	}
	return nil, fmt.Errorf("chunk %s is not a chunk of the plan at %s", id, source)
}

// RecordTicket writes the entries that call.Ticket owes in one write transaction, with the
// steps in TicketSteps order. The amendment step applies when the record plan digest differs
// from the plan at the tip. The chunk step freezes the chunk at the commit that freezeTip
// answers, and the first call of a chunk creates its entry. The verification step appends one
// entry for each call that call.Results answers against the chunk entry after that step. A
// refusal at any step leaves the record bytes unchanged.
func RecordTicket(root, spec string, call TicketCall) (TicketRecord, error) {
	digest, plan, err := atSource(root, spec, call.Tip)
	if err != nil {
		return TicketRecord{}, err
	}
	if _, err := TicketRequirements(plan, call.Tip, call.Chunk, call.Ticket); err != nil {
		return TicketRecord{}, err
	}
	var result TicketRecord
	start := &starter{spec: spec, tip: call.Tip, plan: plan}
	steps := map[string]func(*Record) error{
		FormAmendment: func(record *Record) error {
			if currentPlan(*record, plan, call.Tip) == nil {
				return nil
			}
			entry, err := changeAmendment(record, plan, call.Tip, nil)
			if err != nil {
				return Repair{Form: FormAmendment, Source: call.Tip, Err: err}
			}
			result.Amendment = &entry
			return nil
		},
		FormChunk: func(record *Record) error {
			tip, err := freezeTip(root, spec, record, call.Chunk, call.Tip)
			if err != nil || tip == "" {
				return err
			}
			base, err := chunkBase(root, *record, plan, call.Chunk, call.Base)
			if err == nil {
				result.Action, _, err = changeChunk(record, plan, digest, call.Chunk, base, tip)
			}
			return err
		},
		FormVerification: func(record *Record) error {
			chunk, err := entryOf(record, call.Chunk)
			if err != nil {
				return err
			}
			calls, err := call.Results(*chunk)
			if err != nil {
				return err
			}
			for _, verification := range calls {
				entry, err := changeVerification(root, spec, record, verification)
				if err != nil {
					return err
				}
				result.Verification = append(result.Verification, entry)
			}
			result.Chunk = *chunk
			return nil
		},
	}
	err = write(root, spec, start.create, func(record *Record) error {
		for _, step := range TicketSteps() {
			if err := steps[step.Form](record); err != nil {
				return err
			}
		}
		return nil
	})
	result.Action = start.action(result.Action)
	return result, err
}

// freezeTip answers the commit at which a ticket call at tip freezes chunk id of record, or
// empty when the call leaves the chunk entry as it is. A chunk with no entry freezes at tip. A
// recorded chunk freezes at the commit that boundCommit answers for tip when that commit
// descends from the chunk tip.
func freezeTip(root, spec string, record *Record, id, tip string) (string, error) {
	entry, missing := entryOf(record, id)
	if missing != nil {
		return tip, nil
	}
	bound, err := boundCommit(root, spec, *entry, tip)
	if err != nil {
		return "", err
	}
	if descends, err := isAncestor(root, entry.Tip, bound); err != nil || !descends {
		return "", err
	}
	return bound, nil
}

// chunkBase answers the base of the freeze of chunk id in record: base when it is not empty,
// then the base of the chunk entry, then the commit that coveredAfter answers for the entry of
// the chunk before id in plan. With none of them, it refuses and asks for --base.
func chunkBase(root string, record Record, plan Plan, id, base string) (string, error) {
	if base != "" {
		return base, nil
	}
	if entry, err := entryOf(&record, id); err == nil {
		return entry.Base, nil
	}
	if i := slices.IndexFunc(plan.Chunks, func(chunk PlannedChunk) bool { return chunk.ID == id }); i > 0 {
		if entry, err := entryOf(&record, plan.Chunks[i-1].ID); err == nil {
			_, commit := coveredAfter(root, record, *entry)
			return commit, nil
		}
	}
	return "", fmt.Errorf("chunk %s has no entry and no recorded predecessor; name its base with --base", id)
}
