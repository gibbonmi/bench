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
	branch, tip, err := workingSubject(s.root)
	if err != nil || branch != run.Branch || tip != run.Base || !refAt(s.root, run.Candidate, run.CandidateTip) {
		if err == nil {
			err = errors.New("spec build assignment conflicts with the recorded working subject")
		}
		return Assignment{}, Status{}, err
	}
	ticket, err := resolveTicket(run.Spec, ticketArg)
	if err != nil {
		return Assignment{}, Status{}, err
	}
	requestID := digest(request)
	if existing, ok := run.Assignments[requestID]; ok {
		if existing.Ticket != ticket.Path {
			return Assignment{}, Status{}, errors.New("spec build assignment request conflicts with another ticket")
		}
		return existing.public(), run.status(), nil
	}
	if s.worktrees == nil {
		return Assignment{}, Status{}, errors.New("spec build assign requires a worktree owner")
	}
	owned, err := s.worktrees.Create(ctx, s.root, digest(run.Run+"\x00"+ticket.Path+"\x00"+request), ticket.Title, run.CandidateTip)
	if err != nil {
		return Assignment{}, Status{}, fmt.Errorf("create assignment worktree: %w", err)
	}
	if owned.ID == "" || owned.Path == "" {
		return Assignment{}, Status{}, errors.New("worktree owner returned an incomplete assignment")
	}
	stored := assignment{ID: owned.ID, Path: owned.Path, Base: run.CandidateTip, Request: requestID, Ticket: ticket.Path, Rows: ticket.Rows, Fence: ticket.Fence, Assumptions: ticket.Assumptions}
	run.Assignments[requestID] = stored
	if err := s.save(run); err != nil {
		return Assignment{}, Status{}, err
	}
	return stored.public(), run.status(), nil
}

type ticket struct {
	Path, Title              string
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
	result := ticket{Path: path}
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
