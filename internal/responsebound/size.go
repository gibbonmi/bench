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

// AssignmentScope answers the assignment whose scope a spill of this process takes: the
// assignment of the worktree at root. It answers false for the primary scope, which a
// root outside any assignment worktree and a retiring verb both take. The census output
// record uses the same scope, so the two stores agree on which assignment a verb ran in.
func AssignmentScope(home, root string, retiring bool) (string, bool) {
	_, scope := location(home, root, retiring)
	return scope, scope != primaryScope
}
