# Read managed-block markers through the block reader

Blocked by: 1-add-block-reader.md
Writes: internal/adopt/marker.go, internal/adopt/link_plan_test.go
Covers: MB57, MB58, MB61, MB70

## What to build

Start this ticket only after `cli-desktop-consistency` lands on `main`, because
that branch also writes `internal/adopt/marker.go`. Re-read `marker.go` at that
tip before the first edit.

Make `scanMarkers`, `RewriteAgentsBlock`, and `StripAgentsBlock` read the block
reader's fence classes. Each one still finds a marker in the raw line text,
because each marker is itself an HTML comment. The unbalanced-fence flag comes
from the reader's fence fault. The reader's comment fault gives the same
conflict, because an unterminated comment hides each later fence marker. Both
faults keep the exact AGENTS.md conflict message bytes.

`internal/adopt` is a crowded directory, so the ticket adds no file.

## Acceptance

- [ ] A marker example inside a `~~~` block stays, and `RewriteAgentsBlock` replaces only the live block.
- [ ] `StripAgentsBlock` keeps a tilde-fenced marker example and removes the live block.
- [ ] An unterminated `~~~` block around Bench text refuses with the AGENTS.md conflict message.
- [ ] An unterminated HTML comment opener before a fenced marker example refuses with the same message.
- [ ] Each existing marker test in `internal/adopt/adopt_test.go` stays green without an edit.
