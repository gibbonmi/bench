package specbuild

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

// Integrate replays one verified checkpoint onto the exact current candidate tip.
func (s *Service) Integrate(ctx context.Context, slug, assignmentID string) (Status, error) {
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
	if !ok || (assigned.Checkpoint == "" && assigned.CheckpointRef == "") {
		return Status{}, errors.New("spec build assignment has no verified checkpoint")
	}
	if assigned.Checkpoint == "" || assigned.CheckpointRef == "" || !refAt(s.root, assigned.CheckpointRef, assigned.Checkpoint) {
		return s.routeDelegate(run, key, assigned, errors.New("checkpoint attribution drifted"))
	}
	if assigned.Integrated != "" {
		if !refAt(s.root, run.Candidate, assigned.Integrated) {
			return Status{}, errors.New("spec build integrated candidate drifted")
		}
		return s.releaseIntegrated(ctx, run, key, assigned)
	}
	for attempt := 0; attempt != 2; attempt++ {
		candidate, err := refValue(s.root, run.Candidate)
		if err != nil {
			return Status{}, err
		}
		if attempt == 0 && candidate != run.CandidateTip {
			return Status{}, errors.New("spec build candidate drifted before integration")
		}
		if attempt != 0 {
			run.CandidateTip = candidate
		}
		patch, err := s.verifyIntegration(run, assigned)
		if err != nil {
			return s.routeDelegate(run, key, assigned, err)
		}
		tree, err := replayCheckpoint(s.root, candidate, assigned.Base, assigned.Checkpoint, patch)
		if err != nil {
			return s.routeDelegate(run, key, assigned, err)
		}
		commit, err := benchgit.Output("-C", s.root, "commit-tree", tree, "-p", candidate, "-m", "bench integrate run="+run.Run+" assignment="+assigned.ID+" checkpoint="+assigned.Checkpoint)
		if err != nil {
			return Status{}, fmt.Errorf("create integrated candidate commit: %w", err)
		}
		if s.beforeCandidateCAS != nil {
			s.beforeCandidateCAS()
		}
		if err := updateRef(s.root, run.Candidate, commit, candidate); err == nil {
			run.CandidateTip, assigned.Integrated, assigned.DelegatePending, assigned.CleanupPending, assigned.Released = commit, commit, false, true, false
			run.Review = nil
			run.Assignments[key] = assigned
			if err := s.save(run); err != nil {
				return Status{}, err
			}
			return s.releaseIntegrated(ctx, run, key, assigned)
		} else {
			latest, latestErr := refValue(s.root, run.Candidate)
			if latestErr != nil || latest == candidate {
				return Status{}, fmt.Errorf("advance candidate compare-and-swap: %w", err)
			}
		}
		if attempt == 1 {
			return Status{}, errors.New("spec build candidate changed during integration")
		}
	}
	return Status{}, errors.New("spec build candidate integration retry exhausted")
}

func (s *Service) releaseIntegrated(ctx context.Context, run record, key string, assigned assignment) (Status, error) {
	if assigned.Released {
		return run.status(), nil
	}
	if !assigned.CleanupPending {
		assigned.CleanupPending = true
		run.Assignments[key] = assigned
		if err := s.save(run); err != nil {
			return Status{}, err
		}
	}
	releaser, ok := releaseOwnerFrom(s.worktrees)
	if !ok {
		return run.status(), errors.New("spec build integrate requires a release-capable worktree owner")
	}
	if err := releaser.Release(ctx, s.root, assigned.Request, assigned.Path); err != nil {
		return run.status(), fmt.Errorf("release integrated assignment: %w", err)
	}
	assigned.CleanupPending, assigned.Released = false, true
	run.Assignments[key] = assigned
	if err := s.save(run); err != nil {
		return Status{}, err
	}
	return run.status(), nil
}

func releaseOwnerFrom(owner WorktreeOwner) (ReleaseOwner, bool) {
	releaser, ok := owner.(ReleaseOwner)
	if !ok {
		return nil, false
	}
	value := reflect.ValueOf(releaser)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		if value.IsNil() {
			return nil, false
		}
	}
	return releaser, true
}

func (s *Service) routeDelegate(run record, key string, assigned assignment, cause error) (Status, error) {
	assigned.DelegatePending = true
	run.Assignments[key] = assigned
	if err := s.save(run); err != nil {
		return Status{}, err
	}
	return run.status(), fmt.Errorf("spec build integration refused: %w", cause)
}

