//go:build system

package systemtest

import (
	"path/filepath"
	"strings"
	"testing"
)

// wrapperPoolHomeRefusal is the wrapper's own diagnostic for a missing pool home,
// written here independently of bin/bench.sh so reverting the guard to a bare $HOME
// reds this test rather than passing on a re-derived expectation. bash prefixes the
// message with the script, the line, and the variable name, so only the wrapper-authored
// remainder is asserted.
const wrapperPoolHomeRefusal = "the Bench pool home needs BENCH_HOME set, or HOME set to derive it from"

// TestWrapperNamesTheMissingPoolHome covers the first command an adopter runs. The
// wrapper derives BENCH_HOME from HOME under set -u, so an environment carrying neither
// must name the missing input rather than die on a raw bash diagnostic that gives no
// action. An already-set BENCH_HOME still wins, which is the precedence the guard sits
// behind, so both legs run.
func TestWrapperNamesTheMissingPoolHome(t *testing.T) {
	repo := owner.repos[0]
	wrapper := filepath.Join(owner.kit, "bin", "bench.sh")
	// An override with no "=" removes that variable from the child environment.
	launch := func(overrides ...string) processResult {
		if err := owner.observeSelected(); err != nil {
			t.Fatal(err)
		}
		environment := []string{"BENCH_RUN_BINARY=" + owner.selected.path, "BENCH_KIT=" + owner.kit}
		return owner.runAt(repo, append(environment, overrides...), "bash", wrapper, "version")
	}

	missing := launch("HOME", "BENCH_HOME")
	if missing.code == 0 {
		t.Fatalf("wrapper with no pool-home input = (%d, %q, %q), want a nonzero exit", missing.code, missing.stdout, missing.stderr)
	}
	if !strings.Contains(missing.stderr, wrapperPoolHomeRefusal) {
		t.Fatalf("wrapper does not name its own missing input: %q", missing.stderr)
	}
	if strings.Contains(missing.stderr, "unbound variable") {
		t.Fatalf("wrapper still fails with the raw bash diagnostic: %q", missing.stderr)
	}

	pool := filepath.Join(t.TempDir(), "bench-home")
	derived := launch("HOME", "BENCH_HOME="+pool)
	if derived.code != 0 || !strings.Contains(derived.stdout, "bench ") {
		t.Fatalf("wrapper with BENCH_HOME set and HOME unset = (%d, %q, %q), want the version line", derived.code, derived.stdout, derived.stderr)
	}
	assertPrivateHomeEmpty(t, pool)
}
