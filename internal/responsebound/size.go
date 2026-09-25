package responsebound

// Size is the extent of one response's complete output: its line count, its byte count,
// and whether a spill file holds it.
type Size struct {
	Lines   int
	Bytes   int64
	Spilled bool
}

// Size answers the extent of the response after Finish. A spill file that took only part
// of the output still reads as spilled. A response whose spill file could not be created
// reads as inline, because it printed every byte.
func (o *Owner) Size() Size {
	o.mu.Lock()
	defer o.mu.Unlock()
	return Size{Lines: o.lines.count, Bytes: o.bytes, Spilled: o.path != ""}
}

// AssignmentScope answers the assignment of the worktree at root, the scope a spill of a
// verb that retires nothing takes. It answers false for a root outside any assignment
// worktree. The census output record uses this scope, so the two stores agree on which
// assignment a verb ran in.
func AssignmentScope(home, root string) (string, bool) {
	_, scope := location(home, root, false)
	return scope, scope != primaryScope
}
