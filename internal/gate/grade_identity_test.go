package gate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	benchfreshness "github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/gittest"
)

// GR1, GR2: two checkouts with equal content, one of them under a path with a space, give
// one grade identity. The oracle digest keeps the checkout path, so the gate cache of each
// checkout keeps its own key. An absolute link to a tool inside the checkout leaves the
// checkout path out too.
func TestGradeIdentityLeavesOutTheCheckoutPath(t *testing.T) {
	t.Setenv("BENCH_KIT", "")
	tree := strings.Repeat("a", 40)
	for _, linked := range []bool{false, true} {
		checkout := func(name string) subject {
			root := oracleIdentityFixtureAt(t, filepath.Join(t.TempDir(), name))
			if linked {
				tool, err := filepath.EvalSymlinks(filepath.Join(root, "bin", "tool"))
				if err != nil {
					t.Fatal(err)
				}
				if err := os.Rename(tool, tool+".real"); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(tool+".real", tool); err != nil {
					t.Fatal(err)
				}
			}
			return oracleIdentitySubject(t, root, tree)
		}
		first, second := checkout("checkout"), checkout("other checkout")

		if first.Identity != second.Identity {
			t.Fatalf("identity = %s and %s with linked = %t, want one grade identity across checkouts", first.Identity, second.Identity, linked)
		}
		if first.Oracle == second.Oracle {
			t.Fatalf("oracle = %s for both checkouts with linked = %t, want two oracle digests", first.Oracle, linked)
		}
	}
}

// GR9: a reader table that a sibling checkout stored over equal content answers the
// landing, so the reader set runs.
func TestReaderPlanAdmitsATableOfASiblingCheckout(t *testing.T) {
	f := newLandingFixture(t, nil)
	sibling := filepath.Join(t.TempDir(), "sibling checkout")
	gittest.MustRun(t, f.root, "worktree", "add", "-q", "--detach", sibling, f.marker)
	f.storeTable(currentSubject(t, sibling).Identity)
	f.selects(docsChange)
}

// GR45: a PATH tool inside the checkout keeps one spelling, so an ordinary subject of the
// checkout and a prospective subject of the same tree give one grade identity.
func TestGradeIdentityKeepsTheOrdinaryAndProspectiveSubjectsEqual(t *testing.T) {
	t.Setenv("BENCH_KIT", "")
	tree := strings.Repeat("a", 40)
	checkout, composed := t.TempDir(), t.TempDir()
	for _, root := range []string{checkout, composed} {
		oracleIdentityFixtureAt(t, root)
		writeOracleIdentityManifest(t, root, ``, ``, `"tool"`)
	}
	t.Setenv("PATH", filepath.Join(checkout, "bin")+":/usr/bin:/bin")

	ordinary := oracleIdentitySubjectFor(t, checkout, checkout, policyVersion, tree)
	prospective := oracleIdentitySubjectFor(t, composed, checkout, policyVersion, tree)

	if ordinary.Identity != prospective.Identity {
		t.Fatalf("identity = %s ordinary and %s prospective, want one grade identity", ordinary.Identity, prospective.Identity)
	}
}

// writeRunBinarySource writes a Go source root that declares its build inputs, with main as
// the body of its one command.
func writeRunBinarySource(t *testing.T, root, main string) {
	t.Helper()
	writeGateFixtureFile(t, root, "go.mod", "module fixture\n\ngo 1.21\n", 0o644)
	writeGateFixtureFile(t, root, "cmd/bench/main.go", "package main\n\nfunc main() {"+main+"}\n", 0o644)
	writeGateFixtureFile(t, root, "scripts/go-build.sh", "#!/bin/sh\nexit 0\n", 0o755)
	writeGateFixtureFile(t, root, benchfreshness.BuildInputsManifest, benchfreshness.BuildInputLine("build_script", "scripts/go-build.sh"), 0o644)
}

// GR54: a prospective subject of a tree that builds its own run binary gives one grade
// identity from a worktree whose build inputs equal the tree and from a primary checkout
// whose build inputs differ.
func TestGradeIdentityFollowsTheRunBinarySource(t *testing.T) {
	t.Setenv("BENCH_KIT", "")
	setRunBinaryOwner(t, true)
	tree := strings.Repeat("a", 40)
	composed, worktree, primary := t.TempDir(), t.TempDir(), t.TempDir()
	writeGateFixtureFile(t, composed, ".bench/gate.sh", phaseTableGateScript, 0o755)
	writeGateFixtureFile(t, composed, prospectiveGatePath, "#!/bin/sh\nexit 0\n", 0o755)
	writeOracleIdentityManifest(t, composed, ``, ``, ``)
	writeRunBinarySource(t, composed, "")
	writeRunBinarySource(t, worktree, "")
	writeRunBinarySource(t, primary, " println() ")
	if graded, other := baselineRunnerIdentity(composed), baselineRunnerIdentity(primary); graded == other || graded == unreadableBaselineRunner || other == unreadableBaselineRunner {
		t.Fatalf("build-input digests = %s and %s, want two readable digests", graded, other)
	}

	if fromWorktree, fromPrimary := prospectiveIdentity(t, composed, worktree), prospectiveIdentity(t, composed, primary); fromWorktree != fromPrimary {
		t.Fatalf("identity = %s from the worktree and %s from the primary checkout, want one grade identity", fromWorktree, fromPrimary)
	}

	// An ordinary subject of the same source runs the identity root's kit, so the two build
	// inputs give two identities.
	ordinary := t.TempDir()
	writeGateFixtureFile(t, ordinary, ".bench/gate.sh", phaseTableGateScript, 0o755)
	writeOracleIdentityManifest(t, ordinary, ``, ``, ``)
	writeRunBinarySource(t, ordinary, "")
	if one, other := oracleIdentitySubjectFor(t, ordinary, worktree, policyVersion, tree), oracleIdentitySubjectFor(t, ordinary, primary, policyVersion, tree); one.Identity == other.Identity {
		t.Fatalf("identity = %s for an ordinary subject from both kits, want two identities", one.Identity)
	}
}

// GR55: two run binary source roots that declare no build inputs frame one stable value,
// whatever else they hold.
func TestGradeIdentityGivesAnUndeclaredSourceOneValue(t *testing.T) {
	t.Setenv("BENCH_KIT", "")
	tree := strings.Repeat("a", 40)
	root, first, second := ownedPhaseTableFixture(t), t.TempDir(), t.TempDir()
	writeGateFixtureFile(t, first, "go.mod", "module first\n\ngo 1.21\n", 0o644)
	writeGateFixtureFile(t, second, "go.mod", "module second\n\ngo 1.21\n", 0o644)

	one := oracleIdentitySubjectFor(t, root, first, policyVersion, tree)
	other := oracleIdentitySubjectFor(t, root, second, policyVersion, tree)

	if one.Identity != other.Identity {
		t.Fatalf("identity = %s and %s, want one grade identity for two undeclared sources", one.Identity, other.Identity)
	}
}
