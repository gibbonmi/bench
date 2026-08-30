// This file is the effect boundary: the one place the worktree package reads
// ambient process context — environment variables, the clock, and the user's home.
// Public commands call these adapters once at their boundary and pass the resolved
// values down explicitly. Lower owners never read the environment, the current
// directory, or the clock themselves; the source census test enforces that split.

package worktree

import (
	"os"
	"time"

	"github.com/gibbonmi/bench/internal/benchhome"
)

// currentTime is the boundary clock read. A command resolves one instant here and
// passes it down, so decision code can be tested with an injected time.
func currentTime() time.Time { return time.Now() }

// homeEnv is this package's name for the Bench home variable. internal/benchhome
// declares the string, and a verb writes it onto every child it starts.
const homeEnv = benchhome.Env

// Home resolves the Bench home directory. internal/benchhome owns the read, so this
// function holds none of its own. A command boundary resolves it once and passes the
// value down; a caller in another package resolves it at its own boundary and hands
// it to the verb.
func Home() string { return benchhome.Dir() }

// subshellShell resolves the interactive shell the subshell command launches.
func subshellShell() string { return os.Getenv("SHELL") }

// landingAlreadyRebuilt reports whether this landing process already ran once under
// a rebuild marker, so a second stale verdict refuses instead of rebuilding again.
func landingAlreadyRebuilt() bool { return os.Getenv(rebuiltLandingEnv) != "" }
