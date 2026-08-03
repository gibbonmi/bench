package specbuild

import (
	"context"
	"fmt"
)

func (s *Service) recomposePromotion(ctx context.Context, run *record, subject buildSubject) error {
	old := run.CandidateTip
	patch, err := s.checkpointPatch(ctx, run.Base, old)
	if err != nil {
		return err
	}
	tree, err := s.replayCheckpoint(ctx, subject.tip, run.Base, old, patch)
	if err != nil {
		return err
	}
	commit, err := s.gitOutput(ctx, "commit-tree", tree, "-p", subject.tip, "-m", "bench recompose run="+run.Run+" candidate="+old)
	if err != nil {
		return err
	}
	// Replay must prove conflict-free before the independently valid working tip
	// becomes project-green, while bootstrap must precede candidate and run mutation.
	if err := s.gate.Bootstrap(ctx, s.root, subject.branch, subject.tip, run.Base); err != nil {
		return fmt.Errorf("no exact green evidence: run bench gate --fresh, then retry promote: %w", err)
	}
	if err := updateRef(s.root, run.Candidate, commit, old); err != nil {
		return err
	}
	for key, assigned := range run.Assignments {
		if assigned.Integrated != "" {
			assigned.Integrated = commit
			run.Assignments[key] = assigned
		}
	}
	run.Base, run.CandidateTip, run.Review = subject.tip, commit, nil
	run.PromotionTree, run.PromotionCommit, run.PromotionEvidence, run.PromotionDisposition = "", "", "", ""
	delete(run.Operations, operationID("promote", old))
	return s.save(*run)
}
