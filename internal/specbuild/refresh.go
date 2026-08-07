package specbuild

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/jsonfile"
)

const debugReceiptVersion = 1

var errInvalidDebugReceipt = errors.New("invalid spec build debug receipt")

// debugReceipt is the bounded evidence a write delegate returns when its own
// deterministic repro proves the ticket is blocked by a defect outside its
// ownership fence. It authorizes exactly one lifecycle effect — refreshing the
// blocked assignment onto the repaired candidate — and nothing in it grants a
// write outside the fence.
type debugReceipt struct {
	Version       int        `json:"version"`
	Run           string     `json:"run"`
	Assignment    string     `json:"assignment"`
	Base          string     `json:"base"`
	Repro         debugRepro `json:"repro"`
	Cause         string     `json:"cause"`
	RequiredFence []string   `json:"required_fence"`
	DirtyPaths    []string   `json:"dirty_paths"`
	Resumable     bool       `json:"resumable"`
}

// debugRepro pins the exact failing command the receipt answers for. A zero
// exit is refused: a green repro proves no prerequisite defect.
type debugRepro struct {
	Command      string `json:"command"`
	Exit         int    `json:"exit"`
	OutputDigest string `json:"output_digest"`
	Produced     string `json:"produced"`
}

func readDebugReceipt(path string) (debugReceipt, []byte, error) {
	if path == "" || !filepath.IsAbs(path) {
		return debugReceipt{}, nil, errInvalidDebugReceipt
	}
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() {
		return debugReceipt{}, nil, errInvalidDebugReceipt
	}
	classified := bounds.Classify(path, bounds.ControlRecordLimit)
	if classified.State != bounds.StateParsed || !strings.HasSuffix(string(classified.Data), "\n") {
		return debugReceipt{}, nil, errInvalidDebugReceipt
	}
	var receipt debugReceipt
	if jsonfile.Decode(classified.Data, &receipt) != nil || !validDebugReceipt(receipt) {
		return debugReceipt{}, nil, errInvalidDebugReceipt
	}
	return receipt, classified.Data, nil
}

func validDebugReceipt(receipt debugReceipt) bool {
	if receipt.Version != debugReceiptVersion || receipt.Run == "" || receipt.Assignment == "" || receipt.Base == "" || receipt.Cause == "" {
		return false
	}
	if receipt.Repro.Command == "" || receipt.Repro.Exit == 0 || receipt.Repro.OutputDigest == "" || receipt.Repro.Produced == "" {
		return false
	}
	if _, err := time.Parse(time.RFC3339Nano, receipt.Repro.Produced); err != nil {
		return false
	}
	return len(receipt.RequiredFence) != 0
}

// refreshIdentity is the one construction of an assignment's preservation ref,
// mirroring checkpointIdentity so the writer and the residue enumeration derive
// the same name instead of trusting a stored copy.
func refreshIdentity(run, assignmentID string) string {
	return "refs/bench/specbuild/refresh/" + digest(run+"\x00"+assignmentID)
}

