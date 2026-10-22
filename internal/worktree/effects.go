// This file is the effect boundary: the one place the worktree package reads
// ambient process context — environment variables, the clock, the current directory,
// and the user's home.
// Public commands call these adapters once at their boundary and pass the resolved
// values down explicitly. Lower owners never read the environment, the current
// directory, or the clock themselves; the source census test enforces that split.

package worktree

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/git"
)

// currentTime is the boundary clock read. A command resolves one instant here and
// passes it down, so decision code can be tested with an injected time.
func currentTime() time.Time { return time.Now() }

// baseEnviron is the boundary read of the child base environment. newAmbient carries it
// through a verb run. RunTreeChild holds no ambient value, so it reads it here for its
// chain through runChild to execEnv.
func baseEnviron() []string { return os.Environ() }

// ambient is the process context of one verb run: the Bench home and the warnings writer
// that the verb entry received, the kit value, one instant, the git run value that holds
// the run's location answers, the call directory that a relative operand joins, or the
// error of its failed read, and the base environment of each child that the verb starts.
// Each internal verb form takes it second, after the joins value, so nothing below the
// entry reads the process.
type ambient struct {
	home     string
	kit      string
	now      time.Time
	warnings io.Writer
	run      *git.Run
	dir      string
	dirErr   error
	environ  []string
}

// newAmbient reads the kit value, the clock, the process directory, and the environment
// once, and it builds a fresh run value, so no location answer outlives the verb run. The
// caller supplies the home and the warnings writer.
func newAmbient(home string, stderr io.Writer) ambient {
	dir, dirErr := filepath.Abs(".")
	return ambient{home: home, kit: gate.KitValue(), now: currentTime(), warnings: stderr, run: &git.Run{}, dir: dir, dirErr: dirErr, environ: baseEnviron()}
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
