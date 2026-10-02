package reviewrecord

import (
	"errors"
	"fmt"
	"maps"
	"reflect"
	"slices"

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
	if record.PlanDigest != current.Digest {
		return errors.New("stale plan digest; review and record the plan delta")
	}
	digest, err := SourceDigest(root, tree, record.Spec)
	if err != nil {
		return err
	}
	if !objectID.MatchString(tip) {
		return errors.New("invalid requested source tip")
	}
	if len(record.Chunks) == 0 {
		return errors.New("missing chunk evidence; retain the three native review results")
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
			if err := checkVerification(chunk.Verification, planned.Verification, plan, record, chunk.SourceDigest, "chunk "+chunk.ID, false); err != nil {
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
			if baseDigest != previous.SourceDigest || !benchgit.OK("-C", root, "merge-base", "--is-ancestor", previous.Tip, chunk.Base) {
				return fmt.Errorf("chunk %s: stale review chain gap: expected base %s, the chunk %s tip, got %s; plan commits land before the ticket merge, a default-branch merge lands only before the first chunk, and only record commits follow a chunk tip; review the uncovered delta", chunk.ID, previous.Tip, previous.ID, chunk.Base)
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
		if err := checkReviews(*chunk, reviewExclusions(current, record), reviewerCount); err != nil {
			return err
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
					return fmt.Errorf("missing planned chunk %s before %s", predecessor.ID, id)
				}
				if predecessor.ID == id {
					break
				}
			}
		}
		previous = chunk
		if !complete && contains(ids, chunkID) {
			matched = true
			break
		}
	}
	if !complete && !matched {
		return fmt.Errorf("missing chunk %s; record its completed evidence", chunkID)
	}
	if previous.SourceDigest != digest {
		// The reviewed tip and the requested tip bound the delta no review covers. When they
		// are one commit, the delta is the uncommitted work past that commit.
		uncovered := previous.Tip + ".." + tip
		if previous.Tip == tip {
			uncovered = "the uncommitted work past " + tip
		}
		return fmt.Errorf("chunk %s: stale reviewed source: no chunk review covers %s; cover the later source or repair delta", previous.ID, uncovered)
	}
	if complete {
		for _, planned := range current.Chunks {
			if !covered[planned.ID] {
				return fmt.Errorf("missing planned chunk %s", planned.ID)
			}
		}
	}
	if complete && verification {
		return checkCompletion(record, current, digest)
	}
	return nil
}

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
	err = write(root, spec, nil, func(record *Record) error {
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
				if !contains(keys, id) {
					keys = append(keys, id)
				}
			}
		}
		for _, old := range slices.Sorted(maps.Keys(changes)) {
			if !contains(keys, old) {
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