// Refresh moves one owned, uncheckpointed assignment onto the exact current
// candidate tip, preserving its attributed in-fence work byte-for-byte. The
// authority is a validated debug receipt whose required fence names at least one
// path outside the assignment's own fence — the delegate proved it is blocked by
// a prerequisite repair it must not author. The dirty payload is committed
// against the recorded base and bound to a durable ref before the worktree is
// touched, so an interrupted refresh re-enters and converges rather than losing
// either the repair or the assignment.
func (s *Service) Refresh(ctx context.Context, slug, ticketArg, request, evidence string) (Assignment, Status, error) {
	if strings.TrimSpace(request) == "" {
		return Assignment{}, Status{}, errors.New("spec build assignment request is required")
	}
	if _, err := s.resolve(slug); err != nil {
		return Assignment{}, Status{}, err
	}
	release, err := s.lock(slug)
	if err != nil {
		return Assignment{}, Status{}, err
	}
	defer release()
	run, found, err := s.load(slug)
	if err != nil || !found {
		if err == nil {
			err = errors.New("spec build run does not exist")
		}
		return Assignment{}, Status{}, err
	}
	if _, err := s.preconditions(mutationAssign, slug, run.Spec, &run, "", ""); err != nil {
		return Assignment{}, Status{}, err
	}
	ticket, err := ParseTicket(run.Spec, ticketArg)
	if err != nil {
		return Assignment{}, Status{}, err
	}
	requestID := digest(run.Run + "\x00" + request)
	key := requestID
	assigned, ok := run.Assignments[requestID]
	if !ok {
		return Assignment{}, Status{}, errors.New("spec build refresh names no existing assignment")
	}
	if assigned.Ticket != ticket.Path {
		return Assignment{}, Status{}, errors.New("spec build assignment request conflicts with another ticket")
	}
	if assigned.Released || assigned.Checkpoint != "" || assigned.CheckpointRef != "" {
		return Assignment{}, Status{}, errors.New("spec build refresh requires an uncheckpointed owned assignment")
	}
	receipt, raw, err := readDebugReceipt(evidence)
	if err != nil {
		return Assignment{}, Status{}, err
	}
	refreshRequest := assigned.ID + "\x00" + receipt.Base
	op, opFound := s.operation(run, "refresh", refreshRequest)
	if opFound && op.Input != digest(string(raw)) {
		return Assignment{}, Status{}, errors.New("spec build refresh request conflicts with different inputs")
	}
	if opFound && op.State == "completed" {
		return assigned.public(), run.status(), nil
	}
	// The receipt's base is the one the payload was preserved against. A fresh
	// refresh requires it to be the assignment's recorded base; a re-entry after
	// an interrupted state save also accepts a record already advanced to the
	// candidate, because the journaled input digest above has proven the receipt
	// is the same one that authorized the interrupted run.
	resuming := opFound && op.Result != ""
	base := receipt.Base
	if !resuming && base != assigned.Base {
		return Assignment{}, Status{}, errInvalidDebugReceipt
	}
	if resuming && assigned.Base != base && assigned.Base != run.CandidateTip {
		return Assignment{}, Status{}, errors.New("spec build refresh preservation drifted")
	}
	if err := s.validateDebugReceipt(run, assigned, receipt, evidence); err != nil {
		return Assignment{}, Status{}, err
	}
	if run.CandidateTip == base {
		return Assignment{}, Status{}, errors.New("spec build refresh has no candidate advance to apply")
	}
	if err := s.requireCommittedTicket(ticket); err != nil {
		return Assignment{}, Status{}, err
	}
	if err := requireReciprocalEdges(run.Spec, ticket); err != nil {
		return Assignment{}, Status{}, err
	}
	// Refresh is the one caller allowed to move the recorded pin: a sibling
	// docs-repair may have rewritten the committed ticket since assign, and
	// carrying the preserved assignment onto the advanced candidate means
	// carrying it onto that rewrite too, not refusing forever on the stale
	// digest checkpoint/integrate would otherwise compare against.
	assigned, err = s.refreshTicketPin(run, assigned, ticket)
	if err != nil {
		return Assignment{}, Status{}, err
	}
	var preserved, tree string
	if resuming {
		preserved = op.Result
		parent, parentErr := s.gitOutput(ctx, "rev-parse", preserved+"^")
		if tree, err = s.gitOutput(ctx, "rev-parse", preserved+"^{tree}"); parentErr != nil || err != nil || parent != base {
			return Assignment{}, Status{}, errors.New("spec build refresh preservation drifted")
		}
	} else {
		// The live worktree is only read before any mutation; once the payload is
		// committed, the preservation commit is the single source, so a re-entry
		// after interruption never re-reads a half-refreshed checkout.
		if tree = benchgit.TreeHash(assigned.Path); tree == "none" {
			return Assignment{}, Status{}, errors.New("spec build refresh cannot read the assignment worktree")
		}
	}
	if err := s.validateRefreshPayload(ctx, run, assigned, receipt, base, tree); err != nil {
		return Assignment{}, Status{}, err
	}
	if !resuming {
		if _, completed, err := s.beginOperation(&run, "refresh", refreshRequest, string(raw)); err != nil {
			return Assignment{}, Status{}, err
		} else if completed {
			return Assignment{}, Status{}, errors.New("spec build refresh operation is incomplete")
		}
		preserved, err = s.gitOutput(ctx, "commit-tree", tree, "-p", base, "-m", "bench refresh run="+run.Run+" assignment="+assigned.ID+" base="+base)
		if err != nil {
			return Assignment{}, Status{}, fmt.Errorf("create preservation commit: %w", err)
		}
		if err := s.recordOperation(&run, "refresh", refreshRequest, preserved, false); err != nil {
			return Assignment{}, Status{}, err
		}
	}
	ref := refreshIdentity(run.Run, assigned.ID)
	if !refAt(s.root, ref, preserved) {
		if err := updateRef(s.root, ref, preserved, ""); err != nil {
			return Assignment{}, Status{}, fmt.Errorf("bind preservation commit: %w", err)
		}
	}
	if err := s.faultAt("refresh/preserve"); err != nil {
		return assigned.public(), run.status(), err
	}
	patch, err := s.checkpointPatch(ctx, base, preserved)
	if err != nil {
		return Assignment{}, Status{}, fmt.Errorf("record preservation patch: %w", err)
	}
	newTree, err := s.replayCheckpoint(ctx, run.CandidateTip, base, preserved, patch)
	if err != nil {
		return Assignment{}, Status{}, fmt.Errorf("spec build refresh conflicts with the repaired candidate: %w", err)
	}
	current, err := refValue(s.root, run.Candidate)
	if err != nil || current != run.CandidateTip {
		return Assignment{}, Status{}, errors.New("spec build candidate moved during refresh")
	}
	if _, err := s.worktreeGit(ctx, assigned.Path, "reset", "--hard", run.CandidateTip); err != nil {
		return assigned.public(), run.status(), fmt.Errorf("refresh assignment checkout: %w", err)
	}
	if _, err := s.worktreeGit(ctx, assigned.Path, "clean", "-fd"); err != nil {
		return assigned.public(), run.status(), fmt.Errorf("refresh assignment checkout: %w", err)
	}
	if len(patch) != 0 {
		if _, err := s.worktreeApply(ctx, assigned.Path, patch); err != nil {
			return assigned.public(), run.status(), fmt.Errorf("replay preserved work: %w", err)
		}
	}
	if got := benchgit.TreeHash(assigned.Path); got != newTree {
		return assigned.public(), run.status(), errors.New("spec build refresh replay is not byte-identical")
	}
	if err := s.faultAt("refresh/worktree"); err != nil {
		return assigned.public(), run.status(), err
	}
	assigned.Base = run.CandidateTip
	run.Assignments[key] = assigned
	if err := s.save(run); err != nil {
		return Assignment{}, Status{}, err
	}
	if err := s.faultAt("refresh/state"); err != nil {
		return assigned.public(), run.status(), err
	}
	if err := s.recordOperation(&run, "refresh", refreshRequest, preserved, true); err != nil {
		return Assignment{}, Status{}, err
	}
	return assigned.public(), run.status(), nil
}

