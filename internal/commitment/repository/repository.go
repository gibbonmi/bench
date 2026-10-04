// Package repository adapts commitment decisions to repository state.
package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/jsonfile"
	"github.com/gibbonmi/bench/internal/roadmap"
	"github.com/gibbonmi/bench/internal/sanitize"
	"github.com/gibbonmi/bench/internal/spec"
)

// Store addresses one repository's tracked policy and local receipts.
type Store struct{ Root string }

// InventoryItem is one source available to an adoption proposal.
type InventoryItem struct {
	Kind, ID, State, Identity string
}

// Policy returns the current default-branch policy and whether it exists.
func (store Store) Policy() (commitment.Policy, bool, error) {
	policy, _, err := store.defaultPolicy()
	if err != nil {
		return commitment.Policy{}, false, err
	}
	if policy == nil {
		return commitment.Policy{}, false, nil
	}
	return *policy, true, nil
}

// Inventory projects roadmap obligations, staged specs, and assignment identities.
func (store Store) Inventory() ([]InventoryItem, error) {
	tree := roadmap.LoadTree(store.Root)
	document, failures, diagnostics := roadmap.ParseDocument(tree, nil, true)
	if len(failures) > 0 || len(diagnostics) > 0 {
		return nil, errors.New("commitment inventory refused: roadmap is structurally untrusted")
	}
	items := make([]InventoryItem, 0, len(document.Rows))
	for _, row := range document.Rows {
		items = append(items, InventoryItem{Kind: "roadmap", ID: row.ID, State: row.SpecStatus, Identity: commitment.Identity([]byte(row.Body))})
	}
	facts, err := spec.Facts(store.Root)
	if err != nil {
		return nil, err
	}
	for _, fact := range facts {
		if fact.Status == "staged" {
			items = append(items, InventoryItem{Kind: "deliverable", ID: fact.Path, State: fact.Status, Identity: fact.RoadmapID})
		}
	}
	ledger, err := intent.Read(store.Root)
	if err != nil {
		return nil, err
	}
	for _, assignment := range ledger.Assignments {
		items = append(items, InventoryItem{Kind: "run", ID: assignment.ID, State: string(assignment.State), Identity: assignment.Request})
	}
	return items, nil
}

// Plan validates and records one proposed policy.
func (store Store) Plan(input []byte) (commitment.Plan, error) {
	proposal, err := commitment.ParseProposal(input)
	if err != nil {
		return commitment.Plan{}, err
	}
	proposed := proposal.Policy
	current, _, err := store.defaultPolicy()
	if err != nil {
		return commitment.Plan{}, err
	}
	if err := store.validateDeliverables(proposed); err != nil {
		return commitment.Plan{}, err
	}
	plan, err := commitment.BuildPlan(current, proposal)
	if err != nil {
		return commitment.Plan{}, err
	}
	if err := store.validateSources(plan.Sources); err != nil {
		return commitment.Plan{}, err
	}
	payload, err := json.Marshal(plan)
	if err != nil {
		return commitment.Plan{}, fmt.Errorf("encode commitment plan: %w", err)
	}
	err = intent.Transact(store.Root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
		if err := verifiedCompletions(ledger, current, plan.Predecessor, proposed); err != nil {
			return ledger, false, err
		}
		if err := store.listedRuns(ledger, plan.Continuations); err != nil {
			return ledger, false, err
		}
		for _, receipt := range ledger.CommitmentReceipts {
			if receipt.Plan != plan.ID {
				continue
			}
			if receipt.Payload != string(payload) {
				return ledger, false, errors.New("commitment plan identity collision")
			}
			return ledger, false, nil
		}
		ledger.CommitmentReceipts = append(ledger.CommitmentReceipts, intent.CommitmentReceipt{Plan: plan.ID, Payload: string(payload)})
		return ledger, true, nil
	}, nil)
	return plan, err
}

// Approve stages the exact policy transition named by planID.
func (store Store) Approve(planID, decision string, delayed, removed []string) (bool, error) {
	if strings.TrimSpace(decision) == "" || !sanitize.LineSafe(decision) {
		return false, errors.New("commitment approval refused: decision must be one nonempty line")
	}
	if hasDuplicates(delayed) || hasDuplicates(removed) {
		return false, errors.New("commitment approval refused: effect operands contain duplicates")
	}
	var changed bool
	var restore func()
	err := intent.Transact(store.Root, intent.StrictRead, func(ledger intent.Ledger) (intent.Ledger, bool, error) {
		for i := range ledger.CommitmentReceipts {
			receipt := &ledger.CommitmentReceipts[i]
			if receipt.Plan != planID {
				continue
			}
			var plan commitment.Plan
			if err := jsonfile.DecodeDocument([]byte(receipt.Payload), &plan); err != nil {
				return ledger, false, fmt.Errorf("commitment receipt %q: %w", planID, err)
			}
			for otherIndex := range ledger.CommitmentReceipts {
				other := ledger.CommitmentReceipts[otherIndex]
				if otherIndex == i || !other.Approved {
					continue
				}
				var approved commitment.Plan
				if err := jsonfile.DecodeDocument([]byte(other.Payload), &approved); err != nil {
					return ledger, false, fmt.Errorf("commitment receipt %q: %w", other.Plan, err)
				}
				if approved.Predecessor == plan.Predecessor && approved.ProposalIdentity != plan.ProposalIdentity {
					return ledger, false, fmt.Errorf("commitment approval refused: predecessor changed from %s to %s", plan.Predecessor, approved.ProposalIdentity)
				}
			}
			current, predecessor, err := store.defaultPolicy()
			if err != nil {
				return ledger, false, err
			}
			if predecessor != plan.Predecessor && !(receipt.Approved && predecessor == plan.ProposalIdentity) {
				return ledger, false, fmt.Errorf("commitment approval refused: predecessor changed from %s to %s", plan.Predecessor, predecessor)
			}
			if !sameSet(delayed, plan.Effects.Delayed) || !sameSet(removed, plan.Effects.Removed) {
				return ledger, false, errors.New("commitment approval refused: delayed or removed operands do not match the plan")
			}
			if err := store.validateSources(plan.Sources); err != nil {
				return ledger, false, err
			}
			if receipt.Approved {
				if receipt.Decision != decision {
					return ledger, false, errors.New("commitment approval refused: replay decision differs")
				}
				return ledger, false, nil
			}
			proposed, err := commitment.Parse(plan.Proposed)
			if err != nil {
				return ledger, false, err
			}
			if err := verifiedCompletions(ledger, current, plan.Predecessor, proposed); err != nil {
				return ledger, false, err
			}
			if err := store.listedRuns(ledger, plan.Continuations); err != nil {
				return ledger, false, err
			}
			restore, err = store.stage(plan)
			if err != nil {
				return ledger, false, err
			}
			receipt.Decision = decision
			receipt.Approved = true
			changed = true
			return withContinuations(ledger, plan.Continuations), true, nil
		}
		return ledger, false, fmt.Errorf("commitment approval refused: plan %q is unknown", planID)
	}, func() {
		if restore != nil {
			restore()
		}
	})
	return changed, err
}

