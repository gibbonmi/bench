package gate

import (
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
)

// TestPhasesForModeAcceptsPhaseTableNames grades the owner filter against the resolved
// table rather than a written-down list of names: a fixture whose family names any phase
// the root actually carries runs that phase alone.
func TestPhasesForModeAcceptsPhaseTableNames(t *testing.T) {
	names := []string{"gofmt", "vet", "test", "race", "conformance-suite", "conformance", "contract"}
	var table []Phase
	for _, name := range names {
		table = append(table, Phase{Name: name})
	}
	table = append(table, Phase{Name: "canary"})

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			t.Setenv(canary.PhaseEnv, name)
			got := phasesForMode(table, innerMode)
			if len(got) != 1 || got[0].Name != name {
				t.Fatalf("phasesForMode ran %v, want only %s", phaseNames(got), name)
			}
		})
	}
}

// TestPhasesForModeFallsBackForAbsentOwner keeps the fallback an owner outside the table
// depends on: filtering to a phase this root has no entry for would run nothing and green
// on it, so the fixture gets the full inner gate instead. The canary phase stays skipped
// either way, since an inner gate that swept fixtures would recurse.
func TestPhasesForModeFallsBackForAbsentOwner(t *testing.T) {
	table := []Phase{{Name: "conformance"}, {Name: "gofmt"}, {Name: "canary"}}
	for _, owner := range []string{"", "shellcheck"} {
		t.Setenv(canary.PhaseEnv, owner)
		got := phaseNames(phasesForMode(table, innerMode))
		want := []string{"conformance", "gofmt"}
		if len(got) != len(want) {
			t.Fatalf("owner %q ran %v, want %v", owner, got, want)
		}
		for idx := range want {
			if got[idx] != want[idx] {
				t.Fatalf("owner %q ran %v, want %v", owner, got, want)
			}
		}
	}
}
