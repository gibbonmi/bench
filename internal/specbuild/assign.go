package specbuild

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Assign creates or resumes one owned worktree for ticket in slug.
func (s *Service) Assign(ctx context.Context, slug, ticketArg, request string) (Assignment, Status, error) {
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
	if err := s.requireCommittedTicket(ticket); err != nil {
		return Assignment{}, Status{}, err
	}
	if !ticket.ContractsAnchored() {
		return Assignment{}, Status{}, fmt.Errorf("spec build ticket %s declares a contract crossing no path in its ownership fence", filepath.Base(ticket.Path))
	}
	requestID := digest(run.Run + "\x00" + request)
	op, completed, err := s.beginOperation(&run, "assign", requestID, ticket.Digest)
	if err != nil {
		return Assignment{}, Status{}, err
	}
	if existing, ok := run.Assignments[requestID]; ok {
		if existing.Ticket != ticket.Path {
			return Assignment{}, Status{}, errors.New("spec build assignment request conflicts with another ticket")
		}
		if !completed {
			if err := s.recordOperation(&run, "assign", requestID, existing.ID, true); err != nil {
				return Assignment{}, Status{}, err
			}
		}
		return existing.public(), run.status(), nil
	}
	if completed {
		return Assignment{}, Status{}, errors.New("spec build assignment operation is incomplete")
	}
	owned, err := s.worktrees.Create(ctx, s.root, requestID, ticket.Title, run.CandidateTip)
	if err != nil {
		return Assignment{}, Status{}, fmt.Errorf("create assignment worktree: %w", err)
	}
	if owned.ID == "" || owned.Path == "" {
		return Assignment{}, Status{}, errors.New("worktree owner returned an incomplete assignment")
	}
	ownerResult := owned.ID + "\x00" + owned.Path
	if op.Result != "" && op.Result != ownerResult {
		return Assignment{}, Status{}, errors.New("spec build assignment owner result conflicts with prepared request")
	}
	if op.Result == "" {
		if err := s.recordOperation(&run, "assign", requestID, ownerResult, false); err != nil {
			return Assignment{}, Status{}, err
		}
	}
	if err := s.faultAt("assign/worktree"); err != nil {
		return Assignment{}, run.status(), err
	}
	stored := assignment{ID: owned.ID, Path: owned.Path, Branch: owned.Branch, Base: run.CandidateTip, Request: requestID, OwnerRequest: digest(requestID), Ticket: ticket.Path, TicketDigest: ticket.Digest, Created: time.Now().UTC().Format(time.RFC3339Nano), Rows: ticket.Rows, Fence: ticket.Fence, Assumptions: ticket.Assumptions}
	run.Assignments[requestID] = stored
	if err := s.save(run); err != nil {
		return Assignment{}, Status{}, err
	}
	if err := s.faultAt("assign/state"); err != nil {
		return stored.public(), run.status(), err
	}
	if err := s.recordOperation(&run, "assign", requestID, stored.ID, true); err != nil {
		return Assignment{}, Status{}, err
	}
	return stored.public(), run.status(), nil
}

// Ticket is one parsed ticket file: its title, content digest, charged
// acceptance rows, ownership fence, assumptions, and declared contracts.
type Ticket struct {
	Path, Title              string
	Digest                   string
	Rows, Fence, Assumptions []string
	Contracts                string
}

var ticketRow, packageName, rowRange = regexp.MustCompile(`^\s*-\s+\[[ xX]\]\s+\[([^]]+)\]`), regexp.MustCompile(`\binternal/[A-Za-z0-9_-]+\b`), regexp.MustCompile(`^(R)([0-9]+)-R([0-9]+)$`)

// contractOperand matches one backticked operand of a Contracts crossing.
var contractOperand = regexp.MustCompile("`([^`]+)`")

// ContractsAnchored reports whether the declared crossings name at least one
// path the ticket itself writes. A crossing's far side may name a surface no
// path holds, so one anchored operand is the whole requirement; a crossing
// written entirely in concepts anchors nothing, and the ticket cannot maintain
// the advertisement of what it changes. This is assignment policy, not parse
// validity: ParseTicket stays the shared grammar its other consumers grade
// against, and only the assign path refuses on the answer.
func (t Ticket) ContractsAnchored() bool {
	if t.Contracts == "" || t.Contracts == "none crosses" {
		return true
	}
	for _, match := range contractOperand.FindAllStringSubmatch(t.Contracts, -1) {
		operand := strings.TrimSpace(match[1])
		for _, path := range t.Fence {
			if operand != "" && (strings.HasPrefix(operand, path) || strings.HasPrefix(path, operand)) {
				return true
			}
		}
	}
	return false
}

