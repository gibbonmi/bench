// Package authorization owns prospective gate attribution and opaque recovery proof.
package authorization

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/gibbonmi/bench/internal/gate"
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
	tree, err := benchgit.Output("-C", root, "rev-parse", tip+"^{tree}")
	if err != nil {
		return errors.New("working tree unavailable")
	}
	inspection := gate.InspectTree(root, tree)
	if !inspection.ReusableGreen {
		if inspection.Reason == "" {
			inspection.Reason = "exact evidence unavailable"
		}
		return errors.New(inspection.Reason)
	}
	if !branchAndTipAt(root, branch, tip) {
		return errors.New("working branch or tip changed")
	}
	marker := "refs/bench/green/" + branch
	if existing, err := benchgit.Output("-C", root, "rev-parse", "--verify", marker+"^{commit}"); err == nil {
		if existing == tip {
			return nil
		}
		if existing != expected {
			return errors.New("project-green marker conflicts with another tip")
		}
	} else if expected != "" {
		return errors.New("project-green marker does not match expected prior tip")
	}
	if expected != "" {
		if !benchgit.OK("-C", root, "merge-base", "--is-ancestor", expected, tip) {
			return errors.New("expected project-green marker is not an ancestor of the tip")
		}
	} else {
		expected = "0000000000000000000000000000000000000000"
	}
	if _, err := benchgit.Raw("-C", root, "update-ref", marker, tip, expected); err != nil {
		if existing, checkErr := benchgit.Output("-C", root, "rev-parse", "--verify", marker+"^{commit}"); checkErr == nil && existing == tip {
			return nil
		}
		return fmt.Errorf("record project-green marker: %w", err)
	}
	return nil
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
	inheritedGreen := inheritedSubjectGreen(root)
	execution := gate.ExecuteTree(ctx, root, tree, io.Discard, io.Discard)
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
