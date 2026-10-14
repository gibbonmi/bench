package reviewrecord

import (
	"fmt"
	"maps"
	"slices"
)

// createRecord starts a new record of spec from plan, the plan at tip. Only a version 2
// plan starts a new record.
func createRecord(spec string, plan Plan, tip string) (Record, error) {
	if plan.Version != 2 {
		return Record{}, fmt.Errorf("a new record needs a version 2 completion plan, and the plan at %s has version %d", tip, plan.Version)
	}
	return Record{Version: 2, Spec: spec, PlanDigest: plan.Digest, Completion: Completion{State: "pending", Reconciliation: map[string]string{}, Verification: []Verification{}}}, nil
}

// changeAmendment appends the change from the record plan digest to the digest of plan,
// the plan at source, and moves the record plan digest to the new digest. The keys are
// the recorded chunk IDs under the old digest: mappedIDs maps each recorded chunk through
// the earlier amendments, so the checkpoint and the amendment read one rule. A key maps to
// itself unless changes names it, and every target must be a chunk of the new plan. Every
// recorded chunk must then map through the amended chain to the new digest, or the change
// refuses. The answer is the amendment as written.
func changeAmendment(record *Record, plan Plan, source string, changes map[string][]string) (Amendment, error) {
	if record.PlanDigest == plan.Digest {
		return Amendment{}, fmt.Errorf("the plan at %s is unchanged from the record plan digest %s", source, plan.Digest)
	}
	keys := []string{}
	for _, chunk := range record.Chunks {
		ids, err := mappedIDs(*record, chunk.PlanDigest, record.PlanDigest, chunk.ID)
		if err != nil {
			return Amendment{}, fmt.Errorf("chunk %s: %w", chunk.ID, err)
		}
		for _, id := range ids {
			if !slices.Contains(keys, id) {
				keys = append(keys, id)
			}
		}
	}
	for _, old := range slices.Sorted(maps.Keys(changes)) {
		if !slices.Contains(keys, old) {
			return Amendment{}, fmt.Errorf("chunk %s is not a recorded chunk under the plan digest %s, so no map can name it", old, record.PlanDigest)
		}
	}
	entry := Amendment{From: record.PlanDigest, To: plan.Digest, ChunkIDs: map[string][]string{}}
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
				return entry, fmt.Errorf("recorded chunk %s is not a chunk of the plan at %s; map it to its new chunks", key, source)
			}
			return entry, fmt.Errorf("map target chunk %s is not a chunk of the plan at %s", id, source)
		}
		entry.ChunkIDs[key] = append([]string{}, targets...)
	}
	record.Amendments = append(record.Amendments, entry)
	record.PlanDigest = plan.Digest
	for _, chunk := range record.Chunks {
		if _, err := mappedIDs(*record, chunk.PlanDigest, record.PlanDigest, chunk.ID); err != nil {
			return entry, fmt.Errorf("chunk %s: %w", chunk.ID, err)
		}
	}
	return entry, nil
}

// changeChunk writes the entry of planned chunk id for the frozen pair base..tip. Plan is
// the plan at tip, and digest is the source digest at tip. A recorded chunk keeps its
// results, and a new chunk follows the recorded chunks. A record whose plan digest differs
// from the plan at tip needs its amendment first. The answer is the action, added or
// updated, and the entry as written.
func changeChunk(record *Record, plan Plan, digest, id, base, tip string) (string, Chunk, error) {
	if err := currentPlan(*record, plan, tip); err != nil {
		return "", Chunk{}, err
	}
	planned, err := plannedChunk(plan, tip, id)
	if err != nil {
		return "", Chunk{}, err
	}
	action, index := "added", len(record.Chunks)
	for i := range record.Chunks {
		if record.Chunks[i].ID == id {
			index, action = i, "updated"
		}
	}
	if index == len(record.Chunks) {
		record.Chunks = append(record.Chunks, Chunk{ID: id, Verification: []Verification{}, Reviews: []Review{}})
	}
	chunk := &record.Chunks[index]
	chunk.Base, chunk.Tip, chunk.PlanDigest, chunk.SourceDigest = base, tip, plan.Digest, digest
	chunk.AcceptanceRows = append([]string{}, planned.Rows...)
	return action, *chunk, nil
}

// changeVerification appends one completed verification result to the chunk list or to
// the completion list. The plan at the source tree under root supplies the requirement,
// its command, its probe mutation, and the role; the source tree supplies the source
// digest. A final result can also name the merged-test requirement of a fold package. A
// chunk result must name the source of its chunk entry. The answer is the entry as written.
func changeVerification(root, spec string, record *Record, call VerificationCall) (Verification, error) {
	list, source, requirement := &record.Completion.Verification, call.Source, Requirement{}
	var chunk *Chunk
	var err error
	if call.Chunk != "" {
		if chunk, err = entryOf(record, call.Chunk); err != nil {
			return Verification{}, err
		}
		list = &chunk.Verification
		if source == "" {
			source = chunk.Tip
		}
	}
	digest, plan, err := atSource(root, spec, source)
	if err != nil {
		return Verification{}, err
	}
	if chunk != nil && digest != chunk.SourceDigest {
		return Verification{}, Repair{Form: FormChunk, Chunk: chunk.ID, Err: fmt.Errorf("source %s differs from the source of chunk %s: %w", source, chunk.ID, ErrNoChunkEntry)}
	}
	if chunk != nil {
		requirement, err = ChunkRequirement(plan, source, chunk.ID, call.Requirement)
	} else {
		requirement, err = plannedRequirement(finalRequirements(plan, *record), call.Requirement)
	}
	if err != nil {
		return Verification{}, err
	}
	_, role, err := verifier(plan, *record, requirement, chunk == nil)
	if err != nil {
		return Verification{}, err
	}
	if requirement.Probe == "" && call.Probe != nil {
		return Verification{}, fmt.Errorf("requirement %s plans no probe: %w", requirement.ID, ErrProbe)
	}
	if requirement.Probe != "" && call.Probe == nil {
		return Verification{}, fmt.Errorf("requirement %s plans the probe %q: %w", requirement.ID, requirement.Probe, ErrProbe)
	}
	if err := unclaimed(*record, call.Evidence.ID); err != nil {
		return Verification{}, err
	}
	exitCode, outcome := call.ExitCode, "fail"
	if exitCode == 0 {
		outcome = "pass"
	}
	entry := Verification{Evidence: call.Evidence, Requirement: requirement.ID, Command: requirement.Command, ExitCode: &exitCode}
	entry.Role, entry.SourceDigest, entry.State, entry.Outcome = role, digest, "completed", outcome
	if call.Probe != nil {
		probe := *call.Probe
		probe.Mutation, probe.NativeRef = requirement.Probe, call.Evidence.NativeRef
		entry.Probe = &probe
	}
	*list = append(*list, entry)
	return entry, nil
}
