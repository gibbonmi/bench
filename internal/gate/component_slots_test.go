package gate

// The slot class is what decides whether a gate component may skip, so these tests are
// written as forgery attempts: each one puts something in the store that is not this
// component's slot at this identity and asserts the component runs anyway.
//
// The component family is enumerated from the input registry rather than listed here. A
// scoping decision that held only for the components a fixture happens to carry would leave
// every other component's slot ungraded on the real kit.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/canary"
)

// slotFixtureTime is a fixed authorship instant, so a slot's bytes are reproducible and a
// byte comparison across two authorships is a comparison of what was authored.
var slotFixtureTime = time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)

func slotFixtureNow() time.Time { return slotFixtureTime.Add(time.Hour) }

// componentSlotFixture is a kit-shaped root with every declared component's identity
// resolved against it.
type componentSlotFixture struct {
	root       string
	identities map[string]string
}

func newComponentSlotFixture(t *testing.T) componentSlotFixture {
	t.Helper()
	fixture := newKitShapedFixture(t)
	identities := mustResolveComponentIdentities(t, fixture.root)
	// Two named components carry the per-component assertions below. They are checked here
	// once so a table that stopped materializing either fails as a missing fixture rather
	// than as an assertion passing over slots for empty identities.
	for _, component := range []string{canary.PhaseVet, canary.PhaseTest} {
		if identities[component] == "" {
			t.Fatalf("no identity for %q; resolved %v", component, identities)
		}
	}
	return componentSlotFixture{root: fixture.root, identities: identities}
}

// components returns the resolved component names, sorted, so an assertion over the family
// is an assertion over the registry's answer for this root.
func (f componentSlotFixture) components() []string {
	names := make([]string, 0, len(f.identities))
	for name := range f.identities {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (f componentSlotFixture) slotPath(t *testing.T, component string) string {
	t.Helper()
	return f.slotPathAt(t, component, f.identities[component])
}

// slotPathAt is the store path component's slot occupies at an identity the fixture no
// longer resolves, which is how the identity-move row reaches the slot authored before the
// move.
func (f componentSlotFixture) slotPathAt(t *testing.T, component, identity string) string {
	t.Helper()
	return filepath.Join(commonGitDirOf(t, f.root), "bench-gate-evidence", componentSlotName(component, identity))
}

func (f componentSlotFixture) resolve(t *testing.T, component string) componentSlotInspection {
	t.Helper()
	return resolveComponentSlot(f.root, component, f.identities[component], slotFixtureNow())
}

func (f componentSlotFixture) author(t *testing.T, component string, at time.Time) {
	t.Helper()
	if err := authorComponentSlot(f.root, component, f.identities[component], at); err != nil {
		t.Fatalf("authorComponentSlot(%q) = %v, want a written slot", component, err)
	}
}

// writeSlotFile installs raw bytes at component's slot path, bypassing the author. Every
// refusal test needs a record no author would produce, so the file discipline is applied by
// hand: a regular 0600 file, which is what readStoreRecord requires before the class
// validation this test is actually about is ever reached.
func (f componentSlotFixture) writeSlotFile(t *testing.T, component string, data []byte) {
	t.Helper()
	dir := filepath.Join(commonGitDirOf(t, f.root), "bench-gate-evidence")
	if err := ensureEvidenceDir(filepath.Dir(dir), dir); err != nil {
		t.Fatal(err)
	}
	path := f.slotPath(t, component)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, 0o600); err != nil {
		t.Fatal(err)
	}
}

// slotBytesOf returns each named component's slot bytes, with an absent slot recorded as an
// absent key rather than as empty bytes.
func (f componentSlotFixture) slotBytesOf(t *testing.T, components []string) map[string][]byte {
	t.Helper()
	bytesByComponent := map[string][]byte{}
	for _, component := range components {
		data, err := os.ReadFile(f.slotPath(t, component))
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			t.Fatal(err)
		}
		bytesByComponent[component] = data
	}
	return bytesByComponent
}

