package specbuild

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

// ReleaseEvidence identifies the durable commits that preserve an assignment payload.
type ReleaseEvidence struct {
	Base, CheckpointRef, Checkpoint, IntegratedRef, Integrated string
}

const operationLimit = 64

type operation struct{ Command, Request, Input, Result, State string }

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
	request := assigned.ID
	if op, found := s.operation(run, "integrate", request); found && op.State == "prepared" && op.Result != "" && assigned.Integrated == "" && refAt(s.root, run.Candidate, op.Result) {
		run.CandidateTip, assigned.Integrated, assigned.DelegatePending, assigned.CleanupPending, assigned.Released = op.Result, op.Result, false, true, false
		run.Review, run.PromotionDisposition, run.Assignments[key] = nil, "", assigned
		if err := s.save(run); err != nil {
			return Status{}, err
		}
		return s.releaseIntegrated(ctx, run, key, assigned)
	}
	if _, err := s.preconditions(mutationIntegrate, slug, run.Spec, &run, assignmentID, ""); err != nil {
		return Status{}, err
	}
	if assigned.Integrated == "" {
		if err := s.requireIntegrationTicket(run, assigned); err != nil {
			return Status{}, err
		}
	}
	if assigned.Integrated != "" {
		if !refAt(s.root, run.Candidate, assigned.Integrated) {
			return Status{}, errors.New("spec build integrated candidate drifted")
		}
		return s.releaseIntegrated(ctx, run, key, assigned)
	}
	op, completed, err := s.beginOperation(&run, "integrate", request, assigned.Checkpoint)
	if err != nil {
		return Status{}, err
	} else if completed {
		return Status{}, errors.New("spec build integration operation is incomplete")
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
		patch, err := s.verifyIntegration(ctx, run, assigned)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return run.status(), err
			}
			return s.routeDelegate(run, key, assigned, err)
		}
		tree, err := s.replayCheckpoint(ctx, candidate, assigned.Base, assigned.Checkpoint, patch)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return run.status(), err
			}
			return s.routeDelegate(run, key, assigned, err)
		}
		commit := op.Result
		if commit != "" && !integratedCommitAt(s.root, commit, candidate, tree, assigned.Checkpoint) {
			return Status{}, errors.New("spec build prepared integration result conflicts with replay")
		}
		if commit == "" {
			commit, err = s.gitOutput(ctx, "commit-tree", tree, "-p", candidate, "-m", "bench integrate run="+run.Run+" assignment="+assigned.ID+" checkpoint="+assigned.Checkpoint)
			if err != nil {
				return Status{}, fmt.Errorf("create integrated candidate commit: %w", err)
			}
			if err := s.recordOperation(&run, "integrate", request, commit, false); err != nil {
				return Status{}, err
			}
			op.Result = commit
		}
		if err := s.faultAt("integrate/commit"); err != nil {
			return run.status(), err
		}
		if s.beforeCandidateCAS != nil {
			s.beforeCandidateCAS()
		}
		if err := updateRef(s.root, run.Candidate, commit, candidate); err == nil {
			if err := s.faultAt("integrate/candidate-cas"); err != nil {
				return run.status(), err
			}
			run.CandidateTip, assigned.Integrated, assigned.DelegatePending, assigned.CleanupPending, assigned.Released = commit, commit, false, true, false
			run.Review, run.PromotionDisposition = nil, ""
			run.Assignments[key] = assigned
			if err := s.save(run); err != nil {
				return Status{}, err
			}
			if err := s.faultAt("integrate/state"); err != nil {
				return run.status(), err
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
		if op, found := s.operation(run, "integrate", assigned.ID); found && op.State == "prepared" {
			if err := s.recordOperation(&run, "integrate", assigned.ID, assigned.Integrated, true); err != nil {
				return Status{}, err
			}
		}
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
	evidence := ReleaseEvidence{Base: assigned.Base, CheckpointRef: assigned.CheckpointRef, Checkpoint: assigned.Checkpoint, IntegratedRef: run.Candidate, Integrated: assigned.Integrated}
	if err := releaser.Release(ctx, s.root, assigned.Request, assigned.Path, evidence); err != nil {
		return run.status(), fmt.Errorf("release integrated assignment: %w", err)
	}
	assigned.CleanupPending, assigned.Released = false, true
	run.Assignments[key] = assigned
	if err := s.save(run); err != nil {
		return Status{}, err
	}
	if err := s.faultAt("integrate/release"); err != nil {
		return run.status(), err
	}
	if op, found := s.operation(run, "integrate", assigned.ID); found && op.State == "prepared" {
		if err := s.recordOperation(&run, "integrate", assigned.ID, assigned.Integrated, true); err != nil {
			return Status{}, err
		}
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
	if op, found := s.operation(run, "integrate", assigned.ID); found && op.State == "prepared" {
		delete(run.Operations, operationID("integrate", assigned.ID))
	}
	if err := s.save(run); err != nil {
		return Status{}, err
	}
	return run.status(), fmt.Errorf("spec build integration refused: %w", cause)
}

func (s *Service) verifyIntegration(ctx context.Context, run record, assigned assignment) ([]byte, error) {
	if assigned.CheckpointTree == "" || assigned.CheckpointPatch == "" {
		return nil, errors.New("checkpoint attribution is incomplete")
	}
	tree, err := s.gitOutput(ctx, "rev-parse", assigned.Checkpoint+"^{tree}")
	if err != nil || tree != assigned.CheckpointTree {
		return nil, errors.New("checkpoint patch drifted")
	}
	if err := s.requireIntegrationTicket(run, assigned); err != nil {
		return nil, err
	}
	paths, err := s.changedPaths(ctx, assigned.Base, assigned.Checkpoint)
	if err != nil || !insideFence(paths, assigned.Fence) {
		return nil, errors.New("checkpoint ownership drifted")
	}
	patch, err := s.checkpointPatch(ctx, assigned.Base, assigned.Checkpoint)
	if err != nil || digest(string(patch)) != assigned.CheckpointPatch {
		return nil, errors.New("checkpoint patch drifted")
	}
	return patch, nil
}

func (s *Service) requireIntegrationTicket(run record, assigned assignment) error {
	current, err := validateIntegrationTicket(run, assigned)
	if err != nil {
		return err
	}
	return s.requireCommittedTicket(current)
}

func validateIntegrationTicket(run record, assigned assignment) (Ticket, error) {
	ticketArg, err := filepath.Rel(filepath.Join(filepath.Dir(run.Spec), "tickets"), assigned.Ticket)
	if err != nil {
		return Ticket{}, errors.New("checkpoint ownership drifted")
	}
	current, err := ParseTicket(run.Spec, ticketArg)
	if err != nil || !sameStrings(current.Fence, assigned.Fence) {
		return Ticket{}, errors.New("checkpoint ownership drifted")
	}
	if !sameStrings(current.Assumptions, assigned.Assumptions) {
		return Ticket{}, errors.New("checkpoint assumptions changed")
	}
	if current.Digest != assigned.TicketDigest {
		return Ticket{}, errors.New("checkpoint ticket drifted")
	}
	return current, nil
}

func (s *Service) replayCheckpoint(ctx context.Context, candidate, base, checkpoint string, patch []byte) (string, error) {
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
	if _, err := s.git(ctx, env, nil, "read-tree", candidate); err != nil {
		return "", err
	}
	if _, err := s.git(ctx, env, bytes.NewReader(patch), "apply", "--cached", "--whitespace=nowarn"); err != nil {
		return "", fmt.Errorf("checkpoint patch conflicts with the candidate: %w", err)
	}
	tree, err := s.git(ctx, env, nil, "write-tree")
	if err != nil {
		return "", err
	}
	tree = strings.TrimSpace(tree)
	if candidate == base {
		expected, err := s.gitOutput(ctx, "rev-parse", checkpoint+"^{tree}")
		if err != nil || tree != expected {
			return "", errors.New("checkpoint patch is not byte-identical")
		}
	}
	return tree, nil
}

func (s *Service) checkpointPatch(ctx context.Context, base, checkpoint string) ([]byte, error) {
	patch, err := s.git(ctx, nil, nil, "diff", "--binary", "--full-index", "--no-ext-diff", base, checkpoint)
	if err != nil {
		return nil, err
	}
	return []byte(patch), nil
}

func (s *Service) changedPaths(ctx context.Context, base, tree string) ([]string, error) {
	output, err := s.git(ctx, nil, nil, "diff-tree", "--no-commit-id", "--name-only", "-r", base, tree)
	if err != nil {
		return nil, err
	}
	if output == "" {
		return nil, nil
	}
	return sortedUnique(strings.Split(strings.TrimSpace(output), "\n")), nil
}

func checkpointPatch(root, base, checkpoint string) ([]byte, error) {
	output, err := (processRunner{}).Output(context.Background(), "git", "-C", root, "diff", "--binary", "--full-index", "--no-ext-diff", base, checkpoint)
	return []byte(output), err
}

func refValue(root, ref string) (string, error) {
	return benchgit.Output("-C", root, "rev-parse", "--verify", ref+"^{commit}")
}

func (s *Service) git(ctx context.Context, env []string, input io.Reader, args ...string) (string, error) {
	output, err := s.runner.Run(ctx, Command{Program: "git", Args: append([]string{"-C", s.root}, args...), Env: env, Input: input})
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return "", err
		}
		return "", errors.New(strings.TrimSpace(output))
	}
	return output, nil
}

func integratedCommitAt(root, commit, parent, tree, checkpoint string) bool {
	if commit == "" || !refAt(root, commit, commit) {
		return false
	}
	gotParent, parentErr := benchgit.Output("-C", root, "rev-parse", commit+"^")
	gotTree, treeErr := benchgit.Output("-C", root, "rev-parse", commit+"^{tree}")
	message, messageErr := benchgit.Output("-C", root, "show", "-s", "--format=%B", commit)
	return parentErr == nil && treeErr == nil && messageErr == nil && gotParent == parent && gotTree == tree && strings.Contains(message, "checkpoint="+checkpoint)
}
