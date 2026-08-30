package benchhome

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDirReadsTheEnvironment holds the set case: an operator assignment wins over
// the fallback.
func TestDirReadsTheEnvironment(t *testing.T) {
	want := t.TempDir()
	t.Setenv(Env, want)
	if got := Dir(); got != want {
		t.Fatalf("Dir() = %q, want the assigned home %q", got, want)
	}
}

// TestDirFallsBackToTheUserHome holds the unset case: the home is the .bench
// directory under the user's home.
func TestDirFallsBackToTheUserHome(t *testing.T) {
	t.Setenv(Env, "")
	user, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("user home: %v", err)
	}
	want := filepath.Join(user, ".bench")
	if got := Dir(); got != want {
		t.Fatalf("Dir() = %q, want the fallback %q", got, want)
	}
}

// TestEnvNamesTheBenchHome holds the exported name the child-env injection writes.
func TestEnvNamesTheBenchHome(t *testing.T) {
	if Env != "BENCH_HOME" {
		t.Fatalf("Env = %q, want BENCH_HOME", Env)
	}
}
