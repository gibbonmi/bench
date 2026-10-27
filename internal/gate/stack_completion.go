package gate

import (
	"context"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/gibbonmi/bench/internal/landing/stackpolicy"
)

// StackMember is one member of a landing stack: its source tip and its tickets folder. A
// spec-less member has no folder.
type StackMember struct{ Source, Folder string }

type stackKey struct{}

// WithStack binds the members of a landing stack, in stack order, to the landing value of
// ctx. The gate then proves the graded tree as the composition of those members on the
// landing destination, and the landing runs the whole gate.
func WithStack(ctx context.Context, members []StackMember) context.Context {
	return context.WithValue(ctx, stackKey{}, slices.Clone(members))
}

// stackGrade is the stack context of an evaluation: the landing destination and the members.
type stackGrade struct {
	destination string
	members     []StackMember
}

// stackGradeOf answers the stack context of ctx, or nil when ctx carries none.
func stackGradeOf(ctx context.Context) *stackGrade {
	members, stacked := ctx.Value(stackKey{}).([]StackMember)
	if !stacked {
		return nil
	}
	landing, _ := landingOf(ctx)
	return &stackGrade{destination: landing.destination, members: members}
}

// stackProof proves that graded is exactly the composition of the stack members on the
// destination. Each member's tickets folder is absent, and each member with a folder carries
// its verified closure. Every other path keeps the destination's entry when no member changes
// it, and the entry of its one owner otherwise. A path that two or more members change passes
// only when the settle policy settles it. Only a failed proof is a proof fault.
func (e *gateEvaluation) stackProof(graded *treeGeneration) error {
	_, destination, err := capturedCommit(e.identityRoot, e.stack.destination)
	if err != nil {
		return err
	}
	members := make([]*treeGeneration, len(e.stack.members))
	entries := make([]stackpolicy.Entries, len(e.stack.members))
	proven := map[string]bool{}
	for i, member := range e.stack.members {
		var tree string
		tree, members[i], err = capturedCommit(e.identityRoot, member.Source)
		if err != nil {
			return err
		}
		entries[i] = members[i].snapshot.byPath
		// A spec-less member has no deliverable, and the closure of an empty path fails.
		if member.Folder == "" {
			continue
		}
		if err := closedFolder(members[i], graded, member.Folder, proven); err != nil {
			return err
		}
		closed, err := e.completionClosure(graded, tree, member.Folder, member.Source)
		if err != nil {
			return err
		}
		maps.Copy(proven, closed)
	}
	owners := stackpolicy.Owners(destination.snapshot.byPath, entries)
	overlaps := stackpolicy.Overlaps(owners)
	paths := slices.Concat(slices.Collect(maps.Keys(owners)), slices.Collect(maps.Keys(destination.snapshot.byPath)), slices.Collect(maps.Keys(graded.snapshot.byPath)))
	slices.Sort(paths)
	for _, path := range slices.Compact(paths) {
		if proven[path] || slices.ContainsFunc(e.stack.members, func(member StackMember) bool { return member.Folder != "" && underFolder(path, member.Folder) }) {
			continue
		}
		indexes := owners[path]
		switch len(indexes) {
		case 0:
			if !sameEntry(graded, destination, path) {
				return proofFault("stack composition changes %s; review the destination delta", path)
			}
		case 1:
			if !sameEntry(graded, members[indexes[0]-1], path) {
				return proofFault("stack member %d differs at %s; review the stack composition", indexes[0], path)
			}
		default:
			if overlap, found := overlaps[path]; found {
				return proofFault("stack members %s change %s; land them apart", memberList(overlap), path)
			}
		}
	}
	return nil
}

// sameEntry reports whether path has one entry, or none, in both generations.
func sameEntry(a, b *treeGeneration, path string) bool {
	left, inA := a.entry(path)
	right, inB := b.entry(path)
	return inA == inB && left.Metadata == right.Metadata
}

// memberList spells member indexes as "1 and 2", or "1, 2 and 3".
func memberList(indexes []int) string {
	names := make([]string, len(indexes))
	for i, index := range indexes {
		names[i] = strconv.Itoa(index)
	}
	last := len(names) - 1
	return strings.Join(names[:last], ", ") + " and " + names[last]
}
