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
