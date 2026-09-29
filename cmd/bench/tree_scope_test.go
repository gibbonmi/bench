package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/toon"
)

// TT6 and TT7: each repository-scoped verb refuses a tree target as its first argument,
// before the verb runs. The verbs come from the registry, so no test restates the scope
// table. The flag spelling is authored apart from the dispatcher, so a refusal that reads
// another argument fails. A wrapper-only definition never reaches the dispatcher.
func TestRepositoryVerbRefusesTreeTarget(t *testing.T) {
	t.Setenv(benchhome.Env, t.TempDir())
	// A verb that runs past a missing refusal must fail fast, so `bench shift` finds no
	// harness adapter and starts no loop.
	t.Setenv("BENCH_AGENT", "")
	root := newAXIEnvelopeRepo(t)
	// An ignored inbox lets `bench idea` park its line in this primary checkout, so a verb
	// that runs before the refusal leaves the file behind.
	writeAXIFixture(t, filepath.Join(root, ".gitignore"), "capture/\n")

	refused := 0
	for _, definition := range commandRegistry {
		if definition.Scope != scopeRepository || definition.WrapperOnly {
			continue
		}
		refused++
		want := toon.Usage("bench "+definition.Name, "--in") + "\n"
		for _, args := range [][]string{{"--in", "primary"}, {"--in"}} {
			argv := append([]string{definition.Name}, args...)
			t.Run(strings.Join(argv, " "), func(t *testing.T) {
				result := runAXICommandAt(t, root, argv)
				if result.code != 2 || result.stdout != want || result.stderr != "" {
					t.Fatalf("bench %q = (%d, %q, %q), want (2, %q, \"\")", argv, result.code, result.stdout, result.stderr, want)
				}
			})
		}
	}
	if refused == 0 {
		t.Fatal("the command registry holds no repository-scoped verb")
	}
	t.Run("idea", func(t *testing.T) {
		result := runAXICommandAt(t, root, []string{"idea", "--in", "primary", "x"})
		if result.code != 2 {
			t.Fatalf("idea --in primary x = (%d, %q, %q), want exit 2", result.code, result.stdout, result.stderr)
		}
		if _, err := os.Stat(filepath.Join(root, "capture", "IDEAS.md")); !errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("idea --in primary x left capture/IDEAS.md: %v", err)
		}
	})
}
