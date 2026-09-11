package probe

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/benchhome"
	"github.com/gibbonmi/bench/internal/canonicalpath"
	"github.com/gibbonmi/bench/internal/poolkey"
	"github.com/gibbonmi/bench/internal/toon"
)

// subject is the one file a probe mutates: its absolute path, the repo-relative spelling
// the row and every refusal print, its mode, and the bytes read before any write. Those
// bytes are the restore oracle; the preserved copy is only the restore source.
type subject struct {
	path    string
	display string
	mode    fs.FileMode
	start   []byte
}

// resolveSubject answers the mutation target, or the one refusal line that names why the
// operand is not a regular tree file. The path resolves the way `bench anchors` resolves
// its operand: an absolute path stays, and a relative path joins onto the working
// directory. Lstat decides the kind, so a symlink refuses instead of the verb mutating
// whatever it points at.
func resolveSubject(root, arg string) (subject, string) {
	cwd := ""
	if !filepath.IsAbs(arg) {
		working, err := os.Getwd()
		if err != nil {
			return subject{}, unavailable(filepath.ToSlash(filepath.Clean(arg)), "absent")
		}
		cwd = working
	}
	path, display := canonicalpath.Operand(root, cwd, arg)
	if !withinRoot(root, path) {
		return subject{}, unavailable(display, "outside the repository")
	}
	info, err := os.Lstat(path)
	if err != nil {
		return subject{}, unavailable(display, "absent")
	}
	switch {
	case info.Mode()&fs.ModeSymlink != 0:
		return subject{}, unavailable(display, "a symlink")
	case info.IsDir():
		return subject{}, unavailable(display, "a directory")
	case !info.Mode().IsRegular():
		return subject{}, unavailable(display, "a special file")
	}
	start, err := os.ReadFile(path)
	if err != nil {
		return subject{}, unavailable(display, "unreadable")
	}
	return subject{path: path, display: display, mode: info.Mode().Perm(), start: start}, ""
}

func unavailable(display, reason string) string {
	return toon.Errorf("probe subject unavailable", display+" is "+reason) + "\n"
}

func withinRoot(root, path string) bool {
	relative, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	return relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}

// mutationResult carries mutation classification and bytes together, so callers cannot
// recount the subject to distinguish a miss from the exact-match mutation.
type mutationResult struct {
	mutated []byte
	matches int
}

// mutate classifies the exact-string match once and applies the mutation only for one
// match. Several matches refuse because replacing the first would mutate a site the
// caller did not name.
func mutate(start []byte, old, replacement, kind string) (mutationResult, string) {
	if kind == "swap" && replacement == old {
		return mutationResult{}, toon.Errorf("probe mutation empty", "--with equals --swap") + "\n"
	}
	count := bytes.Count(start, []byte(old))
	switch count {
	case 0:
		return mutationResult{}, ""
	case 1:
	default:
		hint := fmt.Sprintf("the old string matches %d times, want exactly 1", count)
		return mutationResult{}, toon.Errorf("probe mutation ambiguous", hint) + "\n"
	}
	if kind == "unwrap" {
		replacement, ok := unwrapCall(old)
		if !ok {
			return mutationResult{}, toon.Errorf("probe unwrap invalid", "--unwrap must name one outer call") + "\n"
		}
		return mutationResult{mutated: bytes.Replace(start, []byte(old), []byte(replacement), 1), matches: 1}, ""
	}
	return mutationResult{mutated: bytes.Replace(start, []byte(old), []byte(replacement), 1), matches: 1}, ""
}

func unwrapCall(call string) (string, bool) {
	open := strings.IndexByte(call, '(')
	if open <= 0 || !strings.HasSuffix(call, ")") {
		return "", false
	}
	depth := 0
	for i := open; i < len(call); i++ {
		switch call[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 && i != len(call)-1 {
				return "", false
			}
			if depth < 0 {
				return "", false
			}
		}
	}
	return call[open+1 : len(call)-1], depth == 0
}

// preservation is the copy a probe leaves under the Bench home for the run's span. dir
// is removed on a proven restore and kept on a failed one, so a caller whose restore
// failed always has a file to restore by hand. root names the verb's own subtree under
// the home, which is the outermost directory a release may remove.
type preservation struct {
	root string
	dir  string
	file string
}

// preserve writes the start bytes under the Bench home before any tree write. The stamp
// directory is unique per run, so two probes never share a copy, and the copy is 0600 in
// a 0700 directory, because it holds a verbatim copy of a repository file.
func preserve(root string, subject subject) (preservation, string) {
	owned := filepath.Join(benchhome.Dir(), "probe")
	dir := filepath.Join(owned, poolkey.Key(root), stamp())
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return preservation{}, preservationFailure(err)
	}
	file := filepath.Join(dir, filepath.Base(subject.path))
	if err := os.WriteFile(file, subject.start, 0o600); err != nil {
		_ = os.RemoveAll(dir)
		return preservation{}, preservationFailure(err)
	}
	return preservation{root: owned, dir: dir, file: file}, ""
}

func preservationFailure(err error) string {
	return toon.Errorf("probe preservation failed", err.Error()) + "\n"
}

// release removes a proven-restored copy and then every parent it emptied, so a
// completed probe adds no entry to the Bench home at all.
func (p preservation) release() {
	_ = os.RemoveAll(p.dir)
	for parent := filepath.Dir(p.dir); strings.HasPrefix(parent, p.root); parent = filepath.Dir(parent) {
		if err := os.Remove(parent); err != nil {
			return
		}
	}
}

// replaceAtomic writes data over path through a sibling temporary file and one rename,
// so a crash between the two writes leaves the subject whole rather than half-written.
// The temporary carries the subject's mode, so an executable subject stays executable.
func replaceAtomic(path string, data []byte, mode fs.FileMode) error {
	temporary, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".probe-")
	if err != nil {
		return err
	}
	name := temporary.Name()
	defer os.Remove(name)
	if _, err := temporary.Write(data); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	if err := os.Chmod(name, mode); err != nil {
		return err
	}
	return os.Rename(name, path)
}

// restore writes the preserved copy back and proves the result byte-exact against the
// bytes read at the start. The comparison is against those start bytes, never against
// the copy: a copy that changed under the run would otherwise certify its own damage.
func restore(subject subject, preserved preservation) (bool, string) {
	saved, err := os.ReadFile(preserved.file)
	if err != nil {
		return false, err.Error()
	}
	if err := replaceAtomic(subject.path, saved, subject.mode); err != nil {
		return false, err.Error()
	}
	back, err := os.ReadFile(subject.path)
	if err != nil {
		return false, err.Error()
	}
	if !bytes.Equal(back, subject.start) {
		return false, "the restored bytes differ from the bytes read at the start"
	}
	return true, ""
}
