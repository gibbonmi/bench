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
	if run, found, err := s.load(slug); err != nil {
		return Status{}, err
	} else if found {
		if _, err := s.preconditions(mutationStart, slug, resolved, &run, "", ""); err != nil {
			return Status{}, err
		}
		return run.status(), nil
	}
	subject, err := s.preconditions(mutationStart, slug, resolved, nil, "", "")
	if err != nil {
		return Status{}, err
	}
	run := record{
		Version: 1, Slug: slug, Spec: resolved, SpecTip: subject.specTip, Run: digest(resolved), Branch: subject.branch, Base: subject.tip,
		Candidate: "refs/bench/specbuild/candidate/" + digest(resolved), CandidateTip: subject.tip,
		Assignments: map[string]assignment{},
	}
	absent, err := refAbsent(s.root, run.Candidate)
	if err != nil {
		return Status{}, err
	}
	if !absent {
		return Status{}, errors.New("spec build candidate identity already exists")
	}
	if err := s.gate.Bootstrap(ctx, s.root, subject.branch, subject.tip); err != nil {
		return Status{}, fmt.Errorf("no exact green evidence: run bench gate, then retry start: %w", err)
	}
	if err := updateRef(s.root, run.Candidate, subject.tip, zeroObjectID); err != nil {
		return Status{}, fmt.Errorf("create candidate identity: %w", err)
	}
	if err := s.save(run); err != nil {
		return Status{}, err
	}
	return run.status(), nil
}
