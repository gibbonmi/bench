package gate

// Component identities are graded against the kit-shaped fixture, the only root here that
// carries a real module, a resolved kit phase table, and a sealed dist/bench — the three
// things an identity is made of.
//
// Two of the properties cannot be observed by moving the tree: the policy domain and the
// execution closure are what separate components whose declared files are identical, and
// the fixture's components differ in their files as well. Those tests drive
// componentIdentity directly with declarations built to differ in exactly one thing, so a
// hash that ignored that thing has nowhere to hide.

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
)

func mustResolveComponentIdentities(t *testing.T, root string) map[string]string {
	t.Helper()
	identities, err := ResolveComponentIdentities(root)
	if err != nil {
		t.Fatalf("ResolveComponentIdentities = %v, want the fixture's identities", err)
	}
	return identities
}

func mustTreeSnapshot(t *testing.T, root string) treeSnapshot {
	t.Helper()
	snapshot, err := readTreeSnapshot(root)
	if err != nil {
		t.Fatalf("readTreeSnapshot = %v, want the fixture's snapshot", err)
	}
	return snapshot
}

func componentIdentityOf(t *testing.T, identities map[string]string, name string) string {
	t.Helper()
	identity, computed := identities[name]
	if !computed {
		t.Fatalf("no identity for %q; computed %v", name, identities)
	}
	return identity
}

// mustComputeIdentity resolves one declaration against a phase of the caller's choosing.
func mustComputeIdentity(t *testing.T, root string, inputs ComponentInputs, phase Phase, snapshot treeSnapshot) string {
	t.Helper()
	identity, err := componentIdentity(root, inputs, phase, snapshot)
	if err != nil {
		t.Fatalf("componentIdentity(%s) = %v", inputs.Component(), err)
	}
	return identity
}

// undeclaredPath returns a fixture path none of the declarations name, so an edit to it is
// an edit the identities under test are supposed to be blind to.
func undeclaredPath(t *testing.T, sets map[string]ComponentInputs) string {
	t.Helper()
	for _, candidate := range captureSurfacePaths(ReducedScope()) {
		declared := false
		for _, inputs := range sets {
			declared = declared || slices.Contains(inputs.Paths(), candidate)
		}
		if !declared {
			return candidate
		}
	}
	t.Fatal("every capture surface is a declared input; the fixture can no longer express an undeclared edit")
	return ""
}

// [PS15] An identity moves with the content of the files its declaration names and stands
// still for everything else. Hashing the declared names alone would satisfy the second half
// and none of the first: a component whose sources changed would keep inheriting evidence
// taken against the sources it had before.
func TestComponentIdentityTracksDeclaredContent(t *testing.T) {
	fixture := newKitShapedFixture(t)
	sets := mustResolveComponentInputs(t, fixture.root)
	before := mustResolveComponentIdentities(t, fixture.root)

	const declared = "internal/canary/canary.go"
	if paths := componentEntry(t, sets, canary.PhaseTest).Paths(); !slices.Contains(paths, declared) {
		t.Fatalf("%s inputs omit %q, so this test observes nothing: %v", canary.PhaseTest, declared, paths)
	}
	writeGateTestFile(t, fixture.root, declared,
		"package canary\n\n// Name is the surface this package's own test grades.\nfunc Name() string { return \"canary\" }\n\nvar edited = 1\n", 0o644)

	edited := mustResolveComponentIdentities(t, fixture.root)
	if got, was := componentIdentityOf(t, edited, canary.PhaseTest), componentIdentityOf(t, before, canary.PhaseTest); got == was {
		t.Fatalf("%s identity = %s after editing %q, want it to move", canary.PhaseTest, got, declared)
	}

	outside := undeclaredPath(t, sets)
	writeGateTestFile(t, fixture.root, outside, "edited outside every declaration\n", 0o644)
	after := mustResolveComponentIdentities(t, fixture.root)
	if !reflect.DeepEqual(after, edited) {
		t.Fatalf("identities = %v after editing the undeclared %q, want %v unmoved", after, outside, edited)
	}
}

