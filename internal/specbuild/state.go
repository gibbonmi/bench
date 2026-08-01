package specbuild

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/jsonfile"
	"github.com/gibbonmi/bench/internal/spec"
)

type record struct {
	Version              int                   `json:"version"`
	Slug                 string                `json:"slug"`
	Spec                 string                `json:"spec"`
	SpecTip              string                `json:"spec_tip"`
	Run                  string                `json:"run"`
	Branch               string                `json:"branch"`
	Base                 string                `json:"base"`
	Candidate            string                `json:"candidate"`
	CandidateTip         string                `json:"candidate_tip"`
	PromotionTree        string                `json:"promotion_tree,omitempty"`
	PromotionCommit      string                `json:"promotion_commit,omitempty"`
	PromotionEvidence    string                `json:"promotion_evidence,omitempty"`
	PromotionDisposition GateDisposition       `json:"promotion_disposition,omitempty"`
	Terminal             bool                  `json:"terminal,omitempty"`
	History              []json.RawMessage     `json:"history,omitempty"`
	Assignments          map[string]assignment `json:"assignments"`
	Operations           map[string]operation  `json:"operations"`
	Review               *reviewEvidence       `json:"review,omitempty"`
}
type reviewEvidence struct {
	Candidate string       `json:"candidate"`
	Axes      []reviewAxis `json:"axes"`
	Digest    string       `json:"digest"`
}
type reviewAxis struct {
	Axis     string          `json:"axis"`
	Findings []reviewFinding `json:"findings"`
}
type reviewFinding struct {
	ID          string `json:"id"`
	Disposition string `json:"disposition"`
}

const zeroObjectID = "0000000000000000000000000000000000000000"

type assignment struct {
	ID, Path, Base, Request, OwnerRequest, Ticket, TicketDigest, Created string
	Rows, Fence, Assumptions                                             []string
	Checkpoint, CheckpointRef, CheckpointTree, ReceiptDigest             string
	CheckpointPatch, Integrated                                          string
	DelegatePending, CleanupPending, Released                            bool
}

func (a assignment) public() Assignment {
	return Assignment{ID: a.ID, Path: a.Path, Base: a.Base, Rows: append([]string(nil), a.Rows...), Fence: append([]string(nil), a.Fence...), Assumptions: append([]string(nil), a.Assumptions...)}
}
func (r record) status() Status {
	if r.Terminal {
		return Status{Slug: r.Slug, State: "terminal", Subject: r.CandidateTip}
	}
	state, next := "active", "bench spec build assign "+r.Slug
	cleanup := ""
	pending := ""
	for _, assigned := range r.Assignments {
		if assigned.CleanupPending && (cleanup == "" || assigned.ID < cleanup) {
			cleanup = assigned.ID
		}
		if assigned.DelegatePending && (pending == "" || assigned.ID < pending) {
			pending = assigned.ID
		}
	}
	if cleanup != "" {
		next = "release assignment " + cleanup
	} else if pending != "" {
		next = "delegate assignment " + pending
	} else if r.PromotionDisposition != "" {
		next = map[GateDisposition]string{GateCandidate: "delegate candidate gate repair", GateInherited: "diagnose inherited gate", GateInfrastructure: "retry promote", GateCapExhausted: "implementation cap exhausted"}[r.PromotionDisposition]
		if next == "" {
			next = "diagnose gate outcome"
		}
	} else if r.needsReview() {
		next = "bench spec build review " + r.Slug
	} else if r.Review != nil {
		next = "bench spec build promote " + r.Slug
		if r.Review.hasAcceptedFinding() {
			next = "bench spec build assign " + r.Slug
		}
	}
	if cleanup == "" && pending == "" && r.PromotionDisposition == "" {
		for _, op := range r.Operations {
			if op.State == "prepared" {
				next = "resume " + op.Command
				break
			}
		}
	}
	return Status{Slug: r.Slug, State: state, Subject: r.CandidateTip, Next: next}
}
func (e reviewEvidence) hasAcceptedFinding() bool {
	for _, axis := range e.Axes {
		for _, finding := range axis.Findings {
			if finding.Disposition == "accepted" {
				return true
			}
		}
	}
	return false
}
func (r record) needsReview() bool {
	for _, assigned := range r.Assignments {
		if assigned.Integrated != "" {
			return r.Review == nil || r.Review.Candidate != r.CandidateTip
		}
	}
	return false
}
func (s *Service) resolve(slug string) (string, error) {
	if strings.TrimSpace(slug) == "" {
		return "", errors.New("spec build slug is required")
	}
	_, resolved, _, ok, err := spec.Resolve(s.root, slug)
	if err != nil {
		return "", fmt.Errorf("resolve spec: %w", err)
	}
	if !ok {
		return "", errors.New("spec build spec does not exist")
	}
	return resolved, nil
}
func (s *Service) lock(slug string) (func(), error) {
	path, err := s.statePath(slug)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, fmt.Errorf("create spec build state directory: %w", err)
	}
	f, err := os.OpenFile(path+".lock", os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open spec build lock: %w", err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("lock spec build: %w", err)
	}
	return func() { _ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN); _ = f.Close() }, nil
}
func (s *Service) load(slug string) (record, bool, error) {
	path, err := s.statePath(slug)
	if err != nil {
		return record{}, false, err
	}
	b, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return record{}, false, nil
	}
	if err != nil {
		return record{}, false, fmt.Errorf("read spec build state: %w", err)
	}
	var run record
	if err := jsonfile.Decode(b, &run); err != nil || !run.valid(slug) {
		return record{}, false, errors.New("spec build has incomplete prior state")
	}
	return run, true, nil
}
func (s *Service) save(run record) error {
	path, err := s.statePath(run.Slug)
	if err != nil {
		return err
	}
	b, err := json.Marshal(run)
	if err != nil {
		return fmt.Errorf("encode spec build state: %w", err)
	}
	return replaceState(path, append(b, '\n'))
}
func (r record) valid(slug string) bool {
	if !r.validCore(slug) {
		return false
	}
	seen := map[string]bool{r.Run: true}
	for _, raw := range r.History {
		var prior record
		if jsonfile.DecodeDocument(raw, &prior) != nil || len(prior.History) != 0 || !prior.Terminal || !prior.validCore(slug) || prior.Spec != r.Spec || prior.Branch != r.Branch || seen[prior.Run] {
			return false
		}
		seen[prior.Run] = true
	}
	return true
}

