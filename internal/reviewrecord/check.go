package reviewrecord

import (
	"errors"
	"fmt"
)

// Check grades retained occurrences against the requested source and obligation.
func Check(root, tree, tip, spec, chunk string, complete bool) error {
	record, err := ReadTree(root, tree, spec)
	if err != nil {
		return fmt.Errorf("checkpoint %s: %w; retain a valid native result record", spec, err)
	}
	return checkSource(root, tree, tip, record, chunk, complete, true)
}

func checkVerification(items []Verification, requirements []Requirement, performer, source, scope string) error {
	for _, requirement := range requirements {
		var current *Verification
		for i := range items {
			if items[i].Requirement == requirement.ID {
				current = &items[i]
			}
		}
		if current == nil {
			return fmt.Errorf("%s: missing verification %s; execute and retain %s", scope, requirement.ID, requirement.Command)
		}
		if current.Role != "author-verification" || current.Performer != performer || current.SourceDigest != source || current.Command != requirement.Command || current.State != "completed" || current.Outcome != "pass" || current.ExitCode == nil || *current.ExitCode != 0 {
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
	if completion.State != "completed" || completion.SourceDigest != source || completion.Performer != record.ImplementationSession {
		return errors.New("completion is incomplete or stale; execute final verification and reconcile acceptance")
	}
	for _, chunk := range plan.Chunks {
		for _, row := range chunk.Rows {
			if completion.Reconciliation[row] != "covered" {
				return fmt.Errorf("completion: missing acceptance reconciliation %s; reconcile every planned row", row)
			}
		}
	}
	return checkVerification(completion.Verification, plan.FinalVerification, record.ImplementationSession, source, "completion")
}