// [PS15b] A declared directory covers every file beneath it, so a component's identity
// answers for the directory's whole content: a file that lands under it moves the identity,
// and so does an edit to a file already there. Resolving the directory as a path of its own
// finds nothing — git names no directory in a snapshot — so a component declared this way
// would refuse an identity and run on every changeset, which is the whole of what the
// scoping buys.
func TestComponentIdentityCoversDeclaredDirectoryDescendants(t *testing.T) {
	fixture := newKitShapedFixture(t)
	const dir = "internal/canary/"
	const component = "canary"
	if paths := componentEntry(t, mustResolveComponentInputs(t, fixture.root), component).Paths(); !slices.Contains(paths, dir) {
		t.Fatalf("%s inputs omit the directory %q, so this test observes nothing: %v", component, dir, paths)
	}
	before := componentIdentityOf(t, mustResolveComponentIdentities(t, fixture.root), component)

	added := dir + "added.txt"
	writeGateTestFile(t, fixture.root, added, "a surface that landed under the declaration\n", 0o644)
	landed := componentIdentityOf(t, mustResolveComponentIdentities(t, fixture.root), component)
	if landed == before {
		t.Fatalf("%s identity = %s after %q landed under %q, want it to move", component, landed, added, dir)
	}

	writeGateTestFile(t, fixture.root, added, "the same surface, edited\n", 0o644)
	edited := componentIdentityOf(t, mustResolveComponentIdentities(t, fixture.root), component)
	if edited == landed {
		t.Fatalf("%s identity = %s after editing %q, want it to move", component, edited, added)
	}
}

// [PS15c] A declared file the snapshot has no entry for refuses, and names the path it
// could not resolve — the declaration and the tree disagree at a point the declaration
// named exactly, and an identity computed over what happened to resolve would address a
// slot answering for fewer inputs than the component reads.
//
// A declared directory covering nothing is the other judgment: it contributes nothing and
// the identity stands. Git tracks no empty directory, so refusing would leave a component
// naming a not-yet-landed surface permanently unable to compute an identity, while the
// first file to land under it moves the identity anyway.
func TestComponentIdentityRefusesAnAbsentDeclaredFile(t *testing.T) {
	fixture := newKitShapedFixture(t)
	snapshot := mustTreeSnapshot(t, fixture.root)
	const declaredDir = "internal/canary/"
	phase := Phase{Name: "canary", Argv: []string{"bash", "-c", "true"}}
	inputsOver := func(paths ...string) ComponentInputs {
		return ComponentInputs{component: "canary", source: SourceHandDeclared, paths: paths}
	}

	t.Run("absent file", func(t *testing.T) {
		const absent = "internal/canary/never-written.go"
		if _, tracked := snapshot.entry(absent); tracked {
			t.Fatalf("%q is in the snapshot, so this test observes nothing", absent)
		}
		identity, err := componentIdentity(fixture.root, inputsOver(declaredDir, absent), phase, snapshot)
		if err == nil {
			t.Fatalf("componentIdentity over the absent %q = %s, want a refusal", absent, identity)
		}
		if !strings.Contains(err.Error(), absent) {
			t.Fatalf("componentIdentity = %v, want the refusal to name %q", err, absent)
		}
		if identity != "" {
			t.Fatalf("componentIdentity returned %q alongside its refusal, want none", identity)
		}
	})

	t.Run("empty directory", func(t *testing.T) {
		const empty = "tests/never-landed/"
		for _, entry := range snapshot.entries {
			if strings.HasPrefix(entry.Path, empty) {
				t.Fatalf("%q covers %q, so this test observes nothing", empty, entry.Path)
			}
		}
		with := mustComputeIdentity(t, fixture.root, inputsOver(declaredDir, empty), phase, snapshot)
		without := mustComputeIdentity(t, fixture.root, inputsOver(declaredDir), phase, snapshot)
		if with != without {
			t.Fatalf("identity over %v = %s, want the empty %q to contribute nothing and leave %s",
				[]string{declaredDir, empty}, with, empty, without)
		}
	})
}

// [PS16] Two components declaring one file set and running one command still address
// different slots, because the component's name is part of its policy domain. Without it
// the two would share an address, and a component that had never run would inherit the
// other's green.
func TestComponentIdentitiesAreDomainSeparated(t *testing.T) {
	fixture := newKitShapedFixture(t)
	snapshot := mustTreeSnapshot(t, fixture.root)
	shared := componentEntry(t, mustResolveComponentInputs(t, fixture.root), canary.PhaseTest)
	phase := Phase{Name: "shared", Argv: []string{"bash", "-c", "true"}, Env: []string{"SHARED=1"}}

	twin := func(name string) ComponentInputs {
		return ComponentInputs{component: name, source: shared.Source(), paths: shared.Paths(), digests: shared.Digests()}
	}
	first := mustComputeIdentity(t, fixture.root, twin("alpha"), phase, snapshot)
	second := mustComputeIdentity(t, fixture.root, twin("beta"), phase, snapshot)
	if first == second {
		t.Fatalf("two components over one declaration and one command share the identity %s", first)
	}
}

