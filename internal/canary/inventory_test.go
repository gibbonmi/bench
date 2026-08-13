package canary

import (
	"strings"
	"testing"
)

func TestFixturesRefusesEmptyInventoryWithInventoryOnlyDiagnostic(t *testing.T) {
	_, err := Fixtures(t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "canary fixture inventory is empty") {
		t.Fatalf("Fixtures(empty) error = %v, want inventory-only empty diagnostic", err)
	}
}
