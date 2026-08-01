package specbuild

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

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

// preconditions is the single fail-closed owner before a lifecycle mutation writes state, refs, commits, or worktrees.
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
	recompose := run.Base != subject.tip
	if recompose {
		if !recognizedAdvance(s.root, run.Base, subject.tip) {
			return buildSubject{}, errors.New("spec build working checkout does not match recorded subject")
		}
	}
	if run.Spec != subject.spec || run.SpecTip != subject.specTip {
		return buildSubject{}, errors.New("spec build staged spec no longer matches recorded subject")
	}
	if !refAt(s.root, run.Candidate, run.CandidateTip) {
		return buildSubject{}, errors.New("spec build candidate no longer matches durable tip")
	}
	if err := s.ownedAssignments(run); err != nil {
		abandon, found := s.operation(*run, "abandon", "apply")
		if op != mutationAbandon || !found || abandon.State != "prepared" || abandon.Result == "" {
			return buildSubject{}, err
		}
	}
	if err := s.operationEvidence(op, *run, assignmentID, evidence); err != nil {
		return buildSubject{}, err
	}
	if recompose {
		return subject, fmt.Errorf("%w %s", errRecompose, slug)
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
func workingSubject(root string) (string, string, error) {
	branch, err := benchgit.Output("-C", root, "symbolic-ref", "--quiet", "--short", "HEAD")
	if err != nil || branch == "" {
		return "", "", errors.New("spec build start requires a checked-out working branch")
	}
	dirty, err := benchgit.Output("-C", root, "status", "--porcelain", "--untracked-files=all")
	if err != nil {
		return "", "", err
	}
	if dirty != "" {
		return "", "", fmt.Errorf("spec build start requires a clean working checkout: %s", dirty)
	}
	tip, err := benchgit.Output("-C", root, "rev-parse", "HEAD^{commit}")
	if err != nil {
		return "", "", err
	}
	return branch, tip, nil
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

func retainTerminalAttempt(run record) []json.RawMessage {
	history := append([]json.RawMessage(nil), run.History...)
	run.History = nil
	raw, _ := json.Marshal(run)
	return append(history, raw)
}

func runIdentity(specPath, attempt string) (string, string) {
	identity := digest(specPath)
	if attempt != "" {
		identity = digest(specPath + "\x00" + attempt)
	}
	return identity, candidateIdentity(identity)
}

// Abandon returns the read-only inventory that an abandonment apply must match.
func (s *Service) Abandon(ctx context.Context, slug string) (AbandonmentPlan, error) {
	if _, err := s.resolve(slug); err != nil {
		return AbandonmentPlan{}, err
	}
	release, err := s.lock(slug)
	if err != nil {
		return AbandonmentPlan{}, err
	}
	defer release()
	run, found, err := s.load(slug)
	if err != nil || !found {
		if err == nil {
			err = errors.New("spec build run does not exist")
		}
		return AbandonmentPlan{}, err
	}
	if err := s.ownedAssignments(&run); err != nil {
		return AbandonmentPlan{}, err
	}
	return s.abandonmentPlan(ctx, run)
}

// ApplyAbandon releases the exact planned owned worktrees and retains run evidence.
func (s *Service) ApplyAbandon(ctx context.Context, slug, fingerprint string) (Status, error) {
	if fingerprint == "" {
		return Status{}, errors.New("spec build abandon fingerprint is required")
	}
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
	op, found := s.operation(run, "abandon", "apply")
	if found && op.State == "completed" {
		if op.Input != digest(fingerprint) {
			return Status{}, errors.New("spec build abandon request conflicts with different inputs")
		}
		return run.status(), nil
	}
	if _, err := s.preconditions(mutationAbandon, slug, run.Spec, &run, "", ""); err != nil {
		return Status{}, err
	}
	journal := abandonmentJournal{}
	if found {
		if op.Input != digest(fingerprint) {
			return Status{}, errors.New("spec build abandon request conflicts with different inputs")
		}
		if op.Result != "" {
			if err := json.Unmarshal([]byte(op.Result), &journal); err != nil || journal.Original.Fingerprint != fingerprint {
				return Status{}, errors.New("spec build abandonment journal is invalid")
			}
		}
	}
	actual, err := s.abandonmentPlan(ctx, run)
	if err != nil {
		return Status{}, err
	}
	if found && op.Result == "" || !found {
		if actual.Fingerprint != fingerprint {
			return Status{}, errors.New("spec build abandon plan drifted; request a fresh plan")
		}
		journal = abandonmentJournal{Original: actual, Current: actual}
		if !found {
			if _, _, err := s.beginOperation(&run, "abandon", "apply", fingerprint); err != nil {
				return Status{}, err
			}
		}
		if err := s.recordAbandonment(&run, journal, false); err != nil {
			return Status{}, err
		}
	}
	if found && op.Result != "" {
		if err := s.reconcileAbandonment(ctx, &run, &journal); err != nil {
			return Status{}, err
		}
		actual, err = s.abandonmentPlan(ctx, run)
		if err != nil {
			return Status{}, err
		}
	}
	if actual.Fingerprint != journal.Current.Fingerprint {
		return Status{}, errors.New("spec build abandon plan drifted; request a fresh plan")
	}
	plan := journal.Current
	if len(plan.Worktrees) != 0 {
		releaser, ok := abandonOwnerFrom(s.worktrees)
		if !ok {
			return Status{}, errors.New("spec build abandon requires a release-capable worktree owner")
		}
		for len(plan.Worktrees) != 0 {
			worktree := plan.Worktrees[0]
			key, assigned, ok := assignmentFor(run, worktree.ID)
			if !ok {
				return Status{}, errors.New("spec build abandonment assignment drifted")
			}
			if err := releaser.ApplyAbandon(ctx, s.root, worktree.Request, worktree.Path, worktree.OwnerFingerprint); err != nil {
				return run.status(), fmt.Errorf("release abandoned assignment: %w", err)
			}
			if err := s.faultAt("abandon/owner-apply"); err != nil {
				return run.status(), err
			}
			assigned.CleanupPending, assigned.Released = false, true
			run.Assignments[key] = assigned
			if err := s.save(run); err != nil {
				return Status{}, err
			}
			next, err := s.abandonmentPlan(ctx, run)
			if err != nil {
				return Status{}, err
			}
			journal.Current = next
			if err := s.recordAbandonment(&run, journal, false); err != nil {
				return Status{}, err
			}
			plan = next
			if err := s.faultAt("abandon/release"); err != nil {
				return run.status(), err
			}
		}
	}
	run.Terminal = true
	if err := s.save(run); err != nil {
		return Status{}, err
	}
	if err := s.recordAbandonment(&run, journal, true); err != nil {
		return Status{}, err
	}
	return run.status(), nil
}

// AbandonmentPlan identifies every Git-visible fact an abandon apply rechecks.
type AbandonmentPlan struct {
	Fingerprint             string
	Worktrees               []AbandonmentWorktree
	ProvisionalRefs         []AbandonmentRef
	UnintegratedCheckpoints []AbandonmentRef
	RecoveryRefs            []AbandonmentRef
}

// AbandonmentWorktree is one active assignment that the owner may release.
type AbandonmentWorktree struct{ ID, Path, Request, OwnerFingerprint string }

// AbandonmentRef is one retained Git ref and its exact object identity.
type AbandonmentRef struct{ Name, Object string }

func (s *Service) abandonmentPlan(ctx context.Context, run record) (AbandonmentPlan, error) {
	plan := AbandonmentPlan{}
	abandoner, ok := abandonOwnerFrom(s.worktrees)
	if !ok && len(run.Assignments) != 0 {
		return AbandonmentPlan{}, errors.New("spec build abandon requires a plan-capable worktree owner")
	}
	for _, assigned := range run.Assignments {
		if assigned.Integrated == "" && assigned.CheckpointRef != "" {
			ref, err := s.abandonmentRef(ctx, assigned.CheckpointRef)
			if err != nil {
				return AbandonmentPlan{}, err
			}
			plan.UnintegratedCheckpoints = append(plan.UnintegratedCheckpoints, ref)
		}
		if assigned.Released {
			continue
		}
		fingerprint, err := abandoner.PlanAbandon(ctx, s.root, assigned.Request, assigned.Path)
		if err != nil {
			return AbandonmentPlan{}, fmt.Errorf("inventory assignment worktree: %w", err)
		}
		plan.Worktrees = append(plan.Worktrees, AbandonmentWorktree{ID: assigned.ID, Path: assigned.Path, Request: assigned.Request, OwnerFingerprint: fingerprint})
	}
	var err error
	if plan.ProvisionalRefs, err = s.abandonmentRefs(ctx, "refs/bench/specbuild/"); err != nil {
		return AbandonmentPlan{}, err
	}
	if plan.RecoveryRefs, err = s.abandonmentRefs(ctx, "refs/bench/recovery/"); err != nil {
		return AbandonmentPlan{}, err
	}
	sort.Slice(plan.Worktrees, func(i, j int) bool { return plan.Worktrees[i].ID < plan.Worktrees[j].ID })
	sortRefs(plan.ProvisionalRefs)
	sortRefs(plan.UnintegratedCheckpoints)
	sortRefs(plan.RecoveryRefs)
	plan.Fingerprint = digest(abandonmentFacts(plan))
	return plan, nil
}

func (s *Service) abandonmentRefs(ctx context.Context, prefix string) ([]AbandonmentRef, error) {
	out, err := s.gitOutput(ctx, "for-each-ref", "--format=%(refname) %(objectname)", prefix)
	if err != nil {
		return nil, fmt.Errorf("inventory abandonment refs: %w", err)
	}
	refs := make([]AbandonmentRef, 0)
	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 {
			refs = append(refs, AbandonmentRef{Name: fields[0], Object: fields[1]})
		}
	}
	return refs, nil
}
func (s *Service) abandonmentRef(ctx context.Context, name string) (AbandonmentRef, error) {
	object, err := s.gitOutput(ctx, "rev-parse", "--verify", name+"^{commit}")
	if err != nil {
		return AbandonmentRef{}, fmt.Errorf("inventory unintegrated checkpoint: %w", err)
	}
	return AbandonmentRef{Name: name, Object: object}, nil
}