// [PS17] An identity moves when the command or the environment contract behind it moves,
// with every declared file untouched. An identity made of input contents alone would let a
// component inherit evidence produced by a check that no longer runs.
func TestComponentIdentityTracksItsExecutionClosure(t *testing.T) {
	fixture := newKitShapedFixture(t)
	snapshot := mustTreeSnapshot(t, fixture.root)
	inputs := componentEntry(t, mustResolveComponentInputs(t, fixture.root), canary.PhaseTest)
	phase := Phase{Name: canary.PhaseTest, Argv: []string{"go", "test", "./..."}, Env: []string{"GOFLAGS=-count=1"}}
	base := mustComputeIdentity(t, fixture.root, inputs, phase, snapshot)

	for _, moved := range []struct {
		what  string
		phase Phase
	}{
		{"argv", Phase{Name: phase.Name, Argv: []string{"go", "test", "-short", "./..."}, Env: phase.Env}},
		{"env", Phase{Name: phase.Name, Argv: phase.Argv, Env: []string{"GOFLAGS=-count=2"}}},
	} {
		t.Run(moved.what, func(t *testing.T) {
			if got := mustComputeIdentity(t, fixture.root, inputs, moved.phase, snapshot); got == base {
				t.Fatalf("identity = %s with the %s moved and every input file unchanged, want it to move", got, moved.what)
			}
		})
	}
}

// [PS18] The contract component execs the published binary, so its identity moves when the
// seal's source digest moves even though its declared files are byte-identical. The
// snapshot is compared across the republish to prove that is what happened, and the
// toolchain components — which exec nothing — are held still by the same comparison.
func TestContractIdentityTracksTheSeal(t *testing.T) {
	fixture := newKitShapedFixture(t)
	before := mustResolveComponentIdentities(t, fixture.root)
	snapshotBefore := mustTreeSnapshot(t, fixture.root)

	const main = "cmd/bench/main.go"
	original, err := os.ReadFile(filepath.Join(fixture.root, filepath.FromSlash(main)))
	if err != nil {
		t.Fatal(err)
	}
	writeGateTestFile(t, fixture.root, main, "package main\n\nfunc main() { _ = 1 }\n", 0o644)
	sealKitShapedBinary(t, fixture.root)
	writeGateTestFile(t, fixture.root, main, string(original), 0o644)

	if got := mustTreeSnapshot(t, fixture.root); !reflect.DeepEqual(got, snapshotBefore) {
		t.Fatal("the republish moved the tree snapshot, so this test no longer isolates the seal")
	}
	after := mustResolveComponentIdentities(t, fixture.root)
	if got, was := componentIdentityOf(t, after, canary.PhaseContract), componentIdentityOf(t, before, canary.PhaseContract); got == was {
		t.Fatalf("contract identity = %s after the seal's source digest moved, want it to move", got)
	}
	if got, was := componentIdentityOf(t, after, canary.PhaseTest), componentIdentityOf(t, before, canary.PhaseTest); got != was {
		t.Fatalf("%s identity = %s after the republish, want %s — it reads no binary", canary.PhaseTest, got, was)
	}
}

// [PS19] Every way an identity can fail to be computed returns an error and no identity.
// A zero identity handed back instead would be a real address that every component could
// match, which is a skip over work nobody graded.
func TestComponentIdentityFailsClosed(t *testing.T) {
	t.Run("unreadable snapshot", func(t *testing.T) {
		fixture := newKitShapedFixture(t)
		if err := os.RemoveAll(filepath.Join(fixture.root, ".git")); err != nil {
			t.Fatal(err)
		}
		identities, err := ResolveComponentIdentities(fixture.root)
		if err == nil {
			t.Fatalf("ResolveComponentIdentities = %v, want an error with no snapshot to read", identities)
		}
		if identities != nil {
			t.Fatalf("identities = %v alongside the error, want none", identities)
		}
	})

	t.Run("failed derivation", func(t *testing.T) {
		fixture := newKitShapedFixture(t)
		writeGateTestFile(t, fixture.root, "go.mod", "this is not a go.mod\n", 0o644)

		identities, err := ResolveComponentIdentities(fixture.root)
		if err == nil {
			t.Fatalf("ResolveComponentIdentities = %v, want an error on the corrupted module", identities)
		}
		if identities != nil {
			t.Fatalf("identities = %v alongside the error, want none", identities)
		}
	})

	t.Run("declared path outside the repository", func(t *testing.T) {
		fixture := newKitShapedFixture(t)
		snapshot := mustTreeSnapshot(t, fixture.root)
		for _, escaping := range []string{"../outside.go", "/etc/hostname", ""} {
			inputs := ComponentInputs{component: canary.PhaseTest, source: SourceModuleTestClosure, paths: []string{escaping}}
			identity, err := componentIdentity(fixture.root, inputs, Phase{Name: canary.PhaseTest}, snapshot)
			if err == nil {
				t.Fatalf("componentIdentity with the declared path %q = %s, want an error", escaping, identity)
			}
			if identity != "" {
				t.Fatalf("componentIdentity with the declared path %q returned %q alongside its error, want none", escaping, identity)
			}
		}
	})
}