// [PC10a-family] The slot class shares only the store-wide field names with the verdict
// classes. This is the enumerated-family join the class has to hold: a slot field that
// collided with a verdict field would make some verdict record readable as a slot, and the
// exact-field-set refusal every test below rests on would stop separating the two.
func TestComponentSlotFieldsAreDisjointFromTheVerdictClasses(t *testing.T) {
	if shared := componentSlotSharesVerdictFields(); len(shared) != 0 {
		t.Fatalf("slot fields %v share %v with the verdict classes %v", componentSlotFields, shared, partialReadyFields)
	}
	// The shared names are shared in fact, not merely permitted: a name listed as shared that
	// no class carries would silently widen what a slot may borrow from a verdict.
	for _, name := range recordSharedFields {
		if !contains(componentSlotFields, name) || !contains(partialReadyFields, name) {
			t.Errorf("%q is listed as shared but is not carried by both classes", name)
		}
	}
}

// [PC9] A slot a run skipped over is byte-identical afterwards, authorship time included.
// Resolution is the only thing a skip does to a slot, so a resolution that rewrote anything
// would move the recorded authorship away from the run that actually graded the component.
func TestSkippedComponentSlotIsByteIdentical(t *testing.T) {
	fixture := newComponentSlotFixture(t)
	fixture.author(t, canary.PhaseVet, slotFixtureTime)
	before := mustRead(t, fixture.slotPath(t, canary.PhaseVet))

	for i := 0; i < 3; i++ {
		inspection := fixture.resolve(t, canary.PhaseVet)
		if !inspection.Skippable {
			t.Fatalf("resolve #%d = %+v, want a skippable slot", i, inspection)
		}
		if !inspection.AuthoredAt.Equal(slotFixtureTime) {
			t.Fatalf("resolve #%d authored at %s, want the authoring run's %s", i, inspection.AuthoredAt, slotFixtureTime)
		}
		if after := mustRead(t, fixture.slotPath(t, canary.PhaseVet)); !bytes.Equal(before, after) {
			t.Fatalf("slot bytes after resolve #%d = %s, want the authored %s", i, after, before)
		}
	}
}

// [PC9b] A slot answers only for the identity it was authored under. This is the class's
// central claim: the address is the content, so the moment a component's declared inputs
// change its slot stops applying and the component runs. Both halves of the addressing carry
// it — the identity is framed into the slot's name, and the record repeats the identity it
// was authored at — because either one alone leaves a component able to inherit a green
// taken against sources it no longer has.
//
// The move is driven through a real edit to a file the component declares rather than
// through a fabricated identity string, so what is pinned is the whole path from changed
// sources to a running component.
func TestSlotDoesNotAnswerForAMovedIdentity(t *testing.T) {
	fixture := newComponentSlotFixture(t)
	const component = canary.PhaseVet
	const declared = "internal/canary/canary.go"
	if paths := componentEntry(t, mustResolveComponentInputs(t, fixture.root), component).Paths(); !slices.Contains(paths, declared) {
		t.Fatalf("%s inputs omit %q, so this row observes nothing: %v", component, declared, paths)
	}

	authored := fixture.identities[component]
	fixture.author(t, component, slotFixtureTime)
	if inspection := fixture.resolve(t, component); !inspection.Skippable {
		t.Fatalf("resolve at the authoring identity = %+v, want it skippable", inspection)
	}
	slotBytes := mustRead(t, fixture.slotPathAt(t, component, authored))

	writeGateTestFile(t, fixture.root, declared,
		"package canary\n\n// Name is the surface this package's own test grades.\nfunc Name() string { return \"canary\" }\n\nvar edited = 1\n", 0o644)
	moved := mustResolveComponentIdentities(t, fixture.root)[component]
	if moved == authored {
		t.Fatalf("%s identity = %s after editing the declared %q, want it to move", component, moved, declared)
	}

	if inspection := resolveComponentSlot(fixture.root, component, moved, slotFixtureNow()); inspection.Skippable {
		t.Fatalf("resolve %q at the moved identity %s = %+v, want a refusal — its slot was authored at %s",
			component, moved, inspection, authored)
	}
	// The refusal reads and rejects; it does not retire what it rejected. The slot authored
	// at the old identity still answers for the tree that produced it, which is what makes an
	// edit and its revert cost one run rather than two.
	if after := mustRead(t, fixture.slotPathAt(t, component, authored)); !bytes.Equal(after, slotBytes) {
		t.Fatalf("the refusal moved the authored slot's bytes: %s, want %s", after, slotBytes)
	}
}

