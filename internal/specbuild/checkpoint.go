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

const receiptVersion = 1

var (
	errInvalidReceipt   = errors.New("invalid spec build receipt")
	errMalformedReceipt = errors.New("malformed spec build receipt framing")
)

type receipt struct {
	Version      int          `json:"version"`
	Run          string       `json:"run"`
	Assignment   string       `json:"assignment"`
	Base         string       `json:"base"`
	Tree         string       `json:"tree"`
	TicketDigest string       `json:"ticket_digest"`
	Rows         []rowReceipt `json:"rows"`
	Checks       []check      `json:"checks"`
	Probe        probe        `json:"probe"`
	Ownership    []string     `json:"ownership"`
	Assumptions  []string     `json:"assumptions"`
}

type rowReceipt struct {
	Row     string `json:"row"`
	Outcome string `json:"outcome"`
}

type check struct {
	Name   string `json:"name"`
	Passed bool   `json:"passed"`
}

type probe struct {
	Producer     string `json:"producer"`
	Assignment   string `json:"assignment"`
	Tree         string `json:"tree"`
	Command      string `json:"command"`
	Exit         int    `json:"exit"`
	OutputDigest string `json:"output_digest"`
	Produced     string `json:"produced"`
}

// Checkpoint validates coordinator evidence and retains one attributed provisional commit.
func (s *Service) Checkpoint(ctx context.Context, slug, assignmentID, evidence string) (Status, error) {
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
	key, assigned, ok := assignmentFor(run, assignmentID)
	if !ok {
		return Status{}, errors.New("spec build assignment does not exist")
	}
	if assigned.Checkpoint != "" {
		return run.status(), nil
	}
	rec, raw, err := readReceipt(evidence)
	if err != nil {
		return Status{}, err
	}
	if err := s.validateReceipt(run, assigned, rec, evidence); err != nil {
		return Status{}, err
	}
	commit, err := benchgit.Output("-C", s.root, "commit-tree", rec.Tree, "-p", assigned.Base, "-m", "bench checkpoint run="+run.Run+" assignment="+assigned.ID)
	if err != nil {
		return Status{}, fmt.Errorf("create checkpoint commit: %w", err)
	}
	ref := "refs/bench/specbuild/checkpoint/" + digest(run.Run+"\x00"+assigned.ID)
	if err := updateRef(s.root, ref, commit, zeroObjectID); err != nil {
		return Status{}, fmt.Errorf("bind checkpoint commit: %w", err)
	}
	patch, err := checkpointPatch(s.root, assigned.Base, commit)
	if err != nil {
		return Status{}, fmt.Errorf("record checkpoint patch: %w", err)
	}
	assigned.Checkpoint, assigned.CheckpointRef, assigned.CheckpointTree, assigned.ReceiptDigest, assigned.CheckpointPatch = commit, ref, rec.Tree, digest(string(raw)), digest(string(patch))
	run.Assignments[key] = assigned
	if err := s.save(run); err != nil {
		return Status{}, err
	}
	return run.status(), nil
}

func readReceipt(path string) (receipt, []byte, error) {
	if path == "" || !filepath.IsAbs(path) {
		return receipt{}, nil, errInvalidReceipt
	}
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() {
		return receipt{}, nil, errInvalidReceipt
	}
	classified := bounds.Classify(path, bounds.ControlRecordLimit)
	if classified.State != bounds.StateParsed {
		return receipt{}, nil, errInvalidReceipt
	}
	var rec receipt
	// Coordinator receipts are persisted records, so their final newline is canonical framing.
	if jsonfile.Decode(classified.Data, &rec) != nil {
		return receipt{}, nil, errMalformedReceipt
	}
	if !validReceipt(rec) {
		return receipt{}, nil, errInvalidReceipt
	}
	return rec, classified.Data, nil
}

func validReceipt(rec receipt) bool {
	if rec.Version != receiptVersion || rec.Run == "" || rec.Assignment == "" || rec.Base == "" || rec.Tree == "" || rec.TicketDigest == "" || len(rec.Rows) == 0 || len(rec.Checks) == 0 || rec.Probe.Producer == "" || rec.Probe.Assignment == "" || rec.Probe.Tree == "" || rec.Probe.Command == "" || rec.Probe.Exit != 0 || rec.Probe.OutputDigest == "" || rec.Probe.Produced == "" {
		return false
	}
	if _, err := time.Parse(time.RFC3339Nano, rec.Probe.Produced); err != nil {
		return false
	}
	return true
}

