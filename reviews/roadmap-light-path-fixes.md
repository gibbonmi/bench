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

## Repair re-review

Repair range: `1f60a2bf6e2a2d45e1fae168be9744d64d86cd75` through
`fe7ef8051923250c3040c29265054cedae45cb18`.

- Standards: clean. The spec now explicitly serializes overlapping frontier
  `Writes:` entries, and the checkout compatibility mutation-red is durably
  recorded with production restored.
- Spec: clean. Repeating a retrospective slug refuses atomically before
  mutation; distinct slugs, grammar, and primary-local routing remain intact.
- Coverage: clean. The same-slug test proves both refusal and byte-for-byte
  preservation, while the earlier distinct-slug test remains green.

Final build and review preflights are 12/12 green. The final composition gate
passed gofmt, vet, test, race, and system phases with six unchanged capability
skips; shellcheck remained unavailable.

## Terra-high falsification review

Frozen base: `0b8096b18022bf3fc324e36e3afaae4c41151e74`.
Reviewed tip: `820ccf8763e86f0e7a950afda85a37a936334a96`.

### Accepted findings

- Standards P1: `bench retro` follows a live `capture/retros` directory
  symlink and can write outside the repository. A coordinator probe reproduced
  a successful escaped write.
- Standards P1: public help says `retro` replaces an artifact, while the
  accepted implementation refuses an existing slug.
- Standards P2: two new test fixtures duplicate the retrospective heading
  policy without one shared expectation and a durable biting-mutation record.
- Standards P2: `gate-prose` recovers structured prose fields by parsing the
  rendered diagnostic string already owned by `internal/prose`.
- Standards P2: three new anchor-test comments cite LF identifiers instead of
  stating only current-code constraints.
- Coverage P1: the LF2 hostile-input journey exercises only the zero-signal
  scaffold; removing the detected-project hygiene call leaves current tests
  green.
- Standards P2: the independent seeded-input expectation adds `BENCH_KIT` and
  `BENCH_RUN_BINARY` without a durable demonstrated removal red.

All seven findings route to auto-fix.

### Rejected finding

The Spec reviewer reported that the LF2 gate journey failed before reaching
the ignored-input check. The coordinator reran the focused test with the
retained kit and binary, without ambient routing variables, and with an
isolated HOME plus cached modules. All three subcases passed; the isolated run
took 17.65 seconds. The review worktree lacked an authenticated candidate
binary, so its `no packages` failure was review-environment-owned.
