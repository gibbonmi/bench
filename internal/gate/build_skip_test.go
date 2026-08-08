package gate

// The build phase's skip observed from outside the gate: whether the phase's marker ran, what
// the artifact and its seal held afterwards, what the store attested, and what the verdict
// recorded. Build is the one component that skips through artifact reuse rather than an
// ancestor slot, so the failure class here is a run that execs a binary nobody built — which
// only a reading of the artifact itself can catch, never a return value.
//
// Every row seeds a full green first. The seed is what authors the attestation, and without it
// the build phase runs for want of evidence rather than for the reason the row names.

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/canary"
	benchfreshness "github.com/gibbonmi/bench/internal/freshness"
)

// observedBinaryMarker is where the canary phase leaves the bytes it found at dist/bench. It
// sits under .git, which no component declares, so recording what a reader saw does not move
// any identity and cannot itself change which components run.
const observedBinaryMarker = ".git/observed-bench"

// seededBuildFixture is a kit-shaped root whose first full green has attested its dist/bench,
// and whose canary phase copies the artifact it would exec. The copy is what makes "the readers
// exec the sealed binary" an observation rather than an inference: the canary phase declares
// no dependency the runner could satisfy by never starting it, so a byte difference here is a
// reader that got something other than what the seal answers for.
func seededBuildFixture(t *testing.T) kitShapedFixture {
	t.Helper()
	fixture := newKitShapedFixture(t)
	writeGateTestFile(t, fixture.root, ".bench/phase-canary.sh",
		"echo canary >> .git/phase-runs\ncp dist/bench "+observedBinaryMarker+"\n", 0o644)
	// The seed is forced rather than ordinary. A forced run grades every component past every
	// piece of evidence, so it authors the attestation whatever the skip decision would have
	// said — and a row below then reds on its own assertion rather than on this guard, which is
	// what a mutation to that decision has to be able to show.
	if got := executeWithEngineAfterAcquireAtKit(context.Background(), fixture.root, fixture.root, io.Discard, io.Discard, productionGateEngine{}, notifyGateSignals, forceRun); got.ActionExit != 0 {
		t.Fatalf("seed run exit = %d, want 0", got.ActionExit)
	}
	if !buildAttestationOf(t, fixture).Attested {
		t.Fatalf("the seed run left no attestation for %s; nothing below can observe a skip", fixture.binaryPath())
	}
	return fixture
}

func buildAttestationOf(t *testing.T, fixture kitShapedFixture) buildAttestationInspection {
	t.Helper()
	return verifyBuildAttestation(fixture.root, fixture.binaryPath(), time.Now().UTC())
}

// buildAttestationFile is the attestation record the store holds for the fixture's artifact,
// and nil where it holds none.
func buildAttestationFile(t *testing.T, fixture kitShapedFixture) []byte {
	t.Helper()
	path := attestationFixture{root: fixture.root, executable: fixture.binaryPath()}.attestationPath(t)
	data, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			t.Fatal(err)
		}
		return nil
	}
	return data
}

// editCaptureSurfaces makes the changeset a capture-only one, which is what leaves every
// evidence-covered component — build among them — with nothing of its own that moved.
func editCaptureSurfaces(t *testing.T, root, text string) {
	t.Helper()
	for _, path := range captureSurfacePaths(ReducedScope()) {
		writeGateTestFile(t, root, path, text+"\n", 0o644)
	}
}

func mustSealDigests(t *testing.T, executable string) (sources, binary string) {
	t.Helper()
	sources, binary, err := benchfreshness.SealDigests(executable)
	if err != nil {
		t.Fatalf("benchfreshness.SealDigests(%s) = %v", executable, err)
	}
	return sources, binary
}

