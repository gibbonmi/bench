package roadmap

import (
	"fmt"
	"slices"
	"strings"
)

// RowOwner reports whether path is the detail owner of the row id.
func RowOwner(id, path string) bool { return path == rowFilePath(id) }

// Close removes each named row heading from the index, with the blank line that follows
// it, and projects the recommended sequence from the remaining outcomes. A row heading is
// what ParseDocument reads as one, and each named row must appear exactly once. The
// caller removes each detail owner.
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
	return ProjectSequence([]byte(strings.Join(kept, "\n")), outcomes)
}