// validateDebugReceipt binds the receipt to exactly one live blocked assignment.
// The required fence must name at least one path the assignment does not own —
// a receipt entirely inside the fence describes ordinary ticket work, and the
// ordinary checkpoint path is its route.
func (s *Service) validateDebugReceipt(run record, assigned assignment, receipt debugReceipt, evidence string) error {
	if sameOrBelow(assigned.Path, evidence) || receipt.Run != run.Run || receipt.Assignment != assigned.ID {
		return errInvalidDebugReceipt
	}
	created, err := time.Parse(time.RFC3339Nano, assigned.Created)
	if err != nil || !timeAfter(receipt.Repro.Produced, created) {
		return errInvalidDebugReceipt
	}
	outside := false
	for _, path := range receipt.RequiredFence {
		if !insideFence([]string{filepath.ToSlash(path)}, assigned.Fence) {
			outside = true
			break
		}
	}
	if !outside {
		return errors.New("spec build debug receipt names no path outside the ownership fence")
	}
	if err := s.assignedOwnership(assigned); err != nil {
		return errOwnership
	}
	return nil
}

// refreshTicketPin re-pins assigned's recorded ticket digest and rows to the
// current committed text when it changed since assign, so a legitimate
// sibling docs-repair reaches the preserved assignment instead of stranding
// it behind validateIntegrationTicket's stale-digest refusal. An Ownership
// fence change is not this kind of rewrite — it moves the assignment's write
// envelope, which refresh does not own — so it still refuses, naming the
// fence. Any other change re-validates against the same assign-side policy
// (ContractsAnchored, requireClosure, requireCoversMapping) rather than a
// second copy of it; requireReciprocalEdges has already run by the time this
// is called, so a rewrite still missing the reciprocal edge never reaches it.
func (s *Service) refreshTicketPin(run record, assigned assignment, ticket Ticket) (assignment, error) {
	if ticket.Digest == assigned.TicketDigest {
		return assigned, nil
	}
	if !sameStrings(ticket.Fence, assigned.Fence) {
		return assignment{}, fmt.Errorf("spec build refresh cannot re-pin ticket %s: Ownership fence changed to %s", filepath.Base(ticket.Path), strings.Join(ticket.Fence, ", "))
	}
	if !ticket.ContractsAnchored() {
		return assignment{}, fmt.Errorf("spec build ticket %s declares a contract crossing no path in its ownership fence", filepath.Base(ticket.Path))
	}
	if err := requireClosure(ticket); err != nil {
		return assignment{}, err
	}
	if err := requireCoversMapping(run.Spec, ticket); err != nil {
		return assignment{}, err
	}
	assigned.TicketDigest, assigned.Rows = ticket.Digest, ticket.Rows
	return assigned, nil
}