func (r record) validCore(slug string) bool {
	validIdentity := canonicalDigest(r.Run) && r.Candidate == candidateIdentity(r.Run)
	if r.Version != 1 || r.Slug != slug || r.Spec == "" || r.SpecTip == "" || !validIdentity || r.Branch == "" || r.Base == "" || r.CandidateTip == "" || r.Assignments == nil || r.Operations == nil || len(r.Operations) > operationLimit {
		return false
	}
	for key, op := range r.Operations {
		if key != operationID(op.Command, op.Request) || op.Command == "" || op.Request == "" || op.Input == "" || (op.State != "prepared" && op.State != "completed") {
			return false
		}
	}
	return true
}

func replaceState(path string, data []byte) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), ".specbuild-*")
	if err != nil {
		return fmt.Errorf("create spec build state replacement: %w", err)
	}
	name := tmp.Name()
	defer os.Remove(name)
	if err := os.Chmod(name, 0o600); err != nil {
		return err
	}
	if _, err := tmp.Write(data); err != nil || tmp.Sync() != nil || tmp.Close() != nil {
		_ = tmp.Close()
		return errors.New("write spec build state")
	}
	if err := os.Rename(name, path); err != nil {
		return fmt.Errorf("replace spec build state: %w", err)
	}
	installed, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open installed spec build state: %w", err)
	}
	if err := installed.Sync(); err != nil {
		_ = installed.Close()
		return fmt.Errorf("sync installed spec build state: %w", err)
	}
	if err := installed.Close(); err != nil {
		return fmt.Errorf("close installed spec build state: %w", err)
	}
	dir, err := os.Open(filepath.Dir(path))
	if err != nil {
		return fmt.Errorf("open spec build state directory: %w", err)
	}
	defer dir.Close()
	if err := dir.Sync(); err != nil {
		return fmt.Errorf("sync spec build state directory: %w", err)
	}
	return nil
}

func (s *Service) statePath(slug string) (string, error) {
	common, err := benchgit.Output("-C", s.root, "rev-parse", "--path-format=absolute", "--git-common-dir")
	if err != nil {
		return "", fmt.Errorf("resolve common Git directory: %w", err)
	}
	return filepath.Join(common, "bench", "specbuild", digest(slug)+".json"), nil
}

