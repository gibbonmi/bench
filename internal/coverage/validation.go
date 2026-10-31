package coverage

import (
	"strconv"
	"strings"
)

// storyPart is one comma-separated part of a row's story reference: a single number, or a
// range from start to end. The text fields keep the authored digits, which a violation
// message quotes.
type storyPart struct {
	start, end         int
	startText, endText string
	ranged             bool
}

// storyReference reads one row's story cell through the one story-reference grammar. It
// answers the reference without a trailing parenthetical, and its parts. An empty reference
// and an edge reference have no parts. ok is false for a reference that the grammar does not
// recognize. Check and Describe both read a row's stories here.
func storyReference(cell string) (ref string, parts []storyPart, ok bool) {
	ref = parenRe.ReplaceAllString(cell, "")
	if ref == "" || edgeRe.MatchString(ref) {
		return ref, nil, true
	}
	if !storyRefRe.MatchString(ref) {
		return ref, nil, false
	}
	for _, text := range strings.Split(ref, ",") {
		m := storyPartRe.FindStringSubmatch(trimSpaceTab(text))
		part := storyPart{startText: m[1], endText: m[2], ranged: m[2] != ""}
		part.start, _ = strconv.Atoi(m[1])
		part.end = part.start
		if part.ranged {
			part.end, _ = strconv.Atoi(m[2])
		}
		parts = append(parts, part)
	}
	return ref, parts, true
}

// storyNumbers answers every story number that the parts reference, in authored order. A
// range with its end before its start references none.
func storyNumbers(parts []storyPart) []int {
	numbers := []int{}
	for _, part := range parts {
		for n := part.start; n <= part.end; n++ {
			numbers = append(numbers, n)
		}
	}
	return numbers
}