func (s *Service) validateReceipt(run record, assigned assignment, rec receipt, evidence string) error {
	if sameOrBelow(assigned.Path, evidence) || rec.Run != run.Run || rec.Assignment != assigned.ID || rec.Base != assigned.Base || rec.TicketDigest != assigned.TicketDigest || rec.Tree != rec.Probe.Tree || rec.Probe.Assignment != assigned.ID || rec.Probe.Producer != "coordinator" {
		return errInvalidReceipt
	}
	created, err := time.Parse(time.RFC3339Nano, assigned.Created)
	if err != nil || !timeAfter(rec.Probe.Produced, created) {
		return errInvalidReceipt
	}
	status, err := benchgit.Output("-C", assigned.Path, "status", "--porcelain", "--untracked-files=all")
	if err != nil || status != "" {
		return errInvalidReceipt
	}
	tree, err := benchgit.Output("-C", assigned.Path, "rev-parse", "HEAD^{tree}")
	if err != nil || tree != rec.Tree {
		return errInvalidReceipt
	}
	ticketArg, err := filepath.Rel(filepath.Join(filepath.Dir(run.Spec), "tickets"), assigned.Ticket)
	if err != nil {
		return errInvalidReceipt
	}
	current, err := resolveTicket(run.Spec, ticketArg)
	if err != nil || current.Digest != assigned.TicketDigest || !sameStrings(current.Rows, assigned.Rows) || !sameStrings(current.Fence, assigned.Fence) || !sameStrings(current.Assumptions, assigned.Assumptions) {
		return errInvalidReceipt
	}
	if !sameStrings(receiptRows(rec.Rows), assigned.Rows) || !passedChecks(rec.Checks) || !sameStrings(rec.Assumptions, assumptionDigests(assigned.Assumptions)) {
		return errInvalidReceipt
	}
	changed, err := changedPaths(s.root, assigned.Base, rec.Tree)
	if err != nil || !insideFence(changed, assigned.Fence) || !sameStrings(changed, sortedUnique(rec.Ownership)) {
		return errInvalidReceipt
	}
	return nil
}

func assignmentFor(run record, id string) (string, assignment, bool) {
	for key, assigned := range run.Assignments {
		if assigned.ID == id {
			return key, assigned, true
		}
	}
	return "", assignment{}, false
}

func receiptRows(rows []rowReceipt) []string {
	seen := map[string]bool{}
	values := make([]string, 0, len(rows))
	for _, row := range rows {
		if seen[row.Row] || (row.Outcome != "passed" && row.Outcome != "already-covered" && row.Outcome != "not-tdd-able") {
			return nil
		}
		seen[row.Row] = true
		values = append(values, row.Row)
	}
	return sortedUnique(values)
}

func passedChecks(checks []check) bool {
	for _, check := range checks {
		if check.Name == "" || !check.Passed {
			return false
		}
	}
	return true
}
func assumptionDigests(values []string) []string {
	result := make([]string, len(values))
	for i, value := range values {
		result[i] = digest(value)
	}
	return sortedUnique(result)
}
func timeAfter(value string, minimum time.Time) bool {
	parsed, err := time.Parse(time.RFC3339Nano, value)
	return err == nil && !parsed.Before(minimum)
}
func sameOrBelow(root, path string) bool {
	relative, err := filepath.Rel(root, path)
	return err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}

func changedPaths(root, base, tree string) ([]string, error) {
	output, err := benchgit.Output("-C", root, "diff-tree", "--no-commit-id", "--name-only", "-r", base, tree)
	if err != nil {
		return nil, err
	}
	if output == "" {
		return nil, nil
	}
	return sortedUnique(strings.Split(output, "\n")), nil
}

func insideFence(paths, fence []string) bool {
	if len(paths) == 0 || len(fence) == 0 {
		return false
	}
	for _, path := range paths {
		inside := false
		for _, prefix := range fence {
			if path == prefix || strings.HasPrefix(path, strings.TrimSuffix(prefix, "/")+"/") {
				inside = true
				break
			}
		}
		if !inside {
			return false
		}
	}
	return true
}
func sameStrings(left, right []string) bool {
	left, right = sortedUnique(left), sortedUnique(right)
	return strings.Join(left, "\x00") == strings.Join(right, "\x00")
}
func sortedUnique(values []string) []string {
	values = append([]string(nil), values...)
	for i := range values {
		values[i] = filepath.ToSlash(values[i])
	}
	for i := 0; i < len(values); i++ {
		for j := i + 1; j < len(values); j++ {
			if values[j] < values[i] {
				values[i], values[j] = values[j], values[i]
			}
		}
	}
	return unique(values)
}
