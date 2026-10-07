# Read managed-block markers through the block reader

Blocked by: 1-add-block-reader.md
Writes: internal/adopt/marker.go, internal/adopt/link_plan_test.go
Covers: MB57, MB58, MB61, MB70, MB86, MB87, MB91, MB97

## What to build

The `cli-desktop-consistency` prerequisite is satisfied at `37ac80b8`.
Read `marker.go` and its current callers before the first edit. The
`stagedRepairAgents` caller and its link-transaction path are present. Keep that
repair path on the same marker scan and rewrite.

Make `scanMarkers`, `RewriteAgentsBlock`, and `StripAgentsBlock` read the block
reader's fence classes. Each one still finds a marker in the raw line text,
because each marker is itself an HTML comment. The unbalanced-fence flag comes
from the reader's fence fault or its comment fault. Either fault refuses only when
the file holds Bench text, as the fence fault does today at `marker.go:67`. Both
faults give the exact message at `marker.go:68`:
`conflict: AGENTS.md has an unclosed code fence around Bench markers; marker detection cannot be trusted`.

The rewrite and the strip write their output from the source bytes by the
reader's line offsets, so each carriage return survives. The marker scan ignores
the reader's frontmatter fault and reads each line as body, as it does today.

The new rows go in `link_plan_test.go`. `TestRewriteAgentsBlockEdges` lives in
`internal/adopt/adopt_test.go`, but that file is over the line cap.
`internal/adopt` is a crowded directory, so the ticket adds no file.

## Acceptance

- [ ] A marker example inside a `~~~` block stays, and `RewriteAgentsBlock` replaces only the live block.
- [ ] `StripAgentsBlock` keeps a tilde-fenced marker example and removes the live block.
- [ ] An unterminated `~~~` block around Bench text refuses with the `marker.go:68` message.
- [ ] A file that holds only an unterminated comment opener and a fenced marker example refuses with the same message, not the malformed-markers message.
- [ ] A file that holds an unterminated comment opener and no Bench text rewrites with no refusal.
- [ ] `RewriteAgentsBlock` on a CRLF file keeps each carriage return outside the replaced block.
- [ ] `StripAgentsBlock` on a CRLF file keeps each carriage return outside the removed block.
- [ ] An AGENTS.md whose first line is an unclosed `---` rewrites as it does today, with no refusal.
- [ ] `TestRewriteAgentsBlockEdges` stays green without an edit.
- [ ] `RewriteAgentsBlock` of the live `AGENTS.md` gives the same bytes at the base and at the tip.
- [ ] `internal/adopt/marker.go` holds no string literal with a run of three backticks.
- [ ] No file in `Writes:` grows past 400 lines.
