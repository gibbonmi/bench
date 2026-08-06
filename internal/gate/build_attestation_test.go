package gate

// The attestation is what stops a seal from vouching for itself, so these tests are written
// as planting attempts: each one leaves a binary and a seal that agree with each other and
// asserts that the store still refuses, or leaves the store holding something no gate build
// authored and asserts the same.
//
// The record classes are enumerated from the store's own registry rather than listed here. A
// disjointness that held only for the classes a test happened to name would leave the next
// class free to collide with one of them.

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	benchfreshness "github.com/gibbonmi/bench/internal/freshness"
)

// attestationFixtureTime is a fixed authorship instant, so an attestation's bytes are
// reproducible and a byte comparison across two authorships compares what was authored.
var attestationFixtureTime = time.Date(2026, 3, 4, 5, 6, 7, 0, time.UTC)

func attestationFixtureNow() time.Time { return attestationFixtureTime.Add(time.Hour) }

// buildFixtureBinaryTo compiles a fixture root's pkg to staged, the staging path Publish
// consumes. Nothing publishes it — a caller that wants the binary in place goes through
// buildAndPublishAt, the only route here that also leaves a matching seal.
func buildFixtureBinaryTo(t *testing.T, root, pkg, staged string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(staged), 0o755); err != nil {
		t.Fatal(err)
	}
	build := exec.Command("go", "build", "-buildvcs=false", "-o", staged, pkg)
	build.Dir = root
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build %s in the fixture: %v\n%s", pkg, err, out)
	}
}

// attestationFixture is a kit-shaped root together with the dist/bench artifact its gate
// builds publish. A root can carry more than one artifact, so every operation below takes
// the artifact it acts on and the fixture's own is only the default.
type attestationFixture struct {
	root       string
	executable string
}

func newAttestationFixture(t *testing.T) attestationFixture {
	t.Helper()
	fixture := newKitShapedFixture(t)
	return attestationFixture{root: fixture.root, executable: fixture.binaryPath()}
}

// buildAndPublish runs a build and leaves the artifact attested and sealed, as a green gate
// build's build phase does.
func (f attestationFixture) buildAndPublish(t *testing.T, authoredAt time.Time) {
	t.Helper()
	f.buildAndPublishAt(t, "./cmd/bench", f.executable, authoredAt)
}

// buildAndPublishAt attests before it publishes, so the digest every assertion below reads
// cannot have come from a seal: the seal describing those bytes does not exist until Publish
// writes it. The digest is hashed from the staged bytes for the same reason.
func (f attestationFixture) buildAndPublishAt(t *testing.T, pkg, executable string, authoredAt time.Time) {
	t.Helper()
	staged := executable + ".staged"
	buildFixtureBinaryTo(t, f.root, pkg, staged)
	digest, err := benchfreshness.ExecutableDigest(staged)
	if err != nil {
		t.Fatalf("benchfreshness.ExecutableDigest(%s) = %v, want the staged bytes hashed", staged, err)
	}
	if err := authorBuildAttestation(f.root, executable, digest, authoredAt); err != nil {
		t.Fatalf("authorBuildAttestation(%s) = %v, want the staged bytes attested", executable, err)
	}
	if err := benchfreshness.Publish(f.root, staged, executable); err != nil {
		t.Fatalf("benchfreshness.Publish(%s) = %v, want the attested binary published", executable, err)
	}
}

func (f attestationFixture) verify(t *testing.T) buildAttestationInspection {
	t.Helper()
	return f.verifyAt(t, f.executable)
}

func (f attestationFixture) verifyAt(t *testing.T, executable string) buildAttestationInspection {
	t.Helper()
	return verifyBuildAttestation(f.root, executable, attestationFixtureNow())
}

func (f attestationFixture) storeDir(t *testing.T) string {
	t.Helper()
	return filepath.Join(commonGitDirOf(t, f.root), "bench-gate-evidence")
}

func (f attestationFixture) attestationPath(t *testing.T) string {
	t.Helper()
	return f.attestationPathOf(t, f.executable)
}

