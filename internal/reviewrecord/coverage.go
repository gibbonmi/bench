package reviewrecord

import (
	"errors"
	"fmt"
	"maps"
	"reflect"
	"slices"

	"github.com/gibbonmi/bench/internal/commentgap"
	benchgit "github.com/gibbonmi/bench/internal/git"
)

// CheckSource compares retained occurrences with the requested immutable source.
func CheckSource(root, tree, tip string, record Record, chunkID string, complete bool) error {
	return checkSource(root, tree, tip, record, chunkID, complete, false)
}

func checkSource(root, tree, tip string, record Record, chunkID string, complete, verification bool) error {
	current, err := ReadPlan(root, tree, record.Spec)
	if err != nil {
		return err
	}
	if err := currentPlan(record, current, tip); err != nil {
		return err
	}
	digest, err := SourceDigest(root, tree, record.Spec)
	if err != nil {
		return err
	}
	if !objectID.MatchString(tip) {
		return errors.New("invalid requested source tip")
	}
	// No record form writes an unplanned chunk, so its refusal carries no repair.
	if !complete && findChunk(current, chunkID) == nil {
		return fmt.Errorf("chunk %s is not a chunk of the plan at %s", chunkID, tip)
	}
	if len(record.Chunks) == 0 {
		return Repair{Form: FormChunk, Chunk: chunkID, Err: errors.New("missing chunk evidence; retain the three native review results")}
	}
	matched := false
	covered := map[string]bool{}
	var previous *Chunk
	for i := range record.Chunks {
		chunk := &record.Chunks[i]
		if !benchgit.OK("-C", root, "merge-base", "--is-ancestor", chunk.Base, chunk.Tip) || !benchgit.OK("-C", root, "merge-base", "--is-ancestor", chunk.Tip, tip) {
			return fmt.Errorf("chunk %s: stale unrelated base or tip; review the requested source", chunk.ID)
		}
		chunkTree, err := benchgit.Output("-C", root, "rev-parse", "--verify", chunk.Tip+"^{tree}")
		if err != nil {
			return err
		}
		examined, err := SourceDigest(root, chunkTree, record.Spec)
		if err != nil {
			return err
		}
		if examined != chunk.SourceDigest {
			return fmt.Errorf("chunk %s: stale source digest", chunk.ID)
		}
		plan, err := ReadPlan(root, chunkTree, record.Spec)
		if err != nil {
			return err
		}
		if plan.Digest != chunk.PlanDigest {
			return fmt.Errorf("chunk %s: stale plan identity", chunk.ID)
		}
		planned := findChunk(plan, chunk.ID)
		if planned == nil || !sameSet(chunk.AcceptanceRows, planned.Rows) {
			return fmt.Errorf("chunk %s: stale acceptance rows", chunk.ID)
		}
		if verification {
			if err := checkVerification(chunk.Verification, planned.Verification, plan, record, chunk.SourceDigest, Repair{Chunk: chunk.ID}); err != nil {
				return err
			}
		}
		if previous != nil {
			baseTree, err := benchgit.Output("-C", root, "rev-parse", "--verify", chunk.Base+"^{tree}")
			if err != nil {
				return err
			}
			baseDigest, err := SourceDigest(root, baseTree, record.Spec)
			if err != nil {
				return err
			}
			chained := baseDigest == previous.SourceDigest && benchgit.OK("-C", root, "merge-base", "--is-ancestor", previous.Tip, chunk.Base)
			chained = chained || foldCoversBase(root, record, *previous, chunk.Base, baseDigest)
			if !chained {
				return Repair{Form: FormFold, Chunk: chunk.ID, Err: fmt.Errorf("chunk %s: stale review chain gap: expected base %s, the chunk %s tip, got %s; plan commits land before the ticket merge, a default-branch merge lands before the first chunk or as a recorded fold, and only record commits and comment-only corrections follow a chunk tip; record the fold or review the uncovered delta", chunk.ID, previous.Tip, previous.ID, chunk.Base)}
			}
		}
		// Verification ownership reads the frozen plan above, so a historical
		// occurrence keeps its own author. Review exclusion reads the current
		// plan instead: a session that reviewed an early chunk and later became
		// an author must not keep that review.
		reviewerCount, err := current.reviewCardinality()
		if err != nil {
			return err
		}
		excluded := reviewExclusions(current, record)
		if err := checkReviews(chunk.reviewed(), excluded, reviewerCount); err != nil {
			return err
		}
		// A fold after the chunk whose resolution is not empty takes the review rule of a chunk.
		for i := range record.Folds {
			if fold := &record.Folds[i]; fold.After == chunk.ID && len(fold.Resolution) != 0 {
				if err := checkReviews(fold.reviewed(), excluded, reviewerCount); err != nil {
					return err
				}
			}
		}
		ids, err := mappedIDs(record, plan.Digest, current.Digest, chunk.ID)
		if err != nil {
			return err
		}
		for _, id := range ids {
			if findChunk(current, id) == nil {
				return fmt.Errorf("chunk %s: invalid plan mapping target %s", chunk.ID, id)
			}
			covered[id] = true
			for _, predecessor := range current.Chunks {
				if !covered[predecessor.ID] {
					return Repair{Form: FormChunk, Chunk: predecessor.ID, Err: fmt.Errorf("missing planned chunk %s before %s", predecessor.ID, id)}
				}
				if predecessor.ID == id {
					break
				}
			}
		}
		previous = chunk
		if !complete && slices.Contains(ids, chunkID) {
			matched = true
			break
		}
	}
	if !complete && !matched {
		return Repair{Form: FormChunk, Chunk: chunkID, Err: fmt.Errorf("missing chunk %s; record its completed evidence", chunkID)}
	}
	coveredDigest, coveredTip := coveredAfter(root, record, *previous)
	if coveredDigest != digest && commentgap.Prove(root, coveredDigest, digest) != nil {
		// The covered tip and the requested tip bound the delta that no review or fold
		// covers. When they are one commit, the delta is the uncommitted work past that commit.
		uncovered := coveredTip + ".." + tip
		if coveredTip == tip {
			uncovered = "the uncommitted work past " + tip
		}
		err := fmt.Errorf("chunk %s: stale reviewed source: no chunk review covers %s; cover the later source or repair delta", previous.ID, uncovered)
		if fold := uncoveredFold(root, coveredTip, tip); fold != "" {
			return Repair{Form: FormFold, Source: fold, Err: err}
		}
		return err
	}
	if complete {
		for _, planned := range current.Chunks {
			if !covered[planned.ID] {
				return Repair{Form: FormChunk, Chunk: planned.ID, Err: fmt.Errorf("missing planned chunk %s", planned.ID)}
			}
		}
	}
	if complete && verification {
		return checkCompletion(record, current, digest, tip)
	}
	return nil
}

