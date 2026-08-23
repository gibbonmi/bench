package publication

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// approvedReleaseRoot, repoRoot, releasePlanArtifacts, and writeApprovedSet are
// the one fixture harness every test in this package shares that needs a
// verifiable approved release directory. Sharers include the command-surface
// adapter tests, the record-level ordering test, and any behavior canary
// compiled against this package.
//
// A release root is a scratch directory holding a copy of the release plan,
// plus a full green publish-mode approved set. release-plan.mjs answers for
// the fixture from this copy, never from the working tree.

// repoRoot locates the repository root two directories above this package
// (internal/publication). `go test` always runs with the package directory as
// its working directory, regardless of where the test binary was invoked from.
// This is the one way to reach scripts/release-plan.mjs from here.
func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	return filepath.Dir(filepath.Dir(wd))
}

// approvedReleaseRoot builds a scratch release root for version. It copies
// the release plan out of the repository (release-plan.mjs refuses a
// symlinked plan, so these copies are real files). It then adds the
// approved set that plan names.
func approvedReleaseRoot(t *testing.T, version string) string {
	t.Helper()
	root := t.TempDir()
	source := repoRoot(t)
	if err := os.MkdirAll(filepath.Join(root, "scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"release-plan.mjs", "release-plan.json"} {
		data, err := os.ReadFile(filepath.Join(source, "scripts", name))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, "scripts", name), data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	writeApprovedSet(t, root, version, releasePlanArtifacts(t, root, version))
	return root
}

// releasePlanArtifacts returns the platform and wrapper records release-plan.mjs
// names for version under root.
func releasePlanArtifacts(t *testing.T, root, version string) []ArtifactRecord {
	t.Helper()
	out, err := exec.Command("node", filepath.Join(root, "scripts", "release-plan.mjs"), root, "artifact-records", version).Output()
	if err != nil {
		t.Fatalf("release plan artifact-records: %v", err)
	}
	var all []ArtifactRecord
	if err := json.Unmarshal(out, &all); err != nil {
		t.Fatalf("decode artifact-records: %v", err)
	}
	var packages []ArtifactRecord
	for _, record := range all {
		if record.Kind == "wrapper" || record.Kind == "platform" {
			packages = append(packages, record)
		}
	}
	return packages
}

// writeApprovedSet fabricates the approved release directory (dist/preflight +
// dist/artifacts) at root for a full green publish-mode index, the one shape
// VerifyPublishAuthority and VerifyApprovedSet require.
func writeApprovedSet(t *testing.T, root, version string, records []ArtifactRecord) {
	t.Helper()
	artifactDir := filepath.Join(root, "dist", "artifacts")
	if err := os.MkdirAll(artifactDir, 0o755); err != nil {
		t.Fatal(err)
	}
	var sums strings.Builder
	type indexArtifact struct {
		Name   string `json:"name"`
		SHA256 string `json:"sha256"`
	}
	var artifacts []indexArtifact
	for _, record := range records {
		data := []byte(record.Name + " fixture package\n")
		if err := os.WriteFile(filepath.Join(artifactDir, record.Name), data, 0o644); err != nil {
			t.Fatal(err)
		}
		digest := sha256.Sum256(data)
		digestHex := hex.EncodeToString(digest[:])
		sums.WriteString(digestHex + "  " + record.Name + "\n")
		artifacts = append(artifacts, indexArtifact{Name: record.Name, SHA256: digestHex})
	}
	preflightDir := filepath.Join(root, "dist", "preflight")
	if err := os.MkdirAll(preflightDir, 0o755); err != nil {
		t.Fatal(err)
	}
	index := map[string]any{
		"schema_version": 1,
		"mode":           "publish",
		"scope":          "preflight",
		"profile":        "public",
		"status":         "green",
		"marker":         "test-fixture",
		"artifacts":      artifacts,
	}
	indexData, err := json.Marshal(index)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(preflightDir, "release-index.json"), append(indexData, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(preflightDir, "SHA256SUMS"), []byte(sums.String()), 0o644); err != nil {
		t.Fatal(err)
	}
}
