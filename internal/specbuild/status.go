package specbuild

import "sort"

// Status returns the durable compact projection for slug.
func (s *Service) Status(slug string) (Status, error) {
	if _, err := s.resolve(slug); err != nil {
		return Status{}, err
	}
	run, found, err := s.load(slug)
	if err != nil {
		return Status{}, err
	}
	if !found {
		return Status{Slug: slug, State: "empty", Next: "bench spec build start " + slug}, nil
	}
	return run.status(), nil
}

// FullStatus returns the compact projection together with retained provenance.
func (s *Service) FullStatus(slug string) (FullStatus, error) {
	compact, err := s.Status(slug)
	if err != nil || compact.State == "empty" {
		return FullStatus{Status: compact}, err
	}
	run, _, err := s.load(slug)
	if err != nil {
		return FullStatus{}, err
	}
	full := FullStatus{Status: compact, Assignments: make([]RetainedAssignment, 0, len(run.Assignments))}
	for _, assigned := range run.Assignments {
		cleanup := "active"
		if assigned.CleanupPending {
			cleanup = "pending"
		} else if assigned.Released {
			cleanup = "released"
		}
		full.Assignments = append(full.Assignments, RetainedAssignment{ID: assigned.ID, Ticket: assigned.Ticket, TicketDigest: assigned.TicketDigest, Base: assigned.Base, Checkpoint: assigned.Checkpoint, CheckpointRef: assigned.CheckpointRef, Integrated: assigned.Integrated, ReceiptDigest: assigned.ReceiptDigest, Cleanup: cleanup})
	}
	sort.Slice(full.Assignments, func(i, j int) bool { return full.Assignments[i].ID < full.Assignments[j].ID })
	if run.Review != nil {
		full.Review = &RetainedReview{Candidate: run.Review.Candidate, Digest: run.Review.Digest, Axes: publicAxes(run.Review.Axes)}
	}
	return full, nil
}

func publicAxes(axes []reviewAxis) []ReviewAxis {
	result := make([]ReviewAxis, len(axes))
	for index, axis := range axes {
		findings := make([]ReviewFinding, len(axis.Findings))
		for findingIndex, finding := range axis.Findings {
			findings[findingIndex] = ReviewFinding{ID: finding.ID, Disposition: finding.Disposition}
		}
		result[index] = ReviewAxis{Axis: axis.Axis, Findings: findings}
	}
	return result
}
