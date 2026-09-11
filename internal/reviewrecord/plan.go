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
}
type PlannedChunk struct {
	ID           string        `json:"id"`
	Tickets      []string      `json:"tickets"`
	Verification []Requirement `json:"verification"`
	Rows         []string      `json:"-"`
}
type Plan struct {
	Version           int            `json:"version"`
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
	if plan.Version != 1 || len(plan.Chunks) == 0 {
		return plan, errors.New("invalid completion plan version or empty chunks")
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
		if err := requirementsValid(chunk.Verification); err != nil {
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
	if err := requirementsValid(plan.FinalVerification); err != nil {
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

func requirementsValid(items []Requirement) error {
	if len(items) == 0 {
		return errors.New("missing required verification inventory")
	}
	seen := map[string]bool{}
	for _, item := range items {
		if item.ID == "" || item.Command == "" || seen[item.ID] {
			return errors.New("invalid verification requirement")
		}
		seen[item.ID] = true
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
