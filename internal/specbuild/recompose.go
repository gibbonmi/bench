package specbuild

import (
	"context"
	"errors"
	"fmt"
)

func (s *Service) recomposePromotion(ctx context.Context, run *record, subject buildSubject) error {
	old := run.CandidateTip
	candidateTip := subject.tip
	// An empty run has no checkpoint patch, so it rebases directly onto the bootstrapped subject tip.
	if old != run.Base {
		patch, err := s.checkpointPatch(ctx, run.Base, old)
		if err != nil {
			return err
		}
		tree, err := s.replayCheckpoint(ctx, subject.tip, run.Base, old, patch)
		if err != nil {
			return err
		}
		candidateTip, err = s.gitOutput(ctx, "commit-tree", tree, "-p", subject.tip, "-m", "bench recompose run="+run.Run+" candidate="+old)
		if err != nil {
			return err
		}
	}
	// A replay must prove conflict-free before the independently valid working tip
	// becomes project-green, while bootstrap must precede candidate and run mutation.
	if err := s.gate.Bootstrap(ctx, s.root, subject.branch, subject.tip, greenMarker(s.root, subject.branch)); err != nil {
		return fmt.Errorf("no exact green evidence: run bench gate --fresh, then retry promote: %w", err)
	}
	return s.finishRecomposition(run, old, subject.tip, candidateTip)
}

// preconditionsAdvancingEmptyRun runs op's preconditions and, where a run has nothing for
// a recomposition to replay, advances it onto the working tip instead of refusing. Only
// checkpoint and start compose it: promote turns the same refusal into its gated
// recomposition, review and assign refuse on it, and a terminal run restarts on it, so
// widening the fast-forward past this pair would silently retire all four behaviors.
func (s *Service) preconditionsAdvancingEmptyRun(op mutation, slug, specPath string, run *record, assignmentID, evidence string) (buildSubject, error) {
	subject, err := s.preconditions(op, slug, specPath, run, assignmentID, evidence)
	if !errors.Is(err, errRecompose) || run == nil || run.Terminal || !emptyRun(*run) {
		return subject, err
	}
	if err := s.fastForwardEmptyRun(run, subject.tip); err != nil {
		return buildSubject{}, err
	}
	return subject, nil
}

// emptyRun reports whether run holds no work a recomposition would have to replay: no
// assignment has checkpointed, and the candidate still sits on the recorded base.
func emptyRun(run record) bool {
	if run.CandidateTip != run.Base {
		return false
	}
	for _, assigned := range run.Assignments {
		if assigned.Checkpoint != "" || assigned.CheckpointRef != "" {
			return false
		}
	}
	return true
}

// fastForwardEmptyRun moves run onto tip, provisionally: the durable candidate ref is the
// run's identity, so it moves by compare-and-swap from the tip the record claims and a ref
// some other writer already advanced refuses here rather than being overwritten. No gate
// runs — promote remains the lifecycle's only green transition.
func (s *Service) fastForwardEmptyRun(run *record, tip string) error {
	if err := updateRef(s.root, run.Candidate, tip, run.CandidateTip); err != nil {
		return err
	}
	run.Base, run.CandidateTip = tip, tip
	return s.save(*run)
}

func (s *Service) finishRecomposition(run *record, old, base, candidateTip string) error {
	if err := updateRef(s.root, run.Candidate, candidateTip, old); err != nil {
		return err
	}
	for key, assigned := range run.Assignments {
		if assigned.Integrated != "" {
			assigned.Integrated = candidateTip
			run.Assignments[key] = assigned
		}
	}
	run.Base, run.CandidateTip, run.Review = base, candidateTip, nil
	run.PromotionTree, run.PromotionCommit, run.PromotionEvidence, run.PromotionDisposition = "", "", "", ""
	delete(run.Operations, operationID("promote", old))
	return s.save(*run)
}
