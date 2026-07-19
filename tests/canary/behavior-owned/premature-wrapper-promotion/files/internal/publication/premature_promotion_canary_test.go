//go:build bench_canary_publication_premature_promotion

package publication

import (
	"context"
	"fmt"
	"os"
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

// canaryRoot, canaryArtifacts, and canaryApprovedSet live in
// canary_shared_test.go, shared by all four publication behavior canaries
// through tests/canary/behavior-owned/publication-canary-base.txt's BASE
// include.

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
