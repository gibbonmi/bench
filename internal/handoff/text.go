// Package handoff implements the bench handoff command's core: collecting the
// pin facts from the git tree, splitting an existing session-handoff.md into
// its sections, rendering the regenerated document, and writing it back. The
// reviewer-owned State section passes through byte-for-byte; every other
// section is derived.
package handoff

// ShapeSection is the body of the handoff's "## Shape" section, emitted below
// the heading the renderer writes itself. It is the single source of the
// handoff's Shape contract, exported so the conformance check can import it and
// compare it against the tracked session-handoff.md, proving the artifact is
// derived rather than a second source. A double-quoted literal because the
// text carries backticks.
const ShapeSection = "Rewritten in full at every phase close, pruned rather than accreted: a fresh\n" +
	"session pays for every line it reads cold, so drop anything it would not act on.\n" +
	"Operational gotchas are placed by lifetime, not copied here: one that recurs across\n" +
	"phases belongs in `projects/benchkit.md`'s cold-session notes, and one scoped to a\n" +
	"build belongs in that spec's coverage rows. This file names at most when you'll hit\n" +
	"one, never the command — a second copy drifts from the source.\n" +
	"Keep the three sections above — **State** (what is true now, including anything\n" +
	"uncommitted), **Next command** (the exact harness-native invocation, not a\n" +
	"description of it), and this one.\n" +
	"\n" +
	"The handoff carries no date of its own. `bench status` computes its age from the\n" +
	"commit that last wrote this file and reports a `handoff` row once anything has\n" +
	"landed since. Where this document and the tree disagree, the tree wins.\n"

// scaffoldGuidance is rendered between the header block and the "## State"
// heading only while State has no content, so it is what the first session in
// a repo reads and it vanishes once that session writes its State.
const scaffoldGuidance = `This document is a cold-start pin: a fresh session resumes from it with no
conversation history, so it has to carry everything that session needs to act.
The State section is yours — what is true right now and what it means, the one
judgment bench handoff cannot derive and will never overwrite. Everything else
here is regenerated from the git tree on every run of bench handoff, so
hand-editing those sections is wasted work. Write State, run bench handoff
again, and this guidance disappears.
`
