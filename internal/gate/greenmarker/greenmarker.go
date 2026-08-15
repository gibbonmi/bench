// Package greenmarker is the one reader and writer of the project-green marker,
// `refs/bench/green/<branch>` — the gate's durable record of the branch tip whose
// exact tree last held green. Every consumer resolves the marker through Read, so
// the peel and the dangling-symbolic-ref classification are decided once, and every
// writer advances it through Advance, so the compare-and-swap lives once.
package greenmarker

import (
	"errors"
	"fmt"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

// zeroOID is the compare-and-swap expectation for a marker that does not exist yet.
const zeroOID = "0000000000000000000000000000000000000000"

// Ref names the project-green marker for branch.
func Ref(branch string) string {
	return "refs/bench/green/" + branch
}

// Read resolves the marker to the full commit it peels to. present is false when the
// marker does not exist; a marker that exists but does not resolve — a dangling
// symbolic ref, or an object that does not peel to a commit — is present with an
// error, never mistaken for absence.
func Read(root, branch string) (commit string, present bool, err error) {
	marker := Ref(branch)
	if _, err := benchgit.Output("-C", root, "rev-parse", "--verify", marker); err != nil {
		if _, symbolicErr := benchgit.Raw("-C", root, "symbolic-ref", "--quiet", marker); symbolicErr == nil {
			return "", true, errors.New("read project-green marker: dangling symbolic ref")
		}
		return "", false, nil
	}
	commit, err = benchgit.Output("-C", root, "rev-parse", "--verify", marker+"^{commit}")
	if err != nil {
		return "", true, fmt.Errorf("read project-green marker: %w", err)
	}
	return commit, true, nil
}

// Advance moves the marker to destination only if it currently reads expected, where
// an empty expected means the marker must not exist yet. A lost race that still left
// the marker at destination is success.
func Advance(root, branch, destination, expected string) error {
	if expected == "" {
		expected = zeroOID
	}
	if _, err := benchgit.Raw("-C", root, "update-ref", Ref(branch), destination, expected); err != nil {
		if existing, _, readErr := Read(root, branch); readErr == nil && existing == destination {
			return nil
		}
		return fmt.Errorf("record project-green marker: %w", err)
	}
	return nil
}
