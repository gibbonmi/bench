// Package specbuild owns the durable lifecycle of a reviewed spec build.
package specbuild

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

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

// AbandonOwner plans and applies exact recovery-aware owned-worktree cleanup.
type AbandonOwner interface {
	PlanAbandon(context.Context, string, string, string) (string, error)
	ApplyAbandon(context.Context, string, string, string, string) error
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
	runner             Runner
	fault              func(string) error
}

// New constructs a lifecycle service rooted at one working checkout.
func New(root string, gate GateOwner, worktrees WorktreeOwner) *Service {
	return NewWithRunner(root, gate, worktrees, processRunner{})
}

// Runner executes lifecycle subprocesses as one cancellable process group.
// Promote and abandon use the same seam when their long-running work arrives.
type Runner interface {
	Output(context.Context, string, ...string) (string, error)
	Run(context.Context, Command) (string, error)
}

// Command describes one lifecycle subprocess without leaking execution policy.
type Command struct {
	Program   string
	Args, Env []string
	Input     io.Reader
}

// NewWithRunner constructs a lifecycle service with its process runner seam.
func NewWithRunner(root string, gate GateOwner, worktrees WorktreeOwner, runner Runner) *Service {
	if runner == nil {
		runner = processRunner{}
	}
	return &Service{root: root, gate: gate, worktrees: worktrees, runner: runner}
}

type processRunner struct{}

func (processRunner) Output(ctx context.Context, program string, args ...string) (string, error) {
	return (processRunner{}).Run(ctx, Command{Program: program, Args: args})
}

