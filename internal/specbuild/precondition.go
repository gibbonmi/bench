package specbuild

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
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

type mutationPolicy struct {
	op            mutation
	requiresClean bool
}

// lifecycleMutationPolicies is the closed set of preconditioned operations and
// the one source of their strict-versus-provisional checkout policy.
var lifecycleMutationPolicies = []mutationPolicy{
	{mutationStart, true},
	{mutationAssign, false},
	{mutationCheckpoint, false},
	{mutationIntegrate, false},
	{mutationReview, false},
	{mutationPromote, true},
	{mutationAbandon, true},
}

var lifecycleMutations = func() []mutation {
	result := make([]mutation, len(lifecycleMutationPolicies))
	for i, policy := range lifecycleMutationPolicies {
		result[i] = policy.op
	}
	return result
}()

// name renders op for a refusal, falling back to the lifecycle rather than letting an undeclared token name no operation at all.
func (op mutation) name() string {
	for _, declared := range lifecycleMutations {
		if op == declared {
			return string(declared)
		}
	}
	return "lifecycle"
}

var errRecompose = errors.New("spec build requires recomposition: bench spec build promote")

// preconditions is the single fail-closed owner before a lifecycle mutation writes state, refs, commits, or worktrees.
func (s *Service) preconditions(op mutation, slug, specPath string, run *record, assignmentID, evidence string) (buildSubject, error) {
	subject, err := s.subject(op, specPath)
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
	if err := s.ownedAssignments(run, op); err != nil {
		abandon, found := s.operation(*run, "abandon", "apply")
		if op != mutationAbandon || !found || abandon.State != "prepared" || abandon.Result == "" {
			return buildSubject{}, err
		}
	}
	if err := s.operationEvidence(op, *run, assignmentID, evidence); err != nil {
		return buildSubject{}, err
	}
	// Abandon is the escape hatch a mid-repair run needs precisely when the tip has
	// moved, so it alone is exempt from the refusal it would otherwise trigger here.
	// Everything above this line — identity, ownership, evidence — still applies.
	if recompose && op != mutationAbandon {
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

func (s *Service) subject(op mutation, specPath string) (buildSubject, error) {
	branch, tip, err := workingSubject(s.root, op)
	if err != nil {
		return buildSubject{}, err
	}
	relative, ok := checkoutRelativePath(s.root, specPath)
	if !ok {
		return buildSubject{}, errors.New("spec build spec does not belong to working checkout")
	}
	specTip, err := benchgit.Output("-C", s.root, "rev-parse", "HEAD:"+relative)
	if err != nil || specTip == "" {
		return buildSubject{}, errors.New("spec build staged spec has no committed identity")
	}
	return buildSubject{branch: branch, tip: tip, spec: specPath, specTip: specTip}, nil
}
func workingSubject(root string, op mutation) (string, string, error) {
	branch, err := benchgit.Output("-C", root, "symbolic-ref", "--quiet", "--short", "HEAD")
	if err != nil || branch == "" {
		return "", "", fmt.Errorf("spec build %s requires a checked-out working branch", op.name())
	}
	if requiresCleanWorkingCheckout(op) {
		dirty, err := benchgit.Output("-C", root, "status", "--porcelain", "--untracked-files=all")
		if err != nil {
			return "", "", err
		}
		if dirty != "" {
			return "", "", fmt.Errorf("spec build %s requires a clean working checkout: %s", op.name(), dirty)
		}
	}
	tip, err := benchgit.Output("-C", root, "rev-parse", "HEAD^{commit}")
	if err != nil {
		return "", "", err
	}
	return branch, tip, nil
}

func requiresCleanWorkingCheckout(op mutation) bool {
	for _, policy := range lifecycleMutationPolicies {
		if policy.op == op {
			return policy.requiresClean
		}
	}
	return true
}

func checkoutRelativePath(root, path string) (string, bool) {
	relative, err := filepath.Rel(root, path)
	if err != nil || relative == "." || filepath.IsAbs(relative) || !sameOrBelow(root, path) {
		return "", false
	}
	return filepath.ToSlash(relative), true
}

func (s *Service) requireCommittedTicket(ticket Ticket) error {
	path, ok := checkoutRelativePath(s.root, ticket.Path)
	if !ok {
		return errors.New("spec build ticket does not belong to working checkout")
	}
	committed, err := benchgit.Raw("-C", s.root, "show", "HEAD:"+path)
	if err != nil || digest(string(committed)) != ticket.Digest {
		return errors.New("spec build ticket no longer matches committed subject")
	}
	indexed, err := benchgit.Raw("-C", s.root, "show", ":"+path)
	if err != nil || digest(string(indexed)) != ticket.Digest {
		return errors.New("spec build ticket no longer matches committed subject")
	}
	return nil
}

var (
	errOwnership = errors.New("spec build assignment ownership does not match durable state")
	// errAbsentCheckout marks the one ownership fault abandon may proceed through.
	errAbsentCheckout = errors.New("spec build assignment checkout is absent")
)

// ownedAssignments refuses op unless every recorded assignment still matches the
// ownership the repository itself holds. Identity — the uniqueness of ID, path, and
// owner-request, the digest, and the resolved assignment's agreement with the record —
// binds every operation, because softening it would let a hand-edited record drive
// cleanup. Liveness is the checkout on disk, and only abandon proceeds without it.
func (s *Service) ownedAssignments(run *record, op mutation) error {
	seenIDs, seenPaths, seenRequests := map[string]bool{}, map[string]bool{}, map[string]bool{}
	for key, assigned := range run.Assignments {
		if key == "" || key != assigned.Request || assigned.ID == "" || assigned.Path == "" || assigned.OwnerRequest != digest(assigned.Request) || seenIDs[assigned.ID] || seenPaths[assigned.Path] || seenRequests[assigned.OwnerRequest] {
			return errOwnership
		}
		seenIDs[assigned.ID], seenPaths[assigned.Path], seenRequests[assigned.OwnerRequest] = true, true, true
		if assigned.Released {
			continue
		}
		owned, found, err := intent.FindAssignmentByRequest(s.root, assigned.OwnerRequest)
		if err != nil || !found || owned.ID != assigned.ID || filepath.Clean(owned.Worktree) != filepath.Clean(assigned.Path) {
			return errOwnership
		}
		// A checkout that is gone is the exact state abandon exists to clean up, and the
		// recovery refs its plan enumerates hold the payload, so nothing is lost by
		// proceeding. Every other mutation writes into that checkout and must refuse.
		if err := s.liveCheckout(assigned.Path); err != nil {
			if op != mutationAbandon || !errors.Is(err, errAbsentCheckout) {
				return errOwnership
			}
		}
	}
	return nil
}

// liveCheckout reports whether path is still a checkout of this repository, separating
// an absent path from a stranger's checkout. The repository probe cannot make that
// distinction — it fails identically for both — so the path's own existence decides
// first, and a present-but-foreign path is never mistaken for a removed one.
func (s *Service) liveCheckout(path string) error {
	if _, err := os.Lstat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errAbsentCheckout
		}
		return errOwnership
	}
	common, err := benchgit.Output("-C", path, "rev-parse", "--path-format=absolute", "--git-common-dir")
	if err != nil {
		return errOwnership
	}
	rootCommon, err := benchgit.Output("-C", s.root, "rev-parse", "--path-format=absolute", "--git-common-dir")
	if err != nil || filepath.Clean(common) != filepath.Clean(rootCommon) {
		return errOwnership
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
	if err := s.ownedAssignments(&run, mutationAbandon); err != nil {
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
