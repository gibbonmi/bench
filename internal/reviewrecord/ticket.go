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
// verification call for each owed requirement at source, the full commit ID of the chunk
// source. The transaction runs Results before any step changes the record.
type TicketCall struct {
	Chunk, Ticket, Base, Tip string
	Results                  func(source string) ([]VerificationCall, error)
}

// TicketRecord is one ticket record as written. Amendment is nil when the plan digest did not
// move. Action is created, added, or updated for the freeze of the last ticket, and empty for an
// earlier ticket. Chunk is the chunk entry after the write.
type TicketRecord struct {
	Amendment    *Amendment
	Action       string
	Chunk        Chunk
	Verification []Verification
}

// TicketRequirements answers the requirements of chunk in plan, the plan at tip, that ticket
// owes, in plan order, and the last ticket of the chunk. It refuses a chunk that the plan does
// not name, a ticket that the chunk does not list, and chunk requirements that name no ticket.
func TicketRequirements(plan Plan, tip, chunk, ticket string) ([]Requirement, string, error) {
	planned, err := plannedChunk(plan, tip, chunk)
	if err != nil {
		return nil, "", err
	}
	if !slices.Contains(planned.Tickets, ticket) {
		return nil, "", fmt.Errorf("ticket %s is not a ticket of chunk %s in the plan at %s; the chunk tickets are %s", ticket, chunk, tip, strings.Join(planned.Tickets, ", "))
	}
	owed, named := []Requirement{}, false
	for _, requirement := range planned.Verification {
		named = named || requirement.Ticket != ""
		if requirement.Ticket == ticket {
			owed = append(owed, requirement)
		}
	}
	if !named {
		return nil, "", fmt.Errorf("the requirements of chunk %s in the plan at %s name no ticket; the ticket form reads a version 2 plan", chunk, tip)
	}
	return owed, planned.Tickets[len(planned.Tickets)-1], nil
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
// from the plan at the tip, and the chunk step freezes the chunk for the last ticket of the
// chunk. The verification step appends one entry for each owed requirement at the chunk
// source: the freeze tip for the last ticket, and the tip of the chunk entry for an earlier
// ticket. An earlier ticket's call needs that entry and never changes it. A refusal at any
// step leaves the record bytes unchanged.
func RecordTicket(root, spec string, call TicketCall) (TicketRecord, error) {
	digest, plan, err := atSource(root, spec, call.Tip)
	if err != nil {
		return TicketRecord{}, err
	}
	_, last, err := TicketRequirements(plan, call.Tip, call.Chunk, call.Ticket)
	if err != nil {
		return TicketRecord{}, err
	}
	var result TicketRecord
	var calls []VerificationCall
	created := false
	create := func(path string) (Record, error) {
		_, err := absent(call.Chunk)(path)
		return Record{}, lastFreezes(call.Chunk, last, err)
	}
	if call.Ticket == last {
		create = func(string) (Record, error) {
			record, err := createRecord(spec, plan, call.Tip)
			created = err == nil
			return record, err
		}
	}
	bind := func(record *Record) error {
		source := call.Tip
		if call.Ticket != last {
			entry, err := entryOf(record, call.Chunk)
			if err != nil {
				return lastFreezes(call.Chunk, last, err)
			}
			source = entry.Tip
		}
		var err error
		calls, err = call.Results(source)
		return err
	}
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
			if call.Ticket != last {
				return nil
			}
			base, err := chunkBase(root, *record, plan, call.Chunk, call.Base)
			if err == nil {
				result.Action, _, err = changeChunk(record, plan, digest, call.Chunk, base, call.Tip)
			}
			return err
		},
		FormVerification: func(record *Record) error {
			for _, verification := range calls {
				entry, err := changeVerification(root, spec, record, verification)
				if err != nil {
					return err
				}
				result.Verification = append(result.Verification, entry)
			}
			chunk, err := entryOf(record, call.Chunk)
			if err == nil {
				result.Chunk = *chunk
			}
			return err
		},
	}
	sequence := []func(*Record) error{bind}
	for _, step := range TicketSteps() {
		sequence = append(sequence, steps[step.Form])
	}
	err = write(root, spec, create, func(record *Record) error {
		for _, step := range sequence {
			if err := step(record); err != nil {
				return err
			}
		}
		return nil
	})
	if created {
		result.Action = "created"
	}
	return result, err
}

// lastFreezes wraps the refusal of an earlier ticket's call whose chunk has no entry as a
// repair of the ticket form for last, the last ticket of chunk, whose call freezes the chunk.
func lastFreezes(chunk, last string, err error) error {
	return Repair{Form: FormTicket, Chunk: chunk, Ticket: last, Err: err}
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
