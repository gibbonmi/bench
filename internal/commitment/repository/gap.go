package repository

import (
	"fmt"
	"slices"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/reviewrecord"
)

// UnfoldedDefault refuses the first dispatch of a chunk whose source lacks the default-branch
// tip. The build folds that tip first, so the author's charge pins it as the chunk base.
type UnfoldedDefault struct{ Branch, Tip, Assignment string }

func (unfolded UnfoldedDefault) Error() string {
	return fmt.Sprintf("the default branch tip %s is not an ancestor of the source", unfolded.Tip)
}

// OpenGap opens the gap of the build that the store root's assignment binds, after a green
// checkpoint of chunk in plan. The plan's last chunk and an unbound root open no gap.
func (store Store) OpenGap(plan reviewrecord.Plan, chunk string) error {
	order := chunkOrder(plan)
	position := slices.Index(order, chunk)
	if position < 0 {
		return fmt.Errorf("the plan holds no chunk %s", chunk)
	}
	if position == len(order)-1 {
		return nil
	}
	return store.decideGap(func(intent.ChunkGap, bool, intent.Assignment) (*intent.ChunkGap, error) {
		return &intent.ChunkGap{Chunk: chunk, Open: true}, nil
	})
}

// CloseGap closes the gap at the first dispatch of the chunk that holds ticket in plan. With
// no gap record, or a record of an earlier chunk, the dispatch is that chunk's first. The
// first dispatch refuses while the source lacks the default-branch tip. A dispatch to the
// recorded chunk or an earlier one keeps the record, so a repair keeps the build's state.
func (store Store) CloseGap(plan reviewrecord.Plan, ticket string) error {
	position := slices.IndexFunc(plan.Chunks, func(chunk reviewrecord.PlannedChunk) bool { return slices.Contains(chunk.Tickets, ticket) })
	if position < 0 {
		return fmt.Errorf("the plan holds no chunk of ticket %s", ticket)
	}
	order := chunkOrder(plan)
	return store.decideGap(func(gap intent.ChunkGap, recorded bool, owner intent.Assignment) (*intent.ChunkGap, error) {
		if recorded && slices.Index(order, gap.Chunk) >= position {
			return nil, nil
		}
		if err := store.holdsDefault(owner); err != nil {
			return nil, err
		}
		return &intent.ChunkGap{Chunk: order[position], Last: position == len(order)-1}, nil
	})
}

// chunkOrder is the chunk ids of plan in plan order.
func chunkOrder(plan reviewrecord.Plan) []string {
	order := make([]string, len(plan.Chunks))
	for i, chunk := range plan.Chunks {
		order[i] = chunk.ID
	}
	return order
}

// holdsDefault refuses a store root whose HEAD lacks the default-branch tip. The refusal
// names owner, the assignment whose worktree folds that tip.
func (store Store) holdsDefault(owner intent.Assignment) error {
	branch, tip, err := store.defaultRevision()
	if err != nil {
		return err
	}
	head, err := store.checkoutHead()
	if err != nil {
		return err
	}
	holds, err := git.IsAncestor(store.Root, tip, head)
	if err != nil {
		return err
	}
	if !holds {
		return UnfoldedDefault{Branch: branch, Tip: tip, Assignment: owner.ID}
	}
	return nil
}

// decideGap writes the gap record that decide answers for the open build that binds the
// active assignment of the store root. decide reads the build's current record, if any. An
// unbound root and a nil answer write nothing, and an unchanged record keeps the ledger bytes.
func (store Store) decideGap(decide func(gap intent.ChunkGap, recorded bool, owner intent.Assignment) (*intent.ChunkGap, error)) error {
	return intent.Transact(store.Root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
		owner, owned := activeOwner(ledger, store.Root)
		if !owned {
			return ledger, false, nil
		}
		builds, err := store.openBuilds(ledger)
		if err != nil {
			return ledger, false, err
		}
		for _, build := range builds {
			if _, bound := build.bound[owner.ID]; !bound {
				continue
			}
			state := storedState(ledger)
			gap, recorded := chunkGap(state, build.claim)
			next, err := decide(gap, recorded, owner)
			if err != nil || next == nil {
				return ledger, false, err
			}
			next.Milestone, next.Outcome = build.claim.Milestone, build.claim.Outcome
			state.Gaps = slices.Clone(state.Gaps)
			if i := gapIndex(state, build.claim); i >= 0 {
				state.Gaps[i] = *next
			} else {
				state.Gaps = append(state.Gaps, *next)
			}
			return withRuntime(ledger, state)
		}
		return ledger, false, nil
	}, nil)
}
