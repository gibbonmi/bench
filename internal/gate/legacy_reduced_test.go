package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// The reduced verdict class is retired: nothing writes one, and its fields are unknown
// to the loader's exact field sets. A legacy on-disk reduced record must fail safely as
// non-reusable — refused as invalid rather than read as a full green — so old evidence
// forces a fresh run instead of crediting phases nobody ran.
func TestLegacyReducedRecordIsRefusedAsNonReusable(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 8, 3, 12, 0, 0, 0, time.UTC)
	legacy := `{"schema":1,"state":"ready","status":"green",` +
		`"tree":"0123456789abcdef0123456789abcdef01234567",` +
		`"oracle":"` + strings.Repeat("a", 64) + `",` +
		`"recorded_at":"` + now.Add(-time.Minute).Format(time.RFC3339) + `",` +
		`"reduced":true,"phases":["conformance"],` +
		`"ancestor":"89abcdef0123456789abcdef0123456789abcdef",` +
		`"ancestor_recorded_at":"` + now.Add(-30*time.Minute).Format(time.RFC3339) + `"}`
	path := filepath.Join(t.TempDir(), "bench-last-gate")
	if err := os.WriteFile(path, []byte(legacy), 0o600); err != nil {
		t.Fatal(err)
	}
	loaded := loadVerdict(path, now)
	if loaded.state != Invalid {
		t.Fatalf("state = %v, want Invalid for a legacy reduced record", loaded.state)
	}
	if loaded.reason != "invalid cache record" {
		t.Fatalf("reason = %q, want the invalid-record refusal", loaded.reason)
	}
}
