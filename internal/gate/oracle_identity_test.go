package gate

import (
	"strings"
	"testing"

	benchfreshness "github.com/gibbonmi/bench/internal/freshness"
)

// oracleIdentityFixture answers a closed root whose manifest declares one tool, so the
// tool's bytes ride the closure that both digests hash.
func oracleIdentityFixture(t *testing.T) string {
	t.Helper()
	return oracleIdentityFixtureAt(t, t.TempDir())
}

// oracleIdentityFixtureAt writes the fixture of oracleIdentityFixture at root.
func oracleIdentityFixtureAt(t *testing.T, root string) string {
	t.Helper()
	writeGateFixtureFile(t, root, ".bench/gate.sh", "#!/bin/sh\nexit 0\n", 0o755)
	writeOracleIdentityManifest(t, root, ``, ``, `"bin/tool"`)
	writeGateFixtureFile(t, root, "bin/tool", "tool v1\n", 0o755)
	return root
}

// writeOracleIdentityManifest writes a local manifest whose lists hold the given JSON
// elements.
func writeOracleIdentityManifest(t *testing.T, root, environment, paths, tools string) {
	t.Helper()
	writeGateFixtureFile(t, root, ".bench/gate-inputs.json", `{"schema":1,"closure":"local","environment":[`+environment+`],"paths":[`+paths+`],"tools":[`+tools+`]}`+"\n", 0o644)
}

func oracleIdentitySubject(t *testing.T, root, tree string) subject {
	t.Helper()
	return oracleIdentitySubjectFor(t, root, root, policyVersion, tree)
}

func oracleIdentitySubjectFor(t *testing.T, root, identityRoot, policy, tree string) subject {
	t.Helper()
	return oracleIdentitySubjectIn(t, root, identityRoot, policy, tree, false)
}

// oracleIdentitySubjectIn answers the closed subject of root in the ordinary or the
// prospective mode, with identityRoot as its baseline.
func oracleIdentitySubjectIn(t *testing.T, root, identityRoot, policy, tree string, prospective bool) subject {
	t.Helper()
	s, err := buildSubjectInMode(root, identityRoot, identityRoot, policy, tree, prospective)
	if err != nil {
		t.Fatal(err)
	}
	if !s.Closed {
		t.Fatalf("subject open: %s", s.Reason)
	}
	return s
}

// NG1: the identity is tree-free, so a later landing can compare it across trees.
func TestOracleIdentityIgnoresTheTree(t *testing.T) {
	root := oracleIdentityFixture(t)

	first := oracleIdentitySubject(t, root, strings.Repeat("a", 40))
	second := oracleIdentitySubject(t, root, strings.Repeat("b", 40))

	if first.Identity != second.Identity {
		t.Fatalf("identity = %s and %s, want one identity across trees", first.Identity, second.Identity)
	}
	if first.Oracle == second.Oracle {
		t.Fatalf("oracle = %s for both trees, want two oracle digests", first.Oracle)
	}
}

// NG2: the capability posture enters the identity and only the identity, so a posture
// change forces a whole run and a recorded green stays reusable.
func TestOracleIdentityCarriesTheCapabilityPosture(t *testing.T) {
	root := oracleIdentityFixture(t)
	tree := strings.Repeat("a", 40)

	t.Setenv(requireCapabilitiesEnv, "0")
	relaxed := oracleIdentitySubject(t, root, tree)
	t.Setenv(requireCapabilitiesEnv, "1")
	required := oracleIdentitySubject(t, root, tree)

	if relaxed.Identity == required.Identity {
		t.Fatalf("identity = %s under both postures, want two identities", relaxed.Identity)
	}
	if relaxed.Oracle != required.Oracle {
		t.Fatalf("oracle = %s and %s, want one oracle digest across postures", relaxed.Oracle, required.Oracle)
	}
}

// setRunningExecutable makes digest the running executable digest until t ends.
func setRunningExecutable(t *testing.T, digest string) {
	t.Helper()
	previous := runningExecutableDigest
	t.Cleanup(func() { runningExecutableDigest = previous })
	runningExecutableDigest = func() string { return digest }
}

// An unreadable running executable opens the subject, like every other unavailable closure
// member, so it reuses nothing.
func TestOracleIdentityOpensWithoutTheRunningExecutable(t *testing.T) {
	root := oracleIdentityFixture(t)
	setRunningExecutable(t, "")

	s, err := buildSubjectForTree(root, root, policyVersion, strings.Repeat("a", 40))
	if err != nil {
		t.Fatal(err)
	}
	if s.Closed || s.Reason != "running executable unavailable" {
		t.Fatalf("subject closed = %t, reason %q, want open with \"running executable unavailable\"", s.Closed, s.Reason)
	}
}

