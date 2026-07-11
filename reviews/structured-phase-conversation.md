# Structured Bench phase conversation review

## Standards

1 finding. Worst issue: duplicated canonical rule knowledge.

- **Hard violation — one source per fact.** `AGENTS.md` says an enforcement and
  its advertisement, and a parser and its count, must collapse to one source.
  The four shared clauses in `.bench/BENCH.md` are copied as full `needle`
  strings in `internal/conformance/docs_workflow_helpers_test.go`, while
  `internal/conformance/validity_checks_test.go` separately hardcodes the count
  as four. Derive the checked contract and its cardinality from one canonical
  representation.

## Spec

0 findings. All seven acceptance rows were audited, including both fresh-session
dogfood classes.

## Coverage

3 findings. Worst issue: a supported link state can omit the shared rules from
Claude phase sessions.

- **High — project-owned Claude bootstrap misses the contract.**
  `internal/contract/surface/link_test.go` requires `bench link` to preserve a
  project-owned `CLAUDE.md` without adding `@.bench/BENCH.md`. An explicit Bench
  phase in that supported state never loads the new conversation rule, contrary
  to the every-shipped-harness claim in
  `specs/structured-phase-conversation.md`. Resolving this requires a reviewer
  choice: change the safe-link posture or narrow the approved harness claim.
- **Medium — negated or misplaced anchors pass.** A present non-empty shared-rules
  file can retain each exact needle in an HTML comment, quotation, or negating
  sentence outside the communication section. The unscoped `strings.Contains`
  checks in `internal/conformance/docs_workflow_helpers_test.go` stay green even
  though active guidance disappeared; no acceptance row or canary covers that
  state.
- **Medium — the one-sentence/substantial boundary is unexercised.** The spec says
  a routine one-sentence acknowledgement stays untemplated while meaningful
  intermediate state plus continued work uses labels. Neither mapped dogfood
  class exercises a one-sentence update containing both state and continued
  action, so opposite interpretations remain possible.
