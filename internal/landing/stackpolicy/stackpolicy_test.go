package stackpolicy

import (
	"reflect"
	"testing"
)

// Each metadata string is one opaque "<mode> <object>" literal; the rule compares the
// whole string and never parses it.
const (
	blobA = "100644 1111111111111111111111111111111111111111"
	blobB = "100644 2222222222222222222222222222222222222222"
	blobC = "100644 3333333333333333333333333333333333333333"
)

func assertOwners(t *testing.T, got, want map[string][]int) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("owners = %v, want %v", got, want)
	}
}

// TestStackOwnersNamesEachChangingMember proves that a path two members change names
// both members, and a path one member changes names that member alone. (Coverage row
// LS1.)
func TestStackOwnersNamesEachChangingMember(t *testing.T) {
	destination := Entries{"shared.go": blobA, "first.go": blobA, "same.go": blobA}
	members := []Entries{
		{"shared.go": blobB, "first.go": blobB, "same.go": blobA},
		{"shared.go": blobC, "first.go": blobA, "same.go": blobA},
	}
	assertOwners(t, Owners(destination, members), map[string][]int{
		"shared.go": {1, 2},
		"first.go":  {1},
	})
}

// TestStackOverlapsLeaveOutCapture proves that a capture path two members change is
// not an overlap. (Coverage row LS2.)
func TestStackOverlapsLeaveOutCapture(t *testing.T) {
	members := []Entries{
		{"capture/learnings.md": blobB},
		{"capture/learnings.md": blobC},
	}
	if got := Overlaps(Owners(Entries{"capture/learnings.md": blobA}, members)); len(got) != 0 {
		t.Fatalf("overlaps = %v, want none", got)
	}
}

// TestStackOverlapsKeepAnOutsidePath proves that a path outside capture/ that two
// members change is an overlap, and a path one member changes is not. (Coverage row
// LS3.)
func TestStackOverlapsKeepAnOutsidePath(t *testing.T) {
	members := []Entries{
		{"internal/x.go": blobB, "alone.go": blobB},
		{"internal/x.go": blobC},
	}
	assertOwners(t, Overlaps(Owners(Entries{}, members)), map[string][]int{
		"internal/x.go": {1, 2},
	})
}

// TestStackOwnersCountsAModeChange proves that a member that changes only the mode of
// a path owns that path. (Coverage row LS4.)
func TestStackOwnersCountsAModeChange(t *testing.T) {
	destination := Entries{"run.sh": blobA}
	members := []Entries{{"run.sh": "100755 1111111111111111111111111111111111111111"}}
	assertOwners(t, Owners(destination, members), map[string][]int{"run.sh": {1}})
}

// TestStackOwnersCountsBothRenamePaths proves that a member that renames a path owns
// the old path and the new path. (Coverage row LS5.)
func TestStackOwnersCountsBothRenamePaths(t *testing.T) {
	members := []Entries{{"new.go": blobA}}
	assertOwners(t, Owners(Entries{"old.go": blobA}, members), map[string][]int{
		"old.go": {1},
		"new.go": {1},
	})
}

// TestStackOwnersKeepsExactPaths proves that a changed path that holds a space and a
// `*` keeps its exact spelling. (Coverage row LS6.)
func TestStackOwnersKeepsExactPaths(t *testing.T) {
	const path = "dir/a file*.go"
	destination := Entries{path: blobA, "dir/a": blobA, "file*.go": blobA}
	members := []Entries{{path: blobB, "dir/a": blobA, "file*.go": blobA}}
	assertOwners(t, Owners(destination, members), map[string][]int{path: {1}})
}