// [PC3] A valid attested seal skips the build. The artifact and its seal are byte-identical
// afterwards, and the phase that reads the artifact reads exactly those bytes.
//
// dist/bench is backdated before the run so that its mtime predates every source in the tree.
// That is the one shape a timestamp-based freshness rule would call stale, and it is banned by
// decision — so a run that rebuilds here is reading a clock the predicate is not allowed to
// have. The changeset moves the canary fixture rather than a capture surface, because the
// canary is the reader: a capture-only changeset would skip it and leave nothing execing the
// artifact this row is about.
func TestAttestedSealSkipsTheBuild(t *testing.T) {
	t.Parallel()
	fixture := seededBuildFixture(t)
	sealed := mustRead(t, fixture.binaryPath())
	seal := mustRead(t, fixture.binaryPath()+".seal")
	writeGateTestFile(t, fixture.root, "tests/canary/fixture.txt", "moved canary fixture\n", 0o644)
	ancient := time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)
	if err := os.Chtimes(fixture.binaryPath(), ancient, ancient); err != nil {
		t.Fatal(err)
	}

	observation := observeGreenGate(t, fixture.root)

	if observation.ran(canary.PhaseBuild) {
		t.Fatalf("executed %v, want the build skipped on its attested seal:\n%s", observation.executed, observation.stdout)
	}
	if !observation.ran("canary") {
		t.Fatalf("executed %v, want the canary run — nothing else in this run reads the artifact", observation.executed)
	}
	if after := mustRead(t, fixture.binaryPath()); !reflect.DeepEqual(after, sealed) {
		t.Fatalf("dist/bench changed across a skipped build: %d bytes, want the sealed %d", len(after), len(sealed))
	}
	if after := mustRead(t, fixture.binaryPath()+".seal"); !reflect.DeepEqual(after, seal) {
		t.Fatalf("the seal changed across a skipped build: %s, want %s", after, seal)
	}
	observed := mustRead(t, filepath.Join(fixture.root, filepath.FromSlash(observedBinaryMarker)))
	if !reflect.DeepEqual(observed, sealed) {
		t.Fatalf("the canary phase read %d bytes at dist/bench, want the sealed binary's %d", len(observed), len(sealed))
	}
}

// [PC4] Every unsound artifact runs the build. The seal contract and the attestation each
// refuse a different way of holding something that is nearly evidence, and one fail-open case
// anywhere here execs a stale, unknown, or planted binary for the rest of the run.
//
// Each case is a fault applied to a root that has just proved it would otherwise skip — PC3
// owns that baseline — over a capture-only changeset, so the build phase running is the fault's
// doing and not the changeset's. The source-digest case is the exception by construction: a
// changeset that moves a build input is not capture-only, and that is the whole of what it
// asserts.
func TestBuildRunsOnEveryUnsoundArtifact(t *testing.T) {
	t.Parallel()
	for _, testCase := range []struct {
		name  string
		fault func(t *testing.T, fixture kitShapedFixture)
	}{
		{"absent binary", func(t *testing.T, fixture kitShapedFixture) {
			if err := os.Remove(fixture.binaryPath()); err != nil {
				t.Fatal(err)
			}
		}},
		{"absent seal", func(t *testing.T, fixture kitShapedFixture) {
			if err := os.Remove(fixture.binaryPath() + ".seal"); err != nil {
				t.Fatal(err)
			}
		}},
		{"source digest mismatch", func(t *testing.T, fixture kitShapedFixture) {
			// The sources move and the artifact does not, so the seal now answers for a tree
			// that no longer exists while still agreeing with the binary beside it.
			writeGateTestFile(t, fixture.root, "cmd/bench/main.go",
				"package main\n\nvar built = \"moved\"\n\nfunc main() { _ = built }\n", 0o644)
		}},
		{"executable digest mismatch", func(t *testing.T, fixture kitShapedFixture) {
			// The artifact moves and the seal does not, which is the inverse: the seal and the
			// attestation still agree with each other and neither answers for the bytes on disk.
			writeGateTestFile(t, fixture.root, "dist/bench", "#!/usr/bin/env bash\nexit 0\n", 0o755)
		}},
		{"symlinked sidecar", func(t *testing.T, fixture kitShapedFixture) {
			seal := fixture.binaryPath() + ".seal"
			elsewhere := seal + ".moved"
			if err := os.Rename(seal, elsewhere); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(elsewhere, seal); err != nil {
				t.Fatal(err)
			}
		}},
		{"missing attestation", func(t *testing.T, fixture kitShapedFixture) {
			path := attestationFixture{root: fixture.root, executable: fixture.binaryPath()}.attestationPath(t)
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
		}},
		{"attestation names another binary", func(t *testing.T, fixture kitShapedFixture) {
			// A record an author would produce, for bytes no artifact here holds: the seal is
			// untouched and self-consistent, so only the comparison against the store refuses.
			attestationFixture{root: fixture.root, executable: fixture.binaryPath()}.writeAttestationFile(t,
				`{"schema":1,"executable":"`+strings.Repeat("c", 64)+
					`","authored_at":"`+time.Now().UTC().Add(-time.Hour).Truncate(time.Second).Format(time.RFC3339)+`"}`)
		}},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			fixture := seededBuildFixture(t)
			editCaptureSurfaces(t, fixture.root, "capture-only edit")
			testCase.fault(t, fixture)

			observation := observeGate(t, fixture.root)
			if !observation.ran(canary.PhaseBuild) {
				t.Fatalf("executed %v, want the build run past a %s:\n%s", observation.executed, testCase.name, observation.stdout)
			}
		})
	}
}

