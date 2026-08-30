// Fence parsing. The `## Ownership fences` section of a spec authorizes the paths a
// build may write, and two packages read it: internal/preflight bootstraps its
// authorization facts from the token list, and internal/coverage grades the section
// for the review pickup. This file is the one parser for that grammar, so neither
// consumer can disagree with the other about which tokens a section declares.
package spec

import (
	"regexp"
	"strings"
)

// fencesHeading is the exact line that opens the section. fencesEndRe bounds it the
// same way internal/coverage bounds `### Acceptance coverage map`: a level-2-or-deeper
// heading ends it, so an entry under a subsection sits outside the section.
const fencesHeading = "## Ownership fences"

var fencesEndRe = regexp.MustCompile(`^#{2,} `)

// FenceTokens extracts every backticked token in the `## Ownership fences` section
// that is not inside parentheses; parenthetical prose is annotation, never
// authorization. declared reports whether the section heading appeared at all, which
// separates a spec that declares no fences from one whose section is empty.
//
// Paren depth and backtick state carry across line boundaries. A parenthetical that
// opens on one line and closes on a later one still shields every token inside it.
// Depth returns to zero once it closes, so a later real entry authorizes normally.
func FenceTokens(content []byte) (tokens []string, declared bool) {
	inSection := false
	depth := 0
	inTick := false
	depthAtOpen := 0
	var cur strings.Builder
	for _, raw := range strings.Split(string(content), "\n") {
		line := strings.TrimSuffix(raw, "\r")
		if strings.TrimSpace(line) == fencesHeading {
			inSection = true
			declared = true
			continue
		}
		if inSection && fencesEndRe.MatchString(line) {
			inSection = false
		}
		if inSection {
			fenceTokensInLine(line, &depth, &inTick, &depthAtOpen, &cur, &tokens)
		}
	}
	return tokens, declared
}

// fenceTokensInLine is one line's pass through the fence-section state machine. Paren
// depth, backtick state, and the token under construction are threaded in by pointer.
// The caller carries them across every line of the section this way.
//
// A backtick-quoted token is captured into tokens only when the depth at the moment
// its opening backtick appeared was zero. Inside an open paren, whether opened on this
// line or an earlier one, a token never authorizes.
func fenceTokensInLine(line string, depth *int, inTick *bool, depthAtOpen *int, cur *strings.Builder, tokens *[]string) {
	for _, r := range line {
		switch r {
		case '(':
			*depth++
		case ')':
			if *depth > 0 {
				*depth--
			}
		case '`':
			if *inTick {
				if *depthAtOpen == 0 {
					*tokens = append(*tokens, cur.String())
				}
				cur.Reset()
				*inTick = false
			} else {
				*inTick = true
				*depthAtOpen = *depth
			}
		default:
			if *inTick {
				cur.WriteRune(r)
			}
		}
	}
}
