# Review pickup: spec-ticket-fence-reduction

Base `73c97aa5`, reviewed tip `578811f2`. Transient: deleted by the green fix commit that closes these.

## Standards

Findings: 3. Worst: the provenance clause in `internal/coverage/coverage.go`.

- **auto-fix** — `internal/coverage/coverage.go` `projection()` comment narrates the deleted opt-in flag ("the same cells it yielded when the offsets came from an opt-in flag an unmatched header never set"); `craft-comments` forbids change narration. Cut the clause.
- **ask-user → accepted: trim** — `CONTEXT.md` **acceptance row** entry restates `craft-tickets`' slice/grading rule; both copies are anchored. Keep definition + Avoid list, drop the restated grading sentence; retarget the `**acceptance row**` needle and its `context-acceptance-row-vocabulary` canary to the trimmed text.
- no-op — the four-field schema is spelled in parser, craft-spec, CONTEXT.md, docs, CHANGELOG; pre-existing shape, anchors check prose only.

## Spec

Findings: 2, both no-op (the spec's own two-count `Verification log` and "7 slices to 8" are approval-time history, not false claims). Worst: none material.

## Coverage

Findings: 3. Worst: SR29 asserts docs enforcement that does not exist.

- **auto-fix** — `internal/coverage/coverage_test.go`: `TestCommand` case for a control-bearing behavior cell (`b\x1bx`) pinning exit 1 and the `unrepresentable TOON cell` error line — behavior newly reaches the TOON sink through the projection.
- **auto-fix** — `internal/coverage/coverage_test.go`: a rendered-projection case with a delimiter-bearing behavior (`a,b`, tab) asserting the quoted TOON row — the profile's "assert the permitted bytes" rule.
- **ask-user → accepted: anchor** — SR29: `Require` rows + canaries on the reduced-schema text in `docs/field-guide.html` (map card), `docs/reporesident-distillation.md`, and the CHANGELOG entry, mirroring the `context-coverage-row-vocabulary` fixture shape.

Repair targets after de-duplication: 4 (comment clause; two coverage_test cases; CONTEXT.md trim + needle; three doc anchors).
