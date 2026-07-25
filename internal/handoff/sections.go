package handoff

import (
	"strings"

	"github.com/gibbonmi/bench/internal/toon"
)

// stateHeading is the reviewer-owned section's heading. It is the one section the command
// preserves rather than derives, so locating it exactly is what stands between a rewrite
// and data loss.
const stateHeading = "## State"

// headingPrefix marks a level-two heading. The State body ends at the next unfenced one.
const headingPrefix = "## "

// refusal is a structured error that carries the AXI kind/hint pair, so the command
// surfaces one rendering of a refusal rather than re-deriving the line at each exit.
type refusal struct{ kind, hint string }

func (r refusal) Error() string { return toon.Errorf(r.kind, r.hint) }

// splitState returns the body of the document's State section, or refuses.
//
// It refuses on ambiguity rather than choosing: no unfenced heading, more than one, or any
// occurrence inside a fenced block. A fenced heading is prose about the document, not a
// section of it, and an implementation that cannot tell the two apart writes over whatever
// the reviewer had below the real one.
//
// The returned body has its surrounding blank lines stripped and its interior bytes
// untouched. The writer re-attaches exactly one blank line on each side, which is what
// makes a second run byte-identical: any spacing the previous run emitted is normalized
// back to the same form rather than accumulating.
func splitState(content []byte) (string, error) {
	lines := strings.Split(string(content), "\n")

	fenced := false
	unfenced, inFence, start := 0, 0, -1
	for i, raw := range lines {
		if isFence(raw) {
			fenced = !fenced
			continue
		}
		if !isStateHeading(raw) {
			continue
		}
		if fenced {
			inFence++
			continue
		}
		unfenced++
		if start < 0 {
			start = i
		}
	}

	switch {
	case unfenced == 0 && inFence == 0:
		return "", refusal{"handoff has no " + stateHeading + " section",
			"add a " + stateHeading + " heading, or delete the file so the command scaffolds a fresh one"}
	case unfenced != 1 || inFence > 0:
		return "", refusal{"ambiguous " + stateHeading + " section",
			"leave exactly one unfenced " + stateHeading + " heading; the command will not choose between them"}
	}

	fenced = false
	end := len(lines)
	for i := start + 1; i < len(lines); i++ {
		if isFence(lines[i]) {
			fenced = !fenced
			continue
		}
		if !fenced && strings.HasPrefix(lines[i], headingPrefix) {
			end = i
			break
		}
	}
	return strings.Trim(strings.Join(lines[start+1:end], "\n"), "\n"), nil
}

// isFence reports whether a line opens or closes a fenced block. Both markdown fence
// characters count, and leading indentation does not disqualify one.
func isFence(line string) bool {
	trimmed := strings.TrimLeft(line, " \t")
	return strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~")
}

// isStateHeading matches the heading exactly, tolerating only trailing whitespace so a
// CRLF document and a padded line resolve to the same section.
func isStateHeading(line string) bool {
	return strings.TrimRight(line, " \t\r") == stateHeading
}
