# Roadmap light-path fixes review

Base: `0b8096b18022bf3fc324e36e3afaae4c41151e74`
Source: `999512c91e703feb4a79ff330ed4c9607bf74e67`

## Standards

### P1 — Ticket fences contradict the parallel-frontier claim

The spec says blockers serialize shared files and unrelated frontier tickets
remain eligible for parallel authoring. The post-build ticket-evidence closure
instead puts the same command-registry and conformance paths in multiple
unblocked tickets. For example,
`formalize-repair-charge-template.md` and
`make-bare-worktree-usage-safe.md` share five paths, while
`internal/conformance/subcommand_routing_test.go` appears in thirteen ticket
fences. The ticket discipline treats an overlapping `Writes:` note as a
conflict the coordinator must serialize. Repair the ticket metadata so the
declared frontier and its actual shared ownership agree.

### P2 — The checkout layout expectation lacks a durable mutation-red record

`internal/gate/prospectiveartifact/prospectiveartifact_test.go` independently
asserts the literal `"checkout"` while
`internal/gate/prospectiveartifact/prospectiveartifact.go` owns the same value
as `CheckoutName`. The repository permits that otherwise duplicated
expectation only when the named omission or mutation is recorded and
demonstrated red. Preserve the compatibility assertion only with a durable
mutation-red record in the spec evidence.

## Spec

### P1 — Repeating a retrospective slug overwrites earlier capture

LF12 requires repeated writes to preserve earlier capture. `RetroCommand`
currently calls `os.WriteFile` for an existing slug, and its repeated-write test
uses two different slugs. A second write to the same slug therefore destroys
the first body. Refuse an existing target before mutation and add a same-slug
preservation test.

No other LF1–LF28 variance was found.

## Coverage

The Coverage axis mapped every LF1–LF28 row to executable tests or the four
review-owned comment enumerations and reported no additional finding. The six
gate capability skips cover host-specific FIFO, socket, and privilege branches;
none bypasses a changed acceptance seam.

Coordinator correlation: the accepted LF12 Spec finding also identifies the
missing same-slug negative test. Standards and Spec findings remain separate;
the Coverage axis added no independent defect.
