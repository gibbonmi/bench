package axi

import (
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

func TestSubjectRootUsesBenchContractRootInAXISubpackage(t *testing.T) {
	alternateRoot := filepath.Join(t.TempDir(), "subject root [glob]*")
	contract.WriteExecutableAbs(t, filepath.Join(alternateRoot, "bin", "bench.sh"), `#!/usr/bin/env bash
printf 'alternate subject root:%s\n' "$1"
`)
	t.Setenv("BENCH_CONTRACT_ROOT", alternateRoot)

	gotRoot := contract.SubjectRoot(t)
	if gotRoot == contract.KitRoot(t) {
		t.Fatalf("SubjectRoot fell back to KitRoot %q", gotRoot)
	}
	if gotRoot != alternateRoot {
		t.Fatalf("SubjectRoot = %q, want %q", gotRoot, alternateRoot)
	}

	f := contract.NewFixture(t)
	out := f.BenchWrapper("learnings")

	out.RequireExit(0)
	out.RequireContains(out.Stdout, "alternate subject root:learnings")
}
