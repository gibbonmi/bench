package gate

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
)

// TestInnerContractPhaseNarrowsToItsPackage grades the argv a scoped canary fixture's
// inner gate runs: the contract phase tests the one package that owns the fixture's
// EXPECT instead of every package below internal/contract. The nested package is what
// kills the shortcut of keeping the value's last segment.
func TestInnerContractPhaseNarrowsToItsPackage(t *testing.T) {
	kit := contractKit(t, "surface/artifact")
	t.Setenv(canary.ContractPackageEnv, "surface/artifact")

	argv := contractArgv(t, t.TempDir(), kit, innerMode)
	want := []string{"go", "-C", kit, "test", "-count=1", "./internal/contract/surface/artifact"}
	if !reflect.DeepEqual(argv, want) {
		t.Fatalf("contract argv = %#v, want %#v", argv, want)
	}
}

// TestContractPhaseUnnarrowedWithoutAScope keeps every run that is not a scoped fixture
// byte-identical to the table it has today: an adopting repo, an outer gate, and an
// unscoped fixture all grade the whole contract subtree, and an outer run must ignore
// the variable even when an operator exports it.
func TestContractPhaseUnnarrowedWithoutAScope(t *testing.T) {
	kit := contractKit(t, "surface/artifact")
	root := t.TempDir()
	unscoped := []string{"go", "-C", kit, "test", "-count=1", "./internal/contract/..."}

	t.Run("variable absent", func(t *testing.T) {
		if err := os.Unsetenv(canary.ContractPackageEnv); err != nil {
			t.Fatal(err)
		}
		if argv := contractArgv(t, root, kit, innerMode); !reflect.DeepEqual(argv, unscoped) {
			t.Fatalf("contract argv = %#v, want %#v", argv, unscoped)
		}
	})

	t.Run("outer mode", func(t *testing.T) {
		t.Setenv(canary.ContractPackageEnv, "surface/artifact")
		if argv := contractArgv(t, root, kit, outerMode); !reflect.DeepEqual(argv, unscoped) {
			t.Fatalf("contract argv = %#v, want %#v", argv, unscoped)
		}
	})
}

// TestInnerContractScopeRejectsUnusableValues grades the backstop the sweep's own
// structural reds sit in front of. A value that silently ran the whole suite would let
// binding rot pass as green, and one that reached `go test` unchecked would red with a
// toolchain error naming nothing the fixture author can act on — so each class reds here,
// naming the value.
//
// The unclean spellings are the quiet half. Each of them resolves to a real directory, so
// nothing errors: "." grades ./internal/contract/. — the parent package, which holds no
// tests — and the rest build argv a package path was never meant to be. The fixture reds
// as "did not bite" forever, and the binding it was supposed to prove is never proven.
func TestInnerContractScopeRejectsUnusableValues(t *testing.T) {
	kit := contractKit(t, "surface/artifact")
	for name, value := range map[string]string{
		"empty":             "",
		"absolute":          "/internal/contract/axi",
		"traversal":         "surface/../../axi",
		"unknown dir":       "no-such-package",
		"current dir":       ".",
		"dot prefixed":      "./surface/artifact",
		"trailing slash":    "surface/artifact/",
		"doubled separator": "surface//artifact",
	} {
		t.Run(name, func(t *testing.T) {
			t.Setenv(canary.ContractPackageEnv, value)
			_, err := phaseTable(t.TempDir(), kit, innerMode)
			if err == nil {
				t.Fatalf("phaseTable err = nil, want %s rejected", name)
			}
			if !strings.Contains(err.Error(), canary.ContractPackageEnv) || !strings.Contains(err.Error(), fmt.Sprintf("%q", value)) {
				t.Fatalf("phaseTable err = %v, want a diagnostic naming %s=%q", err, canary.ContractPackageEnv, value)
			}
		})
	}
}

// TestDeclaredPhaseTableIgnoresContractScope pins the narrowing to the built-in table it
// belongs to. A root declaring its own phases owns their argv outright, so rewriting one
// would grade something other than what the root asked for.
func TestDeclaredPhaseTableIgnoresContractScope(t *testing.T) {
	kit := contractKit(t, "surface/artifact")
	root := t.TempDir()
	writeFile(t, manifestPath(root), `{"phases":[{"name":"contract","argv":["go","test","./internal/contract/..."]}]}`)
	t.Setenv(canary.ContractPackageEnv, "surface/artifact")

	argv := contractArgv(t, root, kit, innerMode)
	want := []string{"go", "test", "./internal/contract/..."}
	if !reflect.DeepEqual(argv, want) {
		t.Fatalf("declared contract argv = %#v, want %#v", argv, want)
	}
}

// contractKit is a kit checkout carrying pkg under internal/contract, which is the tree
// a scope value is resolved against.
func contractKit(t *testing.T, pkg string) string {
	t.Helper()
	kit := t.TempDir()
	if err := os.MkdirAll(filepath.Join(kit, "internal", "contract", filepath.FromSlash(pkg)), 0o755); err != nil {
		t.Fatal(err)
	}
	return kit
}

// contractArgv is the resolved table's contract phase argv, which is what a scope value
// can only be observed through.
func contractArgv(t *testing.T, root, kit string, mode phaseMode) []string {
	t.Helper()
	phases, err := phaseTable(root, kit, mode)
	if err != nil {
		t.Fatalf("phaseTable: %v", err)
	}
	phase, ok := phaseNamed(phases, canary.PhaseContract)
	if !ok {
		t.Fatalf("resolved table carries no %s phase: %v", canary.PhaseContract, phaseNames(phases))
	}
	return phase.Argv
}
