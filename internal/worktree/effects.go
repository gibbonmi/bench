// This file is the effect boundary: the one place the worktree package reads
// ambient process context — environment variables, the clock, and the user's home.
// Public commands call these adapters once at their boundary and pass the resolved
// values down explicitly. Lower owners never read the environment, the current
// directory, or the clock themselves; the source census test enforces that split.

package worktree

import (
	"os"
	"path/filepath"
	"time"
)

// currentTime is the boundary clock read. A command resolves one instant here and
// passes it down, so decision code can be tested with an injected time.
func currentTime() time.Time { return time.Now() }

// benchHome resolves the Bench home directory from the process environment, with
// the user's home as the fallback. Commands resolve it once at the boundary.
func benchHome() string {
	if h := os.Getenv("BENCH_HOME"); h != "" {
		return h
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".bench")
}

// subshellShell resolves the interactive shell the subshell command launches.
func subshellShell() string { return os.Getenv("SHELL") }

// landingAlreadyRebuilt reports whether this landing process already ran once under
// a rebuild marker, so a second stale verdict refuses instead of rebuilding again.
func landingAlreadyRebuilt() bool { return os.Getenv(rebuiltLandingEnv) != "" }
