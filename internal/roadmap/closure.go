package roadmap

import (
	"fmt"
	"slices"
	"strings"
)

// dependencyHeading opens the board's dependency tables. Each table row names the
// dependent row in its first cell and a comma-separated list of the rows it depends on in
// its second cell.
const dependencyHeading = "## Dependencies"

// RowOwner reports whether path is the detail owner of the row id.
func RowOwner(id, path string) bool { return path == rowFilePath(id) }

// Close removes each named row heading from the index, with the blank line that follows
// it, removes the closed rows from the dependency tables, and projects the recommended
// sequence from the remaining outcomes. A row heading is what ParseDocument reads as one,
// and each named row must appear exactly once. The caller removes each detail owner.
func Close(index []byte, rows, outcomes []string) ([]byte, error) {
	lines := strings.Split(string(index), "\n")
	kept := make([]string, 0, len(lines))
	removed := map[string]int{}
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if m := roadmapStartRe.FindStringSubmatch(line); m != nil && slices.Contains(rows, m[1]) {
			removed[m[1]]++
			if i+1 < len(lines) && strings.TrimSpace(lines[i+1]) == "" {
				i++
			}
			continue
		}
		kept = append(kept, line)
	}
	for _, id := range rows {
		if removed[id] != 1 {
			return nil, fmt.Errorf("closed row %s appears %d times in %s", id, removed[id], RoadmapFile)
		}
	}
	return ProjectSequence([]byte(strings.Join(closeDependencies(kept, rows), "\n")), outcomes)
}

// closeDependencies removes the closed rows from every dependency table row. A table row
// whose dependent is closed goes. A closed row leaves each dependency list, and a table row
// whose list empties goes. Every other line stays byte for byte.
func closeDependencies(lines, closed []string) []string {
	start, end, _ := sectionBounds(lines, dependencyHeading)
	if start < 0 {
		return lines
	}
	kept := slices.Clone(lines[:start+1])
	for _, line := range lines[start+1 : end] {
		cells := strings.SplitN(line, "|", 4)
		if len(cells) < 4 || strings.TrimSpace(cells[0]) != "" {
			kept = append(kept, line)
			continue
		}
		if slices.Contains(closed, strings.TrimSpace(cells[1])) {
			continue
		}
		listed := strings.Split(cells[2], ",")
		var open []string
		for _, dependency := range listed {
			if dependency = strings.TrimSpace(dependency); !slices.Contains(closed, dependency) {
				open = append(open, dependency)
			}
		}
		switch {
		case len(open) == len(listed):
			kept = append(kept, line)
		case len(open) > 0:
			cells[2] = " " + strings.Join(open, ", ") + " "
			kept = append(kept, strings.Join(cells, "|"))
		}
	}
	return append(kept, lines[end:]...)
}