// [PC10a] A record carrying a verdict class's inherited fields is not a slot, and the
// component runs. Reading it as one would let the whole-tree reduced record — which is
// itself evidence about phases it inherited rather than ran — answer for a component.
func TestSlotWithInheritedFieldsRefuses(t *testing.T) {
	fixture := newComponentSlotFixture(t)
	for _, inherited := range []string{"ancestor", "ancestor_recorded_at", "phases", "reduced"} {
		t.Run(inherited, func(t *testing.T) {
			fields := map[string]any{
				"schema":      componentSlotSchema,
				"component":   canary.PhaseVet,
				"identity":    fixture.identities[canary.PhaseVet],
				"authored_at": slotFixtureTime.Format(time.RFC3339),
				inherited:     inheritedFieldValue(inherited),
			}
			fixture.writeSlotFile(t, canary.PhaseVet, mustMarshalSlotFields(t, fields))
			if inspection := fixture.resolve(t, canary.PhaseVet); inspection.Skippable {
				t.Fatalf("resolve with %q present = %+v, want a refusal", inherited, inspection)
			}
		})
	}
}

// inheritedFieldValue returns a well-formed value for a reduced verdict's inherited field,
// so the refusal under test is the field's presence and not a value the decoder choked on.
func inheritedFieldValue(name string) any {
	switch name {
	case "reduced":
		return true
	case "phases":
		return []string{canary.PhaseVet}
	default:
		return "0000000000000000000000000000000000000000"
	}
}

// [PC10b] A record naming a component other than the one the slot resolves for is refused.
// Both halves are checked: the whole of another component's slot copied into this one's
// path, and a record that agrees with the address on identity but names another component —
// the second is what isolates the component comparison, since the identity comparison alone
// already refuses the first.
func TestSlotNamingAnotherComponentRefuses(t *testing.T) {
	fixture := newComponentSlotFixture(t)
	fixture.author(t, canary.PhaseVet, slotFixtureTime)
	vetBytes := mustRead(t, fixture.slotPath(t, canary.PhaseVet))

	t.Run("another components slot copied in", func(t *testing.T) {
		fixture.writeSlotFile(t, canary.PhaseTest, vetBytes)
		if inspection := fixture.resolve(t, canary.PhaseTest); inspection.Skippable {
			t.Fatalf("resolve %q over %q's slot = %+v, want a refusal", canary.PhaseTest, canary.PhaseVet, inspection)
		}
	})

	t.Run("right identity wrong component", func(t *testing.T) {
		fixture.writeSlotFile(t, canary.PhaseTest, mustMarshalSlotFields(t, map[string]any{
			"schema":    componentSlotSchema,
			"component": canary.PhaseVet,
			// The address the record sits at, so nothing but the component name disagrees.
			"identity":    fixture.identities[canary.PhaseTest],
			"authored_at": slotFixtureTime.Format(time.RFC3339),
		}))
		if inspection := fixture.resolve(t, canary.PhaseTest); inspection.Skippable {
			t.Fatalf("resolve %q over a record naming %q = %+v, want a refusal", canary.PhaseTest, canary.PhaseVet, inspection)
		}
	})
}

// [PC10c] A record failing field-set, schema, or time validation is refused and the
// component runs. Each case is a separate way the store could hold something that is nearly
// a slot, and none of them is repaired into one.
func TestMalformedSlotRefuses(t *testing.T) {
	fixture := newComponentSlotFixture(t)
	identity := fixture.identities[canary.PhaseVet]
	valid := map[string]any{
		"schema":      componentSlotSchema,
		"component":   canary.PhaseVet,
		"identity":    identity,
		"authored_at": slotFixtureTime.Format(time.RFC3339),
	}
	// The baseline is asserted skippable first: without it, every refusal below could be
	// coming from the fixture rather than from the mutation the case names.
	fixture.writeSlotFile(t, canary.PhaseVet, mustMarshalSlotFields(t, valid))
	if inspection := fixture.resolve(t, canary.PhaseVet); !inspection.Skippable {
		t.Fatalf("resolve over a well-formed slot = %+v, want it skippable", inspection)
	}

	for _, testCase := range []struct {
		name   string
		mutate func(map[string]any)
	}{
		{"field set missing a name", func(fields map[string]any) { delete(fields, "identity") }},
		{"field set with an extra name", func(fields map[string]any) { fields["tree"] = "abc" }},
		{"schema unset", func(fields map[string]any) { fields["schema"] = 0 }},
		{"schema from another version", func(fields map[string]any) { fields["schema"] = componentSlotSchema + 1 }},
		{"time in the future", func(fields map[string]any) {
			fields["authored_at"] = slotFixtureNow().Add(time.Hour).Format(time.RFC3339)
		}},
		{"time not in the record format", func(fields map[string]any) {
			fields["authored_at"] = slotFixtureTime.Format(time.RFC1123)
		}},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			fields := map[string]any{}
			for name, value := range valid {
				fields[name] = value
			}
			testCase.mutate(fields)
			fixture.writeSlotFile(t, canary.PhaseVet, mustMarshalSlotFields(t, fields))
			if inspection := fixture.resolve(t, canary.PhaseVet); inspection.Skippable {
				t.Fatalf("resolve over a slot with %s = %+v, want a refusal", testCase.name, inspection)
			}
		})
	}
}

