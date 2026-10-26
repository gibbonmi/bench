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
	root := t.TempDir()
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
	s, err := buildSubjectForTree(root, identityRoot, policy, tree)
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

// NG3: one pass feeds the closure to both digests, so a change of any input but the tree
// moves both. Each case changes one input and builds the subject on each side of it.
func TestOracleIdentityAndOracleShareTheClosure(t *testing.T) {
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
		input  string
		change func(t *testing.T) (before, after subject)
	}{
		{"tool bytes", func(t *testing.T) (subject, subject) {
			root := oracleIdentityFixture(t)
			before := build(t, root, root, policyVersion)
			writeGateFixtureFile(t, root, "bin/tool", "tool v2\n", 0o755)
			return before, build(t, root, root, policyVersion)
		}},
		{"policy", func(t *testing.T) (subject, subject) {
			root := oracleIdentityFixture(t)
			return build(t, root, root, policyVersion), build(t, root, root, policyVersion+"/next")
		}},
		{"identity root", func(t *testing.T) (subject, subject) {
			root := bareRoot(t)
			return build(t, root, root, policyVersion), build(t, root, t.TempDir(), policyVersion)
		}},
		{"PATH", func(t *testing.T) (subject, subject) {
			root := oracleIdentityFixture(t)
			t.Setenv("PATH", "/usr/bin:/bin")
			before := build(t, root, root, policyVersion)
			t.Setenv("PATH", "/bin:/usr/bin")
			return before, build(t, root, root, policyVersion)
		}},
		{"baseline runner", func(t *testing.T) (subject, subject) {
			root, identityRoot := bareRoot(t), t.TempDir()
			writeGateFixtureFile(t, root, prospectiveGatePath, "#!/bin/sh\nexit 0\n", 0o755)
			before := build(t, root, identityRoot, policyVersion)
			writeGateFixtureFile(t, identityRoot, benchfreshness.BuildInputsManifest, benchfreshness.BuildInputLine("build_script", "scripts/go-build.sh"), 0o644)
			return before, build(t, root, identityRoot, policyVersion)
		}},
		{"environment value", func(t *testing.T) (subject, subject) {
			root := oracleIdentityFixture(t)
			writeOracleIdentityManifest(t, root, `"NG3_ORACLE_INPUT"`, ``, `"bin/tool"`)
			t.Setenv("NG3_ORACLE_INPUT", "one")
			before := build(t, root, root, policyVersion)
			t.Setenv("NG3_ORACLE_INPUT", "two")
			return before, build(t, root, root, policyVersion)
		}},
		{"declared path bytes", func(t *testing.T) (subject, subject) {
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
			if before.Identity == after.Identity {
				t.Fatalf("identity = %s across %s, want two identities", before.Identity, tc.input)
			}
			if before.Oracle == after.Oracle {
				t.Fatalf("oracle = %s across %s, want two oracle digests", before.Oracle, tc.input)
			}
		})
	}
}
