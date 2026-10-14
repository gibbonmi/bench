package ownership

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"

	"github.com/gibbonmi/bench/internal/consumers"
	"github.com/gibbonmi/bench/internal/git"
)

// beforeState is one changed declaration at the source base: the line that declares it in
// its owner file, zero when that file declares no such name, and every base file that
// uses it.
type beforeState struct {
	line int
	uses []consumers.SourceUse
}

// before reads the pinned before-state of one changed declaration through the consumers
// base-side reader. A relocation or a deletion removes the declaration from its owner file,
// so only the base still shows the files that used it there. The candidate files are the
// base Go sources whose bytes hold the name; the reader then keeps only the files that name
// the declaration itself. An absent source base reads nothing.
func (c *census) before(change Change) (beforeState, error) {
	name := localName(change.Owner.Symbol)
	if c.source.Base == "" || name == "" {
		return beforeState{}, nil
	}
	candidates, err := c.baseCandidates(name)
	if err != nil {
		return beforeState{}, err
	}
	sources := map[string]string{}
	for _, file := range append(candidates, change.Owner.Path) {
		if _, read := sources[file]; read {
			continue
		}
		data, err := git.ReadTreeFile(c.root, c.source.Base, file)
		if err != nil {
			if file == change.Owner.Path {
				return beforeState{}, nil
			}
			return beforeState{}, err
		}
		sources[file] = string(data)
	}
	line, uses := consumers.SourceUses(sources, change.Owner.Path, change.Owner.Package, name)
	return beforeState{line: line, uses: uses}, nil
}

// baseCandidates is every Go source at the source base whose bytes hold name. The fixed
// string search is a superset of the files that name the declaration, so it drops no use.
func (c *census) baseCandidates(name string) ([]string, error) {
	return grepFiles(c.root, c.source.Base, name)
}

// grepFiles is every Go source at rev whose bytes hold one of needles. git grep exits 1 when
// no file matches, which is an empty answer.
func grepFiles(root, rev string, needles ...string) ([]string, error) {
	args := []string{"-C", root, "grep", "-z", "-l", "-F"}
	for _, needle := range needles {
		args = append(args, "-e", needle)
	}
	out, err := git.Raw(append(args, rev, "--", "*.go")...)
	var exit *exec.ExitError
	if errors.As(err, &exit) && exit.ExitCode() == 1 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var files []string
	for _, entry := range bytes.Split(bytes.TrimSuffix(out, []byte{0}), []byte{0}) {
		if file, ok := strings.CutPrefix(string(entry), rev+":"); ok {
			files = append(files, file)
		}
	}
	return files, nil
}
