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
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".bench")
}
