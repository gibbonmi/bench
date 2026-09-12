package reviewrecord

import (
	"errors"
	"fmt"
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
// reads each obligation's owner from its own frozen plan.
func checkVerification(items []Verification, requirements []Requirement, plan Plan, record Record, source, scope string, final bool) error {
	for _, requirement := range requirements {
		performer, role, err := verifier(plan, record, requirement, final)
		if err != nil {
			return fmt.Errorf("%s: %w", scope, err)
		}
		var current *Verification
		for i := range items {
			if items[i].Requirement == requirement.ID {
				current = &items[i]
			}
		}
		if current == nil {
			return fmt.Errorf("%s: missing verification %s; execute and retain %s", scope, requirement.ID, requirement.Command)
		}
		if current.Role != role || current.Performer != performer || current.SourceDigest != source || current.Command != requirement.Command || current.State != "completed" || current.Outcome != "pass" || current.ExitCode == nil || *current.ExitCode != 0 {
			return fmt.Errorf("%s: verification %s is incomplete, failed, or stale; execute and retain %s", scope, requirement.ID, requirement.Command)
		}
		if requirement.Probe != "" {
			probe := current.Probe
			if probe == nil || probe.Mutation != requirement.Probe || probe.Outcome != "bit" || probe.ExitCode == 0 {
				return fmt.Errorf("%s: verification %s requires a successful diagnostic probe %s; run and retain it", scope, requirement.ID, requirement.Probe)
			}
			if probe.Restore != "pass" {
				return fmt.Errorf("%s: verification %s probe restore failed; restore and verify the source", scope, requirement.ID)
			}
		}
	}
	return nil
}

func checkCompletion(record Record, plan Plan, source string) error {
	completion := record.Completion
	reconciler := record.ImplementationSession
	if plan.Delegated() {
		reconciler = plan.Execution.OrchestratorSession
	}
	if completion.State != "completed" || completion.SourceDigest != source || completion.Performer != reconciler {
		return errors.New("completion is incomplete or stale; execute final verification and reconcile acceptance")
	}
	for _, chunk := range plan.Chunks {
		for _, row := range chunk.Rows {
			if completion.Reconciliation[row] != "covered" {
				return fmt.Errorf("completion: missing acceptance reconciliation %s; reconcile every planned row", row)
			}
		}
	}
	return checkVerification(completion.Verification, plan.FinalVerification, plan, record, source, "completion", true)
}
