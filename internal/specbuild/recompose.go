package specbuild

import (
	"context"
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
	if err := s.gate.Bootstrap(ctx, s.root, subject.branch, subject.tip, run.Base); err != nil {
		return fmt.Errorf("no exact green evidence: run bench gate --fresh, then retry promote: %w", err)
	}
	return s.finishRecomposition(run, old, subject.tip, candidateTip)
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
