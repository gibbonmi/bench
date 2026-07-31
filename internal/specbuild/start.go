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
		if op, pending := s.operation(run, "start", "run"); pending && op.State == "prepared" {
			subject, err := s.subject(resolved)
			if err != nil || subject.branch != run.Branch || subject.tip != run.Base || subject.specTip != run.SpecTip {
				return Status{}, errors.New("spec build working checkout does not match recorded subject")
			}
			return s.finishStart(ctx, subject.branch, subject.tip, false, &run)
		}
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
		Assignments: map[string]assignment{}, Operations: map[string]operation{},
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
	if _, _, err := s.beginOperation(&run, "start", "run", resolved+"\x00"+subject.branch+"\x00"+subject.tip); err != nil {
		return Status{}, err
	}
	if err := s.faultAt("start/bootstrap"); err != nil {
		return run.status(), err
	}
	if err := s.faultAt("start/state"); err != nil {
		return run.status(), err
	}
	return s.finishStart(ctx, subject.branch, subject.tip, true, &run)
}

func (s *Service) finishStart(ctx context.Context, branch, tip string, greenReady bool, run *record) (Status, error) {
	if !refAt(s.root, run.Candidate, tip) {
		absent, err := refAbsent(s.root, run.Candidate)
		if err != nil {
			return Status{}, err
		}
		if !absent {
			return Status{}, errors.New("spec build candidate identity already exists")
		}
		if !greenReady && !refAt(s.root, "refs/bench/green/"+branch, tip) {
			if err := s.gate.Bootstrap(ctx, s.root, branch, tip); err != nil {
				return Status{}, fmt.Errorf("no exact green evidence: run bench gate, then retry start: %w", err)
			}
		}
		if err := updateRef(s.root, run.Candidate, tip, zeroObjectID); err != nil {
			return Status{}, fmt.Errorf("create candidate identity: %w", err)
		}
		if err := s.faultAt("start/candidate-ref"); err != nil {
			return run.status(), err
		}
	}
	if err := s.recordOperation(run, "start", "run", run.Candidate, true); err != nil {
		return Status{}, err
	}
	return run.status(), nil
}
