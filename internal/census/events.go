package census

import (
	"fmt"
	"github.com/gibbonmi/bench/internal/bounds"
	"path/filepath"
	"strings"
	"time"
)

// Event identifies a census row by assignment and its one-based record position.
type Event struct {
	ID   string
	Head string
}

// ReadEvents preserves malformed and unfinished coverage alongside valid native rows.
func ReadEvents(home, root, assignment string) ([]Event, []string, error) {
	if !isAssignmentID(assignment) {
		return nil, nil, fmt.Errorf("invalid census assignment")
	}
	path := filepath.Join(Dir(home, root), assignment)
	read := bounds.ClassifyNoFollow(path)
	if read.State != bounds.StateParsed {
		return nil, nil, fmt.Errorf("census %s: %s", read.State, read.Reason)
	}
	text := string(read.Data)
	lines := recordLines(text)
	var events []Event
	var problems []string
	for i, line := range lines {
		id := fmt.Sprintf("%s:%d", assignment, i+1)
		if i == len(lines)-1 && !strings.HasSuffix(text, "\n") {
			problems = append(problems, id+" unfinished")
			continue
		}
		stamp, head, ok := recordFields(line)
		_, stampErr := time.Parse(time.RFC3339, stamp)
		if !ok || stampErr != nil {
			problems = append(problems, id+" malformed")
			continue
		}
		events = append(events, Event{ID: id, Head: head})
	}
	return events, problems, nil
}