func (store Store) defaultPolicy() (*commitment.Policy, string, error) {
	revision, err := store.sourceRevision()
	if err != nil {
		return nil, "", err
	}
	return store.policyAt(revision)
}

func (store Store) policyAt(branch string) (*commitment.Policy, string, error) {
	const path = commitment.PolicyPath
	listing, err := git.Output("-C", store.Root, "ls-tree", "-z", branch, "--", path)
	if err != nil {
		return nil, "", fmt.Errorf("read commitment policy: %w", err)
	}
	if strings.TrimSuffix(listing, "\x00") == "" {
		return nil, "absent", nil
	}
	data, err := git.ReadTreeFile(store.Root, branch, path)
	if err != nil {
		return nil, "", fmt.Errorf("read commitment policy: %w", err)
	}
	policy, err := commitment.Parse(data)
	if err != nil {
		return nil, "", err
	}
	canonical, err := commitment.Bytes(policy)
	if err != nil {
		return nil, "", err
	}
	return &policy, commitment.Identity(canonical), nil
}

func (store Store) validateSources(sources []commitment.SourceBinding) error {
	revision, err := store.sourceRevision()
	if err != nil {
		return err
	}
	for _, source := range sources {
		if err := store.validateSourceAt(revision, source); err != nil {
			return err
		}
	}
	return nil
}

func (store Store) stage(plan commitment.Plan) (func(), error) {
	policyPath := filepath.Join(store.Root, filepath.FromSlash(commitment.PolicyPath))
	roadmapPath := filepath.Join(store.Root, roadmap.RoadmapFile)
	policyBefore, err := snapshotFile(policyPath, true)
	if err != nil {
		return nil, err
	}
	roadmapBefore, err := snapshotFile(roadmapPath, false)
	if err != nil {
		return nil, err
	}
	policy, err := commitment.Parse(plan.Proposed)
	if err != nil {
		return nil, err
	}
	projected, err := roadmap.ProjectSequence(roadmapBefore.data, commitment.Remaining(policy))
	if err != nil {
		return nil, err
	}
	restore := func() {
		_ = policyBefore.restore()
		_ = roadmapBefore.restore()
	}
	if err := atomicWrite(policyPath, plan.Proposed, 0o644); err != nil {
		return nil, err
	}
	if err := atomicWrite(roadmapPath, projected, roadmapBefore.mode); err != nil {
		restore()
		return nil, err
	}
	return restore, nil
}

type fileSnapshot struct {
	path    string
	data    []byte
	mode    os.FileMode
	present bool
}

func snapshotFile(path string, allowAbsent bool) (fileSnapshot, error) {
	info, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) && allowAbsent {
		return fileSnapshot{path: path, mode: 0o644}, nil
	}
	if err != nil {
		return fileSnapshot{}, fmt.Errorf("read commitment output %s: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return fileSnapshot{}, fmt.Errorf("read commitment output %s: not a regular file", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fileSnapshot{}, fmt.Errorf("read commitment output %s: %w", path, err)
	}
	return fileSnapshot{path: path, data: data, mode: info.Mode().Perm(), present: true}, nil
}

func (snapshot fileSnapshot) restore() error {
	if !snapshot.present {
		if err := os.Remove(snapshot.path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
		return nil
	}
	return atomicWrite(snapshot.path, snapshot.data, snapshot.mode)
}

func atomicWrite(path string, data []byte, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".bench-commitment-*")
	if err != nil {
		return err
	}
	name := tmp.Name()
	defer os.Remove(name)
	if err := tmp.Chmod(mode); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(name, path)
}

func sameSet(left, right []string) bool {
	left = append([]string(nil), left...)
	right = append([]string(nil), right...)
	slices.Sort(left)
	slices.Sort(right)
	return slices.Equal(left, right)
}

func hasDuplicates(values []string) bool {
	seen := map[string]bool{}
	for _, value := range values {
		if seen[value] {
			return true
		}
		seen[value] = true
	}
	return false
}
