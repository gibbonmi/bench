package reviewrecord

import (
	"errors"
	"fmt"
	"slices"
)

// Check grades retained occurrences against the requested source and obligation.
func Check(root, tree, tip, spec, chunk string, complete bool) error {
	return CheckTrees(root, tree, tree, tip, spec, chunk, complete)
}

// CheckTrees reads evidence from the graded tree and coverage from its proven source.
// The gate owns proof that the composition preserves that source.
func CheckTrees(root, sourceTree, evidenceTree, tip, spec, chunk string, complete bool) error {
	record, err := ReadTree(root, evidenceTree, spec)
	if err != nil {
		return fmt.Errorf("checkpoint %s: %w; retain a valid native result record", spec, err)
	}
	return checkSource(root, sourceTree, tip, record, chunk, complete, true)
}

// checkVerification grades retained occurrences against one requirement
// inventory. The plan resolves each obligation's owed performer and role, so a
// version 1 record keeps its one implementation session and a version 2 record
// reads each obligation's owner from its own frozen plan. at holds the chunk ID of
// a chunk inventory, or no chunk and the requested source of the final inventory;
// each refusal repairs through the verification form there.
func checkVerification(items []Verification, requirements []Requirement, plan Plan, record Record, source string, at Repair) error {
	scope := "completion"
	if at.Chunk != "" {
		scope = "chunk " + at.Chunk
	}
	for _, requirement := range requirements {
		performer, role, err := verifier(plan, record, requirement, at.Chunk == "")
		if err != nil {
			return fmt.Errorf("%s: %w", scope, err)
		}
		refused := func(format string, args ...any) error {
			repair := at
			repair.Form, repair.Requirement, repair.Err = FormVerification, requirement.ID, fmt.Errorf("%s: "+format, append([]any{scope}, args...)...)
			return repair
		}
		var current *Verification
		for i := range items {
			if items[i].Requirement == requirement.ID {
				current = &items[i]
			}
		}
		if current == nil {
			return refused("missing verification %s; execute and retain %s", requirement.ID, requirement.Command)
		}
		if current.Role != role || current.Performer != performer || current.SourceDigest != source || current.Command != requirement.Command || current.State != "completed" || current.Outcome != "pass" || current.ExitCode == nil || *current.ExitCode != 0 {
			return refused("verification %s is incomplete, failed, or stale; execute and retain %s", requirement.ID, requirement.Command)
		}
		if requirement.Probe != "" {
			probe := current.Probe
			if probe == nil || probe.Mutation != requirement.Probe || probe.Outcome != "bit" || probe.ExitCode == 0 {
				return refused("verification %s requires a successful diagnostic probe %s; run and retain it", requirement.ID, requirement.Probe)
			}
			if probe.Restore != "pass" {
				return refused("verification %s probe restore failed; restore and verify the source", requirement.ID)
			}
		}
	}
	return nil
}

// checkCompletion grades the completion entry and the final verification against the
// current plan at source, the digest of tip. The final list includes the merged-test
// requirements of the folds.
func checkCompletion(record Record, plan Plan, source, tip string) error {
	completion := record.Completion
	reconciler := completionReconciler(record, plan)
	if completion.State != "completed" || completion.SourceDigest != source || completion.Performer != reconciler {
		return Repair{Form: FormCompletion, Source: tip, Err: errors.New("completion is incomplete or stale; execute final verification and reconcile acceptance")}
	}
	for _, chunk := range plan.Chunks {
		for _, row := range chunk.Rows {
			if completion.Reconciliation[row] != "covered" {
				return Repair{Form: FormCompletion, Source: tip, Err: fmt.Errorf("completion: missing acceptance reconciliation %s; reconcile every planned row", row)}
			}
		}
	}
	return checkVerification(completion.Verification, finalRequirements(plan, record), plan, record, source, Repair{Source: tip})
}

// MergedTestRequirement is the final requirement that a fold derives for one merged test
// package.
func MergedTestRequirement(pkg string) Requirement {
	return Requirement{ID: "merged-test:" + pkg, Command: "bench test --package " + pkg}
}

// finalRequirements is the final verification list of record under plan: the planned list,
// then the merged-test requirement of each fold package, each package once.
func finalRequirements(plan Plan, record Record) []Requirement {
	requirements := append([]Requirement{}, plan.FinalVerification...)
	for _, fold := range record.Folds {
		for _, pkg := range fold.TestPackages {
			if merged := MergedTestRequirement(pkg); !slices.Contains(requirements, merged) {
				requirements = append(requirements, merged)
			}
		}
	}
	return requirements
}
