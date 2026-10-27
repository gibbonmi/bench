// Package carry decides whether an owed chunk requirement carries its last passing entry
// across the delta to the chunk source, so a repair session reruns only the reached checks.
package carry

import (
	"path"
	"strings"

	"github.com/gibbonmi/bench/internal/gate"
	rr "github.com/gibbonmi/bench/internal/reviewrecord"
	"github.com/gibbonmi/bench/internal/testreport"
)

// The reasons that a requirement cannot carry, in the order of the decision steps.
const (
	noCurrentEntry = "no current passing entry"
	notTestOnly    = "the delta is not test-only"
	wholeOnly      = "the delta changes a whole-only test file"
	reached        = "the delta reaches the command"
)

// Decide answers the carry candidate of requirement in chunk, the chunk entry after the
// freeze, or the reason that requirement cannot carry; spec names the record that the source
// digest excludes. The reason is empty when the requirement carries, and the error reports a
// failed query. The first step that fails gives the reason.
func Decide(root, spec string, chunk rr.Chunk, requirement rr.Requirement) (rr.Verification, string, error) {
	candidate, ok, err := rr.CarryCandidate(root, spec, chunk, requirement)
	if err != nil {
		return rr.Verification{}, "", err
	}
	if !ok {
		return rr.Verification{}, noCurrentEntry, nil
	}
	// An unchanged source carries before the reach, which holds the conformance operand for
	// every changed source.
	if candidate.SourceDigest == chunk.SourceDigest {
		return candidate, "", nil
	}
	// The list refuses a delta that it cannot prove test-only, an unreadable one included, so
	// every refusal reruns the requirement.
	files, err := rr.TestOnlyFiles(root, spec, candidate.SourceDigest, chunk.SourceDigest)
	if err != nil {
		return rr.Verification{}, notTestOnly, nil
	}
	reach, whole := gate.DeltaReach(root, files)
	if whole {
		return rr.Verification{}, wholeOnly, nil
	}
	if reaches(root, requirement.Command, reach) {
		return rr.Verification{}, reached, nil
	}
	return candidate, "", nil
}

// reaches reports whether reach, the package operands of a delta, reaches command. Two forms
// of one simple `bench test` call that runs a test can escape: the system check, with or
// without a `--run` filter, and a package run that is not `--changed` whose expressions are
// canonical and cover no operand of reach. Every other command is reached, each other named
// check included.
func reaches(root, command string, reach []string) bool {
	selection, ok := testreport.PlanSelection(root, command)
	if !ok || !selection.Runs() {
		return true
	}
	if selection.Form() == testreport.CheckForm {
		return selection.Target() != gate.SystemPhaseName
	}
	if selection.Changed() {
		return true
	}
	for _, expression := range selection.Packages() {
		if !canonical(expression) {
			return true
		}
		for _, operand := range reach {
			if testreport.PackageCovers(expression, operand) {
				return true
			}
		}
	}
	return false
}

// canonical reports whether expression is `./...`, `./.`, `./<path>`, or `./<path>/...`,
// where <path> is clean, relative, and not `.`. A second spelling of a package could hide
// it from the cover test.
func canonical(expression string) bool {
	if expression == "./..." || expression == "./." {
		return true
	}
	rest, ok := strings.CutPrefix(expression, "./")
	if !ok {
		return false
	}
	rest = strings.TrimSuffix(rest, "/...")
	return rest != "." && !path.IsAbs(rest) && path.Clean(rest) == rest
}