// NG3, GR8: one pass feeds the closure to both digests, so a change of a verdict input moves
// both. The policy, the identity root, and the baseline runner move the oracle digest only,
// and the run binary member and the running executable move the grade identity only. Each
// case changes one input and builds the subject on each side of it.
func TestOracleIdentityAndOracleShareTheClosure(t *testing.T) {
	t.Setenv("BENCH_KIT", "")
	tree := strings.Repeat("a", 40)
	build := func(t *testing.T, root, identityRoot, policy string) subject {
		t.Helper()
		return oracleIdentitySubjectFor(t, root, identityRoot, policy, tree)
	}
	// bareRoot holds no executable, so no closure path spells the identity root.
	bareRoot := func(t *testing.T) string {
		root := t.TempDir()
		writeOracleIdentityManifest(t, root, ``, ``, ``)
		return root
	}
	cases := []struct {
		input            string
		identity, oracle bool
		change           func(t *testing.T) (before, after subject)
	}{
		{"tool bytes", true, true, func(t *testing.T) (subject, subject) {
			root := oracleIdentityFixture(t)
			before := build(t, root, root, policyVersion)
			writeGateFixtureFile(t, root, "bin/tool", "tool v2\n", 0o755)
			return before, build(t, root, root, policyVersion)
		}},
		{"policy", false, true, func(t *testing.T) (subject, subject) {
			root := oracleIdentityFixture(t)
			return build(t, root, root, policyVersion), build(t, root, root, policyVersion+"/next")
		}},
		{"identity root", false, true, func(t *testing.T) (subject, subject) {
			root := bareRoot(t)
			return build(t, root, root, policyVersion), build(t, root, t.TempDir(), policyVersion)
		}},
		{"run binary member", true, false, func(t *testing.T) (subject, subject) {
			root := ownedPhaseTableFixture(t)
			writeRunBinarySource(t, root, "")
			before := build(t, root, root, policyVersion)
			writeRunBinarySource(t, root, " println() ")
			return before, build(t, root, root, policyVersion)
		}},
		{"running executable", true, false, func(t *testing.T) (subject, subject) {
			root := oracleIdentityFixture(t)
			setRunningExecutable(t, "one")
			before := build(t, root, root, policyVersion)
			setRunningExecutable(t, "two")
			return before, build(t, root, root, policyVersion)
		}},
		{"PATH", true, true, func(t *testing.T) (subject, subject) {
			root := oracleIdentityFixture(t)
			t.Setenv("PATH", "/usr/bin:/bin")
			before := build(t, root, root, policyVersion)
			t.Setenv("PATH", "/bin:/usr/bin")
			return before, build(t, root, root, policyVersion)
		}},
		// An owned prospective run of a tree that builds its own run binary selects the graded
		// tree as the run binary source, so the baseline's recipe moves the oracle digest and
		// not the run binary member.
		{"baseline runner", false, true, func(t *testing.T) (subject, subject) {
			root, identityRoot := ownedPhaseTableFixture(t), t.TempDir()
			writeGateFixtureFile(t, root, prospectiveGatePath, "#!/bin/sh\nexit 0\n", 0o755)
			writeRunBinarySource(t, root, "")
			before := oracleIdentitySubjectIn(t, root, identityRoot, policyVersion, tree, true)
			writeGateFixtureFile(t, identityRoot, benchfreshness.BuildInputsManifest, benchfreshness.BuildInputLine("build_script", "scripts/go-build.sh"), 0o644)
			return before, oracleIdentitySubjectIn(t, root, identityRoot, policyVersion, tree, true)
		}},
		{"environment value", true, true, func(t *testing.T) (subject, subject) {
			root := oracleIdentityFixture(t)
			writeOracleIdentityManifest(t, root, `"NG3_ORACLE_INPUT"`, ``, `"bin/tool"`)
			t.Setenv("NG3_ORACLE_INPUT", "one")
			before := build(t, root, root, policyVersion)
			t.Setenv("NG3_ORACLE_INPUT", "two")
			return before, build(t, root, root, policyVersion)
		}},
		{"declared path bytes", true, true, func(t *testing.T) (subject, subject) {
			root := oracleIdentityFixture(t)
			writeOracleIdentityManifest(t, root, ``, `"data"`, `"bin/tool"`)
			writeGateFixtureFile(t, root, "data/input", "one\n", 0o644)
			before := build(t, root, root, policyVersion)
			writeGateFixtureFile(t, root, "data/input", "two\n", 0o644)
			return before, build(t, root, root, policyVersion)
		}},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			before, after := tc.change(t)
			if moved := before.Identity != after.Identity; moved != tc.identity {
				t.Fatalf("identity %s to %s across %s, want moved = %t", before.Identity, after.Identity, tc.input, tc.identity)
			}
			if moved := before.Oracle != after.Oracle; moved != tc.oracle {
				t.Fatalf("oracle %s to %s across %s, want moved = %t", before.Oracle, after.Oracle, tc.input, tc.oracle)
			}
		})
	}
}
