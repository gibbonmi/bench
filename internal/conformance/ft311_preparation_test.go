package conformance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/anchors"
)

const preparedBuildGuidanceDiagnosticPrefix = "FT311 prepared build guidance: "

func preparedBuildGuidanceAnchors() []anchors.Anchor {
	var found []anchors.Anchor
	for _, anchor := range anchors.Entries() {
		if strings.HasPrefix(anchor.Diagnostic, preparedBuildGuidanceDiagnosticPrefix) {
			found = append(found, anchor)
		}
	}
	return found
}

func TestPreparedBuildGuidanceAnchorsBiteIndependently(t *testing.T) {
	anchors := preparedBuildGuidanceAnchors()
	if got, want := len(anchors), 4; got != want {
		t.Fatalf("prepared-build guidance anchor count = %d, want %d", got, want)
	}
	for _, anchor := range anchors {
		t.Run(anchor.File, func(t *testing.T) {
			root := t.TempDir()
			path := filepath.Join(root, filepath.FromSlash(anchor.File))
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(anchor.Needle+"\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			if diags := checkWorkflowAnchors(root); containsDiagnostic(diags, anchor.Diagnostic) {
				t.Fatalf("anchor is red while its clause is present: %s", anchor.Diagnostic)
			}
			if err := os.WriteFile(path, []byte("planted omitted clause\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			if diags := checkWorkflowAnchors(root); !containsDiagnostic(diags, anchor.Diagnostic) {
				t.Fatalf("omitting the clause did not bite with %q", anchor.Diagnostic)
			}
		})
	}
}

const preparedReviewDispatchDiagnosticPrefix = "FT311 prepared review dispatch: "

func preparedReviewDispatchAnchors() []anchors.Anchor {
	var found []anchors.Anchor
	for _, anchor := range anchors.Entries() {
		if strings.HasPrefix(anchor.Diagnostic, preparedReviewDispatchDiagnosticPrefix) {
			found = append(found, anchor)
		}
	}
	return found
}

// TestPreparedReviewGuidance grades the review phase's native-dispatch authority. A
// Require row must bite when its clause leaves the file. A Forbid row must bite when
// the retired same-family CLI route or inline-axis route returns to the file. The two
// directions together are what DP26 asks of the canonical review readers.
func TestPreparedReviewGuidance(t *testing.T) {
	dispatchAnchors := preparedReviewDispatchAnchors()
	if got, want := len(dispatchAnchors), 13; got != want {
		t.Fatalf("prepared-review dispatch anchor count = %d, want %d", got, want)
	}
	var required, forbidden int
	for _, anchor := range dispatchAnchors {
		switch anchor.Kind {
		case anchors.Require:
			required++
		case anchors.Forbid:
			forbidden++
		default:
			t.Errorf("prepared-review dispatch anchor %q uses an unsupported kind", anchor.Diagnostic)
		}
	}
	if required != 11 || forbidden != 2 {
		t.Fatalf("prepared-review dispatch kinds = %d Require and %d Forbid, want 11 and 2", required, forbidden)
	}

	for _, anchor := range dispatchAnchors {
		t.Run(anchor.Diagnostic, func(t *testing.T) {
			conformant, contradictory := anchor.Needle, "planted contradictory route"
			if anchor.Kind == anchors.Forbid {
				conformant, contradictory = contradictory, anchor.Needle
			}
			root := t.TempDir()
			path := filepath.Join(root, filepath.FromSlash(anchor.File))
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(conformant+"\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			if diags := checkWorkflowAnchors(root); containsDiagnostic(diags, anchor.Diagnostic) {
				t.Fatalf("anchor is red while the file conforms: %s", anchor.Diagnostic)
			}
			if err := os.WriteFile(path, []byte(contradictory+"\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			if diags := checkWorkflowAnchors(root); !containsDiagnostic(diags, anchor.Diagnostic) {
				t.Fatalf("the contradictory file did not bite with %q", anchor.Diagnostic)
			}
		})
	}
}

// TestPreparedReviewGuidanceHoldsOnTheLiveTree grades the shipped guidance files. The
// synthetic trees above prove that each row can bite. Only the live tree proves that
// the canonical readers carry the finished rules today.
func TestPreparedReviewGuidanceHoldsOnTheLiveTree(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	diags := checkWorkflowAnchors(root)
	for _, anchor := range preparedReviewDispatchAnchors() {
		if containsDiagnostic(diags, anchor.Diagnostic) {
			t.Errorf("live guidance is not conformant: %s", anchor.Diagnostic)
		}
	}
}
