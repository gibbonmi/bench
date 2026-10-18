package staleness

import (
	"maps"
	"slices"
)

// union is every path of before and after in path order, with the object that each side
// holds. A path that one side lacks keeps the zero Object on that side.
func union(before, after map[string]Object) []PathState {
	paths := slices.Collect(maps.Keys(before))
	for path := range after {
		if _, ok := before[path]; !ok {
			paths = append(paths, path)
		}
	}
	slices.Sort(paths)
	states := make([]PathState, len(paths))
	for i, path := range paths {
		states[i] = PathState{Path: path, Before: before[path], After: after[path]}
	}
	return states
}

// drift is one entry for each changed path, in path order. A deleted member and an added
// member with one object ID are one rename that keeps both names, at the deleted path's
// place; the first unpaired added path in path order takes the pair.
func drift(paths []PathState) []DriftEntry {
	added := map[string][]int{}
	for i, state := range paths {
		if state.Before.ID == "" {
			added[state.After.ID] = append(added[state.After.ID], i)
		}
	}
	renamed, paired := map[int]int{}, map[int]bool{}
	for i, state := range paths {
		if state.After.ID != "" || len(added[state.Before.ID]) == 0 {
			continue
		}
		to := added[state.Before.ID][0]
		added[state.Before.ID] = added[state.Before.ID][1:]
		renamed[i], paired[to] = to, true
	}
	var entries []DriftEntry
	for i, state := range paths {
		switch to, isRename := renamed[i]; {
		case paired[i]:
		case isRename:
			entries = append(entries, entry(state.Path, state.Before, paths[to].Path, paths[to].After))
		case state.Before != state.After:
			entries = append(entries, entry(state.Path, state.Before, state.Path, state.After))
		}
	}
	return entries
}

// entry is the drift entry from one side to the other. Absent fills a side that holds no
// member. A mode difference between two present sides adds both modes to the kind, so a
// change of mode alone is never equal bytes.
func entry(beforePath string, before Object, afterPath string, after Object) DriftEntry {
	e := DriftEntry{Before: Absent, After: Absent, BeforeObject: Absent, AfterObject: Absent}
	if before.ID != "" {
		e.Before, e.BeforeObject = beforePath, before.ID
	}
	if after.ID != "" {
		e.After, e.AfterObject = afterPath, after.ID
	}
	switch {
	case before.ID == "":
		e.Kind = "added"
	case after.ID == "":
		e.Kind = "deleted"
	case beforePath != afterPath:
		e.Kind = "renamed"
	case before.ID != after.ID:
		e.Kind = "modified"
	}
	if before.ID != "" && after.ID != "" && before.Mode != after.Mode {
		mode := "mode " + before.Mode + "->" + after.Mode
		if e.Kind != "" {
			mode = e.Kind + " " + mode
		}
		e.Kind = mode
	}
	return e
}
