package specbuild

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"

	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

// mutation names lifecycle operations that must prove their subject before writing.
// The vocabulary intentionally includes transitions that will join this owner later.
type mutation string

const (
	mutationStart      mutation = "start"
	mutationAssign     mutation = "assign"
	mutationCheckpoint mutation = "checkpoint"
	mutationIntegrate  mutation = "integrate"
	mutationReview     mutation = "review"
	mutationPromote    mutation = "promote"
	mutationAbandon    mutation = "abandon"
)

var errRecompose = errors.New("spec build requires recomposition: bench spec build promote")

// preconditions is the single fail-closed owner for lifecycle mutations. It runs
// before a command can write durable state, refs, commits, or worktrees.
func (s *Service) preconditions(op mutation, slug, specPath string, run *record, assignmentID, evidence string) (buildSubject, error) {
	subject, err := s.subject(specPath)
	if err != nil {
		return buildSubject{}, err
	}
	if run == nil {
		if err := s.operationEvidence(op, record{}, assignmentID, evidence); err != nil {
			return buildSubject{}, err
		}
		return subject, nil
	}
	if run.Branch != subject.branch {
		return buildSubject{}, errors.New("spec build working checkout does not match recorded subject")
	}
	if run.Base != subject.tip {
		if !recognizedAdvance(s.root, run.Base, subject.tip) {
			return buildSubject{}, errors.New("spec build working checkout does not match recorded subject")
		}
		return buildSubject{}, fmt.Errorf("%w %s", errRecompose, slug)
	}
	if run.Spec != subject.spec || run.SpecTip != subject.specTip {
		return buildSubject{}, errors.New("spec build staged spec no longer matches recorded subject")
	}
	if !refAt(s.root, run.Candidate, run.CandidateTip) {
		return buildSubject{}, errors.New("spec build candidate no longer matches durable tip")
	}
	if err := s.ownedAssignments(run); err != nil {
		return buildSubject{}, err
	}
	if err := s.operationEvidence(op, *run, assignmentID, evidence); err != nil {
		return buildSubject{}, err
	}
	return subject, nil
}

func recognizedAdvance(root, base, tip string) bool {
	return exec.Command("git", "-C", root, "merge-base", "--is-ancestor", base, tip).Run() == nil
}

type buildSubject struct {
	branch, tip, spec, specTip string
}

func (s *Service) subject(specPath string) (buildSubject, error) {
	branch, tip, err := workingSubject(s.root)
	if err != nil {
		return buildSubject{}, err
	}
	relative, err := filepath.Rel(s.root, specPath)
	if err != nil || relative == "." || filepath.IsAbs(relative) || !sameOrBelow(s.root, specPath) {
		return buildSubject{}, errors.New("spec build spec does not belong to working checkout")
	}
	specTip, err := benchgit.Output("-C", s.root, "rev-parse", "HEAD:"+filepath.ToSlash(relative))
	if err != nil || specTip == "" {
		return buildSubject{}, errors.New("spec build staged spec has no committed identity")
	}
	return buildSubject{branch: branch, tip: tip, spec: specPath, specTip: specTip}, nil
}

func (s *Service) ownedAssignments(run *record) error {
	seenIDs, seenPaths, seenRequests := map[string]bool{}, map[string]bool{}, map[string]bool{}
	for key, assigned := range run.Assignments {
		if key == "" || key != assigned.Request || assigned.ID == "" || assigned.Path == "" || assigned.OwnerRequest != digest(assigned.Request) || seenIDs[assigned.ID] || seenPaths[assigned.Path] || seenRequests[assigned.OwnerRequest] {
			return errors.New("spec build assignment ownership does not match durable state")
		}
		seenIDs[assigned.ID], seenPaths[assigned.Path], seenRequests[assigned.OwnerRequest] = true, true, true
		if assigned.Released {
			continue
		}
		owned, found, err := intent.FindAssignmentByRequest(s.root, assigned.OwnerRequest)
		if err != nil || !found || owned.ID != assigned.ID || filepath.Clean(owned.Worktree) != filepath.Clean(assigned.Path) {
			return errors.New("spec build assignment ownership does not match durable state")
		}
		common, err := benchgit.Output("-C", assigned.Path, "rev-parse", "--path-format=absolute", "--git-common-dir")
		if err != nil {
			return errors.New("spec build assignment ownership does not match durable state")
		}
		rootCommon, err := benchgit.Output("-C", s.root, "rev-parse", "--path-format=absolute", "--git-common-dir")
		if err != nil || filepath.Clean(common) != filepath.Clean(rootCommon) {
			return errors.New("spec build assignment ownership does not match durable state")
		}
	}
	return nil
}

func (s *Service) operationEvidence(op mutation, run record, assignmentID, evidence string) error {
	switch op {
	case mutationStart:
		if s.gate == nil {
			return errors.New("spec build start requires a gate owner")
		}
	case mutationAssign:
		if s.worktrees == nil {
			return errors.New("spec build assign requires a worktree owner")
		}
	case mutationCheckpoint:
		if _, _, ok := assignmentFor(run, assignmentID); !ok {
			return errors.New("spec build assignment does not exist")
		}
	case mutationIntegrate:
		_, assigned, ok := assignmentFor(run, assignmentID)
		if !ok || (assigned.Checkpoint == "" && assigned.CheckpointRef == "") {
			return errors.New("spec build assignment has no verified checkpoint")
		}
	case mutationReview:
		if evidence == "" {
			return errInvalidReviewReceipt
		}
	case mutationPromote, mutationAbandon:
		return nil
	default:
		return errors.New("spec build mutation has no precondition contract")
	}
	return nil
}
