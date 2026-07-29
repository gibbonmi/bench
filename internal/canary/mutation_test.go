package canary

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/racetests"
)

func TestMaterializeMutationFixtureGeneratesRegisteredRaceTests(t *testing.T) {
	fixture := t.TempDir()
	files := filepath.Join(fixture, filesDirName)
	if err := os.MkdirAll(files, 0o755); err != nil {
		t.Fatalf("create fixture files: %v", err)
	}
	if err := os.WriteFile(filepath.Join(files, raceTestsMarker), nil, 0o644); err != nil {
		t.Fatalf("write race-test marker: %v", err)
	}
	dst := t.TempDir()
	if err := materializeMutationFixture(t.TempDir(), fixture, dst); err != nil {
		t.Fatalf("materialize fixture: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dst, raceTestsMarker)); !os.IsNotExist(err) {
		t.Fatalf("race-test marker remained in subject: %v", err)
	}
	for rel, want := range racetests.SyntheticSources() {
		got, err := os.ReadFile(filepath.Join(dst, filepath.FromSlash(rel)))
		if err != nil {
			t.Fatalf("read generated %s: %v", rel, err)
		}
		if string(got) != want {
			t.Fatalf("generated %s = %q, want %q", rel, got, want)
		}
	}
}
