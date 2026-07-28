//go:build bench_canary_publication_order_bypass

package publication

import (
	"context"
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

// TestPublicationOrderBypassCanary is coverage row 8a's order-bypass canary
// (M1): it proves RunFirstPublication still refuses to let the wrapper
// publish ahead of a platform package.
func TestPublicationOrderBypassCanary(t *testing.T) {
	const version = "9.9.100"
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

	registry := newCanaryRegistry()
	if _, err := RunFirstPublication(context.Background(), root, version, "public", registry); err != nil {
		t.Fatalf("first publication failed: %v\ncalls: %v", err, registry.calls)
	}

	wrapperIndex, lastPlatformIndex := -1, -1
	for i, call := range registry.calls {
		if !strings.HasPrefix(call, "PUBLISH ") {
			continue
		}
		name := strings.TrimPrefix(call, "PUBLISH ")
		if name == "redbench" {
			wrapperIndex = i
		} else {
			lastPlatformIndex = i
		}
	}
	if wrapperIndex == -1 || lastPlatformIndex == -1 {
		t.Fatalf("expected both platform and wrapper publish calls: %v", registry.calls)
	}
	if wrapperIndex < lastPlatformIndex {
		t.Fatal("wrapper was published before a platform package")
	}
}
