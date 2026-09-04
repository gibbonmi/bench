// Package handoff implements the bench handoff command's core. It collects the pin
// facts from the git tree and rewrites the one section the caller's checkout owns.
// The document's grammar belongs to internal/handoffdoc, which parses, updates, and
// renders the file; this package supplies the section's content alone. The
// reviewer-owned State text passes through byte-for-byte. Every other line is
// derived.
package handoff

// ShapeSection is the body of the handoff's "## Shape" section, emitted below
// the heading the renderer writes itself. It is the single source of the
// handoff's Shape contract, exported so the conformance check can import it and
// compare it against the tracked capture/session-handoff.md. This proves the artifact
// is derived rather than a second source. A double-quoted literal because the
// text carries backticks.
const ShapeSection = "One section per live assignment, each rewritten in full by the phase that owns it\n" +
	"and pruned rather than accreted. A fresh session pays for every line it reads cold;\n" +
	"drop anything it would not act on.\n" +
	"\n" +
	"Operational gotchas are placed by lifetime, not copied here. One that recurs across\n" +
	"phases belongs in `projects/benchkit.md`'s cold-session notes. One scoped to a build\n" +
	"belongs instead in that spec's coverage rows.\n" +
	"\n" +
	"This file names at most when you'll hit one, never the command — a second copy\n" +
	"drifts from the source.\n" +
	"\n" +
	"A section opens at `## main` or at `## request <digest>`, and this section closes the\n" +
	"document. Below a section heading come its label lines: `Request token`, `Label`,\n" +
	"`Worktree tip`, `Recorded base`, and one `Spec` and `Spec status` pair per live spec.\n" +
	"One `Next command` line and a `### State` body follow them.\n" +
	"\n" +
	"A label line carries no sentence terminator, so the prose lane skips it. **State**\n" +
	"holds what is true now for that assignment, including anything uncommitted. **Next\n" +
	"command** holds the exact harness-native invocation, not a description of it.\n" +
	"\n" +
	"The handoff carries no date of its own. `bench status` dates each assignment section\n" +
	"by the branch commits past its recorded tip. It dates the `## main` section by the\n" +
	"file's last write. That write is the commit that carried the file, or the file's own\n" +
	"timestamp when git ignores it. Where this document and the tree disagree, the tree\n" +
	"wins.\n"

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