// ParseTicket resolves arg against specPath's tickets directory and parses the
// ticket file it names. The conformance example-agreement check is the
// cross-package consumer, grading the taught ticket example with the same parse
// assignment runs; the grammar and every refusal below are the assign path's
// own, so a consumer needing another shape changes its input, not this parse.
func ParseTicket(specPath, arg string) (Ticket, error) {
	if arg == "" || filepath.IsAbs(arg) {
		return Ticket{}, errors.New("spec build ticket must name one regular ticket file")
	}
	clean := filepath.Clean(arg)
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return Ticket{}, errors.New("spec build ticket escapes its spec")
	}
	clean = strings.TrimPrefix(clean, "tickets"+string(filepath.Separator))
	path := filepath.Join(filepath.Dir(specPath), "tickets", clean)
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() {
		return Ticket{}, errors.New("spec build ticket must name one regular ticket file")
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return Ticket{}, fmt.Errorf("read spec build ticket: %w", err)
	}
	result := Ticket{Path: path, Digest: digest(string(b))}
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") && result.Title == "" {
			result.Title = strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
		if m := ticketRow.FindStringSubmatch(line); len(m) == 2 {
			result.Rows = append(result.Rows, expandRows(m[1])...)
		}
		if strings.HasPrefix(line, "Ownership fence:") {
			result.Fence = append(result.Fence, listValue(strings.TrimPrefix(line, "Ownership fence:"))...)
		}
		if strings.HasPrefix(line, "Assumptions:") {
			result.Assumptions = append(result.Assumptions, listValue(strings.TrimPrefix(line, "Assumptions:"))...)
		}
		if strings.HasPrefix(line, "Contracts:") && result.Contracts == "" {
			result.Contracts = strings.TrimSpace(strings.TrimPrefix(line, "Contracts:"))
		}
	}
	if len(result.Rows) == 0 {
		return Ticket{}, errors.New("spec build ticket declares no charged rows")
	}
	if len(result.Fence) == 0 {
		result.Fence = packageName.FindAllString(string(b), -1)
	}
	// Acceptance IDs are the per-obligation identity downstream accounting
	// keys on, so a repeat — literal or via R-range overlap — is malformed
	// input, not something to collapse.
	seen := make(map[string]bool, len(result.Rows))
	for _, row := range result.Rows {
		if seen[row] {
			return Ticket{}, fmt.Errorf("spec build ticket %s declares duplicate acceptance ID %s", filepath.Base(path), row)
		}
		seen[row] = true
	}
	result.Fence = unique(result.Fence)
	result.Assumptions = unique(result.Assumptions)
	if len(result.Fence) == 0 {
		return Ticket{}, errors.New("spec build ticket declares no ownership fence")
	}
	if result.Title == "" {
		return Ticket{}, errors.New("spec build ticket has no title")
	}
	return result, nil
}
func expandRows(raw string) []string {
	match := rowRange.FindStringSubmatch(raw)
	if len(match) != 4 {
		return []string{raw}
	}
	from, fromErr := strconv.Atoi(match[2])
	to, toErr := strconv.Atoi(match[3])
	if fromErr != nil || toErr != nil || to < from {
		return []string{raw}
	}
	width := len(match[2])
	rows := make([]string, 0, to-from+1)
	for value := from; value <= to; value++ {
		rows = append(rows, fmt.Sprintf("%s%0*d", match[1], width, value))
	}
	return rows
}
func unique(values []string) []string {
	seen := make(map[string]bool, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}
func listValue(value string) []string {
	var values []string
	for _, item := range strings.Split(value, ",") {
		item = strings.Trim(strings.TrimSpace(item), "`")
		if item != "" {
			values = append(values, item)
		}
	}
	return values
}

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
			subject, err := s.subject(mutationStart, resolved)
			if err != nil || subject.branch != run.Branch || subject.tip != run.Base || subject.specTip != run.SpecTip {
				return Status{}, errors.New("spec build working checkout does not match recorded subject")
			}
			return s.finishStart(ctx, subject.branch, subject.tip, false, &run)
		}
		subject, preconditionErr := s.preconditionsAdvancingEmptyRun(mutationStart, slug, resolved, &run, "", "")
		if run.Terminal {
			abandon, abandoned := s.operation(run, "abandon", "apply")
			if !abandoned || abandon.State != "completed" || preconditionErr == nil {
				return run.status(), nil
			}
			if !errors.Is(preconditionErr, errRecompose) {
				return Status{}, preconditionErr
			}
			return s.startRun(ctx, slug, resolved, subject, retainTerminalAttempt(run), subject.branch+"\x00"+subject.tip)
		}
		if preconditionErr != nil {
			return Status{}, preconditionErr
		}
		return run.status(), nil
	}
	subject, err := s.preconditions(mutationStart, slug, resolved, nil, "", "")
	if err != nil {
		return Status{}, err
	}
	return s.startRun(ctx, slug, resolved, subject, nil, "")
}
func (s *Service) startRun(ctx context.Context, slug, resolved string, subject buildSubject, history []json.RawMessage, attempt string) (Status, error) {
	runID, candidate := runIdentity(resolved, attempt)
	run := record{Version: 1, Slug: slug, Spec: resolved, SpecTip: subject.specTip, Run: runID, Branch: subject.branch, Base: subject.tip, Candidate: candidate, CandidateTip: subject.tip, History: history, Assignments: map[string]assignment{}, Operations: map[string]operation{}}
	if absent, err := refAbsent(s.root, run.Candidate); err != nil {
		return Status{}, err
	} else if !absent {
		return Status{}, errors.New("spec build candidate identity already exists")
	}
	// Every start — fresh, or a restart after a terminal run — compares against the live
	// marker a predecessor left: a sibling build's promotion moves it legitimately, while
	// any tip recorded earlier would read as a conflict.
	if err := s.gate.Bootstrap(ctx, s.root, subject.branch, subject.tip, greenMarker(s.root, subject.branch)); err != nil {
		return Status{}, fmt.Errorf("no exact green evidence: run bench gate --fresh, then retry start: %w", err)
	}
	if _, _, err := s.beginOperation(&run, "start", "run", resolved+"\x00"+subject.branch+"\x00"+subject.tip); err != nil {
		return Status{}, err
	}
	for _, point := range []string{"start/bootstrap", "start/state"} {
		if err := s.faultAt(point); err != nil {
			return run.status(), err
		}
	}
	return s.finishStart(ctx, subject.branch, subject.tip, true, &run)
}

