package retros

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/markdown"
)

// Recommendation is one drainable paragraph or list item and its source line.
type Recommendation struct {
	Body string
	Line int
}

// Recommendations returns the recommendation units beneath Agent-experience improvements.
// The block reader names the H2 headings, so a fenced heading leaves the section open.
func Recommendations(content []byte) []Recommendation {
	lines := markdown.Read(string(content)).Lines
	inImprovements := false
	var out []Recommendation
	for i := 0; i < len(lines); {
		line := lines[i]
		if line.Heading {
			inImprovements = line.Section == improvementsTitle
			i++
			continue
		}
		if !inImprovements || !unitText(line) {
			i++
			continue
		}
		body, listed := listItem(line.Raw)
		if !listed {
			body = strings.TrimSpace(line.Raw)
		}
		i++
		for i < len(lines) && !lines[i].Heading && unitText(lines[i]) {
			if _, nextListed := listItem(lines[i].Raw); nextListed {
				break
			}
			body += "\n" + strings.TrimSpace(lines[i].Raw)
			i++
		}
		out = append(out, Recommendation{Body: body, Line: line.Number})
	}
	return out
}

// unitText reports whether a line that is no H2 heading holds unit text. A fence marker line,
// a fenced line, an H3 line, and a blank line hold none, so each one ends the open unit.
func unitText(line markdown.Line) bool {
	return line.Class == markdown.Body && !strings.HasPrefix(line.Raw, "### ") && strings.TrimSpace(line.Raw) != ""
}

// feedsMarker is the destination grammar an improvement item ends with: a roadmap row
// ID, the row a drain has to open, or the absence of a row. The anchors are deliberate —
// a value with leading or trailing text is a different sentence, not a marker.
var feedsMarker = regexp.MustCompile(`^Feeds: (FT[1-9][0-9]*|new|none)$`)

// FeedsMarked reports whether the unit ends with one well-formed Feeds line. The marker
// has to be the last line: a person repairing a retro then reads every destination in one
// fixed place, and a Feeds sentence buried mid-item never passes for one.
func (r Recommendation) FeedsMarked() bool {
	lines := strings.Split(r.Body, "\n")
	return feedsMarker.MatchString(lines[len(lines)-1])
}

func listItem(line string) (string, bool) {
	trimmed := strings.TrimLeft(line, " \t")
	for _, marker := range []string{"- ", "* ", "+ "} {
		if after, ok := strings.CutPrefix(trimmed, marker); ok {
			return strings.TrimSpace(after), true
		}
	}
	period := strings.IndexByte(trimmed, '.')
	if period <= 0 || period+1 >= len(trimmed) || trimmed[period+1] != ' ' {
		return "", false
	}
	if _, err := strconv.Atoi(trimmed[:period]); err != nil {
		return "", false
	}
	return strings.TrimSpace(trimmed[period+2:]), true
}