func (f attestationFixture) attestationPathOf(t *testing.T, executable string) string {
	t.Helper()
	name, err := buildAttestationName(executable)
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join(f.storeDir(t), name)
}

// writeAttestationFile installs raw bytes at the attestation's address with the file
// discipline every reader requires, bypassing the author. A refusal test needs a record no
// author would produce, and a laxer file would be refused before the class validation the
// test is about is ever reached.
func (f attestationFixture) writeAttestationFile(t *testing.T, body string) {
	t.Helper()
	path := f.attestationPath(t)
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, 0o600); err != nil {
		t.Fatal(err)
	}
}

// editFixtureMain rewrites the fixture's main package so a rebuild produces different bytes.
func editFixtureMain(t *testing.T, root, marker string) {
	t.Helper()
	writeGateTestFile(t, root, "cmd/bench/main.go",
		"package main\n\nvar built = \""+marker+"\"\n\nfunc main() { _ = built }\n", 0o644)
}

// [PS20] A green build attests the binary it produced, and that attestation verifies against
// the seal the same build published.
//
// The rebuild follows an edit to the main package, so the seal already beside dist/bench
// names a different binary at the moment the build runs. That is what separates the two
// possible sources for the attested digest: taken from the produced bytes it names the new
// binary and agrees with the new seal, taken from a seal it names the previous one.
func TestGreenBuildAttestsItsOwnBinary(t *testing.T) {
	fixture := newAttestationFixture(t)
	if inspection := fixture.verify(t); inspection.Attested {
		t.Fatalf("verify before any gate build = %+v, want a refusal — the fixture seal is not attested", inspection)
	}
	editFixtureMain(t, fixture.root, "green-build")

	fixture.buildAndPublish(t, attestationFixtureTime)

	inspection := fixture.verify(t)
	if !inspection.Attested {
		t.Fatalf("verify after a green build = %+v, want the build's own attestation to hold", inspection)
	}
	if !inspection.AuthoredAt.Equal(attestationFixtureTime) {
		t.Fatalf("attestation authored at %s, want the build's %s", inspection.AuthoredAt, attestationFixtureTime)
	}
	// The seal and the binary agree too, so the attestation is an addition to the seal
	// contract rather than a replacement for any part of it.
	if err := benchfreshness.Verify(fixture.root, fixture.executable); err != nil {
		t.Fatalf("benchfreshness.Verify after a green build = %v, want the published binary to verify", err)
	}
}

// [PC5] A planted binary published with its own recomputed seal fails attestation. Every
// digest agrees with every other — the seal answers for the planted bytes and for the sources
// on disk — and the store still refuses, because nothing in the plant was authored by a gate
// build. This is the whole reason the class exists: seal verification alone cannot tell a
// gate's binary from anyone else's.
func TestPlantedSealFailsAttestation(t *testing.T) {
	fixture := newAttestationFixture(t)
	fixture.buildAndPublish(t, attestationFixtureTime)
	if inspection := fixture.verify(t); !inspection.Attested {
		t.Fatalf("verify after the gate build = %+v, want it attested before anything is planted", inspection)
	}

	// The plant recomputes the seal exactly as a publisher would, which is what anyone able
	// to write beside dist/bench can do.
	editFixtureMain(t, fixture.root, "planted")
	planted := fixture.executable + ".planted"
	buildFixtureBinaryTo(t, fixture.root, "./cmd/bench", planted)
	if err := benchfreshness.Publish(fixture.root, planted, fixture.executable); err != nil {
		t.Fatalf("plant the binary: %v", err)
	}
	if err := benchfreshness.Verify(fixture.root, fixture.executable); err != nil {
		t.Fatalf("benchfreshness.Verify over the plant = %v, want it to pass — the plant is self-consistent, "+
			"so a refusal here would mean this row never reaches the attestation", err)
	}

	if inspection := fixture.verify(t); inspection.Attested {
		t.Fatalf("verify over a planted binary with a recomputed seal = %+v, want a refusal", inspection)
	}
}