// Promote publishes the exact reviewed candidate only after its prospective tree is green.
func (s *Service) Promote(ctx context.Context, slug string) (Status, error) {
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
	if run.Terminal {
		// A promotion that died between the terminal save and its reclamation leaves the
		// refs stranded with no later caller to notice, so re-entry finishes the job.
		s.reclaimProvisionalRefs(run)
		return run.status(), nil
	}
	owner, ok := s.gate.(PromotionGateOwner)
	if run.PromotionCommit != "" && refAt(s.root, "refs/heads/"+run.Branch, run.PromotionCommit) {
		if !ok || owner == nil {
			return run.status(), errors.New("spec build promotion requires a prospective gate owner")
		}
		if err := s.validatePromotionEvidence(run); err != nil {
			return run.status(), err
		}
		valid, validationErr := owner.Validate(ctx, s.root, run.PromotionTree, run.PromotionEvidence)
		if validationErr != nil || !valid || !promotionCommitAt(s.root, run) {
			return run.status(), errors.New("spec build promotion recovery evidence drifted")
		}
		synchronize, checkoutErr := s.validatePromotionRecoveryCheckout(ctx, run)
		if checkoutErr != nil {
			return run.status(), checkoutErr
		}
		if synchronize {
			if _, err := s.git(ctx, nil, nil, "read-tree", "--reset", "-u", run.PromotionCommit); err != nil {
				return run.status(), err
			}
		}
		if err := updateRef(s.root, "refs/bench/green/"+run.Branch, run.PromotionCommit, run.Base); err != nil && !refAt(s.root, "refs/bench/green/"+run.Branch, run.PromotionCommit) {
			return run.status(), err
		}
		if err := s.finishPromotion(&run); err != nil {
			return Status{}, err
		}
		return run.status(), nil
	}
	// Recomposition discards the review, the assignment releases, and the retained
	// promotion evidence, so it is reached before they are graded: a run mid-repair
	// satisfies none of them by construction, and promote is the operation that
	// resolves that state.
	subject, preconditionErr := s.preconditions(mutationPromote, slug, run.Spec, &run, "", "")
	if preconditionErr != nil {
		if !errors.Is(preconditionErr, errRecompose) {
			return Status{}, preconditionErr
		}
		if err := s.recomposePromotion(ctx, &run, subject); err != nil {
			return run.status(), fmt.Errorf("spec build promotion recomposition refused: %w", err)
		}
		return run.status(), nil
	}
	if run.Review == nil || run.Review.Candidate != run.CandidateTip || run.Review.hasAcceptedFinding() {
		return Status{}, errors.New("spec build promotion requires a current clean review")
	}
	for _, assigned := range run.Assignments {
		if assigned.Integrated == "" || !assigned.Released {
			return Status{}, errors.New("spec build promotion requires every assignment integrated and released")
		}
	}
	if err := s.validatePromotionEvidence(run); err != nil {
		return run.status(), err
	}
	if !ok || owner == nil {
		return Status{}, errors.New("spec build promotion requires a prospective gate owner")
	}
	tree, err := s.prospectiveTree(ctx, run)
	if err != nil {
		return Status{}, fmt.Errorf("construct prospective promotion tree: %w", err)
	}
	if _, _, err := s.beginOperation(&run, "promote", run.CandidateTip, tree); err != nil {
		return Status{}, err
	}
	outcome, err := owner.Execute(ctx, s.root, tree)
	if err != nil {
		return run.status(), fmt.Errorf("authorize prospective promotion: %w", err)
	}
	if !validGateOutcome(outcome) {
		return run.status(), errors.New("spec build prospective gate outcome is incomplete")
	}
	run.PromotionTree, run.PromotionEvidence, run.PromotionDisposition = tree, outcome.Evidence, outcome.Disposition
	if err := s.save(run); err != nil {
		return Status{}, err
	}
	if !outcome.Green {
		return run.status(), fmt.Errorf("spec build prospective gate red: %s", outcome.Disposition)
	}
	commit, err := s.gitOutput(ctx, "commit-tree", tree, "-p", run.Base, "-m", "bench promote run="+run.Run+" candidate="+run.CandidateTip)
	if err != nil {
		return Status{}, fmt.Errorf("create promotion squash: %w", err)
	}
	run.PromotionCommit = commit
	if err := s.save(run); err != nil {
		return Status{}, err
	}
	if err := s.faultAt("promote/commit"); err != nil {
		return run.status(), err
	}
	if err := updateRef(s.root, "refs/heads/"+run.Branch, commit, run.Base); err != nil {
		return run.status(), fmt.Errorf("advance working branch compare-and-swap: %w", err)
	}
	if err := s.faultAt("promote/branch"); err != nil {
		return run.status(), err
	}
	if _, err := s.git(ctx, nil, nil, "read-tree", "--reset", "-u", commit); err != nil {
		return run.status(), err
	}
	if err := updateRef(s.root, "refs/bench/green/"+run.Branch, commit, run.Base); err != nil {
		return run.status(), fmt.Errorf("advance project-green marker: %w", err)
	}
	if err := s.faultAt("promote/green"); err != nil {
		return run.status(), err
	}
	if err := s.finishPromotion(&run); err != nil {
		return Status{}, err
	}
	return run.status(), nil
}

func (s *Service) finishPromotion(run *record) error {
	op, found := s.operation(*run, "promote", run.CandidateTip)
	if !found || op.State != "prepared" || op.Result != "" && op.Result != run.PromotionCommit {
		return errors.New("spec build promotion journal is incomplete")
	}
	op.Result, op.State, run.Terminal = run.PromotionCommit, "completed", true
	run.Operations[operationID("promote", run.CandidateTip)] = op
	if err := s.save(*run); err != nil {
		return err
	}
	s.reclaimProvisionalRefs(*run)
	return nil
}
