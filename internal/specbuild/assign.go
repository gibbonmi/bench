package specbuild

import (
	"context"
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
	ticket, err := resolveTicket(run.Spec, ticketArg)
	if err != nil {
		return Assignment{}, Status{}, err
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
	stored := assignment{ID: owned.ID, Path: owned.Path, Base: run.CandidateTip, Request: requestID, OwnerRequest: digest(requestID), Ticket: ticket.Path, TicketDigest: ticket.Digest, Created: time.Now().UTC().Format(time.RFC3339Nano), Rows: ticket.Rows, Fence: ticket.Fence, Assumptions: ticket.Assumptions}
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

type ticket struct {
	Path, Title              string
	Digest                   string
	Rows, Fence, Assumptions []string
}

var ticketRow = regexp.MustCompile(`^\s*-\s+\[[ xX]\]\s+\[([^]]+)\]`)
var packageName = regexp.MustCompile(`\binternal/[A-Za-z0-9_-]+\b`)
var rowRange = regexp.MustCompile(`^(R)([0-9]+)-R([0-9]+)$`)

func resolveTicket(specPath, arg string) (ticket, error) {
	if arg == "" || filepath.IsAbs(arg) {
		return ticket{}, errors.New("spec build ticket must name one regular ticket file")
	}
	clean := filepath.Clean(arg)
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return ticket{}, errors.New("spec build ticket escapes its spec")
	}
	clean = strings.TrimPrefix(clean, "tickets"+string(filepath.Separator))
	path := filepath.Join(filepath.Dir(specPath), "tickets", clean)
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() {
		return ticket{}, errors.New("spec build ticket must name one regular ticket file")
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return ticket{}, fmt.Errorf("read spec build ticket: %w", err)
	}
	result := ticket{Path: path, Digest: digest(string(b))}
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
	}
	if len(result.Rows) == 0 {
		return ticket{}, errors.New("spec build ticket declares no charged rows")
	}
	if len(result.Fence) == 0 {
		result.Fence = packageName.FindAllString(string(b), -1)
	}
	result.Rows = unique(result.Rows)
	result.Fence = unique(result.Fence)
	result.Assumptions = unique(result.Assumptions)
	if len(result.Fence) == 0 {
		return ticket{}, errors.New("spec build ticket declares no ownership fence")
	}
	if result.Title == "" {
		return ticket{}, errors.New("spec build ticket has no title")
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
	run := record{Version: 1, Slug: slug, Spec: resolved, SpecTip: subject.specTip, Run: digest(resolved), Branch: subject.branch, Base: subject.tip, Candidate: "refs/bench/specbuild/candidate/" + digest(resolved), CandidateTip: subject.tip, Assignments: map[string]assignment{}, Operations: map[string]operation{}}
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
