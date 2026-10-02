// This file is the effect boundary: the one place the worktree package reads
// ambient process context — environment variables, the clock, and the user's home.
// Public commands call these adapters once at their boundary and pass the resolved
// values down explicitly. Lower owners never read the environment, the current
// directory, or the clock themselves; the source census test enforces that split.

package worktree

import (
	"io"
	"os"
	"time"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/gate"
)

// currentTime is the boundary clock read. A command resolves one instant here and
// passes it down, so decision code can be tested with an injected time.
func currentTime() time.Time { return time.Now() }

// ambient is the process context of one verb run: the Bench home and the warnings writer
// that the verb entry received, the kit value, and one instant. Each internal verb form
// takes it second, after the joins value, so nothing below the entry reads the process.
type ambient struct {
	home     string
	kit      string
	now      time.Time
	warnings io.Writer
}

// newAmbient reads the kit value and the clock once. The caller supplies the home and the
// warnings writer.
func newAmbient(home string, stderr io.Writer) ambient {
	return ambient{home: home, kit: gate.KitValue(), now: currentTime(), warnings: stderr}
}

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
