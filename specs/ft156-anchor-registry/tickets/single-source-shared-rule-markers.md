# Single-source the shared-rule markers

Blocked by: restate-agents-shared-rule-strength.md
Ownership fence: `internal/anchors`, `internal/conformance/docs_workflow_helpers_test.go`, `internal/conformance/validity_checks_test.go`
Contracts: the two shared-rule marker values cross `internal/anchors`→the registry rows and `checkSharedRuleSingleSource`, asserted by AR4 against one exported owner and the residual-literal sweep

## What to build

Make `internal/anchors` the single owner of the fix-don't-park and source-warrant marker strings. Consume those values from both the registry rows and `checkSharedRuleSingleSource`; remove the stale conformance-local constants and the comment that falsely claims they are still shared with registry data.

## Acceptance

- [ ] [AR4] Each shared-rule marker has one production definition consumed by registry evaluation and the bespoke single-source check, with no residual literal copy in `internal/conformance`.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AR4 | restore either conformance-local marker literal | coordinator residual sweep | restore → enumerate definitions → reject duplicate → reinstate shared owner |
