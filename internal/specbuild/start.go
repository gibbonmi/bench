package specbuild

import (
	"context"
	"errors"
	"fmt"
)

// Start creates or resumes the run for slug.
func (s *Service) Start(ctx context.Context, slug string) (Status, error) {
	resolved, err := s.resolve(slug)
	if err != nil {
		return Status{}, err
	}
	release, err := s.lock(slug)
	if err != nil {
		return Status{}, err
	}
	defer release()
	branch, tip, err := workingSubject(s.root)
	if err != nil {
		return Status{}, err
	}
	if run, found, err := s.load(slug); err != nil {
		return Status{}, err
	} else if found {
		if run.Branch != branch || run.Base != tip || !refAt(s.root, run.Candidate, run.CandidateTip) {
			return Status{}, errors.New("spec build start conflicts with the recorded working subject")
		}
		return run.status(), nil
	}
	run := record{
		Version: 1, Slug: slug, Spec: resolved, Run: digest(resolved), Branch: branch, Base: tip,
		Candidate: "refs/bench/specbuild/candidate/" + digest(resolved), CandidateTip: tip,
		Assignments: map[string]assignment{},
	}
	absent, err := refAbsent(s.root, run.Candidate)
	if err != nil {
		return Status{}, err
	}
	if !absent {
		return Status{}, errors.New("spec build candidate identity already exists")
	}
	if s.gate == nil {
		return Status{}, errors.New("spec build start requires a gate owner")
	}
	if err := s.gate.Bootstrap(ctx, s.root, branch, tip); err != nil {
		return Status{}, fmt.Errorf("no exact green evidence: run bench gate, then retry start: %w", err)
	}
	if err := updateRef(s.root, run.Candidate, tip, zeroObjectID); err != nil {
		return Status{}, fmt.Errorf("create candidate identity: %w", err)
	}
	if err := s.save(run); err != nil {
		return Status{}, err
	}
	return run.status(), nil
}
