//go:build bench_canary_publication_premature_promotion

package publication

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// canaryRegistry is a minimal in-memory Registry the four publication behavior
// canaries share: it records every call in order and tracks live integrity per
// package, without any HTTP or process boundary. It exists only under these
// canary build tags — it never ships and never runs in a normal `go test
// ./internal/publication` pass.
type canaryRegistry struct {
	live  map[string]string
	calls []string
}

func newCanaryRegistry() *canaryRegistry { return &canaryRegistry{live: map[string]string{}} }

func (r *canaryRegistry) key(name, version string) string { return name + "@" + version }

func (r *canaryRegistry) Publish(ctx context.Context, name, version, tag string, tarball []byte) (string, error) {
	integrity := sriIntegrity(tarball)
	r.live[r.key(name, version)] = integrity
	r.calls = append(r.calls, "PUBLISH "+name)
	return integrity, nil
}

func (r *canaryRegistry) StageSubmit(ctx context.Context, name, version string, tarball []byte) (string, error) {
	r.calls = append(r.calls, "STAGE "+name)
	return "stage-" + name, nil
}

func (r *canaryRegistry) Approve(ctx context.Context, stageID string) error { return nil }

func (r *canaryRegistry) Integrity(ctx context.Context, name, version string) (string, bool, error) {
	integrity, live := r.live[r.key(name, version)]
	return integrity, live, nil
}

func (r *canaryRegistry) TagAdd(ctx context.Context, name, tag, version string) error {
	if _, live := r.live[r.key(name, version)]; !live {
		return fmt.Errorf("package version is not live")
	}
	r.calls = append(r.calls, "TAG-ADD "+name+" "+tag)
	return nil
}

func (r *canaryRegistry) TagRemove(ctx context.Context, name, tag string) error {
	r.calls = append(r.calls, "TAG-RM "+name+" "+tag)
	return nil
}

func (r *canaryRegistry) Deprecate(ctx context.Context, name, version, message string) error {
	r.calls = append(r.calls, "DEPRECATE "+name)
	return nil
}

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

// TestPublicationPrematurePromotionCanary is coverage row 8a's premature-
// promotion canary (M4): it proves RunPromotion still refuses to move any
// dist-tag to "latest" — platform or wrapper — until the complete approved
// set reverifies live, even when every platform package is already live and
// only the wrapper is still pending.
func TestPublicationPrematurePromotionCanary(t *testing.T) {
	const version = "9.9.102"
	scriptsRoot := canaryRoot(t)
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, rel := range []string{"scripts/release-plan.mjs", "scripts/release-plan.json"} {
		data, err := os.ReadFile(filepath.Join(scriptsRoot, rel))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, rel), data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	ordered := canaryArtifacts(t, root, version)
	canaryApprovedSet(t, root, version, ordered)

	releaseIndexSHA256, approved, err := VerifyApprovedSet(root, version)
	if err != nil {
		t.Fatalf("verify approved set: %v", err)
	}

	registry := newCanaryRegistry()
	var wrapperName string
	for _, pkg := range approved {
		if pkg.Kind == "wrapper" {
			wrapperName = pkg.Name
			continue
		}
		// Every platform package is already live with matching integrity...
		registry.live[registry.key(pkg.Name, version)] = pkg.Integrity
	}
	if wrapperName == "" {
		t.Fatal("expected a wrapper package in the approved set")
	}
	// ...but the wrapper is deliberately left unpublished: promote must refuse.

	record := Record{
		SchemaVersion:      RecordSchemaVersion,
		ReleaseIndexSHA256: releaseIndexSHA256,
		Path:               "public",
		Profile:            "public",
		Result:             "in_progress",
	}
	if err := SaveRecord(root, record); err != nil {
		t.Fatalf("seed publication record: %v", err)
	}

	if _, err := RunPromotion(context.Background(), root, version, "public", registry); err == nil {
		t.Fatal("promote succeeded before the complete set reverified")
	}
	for _, call := range registry.calls {
		if strings.HasPrefix(call, "TAG-ADD") {
			t.Fatal("promote issued a latest dist-tag-add before the complete set reverified")
		}
	}
}