func (s *Service) prospectiveTree(ctx context.Context, run record) (string, error) {
	relative, err := filepath.Rel(s.root, run.Spec)
	if err != nil || filepath.IsAbs(relative) {
		return "", errors.New("spec build spec does not belong to working checkout")
	}
	checkout, err := os.MkdirTemp("", "bench-specbuild-promote-")
	if err != nil {
		return "", err
	}
	defer func() {
		_, _ = s.git(context.Background(), nil, nil, "worktree", "remove", "--force", checkout)
		_ = os.RemoveAll(checkout)
	}()
	if _, err := s.git(ctx, nil, nil, "worktree", "add", "--quiet", "--detach", checkout, run.CandidateTip); err != nil {
		return "", err
	}
	if _, err := spec.Flip(checkout, filepath.Join(checkout, relative)); err != nil {
		return "", err
	}
	if _, err := s.runner.Output(ctx, "git", "-C", checkout, "add", "--", relative); err != nil {
		return "", err
	}
	tree, err := s.runner.Output(ctx, "git", "-C", checkout, "write-tree")
	return strings.TrimSpace(tree), err
}

func updateRef(root, ref, new, old string) error {
	args := []string{"-C", root, "update-ref", ref, new}
	if old != "" {
		args = append(args, old)
	}
	if out, err := exec.Command("git", args...).CombinedOutput(); err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}

func refAt(root, ref, want string) bool {
	got, err := benchgit.Output("-C", root, "rev-parse", "--verify", ref+"^{commit}")
	return err == nil && got == want
}

func refAbsent(root, ref string) (bool, error) {
	err := exec.Command("git", "-C", root, "show-ref", "--verify", "--quiet", ref).Run()
	if err == nil {
		return false, nil
	}
	var exit *exec.ExitError
	if errors.As(err, &exit) && exit.ExitCode() == 1 {
		return true, nil
	}
	return false, fmt.Errorf("inspect candidate identity: %w", err)
}

func digest(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func canonicalDigest(value string) bool {
	_, err := hex.DecodeString(value)
	return len(value) == sha256.Size*2 && err == nil && strings.ToLower(value) == value
}

func candidateIdentity(run string) string { return "refs/bench/specbuild/candidate/" + run }

// Status returns the durable compact projection for slug.
func (s *Service) Status(slug string) (Status, error) {
	if _, err := s.resolve(slug); err != nil {
		return Status{}, err
	}
	run, found, err := s.load(slug)
	if err != nil {
		return Status{}, err
	}
	if !found {
		return Status{Slug: slug, State: "empty", Next: "bench spec build start " + slug}, nil
	}
	return run.status(), nil
}

// FullStatus returns the compact projection together with retained provenance.
func (s *Service) FullStatus(slug string) (FullStatus, error) {
	compact, err := s.Status(slug)
	if err != nil || compact.State == "empty" {
		return FullStatus{Status: compact}, err
	}
	run, _, err := s.load(slug)
	if err != nil {
		return FullStatus{}, err
	}
	full := FullStatus{Status: compact, Assignments: make([]RetainedAssignment, 0, len(run.Assignments))}
	for _, assigned := range run.Assignments {
		cleanup := "active"
		if assigned.CleanupPending {
			cleanup = "pending"
		} else if assigned.Released {
			cleanup = "released"
		}
		full.Assignments = append(full.Assignments, RetainedAssignment{ID: assigned.ID, Ticket: assigned.Ticket, TicketDigest: assigned.TicketDigest, Base: assigned.Base, Checkpoint: assigned.Checkpoint, CheckpointRef: assigned.CheckpointRef, Integrated: assigned.Integrated, ReceiptDigest: assigned.ReceiptDigest, Cleanup: cleanup})
	}
	sort.Slice(full.Assignments, func(i, j int) bool { return full.Assignments[i].ID < full.Assignments[j].ID })
	if run.Review != nil {
		full.Review = &RetainedReview{Candidate: run.Review.Candidate, Digest: run.Review.Digest, Axes: publicAxes(run.Review.Axes)}
	}
	return full, nil
}

func publicAxes(axes []reviewAxis) []ReviewAxis {
	result := make([]ReviewAxis, len(axes))
	for index, axis := range axes {
		findings := make([]ReviewFinding, len(axis.Findings))
		for findingIndex, finding := range axis.Findings {
			findings[findingIndex] = ReviewFinding{ID: finding.ID, Disposition: finding.Disposition}
		}
		result[index] = ReviewAxis{Axis: axis.Axis, Findings: findings}
	}
	return result
}
