package gate

import (
	"slices"
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
	benchgit "github.com/gibbonmi/bench/internal/git"
)

// A surface that exists in neither the declaration nor the repository today: membership
// must follow location, so a file nobody has enumerated is covered the day it lands.
func TestCaptureDirectoryCoversItsSurfaces(t *testing.T) {
	scope := ReducedScope()
	for _, path := range []string{
		"capture/IDEAS.md",
		"capture/newly-added-surface.md",
		"capture/retros/reduced-gate-phase-set.md",
		"specs/newly-added-spec/spec.md",
		"specs/newly-added-spec/tickets/first.md",
	} {
		if !scope.Member(path) {
			t.Errorf("Member(%q) = false, want true", path)
		}
	}
}

// Confinement is what selects the reduced run, so its boundary is where an over-broad
// match costs a real grading. The near-miss siblings share a declared file's prefix and
// must still fall outside; the mixed set is the all-versus-any error.
func TestConfinementRejectsUndeclaredPath(t *testing.T) {
	scope := ReducedScope()
	confined := map[string]bool{
		"ROADMAP.md":                        true,
		".bench-notes.md":                   true,
		"ROADMAP.md.bak":                    false,
		".bench-notes.md.bak":               false,
		"ROADMAP.mdx":                       false,
		"docs/ROADMAP.md":                   false,
		"internal/gate/scope.go":            false,
		"specs/reduced/spec.md":             true,
		"internal/gate/../capture/IDEAS.md": false,
	}
	for path, want := range confined {
		if got := scope.Confines([]string{path}); got != want {
			t.Errorf("Confines([%q]) = %v, want %v", path, got, want)
		}
	}
	if scope.Confines([]string{"ROADMAP.md", "internal/gate/gate.go"}) {
		t.Error("Confines(allowlisted + unlisted) = true, want false: confinement is all, not any")
	}
	// Nothing to reduce against is not a licence to reduce.
	if scope.Confines(nil) {
		t.Error("Confines(nil) = true, want false")
	}
}

// A declared directory is the only non-exact rule in the declaration, so the path
// boundary is where an over-broad match would hide: a sibling that merely starts with
// the directory's name is a different surface entirely.
func TestDeclaredDirectoryMembership(t *testing.T) {
	scope := ReducedScope()
	member := map[string]bool{
		"capture/learnings.md":                       true,
		"capture/retros/reduced-gate-phase-set.md":   true,
		"specs/reduced-gate-phase-set/spec.md":       true,
		"specs/reduced-gate-phase-set/tickets/t.md":  true,
		"capture-old/x.md":                           false,
		"specs-old/x.md":                             false,
		"captured/x.md":                              false,
		"vendor/capture/x.md":                        false,
		"../capture/x.md":                            false,
		"/home/devuser/workspace/bench/capture/x.md": false,
		"capture/../internal/gate/gate.go":           false,
		"capture//doubled.md":                        false,
	}
	for path, want := range member {
		if got := scope.Member(path); got != want {
			t.Errorf("Member(%q) = %v, want %v", path, got, want)
		}
	}
}

// The reduction is only worth having if it skips the phases that cost — an excludable
// set holding one cheap phase satisfies every other rule while saving nothing.
func TestExcludableSetCoversContractPhase(t *testing.T) {
	scope := ReducedScope()
	if !scope.Excludable(canary.PhaseContract) {
		t.Errorf("Excludable(%q) = false, want true", canary.PhaseContract)
	}
	for _, phase := range scope.IncludedPhases() {
		if scope.Excludable(phase) {
			t.Errorf("phase %q is both included and excludable", phase)
		}
	}
	// The build phase produces the binary the other phases exec, so it runs in both
	// modes: excluding it would break a reduced run, including it would claim it grades
	// the declared paths.
	if scope.Excludable(canary.PhaseBuild) {
		t.Errorf("Excludable(%q) = true, want false", canary.PhaseBuild)
	}
	if slices.Contains(scope.IncludedPhases(), canary.PhaseBuild) {
		t.Errorf("IncludedPhases contains %q, want it declared in neither set", canary.PhaseBuild)
	}
}

// The set a reduced run executes is the table complement IncludedPhaseNames derives;
// IncludedPhases is its readable advertisement. Pinning the two equal over the kit's own
// table is what keeps a non-excludable phase from joining every reduced run while the
// accessor — and the profile row bound to the same derivation — advertises the old pair.
func TestIncludedPhasesMatchKitTableDerivation(t *testing.T) {
	root, err := benchgit.Root()
	if err != nil {
		t.Fatalf("resolve kit root: %v", err)
	}
	derived, err := IncludedPhaseNames(root, root)
	if err != nil {
		t.Fatalf("derive included phases from the kit table: %v", err)
	}
	advertised := ReducedScope().IncludedPhases()
	slices.Sort(derived)
	slices.Sort(advertised)
	if !slices.Equal(derived, advertised) {
		t.Fatalf("IncludedPhases() advertises %v, but the kit table derivation runs %v", advertised, derived)
	}
}
