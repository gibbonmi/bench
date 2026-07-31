// Package specbuild owns the durable lifecycle of a reviewed spec build.
package specbuild

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/jsonfile"
)

// GateOwner validates retained exact green evidence for a working subject.
type GateOwner interface {
	Bootstrap(context.Context, string, string, string) error
}

// WorktreeOwner creates a request-idempotent owned worktree at start.
type WorktreeOwner interface {
	Create(context.Context, string, string, string, string) (OwnedWorktree, error)
}

// ReleaseOwner releases an assignment created by the worktree ownership owner.
type ReleaseOwner interface {
	Release(context.Context, string, string, string) error
}

// OwnedWorktree identifies a worktree created by the existing ownership owner.
type OwnedWorktree struct {
	ID, Path string
}

// Service coordinates spec build transitions from one working checkout.
type Service struct {
	root               string
	gate               GateOwner
	worktrees          WorktreeOwner
	beforeCandidateCAS func()
}

// New constructs a lifecycle service rooted at one working checkout.
func New(root string, gate GateOwner, worktrees WorktreeOwner) *Service {
	return &Service{root: root, gate: gate, worktrees: worktrees}
}

// Status is the compact public projection of one spec build.
type Status struct {
	Slug, State, Subject, Next string
}

// Assignment records the externally useful ownership binding for one ticket.
type Assignment struct {
	ID, Path, Base string
	Rows           []string
	Fence          []string
	Assumptions    []string
}

// FullStatus is the retained evidence projection for one spec build.
type FullStatus struct {
	Status
	Assignments []RetainedAssignment
	Review      *RetainedReview
}

// RetainedAssignment is the durable assignment provenance exposed by full status.
type RetainedAssignment struct {
	ID, Ticket, TicketDigest, Base, Checkpoint, CheckpointRef, Integrated, ReceiptDigest, Cleanup string
}

// RetainedReview is the review evidence exposed without its receipt body.
type RetainedReview struct {
	Candidate, Digest string
	Axes              []ReviewAxis
}

// ReviewAxis is one named semantic review axis and its finding dispositions.
type ReviewAxis struct {
	Axis     string
	Findings []ReviewFinding
}

// ReviewFinding identifies one finding and its reviewer-supplied disposition.
type ReviewFinding struct {
	ID, Disposition string
}

var errInvalidReviewReceipt = errors.New("invalid spec build review receipt")

// Review records one complete semantic review for the exact current candidate.
func (s *Service) Review(_ context.Context, slug, evidence string) (Status, error) {
	if _, err := s.resolve(slug); err != nil {
		return Status{}, err
	}
	release, err := s.lock(slug)
	if err != nil {
		return Status{}, err
	}
	defer release()
	run, found, err := s.load(slug)
	if err != nil || !found {
		if err == nil {
			err = errors.New("spec build run does not exist")
		}
		return Status{}, err
	}
	if _, err := s.preconditions(mutationReview, slug, run.Spec, &run, "", evidence); err != nil {
		return Status{}, err
	}
	receipt, raw, err := readReviewReceipt(evidence)
	if err != nil || receipt.Run != run.Run || receipt.Candidate != run.CandidateTip || !refAt(s.root, run.Candidate, run.CandidateTip) {
		return Status{}, errInvalidReviewReceipt
	}
	digest := digest(string(raw))
	if run.Review != nil {
		if run.Review.Candidate == run.CandidateTip && run.Review.Digest == digest {
			return run.status(), nil
		}
		return Status{}, errors.New("spec build candidate already has review evidence")
	}
	run.Review = &reviewEvidence{Candidate: receipt.Candidate, Axes: receipt.Axes, Digest: digest}
	if err := s.save(run); err != nil {
		return Status{}, err
	}
	return run.status(), nil
}

type reviewReceipt struct {
	Version   int          `json:"version"`
	Run       string       `json:"run"`
	Candidate string       `json:"candidate"`
	Axes      []reviewAxis `json:"axes"`
	Body      string       `json:"body"`
}

func readReviewReceipt(path string) (reviewReceipt, []byte, error) {
	if path == "" || !filepath.IsAbs(path) {
		return reviewReceipt{}, nil, errInvalidReviewReceipt
	}
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() {
		return reviewReceipt{}, nil, errInvalidReviewReceipt
	}
	classified := bounds.Classify(path, bounds.ControlRecordLimit)
	if classified.State != bounds.StateParsed || !strings.HasSuffix(string(classified.Data), "\n") {
		return reviewReceipt{}, nil, errInvalidReviewReceipt
	}
	var receipt reviewReceipt
	if jsonfile.Decode(classified.Data, &receipt) != nil || !validReviewReceipt(receipt) {
		return reviewReceipt{}, nil, errInvalidReviewReceipt
	}
	return receipt, classified.Data, nil
}

func validReviewReceipt(receipt reviewReceipt) bool {
	if receipt.Version != 1 || receipt.Run == "" || receipt.Candidate == "" || len(receipt.Axes) != 3 {
		return false
	}
	want := []string{"Standards", "Spec", "Coverage"}
	for index, axis := range receipt.Axes {
		if axis.Axis != want[index] || !validFindings(axis.Findings) {
			return false
		}
	}
	return true
}

func validFindings(findings []reviewFinding) bool {
	previous := ""
	for _, finding := range findings {
		if finding.ID == "" || finding.Disposition == "" || finding.ID <= previous {
			return false
		}
		previous = finding.ID
	}
	return true
}
