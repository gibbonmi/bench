// Package commentgap proves that a tree delta contains only inert Go changes.
package commentgap

import (
	"errors"
	"fmt"
	"strings"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

var (
	ErrUnreadable    = errors.New("unreadable tree or blob")
	ErrEmptyChanges  = errors.New("different trees have no visible changes")
	ErrStatus        = errors.New("path is not an existing modified file")
	ErrMode          = errors.New("file mode is changed or nonregular")
	ErrNotGo         = errors.New("path is not a Go file")
	ErrScan          = errors.New("Go source cannot scan or parse")
	ErrTokens        = errors.New("Go tokens differ")
	ErrDirective     = errors.New("file contains a directive comment")
	ErrCgo           = errors.New("file imports C")
	ErrExampleOutput = errors.New("test file contains example output")
)

// Prove accepts equal trees or a delta of existing regular Go files with inert changes.
func Prove(root, reviewed, later string) error {
	if reviewed == later {
		return nil
	}
	changes, err := benchgit.TreeChangesIncludingSubmodules(root, reviewed, later)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUnreadable, err)
	}
	if len(changes) == 0 {
		return ErrEmptyChanges
	}
	for _, change := range changes {
		if err := proveFile(root, reviewed, later, change); err != nil {
			return fmt.Errorf("%s: %w", change.Path, err)
		}
	}
	return nil
}

func proveFile(root, reviewed, later string, change benchgit.TreeChange) error {
	if change.Status != "M" {
		return ErrStatus
	}
	if change.SrcMode != change.DstMode || !(benchgit.IndexEntry{Mode: change.SrcMode}).IsRegularFile() {
		return ErrMode
	}
	if !strings.HasSuffix(change.Path, ".go") {
		return ErrNotGo
	}
	before, err := benchgit.ReadTreeFile(root, reviewed, change.Path)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUnreadable, err)
	}
	after, err := benchgit.ReadTreeFile(root, later, change.Path)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUnreadable, err)
	}
	return proveGo(change.Path, before, after)
}
