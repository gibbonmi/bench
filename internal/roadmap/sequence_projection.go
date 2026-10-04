package roadmap

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

func sequenceBounds(lines []string) (start, end int, hasSection bool) {
	inFence := false
	start, end = -1, len(lines)
	for idx, line := range lines {
		trimmed := strings.TrimRight(line, " \t\r")
		if strings.HasPrefix(trimmed, "```") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		if strings.HasPrefix(trimmed, "## ") {
			hasSection = true
		}
		if start < 0 && trimmed == "## Recommended sequence" {
			start = idx
			continue
		}
		if start >= 0 && strings.HasPrefix(trimmed, "## ") {
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
	replacement := []string{"## Recommended sequence", ""}
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