func (s *Service) verifyIntegration(run record, assigned assignment) ([]byte, error) {
	if assigned.CheckpointTree == "" || assigned.CheckpointPatch == "" {
		return nil, errors.New("checkpoint attribution is incomplete")
	}
	tree, err := benchgit.Output("-C", s.root, "rev-parse", assigned.Checkpoint+"^{tree}")
	if err != nil || tree != assigned.CheckpointTree {
		return nil, errors.New("checkpoint patch drifted")
	}
	ticketArg, err := filepath.Rel(filepath.Join(filepath.Dir(run.Spec), "tickets"), assigned.Ticket)
	if err != nil {
		return nil, errors.New("checkpoint ownership drifted")
	}
	current, err := resolveTicket(run.Spec, ticketArg)
	if err != nil || !sameStrings(current.Fence, assigned.Fence) {
		return nil, errors.New("checkpoint ownership drifted")
	}
	if !sameStrings(current.Assumptions, assigned.Assumptions) {
		return nil, errors.New("checkpoint assumptions changed")
	}
	if current.Digest != assigned.TicketDigest {
		return nil, errors.New("checkpoint ticket drifted")
	}
	paths, err := changedPaths(s.root, assigned.Base, assigned.Checkpoint)
	if err != nil || !insideFence(paths, assigned.Fence) {
		return nil, errors.New("checkpoint ownership drifted")
	}
	patch, err := checkpointPatch(s.root, assigned.Base, assigned.Checkpoint)
	if err != nil || digest(string(patch)) != assigned.CheckpointPatch {
		return nil, errors.New("checkpoint patch drifted")
	}
	return patch, nil
}

func replayCheckpoint(root, candidate, base, checkpoint string, patch []byte) (string, error) {
	candidatePaths, err := changedPaths(root, base, candidate)
	if err != nil {
		return "", err
	}
	checkpointPaths, err := changedPaths(root, base, checkpoint)
	if err != nil {
		return "", err
	}
	if overlappingPaths(candidatePaths, checkpointPaths) {
		return "", errors.New("checkpoint patch overlaps the candidate")
	}
	index, err := os.CreateTemp("", "bench-specbuild-index-*")
	if err != nil {
		return "", err
	}
	name := index.Name()
	if err := index.Close(); err != nil {
		return "", err
	}
	defer os.Remove(name)
	env := append(os.Environ(), "GIT_INDEX_FILE="+name)
	if err := gitRun(root, env, nil, "read-tree", candidate); err != nil {
		return "", err
	}
	if err := gitRun(root, env, bytes.NewReader(patch), "apply", "--cached", "--whitespace=nowarn"); err != nil {
		return "", fmt.Errorf("checkpoint patch conflicts with the candidate: %w", err)
	}
	tree, err := gitOutput(root, env, "write-tree")
	if err != nil {
		return "", err
	}
	if candidate == base {
		expected, err := benchgit.Output("-C", root, "rev-parse", checkpoint+"^{tree}")
		if err != nil || tree != expected {
			return "", errors.New("checkpoint patch is not byte-identical")
		}
	}
	return tree, nil
}

func checkpointPatch(root, base, checkpoint string) ([]byte, error) {
	cmd := exec.Command("git", "-C", root, "diff", "--binary", "--full-index", "--no-ext-diff", base, checkpoint)
	patch, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return patch, nil
}

func overlappingPaths(left, right []string) bool {
	for _, leftPath := range left {
		for _, rightPath := range right {
			if leftPath == rightPath {
				return true
			}
		}
	}
	return false
}

func refValue(root, ref string) (string, error) {
	return benchgit.Output("-C", root, "rev-parse", "--verify", ref+"^{commit}")
}

func gitRun(root string, env []string, input io.Reader, args ...string) error {
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	cmd.Env, cmd.Stdin = env, input
	if output, err := cmd.CombinedOutput(); err != nil {
		return errors.New(strings.TrimSpace(string(output)))
	}
	return nil
}

func gitOutput(root string, env []string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	cmd.Env = env
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.New(strings.TrimSpace(string(output)))
	}
	return strings.TrimSpace(string(output)), nil
}
