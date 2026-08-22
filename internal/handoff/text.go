// Package handoff implements the bench handoff command's core. It collects the
// pin facts from the git tree, splits an existing capture/session-handoff.md into
// its sections, renders the regenerated document, and writes it back. The
// reviewer-owned State section passes through byte-for-byte. Every other
// section is derived.
package handoff

// ShapeSection is the body of the handoff's "## Shape" section, emitted below
// the heading the renderer writes itself. It is the single source of the
// handoff's Shape contract, exported so the conformance check can import it and
// compare it against the tracked capture/session-handoff.md. This proves the artifact
// is derived rather than a second source. A double-quoted literal because the
// text carries backticks.
const ShapeSection = "Rewritten in full at every phase close, pruned rather than accreted. A fresh\n" +
	"session pays for every line it reads cold; drop anything it would not act on.\n" +
	"\n" +
	"Operational gotchas are placed by lifetime, not copied here. One that recurs across\n" +
	"phases belongs in `projects/benchkit.md`'s cold-session notes. One scoped to a build\n" +
	"belongs instead in that spec's coverage rows.\n" +
	"\n" +
	"This file names at most when you'll hit one, never the command — a second copy\n" +
	"drifts from the source.\n" +
	"\n" +
	"Keep the three sections above. **State** holds what is true now, including anything\n" +
	"uncommitted. **Next command** holds the exact harness-native invocation, not a\n" +
	"description of it. This section is the third.\n" +
	"\n" +
	"The handoff carries no date of its own. `bench status` computes its age from the\n" +
	"commit that last wrote this file and reports a `handoff` row once anything has\n" +
	"landed since. Where this document and the tree disagree, the tree wins.\n"

// scaffoldGuidance is rendered between the header block and the "## State"
// heading only while State has no content. It is what the first session in
// a repo reads, and it vanishes once that session writes its State.
const scaffoldGuidance = `This document is a cold-start pin. A fresh session resumes from it with no
conversation history, so it has to carry everything that session needs to act.
The State section is yours — what is true right now and what it means, the one
judgment bench handoff cannot derive and will never overwrite. Everything else
here is regenerated from the git tree on every run of bench handoff, so
hand-editing those sections is wasted work. Write State, run bench handoff
again, and this guidance disappears.
`