func (processRunner) Run(ctx context.Context, command Command) (string, error) {
	cmd := exec.Command(command.Program, command.Args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Env, cmd.Stdin = command.Env, command.Input
	output := &strings.Builder{}
	cmd.Stdout, cmd.Stderr = output, output
	if err := cmd.Start(); err != nil {
		return "", err
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		return output.String(), err
	case <-ctx.Done():
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
		select {
		case <-done:
		case <-time.After(250 * time.Millisecond):
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			<-done
		}
		return output.String(), ctx.Err()
	}
}

const operationLimit = 64

type operation struct {
	Command, Request, Input, Result, State string
}

func operationID(command, request string) string { return digest(command + "\x00" + request) }

func (s *Service) beginOperation(run *record, command, request, input string) (operation, bool, error) {
	key, inputDigest := operationID(command, request), digest(input)
	if prior, found := run.Operations[key]; found {
		if prior.Command != command || prior.Request != request || prior.Input != inputDigest {
			return operation{}, false, fmt.Errorf("spec build %s request conflicts with different inputs", command)
		}
		return prior, prior.State == "completed", nil
	}
	if len(run.Operations) >= operationLimit {
		return operation{}, false, errors.New("spec build operation journal is full")
	}
	op := operation{Command: command, Request: request, Input: inputDigest, State: "prepared"}
	run.Operations[key] = op
	if err := s.save(*run); err != nil {
		return operation{}, false, err
	}
	return op, false, nil
}

func (s *Service) recordOperation(run *record, command, request, result string, completed bool) error {
	key := operationID(command, request)
	op, found := run.Operations[key]
	if !found || op.State != "prepared" {
		return errors.New("spec build operation journal is incomplete")
	}
	op.Result = result
	if completed {
		op.State = "completed"
	}
	run.Operations[key] = op
	return s.save(*run)
}

func (s *Service) operation(run record, command, request string) (operation, bool) {
	op, found := run.Operations[operationID(command, request)]
	return op, found
}

func (s *Service) faultAt(point string) error {
	if s.fault == nil {
		return nil
	}
	return s.fault(point)
}

func sortRefs(refs []AbandonmentRef) {
	sort.Slice(refs, func(i, j int) bool { return refs[i].Name < refs[j].Name })
}

func abandonmentFacts(plan AbandonmentPlan) string {
	var facts strings.Builder
	facts.WriteString("bench-specbuild-abandon/v1")
	for _, worktree := range plan.Worktrees {
		facts.WriteString("\x00worktree\x00" + worktree.ID + "\x00" + worktree.Path + "\x00" + worktree.Request + "\x00" + worktree.OwnerFingerprint)
	}
	for _, group := range []struct {
		name string
		refs []AbandonmentRef
	}{{"provisional", plan.ProvisionalRefs}, {"checkpoint", plan.UnintegratedCheckpoints}, {"recovery", plan.RecoveryRefs}} {
		for _, ref := range group.refs {
			facts.WriteString("\x00" + group.name + "\x00" + ref.Name + "\x00" + ref.Object)
		}
	}
	return facts.String()
}

type abandonmentJournal struct {
	Original, Current AbandonmentPlan
}

func (s *Service) recordAbandonment(run *record, journal abandonmentJournal, completed bool) error {
	encoded, err := json.Marshal(journal)
	if err != nil {
		return fmt.Errorf("encode abandonment journal: %w", err)
	}
	return s.recordOperation(run, "abandon", "apply", string(encoded), completed)
}

func (s *Service) reconcileAbandonment(ctx context.Context, run *record, journal *abandonmentJournal) error {
	abandoner, ok := abandonOwnerFrom(s.worktrees)
	if !ok {
		return errors.New("spec build abandon requires a plan-capable worktree owner")
	}
	changed := false
	for _, worktree := range journal.Current.Worktrees {
		if _, err := os.Lstat(worktree.Path); !errors.Is(err, os.ErrNotExist) {
			if err != nil {
				return err
			}
			continue
		}
		fingerprint, err := abandoner.PlanAbandon(ctx, s.root, worktree.Request, worktree.Path)
		if err != nil || fingerprint != worktree.OwnerFingerprint {
			return errors.New("spec build abandonment recovery evidence drifted")
		}
		if err := abandoner.ApplyAbandon(ctx, s.root, worktree.Request, worktree.Path, fingerprint); err != nil {
			return errors.New("spec build abandonment recovery evidence drifted")
		}
		key, assigned, ok := assignmentFor(*run, worktree.ID)
		if !ok || assigned.Released {
			return errors.New("spec build abandonment assignment drifted")
		}
		assigned.CleanupPending, assigned.Released = false, true
		run.Assignments[key] = assigned
		if err := s.save(*run); err != nil {
			return err
		}
		changed = true
	}
	if !changed {
		return nil
	}
	next, err := s.abandonmentPlan(ctx, *run)
	if err != nil {
		return err
	}
	journal.Current = next
	return s.recordAbandonment(run, *journal, false)
}

func abandonOwnerFrom(owner WorktreeOwner) (AbandonOwner, bool) {
	abandoner, ok := owner.(AbandonOwner)
	return abandoner, ok
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
	receipt, raw, err := readReviewReceipt(evidence)
	if err != nil || receipt.Run != run.Run || receipt.Candidate != run.CandidateTip || !refAt(s.root, run.Candidate, run.CandidateTip) {
		return Status{}, errInvalidReviewReceipt
	}
	digest := digest(string(raw))
	if run.Review != nil {
		if run.Review.Candidate == run.CandidateTip && run.Review.Digest == digest {
			if op, found := s.operation(run, "review", run.CandidateTip); found && op.State == "prepared" {
				if err := s.recordOperation(&run, "review", run.CandidateTip, digest, true); err != nil {
					return Status{}, err
				}
			}
			return run.status(), nil
		}
		return Status{}, errors.New("spec build review request conflicts with different inputs")
	}
	if _, err := s.preconditions(mutationReview, slug, run.Spec, &run, "", evidence); err != nil {
		return Status{}, err
	}
	if _, completed, err := s.beginOperation(&run, "review", run.CandidateTip, string(raw)); err != nil {
		return Status{}, err
	} else if completed {
		return Status{}, errors.New("spec build review operation is incomplete")
	}
	run.Review = &reviewEvidence{Candidate: receipt.Candidate, Axes: receipt.Axes, Digest: digest}
	if err := s.save(run); err != nil {
		return Status{}, err
	}
	if err := s.faultAt("review/state"); err != nil {
		return run.status(), err
	}
	if err := s.recordOperation(&run, "review", run.CandidateTip, digest, true); err != nil {
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