// validateRefreshPayload grades the dirty payload the refresh will carry: every
// changed path stays inside the assignment's own fence, and the receipt's
// dirty-path claim matches the payload exactly, so a delegate cannot smuggle an
// out-of-fence edit through the refresh it requested.
func (s *Service) validateRefreshPayload(ctx context.Context, run record, assigned assignment, receipt debugReceipt, base, tree string) error {
	dirty, err := s.changedPaths(ctx, base, tree)
	if err != nil {
		return fmt.Errorf("read assignment payload: %w", err)
	}
	if len(dirty) != 0 && !insideFence(dirty, assigned.Fence) {
		return errors.New("spec build refresh payload leaves the ownership fence")
	}
	if !sameStrings(dirty, sortedUnique(receipt.DirtyPaths)) {
		return errInvalidDebugReceipt
	}
	if _, err := validateIntegrationTicket(run, assigned); err != nil {
		return err
	}
	return nil
}

func (s *Service) worktreeGit(ctx context.Context, path string, args ...string) (string, error) {
	output, err := s.runner.Run(ctx, Command{Program: "git", Args: append([]string{"-C", path}, args...)})
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return "", err
		}
		return "", errors.New(strings.TrimSpace(output))
	}
	return output, nil
}

func (s *Service) worktreeApply(ctx context.Context, path string, patch []byte) (string, error) {
	output, err := s.runner.Run(ctx, Command{Program: "git", Args: []string{"-C", path, "apply", "--whitespace=nowarn"}, Input: strings.NewReader(string(patch))})
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return "", err
		}
		return "", errors.New(strings.TrimSpace(output))
	}
	return output, nil
}
