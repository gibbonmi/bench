package gate

import (
	"strings"
	"testing"
)

// oracleIdentityFixture answers a closed root whose manifest declares one tool, so the
// tool's bytes ride the closure that both digests hash.
func oracleIdentityFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeGateFixtureFile(t, root, ".bench/gate.sh", "#!/bin/sh\nexit 0\n", 0o755)
	writeGateFixtureFile(t, root, ".bench/gate-inputs.json", "{\"schema\":1,\"closure\":\"local\",\"environment\":[],\"paths\":[],\"tools\":[\"bin/tool\"]}\n", 0o644)
	writeGateFixtureFile(t, root, "bin/tool", "tool v1\n", 0o755)
	return root
}

func oracleIdentitySubject(t *testing.T, root, tree string) subject {
	t.Helper()
	s, err := buildSubjectForTree(root, root, policyVersion, tree)
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

// NG3: one pass feeds the closure to both digests, so a tool change moves both.
func TestOracleIdentityAndOracleShareTheClosure(t *testing.T) {
	root := oracleIdentityFixture(t)
	tree := strings.Repeat("a", 40)

	before := oracleIdentitySubject(t, root, tree)
	writeGateFixtureFile(t, root, "bin/tool", "tool v2\n", 0o755)
	after := oracleIdentitySubject(t, root, tree)

	if before.Identity == after.Identity {
		t.Fatalf("identity = %s across tool bytes, want two identities", before.Identity)
	}
	if before.Oracle == after.Oracle {
		t.Fatalf("oracle = %s across tool bytes, want two oracle digests", before.Oracle)
	}
}
