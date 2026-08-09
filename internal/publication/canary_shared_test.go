//go:build bench_canary_publication_order_bypass || bench_canary_publication_unpublish_attempt || bench_canary_publication_premature_promotion || bench_canary_publication_integrity_mismatch

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

// canaryRoot, canaryArtifacts, and canaryApprovedSet are the fixture-building
// helpers the four publication behavior canaries (order-bypass,
// unpublish-attempt, premature-promotion, integrity-mismatch) share verbatim.
// They live here, in the one file reachable from all four fixtures through
// mutation tests can share the same publication state machine input, so
// the shared setup has a single source instead of four pasted copies. Each
// canary's own build tag above gates this file into that canary's isolated
// `go test -tags=...` subprocess; it never compiles into a normal
// `go test ./internal/publication` pass.

// canaryRoot locates the repo root two directories above this package
// (internal/publication) — `go test` always runs with the package directory
// as its working directory, regardless of where the test binary was invoked
// from, so this is the one way to reach scripts/release-plan.mjs from here.
func canaryRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	return filepath.Dir(filepath.Dir(wd))
}

// canaryArtifacts returns the four platform artifact records followed by the
// wrapper record for version, read from the canonical release-plan.mjs.
func canaryArtifacts(t *testing.T, root, version string) []ArtifactRecord {
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

// canaryApprovedSet fabricates the approved release directory (dist/preflight
// + dist/artifacts) at root for a full green publish-mode index, the one shape
// VerifyPublishAuthority and VerifyApprovedSet require.
func canaryApprovedSet(t *testing.T, root, version string, ordered []ArtifactRecord) {
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
	for _, record := range ordered {
		data := []byte(record.Name + " canary fixture package\n")
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
		"marker":         "canary-fixture",
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
