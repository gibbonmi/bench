# Add the handoff document leaf package

Blocked by: none
Writes: internal/handoffdoc (new), internal/handoff/text.go, internal/conformance/handoff_single_source_test.go
Covers: HS1, HS4, HS5

## What to build

Verify the premise first: `splitState` in internal/handoff/sections.go is the
only parser, it knows three headings, and no capture-file lock exists. Then
add `internal/handoffdoc`, a leaf package that imports the standard library
only. It parses and renders the document: a header
block, `## main`, one `## request <digest>` per section, and `## Shape` last.
A section holds label lines for the six pins, one `Next command:` line, and a
`### State` body. Expose parse, render, rewrite-one-section, remove-by-key, and
ensure-main.

Wrap each read-modify-write in an exclusive flock on a lock file
beside the document, in the `lockCleanupRegistration` shape from
internal/worktree/lifecycle.go, with temp-and-rename underneath. A lock held
past two seconds refuses with the lock path, the deadline the intent ledger
uses.

Update `ShapeSection` in internal/handoff/text.go to describe the section
grammar, and keep `checkHandoffShape` green.

## Acceptance

- [ ] A rendered two-section document parses back to the same sections, and `prose.Findings` reports nothing on it.
- [ ] Two goroutines that each rewrite a distinct section fifty times leave both sections present.
- [ ] `TestHandoffShapeSingleSourcedBites` passes with the new Shape text.
- [ ] Self-probe: drop the lock and keep temp-and-rename, and report the parallel test red.