// [PC11] Authorship and invalidation are per component, over the whole registry family.
// Authoring one component's slot leaves every other component's bytes untouched, and
// invalidating one leaves the others present — a run that graded vet has graded nothing
// about test, in either direction.
func TestSlotAuthorshipIsPerComponent(t *testing.T) {
	fixture := newComponentSlotFixture(t)
	family := fixture.components()
	if len(family) < 2 {
		t.Fatalf("resolved components = %v, want the registry family", family)
	}

	fixture.author(t, canary.PhaseVet, slotFixtureTime)
	if present := fixture.slotBytesOf(t, family); len(present) != 1 || present[canary.PhaseVet] == nil {
		t.Fatalf("slots present after authoring only %q = %v, want that one alone", canary.PhaseVet, sortedNames(present))
	}
	for _, component := range family {
		if component == canary.PhaseVet {
			continue
		}
		if inspection := fixture.resolve(t, component); inspection.Skippable {
			t.Fatalf("resolve %q after authoring only %q = %+v, want a refusal", component, canary.PhaseVet, inspection)
		}
	}

	for _, component := range family {
		fixture.author(t, component, slotFixtureTime)
	}
	seeded := fixture.slotBytesOf(t, family)
	if len(seeded) != len(family) {
		t.Fatalf("slots present after authoring the family = %v, want all of %v", sortedNames(seeded), family)
	}

	// A re-authorship at a later instant moves its own slot's bytes and no one else's.
	fixture.author(t, canary.PhaseVet, slotFixtureTime.Add(time.Minute))
	reauthored := fixture.slotBytesOf(t, family)
	if bytes.Equal(reauthored[canary.PhaseVet], seeded[canary.PhaseVet]) {
		t.Fatalf("re-authoring %q left its bytes unchanged: %s", canary.PhaseVet, reauthored[canary.PhaseVet])
	}
	assertOtherSlotsUnchanged(t, family, canary.PhaseVet, seeded, reauthored)

	if err := invalidateComponentSlot(fixture.root, canary.PhaseVet, fixture.identities[canary.PhaseVet]); err != nil {
		t.Fatalf("invalidateComponentSlot(%q) = %v", canary.PhaseVet, err)
	}
	remaining := fixture.slotBytesOf(t, family)
	if remaining[canary.PhaseVet] != nil {
		t.Fatalf("%q's slot survived its own invalidation", canary.PhaseVet)
	}
	assertOtherSlotsUnchanged(t, family, canary.PhaseVet, seeded, remaining)
	for _, component := range family {
		if component == canary.PhaseVet {
			continue
		}
		if inspection := fixture.resolve(t, component); !inspection.Skippable {
			t.Fatalf("resolve %q after invalidating %q = %+v, want it still skippable", component, canary.PhaseVet, inspection)
		}
	}
}

func assertOtherSlotsUnchanged(t *testing.T, family []string, changed string, before, after map[string][]byte) {
	t.Helper()
	for _, component := range family {
		if component == changed {
			continue
		}
		if !bytes.Equal(before[component], after[component]) {
			t.Errorf("%q's slot bytes moved when %q was written: %s, want %s", component, changed, after[component], before[component])
		}
	}
}

func mustMarshalSlotFields(t *testing.T, fields map[string]any) []byte {
	t.Helper()
	data, err := json.Marshal(fields)
	if err != nil {
		t.Fatal(err)
	}
	return append(data, '\n')
}

func sortedNames(byComponent map[string][]byte) []string {
	names := make([]string, 0, len(byComponent))
	for name := range byComponent {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
