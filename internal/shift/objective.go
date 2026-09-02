package shift

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/gibbonmi/bench/internal/sanitize"
)

// objective is the one owner of the reviewer-authored shift objective. It is admitted
// once through validateObjective and then hands out one projection per surface. No
// caller decides how the text is escaped, bounded, or written. The banner and the
// durable commit subject share the sanitizer's escaped preview, bounded at
// bounds.PreviewRuneLimit code points. The
// prompt, the predicate argument, and the scratch bytes carry the text verbatim. The
// scratch file's 0600 mode is the only place that verbatim text persists.
type objective string

// objectiveMaxRunes caps an objective in Unicode code points, not bytes, so a
// multibyte objective is not cut at a fraction of its apparent length. The text
// flows into a commit subject, a scratch file, and status; an unbounded objective
// would carry into all three at once.
const objectiveMaxRunes = 200

// hasControlByte reports whether s carries any byte below 0x20 or the DEL byte 0x7f.
// This is a stricter, single-purpose check than toon.Representable's cell-escaping
// notion, which tolerates \n/\r/\t inside an already-escaped TOON cell. A shift
// objective is one line of operator intent, never a pre-escaped cell, so every control
// byte is rejected outright, with no exceptions.
func hasControlByte(s string) bool {
	for i := 0; i < len(s); i++ {
		if b := s[i]; b < 0x20 || b == 0x7f {
			return true
		}
	}
	return false
}

// validateObjective rejects an empty or whitespace-only objective; the "improve the
// codebase" default is gone, so an operator must state one. It also rejects an objective
// carrying a control byte. Hostile or accidental text is refused at entry, before it can
// reach the intent ledger or the TOON emitter.
func validateObjective(objective string) error {
	if strings.TrimSpace(objective) == "" {
		return errors.New("objective required: bench shift <objective...>")
	}
	if hasControlByte(objective) {
		return errors.New("objective contains a control byte")
	}
	if n := utf8.RuneCountInString(objective); n > objectiveMaxRunes {
		return fmt.Errorf("objective too long: %d runes, maximum %d", n, objectiveMaxRunes)
	}
	return nil
}

// objectiveBanner formats the shift-start line. The objective is reviewer-authored and
// already rejected at intake if it carries a control byte, but the banner renders it
// through the shared sanitizer anyway. This is one policy for every terminal render of
// operator-influenced text, with no render path carrying a raw escape sequence. The
// " — objective:" delimiter is load-bearing: the shift-start parser splits the branch on it.
func objectiveBanner(branch, objective string) string {
	return fmt.Sprintf("▶ shift on %s — objective: %s", branch, sanitize.Preview(objective))
}

// banner is the shift-start line for branch.
func (o objective) banner(branch string) string {
	return objectiveBanner(branch, string(o))
}

// commitSubject is the durable subject for iteration i. It carries the same escaped,
// bounded preview as the banner. History is a terminal render too, and it outlives the
// shift, so it never carries the verbatim text.
func (o objective) commitSubject(i int) string {
	return fmt.Sprintf("shift: iteration %d — %s", i, sanitize.Preview(string(o)))
}

// prompt is the iteration text written to the adapter's stdin, verbatim.
func (o objective) prompt() string {
	return fmt.Sprintf(iterationPrompt, string(o))
}

// predicateArgument is the single argument `.bench/done.sh` receives, verbatim.
func (o objective) predicateArgument() string {
	return string(o)
}

// scratch is the byte content of the worktree's `.bench-objective` file, verbatim
// with a trailing newline; the caller writes it 0600.
func (o objective) scratch() []byte {
	return []byte(string(o) + "\n")
}
