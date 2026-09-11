package conformance

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/gibbonmi/bench/internal/anchors"
	"github.com/gibbonmi/bench/internal/lines"
)

func TestRetainedWorkflow(t *testing.T) {
	h := NewHarness(t)
	binding := lines.ParseBinding([]byte(h.ReadRootFile(".bench", "lines.env")))
	if got, want := binding.Cell("codex", "mid"), "gpt-5.6-sol"; got != want {
		t.Fatalf("Codex mid binding = %q, want %q", got, want)
	}
	if got, want := binding.Cell("codex", "top"), "gpt-6-astra"; got != want {
		t.Fatalf("Codex top binding = %q, want %q", got, want)
	}
	if got, want := binding.Cell("codex", "cheap"), "gpt-5.6-terra"; got != want {
		t.Fatalf("Codex cheap binding = %q, want %q", got, want)
	}

	wantDiagnostics := []string{
		"retained workflow: craft-spec dropped the one implementation-line template",
		"retained workflow: craft-spec dropped the implementation-line reason factors",
		"retained workflow: craft-spec dropped the harder-chunk annotation",
		"retained workflow: craft-line dropped the conditional review stage default",
		"retained workflow: craft-line dropped the Codex mid-to-top review branch",
		"retained workflow: craft-line dropped the default mid review branch",
		"retained workflow: review phase dropped the conditional review-line reference",
		"retained workflow: project profile dropped the conditional review-line reference",
		"retained workflow: craft-line dropped the user-directed model-switch boundary",
	}
	family := anchorsWithDiagnosticPrefix("retained workflow: ")
	for _, want := range wantDiagnostics {
		if !slices.ContainsFunc(family, func(anchor anchors.Anchor) bool { return anchor.Diagnostic == want }) {
			t.Errorf("retained-workflow anchor is absent: %s", want)
		}
	}
	if t.Failed() {
		return
	}
	runRetainedWorkflowAnchorBites(t, family)
}

func runRetainedWorkflowAnchorBites(t *testing.T, family []anchors.Anchor) {
	t.Helper()
	for _, anchor := range family {
		t.Run(anchor.Diagnostic, func(t *testing.T) {
			root := t.TempDir()
			path := filepath.Join(root, filepath.FromSlash(anchor.File))
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			write := func(text string) {
				t.Helper()
				if anchor.Section != "" {
					text = "# Fixture\n\n## " + anchor.Section + "\n\n" + text
				}
				if err := os.WriteFile(path, []byte(text+"\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			write(anchor.Needle)
			if diags := checkWorkflowAnchors(root); slices.Contains(diags, anchor.Diagnostic) {
				t.Fatalf("anchor is red while its file conforms: %s", anchor.Diagnostic)
			}
			write(plantedAnchorContradiction)
			if diags := checkWorkflowAnchors(root); !slices.Contains(diags, anchor.Diagnostic) {
				t.Fatalf("the contradictory file did not bite with %q", anchor.Diagnostic)
			}
		})
	}
}
