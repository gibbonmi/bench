package conformance

import (
	"slices"
	"testing"
)

// TestGuardProjectionOwnerGovernsRgguard plants a private envelope decode in a temporary
// internal/rgguard tree. The registered package-core-guard check must report it, so the rg
// guard decodes its envelope only through the shared owner.
func TestGuardProjectionOwnerGovernsRgguard(t *testing.T) {
	const rel = "internal/rgguard/rgguard.go"
	sources := guardProjectionKit()
	sources[rel] = `package rgguard; import "encoding/json"; func commandOf(data []byte) (e struct{ ToolInput struct{ Command string } }) { _ = json.Unmarshal(data, &e); return e }`
	want := []string{guardProjectionDiagPrefix + rel + ":1: " + guardProjectionEnvelope}
	if got := guardProjectionDiags(t, t.TempDir(), sources); !slices.Equal(got, want) {
		t.Fatalf("want %q, got %q", want, got)
	}
}
