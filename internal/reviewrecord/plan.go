package reviewrecord

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strings"

	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/tickets"
)

type Requirement struct {
	ID      string `json:"id"`
	Command string `json:"command"`
	Probe   string `json:"probe,omitempty"`
	Ticket  string `json:"ticket,omitempty"`
}
type PlannedChunk struct {
	ID           string        `json:"id"`
	Tickets      []string      `json:"tickets"`
	Verification []Requirement `json:"verification"`
	Rows         []string      `json:"-"`
}
type Plan struct {
	Version           int            `json:"version"`
	Execution         *Execution     `json:"execution,omitempty"`
	Chunks            []PlannedChunk `json:"chunks"`
	FinalVerification []Requirement  `json:"final_verification"`
	Digest            string         `json:"-"`
}

func ReadPlan(root, tree, spec string) (Plan, error) {
	var plan Plan
	if _, err := RecordPath(spec); err != nil {
		return plan, err
	}
	if !objectID.MatchString(tree) {
		return plan, errors.New("invalid plan tree")
	}
	data, err := benchgit.ReadTreeFile(root, tree, spec)
	if err != nil {
		return plan, err
	}
	payload, err := fenced(data, "bench-completion-plan")
	if err != nil {
		return plan, fmt.Errorf("missing or invalid completion plan: %w", err)
	}
	if err := decode(payload, &plan); err != nil {
		return plan, err
	}
	if (plan.Version != 1 && plan.Version != 2) || len(plan.Chunks) == 0 {
		return plan, errors.New("invalid completion plan version or empty chunks")
	}
	// The execution declaration decides which form the rest of the plan must
	// take, so it is graded before the per-chunk shape it governs.
	if err := validateExecution(plan); err != nil {
		return plan, err
	}
	inputs := [][]byte{data}
	names := []string{}
	owners := map[string]int{}
	ids := map[string]bool{}
	for i, chunk := range plan.Chunks {
		if chunk.ID == "" || ids[chunk.ID] || len(chunk.Tickets) == 0 {
			return plan, errors.New("invalid duplicate or empty plan chunk")
		}
		ids[chunk.ID] = true
		if err := requirementsValid(chunk.Verification, chunk.Tickets, plan.Version == 2); err != nil {
			return plan, fmt.Errorf("chunk %s: %w", chunk.ID, err)
		}
		for _, name := range chunk.Tickets {
			if !safeRelative(name) || path.Base(name) != name {
				return plan, errors.New("invalid plan ticket path")
			}
			if _, exists := owners[name]; exists {
				return plan, errors.New("duplicate plan ticket")
			}
			names = append(names, name)
			owners[name] = i
		}
	}
	if err := requirementsValid(plan.FinalVerification, nil, plan.Version == 2); err != nil {
		return plan, fmt.Errorf("final verification: %w", err)
	}
	parsed := []tickets.Ticket{}
	for _, name := range names {
		data, err := benchgit.ReadTreeFile(root, tree, path.Join(path.Dir(spec), "tickets", name))
		if err != nil {
			return plan, err
		}
		ticket, diagnostics := tickets.ParseTicket(name, data, names, "")
		if len(diagnostics) != 0 {
			return plan, fmt.Errorf("invalid plan ticket %s: %s", name, strings.Join(diagnostics, "; "))
		}
		for _, blocker := range ticket.Blockers {
			if owners[blocker] > owners[name] {
				return plan, fmt.Errorf("chunk %s precedes dependency %s", plan.Chunks[owners[name]].ID, blocker)
			}
		}
		parsed = append(parsed, ticket)
		inputs = append(inputs, data)
		chunk := &plan.Chunks[owners[name]]
		for _, row := range ticket.Covers {
			if !contains(chunk.Rows, row) {
				chunk.Rows = append(chunk.Rows, row)
			}
		}
	}
	if diagnostics := tickets.Cycles(parsed); len(diagnostics) != 0 {
		return plan, errors.New(strings.Join(diagnostics, "; "))
	}
	framed, _ := json.Marshal(inputs)
	plan.Digest = Digest(framed)
	return plan, nil
}

// requirementsValid grades one verification inventory. A nil owned list marks
// the final inventory, which names no ticket in either version. A version 2
// chunk inventory maps every obligation onto one of its own tickets, and leaves
// no ticket uncovered.
func requirementsValid(items []Requirement, owned []string, delegated bool) error {
	if len(items) == 0 {
		return errors.New("missing required verification inventory")
	}
	seen := map[string]bool{}
	covered := map[string]bool{}
	for _, item := range items {
		if item.ID == "" || item.Command == "" || seen[item.ID] {
			return errors.New("invalid verification requirement")
		}
		seen[item.ID] = true
		if !delegated || owned == nil {
			if item.Ticket != "" {
				return fmt.Errorf("verification %s names a ticket; only a version 2 chunk obligation owns one", item.ID)
			}
			continue
		}
		if !contains(owned, item.Ticket) {
			return fmt.Errorf("verification %s names ticket %q outside this chunk", item.ID, item.Ticket)
		}
		covered[item.Ticket] = true
	}
	for _, name := range owned {
		if delegated && !covered[name] {
			return fmt.Errorf("ticket %s has no verification requirement; every ticket owes one", name)
		}
	}
	return nil
}

func SourceDigest(root, tree, spec string) (string, error) {
	record, err := RecordPath(spec)
	if err != nil {
		return "", err
	}
	if !objectID.MatchString(tree) {
		return "", errors.New("invalid source tree")
	}
	return benchgit.TreeWithoutFile(root, tree, record)
}
