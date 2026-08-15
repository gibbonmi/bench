// Package authorization owns prospective gate attribution and opaque recovery proof.
package authorization

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"time"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/gate/greenmarker"
	benchgit "github.com/gibbonmi/bench/internal/git"
)

// Kind attributes the outcome of authorizing an unpublished tree.
type Kind string

const (
	Green          Kind = "green"
	Candidate      Kind = "candidate"
	Inherited      Kind = "inherited"
	Infrastructure Kind = "infrastructure"
)

// Result is the gate owner's projection. Evidence is opaque to lifecycle callers.
type Result struct {
	Kind     Kind
	Evidence string
}

// Bootstrap imports current reusable green proof with an expected prior marker.
func Bootstrap(root, branch, tip, expected string) error {
	if !branchAndTipAt(root, branch, tip) {
		return errors.New("working branch or tip changed")
	}
	if !gate.ComposedGreen(root) {
		return errors.New("exact whole-tree evidence unavailable")
	}
	if !branchAndTipAt(root, branch, tip) {
		return errors.New("working branch or tip changed")
	}
	return AdvanceMarker(context.Background(), root, branch, tip, expected)
}

// CheckMarker recognizes a compatible project-green marker without changing it.
func CheckMarker(_ context.Context, root, branch, destination, expected string) error {
	return checkMarker(root, branch, destination, expected, false)
}

// AdvanceMarker recognizes a compatible project-green marker and advances it atomically.
func AdvanceMarker(_ context.Context, root, branch, destination, expected string) error {
	return checkMarker(root, branch, destination, expected, true)
}

func checkMarker(root, branch, destination, expected string, advance bool) error {
	actual, present, err := greenmarker.Read(root, branch)
	if err != nil {
		return err
	}
	if present && actual == destination {
		return nil
	}
	if !fullCommit(root, destination) {
		return errors.New("project-green destination is not a full commit ID")
	}
	if expected != "" {
		if !fullCommit(root, expected) {
			return errors.New("project-green marker does not match expected prior tip")
		}
		ancestor, err := isAncestor(root, expected, destination)
		if err != nil {
			return fmt.Errorf("check expected project-green marker ancestry: %w", err)
		}
		if !ancestor {
			return errors.New("expected project-green marker is not an ancestor of the tip")
		}
	}
	if present {
		if expected == "" {
			return errors.New("project-green marker conflicts with another tip")
		}
		for _, ancestorOf := range []string{expected, destination} {
			ancestor, err := isAncestor(root, actual, ancestorOf)
			if err != nil {
				return fmt.Errorf("check project-green marker ancestry: %w", err)
			}
			if !ancestor {
				return errors.New("project-green marker conflicts with another tip")
			}
		}
	} else if expected != "" {
		return errors.New("project-green marker does not match expected prior tip")
	}
	if !advance {
		return nil
	}
	return greenmarker.Advance(root, branch, destination, actual)
}

func fullCommit(root, value string) bool {
	resolved, err := benchgit.Output("-C", root, "rev-parse", "--verify", value+"^{commit}")
	return err == nil && resolved == value
}

func isAncestor(root, ancestor, descendant string) (bool, error) {
	_, err := benchgit.Raw("-C", root, "merge-base", "--is-ancestor", ancestor, descendant)
	if err == nil {
		return true, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
		return false, nil
	}
	return false, err
}

func branchAndTipAt(root, branch, tip string) bool {
	current, branchErr := benchgit.Output("-C", root, "symbolic-ref", "--quiet", "--short", "HEAD")
	head, headErr := benchgit.Output("-C", root, "rev-parse", "HEAD^{commit}")
	return branchErr == nil && headErr == nil && current == branch && head == tip
}

// Authorize executes the gate for tree and attributes only a durable deterministic red
// to project bytes. A green inherited subject makes that red candidate-owned; without
// that proof the red is inherited. Operational outcomes never masquerade as either.
func Authorize(ctx context.Context, root, tree string) Result {
	return AuthorizeWithWriters(ctx, root, tree, io.Discard, io.Discard)
}

// AuthorizeWithWriters executes the exact-tree gate and returns its attributed outcome.
func AuthorizeWithWriters(ctx context.Context, root, tree string, stdout, stderr io.Writer) Result {
	inheritedGreen := inheritedSubjectGreen(root)
	execution := gate.ExecuteTree(ctx, root, tree, stdout, stderr)
	inspection := gate.InspectTree(root, tree)
	kind := Infrastructure
	if execution.ActionExit == 0 && inspection.ReusableGreen {
		kind = Green
	} else if execution.GateExit != 0 && execution.ActionExit == execution.GateExit &&
		execution.Inspection.State == gate.Ready && execution.Inspection.Status == "red" {
		kind = Inherited
		if inheritedGreen {
			kind = Candidate
		}
	}
	return Result{Kind: kind, Evidence: evidenceToken(kind, tree, inspection)}
}

// Validate reports whether evidence still names reusable exact green proof for tree.
func Validate(root, tree, evidence string) bool {
	inspection := gate.InspectTree(root, tree)
	return inspection.ReusableGreen && evidence != "" && evidence == evidenceToken(Green, tree, inspection)
}

func inheritedSubjectGreen(root string) bool {
	branch, err := benchgit.Output("-C", root, "symbolic-ref", "--quiet", "--short", "HEAD")
	return err == nil && gate.ValidateProjectGreen(root, branch).ReusableGreen
}

func evidenceToken(kind Kind, tree string, inspection gate.EvidenceInspection) string {
	fact := "bench-prospective-authorization/v1\x00" + string(kind) + "\x00" + tree + "\x00" +
		inspection.Oracle + "\x00" + inspection.RecordedAt.UTC().Format(time.RFC3339Nano) + "\x00" +
		strconv.FormatBool(inspection.ReusableGreen)
	sum := sha256.Sum256([]byte(fact))
	return "v1:" + hex.EncodeToString(sum[:])
}