// [PC5b] A planted binary published with its own recomputed seal runs the build, and the run
// re-authors the attestation over what it then found on disk.
//
// The plant is what anyone able to write beside dist/bench can produce: every digest agrees
// with every other, and freshness.Check passes on it — asserted first, because a plant the seal
// contract already refuses would make this row a second copy of PC4 rather than the case the
// attestation exists for. What refuses is the store, which never saw a gate build produce those
// bytes.
func TestPlantedArtifactRunsTheBuild(t *testing.T) {
	t.Parallel()
	fixture := seededBuildFixture(t)
	attested := buildAttestationFile(t, fixture)

	editFixtureMain(t, fixture.root, "planted")
	planted := fixture.binaryPath() + ".planted"
	buildFixtureBinaryTo(t, fixture.root, "./cmd/bench", planted)
	if err := benchfreshness.Publish(fixture.root, planted, fixture.binaryPath()); err != nil {
		t.Fatalf("plant the binary: %v", err)
	}
	if err := benchfreshness.Check(fixture.root, fixture.binaryPath()); err != nil {
		t.Fatalf("benchfreshness.Check over the plant = %v, want it to pass — a plant the seal contract "+
			"already refuses never reaches the attestation this row is about", err)
	}

	observation := observeGreenGate(t, fixture.root)

	if !observation.ran(canary.PhaseBuild) {
		t.Fatalf("executed %v, want the build run over an unattested plant:\n%s", observation.executed, observation.stdout)
	}
	after := buildAttestationFile(t, fixture)
	if after == nil || reflect.DeepEqual(after, attested) {
		t.Fatalf("the attestation after the build = %s, want it re-authored past the seed's %s", after, attested)
	}
	// Re-authored over the bytes the run left on disk, which is what makes the next run able to
	// skip: an attestation naming anything else would run the build forever.
	_, sealed := mustSealDigests(t, fixture.binaryPath())
	var record buildAttestationRecord
	if err := json.Unmarshal(after, &record); err != nil {
		t.Fatal(err)
	}
	if record.Executable != sealed {
		t.Fatalf("the re-authored attestation names %s, want the artifact's %s", record.Executable, sealed)
	}
}

// [PS47] A red run authors no attestation, and the planted artifact it leaves behind is still
// refused by the next run.
//
// This is the side door back into PC5b. A plant correctly runs the build, but a build can end
// red for reasons that have nothing to do with the plant — a transient toolchain error, a
// concurrent edit, a full disk — and a failed build never replaced the artifact. A run that
// attested outside its green branch would therefore attest whatever survived on disk, which is
// the plant, and the next run would find a seal and a matching attestation and exec it.
//
// The four steps are driven end to end because the hole is only visible as a sequence. The
// bytes pin the authorship and the second run pins what that authorship would have bought;
// neither half alone says a planted binary stays refused.
func TestRedRunAuthorsNoAttestation(t *testing.T) {
	t.Parallel()
	fixture := seededBuildFixture(t)

	// The plant is self-consistent against the sources present now, so nothing in the seal
	// contract refuses it — asserted, because a plant the seal already refuses would keep the
	// build running for a reason this row is not about.
	editFixtureMain(t, fixture.root, "planted")
	planted := fixture.binaryPath() + ".planted"
	buildFixtureBinaryTo(t, fixture.root, "./cmd/bench", planted)
	if err := benchfreshness.Publish(fixture.root, planted, fixture.binaryPath()); err != nil {
		t.Fatalf("plant the binary: %v", err)
	}
	if err := benchfreshness.Check(fixture.root, fixture.binaryPath()); err != nil {
		t.Fatalf("benchfreshness.Check over the plant = %v, want it to pass", err)
	}
	attested := buildAttestationFile(t, fixture)

	// The build runs — it has to, the plant is unattested — and fails. What moves is the phase's
	// own script, which no component declares, so the red arrives without moving any identity.
	forcePhaseRed(t, fixture.root, canary.PhaseBuild)
	red := observeGate(t, fixture.root)
	if red.exit == 0 {
		t.Fatalf("execution = %+v, want red", red)
	}
	if !red.ran(canary.PhaseBuild) {
		t.Fatalf("executed %v, want the build to have run and failed", red.executed)
	}
	if after := buildAttestationFile(t, fixture); !reflect.DeepEqual(after, attested) {
		t.Fatalf("the attestation moved across a red run:\nafter  %s\nbefore %s", after, attested)
	}

	// What the unmoved attestation buys: the artifact the failed build left behind is still the
	// plant, and the next run rebuilds over it rather than execing it.
	writeGateTestFile(t, fixture.root, ".bench/phase-"+canary.PhaseBuild+".sh",
		"echo "+canary.PhaseBuild+" >> .git/phase-runs\n", 0o644)
	next := observeGreenGate(t, fixture.root)
	if !next.ran(canary.PhaseBuild) {
		t.Fatalf("executed %v, want the build run over the plant the red run left behind:\n%s", next.executed, next.stdout)
	}
}

