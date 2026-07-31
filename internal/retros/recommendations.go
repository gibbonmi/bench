package retros

import (
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
