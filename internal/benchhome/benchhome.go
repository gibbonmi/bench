// Package benchhome owns the one read of the Bench home from the process
// environment. It imports the standard library alone, so any package can resolve
// the home at its own boundary without an import cycle.
package benchhome

import (
	"os"
	"path/filepath"
)

// Env names the Bench home in a process environment. This package reads it here,
// and a verb writes it onto every child it starts, so both halves name it once.
const Env = "BENCH_HOME"

// Dir resolves the Bench home directory from the process environment, with the
// user's home as the fallback. It is the one BENCH_HOME read in the tree. A
// command boundary resolves it once and passes the value down.
func Dir() string {
	if h := os.Getenv(Env); h != "" {
		return h
	}
	return fallbackDir()
}

// Export writes the fallback Bench home into the process environment when
// BENCH_HOME is unset, as the shell wrapper does before it runs the binary. A
// shim can run the binary with no wrapper, and a gate that declares BENCH_HOME
// reads it from the environment, so the binary's entry point calls this once. An
// unresolvable user home exports nothing, so a declared-input check still refuses.
func Export() {
	if os.Getenv(Env) != "" {
		return
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		os.Setenv(Env, fallbackDir())
	}
}

// fallbackDir is the Bench home Dir supplies when BENCH_HOME is unset: the
// user's own home directory joined with .bench. Dir and IsFallback share this
// one join, so a caller never re-derives the fallback path a second way.
func fallbackDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".bench")
}

// IsFallback reports whether home names the user's own Bench home — the value
// Dir supplies when BENCH_HOME is unset — regardless of whether BENCH_HOME
// named that exact path or the fallback supplied it. A caller that must not
// write into the user's real home during a test run grades the resolved home
// here before it writes.
func IsFallback(home string) bool {
	return home == fallbackDir()
}