// [PC20] A run with nothing to inherit from grades everything, build included, and leaves the
// evidence a later run inherits: a slot for every scoped component and an attestation for the
// artifact. The three cases are the three ways a run can have nothing to inherit — no evidence
// was ever written, the store was pruned, or the operator forced a real run.
func TestFirstRunAndFreshBuildEverything(t *testing.T) {
	t.Parallel()
	assertGradedEverything := func(t *testing.T, fixture kitShapedFixture, executed []string) {
		t.Helper()
		want := fixture.phaseNames()
		sort.Strings(executed)
		sort.Strings(want)
		if !reflect.DeepEqual(executed, want) {
			t.Fatalf("executed %v, want the whole resolved table %v", executed, want)
		}
		for component, data := range slotBytes(t, fixture.root) {
			if data == nil {
				t.Fatalf("%s holds no slot after a run that graded it", component)
			}
		}
		if inspection := buildAttestationOf(t, fixture); !inspection.Attested {
			t.Fatalf("the artifact is attested %+v after a run that built it, want the run's own attestation", inspection)
		}
	}

	t.Run("a first run", func(t *testing.T) {
		// The fixture's dist/bench is sealed at construction and attested by nothing, which is
		// exactly the shape a clone arrives in: a self-consistent artifact no gate produced.
		fixture := newKitShapedFixture(t)
		if inspection := buildAttestationOf(t, fixture); inspection.Attested {
			t.Fatalf("the unbuilt fixture is already attested %+v; this row would observe nothing", inspection)
		}
		assertGradedEverything(t, fixture, observeGreenGate(t, fixture.root).executed)
	})

	t.Run("a pruned evidence store", func(t *testing.T) {
		fixture := seededBuildFixture(t)
		dir, err := componentSlotDir(fixture.root)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.RemoveAll(dir); err != nil {
			t.Fatal(err)
		}
		editCaptureSurfaces(t, fixture.root, "capture-only edit")
		assertGradedEverything(t, fixture, observeGreenGate(t, fixture.root).executed)
	})

	t.Run("a forced run", func(t *testing.T) {
		fixture := seededBuildFixture(t)
		editCaptureSurfaces(t, fixture.root, "capture-only edit")
		before := len(phaseRunNames(t, fixture.root))
		if got := executeWithEngineAfterAcquireAtKit(context.Background(), fixture.root, fixture.root, io.Discard, io.Discard, productionGateEngine{}, notifyGateSignals, forceRun); got.ActionExit != 0 {
			t.Fatalf("forced run exit = %d, want 0", got.ActionExit)
		}
		assertGradedEverything(t, fixture, append([]string(nil), phaseRunNames(t, fixture.root)[before:]...))
	})
}

// [PS30] The verdict's evidence for a skipped build names the seal's source digest — the build
// inputs the reused artifact answers for. The executable digest is the other digest in the same
// seal and says nothing about the tree, so the two are asserted to differ first: without that,
// this row would pass over either one.
func TestSkippedBuildEvidenceNamesTheSourceDigest(t *testing.T) {
	t.Parallel()
	fixture := seededBuildFixture(t)
	editCaptureSurfaces(t, fixture.root, "capture-only edit")

	observation := observeGreenGate(t, fixture.root)
	if observation.ran(canary.PhaseBuild) {
		t.Fatalf("executed %v, want the build skipped:\n%s", observation.executed, observation.stdout)
	}

	sources, sealed := mustSealDigests(t, fixture.binaryPath())
	if sources == sealed {
		t.Fatalf("the seal's source and executable digests are both %s; this row cannot tell them apart", sources)
	}
	rec := partialRecord(t, fixture.root)
	evidence := rec.SkipEvidence[canary.PhaseBuild]
	if evidence.Seal != sources {
		t.Fatalf("the build skipped on seal %q, want the seal's source digest %q (its executable digest is %q)",
			evidence.Seal, sources, sealed)
	}
}
