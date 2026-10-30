package roadmap

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/gibbonmi/bench/internal/markdown"
)

// sequenceHeading opens the board's recommended sequence. The parser reads it, and
// ProjectSequence renders it.
const sequenceHeading = markdown.H2Prefix + "Recommended sequence"

func sequenceBounds(lines []string) (start, end int, hasSection bool) {
	return sectionBounds(lines, sequenceHeading)
}

// sectionBounds locates the first H2 section titled heading. The block reader names the H2
// headings, so a fenced heading neither starts nor ends a section. The title match reads the
// raw heading line with its trailing spaces and tabs cut. start is the heading line, or -1
// when the section is absent, and end is the next H2 heading or the end of lines. hasSection
// reports any H2 heading.
func sectionBounds(lines []string, heading string) (start, end int, hasSection bool) {
	start, end = -1, len(lines)
	for idx, line := range markdown.Read(strings.Join(lines, "\n")).Lines {
		if !line.Heading {
			continue
		}
		hasSection = true
		if start < 0 && strings.TrimRight(line.Raw, " \t") == heading {
			start = idx
			continue
		}
		if start >= 0 {
			end = idx
			break
		}
	}
	return start, end, hasSection
}

// ProjectSequence replaces the recognized sequence with the ordered outcome identities.
func ProjectSequence(document []byte, outcomes []string) ([]byte, error) {
	lines := strings.Split(string(document), "\n")
	start, end, _ := sequenceBounds(lines)
	if start < 0 {
		return nil, errors.New("ROADMAP.md has no recommended sequence")
	}
	replacement := []string{sequenceHeading, ""}
	for index, outcome := range outcomes {
		if strings.TrimSpace(outcome) == "" || strings.ContainsAny(outcome, "\r\n") {
			return nil, fmt.Errorf("invalid commitment outcome %q", outcome)
		}
		replacement = append(replacement, fmt.Sprintf("%d. %s", index+1, outcome))
	}
	if (end < len(lines) || bytes.HasSuffix(document, []byte("\n"))) && replacement[len(replacement)-1] != "" {
		replacement = append(replacement, "")
	}
	projected := append([]string{}, lines[:start]...)
	projected = append(projected, replacement...)
	projected = append(projected, lines[end:]...)
	return []byte(strings.Join(projected, "\n")), nil
}

// SequenceText returns the canonical parser's complete recommended-sequence section.
func SequenceText(document []byte) string {
	_, text, _ := parseSequence(strings.Split(string(document), "\n"))
	return text
}
