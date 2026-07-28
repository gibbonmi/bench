//go:build bench_canary_publication_integrity_mismatch

package publication

import (
	"context"
	"os"
	"path/filepath"
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

// canaryRoot, canaryArtifacts, and canaryApprovedSet live in
// canary_shared_test.go, shared by all four publication behavior canaries
// through tests/canary/behavior-owned/publication-canary-base.txt's BASE
// include.

// TestPublicationIntegrityMismatchCanary is coverage row 8a's integrity-
// mismatch canary (M3): it proves RunFirstPublication's resume path still
// refuses to accept an already-live package whose registry integrity does not
// match the approved local tarball.
func TestPublicationIntegrityMismatchCanary(t *testing.T) {
	const version = "9.9.101"
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

	_, approved, err := VerifyApprovedSet(root, version)
	if err != nil {
		t.Fatalf("verify approved set: %v", err)
	}
	registry := newCanaryRegistry()
	var mismatched ApprovedPackage
	for _, pkg := range approved {
		if pkg.Kind == "platform" {
			mismatched = pkg
			break
		}
	}
	if mismatched.Name == "" {
		t.Fatal("expected at least one platform package in the approved set")
	}
	registry.live[registry.key(mismatched.Name, version)] = "sha512-STALE-not-the-approved-tarball-integrity"

	if _, err := RunFirstPublication(context.Background(), root, version, "public", registry); err == nil {
		t.Fatal("release submit passed over an integrity mismatch")
	}
}
