package gate

import (
	"path/filepath"
	"reflect"
	"testing"
)

// BG39: `capability` is the failure table's filing name for the skip reporter's
// diagnoses, so a phase of that name would merge its rows with theirs, under a cap that
// applies to only one of the two halves. The manifest is refused before any phase runs,
// which is the last point where the collision is still preventable.
func TestManifestPhaseNamedCapabilityIsRefused(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".bench", "phases.json")
	writeLaneFile(t, path, `{"phases":[{"name":"build","argv":["true"]},{"name":"`+capabilityPhase+`","argv":["true"]}]}`)

	phases, err := phaseTable(root, t.TempDir())
	if err == nil {
		t.Fatalf("phase table = %#v, want the reserved name refused", phases)
	}
	if want := "gate: " + path + ": reserved phase name: " + capabilityPhase; err.Error() != want {
		t.Fatalf("diagnostic = %q, want %q", err.Error(), want)
	}
}

// The reservation is one exact name. A name that differs in case, or that merely contains
// the reserved one, files its rows under itself and stays a legal declaration.
func TestManifestPhaseNamesBesideTheReservedOneAreAccepted(t *testing.T) {
	root := t.TempDir()
	writeLaneFile(t, filepath.Join(root, ".bench", "phases.json"),
		`{"phases":[{"name":"Capability","argv":["true"]},{"name":"capabilities","argv":["true"]},{"name":"build","argv":["true"]}]}`)

	phases, err := phaseTable(root, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"Capability", "capabilities", "build"}; !reflect.DeepEqual(phaseNames(phases), want) {
		t.Fatalf("phase names = %v, want %v", phaseNames(phases), want)
	}
}
