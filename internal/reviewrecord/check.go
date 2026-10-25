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
// reads each obligation's owner from its own frozen plan. fresh reports whether an
// entry is current for the source of the inventory. at holds the chunk ID of a chunk
// inventory, or no chunk and the requested source of the final inventory; each
// refusal repairs through the verification form there.
func checkVerification(items []Verification, requirements []Requirement, plan Plan, record Record, fresh func(Verification) (bool, error), at Repair) error {
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
		ok, err := fresh(*current)
		if err != nil {
			return fmt.Errorf("%s: verification %s: %w", scope, requirement.ID, err)
		}
		if !ok || current.Role != role || current.Performer != performer || current.Command != requirement.Command || current.State != "completed" || current.Outcome != "pass" || current.ExitCode == nil || *current.ExitCode != 0 {
			return refused("verification %s is incomplete, failed, or stale; execute and retain %s", requirement.ID, requirement.Command)
		}
		if requirement.Probe != "" {
			bit, restored := false, false
			if probe := current.Probe; probe != nil && probe.Mutation == requirement.Probe {
				bit, restored = ProbeBites(*probe)
			}
			if !bit {
				return refused("verification %s requires a successful diagnostic probe %s; run and retain it", requirement.ID, requirement.Probe)
			}
			if !restored {
				return refused("verification %s probe restore failed; restore and verify the source", requirement.ID)
			}
		}
	}
	return nil
}

// EntryCurrent is the one predicate of whether a chunk verification entry is current
// for chunk; spec names the record that the source digest excludes. An entry with a
// source commit is current when that commit strictly descends from the chunk base, is
// an ancestor of the chunk tip, and has the entry digest as the source digest of its
// tree. An entry with no source commit is current when its digest is the chunk digest.
// The error reports a failed query.
func EntryCurrent(root, spec string, chunk Chunk, entry Verification) (bool, error) {
	if entry.SourceCommit == "" {
		return entry.SourceDigest == chunk.SourceDigest, nil
	}
	fault, err := chunkLine(root, chunk, entry.SourceCommit)
	if err != nil || fault != onLine {
		return false, err
	}
	digest, err := commitDigest(root, spec, entry.SourceCommit)
	return err == nil && digest == entry.SourceDigest, err
}

// lineFault names the first condition under which a commit is off the line of a chunk.
type lineFault int

const (
	onLine  lineFault = iota
	atBase            // the commit is the chunk base
	offBase           // the commit does not descend from the chunk base
	offTip            // the commit is not an ancestor of the chunk tip
)

// chunkLine answers onLine, or the first fault of commit against chunk in the order above.
func chunkLine(root string, chunk Chunk, commit string) (lineFault, error) {
	// The ancestry query is reflexive, so it accepts a commit equal to the chunk base.
	if commit == chunk.Base {
		return atBase, nil
	}
	if descends, err := isAncestor(root, chunk.Base, commit); err != nil || !descends {
		return offBase, err
	}
	if precedes, err := isAncestor(root, commit, chunk.Tip); err != nil || !precedes {
		return offTip, err
	}
	return onLine, nil
}

// commitDigest is the source digest of the tree of commit; spec names the record that the
// digest excludes.
func commitDigest(root, spec, commit string) (string, error) {
	tree, err := sourceTree(root, commit)
	if err != nil {
		return "", err
	}
	return SourceDigest(root, tree, spec)
}

// ProbeBites is the one probe-bite predicate of the planned-probe check, the confirmation
// check, and the reused-test premise. A probe bit when its outcome is bit and its focused
// exit is not zero, and it restored when its restore passed.
func ProbeBites(probe Probe) (bit, restored bool) {
	return probe.Outcome == "bit" && probe.ExitCode != 0, probe.Restore == "pass"
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
	atDigest := func(entry Verification) (bool, error) { return entry.SourceDigest == source, nil }
	return checkVerification(completion.Verification, finalRequirements(plan, record), plan, record, atDigest, Repair{Source: tip})
}

// TestCommand is the command whose requirements run test selections.
const TestCommand = "bench test"

// PackageTestCommand is the exact command that runs every test of one package.
func PackageTestCommand(pkg string) string {
	return TestCommand + " --package " + pkg
}

// MergedTestRequirement is the final requirement that a fold derives for one merged test
// package.
func MergedTestRequirement(pkg string) Requirement {
	return Requirement{ID: "merged-test:" + pkg, Command: PackageTestCommand(pkg)}
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
