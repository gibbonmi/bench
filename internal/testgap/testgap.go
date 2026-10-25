// Package testgap proves that a tree delta changes only Go test files.
package testgap

import (
	"errors"
	"fmt"

	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/outline"
)

var (
	ErrUnreadable   = errors.New("unreadable tree")
	ErrEmptyChanges = errors.New("different trees have no visible changes")
	ErrStatus       = errors.New("path is not an added or modified file")
	ErrMode         = errors.New("file mode is changed or nonregular")
	ErrNotTest      = errors.New("path is not a Go test file")
	ErrTestdata     = errors.New("path has a testdata element")
)

// Prove accepts equal trees or a delta of added or modified regular Go test files outside
// every testdata directory. A rename arrives as a deletion and an addition, so the
// deletion refuses it.
func Prove(root, reviewed, later string) error {
	_, err := TestFiles(root, reviewed, later)
	return err
}

// TestFiles answers the changed paths of the delta from reviewed to later in path order
// when Prove accepts the delta, and the refusal of Prove otherwise. Equal trees answer an
// empty list.
func TestFiles(root, reviewed, later string) ([]string, error) {
	files := []string{}
	if reviewed == later {
		return files, nil
	}
	changes, err := benchgit.TreeChangesIncludingSubmodules(root, reviewed, later)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnreadable, err)
	}
	if len(changes) == 0 {
		return nil, ErrEmptyChanges
	}
	for _, change := range changes {
		if err := proveFile(change); err != nil {
			return nil, fmt.Errorf("%s: %w", change.Path, err)
		}
		files = append(files, change.Path)
	}
	return files, nil
}

func proveFile(change benchgit.TreeChange) error {
	if change.Status != "A" && change.Status != "M" {
		return ErrStatus
	}
	if !(benchgit.IndexEntry{Mode: change.DstMode}).IsRegularFile() || (change.Status == "M" && change.SrcMode != change.DstMode) {
		return ErrMode
	}
	if !outline.IsGoTestFile(change.Path) {
		return ErrNotTest
	}
	// Go ignores a testdata directory, so a test file there is fixture data that a
	// production file can embed.
	if outline.IsFixture(change.Path) {
		return ErrTestdata
	}
	return nil
}
