//go:build bench_canary_publication_unpublish_attempt

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
// ./internal/publication` pass. TagRemove enforces the same invariant the real
// offline-registry.mjs fixture does: no tag named after an exact version is a
// legitimate dist-tag — there is no unpublish operation in this port, so any
// such call is rejected and logged distinctly.
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
	if !strings.HasPrefix(tag, "candidate-") {
		r.calls = append(r.calls, "REJECT-UNPUBLISH "+name+" "+tag)
		return fmt.Errorf("unpublish is not supported")
	}
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

// TestPublicationUnpublishAttemptCanary is coverage row 8a's unpublish-attempt
// canary (M5): it proves RunRollback never issues an unpublish-shaped registry
// call (a tag-remove keyed by the exact version rather than the candidate
// tag) in place of deprecation.
func TestPublicationUnpublishAttemptCanary(t *testing.T) {
	const version = "9.9.103"
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

	if _, err := RunRollback(context.Background(), root, version, "public", "recovering from a bad release", registry); err != nil {
		t.Fatalf("rollback failed: %v\ncalls: %v", err, registry.calls)
	}
	for _, call := range registry.calls {
		if strings.HasPrefix(call, "REJECT-UNPUBLISH") {
			t.Fatal("rollback issued an unpublish-shaped request")
		}
	}
}
