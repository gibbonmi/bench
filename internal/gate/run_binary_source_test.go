package gate

import (
	"os"
	"strings"
	"testing"

	benchenv "github.com/gibbonmi/bench/internal/env"
	"github.com/gibbonmi/bench/internal/runbinary"
)

// phaseTableGateScript is a gate that routes through the phase table, so a run with an owner
// selects a run binary.
const phaseTableGateScript = "#!/bin/sh\n" + phaseTableHandoff + "exit 0\n"

// runBinaryKit writes a kit tree at root: its gate script, an empty input manifest, a run
// binary source with main as the body of its one command, and a prospective script when
// prospective is true.
func runBinaryKit(t *testing.T, root, gate, main string, prospective bool) string {
	t.Helper()
	writeGateFixtureFile(t, root, ".bench/gate.sh", gate, 0o755)
	writeOracleIdentityManifest(t, root, ``, ``, ``)
	writeRunBinarySource(t, root, main)
	if prospective {
		writeGateFixtureFile(t, root, prospectiveGatePath, "#!/bin/sh\nexit 0\n", 0o755)
	}
	return root
}

// setRunBinaryOwner gives an ordinary run of this process the production owner, or removes
// every variable that gives it one.
func setRunBinaryOwner(t *testing.T, owned bool) {
	t.Helper()
	if owned {
		t.Setenv(runbinary.Env, "owner")
		return
	}
	for _, name := range []string{runbinary.Env, benchenv.WrapperEnv} {
		t.Setenv(name, "")
		if err := os.Unsetenv(name); err != nil {
			t.Fatal(err)
		}
	}
}

// ownedPhaseTableFixture answers the root of oracleIdentityFixture with a gate that routes
// through the phase table, under the production owner, so an ordinary run selects a run binary.
func ownedPhaseTableFixture(t *testing.T) string {
	t.Helper()
	setRunBinaryOwner(t, true)
	root := oracleIdentityFixture(t)
	writeGateFixtureFile(t, root, ".bench/gate.sh", phaseTableGateScript, 0o755)
	return root
}

// ordinaryIdentity answers the grade identity of an ordinary subject of root under kit.
func ordinaryIdentity(t *testing.T, root, kit string) string {
	t.Helper()
	t.Setenv("BENCH_KIT", kit)
	return oracleIdentitySubject(t, root, strings.Repeat("a", 40)).Identity
}

// prospectiveIdentity answers the grade identity of a prospective subject of composed, stored
// under identityRoot.
func prospectiveIdentity(t *testing.T, composed, identityRoot string) string {
	t.Helper()
	s, err := buildProspectiveSubjectForGeneration(composed, identityRoot, &treeGeneration{tree: strings.Repeat("a", 40)})
	if err != nil {
		t.Fatal(err)
	}
	if !s.Closed {
		t.Fatalf("subject open: %s", s.Reason)
	}
	return s.Identity
}

// BG12: an ordinary kit run whose BENCH_KIT names another checkout frames the build inputs of
// that checkout, so two kits with different build inputs give two identities, and a kit whose
// build inputs equal the graded tree gives the identity of a run with no BENCH_KIT.
func TestIdentityFramesTheRunSelection(t *testing.T) {
	setRunBinaryOwner(t, true)
	root := runBinaryKit(t, t.TempDir(), phaseTableGateScript, "", true)
	same := runBinaryKit(t, t.TempDir(), phaseTableGateScript, "", true)
	other := runBinaryKit(t, t.TempDir(), phaseTableGateScript, " println() ", true)
	if buildInputsIdentity(same) == buildInputsIdentity(other) || buildInputsIdentity(other) == unreadableBaselineRunner {
		t.Fatal("the two kits do not give two readable build-input digests")
	}

	if ordinaryIdentity(t, root, same) == ordinaryIdentity(t, root, other) {
		t.Fatal("identity is one value under two kits with different build inputs, want the identity to frame the kit that BENCH_KIT names")
	}
	if ordinaryIdentity(t, root, same) != ordinaryIdentity(t, root, "") {
		t.Fatal("identity under a kit with the graded tree's build inputs differs from the identity with no BENCH_KIT")
	}
}

// BG14: a run that selects no run binary frames one stable value, so the build inputs of the
// kit that BENCH_KIT names do not move its identity. A gate that does not route through the
// phase table selects none, and so does an ordinary run with no owner.
func TestNoSelectionFramesAStableValue(t *testing.T) {
	cases := []struct {
		name  string
		gate  string
		owned bool
	}{
		{"no phase-table gate", "#!/bin/sh\nexit 0\n", true},
		{"no owner", phaseTableGateScript, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			setRunBinaryOwner(t, tc.owned)
			root := runBinaryKit(t, t.TempDir(), tc.gate, "", false)
			one := runBinaryKit(t, t.TempDir(), tc.gate, "", false)
			other := runBinaryKit(t, t.TempDir(), tc.gate, " println() ", false)

			if ordinaryIdentity(t, root, one) != ordinaryIdentity(t, root, other) {
				t.Fatalf("identity moved with the kit's build inputs under %s, want one stable value", tc.name)
			}
		})
	}
}

// BG109: the grade key of a clean kit tree does not change. The base derivation framed the
// build inputs of the graded tree for a tree that holds a prospective script, and only the
// run-binary member of the identity changed, so each run of the tree under the production
// owner frames that same member.
func TestKitKeyIsUnchanged(t *testing.T) {
	setRunBinaryOwner(t, true)
	tree := strings.Repeat("a", 40)
	root := runBinaryKit(t, t.TempDir(), phaseTableGateScript, "", true)
	composed := runBinaryKit(t, t.TempDir(), phaseTableGateScript, "", true)
	base := buildInputsIdentity(root)
	if base == unboundBaselineRunner || base == unreadableBaselineRunner {
		t.Fatalf("build-input digest = %s, want a readable digest", base)
	}

	for _, kit := range []string{"", root} {
		t.Setenv("BENCH_KIT", kit)
		ordinary := oracleIdentitySubject(t, root, tree)
		if member := runBinaryMember(true, ordinary.Resolution, root, root); member != base {
			t.Fatalf("ordinary member = %s under BENCH_KIT %q, want the graded tree's build inputs %s", member, kit, base)
		}
	}
	t.Setenv("BENCH_KIT", "")
	prospective, err := buildProspectiveSubjectForGeneration(composed, root, &treeGeneration{tree: tree})
	if err != nil {
		t.Fatal(err)
	}
	if member := runBinaryMember(true, prospective.Resolution, composed, root); member != base {
		t.Fatalf("prospective member = %s, want the graded tree's build inputs %s", member, base)
	}
}

// BG15: one clean kit tree gives one run-binary member to its ordinary run and to its
// prospective run, with no BENCH_KIT and with BENCH_KIT naming the tree, so the two runs give
// one grade identity.
func TestOrdinaryAndProspectiveFrameOneSource(t *testing.T) {
	setRunBinaryOwner(t, true)
	root := runBinaryKit(t, t.TempDir(), phaseTableGateScript, "", true)
	composed := runBinaryKit(t, t.TempDir(), phaseTableGateScript, "", true)

	t.Setenv("BENCH_KIT", "")
	prospective := prospectiveIdentity(t, composed, root)
	for _, kit := range []string{"", root} {
		if ordinary := ordinaryIdentity(t, root, kit); ordinary != prospective {
			t.Fatalf("identity = %s ordinary under BENCH_KIT %q and %s prospective, want one grade identity", ordinary, kit, prospective)
		}
	}
}
