package spec

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/bounds"
)

const flatLayoutInstruction = "flat spec layout: move to specs/<slug>/spec.md"

// Resolve finds the readable file backing a spec argument. It tries the argument as
// given first, so a same-named readable file shadows the fallback, then — for a
// separator-free argument only — the folder fallback. The function refuses a live flat
// spec with an explicit migration diagnostic before it considers a folder.
//
// base anchors those fallbacks. Pass the repo root to resolve base repo-root-relative,
// so a cwd deeper than the root still finds it, or pass "" to resolve it relative to
// the process cwd. ok is false when no form resolves; tried holds every form
// attempted, for the not-found error. A non-nil err is a read failure on an existing
// file, for example a permissions error, reported instead of not-found.
func Resolve(base, arg string) (content []byte, resolved string, tried []string, ok bool, err error) {
	return resolve(base, arg, false)
}

// ResolveNoFollow resolves a producer spec without following a final symbolic link.
func ResolveNoFollow(base, arg string) (content []byte, resolved string, tried []string, ok bool, err error) {
	return resolve(base, arg, true)
}

func resolve(base, arg string, noFollow bool) (content []byte, resolved string, tried []string, ok bool, err error) {
	tried = []string{arg}
	if b, err := readCandidate(arg, noFollow); err != nil || b != nil {
		return b, arg, tried, err == nil, err
	}
	if !strings.ContainsRune(arg, '/') {
		slug := strings.TrimSuffix(arg, ".md")
		flat := filepath.Join(base, "specs", slug+".md")
		folder := filepath.Join(base, filepath.FromSlash(LiveSpecPath(arg)))
		tried = append(tried, folder)
		if pathExists(flat) {
			if pathExists(filepath.Dir(folder)) {
				return nil, flat, tried, false, fmt.Errorf("%s; found flat %s and folder form %s", flatLayoutInstruction, flat, folder)
			}
			return nil, flat, tried, false, fmt.Errorf("%s; found flat %s; target %s", flatLayoutInstruction, flat, folder)
		}
		if pathExists(filepath.Dir(folder)) && !pathExists(folder) {
			return nil, folder, tried, false, fmt.Errorf("spec folder is missing %s", folder)
		}
		if b, err := readCandidate(folder, noFollow); err != nil || b != nil {
			return b, folder, tried, err == nil, err
		}
	}
	return nil, "", tried, false, nil
}

func pathExists(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}

// readCandidate reads path as a candidate spec through the classifier, so a FIFO or
// device parked where a spec belongs is rejected before the open rather than blocking
// forever. An absent path or a directory is not a candidate; it returns (nil, nil) so
// the caller can try the next form. Any other failure is a real error to surface,
// never masked as not-found.
func readCandidate(path string, noFollow bool) ([]byte, error) {
	var c bounds.Classified
	if noFollow {
		c = bounds.ClassifyNoFollow(path)
	} else {
		c = bounds.Classify(path, bounds.ControlRecordLimit)
	}
	switch {
	case c.State == bounds.StateAbsent:
		return nil, nil
	case c.State == bounds.StateWrongType && isDir(path, noFollow):
		return nil, nil
	case c.State.Failed():
		if noFollow {
			return nil, fmt.Errorf("%s: %s", c.State, c.Reason)
		}
		return nil, errors.New(c.Reason)
	}
	return c.Data, nil
}

// isDir separates the one non-regular path that means "keep resolving" from the ones
// that mean "this read failed." A directory named like the candidate simply is not the
// spec file. A special file where a spec belongs is a problem the caller must hear
// about.
func isDir(path string, noFollow bool) bool {
	var (
		fi  os.FileInfo
		err error
	)
	if noFollow {
		fi, err = os.Lstat(path)
	} else {
		fi, err = os.Stat(path)
	}
	return err == nil && fi.IsDir()
}