// Chunk returns the planned chunk id, or nil when the plan names no such chunk.
func (p Plan) Chunk(id string) *PlannedChunk { return findChunk(p, id) }

func findChunk(plan Plan, id string) *PlannedChunk {
	for i := range plan.Chunks {
		if plan.Chunks[i].ID == id {
			return &plan.Chunks[i]
		}
	}
	return nil
}

func mappedIDs(record Record, from, to, id string) ([]string, error) {
	ids := []string{id}
	seen := map[string]bool{}
	for from != to {
		if seen[from] {
			return nil, errors.New("invalid plan amendment cycle")
		}
		seen[from] = true
		var change *Amendment
		for i := range record.Amendments {
			if record.Amendments[i].From == from {
				if change != nil {
					return nil, errors.New("invalid ambiguous plan amendment")
				}
				change = &record.Amendments[i]
			}
		}
		if change == nil {
			return nil, errors.New("stale plan amendment; retain an explicit old-to-new chunk ID mapping and review its delta")
		}
		next := []string{}
		for _, old := range ids {
			mapped := change.ChunkIDs[old]
			if len(mapped) == 0 {
				return nil, fmt.Errorf("stale plan amendment: missing chunk mapping %s", old)
			}
			next = append(next, mapped...)
		}
		from, ids = change.To, next
	}
	return ids, nil
}

// RecordAmendment records the change from the record plan digest to the plan digest at
// source, a full commit ID, and moves the record plan digest to the new digest. The keys
// are the recorded chunk IDs under the old digest: mappedIDs maps each recorded chunk
// through the earlier amendments, so the checkpoint and the amendment read one rule. A
// key maps to itself unless changes names it, and every target must be a chunk of the new
// plan. Every recorded chunk must then map through the amended chain to the new digest,
// or the write refuses. The answer is the amendment as written.
func RecordAmendment(root, spec, source string, changes map[string][]string) (Amendment, error) {
	_, plan, err := atSource(root, spec, source)
	if err != nil {
		return Amendment{}, err
	}
	var entry Amendment
	err = write(root, spec, absent(""), func(record *Record) error {
		if record.PlanDigest == plan.Digest {
			return fmt.Errorf("the plan at %s is unchanged from the record plan digest %s", source, plan.Digest)
		}
		keys := []string{}
		for _, chunk := range record.Chunks {
			ids, err := mappedIDs(*record, chunk.PlanDigest, record.PlanDigest, chunk.ID)
			if err != nil {
				return fmt.Errorf("chunk %s: %w", chunk.ID, err)
			}
			for _, id := range ids {
				if !slices.Contains(keys, id) {
					keys = append(keys, id)
				}
			}
		}
		for _, old := range slices.Sorted(maps.Keys(changes)) {
			if !slices.Contains(keys, old) {
				return fmt.Errorf("chunk %s is not a recorded chunk under the plan digest %s, so no map can name it", old, record.PlanDigest)
			}
		}
		entry = Amendment{From: record.PlanDigest, To: plan.Digest, ChunkIDs: map[string][]string{}}
		for _, key := range keys {
			targets, mapped := changes[key]
			if !mapped {
				targets = []string{key}
			}
			for _, id := range targets {
				if findChunk(plan, id) != nil {
					continue
				}
				if !mapped {
					return fmt.Errorf("recorded chunk %s is not a chunk of the plan at %s; map it to its new chunks", key, source)
				}
				return fmt.Errorf("map target chunk %s is not a chunk of the plan at %s", id, source)
			}
			entry.ChunkIDs[key] = append([]string{}, targets...)
		}
		record.Amendments = append(record.Amendments, entry)
		record.PlanDigest = plan.Digest
		for _, chunk := range record.Chunks {
			if _, err := mappedIDs(*record, chunk.PlanDigest, record.PlanDigest, chunk.ID); err != nil {
				return fmt.Errorf("chunk %s: %w", chunk.ID, err)
			}
		}
		return nil
	})
	return entry, err
}

func sameSet(a, b []string) bool {
	set := func(values []string) map[string]bool {
		result := map[string]bool{}
		for _, value := range values {
			result[value] = true
		}
		return result
	}
	return len(a) == len(b) && reflect.DeepEqual(set(a), set(b))
}
