package retros

import (
	"regexp"
	"strconv"
	"strings"
)

// Recommendation is one drainable paragraph or list item and its source line.
type Recommendation struct {
	Body string
	Line int
}

// Recommendations returns the recommendation units beneath Agent-experience improvements.
func Recommendations(content []byte) []Recommendation {
	lines := strings.Split(string(content), "\n")
	inImprovements := false
	var out []Recommendation
	for i := 0; i < len(lines); {
		line := strings.TrimSuffix(lines[i], "\r")
		if strings.HasPrefix(line, "## ") {
			inImprovements = strings.TrimSpace(strings.TrimPrefix(line, "## ")) == "Agent-experience improvements"
			i++
			continue
		}
		if !inImprovements || strings.HasPrefix(line, "### ") || strings.TrimSpace(line) == "" {
			i++
			continue
		}
		body, listed := listItem(line)
		start := i
		if !listed {
			body = strings.TrimSpace(line)
		}
		i++
		for i < len(lines) {
			next := strings.TrimSuffix(lines[i], "\r")
			if strings.HasPrefix(next, "## ") || strings.HasPrefix(next, "### ") || strings.TrimSpace(next) == "" {
				break
			}
			if _, nextListed := listItem(next); nextListed {
				break
			}
			body += "\n" + strings.TrimSpace(next)
			i++
		}
		out = append(out, Recommendation{Body: body, Line: start + 1})
	}
	return out
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
		if strings.HasPrefix(trimmed, marker) {
			return strings.TrimSpace(strings.TrimPrefix(trimmed, marker)), true
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
