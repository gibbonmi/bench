package consumers

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
)

// citationFields is the disclosure schema every success response carries: the commit the
// answer was read at, whether that checkout was replayable, the binary version, the
// replay spelling, and the hash of the answer itself.
var citationFields = []string{"sha", "state", "version", "cmd", "hash"}

// citation is the replay identity of one command run: the checkout it read, the bench
// version that read it, and the argv that asked.
type citation struct {
	root    string
	version string
	args    []string
}

// row renders the disclosure for the response bytes emitted so far. internal/toon has no
// single-object form, so a one-row table is the tree convention for a fixed-shape record.
// The hash covers every prior byte, so a reviewer recomputes it from the printed answer
// without knowing which blocks the response chose.
func (c citation) row(prior string) (string, error) {
	sha, err := git.Output("-C", c.root, "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256([]byte(prior))
	return toon.TableTyped("citation", citationFields, [][]any{{
		sha, c.state(), c.version, c.cmd(), hex.EncodeToString(sum[:]),
	}})
}

// state reports whether the checkout is replayable. Any porcelain entry, an untracked
// file included, makes it dirty, and a status read that fails at all reports dirty too:
// an undeterminable checkout must never promise a replay that cannot match.
func (c citation) state() string {
	_, changes, err := git.AllFilesStatus(c.root)
	if err != nil || len(changes) > 0 {
		return "dirty"
	}
	return "clean"
}

// cmd is the replay spelling of the run. usage.Parse already refused an argument TOON
// cannot carry, so the join asserts nothing more about the cells.
func (c citation) cmd() string {
	return strings.Join(append([]string{"bench", "consumers"}, c.args...), " ")
}
