package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var (
	workflowActionKey = regexp.MustCompile(`^[[:space:]]*(?:-[[:space:]]*)?uses:[[:space:]]*`)
	workflowActionPin = regexp.MustCompile(`^[[:space:]]*(?:-[[:space:]]*)?uses:[[:space:]]*[^@[:space:]]+@[0-9a-f]{40}(?:[[:space:]]+#.*)?[[:space:]]*$`)
	workflowLocalCall = regexp.MustCompile(`^[[:space:]]*(?:-[[:space:]]*)?uses:[[:space:]]*\./[^@[:space:]]+(?:[[:space:]]+#.*)?[[:space:]]*$`)
)

func checkKitCompliance(kitRoot string) []string {
	var diags []string
	license, err := os.ReadFile(filepath.Join(kitRoot, "LICENSE"))
	if os.IsNotExist(err) {
		diags = append(diags, "kit root LICENSE is missing")
	} else if err != nil || !strings.HasPrefix(string(license), "MIT License\n") {
		diags = append(diags, "kit root LICENSE does not begin with the canonical MIT heading")
	}

	workflowDir := filepath.Join(kitRoot, ".github", "workflows")
	entries, err := os.ReadDir(workflowDir)
	if os.IsNotExist(err) {
		return diags
	}
	if err != nil {
		return append(diags, fmt.Sprintf("workflow action pin scan failed: %v", err))
	}
	for _, entry := range entries {
		ext := filepath.Ext(entry.Name())
		if entry.IsDir() || (ext != ".yml" && ext != ".yaml") {
			continue
		}
		workflow := filepath.Join(workflowDir, entry.Name())
		body, err := os.ReadFile(workflow)
		if err != nil {
			diags = append(diags, fmt.Sprintf("workflow action pin scan could not read %s", slashRel(kitRoot, workflow)))
			continue
		}
		for lineNo, line := range strings.Split(string(body), "\n") {
			if workflowActionKey.MatchString(line) && !workflowActionPin.MatchString(line) && !workflowLocalCall.MatchString(line) {
				diags = append(diags, fmt.Sprintf("workflow action is not pinned to a 40-character commit digest: %s:%d", slashRel(kitRoot, workflow), lineNo+1))
			}
		}
	}
	return diags
}

func TestKitComplianceChecksBite(t *testing.T) {
	kitRoot := t.TempDir()
	workflowDir := filepath.Join(kitRoot, ".github", "workflows")
	if err := os.MkdirAll(workflowDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workflowDir, "release.yml"), []byte("steps:\n  - uses: actions/checkout@v4\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	diags := checkKitCompliance(kitRoot)
	if !containsDiagnostic(diags, "LICENSE is missing") {
		t.Fatalf("missing LICENSE did not bite: %v", diags)
	}
	if !containsDiagnostic(diags, "not pinned to a 40-character commit digest") {
		t.Fatalf("mutable workflow action did not bite: %v", diags)
	}
}

func TestKitComplianceScansWorkflowYAMLKeysWithoutGlobOrSubstringFalsePositives(t *testing.T) {
	parent := t.TempDir()
	kitRoot := filepath.Join(parent, "kit[")
	workflowDir := filepath.Join(kitRoot, ".github", "workflows")
	if err := os.MkdirAll(workflowDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(kitRoot, "LICENSE"), []byte("MIT License\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	body := "steps:\n  # This job uses: OIDC\n  - run: 'echo configuration uses: cache'\n  - uses: actions/checkout@v4\n"
	if err := os.WriteFile(filepath.Join(workflowDir, "release.yaml"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	diags := checkKitCompliance(kitRoot)
	if len(diags) != 1 || !strings.Contains(diags[0], "not pinned to a 40-character commit digest") {
		t.Fatalf("workflow diagnostics = %v, want only the mutable .yaml action", diags)
	}
}

func TestKitComplianceAcceptsLocalReusableWorkflowAndRejectsMutableExternalAction(t *testing.T) {
	kitRoot := t.TempDir()
	workflowDir := filepath.Join(kitRoot, ".github", "workflows")
	if err := os.MkdirAll(workflowDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(kitRoot, "LICENSE"), []byte("MIT License\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	body := "jobs:\n  verify:\n    uses: ./.github/workflows/native-runtime.yml\nsteps:\n  - uses: actions/setup-go@v5\n"
	if err := os.WriteFile(filepath.Join(workflowDir, "release.yml"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	diags := checkKitCompliance(kitRoot)
	if len(diags) != 1 || !strings.Contains(diags[0], "not pinned to a 40-character commit digest") || !strings.Contains(diags[0], "release.yml:5") {
		t.Fatalf("workflow diagnostics = %v, want only the mutable external action", diags)
	}
}
