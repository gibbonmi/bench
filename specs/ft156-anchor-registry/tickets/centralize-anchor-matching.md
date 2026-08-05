# Centralize anchor matching

Blocked by: add-workflow-section-relocation-tripwire.md, pin-ticket-contracts-template.md
Ownership fence: `internal/anchors`, `internal/conformance/docs_workflow_helpers_test.go`
Contracts: normalized anchor text and section-resolution results cross `internal/anchors`→the local conformance anchor declarations, asserted by AM1 and AM2 against the real helper callers and recorded parser proofs

## What to build

Move the four-kind matching semantics and Markdown H2 parser into a non-test `internal/anchors` package below the conformance import edge. Keep the anchor declarations local for this ticket, but route their evaluation through the shared package so no matcher or parser copy remains.

## Acceptance

- [ ] [AM1] Whole-file and section-scoped require/forbid declarations retain their current collapse, case-fold, section-resolution, and diagnostic behavior through the shared matcher.
- [ ] [AM2] Fenced and unclosed-fence headings retain the recorded Markdown H2 parser behavior after the parser moves.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AM1 | bypass whitespace collapse for one require kind | the shared matcher unit tests plus real anchor helpers | run the focused matcher tests and root conformance anchor check, expect the hostile-spacing case to fail |
| AM2 | treat a fenced H2 as a real section | the moved recorded-proof parser test | run the focused parser test and expect the fenced-heading assertion to fail |