// [PS21] An attestation the store cannot produce on demand is not evidence. Each case is a
// different way the store could hold something that is nearly an attestation, and none of
// them is read as one or repaired into one.
func TestAttestationRefusals(t *testing.T) {
	for _, testCase := range []struct {
		name    string
		disturb func(*testing.T, attestationFixture)
	}{
		{"absent", func(t *testing.T, f attestationFixture) {
			if err := os.Remove(f.attestationPath(t)); err != nil {
				t.Fatal(err)
			}
		}},
		{"names another binary", func(t *testing.T, f attestationFixture) {
			f.writeAttestationFile(t, `{"schema":1,"executable":"`+strings.Repeat("a", 64)+
				`","authored_at":"`+attestationFixtureTime.Format(time.RFC3339)+`"}`)
		}},
		{"malformed", func(t *testing.T, f attestationFixture) {
			// A slot's own field carried alongside the attestation's: the exact field set is
			// what keeps one class from being read as the other.
			f.writeAttestationFile(t, `{"schema":1,"executable":"`+strings.Repeat("b", 64)+
				`","authored_at":"`+attestationFixtureTime.Format(time.RFC3339)+`","component":"build"}`)
		}},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			fixture := newAttestationFixture(t)
			fixture.buildAndPublish(t, attestationFixtureTime)
			// The baseline is asserted first: without it, the refusal below could be coming
			// from the fixture rather than from the disturbance the case names.
			if inspection := fixture.verify(t); !inspection.Attested {
				t.Fatalf("verify after the gate build = %+v, want it attested before the disturbance", inspection)
			}

			testCase.disturb(t, fixture)

			if inspection := fixture.verify(t); inspection.Attested {
				t.Fatalf("verify with the attestation %s = %+v, want a refusal", testCase.name, inspection)
			}
		})
	}
}

