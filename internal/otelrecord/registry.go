package otelrecord

// This file is the one source for the instrumented set. Each entry names one seam and
// the Go symbol that starts that seam's span. The `go/ast` conformance check enumerates
// Registry and reds when a named symbol starts no span, so a seam cannot stay silently
// uninstrumented. An instrumentation ticket adds its row here, and the check stays red
// until the symbol starts its span. The dispatch and seam declarations in the owning
// packages are reconciled against Registry by set-equality tests, so a drop on either
// side reds.

// SeamEntry binds a seam name to the symbol that opens its span. Package is the module-
// relative package directory, such as "internal/gate", so the check resolves the source
// from a subject root without a build. Function is the unexported or exported function
// name inside that package.
type SeamEntry struct {
	Seam     string
	Package  string
	Function string
}

// Registry is the instrumented set in seam order.
var Registry = []SeamEntry{
	{Seam: "gate", Package: "internal/gate", Function: "beginGateSpan"},
	{Seam: "gate.phase", Package: "internal/gate", Function: "startPhaseSpan"},
	{Seam: "lane", Package: "internal/gate", Function: "beginLaneSpan"},
	{Seam: "commit", Package: "internal/commit", Function: "beginCommitSpan"},
	{Seam: "worktree.land", Package: "internal/worktree", Function: "beginLandingSpan"},
	{Seam: "worktree.create", Package: "internal/worktree", Function: "beginVerbSpan"},
	{Seam: "worktree.exec", Package: "internal/worktree", Function: "beginVerbSpan"},
	{Seam: "worktree.merge", Package: "internal/worktree", Function: "beginVerbSpan"},
	{Seam: "worktree.release", Package: "internal/worktree", Function: "beginVerbSpan"},
	{Seam: "worktree.build", Package: "internal/worktree", Function: "beginVerbSpan"},
	{Seam: "worktree.reauthorize", Package: "internal/worktree", Function: "beginVerbSpan"},
	{Seam: "hook.guard-git", Package: "cmd/bench", Function: "beginHookSpan"},
	{Seam: "hook.guard-bench-follow-on", Package: "cmd/bench", Function: "beginHookSpan"},
	{Seam: "hook.check-agent-line", Package: "cmd/bench", Function: "beginHookSpan"},
	{Seam: "hook.stop-verdict", Package: "cmd/bench", Function: "beginHookSpan"},
	{Seam: "hook.session-inspect", Package: "cmd/bench", Function: "beginHookSpan"},
	{Seam: "hook.worktree-hook", Package: "cmd/bench", Function: "beginHookSpan"},
}