// [PS22] Authoring an attestation replaces the previous record rather than editing it in
// place, and touches nothing else in the store.
//
// The address is seeded with a file no store writer would leave — readable past the store's
// own permissions — so a writer that truncated the existing file would inherit that laxer
// file and the record would stop being readable at all. Every declared component's slot is
// authored first, from the registry rather than from a list here, so a component the registry
// gains is a component this row covers.
func TestAttestationDoesNotDisturbSlots(t *testing.T) {
	fixture := newAttestationFixture(t)
	identities := mustResolveComponentIdentities(t, fixture.root)
	family := make([]string, 0, len(identities))
	for component := range identities {
		family = append(family, component)
	}
	sort.Strings(family)
	if len(family) < 2 {
		t.Fatalf("resolved components = %v, want the registry family", family)
	}

	slotPaths := map[string]string{}
	for _, component := range family {
		if err := authorComponentSlot(fixture.root, component, identities[component], attestationFixtureTime); err != nil {
			t.Fatalf("authorComponentSlot(%q) = %v", component, err)
		}
		slotPaths[component] = filepath.Join(fixture.storeDir(t), componentSlotName(component, identities[component]))
	}
	slotBytes := func() map[string][]byte {
		bytesByComponent := map[string][]byte{}
		for component, path := range slotPaths {
			bytesByComponent[component] = mustRead(t, path)
		}
		return bytesByComponent
	}
	seeded := slotBytes()

	if err := os.WriteFile(fixture.attestationPath(t), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	fixture.buildAndPublish(t, attestationFixtureTime)
	if inspection := fixture.verify(t); !inspection.Attested {
		t.Fatalf("verify over the first attestation = %+v, want it to replace what sat at its address", inspection)
	}
	first := mustRead(t, fixture.attestationPath(t))

	fixture.buildAndPublish(t, attestationFixtureTime.Add(time.Minute))
	if inspection := fixture.verify(t); !inspection.Attested {
		t.Fatalf("verify over the re-authored attestation = %+v, want it to hold", inspection)
	}
	if second := mustRead(t, fixture.attestationPath(t)); string(second) == string(first) {
		t.Fatalf("re-authoring left the attestation bytes unchanged: %s", second)
	}

	if after := slotBytes(); !reflect.DeepEqual(after, seeded) {
		t.Fatalf("component slot bytes moved when the attestation was authored: %v, want %v", after, seeded)
	}
	// One record per slot plus the single attestation: a second attestation left beside the
	// first, or a temporary a failed replacement abandoned, shows up here as an extra entry.
	entries, err := os.ReadDir(fixture.storeDir(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != len(family)+1 {
		names := make([]string, 0, len(entries))
		for _, entry := range entries {
			names = append(names, entry.Name())
		}
		t.Fatalf("store holds %d entries %v, want %d slots and exactly one attestation", len(entries), names, len(family))
	}
}

// [PS22b] An attestation answers for the artifact it was authored against and no other.
//
// The store lives under the common git directory, which every worktree of a repository
// shares while each builds its own dist/bench at its own path. One shared address would have
// each worktree's build clobber the rest, leaving every other worktree reading an attestation
// that names a binary it did not build. That refuses rather than skips, so nothing unsound
// follows — but the build phase then runs on every gate run in every worktree, which is the
// whole saving this class exists to make.
//
// The two artifacts are built from different packages rather than from one edited package, so
// they differ in bytes while a single source digest still answers for both seals. What moves
// between them is the address alone.
func TestAttestationsAreBoundToTheirArtifact(t *testing.T) {
	fixture := newAttestationFixture(t)
	writeGateTestFile(t, fixture.root, "cmd/alt/main.go",
		"package main\n\nvar built = \"alt\"\n\nfunc main() { _ = built }\n", 0o644)
	alt := filepath.Join(fixture.root, "dist", "bench-alt")

	fixture.buildAndPublish(t, attestationFixtureTime)
	primary := mustRead(t, fixture.attestationPath(t))

	// A second sealed artifact no gate build attested: its own address holds nothing, and the
	// first artifact's attestation does not stand in for it.
	unattested := alt + ".unattested"
	buildFixtureBinaryTo(t, fixture.root, "./cmd/alt", unattested)
	if err := benchfreshness.Publish(fixture.root, unattested, alt); err != nil {
		t.Fatalf("publish the second artifact: %v", err)
	}
	if inspection := fixture.verifyAt(t, alt); inspection.Attested {
		t.Fatalf("verify an unattested second artifact = %+v, want a refusal", inspection)
	}

	// The artifacts differ in bytes, without which no assertion below could tell one
	// attestation from the other.
	_, primaryDigest, err := benchfreshness.SealDigests(fixture.executable)
	if err != nil {
		t.Fatal(err)
	}
	_, altDigest, err := benchfreshness.SealDigests(alt)
	if err != nil {
		t.Fatal(err)
	}
	if primaryDigest == altDigest {
		t.Fatalf("both artifacts hash to %s, so this row cannot observe the address", primaryDigest)
	}

	fixture.buildAndPublishAt(t, "./cmd/alt", alt, attestationFixtureTime.Add(time.Minute))

	if inspection := fixture.verifyAt(t, alt); !inspection.Attested {
		t.Fatalf("verify the second artifact after its own gate build = %+v, want it attested", inspection)
	}
	if inspection := fixture.verify(t); !inspection.Attested {
		t.Fatalf("verify the first artifact after the second was attested = %+v, want it undisturbed", inspection)
	}
	if after := mustRead(t, fixture.attestationPath(t)); !bytes.Equal(after, primary) {
		t.Fatalf("attesting %s moved the first artifact's attestation: %s, want %s", alt, after, primary)
	}

	// One artifact keeps one address however its path is spelled. The seal contract refuses
	// every symlinked component, so cleaning is the only thing separating two spellings that
	// can both reach the same file.
	spelled := filepath.Join(fixture.root, "dist") + string(filepath.Separator) + filepath.Join("..", "dist", "bench")
	if inspection := fixture.verifyAt(t, spelled); !inspection.Attested {
		t.Fatalf("verify %q = %+v, want the attestation the artifact's cleaned path resolves", spelled, inspection)
	}
}
